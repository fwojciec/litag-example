package resolvers

import (
	"context"

	"github.com/fwojciec/litag-example/generated/gqlgen" // update username to point to your repo
	"github.com/fwojciec/litag-example/generated/sqlc"   // update username to point to your repo
	"github.com/fwojciec/litag-example/postgres"         // update username to point to your repo
)

// Resolver connects individual resolvers with the datalayer.
type Resolver struct {
	Repo *postgres.Repo
}

// Agent resolver resolves Agent related data.
func (r *Resolver) Agent() gqlgen.AgentResolver {
	return &agentResolver{r}
}

// Author resolver resolves Agent related data.
func (r *Resolver) Author() gqlgen.AuthorResolver {
	return &authorResolver{r}
}

// Book resolver resolves Agent related data.
func (r *Resolver) Book() gqlgen.BookResolver {
	return &bookResolver{r}
}

// Mutation resolver resolves Agent related data.
func (r *Resolver) Mutation() gqlgen.MutationResolver {
	return &mutationResolver{r}
}

// Query resolver resolves Agent related data.
func (r *Resolver) Query() gqlgen.QueryResolver {
	return &queryResolver{r}
}

type agentResolver struct{ *Resolver }

func (r *agentResolver) Authors(ctx context.Context, obj *sqlc.Agent) ([]sqlc.Author, error) {
	return r.Repo.ListAuthorsByAgentID(ctx, obj.ID)
}

type authorResolver struct{ *Resolver }

func (r *authorResolver) Website(ctx context.Context, obj *sqlc.Author) (*string, error) {
	var w string
	if obj.Website.Valid {
		w = obj.Website.String
		return &w, nil
	}
	return nil, nil
}

func (r *authorResolver) Agent(ctx context.Context, obj *sqlc.Author) (*sqlc.Agent, error) {
	agent, err := r.Repo.GetAgent(ctx, obj.AgentID)
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (r *authorResolver) Books(ctx context.Context, obj *sqlc.Author) ([]sqlc.Book, error) {
	return r.Repo.ListBooksByAuthorID(ctx, obj.ID)
}

type bookResolver struct{ *Resolver }

func (r *bookResolver) Authors(ctx context.Context, obj *sqlc.Book) ([]sqlc.Author, error) {
	panic("not implemented")
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateAgent(ctx context.Context, data gqlgen.CreateUpdateAgentInput) (*sqlc.Agent, error) {
	agent, err := r.Repo.CreateAgent(ctx, sqlc.CreateAgentParams{
		Name:  data.Name,
		Email: data.Email,
	})
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (r *mutationResolver) UpdateAgent(ctx context.Context, id int64, data gqlgen.CreateUpdateAgentInput) (*sqlc.Agent, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteAgent(ctx context.Context, id int64) (*sqlc.Agent, error) {
	panic("not implemented")
}
func (r *mutationResolver) CreateAuthor(ctx context.Context, data gqlgen.CreateUpdateAuthorInput) (*sqlc.Author, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateAuthor(ctx context.Context, id int64, data gqlgen.CreateUpdateAuthorInput) (*sqlc.Author, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteAuthor(ctx context.Context, id int64) (*sqlc.Author, error) {
	panic("not implemented")
}
func (r *mutationResolver) CreateBook(ctx context.Context, data gqlgen.CreateUpdateBookInput) (*sqlc.Book, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateBook(ctx context.Context, id int64, data gqlgen.CreateUpdateBookInput) (*sqlc.Book, error) {
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
	return r.Repo.ListAgents(ctx)
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
