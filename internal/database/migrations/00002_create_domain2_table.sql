-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS domain2 (
    id         BIGSERIAL    PRIMARY KEY,
    name       TEXT         NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS domain2;
-- +goose StatementEnd
