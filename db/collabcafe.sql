CREATE TABLE collabs (
    id UUID PRIMARY KEY,
    source TEXT NOT NULL,
    source_url TEXT NOT NULL,
    source_posted_at TIMESTAMP NOT NULL,
    collab_ja JSONB NOT NULL,
    collab_en JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX idx_source_url ON collabs (source_url ASC);
CREATE INDEX idx_source_posted_at ON collabs (source_posted_at DESC);

CREATE TABLE ja_to_en_lookup (
    ja TEXT PRIMARY KEY,
    en TEXT NOT NULL
);

ALTER TABLE collabs
    ADD COLUMN fts_collab_en tsvector
    generated always as (to_tsvector (
        'english', coalesce(collab_en->'content'->>'series', '') || ' ' || coalesce(collab_en->'summary'->>'title', '') || ' ' || coalesce(collab_en->'summary'->>'description', '')
    ))
    stored;

CREATE INDEX collabs_fts_collab_en_idx
    ON collabs
    USING GIN (fts_collab_en);
