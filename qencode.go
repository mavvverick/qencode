package qencode

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	libraryVersion = "1.0.0"
	defaultBaseURL = "https://api.qencode.com"
	userAgent      = "quencode/" + libraryVersion
	mediaType      = "application/x-www-form-urlencoded"
)

var db = ""
var secret = ""

type Client struct {
	Client        *http.Client
	BaseURL       *url.URL
	CallbackURL   string
	WATERMARK     string
	Key           string
	Bucket        string
	Access        *TokenRoot
	StorageKey    string
	StorageSecret string
	UserAgent     string
	secret        string
	Task          TaskService
}

type Response struct {
	Response *http.Response
	monitor  string
}

type ErrorResponse struct {
	Response   *http.Response
	StatusCode string `json:"statusCode"`
	Message    string `json:"message"`
}

type TokenRoot struct {
	Error      int8   `json:"error"`
	Message    string `json:"message"`
	Token      string `json:"token"`
	Expire     string `json:"expire"`
	ExpInMilli int64  `json:"exp"`
}

var FrameRate = "30"

var Resolutions = map[string]string{
	"144p":  "144x256",
	"240p":  "240x352",
	"360p":  "360x480",
	"480p":  "480x858",
	"540p":  "540x960",
	"720p":  "720x1280",
	"1080p": "1080x1920",
	"2160p": "2160x3860",
}

type TaskParams struct {
	PostID      string
	Name        string
	TaskToken   string
	Videos      []string
	SourcePath  string
	FinalPath   string
	Resolutions []string
	Payload     string
	StartTime   string
	Duration    string
	IsLogo      bool
	IsH264      bool
	Upscale     bool
}

// NewClient creates http client for api
func NewClient(httpClient *http.Client, key string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Second * 300,
		}
	}

	baseUrl, _ := url.Parse(defaultBaseURL)
	c := &Client{Client: httpClient, BaseURL: baseUrl, UserAgent: userAgent,
		Key: key, secret: secret, StorageKey: "", StorageSecret: "", CallbackURL: "",
		Bucket: "", WATERMARK: ""}
	c.Task = &TaskServiceOp{client: c}
	return c
}

//NewRequest create a request
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body *strings.Reader) (*http.Request, error) {
	path := fmt.Sprintf("/%v/%v", os.Getenv("QUENCODE_VERSION"), urlStr)
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

func (c *Client) DO(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := DoRequestWithClient(ctx, c.Client, req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if rerr := resp.Body.Close(); err != nil {
			err = rerr
		}
	}()

	response := NewResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return nil, err
	}

	if v != nil {
		e, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(e, &v)
		if err != nil {
			return nil, err
		}
	}

	return response, err
}

func DoRequestWithClient(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req)
}

func NewResponse(r *http.Response) *Response {
	response := Response{Response: r}
	return &response
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 399 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}

	return errorResponse
}

func (r *ErrorResponse) Error() string {
	if r.StatusCode != "" {
		return fmt.Sprintf("%v %v: %d (request %q) %v",
			r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.StatusCode, r.Message)
	}
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message)
}
