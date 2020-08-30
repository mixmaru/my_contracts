-- +migrate Up
alter table bill_details
	add billing_amount decimal not null;

-- +migrate Down
alter table bill_details drop column billing_amount;
