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

// HuaweiCloudProvider implements the DNS provider for Huawei Cloud.
type HuaweiCloudProvider struct {
	domainFilter DomainFilter
	zoneIDFilter ZoneIDFilter
	client       *gophercloud.ServiceClient
	config       *huaweiCloudConfig
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
		ep := endpoint.NewEndpointWithZone(record.Name, record.Type, endpoint.TTL(record.TTL), record.ZoneID, record.ID, record.Records...)
		endpoints = append(endpoints, ep)
		fmt.Println(*ep)
	}
	return endpoints, nil
}

func (p *HuaweiCloudProvider) ApplyChanges(changes *plan.Changes) error {
	fmt.Println("start ApplyChangesssssssssssssssssss")
	fmt.Println(len(changes.Create) + len(changes.Delete) + len(changes.UpdateNew))
	if changes == nil || len(changes.Create)+len(changes.Delete)+len(changes.UpdateNew) == 0 {
		// No op
		return nil
	}
	p.createRecords(changes.Create)
	p.deleteRecords(changes.Delete)
	p.updateRecords(changes.UpdateNew)
	return nil
}

func (p *HuaweiCloudProvider) createRecords(endpoints []*endpoint.Endpoint) error {
	for _, endpoint := range endpoints {
		err := p.createRecord(endpoint)
		if err != nil {
			log.Errorf("HuaweiCloudProvider CreateRecordSet %s error %v", endpoint.DNSName, err)
		}
	}
	return nil
}

func (p *HuaweiCloudProvider) createRecord(endpoint *endpoint.Endpoint) error {
	createOpts := recordsets.CreateOpts{
		Name:    endpoint.DNSName,
		Type:    endpoint.RecordType,
		TTL:     int(endpoint.RecordTTL),
		Records: endpoint.Targets,
	}

	if endpoint.ZoneID == "" {
		zoneID, err := p.getEndpointZoneID(endpoint)
		if err != nil {
			return err
		}
		endpoint.ZoneID = zoneID
	}
	_, err := recordsets.Create(p.client, endpoint.ZoneID, createOpts).Extract()
	if err != nil {
		return err
	}

	return nil
}

func (p *HuaweiCloudProvider) deleteRecords(endpoints []*endpoint.Endpoint) error {
	for _, endpoint := range endpoints {
		p.deleteRecord(endpoint)
	}

	return nil
}

func (p *HuaweiCloudProvider) deleteRecord(endpoint *endpoint.Endpoint) error {
	err := recordsets.Delete(p.client, endpoint.ZoneID, endpoint.RecordsetID).ExtractErr()
	if err != nil {
		return err
	}

	return nil
}

func (p *HuaweiCloudProvider) updateRecords(endpoints []*endpoint.Endpoint) error {
	for _, endpoint := range endpoints {
		err := p.updateRecord(endpoint)
		if err != nil {
			log.Errorf("HuaweiCloudProvider UpdateRecordSet %s error %v", endpoint.DNSName, err)
		}
	}

	return nil
}

func (p *HuaweiCloudProvider) updateRecord(endpoint *endpoint.Endpoint) error {
	updateOpts := recordsets.UpdateOpts{
		TTL:     int(endpoint.RecordTTL),
		Records: endpoint.Targets,
	}
	if endpoint.ZoneID == "" || endpoint.RecordsetID == "" {
		zoneID, recordID, err := p.getEndpointRecordID(endpoint)
		if err != nil {
			return err
		}
		endpoint.ZoneID = zoneID
		endpoint.RecordsetID = recordID
	}
	_, err := recordsets.Update(p.client, endpoint.ZoneID, endpoint.RecordsetID, updateOpts).Extract()
	if err != nil {
		return err
	}

	return nil
}

func (p *HuaweiCloudProvider) getEndpointZoneID(endpoint *endpoint.Endpoint) (string, error) {
	zones, err := p.zones()
	if err != nil {
		return "", err
	}
	for _, zone := range zones {
		if strings.Contains(strings.TrimSuffix(endpoint.DNSName, "."), strings.TrimSuffix(zone.Name, ".")) {
			return zone.ID, nil
		}
	}
	return "", fmt.Errorf("don't find the zoneID of endpoint: %s", endpoint.DNSName)
}

func (p *HuaweiCloudProvider) getEndpointRecordID(endpoint *endpoint.Endpoint) (string, string, error) {
	records, err := p.records()
	if err != nil {
		return "", "", err
	}
	for _, record := range records {
		if record.Name == endpoint.DNSName {
			return record.ZoneID, record.ID, nil
		}
	}
	return "", "", fmt.Errorf("don't find the zoneID/recordID of endpoint: %s", endpoint.DNSName)
}

func (p *HuaweiCloudProvider) records() ([]recordsets.RecordSet, error) {
	zones, err := p.zones()
	if err != nil {
		return nil, err
	}
	fmt.Println("zonezzzzzzzzzzzzzz")
	var recordset []recordsets.RecordSet
	for _, zone := range zones {
		if zone.Name == "" {
			continue
		}
		if !p.domainFilter.Match(zone.Name) {
			continue
		}
		if !p.zoneIDFilter.Match(zone.ID) {
			continue
		}
		rr, err := p.recordsets(zone.ID)
		if err != nil {
			log.Errorf("HuaweiCloudProvider RecordSets %s error %v", zone.Name, err)
			continue
		}
		recordset = append(recordset, rr...)
	}
	return recordset, nil
}

func (p *HuaweiCloudProvider) zones() ([]zones.Zone, error) {
	allPages, err := zones.List(p.client, nil).AllPages()
	if err != nil {
		panic(err)
	}

	allZones, err := zones.ExtractZones(allPages)
	if err != nil {
		return nil, err
	}
	return allZones, nil
}

func (p *HuaweiCloudProvider) recordsets(zoneID string) ([]recordsets.RecordSet, error) {
	allPages, err := recordsets.ListByZone(p.client, zoneID, nil).AllPages()
	if err != nil {
		panic(err)
	}

	allRRs, err := recordsets.ExtractRecordSets(allPages)
	if err != nil {
		panic(err)
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
