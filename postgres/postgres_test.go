package postgres_test

import (
	"context"
	"testing"

	"github.com/fwojciec/litag-example/generated/sqlc"
	"github.com/fwojciec/litag-example/postgres"
)

func TestQueries(t *testing.T) {

	var (
		testName1  = "test name 1"
		testName2  = "test name 2"
		testName3  = "test name 3"
		testEmail1 = "test1@email.com"
		testEmail2 = "test2@email.com"
		testEmail3 = "test3@email.com"
		agentID1   int64
		agentID2   int64
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
