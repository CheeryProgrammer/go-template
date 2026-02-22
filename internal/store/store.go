package store

import "context"

// Store defines the data access layer. Add your domain-specific methods here.
type Store interface {
	Ping(ctx context.Context) error

	// TODO: add your query methods
	// Example:
	//   GetUser(ctx context.Context, id int64) (User, error)
	//   CreateUser(ctx context.Context, params CreateUserParams) (User, error)
}
