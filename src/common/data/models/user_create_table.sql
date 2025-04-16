-- https://jonmeyers.io/blog/automatically-generate-values-for-created-and-updated-columns-in-postgres/

CREATE TABLE IF NOT EXISTS %[1]s (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  name_first TEXT,
  name_last TEXT,
  email TEXT
);

-- ALTER TABLE %[1]s
--   ALTER COLUMN created_at
--     SET DEFAULT NOW();

-- CREATE OR REPLACE FUNCTION set_updated_at_columns() RETURNS TRIGGER AS $$
--   BEGIN
--     NEW.created_at = NOW();
--     RETURN NEW;
--   END;
-- && LANGUAGE plpgsql;

-- CREATE TRIGGER on_create_set_updated_at_columns
-- BEFORE UPDATE ON %[1]s
-- FOR EACH ROW EXECUTE PROCEDURE set_updated_at_columns();
