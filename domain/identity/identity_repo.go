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
	FindAll() ([]db.User, error)
	FindByID(id int32) (db.User, error)
	Save(params db.CreateUserParams) (db.User, error)
	Update(params db.UpdateUserParams) (db.User, error)
	Delete(id int32) error
}

type identityRepo struct {
	queries *db.Queries
}

func NewIdentityRepo(pool *pgxpool.Pool) IdentityRepo {
	return &identityRepo{queries: db.New(pool)}
}

func (r *identityRepo) FindAll() ([]db.User, error) {
	return r.queries.ListUsers(context.Background())
}

func (r *identityRepo) FindByID(id int32) (db.User, error) {
	user, err := r.queries.GetUser(context.Background(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrNotFound
		}
		return db.User{}, err
	}
	return user, nil
}

func (r *identityRepo) Save(params db.CreateUserParams) (db.User, error) {
	return r.queries.CreateUser(context.Background(), params)
}

func (r *identityRepo) Update(params db.UpdateUserParams) (db.User, error) {
	user, err := r.queries.UpdateUser(context.Background(), params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrNotFound
		}
		return db.User{}, err
	}
	return user, nil
}

func (r *identityRepo) Delete(id int32) error {
	return r.queries.DeleteUser(context.Background(), id)
}
