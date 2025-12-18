-- +goose Up
ALTER TABLE treatments
ALTER COLUMN price TYPE DECIMAL(10,2);

-- +goose Down
ALTER TABLE treatments
ALTER COLUMN price TYPE FLOAT;
