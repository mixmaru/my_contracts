-- +migrate Up
alter table customers_customer_properties
    add constraint customers_customer_properties_pk
        primary key (customer_id, customer_property_id);

alter table customer_types_customer_properties
    add constraint customer_types_customer_properties_pk
        primary key (customer_type_id, customer_property_id);

-- +migrate Down
alter table customer_types_customer_properties drop constraint customer_types_customer_properties_pk;
alter table customers_customer_properties drop constraint customers_customer_properties_pk;
