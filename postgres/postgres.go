package postgres

import (
	"context"
	"database/sql"

	"github.com/fwojciec/litag-example/generated/sqlc" // use your own github username
	_ "github.com/lib/pq"                              // required
)

// Repo represents PostgreSQL-backed datalayer functionality.
type Repo struct {
	Querent
	TxQuerent
}

// NewRepo returns a new instance of Repo.
func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		Querent:   sqlc.New(db),
		TxQuerent: &txQuerentService{db},
	}
}

// Querent represents database query methods.
type Querent interface {
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
	ListAuthorsByAgentID(ctx context.Context, agentID int64) ([]sqlc.Author, error)
	ListAuthorsByBookID(ctx context.Context, bookID int64) ([]sqlc.Author, error)

	// book queries
	DeleteBook(ctx context.Context, id int64) (sqlc.Book, error)
	GetBook(ctx context.Context, id int64) (sqlc.Book, error)
	ListBooks(ctx context.Context) ([]sqlc.Book, error)
	ListBooksByAuthorID(ctx context.Context, authorID int64) ([]sqlc.Book, error)
}

// TxQuerent represents database query methods performed using a transaction.
type TxQuerent interface {
	CreateBook(
		ctx context.Context,
		bookArgs sqlc.CreateBookParams,
		authorIDs []int64,
	) (*sqlc.Book, error)
	UpdateBook(
		ctx context.Context,
		bookArgs sqlc.UpdateBookParams,
		authorIDs []int64,
	) (*sqlc.Book, error)
}

type txQuerentService struct {
	db *sql.DB
}

func (txq *txQuerentService) CreateBook(ctx context.Context, bookArgs sqlc.CreateBookParams, authorIDs []int64) (*sqlc.Book, error) {
	// begin the transaction
	tx, err := txq.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	q := sqlc.New(tx)
	book, err := q.CreateBook(ctx, bookArgs)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, authorID := range authorIDs {
		err := q.SetBookAuthor(ctx, sqlc.SetBookAuthorParams{
			BookID:   book.ID,
			AuthorID: authorID,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &book, nil
}

func (txq *txQuerentService) UpdateBook(ctx context.Context, bookArgs sqlc.UpdateBookParams, authorIDs []int64) (*sqlc.Book, error) {
	tx, err := txq.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	q := sqlc.New(tx)
	book, err := q.UpdateBook(ctx, bookArgs)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = q.UnsetBookAuthors(ctx, book.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, authorID := range authorIDs {
		err := q.SetBookAuthor(ctx, sqlc.SetBookAuthorParams{
			BookID:   book.ID,
			AuthorID: authorID,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &book, nil
}
