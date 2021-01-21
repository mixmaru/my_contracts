-- +migrate Up
drop index contract_updates_contract_id_update_num_uindex;

alter table contract_updates drop column update_num;

alter table contract_updates rename to right_to_use;

-- +migrate Down
alter table right_to_use rename to contract_updates;

alter table contract_updates
	add update_num int not null;

create unique index contract_updates_contract_id_update_num_uindex
	on contract_updates (contract_id, update_num);

