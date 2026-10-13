-- Add "status" column to "users" table; only "active" and "disabled" are allowed
ALTER TABLE "users"
  ADD COLUMN "status" text NOT NULL DEFAULT 'active',
  ADD CONSTRAINT "users_status_check" CHECK ("status" IN ('active', 'disabled'));
