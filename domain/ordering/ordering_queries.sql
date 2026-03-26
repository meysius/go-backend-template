-- name: GetProduct :one
SELECT * FROM products WHERE id = $1;

-- name: ListProducts :many
SELECT * FROM products ORDER BY id;

-- name: CreateProduct :one
INSERT INTO products (name, price) VALUES ($1, $2) RETURNING *;

-- name: UpdateProduct :one
UPDATE products SET name = $2, price = $3 WHERE id = $1 RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;
