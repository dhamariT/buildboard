-- Create application user
CREATE USER buildboard_app WITH PASSWORD 'changeme_in_production';

-- Grant privileges on database
GRANT ALL PRIVILEGES ON DATABASE buildboard_db TO buildboard_app;

-- Grant privileges on schema
GRANT ALL PRIVILEGES ON SCHEMA public TO buildboard_app;

-- Grant privileges on all existing tables
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO buildboard_app;

-- Grant privileges on all existing sequences
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO buildboard_app;

-- Grant privileges on future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO buildboard_app;

-- Grant privileges on future sequences
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON SEQUENCES TO buildboard_app;