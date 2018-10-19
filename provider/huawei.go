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
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/recordsets"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/zones"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"k8s.io/apimachinery/pkg/util/wait"

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
	IAMEndpoint string `json:"iamEndpoint" yaml:"iamEndpoint"`
	DNSEndpoint string `json:"dnsEndpoint" yaml:"dnsEndpoint"`
	ProjectID   string `json:"projectId" yaml:"projectId"`
	Username    string `json:"username" yaml:"username"`
	Password    string `json:"password" yaml:"password"`
	Domainname  string `json:"domainname" yaml:"domainname"`
	ZoneID      string `json:"zoneId" yaml:"zoneId"`
	Token       string `json:"-" yaml:"-"`
	Expired     bool   `json:"-" yaml:"-"`
}

// NewAlibabaCloudProvider creates a new Alibaba Cloud provider.
func NewHuaweiCloudProvider(configFile string, domainFilter DomainFilter, zoneIDFileter ZoneIDFilter) (*HuaweiCloudProvider, error) {
	cfg := &huaweiCloudConfig{}
	contents, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to read Alibaba Cloud config file '%s': %v", configFile, err)
	}
	err = yaml.Unmarshal(contents, &cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Alibaba Cloud config file '%s': %v", configFile, err)
	}

	fmt.Println("cfgggggggggggggg", cfg)
	providerClient, err := getProviderClientFromIAM(cfg.Username, cfg.Password, cfg.Domainname, cfg.ProjectID, cfg.IAMEndpoint)
	if err != nil {
		return nil, err
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

	fmt.Println(cfg)

	go wait.Forever(provider.updateProviderClient, 2*time.Hour)
	return provider, nil
}

func (p *HuaweiCloudProvider) Records() ([]*endpoint.Endpoint, error) {
	fmt.Println("start recordssssssssssssss")
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
	managedZones, err := p.zones()
	if err != nil {
		return err
	}

	fmt.Println("start ApplyChangesssssssssssssssssss")
	fmt.Println(len(changes.Create) + len(changes.Delete) + len(changes.UpdateNew))
	fmt.Println(changes.UpdateOld)
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
	for _, rs := range recordSets {
		if err2 := p.upsertRecordSet(rs, managedZones); err == nil {
			err = err2
		}
	}
	return err
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

func (p *HuaweiCloudProvider) createRecord(zoneID string, opts recordsets.CreateOpts) error {
	_, err := recordsets.Create(p.client, zoneID, opts).Extract()
	if err != nil {
		return err
	}

	return nil
}

func (p *HuaweiCloudProvider) deleteRecord(zoneID, recordSetID string) error {
	err := recordsets.Delete(p.client, zoneID, recordSetID).ExtractErr()
	if err != nil {
		return err
	}

	return nil
}

func (p *HuaweiCloudProvider) updateRecord(zoneID, recordSetID string, opts recordsets.UpdateOpts) error {
	_, err := recordsets.Update(p.client, zoneID, recordSetID, opts).Extract()
	if err != nil {
		return err
	}

	return nil
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
	fmt.Println("zonezzzzzzzzzzzzzz")
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
		panic(err)
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

func (p *HuaweiCloudProvider) updateProviderClient() {
	providerClient, err := getProviderClientFromIAM(p.config.Username, p.config.Password, p.config.Domainname, p.config.ProjectID, p.config.IAMEndpoint)
	if err != nil {
		glog.Errorf("Get Token from Huawei Cloud failed: %v", err)
	}
	p.client.ProviderClient = providerClient
}

func getProviderClientFromIAM(username, password, domainname, projectID, iamEndpoint string) (*gophercloud.ProviderClient, error) {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: iamEndpoint,
		Username:         username,
		Password:         password,
		DomainName:       domainname,
		TenantID:         projectID,
	}
	glog.Infof("AuthOptions: %v", opts)
	return openstack.AuthenticatedClient(opts)
}
