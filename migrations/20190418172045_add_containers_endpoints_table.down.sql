ALTER TABLE containers
    ADD COLUMN endpoint VARCHAR(255) NOT NULL;

DROP TABLE containers_endpoints;
