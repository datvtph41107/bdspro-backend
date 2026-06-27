-- public.qh_regions definition

-- Drop table

-- DROP TABLE public.qh_regions;

CREATE TABLE public.qh_regions (
	id bigserial NOT NULL,
	layer_id int8 NOT NULL,
	"name" varchar(255) NULL,
	display_name varchar(255) NULL,
	description text NULL,
	legal_doc text NULL,
	planning_name varchar(255) NULL,
	geometry public.geometry(multipolygon, 4326) NOT NULL,
	centroid public.geometry(point, 4326) NULL,
	area_sqm numeric DEFAULT 0 NULL,
	perimeter_m numeric DEFAULT 0 NULL,
	original_properties jsonb NULL,
	source_file varchar(255) NULL,
	import_batch_id varchar(100) NULL,
	style_config jsonb DEFAULT '{}'::jsonb NULL,
	status int8 DEFAULT 10 NULL,
	processing_status int8 DEFAULT 10 NULL,
	"version" int8 DEFAULT 1 NULL,
	is_latest bool DEFAULT true NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	label_id int8 NULL,
	land_use_id int8 NULL,
	legend_id int8 NULL,
	geom_hash text NULL,
	CONSTRAINT qh_regions_pkey PRIMARY KEY (id),
	CONSTRAINT fk_qh_regions_label FOREIGN KEY (label_id) REFERENCES public.qh_labels(id),
	CONSTRAINT fk_qh_regions_land_use FOREIGN KEY (land_use_id) REFERENCES public.qh_land_use(id),
	CONSTRAINT fk_qh_regions_legend FOREIGN KEY (legend_id) REFERENCES public.qh_legends(id)
);
CREATE INDEX idx_qh_regions_deleted_at ON public.qh_regions USING btree (deleted_at);
CREATE INDEX idx_qh_regions_import_batch_id ON public.qh_regions USING btree (import_batch_id);
CREATE INDEX idx_qh_regions_is_latest ON public.qh_regions USING btree (is_latest);
CREATE INDEX idx_qh_regions_label_id ON public.qh_regions USING btree (label_id);
CREATE INDEX idx_qh_regions_land_use_id ON public.qh_regions USING btree (land_use_id);
CREATE INDEX idx_qh_regions_layer_hash ON public.qh_regions USING btree (layer_id, geom_hash);
CREATE INDEX idx_qh_regions_layer_id ON public.qh_regions USING btree (layer_id);
CREATE INDEX idx_qh_regions_legend_id ON public.qh_regions USING btree (legend_id);
CREATE INDEX idx_qh_regions_processing_status ON public.qh_regions USING btree (processing_status);
CREATE INDEX idx_qh_regions_status ON public.qh_regions USING btree (status);
CREATE INDEX idx_regions_geom_status ON public.qh_regions USING gist (geometry) WHERE ((status = 10) AND (deleted_at IS NULL));

-- public.qh_layers definition

-- Drop table

-- DROP TABLE public.qh_layers;

CREATE TABLE public.qh_layers (
	id bigserial NOT NULL,
	display_name varchar(255) NOT NULL,
	description text NULL,
	"type" int4 NOT NULL,
	status int4 DEFAULT 10 NULL,
	display_order int4 DEFAULT 0 NULL,
	visible int4 DEFAULT 0 NULL,
	min_zoom int4 DEFAULT 6 NULL,
	max_zoom int4 DEFAULT 16 NULL,
	avatar text NULL,
	image_url text NULL,
	thumbnail_url text NULL,
	legal_doc text NULL,
	legal_doc_url text NULL,
	document_number varchar(100) NULL,
	issuing_authority varchar(255) NULL,
	effective_date date NULL,
	expiry_date date NULL,
	created_by int8 DEFAULT 0 NULL,
	updated_by int8 DEFAULT 0 NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	"name" varchar(100) NULL,
	source_type varchar(10) DEFAULT 'raster'::character varying NOT NULL,
	source_code varchar(15) DEFAULT 'qhsdd'::character varying NOT NULL,
	layer_url text NULL,
	style_config jsonb DEFAULT '{}'::jsonb NULL,
	import_status int4 DEFAULT 20 NULL,
	import_batch_id varchar(64) NULL,
	trust_value numeric DEFAULT 0 NULL,
	legal_status int8 DEFAULT 0 NULL,
	replaced_by_id int8 NULL,
	insuance_date date NULL,
	legal_level int8 DEFAULT 0 NULL,
	legal_number varchar(50) NULL,
	authority_issuring_id int8 NULL,
	scope_apply _int4 NULL,
	family_id int8 NULL,
	publish_scopes _int4 NULL,
	parent_id int8 NULL,
	analyse_order int8 DEFAULT 10 NULL,
	visible_default bool DEFAULT false NULL,
	CONSTRAINT qh_layers_pkey PRIMARY KEY (id),
	CONSTRAINT fk_qh_layer_families_layers FOREIGN KEY (family_id) REFERENCES public.qh_layer_families(id),
	CONSTRAINT fk_qh_layers_authority_issuring FOREIGN KEY (authority_issuring_id) REFERENCES public.qh_authority_issuring(id),
	CONSTRAINT fk_qh_layers_parent FOREIGN KEY (parent_id) REFERENCES public.qh_layers(id)
);
CREATE INDEX idx_layers_deleted ON public.qh_layers USING btree (id) WHERE (deleted_at IS NULL);
CREATE INDEX idx_qh_layers_analyse_order ON public.qh_layers USING btree (analyse_order);
CREATE INDEX idx_qh_layers_deleted_at ON public.qh_layers USING btree (deleted_at);
CREATE INDEX idx_qh_layers_display_order ON public.qh_layers USING btree (display_order);
CREATE INDEX idx_qh_layers_effective_date ON public.qh_layers USING btree (effective_date);
CREATE INDEX idx_qh_layers_expiry_date ON public.qh_layers USING btree (expiry_date);
CREATE INDEX idx_qh_layers_family_id ON public.qh_layers USING btree (family_id);
CREATE INDEX idx_qh_layers_import_batch_id ON public.qh_layers USING btree (import_batch_id);
CREATE INDEX idx_qh_layers_import_status ON public.qh_layers USING btree (import_status);
CREATE INDEX idx_qh_layers_insuance_date ON public.qh_layers USING btree (insuance_date);
CREATE INDEX idx_qh_layers_parent_id ON public.qh_layers USING btree (parent_id);
CREATE INDEX idx_qh_layers_replaced_by_id ON public.qh_layers USING btree (replaced_by_id);
CREATE INDEX idx_qh_layers_status ON public.qh_layers USING btree (status);
CREATE INDEX idx_qh_layers_type ON public.qh_layers USING btree (type);
CREATE INDEX idx_qh_layers_visible ON public.qh_layers USING btree (visible);

-- public.parcels definition

-- Drop table

-- DROP TABLE public.parcels;

CREATE TABLE public.parcels (
	created_by int8 NULL,
	updated_by int8 NULL,
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	geometry public.geometry(multipolygon, 4326) NULL,
	lat numeric NULL,
	lon numeric NULL,
	ref_id int8 NULL,
	ref_type int8 NULL,
	centroid_geom public.geometry(point, 4326) NULL,
	seo_id int8 NULL,
	is_seo bool DEFAULT false NOT NULL,
	CONSTRAINT parcels_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_geom_gist ON public.parcels USING gist (geometry);
CREATE INDEX idx_parcels_centroid_geom ON public.parcels USING gist (centroid_geom);
CREATE INDEX idx_parcels_deleted_at ON public.parcels USING btree (deleted_at);
CREATE INDEX idx_parcels_geom_deleted ON public.parcels USING gist (geometry) WHERE (deleted_at IS NULL);
CREATE INDEX idx_parcels_is_seo ON public.parcels USING btree (is_seo);

-- public.qh_land_use definition

-- Drop table

-- DROP TABLE public.qh_land_use;

CREATE TABLE public.qh_land_use (
	created_by int8 NULL,
	updated_by int8 NULL,
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	layer_id int8 NULL,
	land_use_group_id int8 NULL,
	note text NULL,
	"name" varchar(100) NULL,
	code varchar(10) NULL,
	display_order int8 DEFAULT 0 NULL,
	description text NULL,
	is_visible bool DEFAULT true NULL,
	priority int8 DEFAULT 50 NULL,
	build_condition text NULL,
	color varchar(10) NULL,
	can_build bool DEFAULT false NULL,
	is_active bool DEFAULT true NULL,
	warn_level int8 DEFAULT 0 NULL,
	CONSTRAINT qh_land_use_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX idx_layer_land_use_group ON public.qh_land_use USING btree (layer_id, land_use_group_id);
CREATE INDEX idx_qh_land_use_deleted_at ON public.qh_land_use USING btree (deleted_at);

