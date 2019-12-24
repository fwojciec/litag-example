package resolvers_test

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/fwojciec/litag-example/generated/gqlgen"
	"github.com/fwojciec/litag-example/generated/mocks"
	"github.com/fwojciec/litag-example/generated/sqlc"
	"github.com/fwojciec/litag-example/postgres"
	"github.com/fwojciec/litag-example/resolvers"
)

var (
	testAgent   = &sqlc.Agent{ID: 99}
	testAuthor1 = &sqlc.Author{
		ID:      99,
		Name:    "test name 1",
		Website: sql.NullString{String: "https://t.com", Valid: true},
		AgentID: 22,
	}
	testAuthor2 = &sqlc.Author{
		ID:      199,
		Name:    "test name 2",
		Website: sql.NullString{},
		AgentID: 122,
	}
	testBook = &sqlc.Book{
		ID:          88,
		Title:       "test title 1",
		Description: "test description 1",
		Cover:       "cover1.jpg",
	}
	testError = errors.New("test error")
)

func TestAgentResolver(t *testing.T) {
	t.Parallel()
	t.Run("Authors", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			agent *sqlc.Agent
			err   error
		}{
			{"valid", testAgent, nil},
			{"error", testAgent, testError},
		}
		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				var receivedAgentID int64
				r := &resolvers.Resolver{
					Repo: &postgres.Repo{
						Q: &mocks.QMock{
							ListAuthorsByAgentIDFunc: func(ctx context.Context, agentID int64) ([]sqlc.Author, error) {
								receivedAgentID = agentID
								return nil, tc.err
							},
						},
					},
				}
				_, err := r.Agent().Authors(context.Background(), tc.agent)
				if !errors.Is(err, tc.err) {
					t.Errorf("wrong error: expected %v, received %v", tc.err, err)
				}
				if receivedAgentID != tc.agent.ID {
					t.Errorf("wrong id: expected %d, received %d", tc.agent.ID, receivedAgentID)
				}
			})
		}
	})
}

func TestAuthorResolver(t *testing.T) {
	t.Parallel()

	t.Run("Website", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name   string
			author *sqlc.Author
		}{
			{"has website", testAuthor1},
			{"has no website", testAuthor2},
		}
		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				r := &resolvers.Resolver{}
				res, err := r.Author().Website(context.Background(), tc.author)
				if err != nil {
					t.Errorf("expected no error, received: %v", err)
				}
				if !tc.author.Website.Valid && res != nil {
					t.Errorf("expected result to be nil, received: %v", res)
				}
				if tc.author.Website.Valid && res == nil {
					t.Fatalf("expected pointer not to be nil")
				}
				if tc.author.Website.Valid && tc.author.Website.String != *res {
					t.Errorf("expected %s, received %s", *res, tc.author.Website.String)
				}
			})
		}
	})

	t.Run("Agent", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name   string
			author *sqlc.Author
			err    error
		}{
			{"valid", testAuthor1, nil},
			{"error", testAuthor2, testError},
		}
		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				var receivedAgentID int64
				r := &resolvers.Resolver{
					Repo: &postgres.Repo{
						Q: &mocks.QMock{
							GetAgentFunc: func(ctx context.Context, id int64) (sqlc.Agent, error) {
								receivedAgentID = id
								return sqlc.Agent{}, tc.err
							},
						},
					},
				}
				_, err := r.Author().Agent(context.Background(), tc.author)
				if !errors.Is(err, tc.err) {
					t.Errorf("wrong error: expected %v, received %v", tc.err, err)
				}
				if receivedAgentID != tc.author.AgentID {
					t.Errorf("wrong id: expected %d, received %d", tc.author.AgentID, receivedAgentID)
				}
			})
		}
	})

	t.Run("Books", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name   string
			author *sqlc.Author
			err    error
		}{
			{"valid", testAuthor1, nil},
			{"error", testAuthor2, testError},
		}
		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				var receivedAuthorID int64
				r := &resolvers.Resolver{
					Repo: &postgres.Repo{
						Q: &mocks.QMock{
							ListBooksByAuthorIDFunc: func(ctx context.Context, authorID int64) ([]sqlc.Book, error) {
								receivedAuthorID = authorID
								return nil, tc.err
							},
						},
					},
				}
				_, err := r.Author().Books(context.Background(), tc.author)
				if !errors.Is(err, tc.err) {
					t.Errorf("wrong error: expected %v, received %v", tc.err, err)
				}
				if receivedAuthorID != tc.author.ID {
					t.Errorf("wrong id: expected %d, received %d", tc.author.ID, receivedAuthorID)
				}
			})
		}
	})
}

func TestBookResolver(t *testing.T) {
	t.Parallel()
	t.Run("Authors", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name string
			book *sqlc.Book
			err  error
		}{
			{"valid", testBook, nil},
			{"error", testBook, testError},
		}
		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				var receivedBookID int64
				r := &resolvers.Resolver{
					Repo: &postgres.Repo{
						Q: &mocks.QMock{
							ListAuthorsByBookIDFunc: func(ctx context.Context, bookID int64) ([]sqlc.Author, error) {
								receivedBookID = bookID
								return nil, tc.err
							},
						},
					},
				}
				_, err := r.Book().Authors(context.Background(), tc.book)
				if !errors.Is(err, tc.err) {
					t.Errorf("wrong error: expected %v, received %v", tc.err, err)
				}
				if receivedBookID != tc.book.ID {
					t.Errorf("wrong id: expected %d, received %d", tc.book.ID, receivedBookID)
				}
			})
		}
	})
}

func TestMutationResolver(t *testing.T) {
	t.Parallel()

	t.Run("Agent mutations", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name  string
			agent *sqlc.Agent
			err   error
		}{
			{"valid", testAgent, nil},
			{"error", testAgent, testError},
		}

		t.Run("CreateAgent", func(t *testing.T) {
			t.Parallel()
			for _, tc := range tests {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()
					var receivedCreateAgentParams sqlc.CreateAgentParams
					r := &resolvers.Resolver{
						Repo: &postgres.Repo{
							Q: &mocks.QMock{
								CreateAgentFunc: func(ctx context.Context, args sqlc.CreateAgentParams) (sqlc.Agent, error) {
									receivedCreateAgentParams = args
									return sqlc.Agent{}, tc.err
								},
							},
						},
					}
					_, err := r.Mutation().CreateAgent(context.Background(), gqlgen.CreateUpdateAgentInput{
						Name:  tc.agent.Name,
						Email: tc.agent.Email,
					})
					if !errors.Is(err, tc.err) {
						t.Errorf("wrong error: expected %v, received %v", tc.err, err)
					}
					exp := sqlc.CreateAgentParams{
						Name:  tc.agent.Name,
						Email: tc.agent.Email,
					}
					if !reflect.DeepEqual(receivedCreateAgentParams, exp) {
						t.Errorf("wrong params: expected %v, received %v", exp, receivedCreateAgentParams)
					}
				})
			}
		})

		t.Run("UpdateAgent", func(t *testing.T) {
			t.Parallel()
			for _, tc := range tests {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()
					var receivedUpdateAgentParams sqlc.UpdateAgentParams
					r := &resolvers.Resolver{
						Repo: &postgres.Repo{
							Q: &mocks.QMock{
								UpdateAgentFunc: func(ctx context.Context, args sqlc.UpdateAgentParams) (sqlc.Agent, error) {
									receivedUpdateAgentParams = args
									return sqlc.Agent{}, tc.err
								},
							},
						},
					}
					_, err := r.Mutation().UpdateAgent(context.Background(), tc.agent.ID, gqlgen.CreateUpdateAgentInput{
						Name:  tc.agent.Name,
						Email: tc.agent.Email,
					})
					if !errors.Is(err, tc.err) {
						t.Errorf("wrong error: expected %v, received %v", tc.err, err)
					}
					exp := sqlc.UpdateAgentParams{
						ID:    tc.agent.ID,
						Name:  tc.agent.Name,
						Email: tc.agent.Email,
					}
					if !reflect.DeepEqual(receivedUpdateAgentParams, exp) {
						t.Errorf("wrong params: expected %v, received %v", exp, receivedUpdateAgentParams)
					}
				})
			}
		})

		t.Run("DeleteAgent", func(t *testing.T) {
			t.Parallel()
			for _, tc := range tests {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()
					var receivedAgentID int64
					r := &resolvers.Resolver{
						Repo: &postgres.Repo{
							Q: &mocks.QMock{
								DeleteAgentFunc: func(ctx context.Context, id int64) (sqlc.Agent, error) {
									receivedAgentID = id
									return sqlc.Agent{}, tc.err
								},
							},
						},
					}
					_, err := r.Mutation().DeleteAgent(context.Background(), tc.agent.ID)
					if !errors.Is(err, tc.err) {
						t.Errorf("wrong error: expected %v, received %v", tc.err, err)
					}
					if receivedAgentID != tc.agent.ID {
						t.Errorf("wrong id: expected %d, received %d", tc.agent.ID, receivedAgentID)
					}
				})
			}
		})
	})

	t.Run("Author mutations", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name   string
			author *sqlc.Author
			err    error
		}{
			{"valid", testAuthor1, nil},
			{"valid null website", testAuthor2, nil},
			{"error", testAuthor1, testError},
		}

		t.Run("CreateAuthor", func(t *testing.T) {
			t.Parallel()
			for _, tc := range tests {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()
					var receivedCreateAuthorParams sqlc.CreateAuthorParams
					r := &resolvers.Resolver{
						Repo: &postgres.Repo{
							Q: &mocks.QMock{
								CreateAuthorFunc: func(ctx context.Context, args sqlc.CreateAuthorParams) (sqlc.Author, error) {
									receivedCreateAuthorParams = args
									return sqlc.Author{}, tc.err
								},
							},
						},
					}
					_, err := r.Mutation().CreateAuthor(context.Background(), gqlgen.CreateUpdateAuthorInput{
						Name:    tc.author.Name,
						Website: nullStringToPointer(tc.author.Website),
						AgentID: tc.author.AgentID,
					})
					if !errors.Is(err, tc.err) {
						t.Errorf("wrong error: expected %v, received %v", tc.err, err)
					}
					exp := sqlc.CreateAuthorParams{
						Name:    tc.author.Name,
						Website: tc.author.Website,
						AgentID: tc.author.AgentID,
					}
					if !reflect.DeepEqual(receivedCreateAuthorParams, exp) {
						t.Errorf("wrong params: expected %v, received %v", exp, receivedCreateAuthorParams)
					}
				})
			}
		})

		t.Run("UpdateAuthor", func(t *testing.T) {
			t.Parallel()
			for _, tc := range tests {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()
					var receivedUpdateAuthorParams sqlc.UpdateAuthorParams
					r := &resolvers.Resolver{
						Repo: &postgres.Repo{
							Q: &mocks.QMock{
								UpdateAuthorFunc: func(ctx context.Context, args sqlc.UpdateAuthorParams) (sqlc.Author, error) {
									receivedUpdateAuthorParams = args
									return sqlc.Author{}, tc.err
								},
							},
						},
					}
					_, err := r.Mutation().UpdateAuthor(context.Background(), tc.author.ID, gqlgen.CreateUpdateAuthorInput{
						Name:    tc.author.Name,
						Website: nullStringToPointer(tc.author.Website),
						AgentID: tc.author.AgentID,
					})
					if !errors.Is(err, tc.err) {
						t.Errorf("wrong error: expected %v, received %v", tc.err, err)
					}
					exp := sqlc.UpdateAuthorParams{
						ID:      tc.author.ID,
						Name:    tc.author.Name,
						Website: tc.author.Website,
						AgentID: tc.author.AgentID,
					}
					if !reflect.DeepEqual(receivedUpdateAuthorParams, exp) {
						t.Errorf("wrong params: expected %v, received %v", exp, receivedUpdateAuthorParams)
					}
				})
			}
		})

		t.Run("DeleteAuthor", func(t *testing.T) {
			t.Parallel()
			for _, tc := range tests {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()
					var receivedAuthorID int64
					r := &resolvers.Resolver{
						Repo: &postgres.Repo{
							Q: &mocks.QMock{
								DeleteAuthorFunc: func(ctx context.Context, id int64) (sqlc.Author, error) {
									receivedAuthorID = id
									return sqlc.Author{}, tc.err
								},
							},
						},
					}
					_, err := r.Mutation().DeleteAuthor(context.Background(), tc.author.ID)
					if !errors.Is(err, tc.err) {
						t.Errorf("wrong error: expected %v, received %v", tc.err, err)
					}
					if receivedAuthorID != tc.author.ID {
						t.Errorf("wrong id: expected %d, received %d", tc.author.ID, receivedAuthorID)
					}
				})
			}
		})
	})

	t.Run("Book mutations", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name    string
			book    *sqlc.Book
			authors []int64
			err     error
		}{
			{"valid", testBook, []int64{testAuthor1.ID, testAuthor2.ID}, nil},
			{"error", testBook, []int64{testAuthor1.ID}, testError},
		}

		t.Run("CreateBook", func(t *testing.T) {
			t.Parallel()
			for _, tc := range tests {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()
					var receivedCreateBookParams sqlc.CreateBookParams
					var receivedAuthorIDs []int64
					r := &resolvers.Resolver{
						Repo: &postgres.Repo{
							TxQ: &mocks.TxQMock{
								CreateBookFunc: func(ctx context.Context, args sqlc.CreateBookParams, authorIDs []int64) (*sqlc.Book, error) {
									receivedCreateBookParams = args
									receivedAuthorIDs = authorIDs
									return nil, tc.err
								},
							},
						},
					}
					_, err := r.Mutation().CreateBook(context.Background(), gqlgen.CreateUpdateBookInput{
						Title:       tc.book.Title,
						Description: tc.book.Description,
						Cover:       tc.book.Cover,
						AuthorIDs:   tc.authors,
					})
					if !errors.Is(err, tc.err) {
						t.Errorf("wrong error: expected %v, received %v", tc.err, err)
					}
					exp := sqlc.CreateBookParams{
						Title:       tc.book.Title,
						Description: tc.book.Description,
						Cover:       tc.book.Cover,
					}
					if !reflect.DeepEqual(receivedCreateBookParams, exp) {
						t.Errorf("wrong params: expected %v, received %v", exp, receivedCreateBookParams)
					}
					if !reflect.DeepEqual(receivedAuthorIDs, tc.authors) {
						t.Errorf("wrong ids: expected %v, received %v", tc.authors, receivedAuthorIDs)
					}
				})
			}
		})

		t.Run("UpdateBook", func(t *testing.T) {
			t.Parallel()
			for _, tc := range tests {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()
					var receivedUpdateBookParams sqlc.UpdateBookParams
					var receivedAuthorIDs []int64
					r := &resolvers.Resolver{
						Repo: &postgres.Repo{
							TxQ: &mocks.TxQMock{
								UpdateBookFunc: func(ctx context.Context, args sqlc.UpdateBookParams, authorIDs []int64) (*sqlc.Book, error) {
									receivedUpdateBookParams = args
									receivedAuthorIDs = authorIDs
									return nil, tc.err
								},
							},
						},
					}
					_, err := r.Mutation().UpdateBook(context.Background(), tc.book.ID, gqlgen.CreateUpdateBookInput{
						Title:       tc.book.Title,
						Description: tc.book.Description,
						Cover:       tc.book.Cover,
						AuthorIDs:   tc.authors,
					})
					if !errors.Is(err, tc.err) {
						t.Errorf("wrong error: expected %v, received %v", tc.err, err)
					}
					exp := sqlc.UpdateBookParams{
						ID:          tc.book.ID,
						Title:       tc.book.Title,
						Description: tc.book.Description,
						Cover:       tc.book.Cover,
					}
					if !reflect.DeepEqual(receivedUpdateBookParams, exp) {
						t.Errorf("wrong params: expected %v, received %v", exp, receivedUpdateBookParams)
					}
					if !reflect.DeepEqual(receivedAuthorIDs, tc.authors) {
						t.Errorf("wrong ids: expected %v, received %v", tc.authors, receivedAuthorIDs)
					}
				})
			}
		})

		t.Run("DeleteBook", func(t *testing.T) {
			t.Parallel()
			for _, tc := range tests {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()
					var receivedBookID int64
					r := &resolvers.Resolver{
						Repo: &postgres.Repo{
							Q: &mocks.QMock{
								DeleteBookFunc: func(ctx context.Context, id int64) (sqlc.Book, error) {
									receivedBookID = id
									return sqlc.Book{}, tc.err
								},
							},
						},
					}
					_, err := r.Mutation().DeleteBook(context.Background(), tc.book.ID)
					if !errors.Is(err, tc.err) {
						t.Errorf("wrong error: expected %v, received %v", tc.err, err)
					}
					if receivedBookID != tc.book.ID {
						t.Errorf("wrong id: expected %d, received %d", tc.book.ID, receivedBookID)
					}
				})
			}
		})
	})
}

func TestQueryResolver(t *testing.T) {
	t.Parallel()

	t.Run("Agent", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name string
			id   int64
			err  error
		}{
			{"valid", testAgent.ID, nil},
			{"error", testAgent.ID, testError},
		}
		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				var receivedID int64
				r := &resolvers.Resolver{
					Repo: &postgres.Repo{
						Q: &mocks.QMock{
							GetAgentFunc: func(ctx context.Context, id int64) (sqlc.Agent, error) {
								receivedID = id
								return sqlc.Agent{}, tc.err
							},
						},
					},
				}
				_, err := r.Query().Agent(context.Background(), tc.id)
				if !errors.Is(err, tc.err) {
					t.Errorf("wrong error: expected %v, received %v", tc.err, err)
				}
				if receivedID != tc.id {
					t.Errorf("wrong id: expected %d, received %d", tc.id, receivedID)
				}
			})
		}
	})

	t.Run("Agents", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name string
			err  error
		}{
			{"valid", nil},
			{"error", testError},
		}
		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				r := &resolvers.Resolver{
					Repo: &postgres.Repo{
						Q: &mocks.QMock{
							ListAgentsFunc: func(ctx context.Context) ([]sqlc.Agent, error) {
								return nil, tc.err
							},
						},
					},
				}
				_, err := r.Query().Agents(context.Background())
				if !errors.Is(err, tc.err) {
					t.Errorf("wrong error: expected %v, received %v", tc.err, err)
				}
			})
		}
	})

	t.Run("Author", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name string
			id   int64
			err  error
		}{
			{"valid", testAuthor1.ID, nil},
			{"error", testAuthor1.ID, testError},
		}
		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				var receivedID int64
				r := &resolvers.Resolver{
					Repo: &postgres.Repo{
						Q: &mocks.QMock{
							GetAuthorFunc: func(ctx context.Context, id int64) (sqlc.Author, error) {
								receivedID = id
								return sqlc.Author{}, tc.err
							},
						},
					},
				}
				_, err := r.Query().Author(context.Background(), tc.id)
				if !errors.Is(err, tc.err) {
					t.Errorf("wrong error: expected %v, received %v", tc.err, err)
				}
				if receivedID != tc.id {
					t.Errorf("wrong id: expected %d, received %d", tc.id, receivedID)
				}
			})
		}
	})

	t.Run("Authors", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name string
			err  error
		}{
			{"valid", nil},
			{"error", testError},
		}
		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				r := &resolvers.Resolver{
					Repo: &postgres.Repo{
						Q: &mocks.QMock{
							ListAuthorsFunc: func(ctx context.Context) ([]sqlc.Author, error) {
								return nil, tc.err
							},
						},
					},
				}
				_, err := r.Query().Authors(context.Background())
				if !errors.Is(err, tc.err) {
					t.Errorf("wrong error: expected %v, received %v", tc.err, err)
				}
			})
		}
	})

	t.Run("Book", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name string
			id   int64
			err  error
		}{
			{"valid", testBook.ID, nil},
			{"error", testBook.ID, testError},
		}
		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				var receivedID int64
				r := &resolvers.Resolver{
					Repo: &postgres.Repo{
						Q: &mocks.QMock{
							GetBookFunc: func(ctx context.Context, id int64) (sqlc.Book, error) {
								receivedID = id
								return sqlc.Book{}, tc.err
							},
						},
					},
				}
				_, err := r.Query().Book(context.Background(), tc.id)
				if !errors.Is(err, tc.err) {
					t.Errorf("wrong error: expected %v, received %v", tc.err, err)
				}
				if receivedID != tc.id {
					t.Errorf("wrong id: expected %d, received %d", tc.id, receivedID)
				}
			})
		}
	})

	t.Run("Books", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name string
			err  error
		}{
			{"valid", nil},
			{"error", testError},
		}
		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				r := &resolvers.Resolver{
					Repo: &postgres.Repo{
						Q: &mocks.QMock{
							ListBooksFunc: func(ctx context.Context) ([]sqlc.Book, error) {
								return nil, tc.err
							},
						},
					},
				}
				_, err := r.Query().Books(context.Background())
				if !errors.Is(err, tc.err) {
					t.Errorf("wrong error: expected %v, received %v", tc.err, err)
				}
			})
		}
	})
}

func nullStringToPointer(ns sql.NullString) *string {
	var s *string
	if ns.Valid {
		s = &ns.String
	}
	return s
}
