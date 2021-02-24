package curl

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type (
	Request struct {
		URL      string
		URI      string
		Method   string
		Headers  map[string]string
		Params   url.Values
		Json     interface{}
		IsJson   bool
		Byte     []byte `json:"-"`
		ByteStr  string
		IsByte   bool
		IsNotify bool
		Client   HTTPClient `json:"-"`
		Xml      string
		IsXml    bool
	}

	HTTPClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
)

var client *http.Client

func init() {
	client = &http.Client{
		Timeout: time.Second * 1,
	}
}

func NewRequest() *Request {
	return &Request{
		Headers: map[string]string{},
		Client:  client,
	}
}

func (r *Request) DoRequest() (*http.Response, []byte, error) {
	var request *http.Request
	var response *http.Response

	if err := r.validateURL(); err != nil {
		return response, nil, err
	}

	r.Method = strings.ToUpper(r.Method)
	if err := r.validateMethod(); err != nil {
		return response, nil, err
	}

	u, err := url.Parse(r.URL)
	if err != nil {
		return response, nil, err
	}

	if r.Headers == nil {
		r.Headers = map[string]string{}
	}

	if r.Method == "GET" {
		request, err = r.get(u)
	}

	if request == nil {
		return response, nil, errors.New("Failed create new request")
	}

	if err != nil {
		return response, nil, err
	}

	for key, value := range r.Headers {
		request.Header.Set(key, value)
	}

	response, err = r.Client.Do(request)
	if err != nil {
		return response, nil, err
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response, nil, err
	}

	return response, contents, nil
}

func (r *Request) validateURL() error {
	if r.URL == "" || len(r.URL) == 0 {
		return errors.New("URL is required")
	}
	return nil
}

func (r *Request) validateMethod() error {
	if r.Method != "GET" {
		return errors.New("Unsupported method " + r.Method)
	}
	return nil
}

func (r *Request) get(u *url.URL) (*http.Request, error) {
	u.RawQuery = r.Params.Encode()
	req, err := http.NewRequest(r.Method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}
