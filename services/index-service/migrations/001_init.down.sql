-- Index Service Rollback Migration

DROP TRIGGER IF EXISTS trigger_platforms_updated_at ON platforms;
DROP FUNCTION IF EXISTS update_platforms_updated_at();
DROP TABLE IF EXISTS platforms;
DROP TYPE IF EXISTS index_status;
