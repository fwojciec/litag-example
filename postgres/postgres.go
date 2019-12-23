package postgres

import (
	"context"
	"database/sql"

	"github.com/fwojciec/litag-example/generated/sqlc" // use your own github username
	_ "github.com/lib/pq"                              // required
)

// Repo represents PostgreSQL-backed datalayer functionality.
type Repo struct {
	DB
	Q
}

// DB represents additional database functionality required by DataService.
type DB interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// Q represents all available database queries.
type Q interface {
	// agent queries
	CreateAgent(ctx context.Context, args sqlc.CreateAgentParams) (sqlc.Agent, error)
	DeleteAgent(ctx context.Context, id int64) (sqlc.Agent, error)
	GetAgent(ctx context.Context, id int64) (sqlc.Agent, error)
	ListAgents(ctx context.Context) ([]sqlc.Agent, error)
	UpdateAgent(ctx context.Context, args sqlc.UpdateAgentParams) (sqlc.Agent, error)

	// author queries
	CreateAuthor(ctx context.Context, args sqlc.CreateAuthorParams) (sqlc.Author, error)
	DeleteAuthor(ctx context.Context, id int64) (sqlc.Author, error)
	GetAuthor(ctx context.Context, id int64) (sqlc.Author, error)
	ListAuthors(ctx context.Context) ([]sqlc.Author, error)
	UpdateAuthor(ctx context.Context, args sqlc.UpdateAuthorParams) (sqlc.Author, error)

	// book queries
	CreateBook(ctx context.Context, args sqlc.CreateBookParams) (sqlc.Book, error)
	DeleteBook(ctx context.Context, id int64) (sqlc.Book, error)
	GetBook(ctx context.Context, id int64) (sqlc.Book, error)
	ListBooks(ctx context.Context) ([]sqlc.Book, error)
	SetBookAuthor(ctx context.Context, args sqlc.SetBookAuthorParams) error
	UnsetBookAuthors(ctx context.Context, bookID int64) error
	UpdateBook(ctx context.Context, args sqlc.UpdateBookParams) (sqlc.Book, error)

	// transaction support
	WithTx(tx *sql.Tx) *sqlc.Queries
}

// NewRepo returns a pointer to a new instance of Repo.
func NewRepo(connectString string) (*Repo, error) {
	db, err := sql.Open("postgres", connectString)
	if err != nil {
		return nil, err
	}
	q := sqlc.New(db)
	return &Repo{DB: db, Q: q}, nil
}
