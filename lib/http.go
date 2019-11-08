// Package base implements a common library for Abac APIs.
// Copyright 2019 Policy Center Author. All Rights Reserved.
// The license belongs to Platform Team.
// Version 1.0 .
package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"wwwin-github.cisco.com/DevNet/restful/log"
)

var (
	httpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:           1000,
			MaxIdleConnsPerHost:    500,
			IdleConnTimeout:        120 * time.Second,
			TLSHandshakeTimeout:    10 * time.Second,
			ExpectContinueTimeout:  1 * time.Second,
			ResponseHeaderTimeout:  10 * time.Second,
			MaxResponseHeaderBytes: 4 << 20,
		},
		Timeout: 30 * time.Second,
	}
)

// CorrelationID is the name of correlationId in context
const CorrelationID = "X-B3-Traceid"

// HTTPRequest sends a http request with headers
func HTTPRequest(method, uri string, body []byte, ctxRequest *http.Request, headers map[string]string) ([]byte, error) {
	if ctxRequest != nil {
		beforeRequest := time.Now().UnixNano()
		defer logRequest(ctxRequest, method, uri, beforeRequest)
	}
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if ctxRequest != nil {
		for _, cookie := range ctxRequest.Cookies() {
			req.AddCookie(cookie)
		}
		traceRequest(ctxRequest, req)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP %s request failed - %s", method, err.Error())
	}
	if resp.StatusCode >= 400 || resp.StatusCode < 200 {
		return nil, fmt.Errorf("%s %s: %s", method, uri, resp.Status)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return respBody, err
}

func logRequest(ctxRequest *http.Request, method, url string, beginTimeStemp int64) {
	cost := (time.Now().UnixNano() - beginTimeStemp) / 1000000
	log.LoggerFromRequest(ctxRequest).
		WithField("stat", "downstream").
		WithField("url", url).
		WithField("delayms", cost).
		Infof("Completed %s request %s in %d ms.", method, url, cost)
}

func traceRequest(ctxRequest, workRequest *http.Request) {
	wireContext, err := opentracing.GlobalTracer().Extract(
		opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(ctxRequest.Header))
	if v, ok := ctxRequest.Context().Value(CorrelationID).(string); ok && err != nil {
		workRequest.Header.Set(CorrelationID, v)
		return
	}
	_ = opentracing.GlobalTracer().Inject(wireContext,
		opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(workRequest.Header))
}
