package engine

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	httprequest "github.com/luizhenriquees/go-http-rpa/http_request"
)

type PreRequestFunc func() error

type PostExtractFunc func(resp *http.Response, task *HTTPTask) error

// HTTPTask represents a task that makes an HTTP request
type HTTPTask struct {
	name           string
	URL            string
	method         httprequest.HTTPMethod
	Headers        httprequest.Headers
	Params         Parameters
	requiredParams []string
	RequestBody    []byte
	waitTime       time.Duration
	Logger         Logger

	preRequestFunc  PreRequestFunc
	postExtractFunc PostExtractFunc
}

func (t *HTTPTask) Name() string {
	return t.name
}

// Validate checks if all required parameters exist
func (t *HTTPTask) Validate() error {
	if t.URL == "" {
		return fmt.Errorf("task %q: missing URL", t.name)
	}

	if t.Headers == nil {
		return errors.New("missing headers")
	}

	for _, param := range t.requiredParams {
		if _, exists := t.Params[param]; !exists {
			return fmt.Errorf("task %q: missing required parameter %q", t.name, param)
		}
	}
	return nil
}

// NewHTTPTask creates a new HTTP task
func NewHTTPTask(name string, method httprequest.HTTPMethod, url string, headers httprequest.Headers, params Parameters, options ...Option) *HTTPTask {
	task := &HTTPTask{
		name:           name,
		method:         method,
		URL:            url,
		Headers:        headers,
		Params:         params,
		waitTime:       time.Second * 2,
		Logger:         &DefaultLogger{prefix: fmt.Sprintf("HTTP Task - %s", name)},
		requiredParams: []string{},
	}
	if task.Headers == nil {
		task.Headers = make(httprequest.Headers)
	}
	if task.Params == nil {
		task.Params = make(Parameters)
	}

	for _, option := range options {
		option(task)
	}
	return task
}

// Option is a function that configures an HTTPTask
type Option func(*HTTPTask)

// WithRequiredParams sets the list of required parameter keys
func WithRequiredParams(params []string) Option {
	return func(t *HTTPTask) {
		t.requiredParams = append([]string(nil), params...)
	}
}

// WithPreRequestFunc sets a custom function for pre-request
func WithPreRequestFunc(fn PreRequestFunc) Option {
	return func(t *HTTPTask) {
		t.preRequestFunc = fn
	}
}

// WithPostExtractFunc sets a custom function for post-extraction
func WithPostExtractFunc(fn PostExtractFunc) Option {
	return func(t *HTTPTask) {
		t.postExtractFunc = fn
	}
}

// Execute performs the HTTP request
func (t *HTTPTask) Execute() error {
	t.Logger.Info("Initiating task %s...", t.name)
	var resp *http.Response
	var err error

	if t.preRequestFunc != nil {
		t.Logger.Info("Executing pre request function...")
		if err := t.preRequestFunc(); err != nil {
			return fmt.Errorf("pre request failed: %w", err)
		}
	}

	if t.method == httprequest.GET {
		t.Logger.Info("Executing GET request...")
		resp, err = httprequest.DoGet(t.URL, t.Headers)
	} else if t.method == httprequest.POST {
		t.Logger.Info("Executing POST request...")
		resp, err = httprequest.DoPost(t.URL, t.Headers, t.RequestBody)
	} else {
		return fmt.Errorf("task %q: unsupported HTTP method %q", t.name, t.method)
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}

	if t.postExtractFunc != nil {
		t.Logger.Info("Executing post extract function...")
		if err := t.postExtractFunc(resp, t); err != nil {
			return fmt.Errorf("post extraction failed: %w", err)
		}
	}

	if t.waitTime > 0 {
		time.Sleep(t.waitTime)
	}
	t.Logger.Info("HTTP Task executed successfully")
	return nil
}
