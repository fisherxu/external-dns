/*
Copyright 2018 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/recordsets"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/zones"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/kubernetes-incubator/external-dns/endpoint"
	"github.com/kubernetes-incubator/external-dns/plan"
)

const (
	// ID of the RecordSet from which endpoint was created
	RecordSetID = "huawei-recordset-id"
	// Zone ID of the RecordSet
	ZoneID = "huawei-zone-id"

	// Initial records values of the RecordSet. This label is required in order not to loose records that haven't
	// changed where there are several targets per domain and only some of them changed.
	// Values are joined by zero-byte to in order to get a single string
	OriginalRecords = "huawei-original-records"
)

// HuaweiCloudProvider implements the DNS provider for Huawei Cloud.
type HuaweiCloudProvider struct {
	domainFilter DomainFilter
	zoneIDFilter ZoneIDFilter
	client       *gophercloud.ServiceClient
	config       *huaweiCloudConfig
	dryRun       bool
}

type huaweiCloudConfig struct {
	AuthType    string `json:"authType" yaml:"authType"`
	DNSEndpoint string `json:"dnsEndpoint" yaml:"dnsEndpoint"`
	ProjectID   string `json:"projectId" yaml:"projectId"`

	// token config
	IAMEndpoint string `json:"iamEndpoint" yaml:"iamEndpoint"`
	Username    string `json:"username" yaml:"username"`
	Password    string `json:"password" yaml:"password"`
	Domainname  string `json:"domainname" yaml:"domainname"`

	// aksk config
	Region        string `json:"region" yaml:"region"`
	AccessKeyPath string `json:"accessKeyPath" yaml:"accessKeyPath"`
	SecretKeyPath string `json:"secretKeyPath" yaml:"secretKeyPath"`
	ServiceType   string `json:"serviceType" yaml:"serviceType"`

	ZoneID  string `json:"zoneId" yaml:"zoneId"`
	Token   string `json:"-" yaml:"-"`
	Expired bool   `json:"-" yaml:"-"`
}

var defaultHuaweiCloudConfig = &huaweiCloudConfig{
	AuthType:      "aksk",
	DNSEndpoint:   "https://dns.myhuaweicloud.com/v2",
	IAMEndpoint:   "https://iam.myhuaweicloud.com/v3/auth/tokens/",
	Region:        "cn-north-1",
	ServiceType:   "all",
	AccessKeyPath: "/etc/hwcloud/accesskey",
	SecretKeyPath: "/etc/hwcloud/secretkey",
}

var fileChanged = false

var mutex sync.Mutex

// NewHuaweiCloudProvider creates a new Huawei Cloud provider.
func NewHuaweiCloudProvider(configFile string, domainFilter DomainFilter, zoneIDFileter ZoneIDFilter) (*HuaweiCloudProvider, error) {
	cfg := defaultHuaweiCloudConfig
	contents, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to read Huawei Cloud config file '%s': %v", configFile, err)
	}
	err = yaml.Unmarshal(contents, &cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Huawei Cloud config file '%s': %v", configFile, err)
	}

	var providerClient *gophercloud.ProviderClient
	if cfg.AuthType == "token" {
		providerClient, err = getProviderClientFromIAM(cfg.Username, cfg.Password, cfg.Domainname, cfg.ProjectID, cfg.IAMEndpoint)
		if err != nil {
			return nil, err
		}
	} else {
		providerClient, err = getProviderClientFromAksk(cfg.Region, cfg.AccessKeyPath, cfg.SecretKeyPath, cfg.ServiceType, cfg.ProjectID)
		if err != nil {
			return nil, err
		}
	}

	var dnsClient *gophercloud.ServiceClient
	dnsClient = &gophercloud.ServiceClient{
		ProviderClient: providerClient,
		Endpoint:       strings.TrimSuffix(cfg.DNSEndpoint, "/") + "/",
	}

	provider := &HuaweiCloudProvider{
		client:       dnsClient,
		config:       cfg,
		domainFilter: domainFilter,
		zoneIDFilter: zoneIDFileter,
	}

	err = Newfilewatcher(cfg.AccessKeyPath, cfg.SecretKeyPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to watch aksk files, err: %v", err)
	}
	return provider, nil
}

func (p *HuaweiCloudProvider) Records() ([]*endpoint.Endpoint, error) {
	mutex.Lock()
	if fileChanged {
		err := p.updateProviderClient()
		if err != nil {
			mutex.Unlock()
			return nil, err
		}
		fileChanged = false
	}
	mutex.Unlock()

	records, err := p.records()
	if err != nil {
		return nil, err
	}
	var endpoints []*endpoint.Endpoint

	for _, record := range records {
		ep := endpoint.NewEndpointWithTTL(record.Name, record.Type, endpoint.TTL(record.TTL), record.Records...)
		ep.ZoneID = record.ZoneID
		ep.RecordsetID = record.ID
		endpoints = append(endpoints, ep)
	}
	return endpoints, nil
}

func (p *HuaweiCloudProvider) ApplyChanges(changes *plan.Changes) error {
	mutex.Lock()
	if fileChanged {
		err := p.updateProviderClient()
		if err != nil {
			mutex.Unlock()
			return err
		}
		fileChanged = false
	}
	mutex.Unlock()

	managedZones, err := p.zones()
	if err != nil {
		return err
	}

	recordSets := map[string]*recordSet{}
	for _, ep := range changes.Create {
		addEndpoint(ep, recordSets, false)
	}
	for _, ep := range changes.UpdateNew {
		addEndpoint(ep, recordSets, false)
	}
	for _, ep := range changes.UpdateOld {
		addEndpoint(ep, recordSets, true)
	}
	for _, ep := range changes.Delete {
		addEndpoint(ep, recordSets, true)
	}

	failedrecordSets := []string{}
	for _, rs := range recordSets {
		if err = p.upsertRecordSet(rs, managedZones); err != nil {
			log.Errorf("apply recordset changes failed, recordset: %v, err: %v", rs.dnsName, err)
			failedrecordSets = append(failedrecordSets, rs.dnsName)
		}
	}
	if len(failedrecordSets) > 0 {
		return fmt.Errorf("Failed to submit all changes for the following recordsets: %v", failedrecordSets)
	}
	return nil
}

// apply recordset changes by inserting/updating/deleting recordsets
func (p HuaweiCloudProvider) upsertRecordSet(rs *recordSet, managedZones map[string]string) error {
	if rs.zoneID == "" {
		var err error
		rs.zoneID, err = p.getEndpointZoneID(rs.dnsName, managedZones)
		if err != nil {
			return err
		}
		if rs.zoneID == "" {
			log.Debugf("Skipping record %s because no hosted zone matching record DNS Name was detected ", rs.dnsName)
			return nil
		}
	}
	var records []string
	for rec, v := range rs.names {
		if v {
			records = append(records, rec)
		}
	}
	if rs.recordSetID == "" && records == nil {
		return nil
	}
	if rs.recordSetID == "" {
		opts := recordsets.CreateOpts{
			Name:    rs.dnsName,
			Type:    rs.recordType,
			TTL:     rs.ttl,
			Records: records,
		}
		log.Infof("Creating records: %s/%s: %s", rs.dnsName, rs.recordType, strings.Join(records, ","))
		if p.dryRun {
			return nil
		}
		err := p.createRecord(rs.zoneID, opts)
		return err
	} else if len(records) == 0 {
		log.Infof("Deleting records for %s/%s", rs.dnsName, rs.recordType)
		if p.dryRun {
			return nil
		}
		return p.deleteRecord(rs.zoneID, rs.recordSetID)
	} else {
		opts := recordsets.UpdateOpts{
			Records: records,
			TTL:     rs.ttl,
		}
		log.Infof("Updating records: %s/%s: %s", rs.dnsName, rs.recordType, strings.Join(records, ","))
		if p.dryRun {
			return nil
		}
		return p.updateRecord(rs.zoneID, rs.recordSetID, opts)
	}
}

func (p *HuaweiCloudProvider) createRecord(zoneID string, opts recordsets.CreateOpts) (err error) {
	_, err = recordsets.Create(p.client, zoneID, opts).Extract()
	return
}

func (p *HuaweiCloudProvider) deleteRecord(zoneID, recordSetID string) error {
	return recordsets.Delete(p.client, zoneID, recordSetID).ExtractErr()
}

func (p *HuaweiCloudProvider) updateRecord(zoneID, recordSetID string, opts recordsets.UpdateOpts) (err error) {
	_, err = recordsets.Update(p.client, zoneID, recordSetID, opts).Extract()
	return
}

func (p *HuaweiCloudProvider) getEndpointZoneID(hostname string, managedZones map[string]string) (string, error) {
	longestZoneLength := 0
	resultID := ""

	for zoneID, zoneName := range managedZones {
		if !strings.HasSuffix(hostname, zoneName) {
			continue
		}
		ln := len(zoneName)
		if ln > longestZoneLength {
			resultID = zoneID
			longestZoneLength = ln
		}
	}

	return resultID, nil
}

func (p *HuaweiCloudProvider) getEndpointRecordID(endpoint *endpoint.Endpoint) (string, string, error) {
	records, err := p.records()
	if err != nil {
		return "", "", err
	}
	for _, record := range records {
		if strings.TrimSuffix(record.Name, ".") == strings.TrimSuffix(endpoint.DNSName, ".") && record.Type != "TXT" {
			return record.ZoneID, record.ID, nil
		}
	}
	return "", "", fmt.Errorf("don't find the zoneID/recordID of endpoint: %s", endpoint.DNSName)
}

func (p *HuaweiCloudProvider) records() ([]recordsets.RecordSet, error) {
	manageZones, err := p.zones()
	if err != nil {
		return nil, err
	}

	var recordset []recordsets.RecordSet
	for zoneID := range manageZones {
		if !p.zoneIDFilter.Match(zoneID) {
			continue
		}
		rr, err := p.recordsets(zoneID)
		if err != nil {
			log.Errorf("HuaweiCloudProvider RecordSets %s error %v", manageZones[zoneID], err)
			continue
		}
		recordset = append(recordset, rr...)
	}
	return recordset, nil
}

func (p *HuaweiCloudProvider) zones() (map[string]string, error) {
	allPages, err := zones.List(p.client, nil).AllPages()
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	allZones, err := zones.ExtractZones(allPages)
	if err != nil {
		return nil, err
	}
	for _, zone := range allZones {
		zoneName := canonicalizeDomainName(zone.Name)
		if !p.domainFilter.Match(zoneName) {
			continue
		}
		result[zone.ID] = zoneName
	}
	return result, nil
}

func (p *HuaweiCloudProvider) recordsets(zoneID string) ([]recordsets.RecordSet, error) {
	allPages, err := recordsets.ListByZone(p.client, zoneID, nil).AllPages()
	if err != nil {
		panic(err)
	}

	allRRs, err := recordsets.ExtractRecordSets(allPages)
	if err != nil {
		return nil, err
	}
	return allRRs, nil
}

func (p *HuaweiCloudProvider) updateProviderClient() error {
	var providerClient *gophercloud.ProviderClient
	var err error
	if p.config.AuthType == "token" {
		providerClient, err = getProviderClientFromIAM(p.config.Username, p.config.Password, p.config.Domainname, p.config.ProjectID, p.config.IAMEndpoint)
		if err != nil {
			log.Errorf("Get Token from Huawei Cloud failed: %v", err)
			return err
		}
	} else {
		providerClient, err = getProviderClientFromAksk(p.config.Region, p.config.AccessKeyPath, p.config.SecretKeyPath, p.config.ServiceType, p.config.ProjectID)
		if err != nil {
			log.Errorf("Update aksk client failed: %v", err)
			return err
		}
	}

	p.client.ProviderClient = providerClient
	return nil
}

func getProviderClientFromIAM(username, password, domainname, projectID, iamEndpoint string) (*gophercloud.ProviderClient, error) {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: iamEndpoint,
		Username:         username,
		Password:         password,
		DomainName:       domainname,
		TenantID:         projectID,
	}
	log.Infof("AuthOptions: %v", opts)
	return openstack.AuthenticatedClient(opts)
}

func getProviderClientFromAksk(region, accessKeyPath, secretKeyPath, serviceType, id string) (*gophercloud.ProviderClient, error) {
	accessKey, err := ioutil.ReadFile(accessKeyPath)
	if err != nil {
		return nil, err
	}
	secretKey, err := ioutil.ReadFile(secretKeyPath)
	if err != nil {
		return nil, err
	}

	access := &gophercloud.AccessInfo{
		AccessKey:   strings.TrimSpace(string(accessKey)),
		SecretKey:   strings.TrimSpace(string(secretKey)),
		Region:      region,
		ServiceType: serviceType,
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*15)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 15,
		},
	}

	akskClient := &gophercloud.AkskClient{
		Client:   httpClient,
		Access:   access,
		TenantId: id,
	}
	return &gophercloud.ProviderClient{AkskClient: akskClient}, nil
}

func Newfilewatcher(accessKeyPath, secretKeyPath string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case <-watcher.Events:
				mutex.Lock()
				fileChanged = true
				mutex.Unlock()
			case err = <-watcher.Errors:
				log.Errorf("Received an error from file watcher: %s", err)
			}
		}
	}()

	err = watcher.Add(accessKeyPath)
	if err != nil {
		return err
	}
	err = watcher.Add(secretKeyPath)
	if err != nil {
		return err
	}
	return nil
}
