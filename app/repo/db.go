package repo

import (
	"context"

	categorydb "medisuite-api/pkg/db/categories"
	permissiondb "medisuite-api/pkg/db/permissions"
	rolepermissiondb "medisuite-api/pkg/db/role_permissions"
	roledb "medisuite-api/pkg/db/roles"
	treatmentdb "medisuite-api/pkg/db/treatments"
	sessiondb "medisuite-api/pkg/db/user_sessions"
	userdb "medisuite-api/pkg/db/users"

	"github.com/jackc/pgx/v5"
)

// Queries groups all sqlc-generated query sets
// so they can be injected into the repository layer.
type Queries struct {
	Permissions     *permissiondb.Queries
	Roles           *roledb.Queries
	RolePermissions *rolepermissiondb.Queries
	Users           *userdb.Queries
	Sessions        *sessiondb.Queries
	Categories      *categorydb.Queries
	Treatments      *treatmentdb.Queries
}

// Store is the common abstraction for database access at the repository layer.
type Store interface {
	// Queries returns the set of query structs that can be used directly.
	Queries() *Queries
	// ExecTx runs fn within a single database transaction.
	ExecTx(ctx context.Context, fn func(*Queries) error) error
}

// SQLStore is a Store implementation backed by a pgx.Conn.
type SQLStore struct {
	queries *Queries
	conn    *pgx.Conn
}

// NewStore creates a new Store instance from a pgx connection.
func NewStore(conn *pgx.Conn) Store {
	return &SQLStore{
		queries: &Queries{
			Permissions:     permissiondb.New(conn),
			Roles:           roledb.New(conn),
			RolePermissions: rolepermissiondb.New(conn),
			Users:           userdb.New(conn),
			Sessions:        sessiondb.New(conn),
			Categories:      categorydb.New(conn),
			Treatments:      treatmentdb.New(conn),
		},
		conn: conn,
	}
}

// Queries returns non-transactional query structs.
func (s *SQLStore) Queries() *Queries {
	return s.queries
}

// ExecTx runs the provided callback inside a database transaction.
// Inside fn, the received Queries are already bound to the transaction.
func (s *SQLStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}

	q := &Queries{
		Permissions:     s.queries.Permissions.WithTx(tx),
		Roles:           s.queries.Roles.WithTx(tx),
		RolePermissions: s.queries.RolePermissions.WithTx(tx),
		Users:           s.queries.Users.WithTx(tx),
		Sessions:        s.queries.Sessions.WithTx(tx),
		Categories:      s.queries.Categories.WithTx(tx),
		Treatments:      s.queries.Treatments.WithTx(tx),
	}

	if err := fn(q); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
