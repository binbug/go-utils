package httputils

import (
	"io"
	"net/http"
	"testing"
)

func Test_GetThenReturnResponse(t *testing.T) {
	httpResult := Get[http.Response]("https://jsonplaceholder.typicode.com/posts/1")

	if httpResult.Err() != nil {
		t.Error(httpResult.Err())
	}

	resp := httpResult.ToObject()

	t.Log(resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	t.Log(string(data))
}

func Test_GetThenPrintBody(t *testing.T) {
	httpResult := Get[string]("https://jsonplaceholder.typicode.com/posts/1")
	if httpResult.Err() != nil {
		t.Error(httpResult.Err())
	}

	t.Log(httpResult.ToObject())
}

func Test_GetThenReturnObject(t *testing.T) {
	type Post struct {
		UserId int    `json:"userId"`
		Id     int    `json:"id"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}

	httpResult := Get[Post]("https://jsonplaceholder.typicode.com/posts/1")
	if httpResult.Err() != nil {
		t.Error(httpResult.Err())
	}

	t.Log(httpResult.ToObject())
	if httpResult.ToObject().Id != 1 {
		t.Error("Id should be 1")
	}

	if httpResult.ToObject().Id != 1 {
		t.Error("Id should be 1")
	}

	if httpResult.ToObject().Title == "" {
		t.Error("Title should not be empty")
	}
}

func Test_PostJSON(t *testing.T) {
	type PostReq struct {
		UserId int    `json:"userId"`
		Body   string `json:"body"`
		Title  string `json:"title"`
	}

	httpResult := PostJSON[string]("https://jsonplaceholder.typicode.com/posts", PostReq{
		UserId: 1,
		Body:   "Hello World",
		Title:  "Hello World",
	})

	if httpResult.Err() != nil {
		t.Error(httpResult.Err())
	}

	t.Log(httpResult.ToObject())
}
