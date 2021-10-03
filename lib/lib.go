package lib

import (
	"posts/v1/lib/hatchways"
	"sync"
)

type Post struct {
	*hatchways.Client
}

var (
	Once sync.Once

	// Hatchways is a pointer to the global instance of the hatchways client
	HW *Post
)

// SetupHatchways sets up the hatchways client and returns a pointer to the global instance
// Ensures only one instance of the hatchways api client is created (Singleton)
func SetupHatchways() *Post {
	Once.Do(func() {
		client, _ := hatchways.NewClient(nil)
		HW = &Post{client}
	})
	return HW
}
