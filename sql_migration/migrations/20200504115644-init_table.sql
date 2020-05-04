
-- +migrate Up
create table users
(
    id bigint,
    user_type smallint
);

-- +migrate Down
DROP TABLE users;
