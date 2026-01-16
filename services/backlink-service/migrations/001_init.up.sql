-- Create projects table
CREATE TABLE IF NOT EXISTS projects (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    user_id BIGINT NOT NULL,
    google_sheet_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for projects
CREATE INDEX IF NOT EXISTS idx_projects_user_id ON projects(user_id);

-- Create link_status enum type
DO $$ BEGIN
    CREATE TYPE link_status AS ENUM ('pending', 'active', 'broken', 'removed', 'nofollow');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Create link_type enum type
DO $$ BEGIN
    CREATE TYPE link_type AS ENUM ('dofollow', 'nofollow', 'sponsored', 'ugc');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Create backlinks table
CREATE TABLE IF NOT EXISTS backlinks (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    source_url TEXT NOT NULL,
    target_url TEXT NOT NULL,
    anchor_text VARCHAR(500) NOT NULL DEFAULT '',
    status link_status NOT NULL DEFAULT 'pending',
    link_type link_type NOT NULL DEFAULT 'dofollow',
    http_status SMALLINT,
    last_checked_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for backlinks
CREATE INDEX IF NOT EXISTS idx_backlinks_project_id ON backlinks(project_id);
CREATE INDEX IF NOT EXISTS idx_backlinks_status ON backlinks(status);
CREATE INDEX IF NOT EXISTS idx_backlinks_source_url ON backlinks(source_url);
CREATE INDEX IF NOT EXISTS idx_backlinks_target_url ON backlinks(target_url);
CREATE INDEX IF NOT EXISTS idx_backlinks_last_checked_at ON backlinks(last_checked_at);
