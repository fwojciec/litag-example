package gqlgen

import (
	"context"

	"github.com/fwojciec/litag-example/generated/sqlc"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Agent() AgentResolver {
	return &agentResolver{r}
}
func (r *Resolver) Author() AuthorResolver {
	return &authorResolver{r}
}
func (r *Resolver) Book() BookResolver {
	return &bookResolver{r}
}
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type agentResolver struct{ *Resolver }

func (r *agentResolver) Authors(ctx context.Context, obj *sqlc.Agent) ([]sqlc.Author, error) {
	panic("not implemented")
}

type authorResolver struct{ *Resolver }

func (r *authorResolver) Website(ctx context.Context, obj *sqlc.Author) (*string, error) {
	panic("not implemented")
}
func (r *authorResolver) Agent(ctx context.Context, obj *sqlc.Author) (*sqlc.Agent, error) {
	panic("not implemented")
}
func (r *authorResolver) Books(ctx context.Context, obj *sqlc.Author) ([]sqlc.Book, error) {
	panic("not implemented")
}

type bookResolver struct{ *Resolver }

func (r *bookResolver) Authors(ctx context.Context, obj *sqlc.Book) ([]sqlc.Author, error) {
	panic("not implemented")
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateAgent(ctx context.Context, data CreateUpdateAgentInput) (*sqlc.Agent, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateAgent(ctx context.Context, id int64, data CreateUpdateAgentInput) (*sqlc.Agent, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteAgent(ctx context.Context, id int64) (*sqlc.Agent, error) {
	panic("not implemented")
}
func (r *mutationResolver) CreateAuthor(ctx context.Context, data CreateUpdateAuthorInput) (*sqlc.Author, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateAuthor(ctx context.Context, id int64, data CreateUpdateAuthorInput) (*sqlc.Author, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteAuthor(ctx context.Context, id int64) (*sqlc.Author, error) {
	panic("not implemented")
}
func (r *mutationResolver) CreateBook(ctx context.Context, data CreateUpdateBookInput) (*sqlc.Book, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateBook(ctx context.Context, id int64, data CreateUpdateBookInput) (*sqlc.Book, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteBook(ctx context.Context, id int64) (*sqlc.Book, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Agent(ctx context.Context, id int64) (*sqlc.Agent, error) {
	panic("not implemented")
}
func (r *queryResolver) Agents(ctx context.Context) ([]sqlc.Agent, error) {
	panic("not implemented")
}
func (r *queryResolver) Author(ctx context.Context, id int64) (*sqlc.Author, error) {
	panic("not implemented")
}
func (r *queryResolver) Authors(ctx context.Context) ([]sqlc.Author, error) {
	panic("not implemented")
}
func (r *queryResolver) Book(ctx context.Context, id int64) (*sqlc.Book, error) {
	panic("not implemented")
}
func (r *queryResolver) Books(ctx context.Context) ([]sqlc.Book, error) {
	panic("not implemented")
}
