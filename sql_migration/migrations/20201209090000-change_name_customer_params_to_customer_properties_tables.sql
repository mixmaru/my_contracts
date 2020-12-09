-- +migrate Up
alter table customer_params rename to customer_properties;
alter table customer_types_customer_params rename to customer_types_customer_properties;
alter table customers_customer_params rename to customers_customer_properties;

alter table customer_types_customer_properties rename column customer_param_id to customer_property_id;
alter table customers_customer_properties rename column customer_param_id to customer_property_id;

-- +migrate Down
alter table customers_customer_properties rename column customer_property_id to customer_param_id;
alter table customer_types_customer_properties rename column customer_property_id to customer_param_id;

alter table customers_customer_properties rename to customers_customer_params;
alter table customer_types_customer_properties rename to customer_types_customer_params;
alter table customer_properties rename to customer_params;
