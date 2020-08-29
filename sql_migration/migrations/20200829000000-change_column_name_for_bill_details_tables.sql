-- +migrate Up
alter table bill_details rename column contract_update_id to right_to_use_id;
alter table bill_details rename constraint bill_details_contract_updates_id_fk to bill_details_right_to_use_id_fk;

-- +migrate Down
alter table bill_details rename column right_to_use_id to contract_update_id ;
alter table bill_details rename constraint bill_details_right_to_use_id_fk to bill_details_contract_updates_id_fk;
