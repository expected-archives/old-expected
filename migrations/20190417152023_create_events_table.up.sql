CREATE TYPE event_resource AS ENUM ('account', 'container', 'image');
CREATE TYPE event_action AS ENUM ('create', 'update', 'delete');
CREATE TYPE event_issuer AS ENUM ('robot', 'account');

CREATE TABLE events
(
    id          UUID           NOT NULL PRIMARY KEY,
    resource    event_resource NOT NULL,
    resource_id UUID,
    action      event_action   NOT NULL,
    issuer      event_issuer   NOT NULL,
    issuer_id   UUID,
    metadata    JSON           NOT NULL,
    created_at  TIMESTAMP      NOT NULL DEFAULT NOW()
)
