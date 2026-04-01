package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-starter-template/db"
	"go-starter-template/domain/ordering"
)

type ProductResponse struct {
	ID    int32   `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func toProductResponse(p db.Product) ProductResponse {
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

// list godoc
// @Summary      List products
// @Tags         products
// @Produce      json
// @Success      200  {array}   ProductResponse
// @Failure      500  {object}  map[string]string
// @Router       /products [get]
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

// get godoc
// @Summary      Get a product by ID
// @Tags         products
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  ProductResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /products/{id} [get]
func (pc *ProductsController) get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	product, err := pc.service.GetProduct(int32(id))
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

// CreateProductRequest is the request body for creating a product.
type CreateProductRequest struct {
	Name  string  `json:"name"  binding:"required"`
	Price float64 `json:"price" binding:"required"`
}

// create godoc
// @Summary      Create a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        body  body      CreateProductRequest  true  "Product payload"
// @Success      201   {object}  ProductResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /products [post]
func (pc *ProductsController) create(c *gin.Context) {
	var body CreateProductRequest
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

// UpdateProductRequest is the request body for updating a product.
type UpdateProductRequest struct {
	Name  string  `json:"name"  binding:"required"`
	Price float64 `json:"price" binding:"required"`
}

// update godoc
// @Summary      Update a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id    path      int                  true  "Product ID"
// @Param        body  body      UpdateProductRequest  true  "Product payload"
// @Success      200   {object}  ProductResponse
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /products/{id} [put]
func (pc *ProductsController) update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var body UpdateProductRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := pc.service.UpdateProduct(int32(id), body.Name, body.Price)
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

// delete godoc
// @Summary      Delete a product
// @Tags         products
// @Param        id   path  int  true  "Product ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /products/{id} [delete]
func (pc *ProductsController) delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := pc.service.DeleteProduct(int32(id)); err != nil {
		if errors.Is(err, ordering.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
