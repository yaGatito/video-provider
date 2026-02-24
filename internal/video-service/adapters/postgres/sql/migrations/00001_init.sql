-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS videos (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        publisherid UUID NOT NULL,
        topic text NOT NULL,
        description text,
        createdAt time,
        status varchar(50)
    );

CREATE TABLE
    IF NOT EXISTS comments (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        publisherid UUID NOT NULL,
        comment text NOT NULL,
        createdAt time,
        status varchar(50)
    );

-- Insert initial test data for videos
INSERT INTO videos (publisherid, topic, description, createdAt, status) VALUES
('123e4567-e89b-12d3-a456-426614174000', 'Tech Talk', 'A discussion on the latest tech trends.', CURRENT_TIME, 'active'),
('123e4567-e89b-12d3-a456-426614174001', 'Travel Vlogs', 'Exploring new places around the world.', CURRENT_TIME, 'active'),
('123e4567-e89b-12d3-a456-426614174002', 'Food Reviews', 'Tasting and reviewing delicious dishes.', CURRENT_TIME, 'active');

-- Insert initial test data for comments
INSERT INTO comments (publisherid, comment, createdAt, status) VALUES
('123e4567-e89b-12d3-a456-426614174000', 'Great video!', CURRENT_TIME, 'active'),
('123e4567-e89b-12d3-a456-426614174001', 'Love the travel tips.', CURRENT_TIME, 'active'),
('123e4567-e89b-12d3-a456-426614174002', 'Yummy food!', CURRENT_TIME, 'active');

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS videos;

DROP TABLE IF EXISTS comments;

-- +goose StatementEnd