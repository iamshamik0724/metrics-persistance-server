CREATE TABLE api_metrics (
    time TIMESTAMPTZ NOT NULL,
    route TEXT NOT NULL,
    method VARCHAR(50) NOT NULL,
    status_code INT NOT NULL,
    response_time FLOAT8 NOT NULL,
    PRIMARY KEY (time, route, method)
);

-- Convert the table to a hypertable
SELECT create_hypertable('api_metrics', 'time');