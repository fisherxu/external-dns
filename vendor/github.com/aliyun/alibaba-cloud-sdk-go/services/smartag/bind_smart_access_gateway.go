package smartag

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// BindSmartAccessGateway invokes the smartag.BindSmartAccessGateway API synchronously
// api document: https://help.aliyun.com/api/smartag/bindsmartaccessgateway.html
func (client *Client) BindSmartAccessGateway(request *BindSmartAccessGatewayRequest) (response *BindSmartAccessGatewayResponse, err error) {
	response = CreateBindSmartAccessGatewayResponse()
	err = client.DoAction(request, response)
	return
}

// BindSmartAccessGatewayWithChan invokes the smartag.BindSmartAccessGateway API asynchronously
// api document: https://help.aliyun.com/api/smartag/bindsmartaccessgateway.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) BindSmartAccessGatewayWithChan(request *BindSmartAccessGatewayRequest) (<-chan *BindSmartAccessGatewayResponse, <-chan error) {
	responseChan := make(chan *BindSmartAccessGatewayResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.BindSmartAccessGateway(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// BindSmartAccessGatewayWithCallback invokes the smartag.BindSmartAccessGateway API asynchronously
// api document: https://help.aliyun.com/api/smartag/bindsmartaccessgateway.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) BindSmartAccessGatewayWithCallback(request *BindSmartAccessGatewayRequest, callback func(response *BindSmartAccessGatewayResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *BindSmartAccessGatewayResponse
		var err error
		defer close(result)
		response, err = client.BindSmartAccessGateway(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// BindSmartAccessGatewayRequest is the request struct for api BindSmartAccessGateway
type BindSmartAccessGatewayRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	CcnId                string           `position:"Query" name:"CcnId"`
	SmartAGId            string           `position:"Query" name:"SmartAGId"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}

// BindSmartAccessGatewayResponse is the response struct for api BindSmartAccessGateway
type BindSmartAccessGatewayResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateBindSmartAccessGatewayRequest creates a request to invoke BindSmartAccessGateway API
func CreateBindSmartAccessGatewayRequest() (request *BindSmartAccessGatewayRequest) {
	request = &BindSmartAccessGatewayRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Smartag", "2018-03-13", "BindSmartAccessGateway", "smartag", "openAPI")
	return
}

// CreateBindSmartAccessGatewayResponse creates a response to parse from BindSmartAccessGateway response
func CreateBindSmartAccessGatewayResponse() (response *BindSmartAccessGatewayResponse) {
	response = &BindSmartAccessGatewayResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
