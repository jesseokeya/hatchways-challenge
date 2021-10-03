package data

import (
	"context"
	"errors"
	"posts/v1/lib"
	"posts/v1/lib/hatchways"
	"strings"
	"sync"
)

var (
	// errors
	ErrMissingTags = errors.New("Tags parameter is required")

	// Store all posts recieved
	Cache = sync.Map{}

	// SeenTags keeps track of previous tags to avoid duplicated work
	SeenTags = sync.Map{}
)

// GetPosts returns all posts from the cache or retrieves them from the hatchways API
func GetPosts(ctx context.Context, r *hatchways.FilterQueryParams) (*hatchways.PostResponse, error) {
	if r == nil || strings.Trim(r.Tags, " ") == "" {
		return nil, ErrMissingTags
	}

	// Response to return which contains all posts based on the query above
	response := &hatchways.PostResponse{}

	// Split all tags into an array
	tags := strings.Split(r.Tags, ",")

	// Keep track of all posts in the tags array to be done asyncronoulsy after resolving request
	relatedTags := []string{}

	// Removes duplicate when appending the post to the PostResponse
	lookUp := make(map[int64]bool)

	// Iterate through all tags
	for _, tag := range tags {
		keys, exists := SeenTags.Load(tag)
		// If not seen fetch posts from hatchways API and add to cache
		if !exists {
			// Fetch and put in cache for next time
			posts, err := lib.HW.Posts.GetPosts(ctx, &hatchways.PostRequest{Tag: tag})
			if err != nil {
				return nil, err
			}

			if keys == nil {
				SeenTags.Store(tag, make(map[int64]bool))
			}

			// Loop through all posts and add to cache and seen tags for future asyncronous jobs
			for _, post := range posts.Posts {
				Cache.Store(post.ID, post)
				seenTags, ok := SeenTags.Load(tag)
				if ok {
					// Make a deep copy of map to prevent deadlock caused by concurrent access
					newSeenTags := CopyMap(seenTags.(map[int64]bool))
					newSeenTags[post.ID] = true
					SeenTags.Store(tag, newSeenTags)
				}
				response.Posts = append(response.Posts, post)
				_, hasApended := lookUp[post.ID]
				if !hasApended {
					relatedTags = append(relatedTags, post.Tags...)
					lookUp[post.ID] = true
				}
			}
		}

		// If seen load from cache and append to response
		seenTags, ok := SeenTags.Load(tag)
		if ok {
			for k := range seenTags.(map[int64]bool) {
				if result, exists := Cache.Load(k); exists {
					p := result.(*hatchways.Post)
					_, hasApended := lookUp[p.ID]
					if !hasApended {
						response.Posts = append(response.Posts, p)
						lookUp[p.ID] = true
					}
				}
			}
		}
	}

	// Go routine to fetch posts from hatchways API for all tags in the relatedTags array
	go func(t []string) {
		for _, tag := range t {
			_, _ = GetPosts(context.Background(), &hatchways.FilterQueryParams{Tags: tag})
		}
	}(relatedTags)

	// Return and apply sort and direction filters on the posts
	return hatchways.ApplyFilters(response, r), nil
}

// CopyMap makes a deep copy of a map
func CopyMap(m map[int64]bool) map[int64]bool {
	cp := make(map[int64]bool)
	for k, v := range m {
		cp[k] = v
	}
	return cp
}
