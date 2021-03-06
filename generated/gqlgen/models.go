// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gqlgen

type CreateUpdateAgentInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUpdateAuthorInput struct {
	Name    string  `json:"name"`
	Website *string `json:"website"`
	AgentID int64   `json:"agent_id"`
}

type CreateUpdateBookInput struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Cover       string  `json:"cover"`
	AuthorIDs   []int64 `json:"authorIDs"`
}
