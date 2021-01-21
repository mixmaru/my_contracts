-- +migrate Up
alter table customer_types_customer_properties
    add "order" smallint not null;

-- +migrate Down
alter table customer_types_customer_properties drop column "order";
