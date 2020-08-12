-- +migrate Up
create table product_price_monthlies
(
	product_id bigint not null
		constraint product_price_monthlies_products_id_fk
			references products (id),
	price numeric not null,
	created_at timestamptz not null,
	updated_at timestamptz not null
);
create unique index product_price_monthlies_product_id_uindex
	on product_price_monthlies (product_id);

insert into product_price_monthlies (product_id, price, created_at, updated_at) select id, price, created_at, updated_at from products;
alter table products drop column price;

create table product_price_yearlies
(
	product_id bigint not null
		constraint product_price_yearlies_products_id_fk
			references products (id),
	price numeric not null,
	created_at timestamptz not null,
	updated_at timestamptz not null
);
create unique index product_price_yearlies_product_id_uindex
	on product_price_yearlies (product_id);

create table product_price_lumps
(
	product_id bigint not null
		constraint product_price_lump_products_id_fk
			references products (id),
	price numeric not null,
	created_at timestamptz not null,
	updated_at timestamptz not null
);
create unique index product_price_lumps_product_id_uindex
	on product_price_lumps (product_id);

create table product_price_custom_terms
(
	product_id bigint not null
		constraint product_price_custom_terms_products_id_fk
			references products (id),
	price numeric not null,
	term int not null,
	created_at timestamptz not null,
	updated_at timestamptz not null
);
create unique index product_price_custom_terms_product_id_uindex
	on product_price_custom_terms (product_id);

-- +migrate Down
drop table product_price_monthlies;
drop table product_price_yearlies;
drop table product_price_lumps;
drop table product_price_custom_terms;
