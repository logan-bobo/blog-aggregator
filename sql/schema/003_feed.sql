-- +goose Up
CREATE TABLE feeds(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL, 
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL, 
    name VARCHAR(50) NOT NULL,
    url VARCHAR(250) NOT NULL, 
    UNIQUE(url),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id) 
            REFERENCES users(id)
            ON DELETE CASCADE
);
