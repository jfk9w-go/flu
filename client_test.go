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

func (p *Post) String() string {
	return fmt.Sprintf("ID: %d\nUserID: %d\nTitle: %s\nBody: %s\n", p.ID, p.UserID, p.Title, p.Body)
}

// Example_Get provides an example of performing a GET request.
// See https://jsonplaceholder.typicode.com/ for endpoint description.
func Example_Get() {
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

	fmt.Println(post)

	// Output:
	// ID: 1
	// UserID: 1
	// Title: sunt aut facere repellat provident occaecati excepturi optio reprehenderit
	// Body: quia et suscipit
	// suscipit recusandae consequuntur expedita et cum
	// reprehenderit molestiae ut ut quas totam
	// nostrum rerum est autem sunt rem eveniet architecto
}

// Example_Get_QueryParams provides an example of performing a GET request with query parameters.
// See https://jsonplaceholder.typicode.com/ for endpoint description.
func Example_Get_QueryParams() {
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

// Example_Post provides an example of performing a POST request.
// See https://jsonplaceholder.typicode.com/ for endpoint description.
func Example_Post() {
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

	fmt.Println(post)

	// Output:
	// ID: 101
	// UserID: 10
	// Title: lorem ipsum
	// Body: body
}

// Example_Put provides an example of performing a PUT request.
// See https://jsonplaceholder.typicode.com/ for endpoint description.
func Example_Put() {
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

	fmt.Println(post)

	// Output:
	// ID: 1
	// UserID: 10
	// Title: lorem ipsum
	// Body: body
}

// Example_Patch provides an example of performing a PATCH request.
// See https://jsonplaceholder.typicode.com/ for endpoint description.
func Example_Patch() {
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

	fmt.Println(post)

	// Output:
	// ID: 1
	// UserID: 10
	// Title: lorem ipsum
	// Body: body
}

// Example_Delete provides an example of performing a DELETE request.
// See https://jsonplaceholder.typicode.com/ for endpoint description.
func Example_Delete() {
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
