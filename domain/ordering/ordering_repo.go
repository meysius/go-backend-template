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
	FindAll(ctx context.Context) ([]db.Product, error)
	FindByID(ctx context.Context, id int32) (db.Product, error)
	Save(ctx context.Context, params db.CreateProductParams) (db.Product, error)
	Update(ctx context.Context, params db.UpdateProductParams) (db.Product, error)
	Delete(ctx context.Context, id int32) error
}

type orderingRepo struct {
	queries *db.Queries
}

func NewOrderingRepo(pool *pgxpool.Pool) OrderingRepo {
	return &orderingRepo{queries: db.New(pool)}
}

func (r *orderingRepo) FindAll(ctx context.Context) ([]db.Product, error) {
	return r.queries.ListProducts(ctx)
}

func (r *orderingRepo) FindByID(ctx context.Context, id int32) (db.Product, error) {
	product, err := r.queries.GetProduct(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Product{}, ErrNotFound
		}
		return db.Product{}, err
	}
	return product, nil
}

func (r *orderingRepo) Save(ctx context.Context, params db.CreateProductParams) (db.Product, error) {
	return r.queries.CreateProduct(ctx, params)
}

func (r *orderingRepo) Update(ctx context.Context, params db.UpdateProductParams) (db.Product, error) {
	product, err := r.queries.UpdateProduct(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Product{}, ErrNotFound
		}
		return db.Product{}, err
	}
	return product, nil
}

func (r *orderingRepo) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteProduct(ctx, id)
}
