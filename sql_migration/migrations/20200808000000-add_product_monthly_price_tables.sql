-- +migrate Up
create table product_price_monthlies
(
	product_id bigint not null
		constraint product_price_monthlies_products_id_fk
			references products (id),
	created_at timestamptz not null,
	updated_at timestamptz not null
);
insert into product_price_monthlies (product_id, created_at, updated_at) select id, created_at, updated_at from products;

create table product_price_yearlies
(
	product_id bigint not null
		constraint product_price_yearlies_products_id_fk
			references products (id),
	created_at timestamptz not null,
	updated_at timestamptz not null
);

create table product_price_lumps
(
	product_id bigint not null
		constraint product_price_lump_products_id_fk
			references products (id),
	created_at timestamptz not null,
	updated_at timestamptz not null
);

create table product_price_custom_terms
(
	product_id bigint not null
		constraint product_price_lump_products_id_fk
			references products (id),
	term int not null,
	created_at timestamptz not null,
	updated_at timestamptz not null
);

-- +migrate Down
drop table product_price_monthlies;
drop table product_price_yearlies;
drop table product_price_lumps;
drop table product_price_custom_terms;
