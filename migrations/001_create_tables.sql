CREATE TABLE IF NOT EXISTS accounts
(
    id                  UUID         NOT NULL PRIMARY KEY,
    name                VARCHAR(255) NOT NULL,
    email               VARCHAR(255) NOT NULL,
    avatar_url          VARCHAR(255) NOT NULL,
    github_id           BIGINT       NOT NULL,
    github_access_token VARCHAR(255) NOT NULL,
    api_key             VARCHAR(32)  NOT NULL,
    admin               BOOLEAN   DEFAULT FALSE,
    created_at          TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS containers
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
