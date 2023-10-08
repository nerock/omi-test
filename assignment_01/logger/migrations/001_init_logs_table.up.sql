CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    type TEXT,
    timestamp INTEGER,
    user_ip TEXT,
    data JSONB
);