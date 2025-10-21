package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/xanderxampp-be/franco/contextwrap"
	"github.com/xanderxampp-be/franco/log"
	"github.com/xanderxampp-be/franco/trace"

	"go.elastic.co/apm/module/apmhttp"
)

var client *http.Client

func Init() {
	client = &http.Client{
		Timeout: 20 * time.Second,
	}

	client = apmhttp.WrapClient(client)
}

func InitWithParam(c *http.Client) {
	client = apmhttp.WrapClient(c)
}

func Call(ctx context.Context, requestBody map[string]interface{}, header http.Header, endpoint string) (context.Context, []byte, http.Header, error) {
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

	request, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		return ctx, nil, nil, err
	}

	request.Header = header

	response, err := client.Do(request.WithContext(ctx))
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
		Url:     endpoint,
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
