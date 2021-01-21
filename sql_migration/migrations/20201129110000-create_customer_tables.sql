-- +migrate Up
create table customer_types
(
	id int
		constraint customer_types_pk
			primary key,
	name text
);
alter table customer_types alter column name set not null;
alter table customer_types
	add created_at timestamptz not null;
alter table customer_types
	add updated_at timestamptz not null;
create unique index customer_types_name_uindex
	on customer_types (name);

create table customers
(
	id bigint
		constraint customers_pk
			primary key,
	customer_type_id int not null
		constraint customers_customer_types_id_fk
			references customer_types,
	name text not null,
	created_at timestamptz not null,
	updated_at timestamptz not null
);


-- create table customer_param_types
-- (
-- 	id int
-- 		constraint customer_param_types_pk
-- 			primary key,
-- 	name text not null,
-- 	type smallint not null,
-- 	crated_at timestamptz not null,
-- 	updated_at timestamptz not null
-- );
--
-- comment on column customer_param_types.type is '1: string, 2: numeric';
--
-- create unique index customer_param_types_name_uindex
-- 	on customer_param_types (name);

create table customer_params
(
	id bigint
		constraint customer_params_pk
			primary key,
	name text not null,
	type text not null,
	created_at timestamptz not null,
	updated_at timestamptz not null
);
comment on table customer_params is '顧客に設定する属性';
comment on column customer_params.type is '1: text, 2: numeric';
create unique index customer_params_name_uindex
	on customer_params (name);


create table customers_customer_params
(
    customer_id       bigint                   not null
        constraint customers_customer_params_customers_id_fk
            references customers,
    customer_param_id bigint                   not null
        constraint customers_customer_params_customer_params_id_fk
            references customer_params,
    value             text                     not null,
    created_at        timestamp with time zone,
    upadted_at        timestamp with time zone not null
);

create table customer_types_customer_params
(
	customer_type_id int not null
		constraint customer_types_customer_params_customer_types_id_fk
			references customer_types,
	customer_param_id bigint not null
		constraint customer_types_customer_params_customer_params_id_fk
			references customer_params,
	created_at timestamptz not null,
	updated_at timestamptz not null
);


-- +migrate Down
drop table customer_types_customer_params;
drop table customers_customer_params;
drop table customer_params;
drop table customers;
drop table customer_types;
