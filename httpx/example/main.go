package main

import (
	"fmt"
	"log"

	. "github.com/jfk9w-go/flu"
	. "github.com/jfk9w-go/flu/httpx"
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

var ExampleClient = NewClient(nil)

func main() {
	log.Printf("GET")
	Example_GET()
	log.Printf("GET with query parameters")
	Example_GET_QueryParams()
	log.Printf("POST")
	Example_POST()
	log.Printf("PUT")
	Example_PUT()
	log.Printf("PATCH")
	Example_PATCH()
	log.Printf("DELETE")
	Example_DELETE()
}

// Example_GET provides an example of performing a GET request.
// See https://jsonplaceholder.typicode.com/ for resource description.
func Example_GET() {
	response := new(Post)
	err := ExampleClient.
		GET("https://jsonplaceholder.typicode.com/posts/1").
		Execute().
		DecodeBody(JSON(response)).
		Error
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response)
	// Output:
	// ID: 1
	// UserID: 1
	// Title: sunt aut facere repellat provident occaecati excepturi optio reprehenderit
	// Body: quia et suscipit
	// suscipit recusandae consequuntur expedita et cum
	// reprehenderit molestiae ut ut quas totam
	// nostrum rerum est autem sunt rem eveniet architecto
}

// Example_GET_QueryParams provides an example of performing a GET request with query parameters.
// See https://jsonplaceholder.typicode.com/ for resource description.
func Example_GET_QueryParams() {
	response := make([]Post, 0)
	err := ExampleClient.
		GET("https://jsonplaceholder.typicode.com/posts").
		QueryParam("userId", "1").
		Execute().
		DecodeBody(JSON(&response)).
		Error
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Posts count: %d\n", len(response))
	// Output:
	// Posts count: 10
}

// Example_POST provides an example of performing a POST request.
// See https://jsonplaceholder.typicode.com/ for resource description.
func Example_POST() {
	post := &Post{
		UserID: 10,
		Title:  "lorem ipsum",
		Body:   "body",
	}
	response := new(Post)
	err := ExampleClient.
		POST("https://jsonplaceholder.typicode.com/posts").
		Body(JSON(post)).
		Execute().
		DecodeBody(JSON(response)).
		Error
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response)
	// Output:
	// ID: 101
	// UserID: 10
	// Title: lorem ipsum
	// Body: body
}

// Example_PUT provides an example of performing a PUT request.
// See https://jsonplaceholder.typicode.com/ for resource description.
func Example_PUT() {
	post := &Post{
		UserID: 10,
		Title:  "lorem ipsum",
		Body:   "body",
	}
	response := new(Post)
	err := ExampleClient.
		PUT("https://jsonplaceholder.typicode.com/posts/1").
		Body(JSON(post)).
		Execute().
		DecodeBody(JSON(response)).
		Error
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response)
	// Output:
	// ID: 1
	// UserID: 10
	// Title: lorem ipsum
	// Body: body
}

// Example_PATCH provides an example of performing a PATCH request.
// See https://jsonplaceholder.typicode.com/ for resource description.
func Example_PATCH() {
	post := &Post{
		UserID: 10,
		Title:  "lorem ipsum",
		Body:   "body",
	}
	response := new(Post)
	err := ExampleClient.
		PATCH("https://jsonplaceholder.typicode.com/posts/1").
		Body(JSON(post)).
		Execute().
		DecodeBody(JSON(response)).
		Error
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response)
	// Output:
	// ID: 1
	// UserID: 10
	// Title: lorem ipsum
	// Body: body
}

// Example_DELETE provides an example of performing a DELETE request.
// See https://jsonplaceholder.typicode.com/ for resource description.
func Example_DELETE() {
	response := PlainText("")
	err := ExampleClient.
		DELETE("https://jsonplaceholder.typicode.com/posts/1").
		Execute().
		Decode(response).
		Error
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Response: %s\n", response.Value)
	// Output:
	// Response: {}
}
