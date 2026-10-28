-- public.tb_contact definition (reference for modules/crm entity Contact).
-- Not executed by the service: schema changes go through v4/migrate (000046–000048 for tb_contact).

-- Drop table

-- DROP TABLE public.tb_contact;

CREATE TABLE public.tb_contact (
	created_by int8 NULL,
	updated_by int8 NULL,
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	profile_id int8 NULL,
	phone text NULL,
	full_name text NULL,
	avatar text NULL,
	email text NULL,
	zalo text NULL,
	company text NULL,
	note text NULL,
	birthday date NULL,
	owner_id int8 NULL,
	owner_of int2 DEFAULT 30 NULL,
	friend_id int8 NULL,
	origin_profile_id int8 NULL,
	address text NULL,
	visibility int2 DEFAULT 30 NULL,
	status int2 DEFAULT 10 NULL,
	lead_id int8 NULL,
	contact_tags _int4 NULL,
	priority_pin_at timestamptz NULL,
	last_interaction_at timestamptz NULL,
	CONSTRAINT tb_contact_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_tb_contact_last_interaction_at ON public.tb_contact USING btree (last_interaction_at) WHERE (deleted_at IS NULL);
CREATE INDEX idx_tb_contact_origin_profile_id ON public.tb_contact USING btree (origin_profile_id) WHERE (deleted_at IS NULL);
CREATE INDEX idx_tb_contact_owner ON public.tb_contact USING btree (owner_id, owner_of) WHERE (deleted_at IS NULL);
CREATE INDEX idx_tb_contact_priority_pin_at ON public.tb_contact USING btree (priority_pin_at) WHERE (deleted_at IS NULL);
