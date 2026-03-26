# Huma v2 Migration Plan

## Goal
Replace manual Gin handlers with [huma v2](https://huma.rocks/) to get auto-generated OpenAPI 3.1 spec and Swagger UI at `/docs`, similar to FastAPI.

## What stays the same
- Everything under `domain/` (services, repositories, SQL, generated code)
- Gin as the underlying HTTP router
- `app.go` composition root and wiring

## What changes

### 1. Add dependency
```
go get github.com/danielgtaylor/huma/v2
```
Uses the built-in `humaginv2` adapter to wrap the existing Gin engine.

### 2. `app.go` — update `Mount()`
- Create a `huma.API` instance wrapping the Gin engine via `humaginv2.New(r, huma.DefaultConfig(...))`
- Pass `huma.API` to controllers instead of `*gin.RouterGroup`

### 3. `controllers/users_controller.go` — full rewrite
- Change `RegisterRoutes(r *gin.RouterGroup)` → `RegisterRoutes(api huma.API, prefix string)`
- Convert each handler to a `huma.Register(...)` call with typed input/output structs
- Path params, query params, and request body declared via struct field tags
- Replace `c.JSON(...)` error responses with huma error helpers

### 4. `controllers/products_controller.go` — full rewrite
Same changes as users controller.

### 5. Error mapping
| Before | After |
|---|---|
| `c.JSON(http.StatusNotFound, ...)` | `huma.Error404NotFound("...", nil)` |
| `c.JSON(http.StatusBadRequest, ...)` | `huma.Error400BadRequest("...", nil)` |
| `c.JSON(http.StatusInternalServerError, ...)` | `huma.Error500InternalServerError("...", nil)` |

### 6. `main.go`
- No changes needed — `/test` stays as a plain Gin route

## Handler pattern (before → after)

**Before:**
```go
func (uc *UsersController) RegisterRoutes(r *gin.RouterGroup) {
    r.GET("/:id", uc.get)
}

func (uc *UsersController) get(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    ...
    c.JSON(http.StatusOK, toUserResponse(user))
}
```

**After:**
```go
func (uc *UsersController) RegisterRoutes(api huma.API, prefix string) {
    huma.Register(api, huma.Operation{
        Method:  http.MethodGet,
        Path:    prefix + "/{id}",
        Summary: "Get user",
        Tags:    []string{"users"},
    }, func(ctx context.Context, input *struct {
        ID int `path:"id"`
    }) (*struct{ Body UserResponse }, error) {
        user, err := uc.service.GetUser(input.ID)
        if errors.Is(err, identity.ErrNotFound) {
            return nil, huma.Error404NotFound("user not found", nil)
        }
        if err != nil {
            return nil, huma.Error500InternalServerError("internal error", nil)
        }
        resp := toUserResponse(user)
        return &struct{ Body UserResponse }{Body: resp}, nil
    })
}
```

## File change summary
| File | Change |
|---|---|
| `go.mod` / `go.sum` | Add `github.com/danielgtaylor/huma/v2` |
| `app.go` | Update `Mount()` to create and pass `huma.API` |
| `controllers/users_controller.go` | Full rewrite |
| `controllers/products_controller.go` | Full rewrite |
| `main.go` | No changes |
| `domain/**` | No changes |

## Result
- Swagger UI at `http://localhost:8080/docs`
- OpenAPI 3.1 JSON spec at `http://localhost:8080/openapi.json`
- Request validation handled automatically by huma
- No comment annotations required
