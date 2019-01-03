package flu_test

import (
	"fmt"

	. "github.com/jfk9w-go/flu"
)

// ExampleGET provides an example of performing a GET request.
// See https://jsonplaceholder.typicode.com/ for endpoint description.
func ExampleGET() {
	resp := new(struct {
		UserID    int    `json:"userId"`
		ID        int    `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	})

	err := NewClient(nil).NewRequest().
		Endpoint("https://jsonplaceholder.typicode.com/todos/1").
		GET().Retrieve().
		ReadJSON(resp).
		Error

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("UserID: %d\n", resp.UserID)
	fmt.Printf("ID: %d\n", resp.ID)
	fmt.Printf("Title: %s\n", resp.Title)
	fmt.Printf("Completed: %t\n", resp.Completed)

	// Output:
	// UserID: 1
	// ID: 1
	// Title: delectus aut autem
	// Completed: false
}

// ExamplePOST provides an example of performing a POST request.
// See https://jsonplaceholder.typicode.com/ for endpoint description.
func ExamplePOST() {
	var req, resp struct {
		ID     int    `json:"id,omitempty"`
		UserID int    `json:"userId"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}

	req.UserID = 10
	req.Title = "lorem ipsum"
	req.Body = "body"

	err := NewClient(nil).NewRequest().
		Endpoint("https://jsonplaceholder.typicode.com/posts").
		POST().Body(JSON(req)).Retrieve().
		ReadJSON(&resp).
		Error

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("ID: %d\n", resp.ID)
	fmt.Printf("UserID: %d\n", resp.UserID)
	fmt.Printf("Title: %s\n", resp.Title)
	fmt.Printf("Body: %s\n", resp.Body)

	// Output:
	// ID: 101
	// UserID: 10
	// Title: lorem ipsum
	// Body: body
}
