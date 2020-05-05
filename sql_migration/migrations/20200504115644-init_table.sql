
-- +migrate Up
create table users
(
    id bigserial not null
        constraint users_pk
            primary key
);

create table "users_individual"
(
    user_id bigserial not null
        constraint users_individual_pk
            primary key
        constraint users_individual_users_id_fk
            references users (id),
    name text not null
);

create table users_corporation
(
    user_id bigserial not null
        constraint users_corporation_pk
            primary key
        constraint users_corporation_users_id_fk
            references users (id),
    name text not null
);

create table products
(
    id bigserial not null
        constraint products_pk
            primary key,
    name text not null,
    price int not null
);

create unique index products_name_uindex
    on products (name);


create table contracts
(
    id bigserial not null
        constraint contracts_pk
            primary key,
    user_id bigint not null
        constraint contracts_users_id_fk
            references users,
    product_id bigint not null
        constraint contracts_products_id_fk
            references products
);

create table contract_updates
(
    contract_id bigint not null
        constraint contract_updates_contracts_id_fk
            references contracts,
    update_num int not null,
    start timestamp not null,
    "end" timestamp not null,
    constraint contract_updates_pk
        unique (contract_id, update_num)
);


-- +migrate Down
DROP TABLE contract_updates;
DROP TABLE contracts;

DROP TABLE products;

DROP TABLE users_corporation;
DROP TABLE users_individual;
DROP TABLE users;
