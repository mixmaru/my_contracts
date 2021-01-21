-- +migrate Up
alter table contracts
	add contract_date timestamptz not null;

comment on column contracts.contract_date is '契約日';

alter table contracts
	add billing_start_date timestamptz not null;

comment on column contracts.billing_start_date is '課金開始日';


-- +migrate Down
alter table contracts drop column contract_date;

alter table contracts drop column billing_start_date;
