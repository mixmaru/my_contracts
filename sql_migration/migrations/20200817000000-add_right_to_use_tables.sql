-- +migrate Up
alter table right_to_use alter column valid_from type timestamptz using valid_from::timestamptz;
alter table right_to_use alter column valid_to type timestamptz using valid_to::timestamptz;

-- +migrate Down
alter table right_to_use alter column valid_from type timestamp using valid_from::timestamp;
alter table right_to_use alter column valid_to type timestamp using valid_to::timestamp;
