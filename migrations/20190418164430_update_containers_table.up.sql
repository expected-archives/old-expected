DROP TABLE containers;

CREATE TYPE container_state AS ENUM ('stopped', 'starting', 'running');

CREATE TABLE containers
(
    id           UUID            NOT NULL PRIMARY KEY,
    name         VARCHAR(255)    NOT NULL,
    image        VARCHAR(255)    NOT NULL,
    endpoint     VARCHAR(255)    NOT NULL,
    plan_id      UUID            NOT NULL REFERENCES plans (id),
    environment  JSON            NOT NULL,
    tags         JSON            NOT NULL,
    namespace_id UUID            NOT NULL,
    state        container_state NOT NULL DEFAULT 'stopped',
    created_at   TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP       NOT NULL DEFAULT NOW()
);
