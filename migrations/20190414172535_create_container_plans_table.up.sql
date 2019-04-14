CREATE TABLE container_plans
(
    id        UUID         NOT NULL PRIMARY KEY,
    name      VARCHAR(255) NOT NULL,
    price     FLOAT        NOT NULL,
    cpu       INT          NOT NULL,
    memory    INT          NOT NULL,
    available BOOLEAN DEFAULT FALSE
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO container_plans (id, name, price, cpu, memory, available)
VALUES (uuid_generate_v4(), 'Small', 4, 1, 256, true),
       (uuid_generate_v4(), 'Standard', 7, 1, 512, true),
       (uuid_generate_v4(), 'Standard X', 12, 2, 1024, true),
       (uuid_generate_v4(), 'Performance', 21, 4, 2048, false),
       (uuid_generate_v4(), 'Performance X', 42, 8, 4096, false);
