package flu_test

import (
	"fmt"
	"net/http"

	. "github.com/jfk9w-go/flu"
)

type Post struct {
	ID     int    `json:"id,omitempty"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// ExampleGet provides an example of performing a GET request.
// See https://jsonplaceholder.typicode.com/ for endpoint description.
func ExampleGet() {
	post := new(Post)
	err := NewClient(nil).NewRequest().
		Endpoint("https://jsonplaceholder.typicode.com/posts/1").
		Get().
		ReadBody(JSON(post)).
		Error

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("UserID: %d\n", post.UserID)
	fmt.Printf("ID: %d\n", post.ID)
	fmt.Printf("Title: %s\n", post.Title)
	fmt.Printf("ReadBodyFunc: %s\n", post.Body)

	// Output:
	// UserID: 1
	// ID: 1
	// Title: sunt aut facere repellat provident occaecati excepturi optio reprehenderit
	// ReadBodyFunc: quia et suscipit
	// suscipit recusandae consequuntur expedita et cum
	// reprehenderit molestiae ut ut quas totam
	// nostrum rerum est autem sunt rem eveniet architecto
}

// ExampleGet provides an example of performing a Get request.
// See https://jsonplaceholder.typicode.com/ for endpoint description.
func ExampleGetQueryParams() {
	posts := make([]Post, 0)
	err := NewClient(nil).NewRequest().
		Endpoint("https://jsonplaceholder.typicode.com/posts").
		QueryParam("userId", "1").
		Get().
		ReadBody(JSON(&posts)).
		Error

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Posts count: %d\n", len(posts))

	// Output:
	// Posts count: 10
}

// ExamplePost provides an example of performing a POST request.
// See https://jsonplaceholder.typicode.com/ for endpoint description.
func ExamplePost() {
	post := new(Post)
	err := NewClient(nil).NewRequest().
		Endpoint("https://jsonplaceholder.typicode.com/posts").
		Body(JSON(&Post{
			UserID: 10,
			Title:  "lorem ipsum",
			Body:   "body",
		})).
		Post().
		ReadBody(JSON(post)).
		Error

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("ID: %d\n", post.ID)
	fmt.Printf("UserID: %d\n", post.UserID)
	fmt.Printf("Title: %s\n", post.Title)
	fmt.Printf("ReadBodyFunc: %s\n", post.Body)

	// Output:
	// ID: 101
	// UserID: 10
	// Title: lorem ipsum
	// ReadBodyFunc: body
}

// ExamplePut provides an example of performing a PUT request.
// See https://jsonplaceholder.typicode.com/ for endpoint description.
func ExamplePut() {
	post := new(Post)
	err := NewClient(nil).NewRequest().
		Endpoint("https://jsonplaceholder.typicode.com/posts/1").
		Body(JSON(&Post{
			UserID: 10,
			Title:  "lorem ipsum",
			Body:   "body",
		})).
		Put().
		ReadBody(JSON(post)).
		Error

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("ID: %d\n", post.ID)
	fmt.Printf("UserID: %d\n", post.UserID)
	fmt.Printf("Title: %s\n", post.Title)
	fmt.Printf("ReadBodyFunc: %s\n", post.Body)

	// Output:
	// ID: 1
	// UserID: 10
	// Title: lorem ipsum
	// ReadBodyFunc: body
}

// ExamplePatch provides an example of performing a PATCH request.
// See https://jsonplaceholder.typicode.com/ for endpoint description.
func ExamplePatch() {
	post := new(Post)
	err := NewClient(nil).NewRequest().
		Endpoint("https://jsonplaceholder.typicode.com/posts/1").
		Body(JSON(&Post{
			UserID: 10,
			Title:  "lorem ipsum",
			Body:   "body",
		})).
		Patch().
		ReadBody(JSON(post)).
		Error

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("ID: %d\n", post.ID)
	fmt.Printf("UserID: %d\n", post.UserID)
	fmt.Printf("Title: %s\n", post.Title)
	fmt.Printf("ReadBodyFunc: %s\n", post.Body)

	// Output:
	// ID: 1
	// UserID: 10
	// Title: lorem ipsum
	// ReadBodyFunc: body
}

// ExampleDelete provides an example of performing a DELETE request.
// See https://jsonplaceholder.typicode.com/ for endpoint description.
func ExampleDelete() {
	resp := new(string)
	err := NewClient(nil).NewRequest().
		Endpoint("https://jsonplaceholder.typicode.com/posts/1").
		Delete().
		StatusCodes(http.StatusOK).
		ReadBytesFunc(func(data []byte) error {
			*resp = string(data)
			return nil
		}).
		Error

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Response: %s\n", *resp)

	// Output:
	// Response: {}
}
