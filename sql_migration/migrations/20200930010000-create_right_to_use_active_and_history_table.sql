-- +migrate Up
create table right_to_use_active
(
	right_to_use_id int
		constraint right_to_use_active_pk
			primary key
		constraint right_to_use_active_right_to_use_id_fk
			references right_to_use (id)
);
INSERT INTO right_to_use_active (right_to_use_id) SELECT id FROM right_to_use;

create table right_to_use_history
(
	right_to_use_id int
		constraint right_to_use_history_pk
			primary key
		constraint right_to_use_history_right_to_use_id_fk
			references right_to_use (id)
);

-- +migrate Down
drop table right_to_use_active;
drop table right_to_use_history;
