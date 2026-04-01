package ordering

import "go-starter-template/db"

type OrderingService struct {
	repo OrderingRepo
}

func NewOrderingService(repo OrderingRepo) *OrderingService {
	return &OrderingService{repo: repo}
}

func (s *OrderingService) ListProducts() ([]db.Product, error) {
	return s.repo.FindAll()
}

func (s *OrderingService) GetProduct(id int32) (db.Product, error) {
	return s.repo.FindByID(id)
}

func (s *OrderingService) CreateProduct(name string, price float64) (db.Product, error) {
	return s.repo.Save(db.CreateProductParams{Name: name, Price: price})
}

func (s *OrderingService) UpdateProduct(id int32, name string, price float64) (db.Product, error) {
	return s.repo.Update(db.UpdateProductParams{ID: id, Name: name, Price: price})
}

func (s *OrderingService) DeleteProduct(id int32) error {
	return s.repo.Delete(id)
}
