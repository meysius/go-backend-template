package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"go-starter-template/db"
	"go-starter-template/domain/identity"
)

type mockRepo struct {
	users  []db.User
	nextID int32
	err    error
}

func newMockRepo() *mockRepo {
	return &mockRepo{nextID: 1}
}

func (m *mockRepo) FindAll(_ context.Context) ([]db.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.users, nil
}

func (m *mockRepo) FindByID(_ context.Context, id int32) (db.User, error) {
	if m.err != nil {
		return db.User{}, m.err
	}
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return db.User{}, identity.ErrNotFound
}

func (m *mockRepo) Save(_ context.Context, params db.CreateUserParams) (db.User, error) {
	if m.err != nil {
		return db.User{}, m.err
	}
	u := db.User{ID: m.nextID, Name: params.Name, Email: params.Email}
	m.nextID++
	m.users = append(m.users, u)
	return u, nil
}

func (m *mockRepo) Update(_ context.Context, params db.UpdateUserParams) (db.User, error) {
	if m.err != nil {
		return db.User{}, m.err
	}
	for i, u := range m.users {
		if u.ID == params.ID {
			m.users[i].Name = params.Name
			m.users[i].Email = params.Email
			return m.users[i], nil
		}
	}
	return db.User{}, identity.ErrNotFound
}

func (m *mockRepo) Delete(_ context.Context, id int32) error {
	if m.err != nil {
		return m.err
	}
	for i, u := range m.users {
		if u.ID == id {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return nil
}

func setupRouter(repo *mockRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := identity.NewIdentityService(repo)
	ctrl := NewUsersController(svc)
	ctrl.RegisterRoutes(r.Group("/users"))
	return r
}

func TestListUsers_OK(t *testing.T) {
	repo := newMockRepo()
	repo.users = []db.User{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
	}
	r := setupRouter(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp []UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(resp) != 1 || resp[0].Name != "Alice" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestListUsers_Empty(t *testing.T) {
	r := setupRouter(newMockRepo())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestListUsers_Error(t *testing.T) {
	repo := newMockRepo()
	repo.err = errors.New("db down")
	r := setupRouter(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestGetUser_OK(t *testing.T) {
	repo := newMockRepo()
	repo.users = []db.User{{ID: 5, Name: "Bob", Email: "bob@example.com"}}
	r := setupRouter(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/5", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.Name != "Bob" {
		t.Fatalf("expected Bob, got %s", resp.Name)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	r := setupRouter(newMockRepo())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetUser_InvalidID(t *testing.T) {
	r := setupRouter(newMockRepo())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/abc", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateUser_OK(t *testing.T) {
	r := setupRouter(newMockRepo())

	body, _ := json.Marshal(CreateUserRequest{Name: "Carol", Email: "carol@example.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.Name != "Carol" || resp.ID != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestCreateUser_MissingFields(t *testing.T) {
	r := setupRouter(newMockRepo())

	body, _ := json.Marshal(map[string]string{"name": "NoEmail"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateUser_Error(t *testing.T) {
	repo := newMockRepo()
	repo.err = errors.New("duplicate email")
	r := setupRouter(repo)

	body, _ := json.Marshal(CreateUserRequest{Name: "Carol", Email: "carol@example.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestUpdateUser_OK(t *testing.T) {
	repo := newMockRepo()
	repo.users = []db.User{{ID: 1, Name: "Old", Email: "old@example.com"}}
	r := setupRouter(repo)

	body, _ := json.Marshal(UpdateUserRequest{Name: "New", Email: "new@example.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.Name != "New" {
		t.Fatalf("expected New, got %s", resp.Name)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	r := setupRouter(newMockRepo())

	body, _ := json.Marshal(UpdateUserRequest{Name: "X", Email: "x@example.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateUser_InvalidID(t *testing.T) {
	r := setupRouter(newMockRepo())

	body, _ := json.Marshal(UpdateUserRequest{Name: "X", Email: "x@example.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/abc", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateUser_MissingFields(t *testing.T) {
	r := setupRouter(newMockRepo())

	body, _ := json.Marshal(map[string]string{"name": "NoEmail"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteUser_OK(t *testing.T) {
	repo := newMockRepo()
	repo.users = []db.User{{ID: 1, Name: "Del", Email: "del@example.com"}}
	r := setupRouter(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestDeleteUser_InvalidID(t *testing.T) {
	r := setupRouter(newMockRepo())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/abc", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteUser_Error(t *testing.T) {
	repo := newMockRepo()
	repo.err = errors.New("db error")
	r := setupRouter(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
