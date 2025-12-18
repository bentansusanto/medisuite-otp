-- +goose Up

-- Add deleted_at column for soft delete (only if it doesn't exist)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'users' AND column_name = 'deleted_at'
    ) THEN
        ALTER TABLE users ADD COLUMN deleted_at TIMESTAMPTZ NULL;
    END IF;
END $$;

-- Drop the old UNIQUE constraint on email (only if it exists)
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'users_email_key'
    ) THEN
        ALTER TABLE users DROP CONSTRAINT users_email_key;
    END IF;
END $$;

-- Create partial unique index for email (only for non-deleted users)
-- This allows the same email to be reused after soft delete
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_active
ON users(email)
WHERE deleted_at IS NULL;

-- +goose Down

-- Drop the partial index
DROP INDEX IF EXISTS idx_users_email_active;

-- Restore the UNIQUE constraint on email
ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);

-- Drop the deleted_at column
ALTER TABLE users DROP COLUMN deleted_at;
