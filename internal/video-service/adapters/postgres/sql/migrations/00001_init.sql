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

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS videos;

DROP TABLE IF EXISTS comments;

-- +goose StatementEnd