ALTER TABLE containers
    DROP COLUMN endpoint;

CREATE TABLE containers_endpoints
(
    id           UUID         NOT NULL PRIMARY KEY,
    container_id UUID         NOT NULL REFERENCES containers (id),
    endpoint     VARCHAR(255) NOT NULL,
    is_default   BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMP    NOT NULL DEFAULT NOW()
);
