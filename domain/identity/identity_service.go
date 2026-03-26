package identity

type IdentityService struct {
	repo IdentityRepo
}

func NewIdentityService(repo IdentityRepo) *IdentityService {
	return &IdentityService{repo: repo}
}

func (s *IdentityService) ListUsers() ([]User, error) {
	return s.repo.FindAll()
}

func (s *IdentityService) GetUser(id int) (User, error) {
	return s.repo.FindByID(id)
}

func (s *IdentityService) CreateUser(name, email string) (User, error) {
	return s.repo.Save(User{Name: name, Email: email})
}

func (s *IdentityService) UpdateUser(id int, name, email string) (User, error) {
	return s.repo.Update(User{ID: id, Name: name, Email: email})
}

func (s *IdentityService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}
