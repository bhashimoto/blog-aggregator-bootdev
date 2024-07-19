-- +goose Up
CREATE TABLE posts (
	id UUID UNIQUE NOT NULL,
	created_at	TIMESTAMP NOT NULL,
	updated_at	TIMESTAMP NOT NULL,
	title		TEXT,
	url		TEXT UNIQUE,
	description	TEXT,
	published_at	TIMESTAMP,
	feed_id		UUID REFERENCES(feeds) ON CASCADE DELETE
);

-- +goose Down
DROP TABLE posts;
