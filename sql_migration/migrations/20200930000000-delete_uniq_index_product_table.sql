-- +migrate Up
drop index products_name_uindex;

create index products_name_index
	on products (name);

-- +migrate Down
drop index products_name_index;

create unique index products_name_uindex
	on products (name);
