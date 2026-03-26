package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-starter-template/domain/ordering"
)

type ProductResponse struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func toProductResponse(p ordering.Product) ProductResponse {
	return ProductResponse{ID: p.ID, Name: p.Name, Price: p.Price}
}

type ProductsController struct {
	service *ordering.OrderingService
}

func NewProductsController(service *ordering.OrderingService) *ProductsController {
	return &ProductsController{service: service}
}

func (pc *ProductsController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", pc.list)
	r.GET("/:id", pc.get)
	r.POST("", pc.create)
	r.PUT("/:id", pc.update)
	r.DELETE("/:id", pc.delete)
}

func (pc *ProductsController) list(c *gin.Context) {
	products, err := pc.service.ListProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]ProductResponse, len(products))
	for i, p := range products {
		resp[i] = toProductResponse(p)
	}
	c.JSON(http.StatusOK, resp)
}

func (pc *ProductsController) get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	product, err := pc.service.GetProduct(id)
	if err != nil {
		if errors.Is(err, ordering.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toProductResponse(product))
}

func (pc *ProductsController) create(c *gin.Context) {
	var body struct {
		Name  string  `json:"name" binding:"required"`
		Price float64 `json:"price" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := pc.service.CreateProduct(body.Name, body.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, toProductResponse(product))
}

func (pc *ProductsController) update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var body struct {
		Name  string  `json:"name" binding:"required"`
		Price float64 `json:"price" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := pc.service.UpdateProduct(id, body.Name, body.Price)
	if err != nil {
		if errors.Is(err, ordering.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toProductResponse(product))
}

func (pc *ProductsController) delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := pc.service.DeleteProduct(id); err != nil {
		if errors.Is(err, ordering.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
