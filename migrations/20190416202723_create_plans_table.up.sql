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

CREATE TABLE custom_plans
(
    plan_id      UUID NOT NULL REFERENCES plans (id),
    namespace_id UUID NOT NULL REFERENCES accounts (id)
)