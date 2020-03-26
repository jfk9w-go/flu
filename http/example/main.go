package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	. "github.com/jfk9w-go/flu"
	. "github.com/jfk9w-go/flu/http"
)

type Post struct {
	ID     int    `json:"id,omitempty"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func (p *Post) String() string {
	return fmt.Sprintf("ID: %d\nUserID: %d\nTitle: %s\nBody: %s\n", p.ID, p.UserID, p.Title, p.Body)
}

var ExampleClient = NewTransport().
	ResponseHeaderTimeout(10*time.Second).
	NewClient().
	AcceptStatus(http.StatusOK, http.StatusCreated)

func main() {
	exampleGet()
	exampleGetQueryParams()
	examplePost()
	examplePut()
	examplePatch()
	exampleDelete()
}

// exampleGet provides an example of performing a GET request.
// See https://jsonplaceholder.typicode.com/ for resource description.
func exampleGet() {
	response := new(Post)
	err := ExampleClient.
		GET("https://jsonplaceholder.typicode.com/posts/1").
		Execute().
		DecodeBody(JSON{response}).
		Error
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("GET\n%s", response)
	// Output:
	// GET
	// ID: 1
	// UserID: 1
	// Title: sunt aut facere repellat provident occaecati excepturi optio reprehenderit
	// ContentType: quia et suscipit
	// suscipit recusandae consequuntur expedita et cum
	// reprehenderit molestiae ut ut quas totam
	// nostrum rerum est autem sunt rem eveniet architecto
}

// exampleGetQueryParams provides an example of performing a GET request with query parameters.
// See https://jsonplaceholder.typicode.com/ for resource description.
func exampleGetQueryParams() {
	response := make([]Post, 0)
	err := ExampleClient.
		GET("https://jsonplaceholder.typicode.com/posts").
		QueryParam("userId", "1").
		Execute().
		DecodeBody(JSON{&response}).
		Error
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("GET query params\nPosts count: %d\n", len(response))
	// Output:
	// GET
	// Posts count: 10
}

// examplePost provides an example of performing a POST request.
// See https://jsonplaceholder.typicode.com/ for resource description.
func examplePost() {
	post := &Post{
		UserID: 10,
		Title:  "lorem ipsum",
		Body:   "body",
	}
	response := new(Post)
	err := ExampleClient.
		POST("https://jsonplaceholder.typicode.com/posts").
		BodyEncoder(JSON{post}).
		Execute().
		DecodeBody(JSON{response}).
		Error
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("POST\n%s", response)
	// Output:
	// POST
	// ID: 101
	// UserID: 10
	// Title: lorem ipsum
	// ContentType: body
}

// examplePut provides an example of performing a PUT request.
// See https://jsonplaceholder.typicode.com/ for resource description.
func examplePut() {
	post := &Post{
		UserID: 10,
		Title:  "lorem ipsum",
		Body:   "body",
	}
	response := new(Post)
	err := ExampleClient.
		PUT("https://jsonplaceholder.typicode.com/posts/1").
		BodyEncoder(JSON{post}).
		Execute().
		DecodeBody(JSON{response}).
		Error
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("PUT\n%s", response)
	// Output:
	// PUT
	// ID: 1
	// UserID: 10
	// Title: lorem ipsum
	// ContentType: body
}

// examplePatch provides an example of performing a PATCH request.
// See https://jsonplaceholder.typicode.com/ for resource description.
func examplePatch() {
	post := &Post{
		UserID: 10,
		Title:  "lorem ipsum",
		Body:   "body",
	}
	response := new(Post)
	err := ExampleClient.
		PATCH("https://jsonplaceholder.typicode.com/posts/1").
		BodyEncoder(JSON{post}).
		Execute().
		DecodeBody(JSON{response}).
		Error
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("PATCH\n%s", response)
	// Output:
	// PATCH
	// ID: 1
	// UserID: 10
	// Title: lorem ipsum
	// ContentType: body
}

// exampleDelete provides an example of performing a DELETE request.
// See https://jsonplaceholder.typicode.com/ for resource description.
func exampleDelete() {
	response := &PlainText{""}
	err := ExampleClient.
		DELETE("https://jsonplaceholder.typicode.com/posts/1").
		Execute().
		DecodeBody(response).
		Error
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("DELETE\nResponse: %s\n", response.Value)
	// Output:
	// DELETE
	// Response: {}
}
