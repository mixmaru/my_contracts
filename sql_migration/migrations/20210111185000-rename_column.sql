-- +migrate Up
alter table customers_customer_properties rename column upadted_at to updated_at;

-- +migrate Down
alter table customers_customer_properties rename column updated_at to upadted_at;
