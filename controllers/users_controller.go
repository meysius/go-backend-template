package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-starter-template/db"
	"go-starter-template/domain/identity"
)

type UserResponse struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func toUserResponse(u db.User) UserResponse {
	return UserResponse{ID: u.ID, Name: u.Name, Email: u.Email}
}

type UsersController struct {
	service *identity.IdentityService
}

func NewUsersController(service *identity.IdentityService) *UsersController {
	return &UsersController{service: service}
}

func (uc *UsersController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", uc.list)
	r.GET("/:id", uc.get)
	r.POST("", uc.create)
	r.PUT("/:id", uc.update)
	r.DELETE("/:id", uc.delete)
}

// list godoc
// @Summary      List users
// @Tags         users
// @Produce      json
// @Success      200  {array}   UserResponse
// @Failure      500  {object}  map[string]string
// @Router       /users [get]
func (uc *UsersController) list(c *gin.Context) {
	users, err := uc.service.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]UserResponse, len(users))
	for i, u := range users {
		resp[i] = toUserResponse(u)
	}
	c.JSON(http.StatusOK, resp)
}

// get godoc
// @Summary      Get a user by ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  UserResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id} [get]
func (uc *UsersController) get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := uc.service.GetUser(int32(id))
	if err != nil {
		if errors.Is(err, identity.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(user))
}

// CreateUserRequest is the request body for creating a user.
type CreateUserRequest struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required"`
}

// create godoc
// @Summary      Create a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body      CreateUserRequest  true  "User payload"
// @Success      201   {object}  UserResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /users [post]
func (uc *UsersController) create(c *gin.Context) {
	var body CreateUserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.service.CreateUser(body.Name, body.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, toUserResponse(user))
}

// UpdateUserRequest is the request body for updating a user.
type UpdateUserRequest struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required"`
}

// update godoc
// @Summary      Update a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      int               true  "User ID"
// @Param        body  body      UpdateUserRequest  true  "User payload"
// @Success      200   {object}  UserResponse
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /users/{id} [put]
func (uc *UsersController) update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var body UpdateUserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.service.UpdateUser(int32(id), body.Name, body.Email)
	if err != nil {
		if errors.Is(err, identity.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(user))
}

// delete godoc
// @Summary      Delete a user
// @Tags         users
// @Param        id   path  int  true  "User ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id} [delete]
func (uc *UsersController) delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := uc.service.DeleteUser(int32(id)); err != nil {
		if errors.Is(err, identity.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
