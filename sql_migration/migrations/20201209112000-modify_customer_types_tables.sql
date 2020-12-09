-- +migrate Up
create sequence customer_types_id_seq;
alter table customer_types alter column id set default nextval('public.customer_types_id_seq');
alter sequence customer_types_id_seq owned by customer_types.id;

create sequence customer_properties_id_seq;
alter table customer_properties alter column id set default nextval('public.customer_properties_id_seq');
alter sequence customer_properties_id_seq owned by customer_properties.id;

-- +migrate Down
alter table customer_properties alter column id drop default;
drop sequence customer_properties_id_seq;

alter table customer_types alter column id drop default;
drop sequence customer_types_id_seq;
