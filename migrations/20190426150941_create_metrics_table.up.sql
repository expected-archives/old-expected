CREATE TABLE metrics
(
    time         TIMESTAMPTZ      NOT NULL,
    id           TEXT             NOT NULL,
    memory       BIGINT           NOT NULL,
    net_input    DOUBLE PRECISION NOT NULL,
    net_output   DOUBLE PRECISION NOT NULL,
    block_input  BIGINT           NOT NULL,
    block_output BIGINT           NOT NULL,
    cpu          DOUBLE PRECISION NULL
);

SELECT create_hypertable('metrics', 'time');
CREATE INDEX ON metrics (id, time DESC);