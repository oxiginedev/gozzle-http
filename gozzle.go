package gozzle

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultUserAgent    string = "gozzle-client-v1.0"
	defaultContentType  string = "application/json; charset=utf-8"
	defaultTimeout      int    = 5
	defaultMaxRedirects int    = 5
)

type Map map[string]interface{}

type Config struct {
	URL          string
	Method       string
	BaseURL      string
	BasicAuth    map[string]string
	BearerToken  string
	Headers      map[string]string
	Timeout      int
	Body         map[string]interface{}
	Params       map[string]interface{}
	UserAgent    string
	AcceptJSON   bool
	Accepts      string
	AsMultipart  bool
	AsURLEncoded bool
	ContentType  string
	StructScan   interface{}
	MaxRedirects int
	Retries      int
	RetrySleep   int
}

type Response struct {
	Data       interface{}
	Status     int
	StatusText string
	Headers    map[string]string
	Config     interface{}
	Request    *http.Request
}

func Send(config *Config) (*Response, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	var buf io.Reader

	if config.Body != nil {
		b, err := json.Marshal(config.Body)
		if err != nil {
			return nil, err
		}

		buf = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(strings.ToUpper(config.Method), config.URL, buf)
	if err != nil {
		return nil, err
	}

	if config.Body != nil {
		req.Header.Set("Content-Type", config.ContentType)
	}

	if config.BasicAuth != nil {
		req.SetBasicAuth(config.BasicAuth["username"], config.BasicAuth["password"])
	}

	if !IsStringEmpty(config.BearerToken) {
		req.Header.Set("Authorization", "Bearer "+config.BearerToken)
	}

	// set necessary headers
	req.Header.Set("Accept", config.Accepts)
	req.Header.Set("User-Agent", config.UserAgent)

	// this overrides any previously set headers
	if config.Headers != nil {
		for k, v := range config.Headers {
			req.Header.Set(k, v)
		}
	}

	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > config.MaxRedirects {
				return errors.New("too many redirects")
			}
			return nil
		},
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Data:       data,
		Status:     res.StatusCode,
		StatusText: http.StatusText(res.StatusCode),
		Config:     config,
		Request:    req,
	}, nil
}

func (c *Config) validate() error {
	c.AcceptJSON = true

	if IsStringEmpty(c.UserAgent) {
		c.UserAgent = defaultUserAgent
	}

	if c.Timeout == 0 {
		c.Timeout = defaultTimeout
	}

	if c.MaxRedirects == 0 {
		c.MaxRedirects = defaultMaxRedirects
	}

	if !IsStringEmpty(c.BearerToken) && c.BasicAuth != nil {
		return errors.New("cannot authenticate with bearer and basic auth")
	}

	if c.Body != nil && c.Method == http.MethodGet {
		return errors.New("body not allowed for 'GET' method")
	}

	if c.AcceptJSON {
		c.Accepts = "application/json"
	}

	if !IsStringEmpty(c.ContentType) {
		c.ContentType = defaultContentType
	}

	if c.AsMultipart {
		c.ContentType = "multipart/form-data"
	}

	if c.AsURLEncoded {
		c.ContentType = "application/x-www-form-url-encoded"
	}

	if !IsStringEmpty(c.BaseURL) {
		c.BaseURL = strings.Trim(c.BaseURL, "/")
	}

	return nil
}
