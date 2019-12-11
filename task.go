package qencode

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type TaskService interface {
	SessionToken(context.Context) (*TokenRoot, *Response, error)
	Create(context.Context) (*CreateRoot, *Response, error)
	Encode(context.Context, *TaskParams) (*EncodeRoot, *Response, error)
}

type TaskServiceOp struct {
	client *Client
}

type CreateRoot struct {
	Error     int8   `json:"error,omitempty"`
	Message   string `json:"message,omitempty"`
	UploadURL string `json:"upload_url,omitempty"`
	TaskToken string `json:"task_token,omitempty"`
}

type EncodeRoot struct {
	Error     int8   `json:"error,omitempty"`
	Message   string `json:"message,omitempty"`
	StatusURL string `json:"status_url,omitempty"`
}

var _ TaskService = &TaskServiceOp{}

func (t *TaskServiceOp) SessionToken(ctx context.Context) (*TokenRoot, *Response, error) {
	path := "access_token"

	if t.client.Access != nil {
		current := time.Now().Unix()
		if t.client.Access.ExpInMilli > current {
			return t.client.Access, nil, nil
		}
	}

	fmt.Println("Request")
	payload := strings.NewReader(fmt.Sprintf("api_key=%v", t.client.Key))
	req, err := t.client.NewRequest(ctx, http.MethodPost, path, payload)
	if err != nil {
		return nil, nil, err
	}

	root := new(TokenRoot)

	resp, err := t.client.DO(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if root.Error > 0 {
		return root, resp, errors.New(root.Message)
	}

	t.client.Access = root
	_, err = t.setUnixTime()
	if err != nil {
		return nil, nil, err
	}

	return root, resp, nil

}

func (t *TaskServiceOp) Create(ctx context.Context) (*CreateRoot, *Response, error) {

	access, _, err := t.SessionToken(ctx)
	if err != nil {
		return nil, nil, err
	}

	path := "create_task"
	payload := strings.NewReader(fmt.Sprintf("token=%v", access.Token))
	req, err := t.client.NewRequest(ctx, http.MethodPost, path, payload)
	if err != nil {
		return nil, nil, err
	}

	root := new(CreateRoot)

	resp, err := t.client.DO(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if root.Error > 0 {
		return root, resp, errors.New(root.Message)
	}

	return root, resp, nil
}

//Encode beigns the custom task
func (t *TaskServiceOp) Encode(ctx context.Context, params *TaskParams) (*EncodeRoot, *Response, error) {
	path := "start_encode2"

	//query := fmt.Sprintf(schema2, params.SourcePath, t.client.CallbackURL, params.FinalPath, t.client.StorageKey, t.client.StorageSecret, Resolutions["540p"])
	query, err := QueryBuilder(params, t)
	if err != nil {
		return nil, nil, err
	}

	payload := strings.NewReader(fmt.Sprintf(`task_token=%v&payload=%v&query=%v`, params.TaskToken, params.Payload, query))
	fmt.Println(payload)
	// fmt.Println("+++++++++ NIL ++++++++")
	// return nil, nil, nil

	req, err := t.client.NewRequest(ctx, http.MethodPost, path, payload)
	if err != nil {
		return nil, nil, err
	}

	root := new(EncodeRoot)

	resp, err := t.client.DO(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if root.Error > 0 {
		return root, resp, errors.New(root.Message)
	}
	return root, resp, nil
}

func (t *TaskServiceOp) Status(ctx context.Context, token string) (*CreateRoot, *Response, error) {
	path := "status"
	payload := strings.NewReader(fmt.Sprintf("task_tokens=%v", token))
	req, err := t.client.NewRequest(ctx, http.MethodPost, path, payload)
	if err != nil {
		return nil, nil, err
	}

	root := new(CreateRoot)

	resp, err := t.client.DO(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

func (t *TaskServiceOp) setUnixTime() (bool, error) {
	te, err := time.Parse(time.RFC3339, fmt.Sprintf("%v+00:00", t.client.Access.Expire))
	if err != nil {
		return false, err
	}

	t.client.Access.ExpInMilli = te.Unix()
	return true, nil
}
