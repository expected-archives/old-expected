CREATE TABLE images
(
    id           UUID         NOT NULL PRIMARY KEY,
    namespace_id UUID         NOT NULL,
    digest       TEXT         NOT NULL,
    name         VARCHAR(255) NOT NULL,
    tag          VARCHAR(255) NOT NULL,
    created_at   TIMESTAMP DEFAULT now()
);

CREATE TABLE layers
(
    digest     TEXT   NOT NULL PRIMARY KEY,
    repository TEXT   NOT NULL,
    size       BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE image_layer
(
    image_id     UUID NOT NULL REFERENCES images (id),
    layer_digest TEXT NOT NULL REFERENCES layers (digest),
    created_at   TIMESTAMP DEFAULT now()
);
