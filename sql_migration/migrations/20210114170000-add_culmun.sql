-- +migrate Up
-- カラム追加
alter table customers_customer_properties
    add customer_type_id int not null;

-- プライマリーキー変更
alter table customers_customer_properties drop constraint customers_customer_properties_pk;
alter table customers_customer_properties
    add constraint customers_customer_properties_pk
        primary key (customer_id, customer_type_id, customer_property_id);

-- 一旦外部キー削除
alter table customers_customer_properties drop constraint customers_customer_params_customers_id_fk;

-- 外部キーのためにuniqインデックスが必要だった
alter table customers drop constraint customers_pk;
alter table customers
    add constraint customers_pk
        primary key (id, customer_type_id);

-- 複合外部キー追加
alter table customers_customer_properties
    add constraint customers_customer_properties_customers_fk
        foreign key (customer_id, customer_type_id) references customers (id, customer_type_id);

alter table customers_customer_properties drop constraint customers_customer_params_customer_params_id_fk;
alter table customers_customer_properties
    add constraint customers_customer_properties_customer_types_customer_properties_customer_type_id_fk
        foreign key (customer_type_id, customer_property_id) references customer_types_customer_properties (customer_type_id, customer_property_id);

-- +migrate Down
