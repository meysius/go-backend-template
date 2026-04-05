package identity

import (
	"context"
	"errors"
	"testing"

	"go-starter-template/db"
)

func TestListUsers_Empty(t *testing.T) {
	svc := NewIdentityService(newMockRepo())
	users, err := svc.ListUsers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("expected 0 users, got %d", len(users))
	}
}

func TestListUsers_ReturnsAll(t *testing.T) {
	repo := newMockRepo()
	repo.users = []db.User{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
	}
	svc := NewIdentityService(repo)

	users, err := svc.ListUsers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}

func TestListUsers_Error(t *testing.T) {
	repo := newMockRepo()
	repo.err = errors.New("db down")
	svc := NewIdentityService(repo)

	_, err := svc.ListUsers(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetUser_Found(t *testing.T) {
	repo := newMockRepo()
	repo.users = []db.User{{ID: 5, Name: "Carol", Email: "carol@example.com"}}
	svc := NewIdentityService(repo)

	user, err := svc.GetUser(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Name != "Carol" {
		t.Fatalf("expected Carol, got %s", user.Name)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	svc := NewIdentityService(newMockRepo())

	_, err := svc.GetUser(context.Background(), 999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCreateUser(t *testing.T) {
	repo := newMockRepo()
	svc := NewIdentityService(repo)

	user, err := svc.CreateUser(context.Background(), "Dave", "dave@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != 1 || user.Name != "Dave" || user.Email != "dave@example.com" {
		t.Fatalf("unexpected user: %+v", user)
	}
	if len(repo.users) != 1 {
		t.Fatalf("expected 1 stored user, got %d", len(repo.users))
	}
}

func TestCreateUser_Error(t *testing.T) {
	repo := newMockRepo()
	repo.err = errors.New("duplicate email")
	svc := NewIdentityService(repo)

	_, err := svc.CreateUser(context.Background(), "Dave", "dave@example.com")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateUser_Found(t *testing.T) {
	repo := newMockRepo()
	repo.users = []db.User{{ID: 1, Name: "Old", Email: "old@example.com"}}
	svc := NewIdentityService(repo)

	user, err := svc.UpdateUser(context.Background(), 1, "New", "new@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Name != "New" || user.Email != "new@example.com" {
		t.Fatalf("unexpected user: %+v", user)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	svc := NewIdentityService(newMockRepo())

	_, err := svc.UpdateUser(context.Background(), 999, "X", "x@example.com")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	repo := newMockRepo()
	repo.users = []db.User{{ID: 1, Name: "Del", Email: "del@example.com"}}
	svc := NewIdentityService(repo)

	if err := svc.DeleteUser(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.users) != 0 {
		t.Fatalf("expected 0 users after delete, got %d", len(repo.users))
	}
}

func TestDeleteUser_Error(t *testing.T) {
	repo := newMockRepo()
	repo.err = errors.New("db error")
	svc := NewIdentityService(repo)

	err := svc.DeleteUser(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
