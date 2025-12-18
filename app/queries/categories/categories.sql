-- name: CreateCategory :one
INSERT INTO categories (name_category) VALUES ($1) RETURNING *;

-- name: FindCategories :many
SELECT id, name_category, created_at, updated_at FROM categories;

-- name: FindCategoryById :one
SELECT id, name_category, created_at, updated_at FROM categories WHERE id = $1;

-- name: FindCategoryByName :one
SELECT id, name_category, created_at, updated_at FROM categories WHERE name_category = $1;

-- name: UpdateCategory :one
UPDATE categories SET name_category = $2 WHERE id = $1 RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;
