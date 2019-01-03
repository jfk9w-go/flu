# flu

**flu** is a fluent net/http client wrapper. It provides
a developer with a convenience suite for writing HTTP clients.

#### Example
```go
package main

import (
	"fmt"
	
	"github.com/jfk9w-go/flu"
)

type Post struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func main() {
	// Create a response container.
	post := new(Post)
	// Create a client and execute a GET request.
	// Unmarshal response body from JSON into the post.
	err := flu.NewClient(nil).NewRequest().
        Endpoint("https://jsonplaceholder.typicode.com/posts/1").
        Get().
        ReadBody(flu.JSON(post)).
        Error
	
	// Check the error.
	if err != nil {
		panic(err)
	}
	
	// Print out the response.
	fmt.Printf("Response: %+v\n", post)
}
```

More examples [here](https://github.com/jfk9w-go/flu/blob/master/client_test.go).