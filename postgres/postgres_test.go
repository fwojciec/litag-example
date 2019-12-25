package postgres_test

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/fwojciec/litag-example/generated/sqlc"
	"github.com/fwojciec/litag-example/postgres"
)

/*
	1) Test create queries
	2) Test list queries
	3) Test update queries
	4) Test get queries (on updated queries)
	5) Test delete queries
*/

func TestQueries(t *testing.T) {

	var (
		testAgent1 = sqlc.Agent{
			Name:  "test agent name 1",
			Email: "agent1@test.com",
		}
		testAgent2 = sqlc.Agent{
			Name:  "test agent name 2",
			Email: "agent2@test.com",
		}
		testAgentUpdated = sqlc.Agent{
			// agent 2
			Name:  "test agent name updated",
			Email: "updated@test.com",
		}
		testAuthor1 = sqlc.Author{
			// agent 1
			Name:    "test author name 1",
			Website: sql.NullString{String: "https://author1.com", Valid: true},
		}
		testAuthor2 = sqlc.Author{
			// agent 2
			Name:    "test author name 2",
			Website: sql.NullString{},
		}
		testAuthorUpdated = sqlc.Author{
			// author 2
			// agent 1
			Name:    "test author name updated",
			Website: sql.NullString{String: "https://authorupdated.com", Valid: true},
		}
		testBook1 = sqlc.Book{
			// authors 1 and 2
			Title:       "book 1 title",
			Description: "book 1 description",
			Cover:       "book1.jpg",
		}
		testBook2 = sqlc.Book{
			// author 2
			Title:       "book 2 title",
			Description: "book 2 description",
			Cover:       "book2.jpg",
		}
		testBookUpdated = sqlc.Book{
			// book 1
			// author 1
			Title:       "book title updated",
			Description: "book description updated",
			Cover:       "bookupdated.jpg",
		}
	)

	runner(t, func(ctx context.Context, r *postgres.Repo, t *testing.T) {
		t.Run("Create queries", func(t *testing.T) {
			t.Run("CreateAgent 1", func(t *testing.T) {
				a, err := r.CreateAgent(ctx, sqlc.CreateAgentParams{
					Name:  testAgent1.Name,
					Email: testAgent1.Email,
				})
				if err != nil {
					t.Fatalf("failed to create agent: %s", err)
				}
				testAgent1.ID = a.ID
				testAuthor1.AgentID = a.ID
				if !reflect.DeepEqual(testAgent1, a) {
					t.Errorf("expected %v, received %v", testAgent1, a)
				}
			})

			t.Run("CreateAgent 2", func(t *testing.T) {
				a, err := r.CreateAgent(ctx, sqlc.CreateAgentParams{
					Name:  testAgent2.Name,
					Email: testAgent2.Email,
				})
				if err != nil {
					t.Fatalf("failed to create agent: %s", err)
				}
				testAgent2.ID = a.ID
				testAuthor2.AgentID = a.ID
			})

			t.Run("CreateAuthor 1", func(t *testing.T) {
				a, err := r.CreateAuthor(ctx, sqlc.CreateAuthorParams{
					Name:    testAuthor1.Name,
					Website: testAuthor1.Website,
					AgentID: testAuthor1.AgentID,
				})
				if err != nil {
					t.Fatalf("failed to create author: %s", err)
				}
				testAuthor1.ID = a.ID
				if !reflect.DeepEqual(testAuthor1, a) {
					t.Errorf("expected %v, received %v", testAuthor1, a)
				}
			})

			t.Run("CreateAuthor 2", func(t *testing.T) {
				a, err := r.CreateAuthor(ctx, sqlc.CreateAuthorParams{
					Name:    testAuthor2.Name,
					Website: testAuthor2.Website,
					AgentID: testAuthor2.AgentID,
				})
				if err != nil {
					t.Fatalf("failed to create author: %s", err)
				}
				testAuthor2.ID = a.ID
			})

			t.Run("CreateBook 1", func(t *testing.T) {
				b, err := r.CreateBook(ctx, sqlc.CreateBookParams{
					Title:       testBook1.Title,
					Description: testBook1.Description,
					Cover:       testBook1.Cover,
				}, []int64{testAuthor1.ID, testAuthor2.ID})
				if err != nil {
					t.Fatalf("failed to create book: %s", err)
				}
				testBook1.ID = b.ID
				if !reflect.DeepEqual(&testBook1, b) {
					t.Errorf("expected %v, received %v", testBook1, b)
				}
			})

			t.Run("CreateBook 2", func(t *testing.T) {
				b, err := r.CreateBook(ctx, sqlc.CreateBookParams{
					Title:       testBook2.Title,
					Description: testBook2.Description,
					Cover:       testBook2.Cover,
				}, []int64{testAuthor2.ID})
				if err != nil {
					t.Fatalf("failed to create book: %s", err)
				}
				testBook2.ID = b.ID
			})
		})

		t.Run("List queries", func(t *testing.T) {
			t.Run("ListAgents", func(t *testing.T) {
				l, err := r.ListAgents(ctx)
				if err != nil {
					t.Fatalf("failed to list agents: %s", err)
				}
				exp := []sqlc.Agent{testAgent1, testAgent2}
				if !reflect.DeepEqual(exp, l) {
					t.Errorf("expected %v, received %v", exp, l)
				}
			})

			t.Run("ListAuthors", func(t *testing.T) {
				l, err := r.ListAuthors(ctx)
				if err != nil {
					t.Fatalf("failed to list authors: %s", err)
				}
				exp := []sqlc.Author{testAuthor1, testAuthor2}
				if !reflect.DeepEqual(exp, l) {
					t.Errorf("expected %v, received %v", exp, l)
				}
			})

			t.Run("ListAgentsByAuthorID", func(t *testing.T) {
				l, err := r.ListAuthorsByAgentID(ctx, testAgent1.ID)
				if err != nil {
					t.Fatalf("failed to list authors by agent id: %s", err)
				}
				exp := []sqlc.Author{testAuthor1}
				if !reflect.DeepEqual(exp, l) {
					t.Errorf("expected %v, received %v", exp, l)
				}
			})

			t.Run("ListBooks", func(t *testing.T) {
				l, err := r.ListBooks(ctx)
				if err != nil {
					t.Fatalf("failed to list books: %s", err)
				}
				exp := []sqlc.Book{testBook1, testBook2}
				if !reflect.DeepEqual(exp, l) {
					t.Errorf("expected %v, received %v", exp, l)
				}
			})

			t.Run("ListAuthorsByBookIDs", func(t *testing.T) {
				l, err := r.ListAuthorsByBookID(ctx, testBook1.ID)
				if err != nil {
					t.Fatalf("failed to get authors by book id: %s", err)
				}
				exp := []sqlc.Author{testAuthor1, testAuthor2}
				if !reflect.DeepEqual(exp, l) {
					t.Errorf("expected %v, received %v", exp, l)
				}
			})

			t.Run("ListBooksByAuthorIDs", func(t *testing.T) {
				l, err := r.ListBooksByAuthorID(ctx, testAuthor2.ID)
				if err != nil {
					t.Fatalf("failed to list books by author ids: %s", err)
				}
				exp := []sqlc.Book{testBook1, testBook2}
				if !reflect.DeepEqual(exp, l) {
					t.Errorf("expected %v, received %v", exp, l)
				}
			})
		})

		t.Run("Update queries", func(t *testing.T) {
			t.Run("UpdateAgent", func(t *testing.T) {
				a, err := r.UpdateAgent(ctx, sqlc.UpdateAgentParams{
					ID:    testAgent2.ID,
					Name:  testAgentUpdated.Name,
					Email: testAgentUpdated.Email,
				})
				if err != nil {
					t.Fatalf("failed to update agent: %s", err)
				}
				testAgentUpdated.ID = a.ID
				if !reflect.DeepEqual(testAgentUpdated, a) {
					t.Errorf("expected %v, received %v", testAgentUpdated, a)
				}
			})

			t.Run("UpdateAuthor", func(t *testing.T) {
				a, err := r.UpdateAuthor(ctx, sqlc.UpdateAuthorParams{
					ID:      testAuthor2.ID,
					Name:    testAuthorUpdated.Name,
					Website: testAuthorUpdated.Website,
					AgentID: testAgent1.ID,
				})
				if err != nil {
					t.Fatal("failed to update agent")
				}
				testAuthorUpdated.ID = a.ID
				testAuthorUpdated.AgentID = testAgent1.ID
				if !reflect.DeepEqual(testAuthorUpdated, a) {
					t.Errorf("expected %v, received %v", testAuthorUpdated, a)
				}
			})

			t.Run("UpdateBook", func(t *testing.T) {
				b, err := r.UpdateBook(ctx, sqlc.UpdateBookParams{
					ID:          testBook1.ID,
					Title:       testBookUpdated.Title,
					Description: testBookUpdated.Description,
					Cover:       testBookUpdated.Cover,
				}, []int64{testAuthor1.ID})
				if err != nil {
					t.Fatalf("failed to update book: %s", err)
				}
				testBookUpdated.ID = b.ID
				if !reflect.DeepEqual(&testBookUpdated, b) {
					t.Errorf("expected %v, received %v", testBookUpdated, b)
				}
				l, err := r.ListAuthorsByBookID(ctx, testBookUpdated.ID)
				if err != nil {
					t.Fatalf("failed to list authors by book id: %s", err)
				}
				exp := []sqlc.Author{testAuthor1}
				if !reflect.DeepEqual(exp, l) {
					t.Errorf("expected %v, received %v", exp, l)
				}
			})
		})

		t.Run("Get queries", func(t *testing.T) {
			t.Run("GetAgent", func(t *testing.T) {
				a, err := r.GetAgent(ctx, testAgentUpdated.ID)
				if err != nil {
					t.Fatalf("failed to get agent: %s", err)
				}
				if !reflect.DeepEqual(testAgentUpdated, a) {
					t.Errorf("expected %v, received %v", testAgentUpdated, a)
				}
			})

			t.Run("GetAuthor", func(t *testing.T) {
				a, err := r.GetAuthor(ctx, testAuthorUpdated.ID)
				if err != nil {
					t.Fatalf("failed to get author: %s", err)
				}
				if !reflect.DeepEqual(testAuthorUpdated, a) {
					t.Errorf("expected %v, received %v", testAuthorUpdated, a)
				}
			})

			t.Run("GetBook", func(t *testing.T) {
				b, err := r.GetBook(ctx, testBookUpdated.ID)
				if err != nil {
					t.Fatalf("failed to get book: %s", err)
				}
				if !reflect.DeepEqual(testBookUpdated, b) {
					t.Errorf("expected %v, received %v", testBookUpdated, b)
				}
			})
		})

		t.Run("Delete queries", func(t *testing.T) {
			t.Run("DeleteAgent ForeignKey Constraint", func(t *testing.T) {
				_, err := r.DeleteAgent(ctx, testAgent1.ID)
				expError := `pq: update or delete on table "agents" violates foreign key constraint "authors_agent_id_fkey" on table "authors"`
				if err.Error() != expError {
					t.Fatalf("expected %s message, received %s", expError, err)
				}
			})

			t.Run("DeleteBook", func(t *testing.T) {
				b, err := r.DeleteBook(ctx, testBook2.ID)
				if err != nil {
					t.Fatalf("failed to delete book: %s", err)
				}
				if !reflect.DeepEqual(testBook2, b) {
					t.Errorf("expected %v, received %v", testBook2, b)
				}
				l1, err := r.ListAuthorsByBookID(ctx, testBook2.ID)
				if err != nil {
					t.Fatalf("failed to list authors by book id: %s", err)
				}
				if len(l1) != 0 {
					t.Errorf("expected length of 0, received %d", len(l1))
				}
				l2, err := r.ListBooks(ctx)
				if err != nil {
					t.Fatalf("failed to list books: %s", err)
				}
				if len(l2) != 1 {
					t.Errorf("expected length of 1, received %d", len(l2))
				}
			})

			t.Run("DeleteAuthor", func(t *testing.T) {
				a, err := r.DeleteAuthor(ctx, testAuthor1.ID)
				if err != nil {
					t.Fatalf("failed to delete author: %s", err)
				}
				if !reflect.DeepEqual(testAuthor1, a) {
					t.Errorf("expected %v, received %v", testAuthor1, a)
				}
				l, err := r.ListAuthors(ctx)
				if err != nil {
					t.Fatalf("failed to list books: %s", err)
				}
				if len(l) != 1 {
					t.Errorf("expected length of 1, received %d", len(l))
				}
			})

			t.Run("DeleteAgent", func(t *testing.T) {
				a, err := r.DeleteAgent(ctx, testAgentUpdated.ID)
				if err != nil {
					t.Fatalf("failed to delete agent: %s", err)
				}
				if !reflect.DeepEqual(testAgentUpdated, a) {
					t.Errorf("expected %v, received %v", testAgentUpdated, a)
				}
				l, err := r.ListAgents(ctx)
				if err != nil {
					t.Fatalf("failed to list agents after delete: %s", err)
				}
				if len(l) != 1 {
					t.Errorf("expected length of 1, received %d", len(l))
				}
			})
		})
	})
}

func runner(t *testing.T, test func(context.Context, *postgres.Repo, *testing.T)) {
	ctx := context.Background()

	db, err := sql.Open("postgres", "dbname=test_db sslmode=disable")
	if err != nil {
		t.Fatalf("failed to connect to the db: %s\n", err)
	}

	repo := postgres.NewRepo(db)

	// create and drop schema (defer)
	defer func() {
		if err := dropSchema(ctx, db); err != nil {
			t.Fatalf("failed to drop queue schema: %s\n", err)
		}
	}()
	if err := createSchema(ctx, db); err != nil {
		t.Fatalf("failed to create queue schema: %s\n", err)
	}

	test(ctx, repo, t)
}

func createSchema(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS agents (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL
		);
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS authors (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			website TEXT,
			agent_id BIGINT NOT NULL,
			FOREIGN KEY (agent_id) REFERENCES agents(id) 
		);
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS books (
			id BIGSERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			cover TEXT NOT NULL
		);
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS book_authors (
			id BIGSERIAL PRIMARY KEY,
			book_id BIGINT NOT NULL,
			author_id BIGINT NOT NULL,
			FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
			FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE CASCADE,
			UNIQUE (book_id,author_id)
		);
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func dropSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		DROP TABLE IF EXISTS book_authors, books, authors, agents;
	`)
	if err != nil {
		return err
	}
	return nil
}
