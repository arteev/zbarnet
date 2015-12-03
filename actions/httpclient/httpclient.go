package httpclient

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/arteev/zbarnet/barcode"
	"github.com/arteev/zbarnet/config"
	"github.com/arteev/zbarnet/logger"
)

//A HTTPClient client for HTTP request with API key
type HTTPClient struct {
	method       config.HTTPMethod
	url          string
	apikey       string
	apiKeyHeader bool
}

//New HttpClient for request with a barcode
func New(method config.HTTPMethod, url, apikey string, apiKeyHeader bool) *HTTPClient {
	result := &HTTPClient{
		//req: http.new
		method:       method,
		url:          url,
		apikey:       apikey,
		apiKeyHeader: apiKeyHeader,
	}
	return result
}

//buildUrl replace the values in the template
func (h *HTTPClient) buildURL(bc *barcode.BarCode) (result string) {
	logger.Trace.Println("buildUrl start")
	defer logger.Trace.Println("buildUrl done")
	result = h.url
	if h.apikey != "" && !h.apiKeyHeader {
		result = strings.Replace(result, "${apikey}", h.apikey, -1)
	}
	result = strings.Replace(result, "${barCodeType}", url.QueryEscape(bc.Type), -1)
	result = strings.Replace(result, "${quality}", strconv.Itoa(bc.Quality), -1)
	if h.method == config.HTTPGET {
		result = strings.Replace(result, "${barCodeRaw}", url.QueryEscape(string(bc.Data)), -1)
		result = strings.Replace(result, "${barCode}", url.QueryEscape(base64.StdEncoding.EncodeToString(bc.Data)), -1)
	}
	logger.Debug.Printf("buildUrl url = %q\n", result)
	return result
}

//Send barcode on server
func (h *HTTPClient) Send(bc *barcode.BarCode) error {
	logger.Trace.Println("Send start")
	defer logger.Trace.Println("Send done")
	url := h.buildURL(bc)
	databc := new(bytes.Buffer)
	if h.method == config.HTTPPOST {
		js, err := bc.ToJSON()
		if err != nil {
			return err
		}
		databc = bytes.NewBufferString(js)
	}
	logger.Debug.Printf("Send method:%q, url:%q, data:%v\n", string(h.method), url, databc)
	req, err := http.NewRequest(string(h.method), url, databc)
	if err != nil {
		return err
	}
	if h.apiKeyHeader {
		req.Header.Set("Authorization", h.apikey)
	}
	client := &http.Client{}
	logger.Trace.Println("Send request")
	res, err := client.Do(req)
	logger.Trace.Println("Send request done")
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Warn.Println(err)
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			logger.Warn.Println(err)
		}
	}()
	logger.Debug.Printf("Send response  status: %s data: %s\n", res.Status, string(data))
	logger.Info.Println("Barcode send")
	return nil
}
