
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

-- +migrate Down
DROP TABLE products;

DROP TABLE users_corporation;
DROP TABLE users_individual;
DROP TABLE users;
