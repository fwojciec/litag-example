package postgres

import (
	"context"
	"database/sql"

	"github.com/fwojciec/litag-example/generated/sqlc" // use your own github username
	_ "github.com/lib/pq"                              // required
)

// Repo represents PostgreSQL-backed datalayer functionality.
type Repo struct {
	Q
	TxQ
}

// NewRepo returns a pointer to a new instance of Repo.
func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		Q:   sqlc.New(db),
		TxQ: &txqService{db},
	}
}

// Q represents database queries.
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
	ListAuthorsByAgentID(ctx context.Context, agentID int64) ([]sqlc.Author, error)
	ListAuthorsByBookID(ctx context.Context, bookID int64) ([]sqlc.Author, error)

	// book queries
	DeleteBook(ctx context.Context, id int64) (sqlc.Book, error)
	GetBook(ctx context.Context, id int64) (sqlc.Book, error)
	ListBooks(ctx context.Context) ([]sqlc.Book, error)
	ListBooksByAuthorID(ctx context.Context, authorID int64) ([]sqlc.Book, error)
}

// TxQ represents queries performed using a transaction.
type TxQ interface {
	CreateBook(ctx context.Context, bookArgs sqlc.CreateBookParams, authorIDs []int64) (*sqlc.Book, error)
	UpdateBook(ctx context.Context, bookArgs sqlc.UpdateBookParams, authorIDs []int64) (*sqlc.Book, error)
}

type txqService struct {
	db *sql.DB
}

func (txq *txqService) CreateBook(ctx context.Context, bookArgs sqlc.CreateBookParams, authorIDs []int64) (*sqlc.Book, error) {
	tx, q, err := txq.init(ctx)
	if err != nil {
		return nil, err
	}
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

func (txq *txqService) UpdateBook(ctx context.Context, bookArgs sqlc.UpdateBookParams, authorIDs []int64) (*sqlc.Book, error) {
	tx, q, err := txq.init(ctx)
	if err != nil {
		return nil, err
	}
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

func (txq *txqService) init(ctx context.Context) (*sql.Tx, *sqlc.Queries, error) {
	tx, err := txq.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	q := sqlc.New(tx)
	return tx, q, nil
}
