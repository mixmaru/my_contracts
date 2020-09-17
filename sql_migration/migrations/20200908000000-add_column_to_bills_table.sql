-- +migrate Up
alter table bills
	add user_id bigint;

alter table bills
	add constraint bills_users_id_fk
		foreign key (user_id) references users;

UPDATE bills SET user_id = c.user_id
FROM bill_details bd
INNER JOIN right_to_use rtu on bd.right_to_use_id = rtu.id
INNER JOIN contracts c on rtu.contract_id = c.id
WHERE bills.id = bd.bill_id
;

alter table bills alter column user_id set not null;

-- +migrate Down
alter table bills drop constraint bills_users_id_fk;

alter table bills drop column user_id;
