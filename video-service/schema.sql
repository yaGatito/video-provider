CREATE TABLE
  videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    publisherid UUID NOT NULL,
    topic text NOT NULL,
    description text,
    createdAt time,
    status string
  );

CREATE TABLE
  comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    publisherid UUID NOT NULL,
    comment text NOT NULL,
    createdAt time,
    status string
  );