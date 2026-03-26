package ordering

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-starter-template/db"
)

var ErrNotFound = errors.New("product not found")

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type OrderingRepo interface {
	FindAll() ([]Product, error)
	FindByID(id int) (Product, error)
	Save(product Product) (Product, error)
	Update(product Product) (Product, error)
	Delete(id int) error
}

type orderingRepo struct {
	queries *db.Queries
}

func NewOrderingRepo(pool *pgxpool.Pool) OrderingRepo {
	return &orderingRepo{queries: db.New(pool)}
}

func (r *orderingRepo) FindAll() ([]Product, error) {
	rows, err := r.queries.ListProducts(context.Background())
	if err != nil {
		return nil, err
	}
	products := make([]Product, len(rows))
	for i, row := range rows {
		products[i] = toProduct(row)
	}
	return products, nil
}

func (r *orderingRepo) FindByID(id int) (Product, error) {
	row, err := r.queries.GetProduct(context.Background(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Product{}, ErrNotFound
		}
		return Product{}, err
	}
	return toProduct(row), nil
}

func (r *orderingRepo) Save(product Product) (Product, error) {
	row, err := r.queries.CreateProduct(context.Background(), db.CreateProductParams{
		Name:  product.Name,
		Price: product.Price,
	})
	if err != nil {
		return Product{}, err
	}
	return toProduct(row), nil
}

func (r *orderingRepo) Update(product Product) (Product, error) {
	row, err := r.queries.UpdateProduct(context.Background(), db.UpdateProductParams{
		ID:    int32(product.ID),
		Name:  product.Name,
		Price: product.Price,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Product{}, ErrNotFound
		}
		return Product{}, err
	}
	return toProduct(row), nil
}

func (r *orderingRepo) Delete(id int) error {
	return r.queries.DeleteProduct(context.Background(), int32(id))
}

func toProduct(p db.Product) Product {
	return Product{ID: int(p.ID), Name: p.Name, Price: p.Price}
}
