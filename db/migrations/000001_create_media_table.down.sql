-- Drop trigger and function
DROP TRIGGER IF EXISTS trg_set_updated_at On media;
DROP FUNCTION IF EXISTS set_updated_at;

-- Drop table
DROP TABLE IF EXISTS media;

-- Drop enum
DROP TYPE IF EXISTS media_type;