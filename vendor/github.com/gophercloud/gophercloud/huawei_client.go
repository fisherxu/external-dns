package gophercloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang/glog"

	"k8s.io/client-go/util/flowcontrol"
)

const (
	// longThrottleLatency defines threshold for logging requests. All requests being
	// throttle for more than longThrottleLatency will be logged.
	longThrottleLatency = 50 * time.Millisecond
)

type AkskClient struct {
	Client   *http.Client
	Access   *AccessInfo
	TenantId string
}

type AccessInfo struct {
	Region      string
	AccessKey   string
	SecretKey   string
	ServiceType string
}

// request is used to help build up a request
type request struct {
	method  string
	url     string
	params  url.Values
	body    io.ReadCloser
	headers map[string]string
}

// Add by HUAWEI provider
// newRequest is used to create a new request
// if accessIn == nil mean not to sign header
func (client *AkskClient) NewRequest(method, url string, headersIn map[string]string, body io.Reader) *request {
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	r := &request{
		method:  method,
		url:     url,
		params:  make(map[string][]string),
		headers: headersIn,
		body:    rc,
	}
	return r
}

// decodeBody is used to JSON decode a body
func (client *AkskClient) DecodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("request failed: %s, status code: %d", string(resBody), resp.StatusCode)
	}

	if len(strings.Replace(string(resBody), " ", "", -1)) <= 2 {
		return nil
	}

	dec := json.NewDecoder(bytes.NewReader(resBody))

	err = dec.Decode(out)
	if err != nil {
		return fmt.Errorf("Decode failed: %s, err: %v", string(resBody), err)
	}

	return nil
}

// encodeBody is used to encode a request body
func encodeBody(obj interface{}) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return nil, fmt.Errorf("encode obj error")
	}
	return buf, nil
}

// doRequest runs a request with our client
func (client *AkskClient) DoRequest(throttle flowcontrol.RateLimiter, r *request) (*http.Response, error) {
	tryThrottle(throttle, r)

	// Create the HTTP request
	req, err := http.NewRequest(r.method, r.url, r.body)
	if err != nil {
		return nil, fmt.Errorf("http new request error")
	}
	req.Close = true

	// add the sign to request header if needed.
	if client.Access != nil {
		sign := Signer{
			AccessKey: client.Access.AccessKey,
			SecretKey: client.Access.SecretKey,
			Region:    client.Access.Region,
			Service:   client.Access.ServiceType,
		}
		req.Header.Set(HeaderProject, client.TenantId)
		if err := sign.Sign(req); err != nil {
			return nil, fmt.Errorf("DoRequest failed to get sign key %v", err)
		}
	}

	resp, err := client.Client.Do(req)
	if err != nil {
		return resp, fmt.Errorf("http client do request error. %v", err)
	}

	return resp, nil
}

func tryThrottle(throttle flowcontrol.RateLimiter, r *request) {
	now := time.Now()
	if throttle != nil {
		throttle.Accept()
	}
	if latency := time.Since(now); latency > longThrottleLatency {
		glog.V(2).Infof("Throttling request took %v, request: %s:%s", latency, r.method, r.url)
	}
}
