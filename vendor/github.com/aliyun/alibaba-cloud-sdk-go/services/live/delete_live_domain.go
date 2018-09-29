package live

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

// DeleteLiveDomain invokes the live.DeleteLiveDomain API synchronously
// api document: https://help.aliyun.com/api/live/deletelivedomain.html
func (client *Client) DeleteLiveDomain(request *DeleteLiveDomainRequest) (response *DeleteLiveDomainResponse, err error) {
	response = CreateDeleteLiveDomainResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteLiveDomainWithChan invokes the live.DeleteLiveDomain API asynchronously
// api document: https://help.aliyun.com/api/live/deletelivedomain.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteLiveDomainWithChan(request *DeleteLiveDomainRequest) (<-chan *DeleteLiveDomainResponse, <-chan error) {
	responseChan := make(chan *DeleteLiveDomainResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteLiveDomain(request)
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

// DeleteLiveDomainWithCallback invokes the live.DeleteLiveDomain API asynchronously
// api document: https://help.aliyun.com/api/live/deletelivedomain.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteLiveDomainWithCallback(request *DeleteLiveDomainRequest, callback func(response *DeleteLiveDomainResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteLiveDomainResponse
		var err error
		defer close(result)
		response, err = client.DeleteLiveDomain(request)
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

// DeleteLiveDomainRequest is the request struct for api DeleteLiveDomain
type DeleteLiveDomainRequest struct {
	*requests.RpcRequest
	SecurityToken string           `position:"Query" name:"SecurityToken"`
	OwnerAccount  string           `position:"Query" name:"OwnerAccount"`
	DomainName    string           `position:"Query" name:"DomainName"`
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
}

// DeleteLiveDomainResponse is the response struct for api DeleteLiveDomain
type DeleteLiveDomainResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteLiveDomainRequest creates a request to invoke DeleteLiveDomain API
func CreateDeleteLiveDomainRequest() (request *DeleteLiveDomainRequest) {
	request = &DeleteLiveDomainRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("live", "2016-11-01", "DeleteLiveDomain", "live", "openAPI")
	return
}

// CreateDeleteLiveDomainResponse creates a response to parse from DeleteLiveDomain response
func CreateDeleteLiveDomainResponse() (response *DeleteLiveDomainResponse) {
	response = &DeleteLiveDomainResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
