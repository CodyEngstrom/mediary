-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Define enum for media types
CREATE TYPE media_type AS ENUM ('movie', 'tv_show', 'anime', 'book', 'game');

-- Create media table
CREATE TABLE media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
    media_type media_type NOT NULL,
    release_date DATE,
    language VARCHAR(10),
    external_id TEXT,
    external_source VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- auto-update updated_at
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_set_updated_at
BEFORE UPDATE ON media
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();