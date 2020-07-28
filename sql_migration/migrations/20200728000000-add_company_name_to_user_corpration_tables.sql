
-- +migrate Up
alter table users_corporation
	add company_name text not null;

-- +migrate Down
alter table users_corporation drop column company_name;
