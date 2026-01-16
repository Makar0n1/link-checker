-- Health Service Initial Migration

CREATE TABLE monitored_sites (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    url TEXT NOT NULL,
    domain VARCHAR(255),
    http_status INTEGER,
    is_alive BOOLEAN DEFAULT FALSE,
    response_time_ms INTEGER,
    allows_indexing BOOLEAN,
    robots_txt_status VARCHAR(50),
    has_noindex BOOLEAN DEFAULT FALSE,
    pages_indexed INTEGER DEFAULT 0,
    last_checked_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE site_check_history (
    id BIGSERIAL PRIMARY KEY,
    site_id BIGINT REFERENCES monitored_sites(id) ON DELETE CASCADE,
    http_status INTEGER,
    is_alive BOOLEAN,
    response_time_ms INTEGER,
    checked_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_sites_user_id ON monitored_sites(user_id);
CREATE INDEX idx_sites_domain ON monitored_sites(domain);
CREATE INDEX idx_sites_is_alive ON monitored_sites(is_alive);
CREATE INDEX idx_history_site_id ON site_check_history(site_id);
CREATE INDEX idx_history_checked_at ON site_check_history(checked_at);
