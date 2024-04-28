-- +goose Up
ALTER TABLE users
ADD api_key VARCHAR(64) NOT NULL
DEFAULT encode(sha256(random()::text::bytea), 'hex');

ALTER TABLE users
ADD CONSTRAINT unique_api_key UNIQUE (api_key);

-- +goose Down
ALTER TABLE users
DROP COLUMN api_key;
