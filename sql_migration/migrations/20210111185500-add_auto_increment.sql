-- +migrate Up
create sequence customers_id_seq;
alter table customers alter column id set default nextval('public.customers_id_seq');
alter sequence customers_id_seq owned by customers.id;

-- +migrate Down
alter table customers alter column id drop default;
drop sequence customers_id_seq;
