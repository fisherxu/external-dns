package mediaservices

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator 1.0.1.0
// Changes may cause incorrect behavior and will be lost if the code is
// regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"net/http"
)

// Client is the media Services resource management APIs.
type Client struct {
	ManagementClient
}

// NewClient creates an instance of the Client client.
func NewClient(subscriptionID string) Client {
	return NewClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewClientWithBaseURI creates an instance of the Client client.
func NewClientWithBaseURI(baseURI string, subscriptionID string) Client {
	return Client{NewWithBaseURI(baseURI, subscriptionID)}
}

// CheckNameAvailability checks whether the Media Service resource name is
// available. The name must be globally unique.
//
// checkNameAvailabilityInput is properties needed to check the availability of
// a name.
func (client Client) CheckNameAvailability(checkNameAvailabilityInput CheckNameAvailabilityInput) (result CheckNameAvailabilityOutput, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: checkNameAvailabilityInput,
			Constraints: []validation.Constraint{{Target: "checkNameAvailabilityInput.Name", Name: validation.Null, Rule: true,
				Chain: []validation.Constraint{{Target: "checkNameAvailabilityInput.Name", Name: validation.MaxLength, Rule: 24, Chain: nil},
					{Target: "checkNameAvailabilityInput.Name", Name: validation.MinLength, Rule: 3, Chain: nil},
					{Target: "checkNameAvailabilityInput.Name", Name: validation.Pattern, Rule: `^[a-z0-9]`, Chain: nil},
				}},
				{Target: "checkNameAvailabilityInput.Type", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "mediaservices.Client", "CheckNameAvailability")
	}

	req, err := client.CheckNameAvailabilityPreparer(checkNameAvailabilityInput)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "CheckNameAvailability", nil, "Failure preparing request")
		return
	}

	resp, err := client.CheckNameAvailabilitySender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "CheckNameAvailability", resp, "Failure sending request")
		return
	}

	result, err = client.CheckNameAvailabilityResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "CheckNameAvailability", resp, "Failure responding to request")
	}

	return
}

// CheckNameAvailabilityPreparer prepares the CheckNameAvailability request.
func (client Client) CheckNameAvailabilityPreparer(checkNameAvailabilityInput CheckNameAvailabilityInput) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-10-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Media/CheckNameAvailability", pathParameters),
		autorest.WithJSON(checkNameAvailabilityInput),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// CheckNameAvailabilitySender sends the CheckNameAvailability request. The method will close the
// http.Response Body if it receives an error.
func (client Client) CheckNameAvailabilitySender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// CheckNameAvailabilityResponder handles the response to the CheckNameAvailability request. The method always
// closes the http.Response Body.
func (client Client) CheckNameAvailabilityResponder(resp *http.Response) (result CheckNameAvailabilityOutput, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Create creates a Media Service.
//
// resourceGroupName is name of the resource group within the Azure
// subscription. mediaServiceName is name of the Media Service. mediaService is
// media Service properties needed for creation.
func (client Client) Create(resourceGroupName string, mediaServiceName string, mediaService MediaService) (result MediaService, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: mediaServiceName,
			Constraints: []validation.Constraint{{Target: "mediaServiceName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "mediaServiceName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "mediaServiceName", Name: validation.Pattern, Rule: `^[a-z0-9]`, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "mediaservices.Client", "Create")
	}

	req, err := client.CreatePreparer(resourceGroupName, mediaServiceName, mediaService)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "Create", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "Create", resp, "Failure sending request")
		return
	}

	result, err = client.CreateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "Create", resp, "Failure responding to request")
	}

	return
}

// CreatePreparer prepares the Create request.
func (client Client) CreatePreparer(resourceGroupName string, mediaServiceName string, mediaService MediaService) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"mediaServiceName":  autorest.Encode("path", mediaServiceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-10-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaservices/{mediaServiceName}", pathParameters),
		autorest.WithJSON(mediaService),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// CreateSender sends the Create request. The method will close the
// http.Response Body if it receives an error.
func (client Client) CreateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// CreateResponder handles the response to the Create request. The method always
// closes the http.Response Body.
func (client Client) CreateResponder(resp *http.Response) (result MediaService, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete deletes a Media Service.
//
// resourceGroupName is name of the resource group within the Azure
// subscription. mediaServiceName is name of the Media Service.
func (client Client) Delete(resourceGroupName string, mediaServiceName string) (result autorest.Response, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: mediaServiceName,
			Constraints: []validation.Constraint{{Target: "mediaServiceName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "mediaServiceName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "mediaServiceName", Name: validation.Pattern, Rule: `^[a-z0-9]`, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "mediaservices.Client", "Delete")
	}

	req, err := client.DeletePreparer(resourceGroupName, mediaServiceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client Client) DeletePreparer(resourceGroupName string, mediaServiceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"mediaServiceName":  autorest.Encode("path", mediaServiceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-10-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaservices/{mediaServiceName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client Client) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client Client) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets a Media Service.
//
// resourceGroupName is name of the resource group within the Azure
// subscription. mediaServiceName is name of the Media Service.
func (client Client) Get(resourceGroupName string, mediaServiceName string) (result MediaService, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: mediaServiceName,
			Constraints: []validation.Constraint{{Target: "mediaServiceName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "mediaServiceName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "mediaServiceName", Name: validation.Pattern, Rule: `^[a-z0-9]`, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "mediaservices.Client", "Get")
	}

	req, err := client.GetPreparer(resourceGroupName, mediaServiceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client Client) GetPreparer(resourceGroupName string, mediaServiceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"mediaServiceName":  autorest.Encode("path", mediaServiceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-10-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaservices/{mediaServiceName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client Client) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client Client) GetResponder(resp *http.Response) (result MediaService, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByResourceGroup lists all of the Media Services in a resource group.
//
// resourceGroupName is name of the resource group within the Azure
// subscription.
func (client Client) ListByResourceGroup(resourceGroupName string) (result Collection, err error) {
	req, err := client.ListByResourceGroupPreparer(resourceGroupName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "ListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByResourceGroupSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "ListByResourceGroup", resp, "Failure sending request")
		return
	}

	result, err = client.ListByResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "ListByResourceGroup", resp, "Failure responding to request")
	}

	return
}

// ListByResourceGroupPreparer prepares the ListByResourceGroup request.
func (client Client) ListByResourceGroupPreparer(resourceGroupName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-10-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaservices", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListByResourceGroupSender sends the ListByResourceGroup request. The method will close the
// http.Response Body if it receives an error.
func (client Client) ListByResourceGroupSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListByResourceGroupResponder handles the response to the ListByResourceGroup request. The method always
// closes the http.Response Body.
func (client Client) ListByResourceGroupResponder(resp *http.Response) (result Collection, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListKeys lists the keys for a Media Service.
//
// resourceGroupName is name of the resource group within the Azure
// subscription. mediaServiceName is name of the Media Service.
func (client Client) ListKeys(resourceGroupName string, mediaServiceName string) (result ServiceKeys, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: mediaServiceName,
			Constraints: []validation.Constraint{{Target: "mediaServiceName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "mediaServiceName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "mediaServiceName", Name: validation.Pattern, Rule: `^[a-z0-9]`, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "mediaservices.Client", "ListKeys")
	}

	req, err := client.ListKeysPreparer(resourceGroupName, mediaServiceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "ListKeys", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListKeysSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "ListKeys", resp, "Failure sending request")
		return
	}

	result, err = client.ListKeysResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "ListKeys", resp, "Failure responding to request")
	}

	return
}

// ListKeysPreparer prepares the ListKeys request.
func (client Client) ListKeysPreparer(resourceGroupName string, mediaServiceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"mediaServiceName":  autorest.Encode("path", mediaServiceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-10-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaservices/{mediaServiceName}/listKeys", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListKeysSender sends the ListKeys request. The method will close the
// http.Response Body if it receives an error.
func (client Client) ListKeysSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListKeysResponder handles the response to the ListKeys request. The method always
// closes the http.Response Body.
func (client Client) ListKeysResponder(resp *http.Response) (result ServiceKeys, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// RegenerateKey regenerates a primary or secondary key for a Media Service.
//
// resourceGroupName is name of the resource group within the Azure
// subscription. mediaServiceName is name of the Media Service.
// regenerateKeyInput is properties needed to regenerate the Media Service key.
func (client Client) RegenerateKey(resourceGroupName string, mediaServiceName string, regenerateKeyInput RegenerateKeyInput) (result RegenerateKeyOutput, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: mediaServiceName,
			Constraints: []validation.Constraint{{Target: "mediaServiceName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "mediaServiceName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "mediaServiceName", Name: validation.Pattern, Rule: `^[a-z0-9]`, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "mediaservices.Client", "RegenerateKey")
	}

	req, err := client.RegenerateKeyPreparer(resourceGroupName, mediaServiceName, regenerateKeyInput)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "RegenerateKey", nil, "Failure preparing request")
		return
	}

	resp, err := client.RegenerateKeySender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "RegenerateKey", resp, "Failure sending request")
		return
	}

	result, err = client.RegenerateKeyResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "RegenerateKey", resp, "Failure responding to request")
	}

	return
}

// RegenerateKeyPreparer prepares the RegenerateKey request.
func (client Client) RegenerateKeyPreparer(resourceGroupName string, mediaServiceName string, regenerateKeyInput RegenerateKeyInput) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"mediaServiceName":  autorest.Encode("path", mediaServiceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-10-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaservices/{mediaServiceName}/regenerateKey", pathParameters),
		autorest.WithJSON(regenerateKeyInput),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// RegenerateKeySender sends the RegenerateKey request. The method will close the
// http.Response Body if it receives an error.
func (client Client) RegenerateKeySender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// RegenerateKeyResponder handles the response to the RegenerateKey request. The method always
// closes the http.Response Body.
func (client Client) RegenerateKeyResponder(resp *http.Response) (result RegenerateKeyOutput, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// SyncStorageKeys synchronizes storage account keys for a storage account
// associated with the Media Service account.
//
// resourceGroupName is name of the resource group within the Azure
// subscription. mediaServiceName is name of the Media Service.
// syncStorageKeysInput is properties needed to synchronize the keys for a
// storage account to the Media Service.
func (client Client) SyncStorageKeys(resourceGroupName string, mediaServiceName string, syncStorageKeysInput SyncStorageKeysInput) (result autorest.Response, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: mediaServiceName,
			Constraints: []validation.Constraint{{Target: "mediaServiceName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "mediaServiceName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "mediaServiceName", Name: validation.Pattern, Rule: `^[a-z0-9]`, Chain: nil}}},
		{TargetValue: syncStorageKeysInput,
			Constraints: []validation.Constraint{{Target: "syncStorageKeysInput.ID", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "mediaservices.Client", "SyncStorageKeys")
	}

	req, err := client.SyncStorageKeysPreparer(resourceGroupName, mediaServiceName, syncStorageKeysInput)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "SyncStorageKeys", nil, "Failure preparing request")
		return
	}

	resp, err := client.SyncStorageKeysSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "SyncStorageKeys", resp, "Failure sending request")
		return
	}

	result, err = client.SyncStorageKeysResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "SyncStorageKeys", resp, "Failure responding to request")
	}

	return
}

// SyncStorageKeysPreparer prepares the SyncStorageKeys request.
func (client Client) SyncStorageKeysPreparer(resourceGroupName string, mediaServiceName string, syncStorageKeysInput SyncStorageKeysInput) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"mediaServiceName":  autorest.Encode("path", mediaServiceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-10-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaservices/{mediaServiceName}/syncStorageKeys", pathParameters),
		autorest.WithJSON(syncStorageKeysInput),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// SyncStorageKeysSender sends the SyncStorageKeys request. The method will close the
// http.Response Body if it receives an error.
func (client Client) SyncStorageKeysSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// SyncStorageKeysResponder handles the response to the SyncStorageKeys request. The method always
// closes the http.Response Body.
func (client Client) SyncStorageKeysResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Update updates a Media Service.
//
// resourceGroupName is name of the resource group within the Azure
// subscription. mediaServiceName is name of the Media Service. mediaService is
// media Service properties needed for update.
func (client Client) Update(resourceGroupName string, mediaServiceName string, mediaService MediaService) (result MediaService, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: mediaServiceName,
			Constraints: []validation.Constraint{{Target: "mediaServiceName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "mediaServiceName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "mediaServiceName", Name: validation.Pattern, Rule: `^[a-z0-9]`, Chain: nil}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "mediaservices.Client", "Update")
	}

	req, err := client.UpdatePreparer(resourceGroupName, mediaServiceName, mediaService)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "Update", nil, "Failure preparing request")
		return
	}

	resp, err := client.UpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "Update", resp, "Failure sending request")
		return
	}

	result, err = client.UpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mediaservices.Client", "Update", resp, "Failure responding to request")
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client Client) UpdatePreparer(resourceGroupName string, mediaServiceName string, mediaService MediaService) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"mediaServiceName":  autorest.Encode("path", mediaServiceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-10-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaservices/{mediaServiceName}", pathParameters),
		autorest.WithJSON(mediaService),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client Client) UpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client Client) UpdateResponder(resp *http.Response) (result MediaService, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
