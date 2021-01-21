-- +migrate Up
alter table customer_properties alter column type type smallint using type::smallint;

-- +migrate Down
alter table customer_properties alter column type type text using type::text;
