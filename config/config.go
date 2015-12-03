package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/arteev/zbarnet/barcode"
	"fmt"
	"os/user"
)

// Type of sources
const (
	SourceUnknown = ""
	SourceZBar    = "zbar"
)

//Errors for config
var (
	ErrUnknownOutput   = errors.New("Config error: unknown type output")
	ErrUnknownSource   = errors.New("Config error: unknown source")
	ErrUnknownLocation = errors.New("Config error: unknown location zbarcam")

	ErrUnknownHTTPMethod = errors.New("Config error: unknown http method")
	ErrUnknownHTTPUrl    = errors.New("Config error: unknown http url")
	ErrEmptyHTTPKeyHdr   = errors.New("Config error: empty apikey but use request header authorization")
)

const defaultConfig = ".zbarnet.json"

//Pre defined of the http methods
const (
	HTTPUnknown = HTTPMethod("")
	HTTPGET     = HTTPMethod("GET")
	HTTPPOST    = HTTPMethod("POST")
)

//A HTTPMethod type request on server
type HTTPMethod string

//A Config it setting of application from .zbarnet.json
type Config struct {
	Source string
	Output string
	Once   bool
	ZBar   *zbarconfig
	HTTP   *http
}

//A zbarconfig parameters for zbarcam
type zbarconfig struct {
	mode     barcode.Mode
	Enabled  bool
	Location string
	Device   string
	Args     []string
}

type http struct {
	Enabled      bool
	APIKeyHeader bool
	URL          string
	APIKey       string
	Method       HTTPMethod
}

//MustConfig Load from file "path" in buffer []byte
func MustConfig(path string) *Config {
	if path == "" {
		path = defaultConfig
		if _, e := os.Stat(path); e != nil && os.IsNotExist(e) {
			cd,_:=filepath.Abs(filepath.Dir(os.Args[0]))
			path = filepath.Join(cd, defaultConfig)
		}
		if _, e := os.Stat(path); e != nil && os.IsNotExist(e) {
			u,e:=user.Current();
			if e==nil {
				path = filepath.Join(u.HomeDir, defaultConfig)
			}
		}
	}
	data, e := ioutil.ReadFile(path)
	if e != nil {
		panic(e)
	}
	result, err := parse(data)
	if err != nil {
		panic(err)
	}
	if err := result.check(); err != nil {
		panic(err)
	}
	return result
}

func (z *zbarconfig) GetMode() barcode.Mode {
	return z.mode
}

//check parameters from configuration file
func (c *Config) check() error {
	if c.Output != "json" && c.Output != "raw" && c.Output != "" {
		return ErrUnknownOutput
	}
	if c.Source != "zbar" {
		return ErrUnknownSource
	}
	if c.ZBar != nil && c.ZBar.Enabled {
		if c.ZBar.Location == "" {
			return ErrUnknownLocation
		}
	}
	if e := c.checkHTTP(); e != nil {
		return e
	}
	return nil
}

func (c *Config) checkHTTP() error {
	if c.HTTP.Enabled {
		if c.HTTP.Method == HTTPUnknown {
			return ErrUnknownHTTPMethod
		}
		if c.HTTP.URL == "" {
			return ErrUnknownHTTPUrl
		}
		if c.HTTP.APIKeyHeader && c.HTTP.APIKey == "" {
			return ErrEmptyHTTPKeyHdr
		}
	}
	return nil
}

//parse Parse data and create *Config
func parse(data []byte) (*Config, error) {
	result := &Config{}
	var str map[string]interface{}
	if err := json.Unmarshal(data, &str); err != nil {
		return nil, err
	}
	result.Source = valuedef(str["source"],"zbar").(string)
	result.Output = valuedef(str["output"],"").(string)
	result.Once = valuedef(str["once"],false).(bool)

	result.ZBar = &zbarconfig{}
	if result.Source == SourceZBar {
		// parse parameters for source ZBar
		if _, ok := str["zbar"]; ok {
			zb := str["zbar"].(map[string]interface{})
			result.ZBar.Enabled = valuedef(zb["enabled"], true).(bool)
			result.ZBar.Device = valuedef(zb["device"], "").(string)
			result.ZBar.Location = valuedef(zb["location"], "zbarcam").(string)
			result.ZBar.Args = iArr2sArr(zb["args"].([]interface{}))
			result.ZBar.mode = barcode.ModeNative
			for _, s := range result.ZBar.Args {
				if s == "--raw" {
					result.ZBar.mode = barcode.ModeRaw
				}
				if s == "--xml" {
					result.ZBar.mode = barcode.ModeXML
				}
			}
		}
	}

	result.HTTP = &http{Enabled: false}
	if _, ok := str["http"]; ok {
		hc := str["http"].(map[string]interface{})
		result.HTTP.APIKey = valuedef(hc["apikey"], "").(string)
		result.HTTP.Enabled = valuedef(hc["enabled"], false).(bool)
		result.HTTP.URL = strings.TrimSpace(valuedef(hc["url"], "").(string))
		result.HTTP.APIKeyHeader = valuedef(hc["apikeyhdr"], false).(bool)
		result.HTTP.Method = HTTPUnknown
		method := HTTPMethod(valuedef(hc["method"], HTTPUnknown).(string))
		switch method {
		case HTTPGET:
			result.HTTP.Method = HTTPGET
			break
		case HTTPPOST:
			result.HTTP.Method = HTTPPOST
			break
		}
	}

	//TODO: read commands
	//TODO: read kvdb
	return result, nil
}

//iArr2sArr cast array of interface{} to array of string
func iArr2sArr(arr []interface{}) []string {
	var result []string
	if arr != nil {
		for _, item := range arr {
			result = append(result, item.(string))
		}
	}
	return result
}

func valuedef(value, def interface{}) interface{} {
	if value == nil {
		return def
	}
	return value
}
