package resolvers

//go:generate go run github.com/matryer/moq -out db_mock_test.go -pkg resolvers_test ../postgres DB
//go:generate go run github.com/matryer/moq -out q_mock_test.go -pkg resolvers_test ../postgres Q

import (
	"context"
	"database/sql"

	"github.com/fwojciec/litag-example/generated/gqlgen" // update the username
	"github.com/fwojciec/litag-example/generated/sqlc"   // update the username
	"github.com/fwojciec/litag-example/postgres"         // update the username
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
	return r.Repo.ListAuthorsByBookID(ctx, obj.ID)
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
	agent, err := r.Repo.UpdateAgent(ctx, sqlc.UpdateAgentParams{
		ID:    id,
		Name:  data.Name,
		Email: data.Email,
	})
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (r *mutationResolver) DeleteAgent(ctx context.Context, id int64) (*sqlc.Agent, error) {
	agent, err := r.Repo.DeleteAgent(ctx, id)
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (r *mutationResolver) CreateAuthor(ctx context.Context, data gqlgen.CreateUpdateAuthorInput) (*sqlc.Author, error) {
	author, err := r.Repo.CreateAuthor(ctx, sqlc.CreateAuthorParams{
		Name:    data.Name,
		Website: stringPtrToNullString(data.Website),
		AgentID: data.AgentID,
	})
	if err != nil {
		return nil, err
	}
	return &author, nil
}

func (r *mutationResolver) UpdateAuthor(ctx context.Context, id int64, data gqlgen.CreateUpdateAuthorInput) (*sqlc.Author, error) {
	author, err := r.Repo.UpdateAuthor(ctx, sqlc.UpdateAuthorParams{
		ID:      id,
		Name:    data.Name,
		Website: stringPtrToNullString(data.Website),
		AgentID: data.AgentID,
	})
	if err != nil {
		return nil, err
	}
	return &author, nil
}

func (r *mutationResolver) DeleteAuthor(ctx context.Context, id int64) (*sqlc.Author, error) {
	author, err := r.Repo.DeleteAuthor(ctx, id)
	if err != nil {
		return nil, err
	}
	return &author, nil
}

func (r *mutationResolver) CreateBook(ctx context.Context, data gqlgen.CreateUpdateBookInput) (*sqlc.Book, error) {
	return r.Repo.CreateBook(ctx, sqlc.CreateBookParams{
		Title:       data.Title,
		Description: data.Description,
		Cover:       data.Cover,
	}, data.AuthorIDs)
}

func (r *mutationResolver) UpdateBook(ctx context.Context, id int64, data gqlgen.CreateUpdateBookInput) (*sqlc.Book, error) {
	return r.Repo.UpdateBook(ctx, sqlc.UpdateBookParams{
		ID:          id,
		Title:       data.Title,
		Description: data.Description,
		Cover:       data.Cover,
	}, data.AuthorIDs)
}

func (r *mutationResolver) DeleteBook(ctx context.Context, id int64) (*sqlc.Book, error) {
	// BookAuthors associations will cascade automatically.
	book, err := r.Repo.DeleteBook(ctx, id)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Agent(ctx context.Context, id int64) (*sqlc.Agent, error) {
	agent, err := r.Repo.GetAgent(ctx, id)
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (r *queryResolver) Agents(ctx context.Context) ([]sqlc.Agent, error) {
	return r.Repo.ListAgents(ctx)
}

func (r *queryResolver) Author(ctx context.Context, id int64) (*sqlc.Author, error) {
	author, err := r.Repo.GetAuthor(ctx, id)
	if err != nil {
		return nil, err
	}
	return &author, nil
}

func (r *queryResolver) Authors(ctx context.Context) ([]sqlc.Author, error) {
	return r.Repo.ListAuthors(ctx)
}

func (r *queryResolver) Book(ctx context.Context, id int64) (*sqlc.Book, error) {
	book, err := r.Repo.GetBook(ctx, id)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *queryResolver) Books(ctx context.Context) ([]sqlc.Book, error) {
	return r.Repo.ListBooks(ctx)
}

func stringPtrToNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{}
}
