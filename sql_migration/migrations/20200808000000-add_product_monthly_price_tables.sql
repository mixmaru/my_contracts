-- +migrate Up
create table product_monthly_price
(
	id bigint not null
		constraint product_monthly_price_pk
			primary key,
	product_id bigint not null
		constraint product_monthly_price_products_id_fk
			references products (id),
	price decimal not null
);
alter table product_monthly_price
	add created_at timestamptz not null;

alter table product_monthly_price
	add updated_at timestamptz not null;

create sequence product_monthly_price_id_seq;

alter table product_monthly_price alter column id set default nextval('public.product_monthly_price_id_seq');

alter sequence product_monthly_price_id_seq owned by product_monthly_price.id;


insert into product_monthly_price (product_id, price, created_at, updated_at)
select id, price, created_at, created_at
from products;

alter table products drop column price;


-- +migrate Down
alter table products
	add price numeric default 0 not null;

update products p set price = pm.price from product_monthly_price pm where pm.product_id = p.id;

alter table products alter column price drop default;
drop table product_monthly_price;
