package http_test

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jfk9w-go/flu"
	fluhttp "github.com/jfk9w-go/flu/http"
	"github.com/stretchr/testify/assert"
)

func TestClient_GET_Basic(t *testing.T) {
	server := httptest.NewServer(RequestHandlerFunc(func(req *http.Request) {
		assert.Equal(t, http.MethodGet, req.Method, "handler 1 method")
		assert.Equal(t, "Basic "+base64.StdEncoding.EncodeToString([]byte("username:password")), req.Header.Get("Authorization"), "handler 1 auth")
		assert.Equal(t, "/test/path", req.URL.Path, "handler 1 path")
		assert.Equal(t, "a=1&b=2&c=3", req.URL.Query().Encode(), "handler 1 query")
	}))

	defer server.Close()

	text := new(flu.PlainText)
	err := fluhttp.NewClient(nil).
		GET(server.URL+"/test/path").
		QueryParam("a", "1").
		QueryParam("b", "2").
		QueryParam("c", "3").
		Auth(fluhttp.Basic("username", "password")).
		Execute().
		CheckStatus(http.StatusOK).
		DecodeBody(text).
		Error
	assert.Nil(t, err, "client 1 error")
	assert.Equal(t, "OK", text.Value, "client 1 response")
}

func TestClient_GET_ResponseJSON(t *testing.T) {
	server := httptest.NewServer(ConstHandler{
		StatusCode: http.StatusOK,
		Response:   `{"status": "OK"}`,
	})

	defer server.Close()

	type StatusResponse struct {
		Status string `json:"status"`
	}

	status := new(StatusResponse)
	err := fluhttp.NewClient(nil).
		GET(server.URL).
		Execute().
		CheckStatus(http.StatusOK).
		DecodeBody(flu.JSON{status}).
		Error
	assert.Nil(t, err, "client 2 error")
	assert.Equal(t, "OK", status.Status)
}

func TestClient_GET_StatusCodeError(t *testing.T) {
	server := httptest.NewServer(ConstHandler{
		StatusCode: http.StatusInternalServerError,
		Response:   "request failed",
	})

	defer server.Close()

	err := fluhttp.NewClient(nil).
		GET(server.URL).
		Execute().
		CheckStatus(http.StatusOK).
		Error
	assert.Equal(t, fluhttp.StatusCodeError{
		Code: http.StatusInternalServerError,
		Text: "request failed",
	}, err, "client 3 error")
}

func TestClient_POST_JSON(t *testing.T) {
	client := fluhttp.NewClient(nil).
		AcceptStatus(http.StatusOK, http.StatusCreated)

	type Post struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	server := httptest.NewServer(ConstHandler{
		RequestHandler: func(req *http.Request) {
			assert.Equal(t, http.MethodPost, req.Method, "handler 1 method")
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"), "handler 1 content type")
			post := new(Post)
			err := flu.DecodeFrom(flu.IO{R: req.Body}, flu.JSON{post})
			assert.Nil(t, err, "handler 1 decode error")
			assert.Equal(t, 1, post.ID, "handler 1 post id")
			assert.Equal(t, "Test Post", post.Name, "handler 1 post name")
		},
		StatusCode: http.StatusCreated,
		Response:   `{"id": 1}"`,
	})

	defer server.Close()

	request := Post{
		ID:   1,
		Name: "Test Post",
	}
	response := new(Post)
	err := client.POST(server.URL).
		BodyEncoder(flu.JSON{request}).
		Execute().
		DecodeBody(flu.JSON{response}).
		Error
	assert.Nil(t, err, "client 1 error")
	assert.Equal(t, &Post{ID: 1}, response, "client 1 response")
}

func TestClient_POST_Form(t *testing.T) {
	client := fluhttp.NewClient(nil).
		AcceptStatus(http.StatusOK, http.StatusCreated)

	type Post struct {
		ID   int    `url:"id"`
		Name string `url:"name"`
	}

	server := httptest.NewServer(ConstHandler{
		RequestHandler: func(req *http.Request) {
			assert.Equal(t, http.MethodPost, req.Method, "handler 1 method")
			assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"), "handler 1 content type")
			body := new(flu.PlainText)
			err := flu.DecodeFrom(flu.IO{R: req.Body}, body)
			assert.Nil(t, err, "handler 1 decode error")
			assert.Equal(t, "id=1&name=Test+Post&option=check", body.Value, "handler 1 body")
		},
		StatusCode: http.StatusCreated,
	})

	defer server.Close()

	err := client.POST(server.URL).
		BodyEncoder(fluhttp.Form{}.
			Value(Post{
				ID:   1,
				Name: "Test Post",
			}).
			Set("option", "check")).
		Execute().
		Error
	assert.Nil(t, err, "client 1 error")
}

func TestClient_POST_MultipartFormData(t *testing.T) {
	client := fluhttp.NewClient(nil).
		AcceptStatus(http.StatusOK, http.StatusCreated)

	type Post struct {
		ID   int    `url:"id"`
		Name string `url:"name"`
	}

	server := httptest.NewServer(ConstHandler{
		RequestHandler: func(req *http.Request) {
			assert.Equal(t, http.MethodPost, req.Method, "handler 1 method")
			contentType := strings.Split(req.Header.Get("Content-Type"), "; boundary=")
			assert.Equal(t, "multipart/form-data", contentType[0], "handler 1 content type")
			err := req.ParseMultipartForm(1 << 20)
			assert.Nil(t, err, "handler 1 multipart parse")
			assert.Equal(t, map[string][]string{"id": {"1"}, "name": {"Test Post"}, "option": {"check"}}, req.MultipartForm.Value, "handler 1 multipart values")
			fileHeader := req.MultipartForm.File["file"][0]
			assert.Equal(t, "photo.jpg", fileHeader.Filename, "handler 1 file name")
			assert.Equal(t, int64(5), fileHeader.Size, "handler 1 file size")
			file, err := fileHeader.Open()
			assert.Nil(t, err, "handler 1 file open")
			text := new(flu.PlainText)
			err = flu.DecodeFrom(flu.IO{R: file}, text)
			assert.Nil(t, err, "handler 1 file decode")
			assert.Equal(t, "TESTE", text.Value, "handler 1 file content")
		},
		StatusCode: http.StatusCreated,
	})

	defer server.Close()

	buf := flu.NewBuffer()
	_, _ = buf.WriteString("TESTE")
	err := client.POST(server.URL).
		BodyEncoder(fluhttp.NewMultipartForm().
			Value(Post{
				ID:   1,
				Name: "Test Post",
			}).
			Set("option", "check").
			File("file", "photo.jpg", buf)).
		Execute().
		Error
	assert.Nil(t, err, "client 1 error")
}

type ConstHandler struct {
	RequestHandler RequestHandlerFunc
	StatusCode     int
	Response       string
}

func (h ConstHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	if h.RequestHandler != nil {
		h.RequestHandler(req)
	}

	writer.WriteHeader(h.StatusCode)
	_, _ = writer.Write([]byte(h.Response))
}

type RequestHandlerFunc func(req *http.Request)

func (f RequestHandlerFunc) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	f(req)
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("OK"))
}
