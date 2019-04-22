    CREATE TYPE plans_type AS ENUM ('container', 'image');

CREATE TABLE plans
(
    id         UUID         NOT NULL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    type       plans_type   NOT NULL,
    price      FLOAT        NOT NULL,
    metadata   JSON         NOT NULL,
    public     BOOLEAN      NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE plans_authorizations
(
    plan_id      UUID NOT NULL REFERENCES plans (id),
    namespace_id UUID NOT NULL REFERENCES accounts (id)
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
INSERT INTO plans (id, name, type, price, metadata, public)
VALUES  (uuid_generate_v4(), 'Small', 'container', 4, '{"cpu": 1, "memory": 256}', true),
        (uuid_generate_v4(), 'Standard', 'container', 7, '{"cpu": 1, "memory": 512}', true),
        (uuid_generate_v4(), 'Standard X', 'container', 12, '{"cpu": 2, "memory": 1024}', true),
        (uuid_generate_v4(), 'Performance', 'container', 21, '{"cpu": 4, "memory": 2048}', false),
        (uuid_generate_v4(), 'Performance X', 'container', 42, '{"cpu": 8, "memory": 4096}', false);
