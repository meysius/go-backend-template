package ordering

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-starter-template/db"
)

var ErrNotFound = errors.New("product not found")

type OrderingRepo interface {
	FindAll() ([]db.Product, error)
	FindByID(id int32) (db.Product, error)
	Save(params db.CreateProductParams) (db.Product, error)
	Update(params db.UpdateProductParams) (db.Product, error)
	Delete(id int32) error
}

type orderingRepo struct {
	queries *db.Queries
}

func NewOrderingRepo(pool *pgxpool.Pool) OrderingRepo {
	return &orderingRepo{queries: db.New(pool)}
}

func (r *orderingRepo) FindAll() ([]db.Product, error) {
	return r.queries.ListProducts(context.Background())
}

func (r *orderingRepo) FindByID(id int32) (db.Product, error) {
	product, err := r.queries.GetProduct(context.Background(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Product{}, ErrNotFound
		}
		return db.Product{}, err
	}
	return product, nil
}

func (r *orderingRepo) Save(params db.CreateProductParams) (db.Product, error) {
	return r.queries.CreateProduct(context.Background(), params)
}

func (r *orderingRepo) Update(params db.UpdateProductParams) (db.Product, error) {
	product, err := r.queries.UpdateProduct(context.Background(), params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Product{}, ErrNotFound
		}
		return db.Product{}, err
	}
	return product, nil
}

func (r *orderingRepo) Delete(id int32) error {
	return r.queries.DeleteProduct(context.Background(), id)
}
