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

func newClient() *Client {
	return NewClient(http.DefaultClient).
		AcceptResponseCodes(http.StatusOK, http.StatusCreated)
}

// Example_GET provides an example of performing a GET request.
// See https://jsonplaceholder.typicode.com/ for resource description.
func Example_GET() {
	post := new(Post)
	err := newClient().NewRequest().
		GET().
		Resource("https://jsonplaceholder.typicode.com/posts/1").
		Send().
		DecodeBody(JSON(post)).
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

// Example_GET_QueryParams provides an example of performing a GET request with query parameters.
// See https://jsonplaceholder.typicode.com/ for resource description.
func Example_GET_QueryParams() {
	posts := make([]Post, 0)
	err := newClient().NewRequest().
		GET().
		Resource("https://jsonplaceholder.typicode.com/posts").
		QueryParam("userId", "1").
		Send().
		DecodeBody(JSON(&posts)).
		Error

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Posts count: %d\n", len(posts))

	// Output:
	// Posts count: 10
}

// Example_POST provides an example of performing a POST request.
// See https://jsonplaceholder.typicode.com/ for resource description.
func Example_POST() {
	post := new(Post)
	err := newClient().NewRequest().
		POST().
		Resource("https://jsonplaceholder.typicode.com/posts").
		Body(JSON(&Post{
			UserID: 10,
			Title:  "lorem ipsum",
			Body:   "body",
		})).
		Send().
		DecodeBody(JSON(post)).
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

// Example_PUT provides an example of performing a PUT request.
// See https://jsonplaceholder.typicode.com/ for resource description.
func Example_PUT() {
	post := new(Post)
	err := newClient().NewRequest().
		PUT().
		Resource("https://jsonplaceholder.typicode.com/posts/1").
		Body(JSON(&Post{
			UserID: 10,
			Title:  "lorem ipsum",
			Body:   "body",
		})).
		Send().
		DecodeBody(JSON(post)).
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

// Example_PATCH provides an example of performing a PATCH request.
// See https://jsonplaceholder.typicode.com/ for resource description.
func Example_PATCH() {
	post := new(Post)
	err := newClient().NewRequest().
		PATCH().
		Resource("https://jsonplaceholder.typicode.com/posts/1").
		Body(JSON(&Post{
			UserID: 10,
			Title:  "lorem ipsum",
			Body:   "body",
		})).
		Send().
		DecodeBody(JSON(post)).
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

// Example_DELETE provides an example of performing a DELETE request.
// See https://jsonplaceholder.typicode.com/ for resource description.
func Example_DELETE() {
	response := PlainText("")
	err := newClient().NewRequest().
		DELETE().
		Resource("https://jsonplaceholder.typicode.com/posts/1").
		Send().
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
