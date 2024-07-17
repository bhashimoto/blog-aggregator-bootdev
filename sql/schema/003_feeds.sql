-- +goose Up
CREATE TABLE feeds (
	name VARCHAR(255) NOT NULL,
	url TEXT UNIQUE NOT NULL,
	user_id UUID NOT NULL 
		REFERENCES users(id) 
		ON DELETE CASCADE
);


-- +goose Down
DROP TABLE feeds;
