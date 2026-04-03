package identity

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-starter-template/db"
)

var ErrNotFound = errors.New("user not found")

type IdentityRepo interface {
	FindAll(ctx context.Context) ([]db.User, error)
	FindByID(ctx context.Context, id int32) (db.User, error)
	Save(ctx context.Context, params db.CreateUserParams) (db.User, error)
	Update(ctx context.Context, params db.UpdateUserParams) (db.User, error)
	Delete(ctx context.Context, id int32) error
}

type identityRepo struct {
	queries *db.Queries
}

func NewIdentityRepo(pool *pgxpool.Pool) IdentityRepo {
	return &identityRepo{queries: db.New(pool)}
}

func (r *identityRepo) FindAll(ctx context.Context) ([]db.User, error) {
	return r.queries.ListUsers(ctx)
}

func (r *identityRepo) FindByID(ctx context.Context, id int32) (db.User, error) {
	user, err := r.queries.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrNotFound
		}
		return db.User{}, err
	}
	return user, nil
}

func (r *identityRepo) Save(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	return r.queries.CreateUser(ctx, params)
}

func (r *identityRepo) Update(ctx context.Context, params db.UpdateUserParams) (db.User, error) {
	user, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrNotFound
		}
		return db.User{}, err
	}
	return user, nil
}

func (r *identityRepo) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteUser(ctx, id)
}
