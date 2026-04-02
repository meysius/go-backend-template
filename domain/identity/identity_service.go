package identity

import (
	"go-starter-template/db"
)

type IdentityService struct {
	repo IdentityRepo
}

func NewIdentityService(repo IdentityRepo) *IdentityService {
	return &IdentityService{repo: repo}
}

func (s *IdentityService) ListUsers() ([]db.User, error) {
	return s.repo.FindAll()
}

func (s *IdentityService) GetUser(id int32) (db.User, error) {
	return s.repo.FindByID(id)
}

func (s *IdentityService) CreateUser(name, email string) (db.User, error) {
	return s.repo.Save(db.CreateUserParams{Name: name, Email: email})
}

func (s *IdentityService) UpdateUser(id int32, name, email string) (db.User, error) {
	return s.repo.Update(db.UpdateUserParams{ID: id, Name: name, Email: email})
}

func (s *IdentityService) DeleteUser(id int32) error {
	return s.repo.Delete(id)
}
