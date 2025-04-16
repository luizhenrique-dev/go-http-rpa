package engine

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	httprequest "github.com/luizhenriquees/go-http-rpa/http_request"
)

const (
	ParamPath = "path"
)

// HTTPTask represents a task that makes an HTTP request
type HTTPTask struct {
	name           string
	method         httprequest.HTTPMethod
	Headers        httprequest.Headers
	Params         Parameters
	requiredParams []string
	urlParam       string
	bodyParam      string
	waitTime       time.Duration
	Logger         Logger
	responseKey    string
}

func (t *HTTPTask) Name() string {
	return t.name
}

// Validate checks if all required parameters exist
func (t *HTTPTask) Validate() error {
	if t.Headers == nil {
		return errors.New("missing headers")
	}

	for _, param := range t.requiredParams {
		if _, exists := t.Params[param]; !exists {
			return errors.New("missing required parameter: " + param)
		}
	}
	return nil
}

// NewHTTPTask creates a new HTTP task
func NewHTTPTask(name string, method httprequest.HTTPMethod, headers httprequest.Headers, params Parameters, options ...Option) *HTTPTask {
	task := &HTTPTask{
		name:        name,
		method:      method,
		Headers:     headers,
		Params:      params,
		waitTime:    time.Second * 2,
		Logger:      &DefaultLogger{prefix: fmt.Sprintf("HTTP Task - %s", name)},
		responseKey: "response_" + name,
	}

	for _, option := range options {
		option(task)
	}

	var requiredParams []string
	if method == httprequest.POST && task.bodyParam != "" {
		requiredParams = append(requiredParams, task.bodyParam)
	}
	task.requiredParams = requiredParams

	return task
}

// Option is a function that configures an HTTPTask
type Option func(*HTTPTask)

// WithURLParam sets the path parameter for the URL
func WithURLParam(param string) Option {
	return func(t *HTTPTask) {
		t.urlParam = param
	}
}

// WithBodyParam sets the content for the request body
func WithBodyParam(param string) Option {
	return func(t *HTTPTask) {
		t.bodyParam = param
	}
}

// WithWaitTime sets the wait time after the request
func WithWaitTime(duration time.Duration) Option {
	return func(t *HTTPTask) {
		t.waitTime = duration
	}
}

// WithResponseKey sets the key where the response will be stored
func WithResponseKey(key string) Option {
	return func(t *HTTPTask) {
		t.responseKey = key
	}
}

// Execute performs the HTTP request
func (t *HTTPTask) Execute() error {
	t.Logger.Info("Initiating task...")
	finalURL := t.BuildURL()
	var resp *http.Response
	var err error

	if t.method == httprequest.GET {
		resp, err = httprequest.DoGet(finalURL, t.Headers)
	} else {
		var body []byte
		if t.bodyParam != "" {
			if bodyValue, ok := t.Params[t.bodyParam].([]byte); ok {
				body = bodyValue
			}
		}
		resp, err = httprequest.DoPost(finalURL, t.Headers, body)
	}

	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	t.Params[t.responseKey] = resp
	if t.waitTime > 0 {
		time.Sleep(t.waitTime)
	}
	t.Logger.Info("Task completed")
	return nil
}

func (t *HTTPTask) BuildURL() string {
	finalURL := t.Params.Get(ParamBaseURL).(string) + t.Params.Get(ParamPath).(string)
	if t.urlParam != "" {
		if urlValue, ok := t.Params.Get(t.urlParam).(string); ok {
			finalURL = finalURL + urlValue
		}
	}
	return finalURL
}
