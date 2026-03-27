CREATE TABLE IF NOT EXISTS links (
    id BIGSERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    code VARCHAR(20) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);
CREATE INDEX idx_links_code ON links(code);