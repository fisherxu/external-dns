package imm

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

// DeleteTagSet invokes the imm.DeleteTagSet API synchronously
// api document: https://help.aliyun.com/api/imm/deletetagset.html
func (client *Client) DeleteTagSet(request *DeleteTagSetRequest) (response *DeleteTagSetResponse, err error) {
	response = CreateDeleteTagSetResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteTagSetWithChan invokes the imm.DeleteTagSet API asynchronously
// api document: https://help.aliyun.com/api/imm/deletetagset.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteTagSetWithChan(request *DeleteTagSetRequest) (<-chan *DeleteTagSetResponse, <-chan error) {
	responseChan := make(chan *DeleteTagSetResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteTagSet(request)
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

// DeleteTagSetWithCallback invokes the imm.DeleteTagSet API asynchronously
// api document: https://help.aliyun.com/api/imm/deletetagset.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteTagSetWithCallback(request *DeleteTagSetRequest, callback func(response *DeleteTagSetResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteTagSetResponse
		var err error
		defer close(result)
		response, err = client.DeleteTagSet(request)
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

// DeleteTagSetRequest is the request struct for api DeleteTagSet
type DeleteTagSetRequest struct {
	*requests.RpcRequest
	LazyMode   string `position:"Query" name:"LazyMode"`
	Project    string `position:"Query" name:"Project"`
	SetId      string `position:"Query" name:"SetId"`
	CheckEmpty string `position:"Query" name:"CheckEmpty"`
}

// DeleteTagSetResponse is the response struct for api DeleteTagSet
type DeleteTagSetResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteTagSetRequest creates a request to invoke DeleteTagSet API
func CreateDeleteTagSetRequest() (request *DeleteTagSetRequest) {
	request = &DeleteTagSetRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("imm", "2017-09-06", "DeleteTagSet", "imm", "openAPI")
	return
}

// CreateDeleteTagSetResponse creates a response to parse from DeleteTagSet response
func CreateDeleteTagSetResponse() (response *DeleteTagSetResponse) {
	response = &DeleteTagSetResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
