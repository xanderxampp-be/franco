package httpclient

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/xanderxampp-be/franco/contextwrap"
	"github.com/xanderxampp-be/franco/log"
	"github.com/xanderxampp-be/franco/log/entity"
	"github.com/xanderxampp-be/franco/trace"

	"go.elastic.co/apm/module/apmhttp"
)

type HttpMicro interface {
	Call(ctx context.Context, requestBody map[string]interface{}, header http.Header, path string) (context.Context, []byte, http.Header, error)
	GenerateHeaderLivvik(body, clientID, clientKey string) http.Header
	GenerateHeaderGenMicro(device *entity.Device, ipSource, agent string) http.Header
}

type HttpMicroImpl struct {
	httpc   *http.Client
	baseUrl string
}

func NewMicro(timeout int, baseUrl string) HttpMicro {
	hclientWConfig := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	client := apmhttp.WrapClient(hclientWConfig)

	custom := &HttpMicroImpl{
		httpc:   client,
		baseUrl: baseUrl,
	}

	return custom
}

func (c *HttpMicroImpl) Call(ctx context.Context, requestBody map[string]interface{}, header http.Header, path string) (context.Context, []byte, http.Header, error) {
	start := time.Now()
	jsonRequest, _ := json.Marshal(requestBody)

	var payload *bytes.Reader

	if _, ok := header["X-CLIENT-ID"]; ok {
		var param = url.Values{}
		param.Set("request", string(jsonRequest))
		payload = bytes.NewReader([]byte(param.Encode()))
	} else {
		payload = bytes.NewReader(jsonRequest)
	}

	currentTrace := contextwrap.GetTraceFromContext(ctx)

	request, err := http.NewRequest("POST", c.baseUrl+path, payload)
	if err != nil {
		return ctx, nil, nil, err
	}

	request.Header = header

	response, err := c.httpc.Do(request.WithContext(ctx))
	if err != nil {
		return ctx, nil, nil, err
	}

	defer response.Body.Close()

	responseByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		//log.LogWarn(err.Error(), "read esb response")
		return ctx, nil, nil, err
	}

	elapsed := time.Since(start).String()

	tr := &trace.TraceHttp{
		Url:     c.baseUrl + path,
		Request: log.Minify(requestBody),
		Elapsed: elapsed,
	}

	currentTrace = append(currentTrace, tr)

	ctx = context.WithValue(ctx, contextwrap.TraceKey, currentTrace)

	// check for valid json response
	var js map[string]interface{}
	err = json.Unmarshal(responseByte, &js)
	if err != nil {
		//log.LogWarn("invalid json", "invalid json")
		return ctx, nil, nil, err
	}

	tr.Response = log.Minify(js)

	return ctx, responseByte, response.Header, nil
}

func (c *HttpMicroImpl) GenerateHeaderLivvik(body, clientID, clientKey string) http.Header {
	contentType := "application/x-www-form-urlencoded"
	randomKey := strconv.Itoa(1000000000000000 + rand.Intn(9000000000000000))
	xApiKeyBytes := sha256.Sum256([]byte(body + clientKey + randomKey))
	xApiKeyString := hex.EncodeToString(xApiKeyBytes[:])

	headers := http.Header{
		"Content-Type": []string{contentType},
		"X-CLIENT-ID":  []string{clientID},
		"X-RANDOM-KEY": []string{randomKey},
		"X-API-KEY":    []string{xApiKeyString[0:16]},
	}

	return headers
}

func (c *HttpMicroImpl) GenerateHeaderGenMicro(device *entity.Device, ipSource, agent string) http.Header {
	contentType := "application/json"

	headers := http.Header{
		"Content-Type": []string{contentType},
		"User-Agent":   []string{agent},
		"IP-Address":   []string{ipSource},
		"Device-Id":    []string{device.DeviceID},
		"Device-Type":  []string{device.DeviceType},
		"Device-Name":  []string{device.DeviceVersion},
	}

	return headers
}
