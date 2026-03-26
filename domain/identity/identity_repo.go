package identity

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-starter-template/db"
)

var ErrNotFound = errors.New("user not found")

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type IdentityRepo interface {
	FindAll() ([]User, error)
	FindByID(id int) (User, error)
	Save(user User) (User, error)
	Update(user User) (User, error)
	Delete(id int) error
}

type identityRepo struct {
	queries *db.Queries
}

func NewIdentityRepo(pool *pgxpool.Pool) IdentityRepo {
	return &identityRepo{queries: db.New(pool)}
}

func (r *identityRepo) FindAll() ([]User, error) {
	rows, err := r.queries.ListUsers(context.Background())
	if err != nil {
		return nil, err
	}
	users := make([]User, len(rows))
	for i, row := range rows {
		users[i] = toUser(row)
	}
	return users, nil
}

func (r *identityRepo) FindByID(id int) (User, error) {
	row, err := r.queries.GetUser(context.Background(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}
	return toUser(row), nil
}

func (r *identityRepo) Save(user User) (User, error) {
	row, err := r.queries.CreateUser(context.Background(), db.CreateUserParams{
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		return User{}, err
	}
	return toUser(row), nil
}

func (r *identityRepo) Update(user User) (User, error) {
	row, err := r.queries.UpdateUser(context.Background(), db.UpdateUserParams{
		ID:    int32(user.ID),
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}
	return toUser(row), nil
}

func (r *identityRepo) Delete(id int) error {
	return r.queries.DeleteUser(context.Background(), int32(id))
}

func toUser(u db.User) User {
	return User{ID: int(u.ID), Name: u.Name, Email: u.Email}
}
