DROP TABLE fetch_data;

CREATE TABLE "fetch_data"
(
    id SERIAL,
    created TIMESTAMPTZ NOT NULL,
    data BYTEA NOT NULL,
    type fetch_data_type NOT NULL,
    PRIMARY KEY (id,created)
) PARTITION BY RANGE (created);

CREATE TABLE fetch_data_default PARTITION OF fetch_data DEFAULT;
CREATE INDEX "fetch_data_created_index" ON "fetch_data" (created);