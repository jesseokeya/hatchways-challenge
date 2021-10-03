package hatchways

import (
	"context"
	"sort"
	"strings"
)

// PostService is the service for the hatchways post API
type PostService service

// Post is the response structure for Hatchways post
type Post struct {
	ID         int64    `json:"id"`
	Author     string   `json:"author"`
	AuthorID   int64    `json:"authorId"`
	Likes      int64    `json:"likes"`
	Popularity float64  `json:"popularity"`
	Reads      int64    `json:"reads"`
	Tags       []string `json:"tags"`
}

// PostResponse is the request structure for Hatchways post api response
type PostResponse struct {
	Posts []*Post `json:"posts"`
}

// PostRequest contains the request query parameters for the post api
type PostRequest struct {
	Tag string `url:"tag"`
}

// FilterQueryParams is the request structure for client facing the filter query parameters
type FilterQueryParams struct {
	Tags      string
	SortBy    string
	Direction string
}

// GetPosts returns the posts for the given tag
func (c *PostService) GetPosts(ctx context.Context, request *PostRequest) (*PostResponse, error) {
	req, err := c.client.NewRequest("GET", "assessment/blog/posts", nil)
	if err != nil {
		return nil, err
	}

	query, err := ToURLQuery(request)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query

	var response PostResponse
	_, err = c.client.Do(ctx, req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// EnumSortBy is the enum for the sort by
type EnumSortBy string

// EnumDirection is the enum for the direction
type EnumDirection string

var (
	EnumSortByID         EnumSortBy = "id"
	EnumSortByReads      EnumSortBy = "reads"
	EnumSortByLikes      EnumSortBy = "likes"
	EnumSortByPopularity EnumSortBy = "popularity"
)

var (
	EnumDirectionAsc  EnumDirection = "asc"
	EnumDirectionDesc EnumDirection = "desc"
)

// PostSorter is the struct for sorting the posts
type PostSorter struct {
	Posts     []*Post
	By        EnumSortBy
	Direction EnumDirection
}

// Len retieves the length of the posts. And is a function of the go sort interface
func (a PostSorter) Len() int { return len(a.Posts) }

// Swap swaps 2 posts in place. And is a function of the go sort interface
func (a PostSorter) Swap(i, j int) { a.Posts[i], a.Posts[j] = a.Posts[j], a.Posts[i] }

// Less determines the smaller post based on a field . And is a function of the go sort interface
func (a PostSorter) Less(i, j int) bool {
	if a.By == EnumSortByReads {
		if a.Direction == "asc" {
			return a.Posts[i].Reads < a.Posts[j].Reads
		}
		return a.Posts[i].Reads > a.Posts[j].Reads
	}
	if a.By == EnumSortByLikes {
		if a.Direction == EnumDirectionAsc {
			return a.Posts[i].Likes < a.Posts[j].Likes
		}
		return a.Posts[i].Likes > a.Posts[j].Likes
	}
	if a.By == EnumSortByPopularity {
		if a.Direction == EnumDirectionAsc {
			return a.Posts[i].Popularity < a.Posts[j].Popularity
		}
		return a.Posts[i].Popularity > a.Posts[j].Popularity
	}

	if a.Direction == EnumDirectionAsc {
		return a.Posts[i].ID < a.Posts[j].ID
	}
	return a.Posts[i].ID > a.Posts[j].ID
}

// ApplyFilters	applies the filters to the posts to be returned
func ApplyFilters(p *PostResponse, r *FilterQueryParams) *PostResponse {
	direction := strings.Trim(r.Direction, " ")
	if direction == "" {
		// asc is default if no direction is specified
		direction = "asc"
	}

	// Sort the posts based on the sort by and direction
	sort.Sort(PostSorter{
		Posts:     p.Posts,
		By:        EnumSortBy(r.SortBy),
		Direction: EnumDirection(direction),
	})

	return p
}
