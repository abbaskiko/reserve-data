DROP TABLE fetch_data;

CREATE TABLE "fetch_data"
(
    id SERIAL PRIMARY KEY,
    created TIMESTAMPTZ NOT NULL,
    data JSON NOT NULL,
    type fetch_data_type NOT NULL
);
CREATE INDEX "fetch_data_created_index" ON "fetch_data" (created);