CREATE TABLE containers
(
    id          UUID         NOT NULL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    image       VARCHAR(255) NOT NULL,
    endpoint    VARCHAR(255) NOT NULL,
    memory      INT          NOT NULL,
    environment JSON         NOT NULL,
    tags        JSON         NOT NULL,
    owner_id    UUID         NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW()
);
