package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

const (
	TYPE_QUERY = 1
	TYPE_JSON  = 2
)

type RequestAdaptorInterface interface {
	RequestWithJSON(reqType string, address string, data interface{}, header map[string]string) ([]byte, error)
	RequestWithQuery(reqType string, address string, data, header map[string]string) ([]byte, error)
}

type RequestAdaptor struct {
	debug  bool
	log    *zap.Logger
	client *http.Client
}

func NewRequestAdaptor(rt http.RoundTripper, timeOut time.Duration, log *zap.Logger, debug bool) RequestAdaptor {
	return RequestAdaptor{
		debug: debug,
		log:   log,
		client: &http.Client{
			Transport: rt,
			Timeout:   timeOut,
		},
	}
}

func makeHeader(data map[string]string) http.Header {
	h := http.Header{}
	h.Set("Content-Type", "application/json")
	for i, val := range data {
		h.Set(i, val)
	}
	return h
}

func (n *RequestAdaptor) RequestWithJSON(reqType string, address string, data interface{}, header map[string]string) ([]byte, error) {
	errHandle := func(err error) ([]byte, error) {
		n.log.Error("request with json",
			zap.String("address", address),
			zap.Any("data", data),
			zap.Any("header", header),
			zap.Error(err),
		)
		return nil, err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return errHandle(fmt.Errorf("marshal %w", err))
	}

	//parse url to safety string
	host, err := url.Parse(address)
	if err != nil {
		return errHandle(fmt.Errorf("parse url %w", err))
	}

	//init http request
	req, err := http.NewRequest(reqType, host.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return errHandle(fmt.Errorf("init request %w", err))
	}

	//work with header
	req.Header = makeHeader(header)

	//do request
	resp, err := n.client.Do(req)
	if err != nil {
		return errHandle(fmt.Errorf("do request %w", err))
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errHandle(fmt.Errorf("read body %w", err))
	}

	if n.debug {
		n.log.Info("request with json",
			zap.String("address", address),
			zap.Any("data", data),
			zap.Any("header", header),
			zap.String("response", string(body)),
		)
	}

	return body, nil
}

func (n *RequestAdaptor) RequestWithQuery(reqType string, address string, data, header map[string]string) ([]byte, error) {
	errHandle := func(err error) ([]byte, error) {
		n.log.Error("request with query",
			zap.String("address", address),
			zap.Any("data", data),
			zap.Any("header", header),
			zap.Error(err),
		)
		return nil, err
	}

	//parse url to safety string
	host, err := url.Parse(address)
	if err != nil {
		return errHandle(fmt.Errorf("parse url %w", err))
	}

	//init http GET request
	req, err := http.NewRequest(reqType, host.String(), nil)
	if err != nil {
		return errHandle(fmt.Errorf("init request %w", err))
	}

	//work with query
	q := req.URL.Query()
	for k, v := range data {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	//work with header
	req.Header = makeHeader(header)

	//do request
	resp, err := n.client.Do(req)
	if err != nil {
		return errHandle(fmt.Errorf("do request %w", err))
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errHandle(fmt.Errorf("read body %w", err))
	}

	if n.debug {
		n.log.Info("request with json",
			zap.String("address", address),
			zap.Any("data", data),
			zap.Any("header", header),
			zap.String("response", string(body)),
		)
	}

	return body, nil
}
