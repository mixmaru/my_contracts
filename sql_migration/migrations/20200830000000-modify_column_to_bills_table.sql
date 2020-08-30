-- +migrate Up
alter table bills alter column billing_date type timestamptz using billing_date::timestamptz;
alter table bills alter column payment_confirmed_at type timestamptz using payment_confirmed_at::timestamptz;

-- +migrate Down
alter table bills alter column billing_date type date using billing_date::date;
alter table bills alter column payment_confirmed_at type timestamp using payment_confirmed_at::timestamp;
