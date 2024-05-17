-- +goose Up
CREATE TABLE posts(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title VARCHAR(250),
    url VARCHAR(250),
    description VARCHAR(500),
    published_at TIMESTAMP,
    feed_id UUID, 
    CONSTRAINT unique_url UNIQUE (url),
    CONSTRAINT fk_feed
        FOREIGN KEY(feed_id)
            REFERENCES feeds(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;
