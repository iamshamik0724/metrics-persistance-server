CREATE DATABASE metrics_persistance_server;

CREATE USER metrics_persistance_server_user WITH PASSWORD 'paswword123';

GRANT ALL PRIVILEGES ON DATABASE metrics_persistance_server TO metrics_persistance_server_user;

GRANT USAGE ON SCHEMA public TO metrics_persistance_server_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO metrics_persistance_server_user;

CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

