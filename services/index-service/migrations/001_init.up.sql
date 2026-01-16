-- Index Service Initial Migration

CREATE TYPE index_status AS ENUM ('pending', 'indexed', 'not_indexed', 'error');

CREATE TABLE platforms (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    url TEXT NOT NULL,
    domain VARCHAR(255),
    index_status index_status DEFAULT 'pending',
    is_indexed BOOLEAN DEFAULT FALSE,
    first_indexed_at TIMESTAMP,
    last_checked_at TIMESTAMP,
    check_count INTEGER DEFAULT 0,
    potential_score INTEGER DEFAULT 0,
    is_must_have BOOLEAN DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_platforms_user_id ON platforms(user_id);
CREATE INDEX idx_platforms_index_status ON platforms(index_status);
CREATE INDEX idx_platforms_domain ON platforms(domain);
CREATE INDEX idx_platforms_is_indexed ON platforms(is_indexed);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_platforms_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_platforms_updated_at
    BEFORE UPDATE ON platforms
    FOR EACH ROW
    EXECUTE FUNCTION update_platforms_updated_at();
