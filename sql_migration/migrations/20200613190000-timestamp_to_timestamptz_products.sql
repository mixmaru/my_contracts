
-- +migrate Up
alter table products alter column created_at type timestamptz using created_at::timestamptz;
alter table products alter column updated_at type timestamptz using updated_at::timestamptz;

-- +migrate Down
alter table products alter column created_at type timestamp using created_at::timestamp;
alter table products alter column updated_at type timestamp using updated_at::timestamp;
