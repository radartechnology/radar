CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    token TEXT NOT NULL,
    username VARCHAR(255) NOT NULL
);
