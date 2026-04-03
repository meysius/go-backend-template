package identity

import (
	"context"

	"go-starter-template/db"
)

type IdentityService struct {
	repo IdentityRepo
}

func NewIdentityService(repo IdentityRepo) *IdentityService {
	return &IdentityService{repo: repo}
}

func (s *IdentityService) ListUsers(ctx context.Context) ([]db.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *IdentityService) GetUser(ctx context.Context, id int32) (db.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *IdentityService) CreateUser(ctx context.Context, name, email string) (db.User, error) {
	return s.repo.Save(ctx, db.CreateUserParams{Name: name, Email: email})
}

func (s *IdentityService) UpdateUser(ctx context.Context, id int32, name, email string) (db.User, error) {
	return s.repo.Update(ctx, db.UpdateUserParams{ID: id, Name: name, Email: email})
}

func (s *IdentityService) DeleteUser(ctx context.Context, id int32) error {
	return s.repo.Delete(ctx, id)
}
