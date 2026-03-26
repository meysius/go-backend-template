package ordering

type OrderingService struct {
	repo OrderingRepo
}

func NewOrderingService(repo OrderingRepo) *OrderingService {
	return &OrderingService{repo: repo}
}

func (s *OrderingService) ListProducts() ([]Product, error) {
	return s.repo.FindAll()
}

func (s *OrderingService) GetProduct(id int) (Product, error) {
	return s.repo.FindByID(id)
}

func (s *OrderingService) CreateProduct(name string, price float64) (Product, error) {
	return s.repo.Save(Product{Name: name, Price: price})
}

func (s *OrderingService) UpdateProduct(id int, name string, price float64) (Product, error) {
	return s.repo.Update(Product{ID: id, Name: name, Price: price})
}

func (s *OrderingService) DeleteProduct(id int) error {
	return s.repo.Delete(id)
}
