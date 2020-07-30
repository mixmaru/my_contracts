
-- +migrate Up
alter table users_corporation rename column company_name to corporation_name;

-- +migrate Down
alter table users_corporation rename column corporation_name to company_name;
