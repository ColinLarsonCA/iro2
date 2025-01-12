CREATE TABLE collabs (
    id UUID PRIMARY KEY,
    source TEXT NOT NULL,
    source_url TEXT NOT NULL,
    source_posted_at TIMESTAMP NOT NULL,
    collab_jp JSONB NOT NULL,
    collab_en JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_source_url ON collabs (source_url ASC);
CREATE INDEX source_posted_at ON collabs (source_posted_at DESC);