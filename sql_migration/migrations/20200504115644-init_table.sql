
-- +migrate Up
create table users
(
    id serial not null
        constraint users_pk
            primary key
);

create table "users_individual"
(
    user_id int not null
        constraint users_individual_pk
            primary key
        constraint users_individual_users_id_fk
            references users (id),
    name text not null
);

create table users_corporation
(
    user_id int not null
        constraint users_corporation_pk
            primary key
        constraint users_corporation_users_id_fk
            references users (id),
    name text not null
);

create table products
(
    id serial not null
        constraint products_pk
            primary key,
    name text not null,
    price decimal not null
);
create unique index products_name_uindex
    on products (name);

create table contracts
(
    id serial not null
        constraint contracts_pk
            primary key,
    user_id int not null
        constraint contracts_users_id_fk
            references users,
    product_id int not null
        constraint contracts_products_id_fk
            references products
);

create table contract_updates
(
    id bigserial not null
        constraint contract_updates_pk
            primary key,
    contract_id int not null
        constraint contract_updates_contracts_id_fk
            references contracts,
    update_num int not null,
    valid_from timestamp not null,
    valid_to timestamp not null
);
create unique index contract_updates_contract_id_update_num_uindex
    on contract_updates (contract_id, update_num);

create table bills
(
    id bigserial not null
        constraint bills_pk
            primary key,
    billing_date date not null,
    payment_date date
);

create table bill_details
(
    id bigserial not null
        constraint bill_details_pk
            primary key,
    bill_id bigint not null
        constraint bill_details_bills_id_fk
            references bills,
    order_num smallint not null,
    contract_update_id bigint not null
        constraint bill_details_contract_updates_id_fk
            references contract_updates
);
create unique index bill_details_bill_id_order_num_uindex
    on bill_details (bill_id, order_num);

create table discounts
(
    id serial not null
        constraint discounts_pk
            primary key,
    discount_type smallint not null,
    amount decimal not null,
    condition text not null,
    condition_values json not null,
    activate_from timestamp not null,
    activate_to timestamp not null
);

create table discount_individual
(
    discount_id int not null
        constraint discount_individual_pk
            primary key
        constraint discount_individual_discounts_id_fk
            references discounts,
    user_id int not null
        constraint discount_individual_users_id_fk
            references users
);

create table discount_products
(
    discount_id int not null
        constraint discount_products_pk
            primary key
        constraint discount_products_discounts_id_fk
            references discounts,
    product_id int not null
        constraint discount_products_products_id_fk
            references products
);

create table discount_contracts
(
    discount_id int not null
        constraint discount_contracts_pk
            primary key
        constraint discount_contracts_discounts_id_fk
            references discounts,
    contract_id int not null
        constraint discount_contracts_contracts_id_fk
            references contracts
);

create table discount_applies
(
    id bigserial not null
        constraint discount_applies_pk
            primary key,
    discount_id int not null
        constraint discount_applies_discounts_id_fk
            references discounts
);

create table discount_apply_bills
(
    discount_apply_id int not null
        constraint discount_apply_bills_pk
            primary key
        constraint discount_apply_bills_discount_applies_id_fk
            references discount_applies,
    bill_id bigint not null
        constraint discount_apply_bills_bills_id_fk
            references bills
);

create table discount_apply_contract_updates
(
    discount_apply_id int not null
        constraint discount_apply_contract_updates_pk
            primary key
        constraint discount_apply_contract_updates_discount_applies_id_fk
            references discount_applies,
    contract_update_id bigint not null
        constraint discount_apply_contract_updates_contract_updates_id_fk
            references contract_updates
);

-- +migrate Down
DROP TABLE discount_apply_contract_updates;
DROP TABLE discount_apply_bills;
DROP TABLE discount_applies;
DROP TABLE discount_contracts;
DROP TABLE discount_products;
DROP TABLE discount_individual;
DROP TABLE discounts;
DROP TABLE bill_details;
DROP TABLE bills;
DROP TABLE contract_updates;
DROP TABLE contracts;
DROP TABLE products;
DROP TABLE users_corporation;
DROP TABLE users_individual;
DROP TABLE users;
