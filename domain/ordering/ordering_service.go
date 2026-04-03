package ordering

import (
	"context"

	"go-starter-template/db"
)

type OrderingService struct {
	repo OrderingRepo
}

func NewOrderingService(repo OrderingRepo) *OrderingService {
	return &OrderingService{repo: repo}
}

func (s *OrderingService) ListProducts(ctx context.Context) ([]db.Product, error) {
	return s.repo.FindAll(ctx)
}

func (s *OrderingService) GetProduct(ctx context.Context, id int32) (db.Product, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *OrderingService) CreateProduct(ctx context.Context, name string, price float64) (db.Product, error) {
	return s.repo.Save(ctx, db.CreateProductParams{Name: name, Price: price})
}

func (s *OrderingService) UpdateProduct(ctx context.Context, id int32, name string, price float64) (db.Product, error) {
	return s.repo.Update(ctx, db.UpdateProductParams{ID: id, Name: name, Price: price})
}

func (s *OrderingService) DeleteProduct(ctx context.Context, id int32) error {
	return s.repo.Delete(ctx, id)
}
