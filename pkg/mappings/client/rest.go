package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const timeoutDefault = 5 // Seconds
var timeout *uint8

func init() {
	if timeout == nil {
		SetTimeout(timeoutDefault)
	}
}

func SetTimeout(seconds uint8) {
	timeout = &seconds
}

func NewClient(u string) Client {
	var err error
	cl := Client{}
	temp, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}

	cl.u = *temp

	if !strings.HasSuffix(cl.u.Path, "api") {
		log.Fatal("Please set the path as `/api`")
	}

	cl.client = &http.Client{
		Transport: http.DefaultTransport,
	}

	return cl
}

type Client struct {
	u      url.URL
	client *http.Client
}

type MappingRequest struct {
	Key string `json:"key"`
}

type MappingResponse struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

type CreateMappingResponse struct {
	Key string `json:"key"`
}

type CreateMappingRequest struct {
	URL string `json:"url"`
}

type DeleteMappingResponse struct {
	Deleted string `json:"deleted"`
}

type DeleteMappingRequest struct {
	Key string `json:"key"`
}

type MappingCounterResponse struct {
	Key     string `json:"key"`
	Counter uint32 `json:"counter"`
}

type ErrorResponse struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

func (er ErrorResponse) Error() string {
	return er.Key
}

func (c Client) GetMapping(req MappingRequest) (m MappingResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()
	return c.GetMappingWithContext(ctx, req)
}

func (c Client) GetMappingWithContext(ctx context.Context, req MappingRequest) (m MappingResponse, err error) {
	c.u.Path = path.Join(c.u.Path, "mappings", req.Key)
	err = c.executeRequest(request{
		ctx:    ctx,
		method: http.MethodGet,
		url:    c.u.String(),
	}, &m)
	if err != nil {
		return
	}

	return
}

func (c Client) GetMappingCounter(req MappingRequest) (MappingCounterResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()
	return c.GetMappingCounterWithContext(ctx, req)
}

func (c Client) GetMappingCounterWithContext(ctx context.Context, req MappingRequest) (m MappingCounterResponse, err error) {
	c.u.Path = path.Join(c.u.Path, "mappings", req.Key, "redirects")
	err = c.executeRequest(request{
		ctx:    ctx,
		method: http.MethodGet,
		url:    c.u.String(),
	}, &m)
	if err != nil {
		return
	}

	return
}

func (c Client) CreateMapping(req CreateMappingRequest) (CreateMappingResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()
	return c.CreateMappingWithContext(ctx, req)
}

func (c Client) CreateMappingWithContext(ctx context.Context, req CreateMappingRequest) (m CreateMappingResponse, err error) {
	c.u.Path = path.Join(c.u.Path, "mappings")
	err = c.executeRequest(request{
		ctx:    ctx,
		method: http.MethodPost,
		url:    c.u.String(),
		body:   req,
	}, &m)
	if err != nil {
		return
	}

	return
}

func (c Client) DeleteMapping(req DeleteMappingRequest) (DeleteMappingResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()
	return c.DeleteMappingWithContext(ctx, req)
}

func (c Client) DeleteMappingWithContext(ctx context.Context, req DeleteMappingRequest) (m DeleteMappingResponse, err error) {
	c.u.Path = path.Join(c.u.Path, "mappings", req.Key)
	err = c.executeRequest(request{
		ctx:    ctx,
		method: http.MethodDelete,
		url:    c.u.String(),
	}, &m)
	if err != nil {
		return
	}

	return
}

type request struct {
	ctx    context.Context
	method string
	url    string
	body   interface{}
}

func (c Client) executeRequest(r request, resp interface{}) error {
	var body io.Reader
	if r.body != nil {
		b, err := json.Marshal(r.body)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}

	httpReq, err := http.NewRequestWithContext(r.ctx, r.method, r.url, body)
	if err != nil {
		return err
	}

	rs, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer rs.Body.Close()
	result, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(result, resp)
}
