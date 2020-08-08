-- +migrate Up
create table product_yearly_price
(
	id bigserial not null
		constraint product_yearly_price_pk
			primary key,
	product_id bigint not null
		constraint product_yearly_price_products_id_fk
			references products,
	price numeric not null,
	created_at timestamptz not null,
	updated_at timestamptz not null
);

create index product_yearly_price_product_id_index
	on product_yearly_price (product_id);

-- +migrate Down
drop table product_yearly_price;
