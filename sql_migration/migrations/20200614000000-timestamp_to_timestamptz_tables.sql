
-- +migrate Up
alter table bill_details alter column created_at type timestamptz using created_at::timestamptz;
alter table bill_details alter column updated_at type timestamptz using updated_at::timestamptz;

alter table bills alter column created_at type timestamptz using created_at::timestamptz;
alter table bills alter column updated_at type timestamptz using updated_at::timestamptz;

alter table contract_updates alter column created_at type timestamptz using created_at::timestamptz;
alter table contract_updates alter column updated_at type timestamptz using updated_at::timestamptz;

alter table contracts alter column created_at type timestamptz using created_at::timestamptz;
alter table contracts alter column updated_at type timestamptz using updated_at::timestamptz;

alter table discount_applies alter column created_at type timestamptz using created_at::timestamptz;
alter table discount_applies alter column updated_at type timestamptz using updated_at::timestamptz;

alter table discount_apply_bills alter column created_at type timestamptz using created_at::timestamptz;
alter table discount_apply_bills alter column updated_at type timestamptz using updated_at::timestamptz;

alter table discount_apply_contract_updates alter column created_at type timestamptz using created_at::timestamptz;
alter table discount_apply_contract_updates alter column updated_at type timestamptz using updated_at::timestamptz;

alter table discount_contracts alter column created_at type timestamptz using created_at::timestamptz;
alter table discount_contracts alter column updated_at type timestamptz using updated_at::timestamptz;

alter table discount_individual alter column created_at type timestamptz using created_at::timestamptz;
alter table discount_individual alter column updated_at type timestamptz using updated_at::timestamptz;

alter table discount_products alter column created_at type timestamptz using created_at::timestamptz;
alter table discount_products alter column updated_at type timestamptz using updated_at::timestamptz;

alter table discounts alter column created_at type timestamptz using created_at::timestamptz;
alter table discounts alter column updated_at type timestamptz using updated_at::timestamptz;

alter table users alter column created_at type timestamptz using created_at::timestamptz;
alter table users alter column updated_at type timestamptz using updated_at::timestamptz;

alter table users_corporation alter column created_at type timestamptz using created_at::timestamptz;
alter table users_corporation alter column updated_at type timestamptz using updated_at::timestamptz;

alter table users_individual alter column created_at type timestamptz using created_at::timestamptz;
alter table users_individual alter column updated_at type timestamptz using updated_at::timestamptz;

-- +migrate Down
alter table bill_details alter column created_at type timestamp using created_at::timestamp;
alter table bill_details alter column updated_at type timestamp using updated_at::timestamp;

alter table bills alter column created_at type timestamp using created_at::timestamp;
alter table bills alter column updated_at type timestamp using updated_at::timestamp;

alter table contract_updates alter column created_at type timestamp using created_at::timestamp;
alter table contract_updates alter column updated_at type timestamp using updated_at::timestamp;

alter table contracts alter column created_at type timestamp using created_at::timestamp;
alter table contracts alter column updated_at type timestamp using updated_at::timestamp;

alter table discount_applies alter column created_at type timestamp using created_at::timestamp;
alter table discount_applies alter column updated_at type timestamp using updated_at::timestamp;

alter table discount_apply_bills alter column created_at type timestamp using created_at::timestamp;
alter table discount_apply_bills alter column updated_at type timestamp using updated_at::timestamp;

alter table discount_apply_contract_updates alter column created_at type timestamp using created_at::timestamp;
alter table discount_apply_contract_updates alter column updated_at type timestamp using updated_at::timestamp;

alter table discount_contracts alter column created_at type timestamp using created_at::timestamp;
alter table discount_contracts alter column updated_at type timestamp using updated_at::timestamp;

alter table discount_individual alter column created_at type timestamp using created_at::timestamp;
alter table discount_individual alter column updated_at type timestamp using updated_at::timestamp;

alter table discount_products alter column created_at type timestamp using created_at::timestamp;
alter table discount_products alter column updated_at type timestamp using updated_at::timestamp;

alter table discounts alter column created_at type timestamp using created_at::timestamp;
alter table discounts alter column updated_at type timestamp using updated_at::timestamp;

alter table users alter column created_at type timestamp using created_at::timestamp;
alter table users alter column updated_at type timestamp using updated_at::timestamp;

alter table users_corporation alter column created_at type timestamp using created_at::timestamp;
alter table users_corporation alter column updated_at type timestamp using updated_at::timestamp;

alter table users_individual alter column created_at type timestamp using created_at::timestamp;
alter table users_individual alter column updated_at type timestamp using updated_at::timestamp;
