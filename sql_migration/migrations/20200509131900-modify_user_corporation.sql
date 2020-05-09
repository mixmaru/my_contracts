
-- +migrate Up
alter table users_corporation rename column name to contact_person_name;

alter table users_corporation
	add president_name text not null;

-- +migrate Down
alter table users_corporation rename column contact_person_name to name;

alter table users_corporation drop column president_name;
