package postgres_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/fwojciec/litag-example/generated/sqlc"
	"github.com/fwojciec/litag-example/postgres"
)

func TestQueries(t *testing.T) {

	var (
		testName1        = "test name 1"
		testName2        = "test name 2"
		testName3        = "test name 3"
		testEmail1       = "test1@email.com"
		testEmail2       = "test2@email.com"
		testEmail3       = "test3@email.com"
		testWebsite1     = "https://test.com"
		testTitle1       = "book title 1"
		testTitle2       = "book title 2"
		testTitle3       = "book title 3"
		testDescription1 = "description 1"
		testDescription2 = "description 2"
		testDescription3 = "description 3"
		testCover1       = "cover1.jpg"
		testCover2       = "cover2.jpg"
		testCover3       = "cover3.jpg"
		agentID1         int64
		agentID2         int64
		authorID1        int64
		authorID2        int64
		bookID1          int64
		bookID2          int64
	)

	runner(t, func(ctx context.Context, q postgres.Q, t *testing.T) {
		t.Run("CreateAgent", func(t *testing.T) {
			// create agent 1
			a1, err := q.CreateAgent(ctx, sqlc.CreateAgentParams{
				Name:  testName1,
				Email: testEmail1,
			})
			if err != nil {
				t.Fatalf("failed to create agent: %s", err)
			}
			if a1.Name != testName1 {
				t.Errorf("expected %q, received %q", testName1, a1.Name)
			}
			if a1.Email != testEmail1 {
				t.Errorf("expected %q, received %q", testEmail1, a1.Email)
			}
			agentID1 = a1.ID

			// create agent 2
			a2, err := q.CreateAgent(ctx, sqlc.CreateAgentParams{
				Name:  testName2,
				Email: testEmail2,
			})
			if err != nil {
				t.Fatalf("failed to create agent: %s", err)
			}
			agentID2 = a2.ID
		})

		t.Run("ListAgents", func(t *testing.T) {
			l, err := q.ListAgents(ctx)
			if err != nil {
				t.Fatalf("failed to list agents: %s", err)
			}
			if len(l) != 2 {
				t.Errorf("expected length of 2, received %d", len(l))
			}
			if l[0].ID != agentID1 {
				t.Errorf("expected %d, received %d", agentID1, l[0].ID)
			}
			if l[0].Name != testName1 {
				t.Errorf("expected %q, received %q", testName1, l[0].Name)
			}
			if l[0].Email != testEmail1 {
				t.Errorf("expected %q, received %q", testEmail1, l[0].Email)
			}
			if l[1].ID != agentID2 {
				t.Errorf("expected %d, received %d", agentID2, l[1].ID)
			}
		})

		t.Run("UpdateAgent", func(t *testing.T) {
			a, err := q.UpdateAgent(ctx, sqlc.UpdateAgentParams{
				ID:    agentID2,
				Name:  testName3,
				Email: testEmail3,
			})
			if err != nil {
				t.Fatalf("failed to update agent: %s", err)
			}
			if a.ID != agentID2 {
				t.Errorf("expected %d, received %d", agentID2, a.ID)
			}
			if a.Name != testName3 {
				t.Errorf("expected %q, received %q", testName3, a.Name)
			}
			if a.Email != testEmail3 {
				t.Errorf("expected %q, received %q", testEmail3, a.Email)
			}
		})

		t.Run("GetAgent", func(t *testing.T) {
			g, err := q.GetAgent(ctx, agentID2)
			if err != nil {
				t.Fatalf("failed to get agent: %s", err)
			}
			// should have the data updated in UpdateAgent above
			if g.ID != agentID2 {
				t.Errorf("expected %d, received %d", agentID2, g.ID)
			}
			if g.Name != testName3 {
				t.Errorf("expected %q, received %q", testName3, g.Name)
			}
			if g.Email != testEmail3 {
				t.Errorf("expected %q, received %q", testEmail3, g.Email)
			}
		})

		t.Run("CreateAuthor", func(t *testing.T) {
			a1, err := q.CreateAuthor(ctx, sqlc.CreateAuthorParams{
				Name:    testName1,
				Website: sql.NullString{String: testWebsite1, Valid: true},
				AgentID: agentID1,
			})
			if err != nil {
				t.Fatalf("failed to create author: %s", err)
			}
			if a1.Name != testName1 {
				t.Errorf("expected %q, received %q", testName1, a1.Name)
			}
			if !a1.Website.Valid {
				t.Errorf("expected true, received %v", a1.Website.Valid)
			}
			if a1.Website.String != testWebsite1 {
				t.Errorf("expected %q, received %q", testWebsite1, a1.Website.String)
			}
			if a1.AgentID != agentID1 {
				t.Errorf("expected %d, received %d", agentID1, a1.AgentID)
			}
			authorID1 = a1.ID
			a2, err := q.CreateAuthor(ctx, sqlc.CreateAuthorParams{
				Name:    testName2,
				Website: sql.NullString{},
				AgentID: agentID2,
			})
			if err != nil {
				t.Fatalf("failed to create author: %s", err)
			}
			authorID2 = a2.ID
		})

		t.Run("ListAuthors", func(t *testing.T) {
			l, err := q.ListAuthors(ctx)
			if err != nil {
				t.Fatalf("failed to list authors: %s", err)
			}
			if len(l) != 2 {
				t.Errorf("expected length of 2, received %d", len(l))
			}
			if l[0].ID != authorID1 {
				t.Errorf("expected %d, received %d", authorID1, l[0].ID)
			}
			if l[0].Name != testName1 {
				t.Errorf("expected %q, received %q", testName1, l[0].Name)
			}
			if l[0].Website.String != testWebsite1 {
				t.Errorf("expected %q, received %q", testWebsite1, l[0].Website.String)
			}
			if l[1].ID != authorID2 {
				t.Errorf("expected %d, received %d", authorID2, l[1].ID)
			}
		})

		t.Run("ListAgentsByAuthorID", func(t *testing.T) {
			l, err := q.ListAuthorsByAgentID(ctx, authorID1)
			if err != nil {
				t.Fatalf("failed to list agents by author id: %s", err)
			}
			if len(l) != 1 {
				t.Errorf("expected length of 1, received %d", len(l))
			}
			if l[0].ID != agentID1 {
				t.Errorf("expected %d, received %d", agentID1, l[0].ID)
			}
		})

		t.Run("UpdateAuthor", func(t *testing.T) {
			a, err := q.UpdateAuthor(ctx, sqlc.UpdateAuthorParams{
				ID:      authorID1,
				Name:    testName3,
				Website: sql.NullString{},
				AgentID: agentID2,
			})
			if err != nil {
				t.Fatal("failed to update agent")
			}
			if a.ID != authorID1 {
				t.Errorf("expected %d, received %d", authorID1, a.ID)
			}
			if a.Name != testName3 {
				t.Errorf("expected %q, received %q", testName3, a.Name)
			}
			if a.Website.Valid {
				t.Errorf("expected false, received %v", a.Website.Valid)
			}
			if a.AgentID != agentID2 {
				t.Errorf("expected %d, received %d", agentID2, a.AgentID)
			}
		})

		t.Run("GetAuthor", func(t *testing.T) {
			g, err := q.GetAuthor(ctx, authorID1)
			if err != nil {
				t.Fatalf("failed to get author: %s", err)
			}
			// should have the data updated in UpdateAuthor above
			if g.ID != authorID1 {
				t.Errorf("expected %d, received %d", authorID1, g.ID)
			}
			if g.Name != testName3 {
				t.Errorf("expected %q, received %q", testName3, g.Name)
			}
			if g.Website.Valid {
				t.Errorf("expected false, received %v", g.Website.Valid)
			}
			if g.AgentID != agentID2 {
				t.Errorf("expected %d, received %d", agentID2, g.AgentID)
			}
		})

		t.Run("CreateBook", func(t *testing.T) {
			b1, err := q.CreateBook(ctx, sqlc.CreateBookParams{
				Title:       testTitle1,
				Description: testDescription1,
				Cover:       testCover1,
			})
			if err != nil {
				t.Fatalf("failed to create book: %s", err)
			}
			if b1.Title != testTitle1 {
				t.Errorf("expected %q, received %q", testTitle1, b1.Title)
			}
			if b1.Description != testDescription1 {
				t.Errorf("expected %q, received %q", testDescription1, b1.Description)
			}
			if b1.Cover != testCover1 {
				t.Errorf("expected %q, received %q", testCover1, b1.Cover)
			}
			bookID1 = b1.ID
			b2, err := q.CreateBook(ctx, sqlc.CreateBookParams{
				Title:       testTitle2,
				Description: testDescription2,
				Cover:       testCover2,
			})
			if err != nil {
				t.Fatalf("failed to create book: %s", err)
			}
			bookID2 = b2.ID
		})

		t.Run("SetAuthor", func(t *testing.T) {
			err := q.SetBookAuthor(ctx, sqlc.SetBookAuthorParams{
				BookID:   bookID1,
				AuthorID: authorID1,
			})
			if err != nil {
				t.Fatalf("failed to set author book: %s", err)
			}
			err = q.SetBookAuthor(ctx, sqlc.SetBookAuthorParams{
				BookID:   bookID2,
				AuthorID: authorID2,
			})
			if err != nil {
				t.Fatalf("failed to set author book: %s", err)
			}
		})

		t.Run("ListBooks", func(t *testing.T) {
			l, err := q.ListBooks(ctx)
			if err != nil {
				t.Fatalf("failed to list books: %s", err)
			}
			if len(l) != 2 {
				t.Errorf("expected length of 2, received %d", len(l))
			}
			if l[0].ID != bookID1 {
				t.Errorf("expected %d, received %d", authorID1, l[0].ID)
			}
			if l[0].Title != testTitle1 {
				t.Errorf("expected %q, received %q", testTitle1, l[0].Title)
			}
			if l[0].Description != testDescription1 {
				t.Errorf("expected %q, received %q", testDescription1, l[0].Description)
			}
			if l[0].Cover != testCover1 {
				t.Errorf("expected %q, received %q", testCover1, l[0].Cover)
			}
			if l[1].ID != bookID2 {
				t.Errorf("expected %d, received %d", bookID2, l[1].ID)
			}
		})

		t.Run("ListAuthorsByBookIDs", func(t *testing.T) {
			l, err := q.ListAuthorsByBookID(ctx, bookID1)
			if err != nil {
				t.Fatalf("failed to list authors by book ids: %s", err)
			}
			if len(l) != 1 {
				t.Errorf("expected length of 1, received %d", len(l))
			}
			if l[0].ID != authorID1 {
				t.Errorf("expected %d, received %d", authorID1, l[0].ID)
			}
		})

		t.Run("ListBooksByAuthorIDs", func(t *testing.T) {
			l, err := q.ListBooksByAuthorID(ctx, authorID1)
			if err != nil {
				t.Fatalf("failed to list books by author ids: %s", err)
			}
			if len(l) != 1 {
				t.Errorf("expected length of 1, received %d", len(l))
			}
			if l[0].ID != bookID1 {
				t.Errorf("expected %d, received %d", bookID1, l[0].ID)
			}
		})

		t.Run("UpdateBook", func(t *testing.T) {
			b, err := q.UpdateBook(ctx, sqlc.UpdateBookParams{
				ID:          bookID2,
				Title:       testTitle3,
				Description: testDescription3,
				Cover:       testCover3,
			})
			if err != nil {
				t.Fatalf("failed to update book: %s", err)
			}
			if b.Title != testTitle3 {
				t.Errorf("expected %q, received %q", testTitle3, b.Title)
			}
			if b.Description != testDescription3 {
				t.Errorf("expected %q, received %q", testDescription3, b.Description)
			}
			if b.Cover != testCover3 {
				t.Errorf("expected %q, received %q", testCover3, b.Cover)
			}
		})

		t.Run("GetBook", func(t *testing.T) {
			b, err := q.GetBook(ctx, bookID2)
			if err != nil {
				t.Fatalf("failed to get book: %s", err)
			}
			if b.Title != testTitle3 {
				t.Errorf("expected %q, received %q", testTitle3, b.Title)
			}
			if b.Description != testDescription3 {
				t.Errorf("expected %q, received %q", testDescription3, b.Description)
			}
			if b.Cover != testCover3 {
				t.Errorf("expected %q, received %q", testCover3, b.Cover)
			}
		})

		t.Run("UnsetBookAuthors", func(t *testing.T) {
			err := q.UnsetBookAuthors(ctx, bookID2)
			if err != nil {
				t.Fatalf("failed to unset book authors: %s", err)
			}
			l, err := q.ListAuthorsByBookID(ctx, bookID2)
			if err != nil {
				t.Fatalf("failed to list authors by book ids: %s", err)
			}
			if len(l) != 0 {
				t.Errorf("expected length of 0, received %d", len(l))
			}
		})
	})
}

func runner(t *testing.T, test func(context.Context, postgres.Q, *testing.T)) {
	ctx := context.Background()

	repo, err := postgres.NewRepo("dbname=test_db sslmode=disable")
	if err != nil {
		t.Fatalf("failed to connect to the db: %s\n", err)
	}

	// create and drop schema (defer)
	defer func() {
		if err := dropSchema(ctx, repo); err != nil {
			t.Fatalf("failed to drop queue schema: %s\n", err)
		}
	}()
	if err := createSchema(ctx, repo); err != nil {
		t.Fatalf("failed to create queue schema: %s\n", err)
	}

	test(ctx, repo, t)
}

func createSchema(ctx context.Context, db postgres.DB) error {
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

func dropSchema(ctx context.Context, db postgres.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		DROP TABLE IF EXISTS book_authors, books, authors, agents;
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
