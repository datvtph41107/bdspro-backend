# BDSPro - Tổng hợp toàn bộ bảng (Database Schema Reference)

Tài liệu tham khảo nhanh khi viết SQL thuần. Được sinh từ `bdspro-service/internal/domain`.

## Base fields (dùng chung)

Các entity dùng `_models.BaseEntity` có:

- `id` (bigint, PK)
- `created_at`, `updated_at` (timestamp)
- `deleted_at` (timestamp, soft delete)
- `created_by`, `updated_by` (bigint, audit)

---

## 1. PROPERTY (BĐS định danh)

### property_identify

Định danh BĐS gốc (PID + version).


| Column                                                     | Type   | Note            |
| ---------------------------------------------------------- | ------ | --------------- |
| id                                                         | bigint | PK              |
| pid                                                        | int8   | NOT NULL, index |
| version                                                    | int4   | NOT NULL        |
| created_at, updated_at, deleted_at, created_by, updated_by |        | BaseEntity      |


### property_lineage

Phiên bản vòng đời BĐS (Property chính để query).


| Column               | Type        | Note                          |
| -------------------- | ----------- | ----------------------------- |
| id                   | bigint      | PK                            |
| property_identify_id | int8        | index, FK → property_identify |
| location_id          | bigint      | index, FK → property_location |
| property_info_id     | bigint      | index, FK → property_info     |
| land_info_id         | bigint      | FK → property_land_info       |
| edvidence_id         | bigint      | FK → property_edvidence       |
| external_ref_id      | bigint      | FK → property_external_ref    |
| building_info_id     | bigint      | FK → property_building_info   |
| national_id          | varchar(50) | index                         |


### property_info

Thông tin cơ bản BĐS.


| Column               | Type         | Note                          |
| -------------------- | ------------ | ----------------------------- |
| id                   | bigint       | PK                            |
| property_identify_id | int8         | index, FK → property_identify |
| title                | varchar(255) |                               |
| source_type          | smallint     | EPropertySourceType           |
| avatar_id            | bigint       | FK → property_media           |
| property_type_id     | bigint       | index                         |
| project_id           | bigint       | index                         |
| legal_status         | smallint     | EHouseCertificate             |
| privacy_level        | smallint     | default 30                    |
| record_status        | smallint     | default 10                    |
| unit_code            | varchar(50)  |                               |
| identifier           | varchar(50)  |                               |
| level                | varchar(100) |                               |
| scope                | smallint     | default 10, index             |
| visibility           | smallint     | default 10                    |
| origin_profile_id    | bigint       | index                         |


### property_location

Địa chỉ/vị trí (không kế thừa BaseEntity).


| Column                                       | Type          | Note  |
| -------------------------------------------- | ------------- | ----- |
| id                                           | bigint        | PK    |
| property_identify_id                         | int8          | index |
| address_detail                               | varchar(100)  |       |
| region_id, province_id, district_id, ward_id | bigint        | index |
| latitude                                     | decimal(10,8) |       |
| longitude                                    | decimal(11,8) |       |
| map_url                                      | varchar(100)  |       |


### property_land_info

Thông tin đất/GCN.


| Column                            | Type          | Note       |
| --------------------------------- | ------------- | ---------- |
| id                                | bigint        | PK         |
| property_identify_id              | int8          | index      |
| document_type                     | smallint      |            |
| document_no                       | varchar(50)   |            |
| issuring_auth                     | varchar(50)   |            |
| plot, sheet                       | int4          |            |
| area_total, area_land, area_plant | decimal(15,2) |            |
| expired_land, expired_plant       | timestamp     |            |
| purpose_used                      | smallint      | default 10 |
| front_width, depth, street_width  | decimal(10,2) |            |
| land_note                         | text          |            |


### property_building_info

Thông tin nhà/công trình.


| Column                                     | Type          | Note       |
| ------------------------------------------ | ------------- | ---------- |
| id                                         | bigint        | PK         |
| property_identify_id                       | int8          | index      |
| area_actual, area_floor, area_construction | decimal(15,2) |            |
| floors, room_number, bedrooms, bathrooms   | smallint      |            |
| build_status                               | int8          | default 10 |
| building_type                              | int8          | default 10 |
| direction, balcony_direction               | int8          |            |


### property_media

Ảnh/video BĐS.


| Column               | Type         | Note      |
| -------------------- | ------------ | --------- |
| id                   | bigint       | PK        |
| property_identify_id | int8         | index     |
| lineage_id           | bigint       |           |
| media_type           | varchar(20)  | NOT NULL  |
| media_url            | varchar(500) | NOT NULL  |
| thumb_url            | varchar(500) |           |
| sort_order           | int          | default 0 |


### property_external_ref

Tham chiếu ngoài.


| Column               | Type         | Note            |
| -------------------- | ------------ | --------------- |
| id                   | bigint       | PK              |
| property_identify_id | int8         | index           |
| external_ref_id      | bigint       | NOT NULL, index |
| source_system        | smallint     | index           |
| source_code          | varchar(255) |                 |
| confidence           | smallint     | default 0       |
| sync_status          | smallint     | default 10      |
| note                 | text         |                 |


### property_edvidence

Chứng minh thực tế.


| Column               | Type         | Note  |
| -------------------- | ------------ | ----- |
| id                   | bigint       | PK    |
| property_identify_id | int8         | index |
| title                | varchar(255) |       |
| file_id              | bigint       | index |
| description          | text         |       |


### property_tag_link

Liên kết Property (lineage) với Tag.


| Column      | Type | Note                      |
| ----------- | ---- | ------------------------- |
| property_id | int8 | PK, FK → property_lineage |
| tag_id      | int8 | PK, FK → tag              |


### property_relation

Liên kết Property với Product/Asset.


| Column                 | Type | Note                        |
| ---------------------- | ---- | --------------------------- |
| property_id            | int8 | PK, FK → property_lineage   |
| relation_id            | int8 | index (Product/Asset ID)    |
| relation_type          | int8 | index, 10=Product, 20=Asset |
| updated_at, deleted_at |      |                             |


### property_user

User sở hữu Property.


| Column              | Type      | Note     |
| ------------------- | --------- | -------- |
| id                  | bigint    | PK       |
| owner_origin_id     | bigint    | NOT NULL |
| property_lineage_id | bigint    | NOT NULL |
| owner_at            | timestamp |          |
| role_id             | bigint    |          |
| origin_profile_id   | bigint    | index    |


### property_product

Liên kết Property với Product.


| Column      | Type   | Note                         |
| ----------- | ------ | ---------------------------- |
| id          | bigint | PK                           |
| property_id | int8   | index, FK → property_lineage |
| product_id  | int8   | index, FK → products         |


### tag

Tag tiện ích/đặc điểm.


| Column      | Type         | Note         |
| ----------- | ------------ | ------------ |
| id          | bigint       | PK           |
| name        | varchar(100) | NOT NULL     |
| type        | varchar(50)  |              |
| description | text         |              |
| icon        | varchar(255) |              |
| active      | bool         | default true |


### property_amenity

Bảng liên kết Property–Amenity (many2many).


| Column      | Type   | Note |
| ----------- | ------ | ---- |
| property_id | bigint | PK   |
| amenity_id  | bigint | PK   |


### lineage_media

Bảng liên kết property_lineage–property_media (many2many).


| Column              | Type   | Note |
| ------------------- | ------ | ---- |
| property_lineage_id | bigint |      |
| property_media_id   | bigint |      |


### property_spatial


| Column  | Type         | Note |
| ------- | ------------ | ---- |
| id      | bigint       | PK   |
| title   | varchar(255) |      |
| version | int2         |      |


### property_identifiers (CountryIdentifier)

Định danh quốc gia.


| Column                            | Type        | Note  |
| --------------------------------- | ----------- | ----- |
| id                                | varchar(40) | PK    |
| land_parcel_code                  | varchar(30) | index |
| project_code                      | varchar(30) | index |
| province_id, district_id, ward_id | varchar     | index |
| latitude, longitude               | decimal     |       |
| type, legal_status                | varchar     | index |
| current_owner_id                  | varchar(20) | index |


---

## 2. PRODUCT

### products


| Column                                | Type         | Note          |
| ------------------------------------- | ------------ | ------------- |
| id                                    | bigint       | PK            |
| parent_id, asset_id, property_id      | int8         |               |
| image_id                              | int8         |               |
| property_type_id                      | bigint       |               |
| doc_type_id                           | bigint       |               |
| doc_type_note                         | text         |               |
| project_id                            | int8         |               |
| visibility                            | smallint     | default 10    |
| name                                  | varchar(255) |               |
| code                                  | varchar(15)  |               |
| category_id                           | bigint       |               |
| area                                  | float        |               |
| position_url                          | varchar      |               |
| description                           | text         |               |
| last_price_id                         | bigint       |               |
| apartment_id                          | bigint       |               |
| source_type, source_status            |              |               |
| source_contact_id                     | int8         |               |
| priority                              | smallint     | default 20    |
| transaction_type                      | smallint     | default 10    |
| sale_status                           | smallint     | default 10    |
| sale_transaction_id                   | bigint       |               |
| sale_visibility                       | smallint     | default 10    |
| rent_status                           | smallint     | default 40    |
| rent_transaction_id                   | bigint       |               |
| rent_visibility                       | smallint     | default 10    |
| owner_id                              | int8         |               |
| owner_of                              | smallint     | default 10    |
| province_id, ward_id                  | bigint       |               |
| address                               | varchar(255) |               |
| google_map_link                       | varchar(255) |               |
| archived                              | bool         | default false |
| share_count, asset_count, build_count | int8         |               |
| post_count, deal_count, contact_count | int8         |               |
| appointment_count                     | int8         |               |
| mining_mode                           | smallint     | default 10    |
| mining_scope, mining_description      | text         |               |


### product_price


| Column                                   | Type          | Note               |
| ---------------------------------------- | ------------- | ------------------ |
| id                                       | bigint        | PK                 |
| product_id                               | bigint        |                    |
| currency                                 | varchar       |                    |
| sale_price, sale_commission              | decimal(15,0) |                    |
| sale_commission_type                     | int           |                    |
| deposite                                 | decimal(15,0) |                    |
| rent_price, rent_commission              | decimal(15,0) |                    |
| rent_commission_type, rent_payment_cycle | int           |                    |
| distribute_id                            | bigint        |                    |
| channel_price                            | bool          | default false      |
| price_status                             |               | ProductPriceStatus |
| change_note                              | text          |                    |
| changed_by                               | bigint        |                    |


### product_media


| Column             | Type         | Note          |
| ------------------ | ------------ | ------------- |
| id                 | bigint       | PK            |
| product_id         | bigint       | NOT NULL      |
| media_url          | varchar(500) | NOT NULL      |
| media_type         | varchar(20)  | NOT NULL      |
| is_main            | bool         | default false |
| sort_order / order | int          |               |


### product_user


| Column            | Type     | Note          |
| ----------------- | -------- | ------------- |
| id                | bigint   | PK            |
| profile_id        | bigint   | NOT NULL      |
| product_id        | bigint   | NOT NULL      |
| is_owner          | bool     | default false |
| role_id           | bigint   | default 0     |
| distribute_id     | bigint   | index         |
| owner_of          | smallint | default 10    |
| price_id          | bigint   |               |
| origin_profile_id | bigint   | index         |


### product_organization


| Column          | Type   | Note            |
| --------------- | ------ | --------------- |
| id              | bigint | PK              |
| organization_id | bigint | NOT NULL, index |
| product_id      | bigint | NOT NULL, index |
| is_owner        | bool   | default false   |
| role_id         | bigint | default 0       |


### product_asset


| Column     | Type   | Note            |
| ---------- | ------ | --------------- |
| id         | bigint | PK              |
| product_id | bigint | NOT NULL, index |
| asset_id   | bigint | NOT NULL, index |


### product_amenity


| Column     | Type   | Note |
| ---------- | ------ | ---- |
| product_id | bigint | PK   |
| amenity_id | bigint | PK   |


### product_stats


| Column                                | Type   | Note      |
| ------------------------------------- | ------ | --------- |
| product_id                            | bigint | PK        |
| total_views                           | bigint | default 0 |
| duration_views                        | int    | default 0 |
| last_view_event_time                  | bigint | index     |
| last_updated_at                       | bigint | index     |
| num_of_interested                     | bigint |           |
| post_count, deal_count                | bigint |           |
| share_count, asset_count, build_count | bigint |           |


### product_notes


| Column     | Type      | Note            |
| ---------- | --------- | --------------- |
| id         | bigint    | PK              |
| product_id | bigint    | NOT NULL, index |
| author_id  | bigint    | NOT NULL, index |
| content    | text      | NOT NULL        |
| pinned_at  | timestamp |                 |


### product_market (view/materialized)

Bảng/view market với nhiều cột join.

### house_info


| Column                                 | Type          | Note |
| -------------------------------------- | ------------- | ---- |
| product_id                             | bigint        | PK   |
| num_bedroom, num_bathroom, num_floor   | int           |      |
| num_front, num_car_park, num_toilet    | int           |      |
| furniture                              | varchar(50)   |      |
| orientation                            | varchar(50)   |      |
| orientation_house                      | int8          |      |
| certificate_house                      | int8          |      |
| road_width                             | decimal(10,2) |      |
| front_width, back_width, width, height | float         |      |


### distributions


| Column               | Type      | Note            |
| -------------------- | --------- | --------------- |
| id                   | bigint    | PK              |
| product_id           | bigint    | NOT NULL, index |
| can_deal             | bool      | default false   |
| distributed_at       | timestamp | NOT NULL        |
| revoked_at           | timestamp |                 |
| reason, reason_other |           |                 |
| from_date, to_date   | timestamp | NOT NULL        |
| price_id             | bigint    |                 |


---

## 3. ASSET

### assets


| Column               | Type      | Note          |
| -------------------- | --------- | ------------- |
| id                   | bigint    | PK            |
| name                 | varchar   | NOT NULL      |
| ward_id, province_id | bigint    |               |
| product_id           | bigint    |               |
| parent_asset_id      | bigint    |               |
| split_merge_status   | varchar   |               |
| purchase_price       | float     | NOT NULL      |
| purchase_date        | timestamp |               |
| legal_status         | smallint  | default 10    |
| archived             | bool      | default false |
| rent_status          | smallint  | default 10    |
| address, description |           |               |
| area                 | float     | NOT NULL      |
| property_type_id     | bigint    |               |
| image_id             | bigint    |               |
| owner_id             | bigint    |               |
| owner_of             | smallint  |               |


### asset_costs


| Column       | Type      | Note       |
| ------------ | --------- | ---------- |
| id           | bigint    | PK         |
| asset_id     | bigint    | NOT NULL   |
| owner_id     | bigint    |            |
| owner_type   | smallint  | default 10 |
| type         | smallint  | ECostType  |
| cost_type_id | bigint    |            |
| amount       | float     | NOT NULL   |
| date         | timestamp |            |
| description  | text      |            |


### asset_cost_types


| Column      | Type     | Note       |
| ----------- | -------- | ---------- |
| id          | bigint   | PK         |
| owner_id    | bigint   |            |
| owner_type  | smallint | default 10 |
| type_name   | varchar  | NOT NULL   |
| type        | smallint | default 10 |
| is_custom   | bool     |            |
| description | text     |            |


### asset_legals


| Column                                     | Type      | Note |
| ------------------------------------------ | --------- | ---- |
| id                                         | bigint    | PK   |
| asset_id                                   | bigint    |      |
| document_name, document_url, document_type |           |      |
| issued_date, expiry_date                   | timestamp |      |
| description                                | text      |      |
| related_split_merge_id                     | bigint    |      |


### asset_exploitation


| Column                     | Type      | Note       |
| -------------------------- | --------- | ---------- |
| id                         | bigint    | PK         |
| asset_id                   | bigint    | NOT NULL   |
| contract_name              | varchar   | NOT NULL   |
| type_id                    | bigint    |            |
| contract_value             | float     | NOT NULL   |
| customer                   | varchar   | NOT NULL   |
| cycle_pay                  | int       |            |
| income_amount              | float     |            |
| start_date, end_date       | timestamp |            |
| description, file_contract |           |            |
| status                     | int       | default 10 |
| owner_id                   | bigint    |            |
| owner_type                 | smallint  | default 10 |


### asset_user


| Column     | Type   | Note            |
| ---------- | ------ | --------------- |
| id         | bigint | PK              |
| profile_id | bigint | NOT NULL, index |
| asset_id   | bigint | NOT NULL, index |
| is_owner   | bool   | default false   |
| role_id    | bigint | default 0       |


### asset_organization


| Column          | Type   | Note            |
| --------------- | ------ | --------------- |
| id              | bigint | PK              |
| organization_id | bigint | NOT NULL, index |
| asset_id        | bigint | NOT NULL, index |
| is_owner        | bool   | default false   |
| role_id         | bigint | default 0       |


### asset_share


| Column | Type   | Note |
| ------ | ------ | ---- |
| id     | bigint | PK   |


---

## 4. POST

### posts


| Column                     | Type      | Note                       |
| -------------------------- | --------- | -------------------------- |
| id                         | bigint    | PK                         |
| product_id                 | bigint    | NOT NULL                   |
| expired_at                 | timestamp |                            |
| visibility                 | smallint  |                            |
| hidden                     | bool      | default false              |
| content                    | text      |                            |
| transaction_type           | smallint  |                            |
| title                      | varchar   |                            |
| num_date                   | int       |                            |
| price                      | float     | (post_price)               |
| price_type                 | int       | 1 chính thức, 2 thỏa thuận |
| package_visible            | int       |                            |
| like, comment, num_view    | bigint    |                            |
| owner_id                   | bigint    |                            |
| owner_of                   | smallint  | default 10                 |
| published_at, published_by |           |                            |
| code                       | varchar   |                            |
| status                     | smallint  | default 20                 |


### post_media


| Column     | Type         | Note          |
| ---------- | ------------ | ------------- |
| id         | bigint       | PK            |
| post_id    | bigint       | NOT NULL      |
| media_url  | varchar(500) | NOT NULL      |
| media_type | varchar(20)  | NOT NULL      |
| is_main    | bool         | default false |
| sort_order | int          | default 0     |


### post_organization, post_user


| Column | Type | Note                       |
| ------ | ---- | -------------------------- |
| ...    |      | Liên kết post với org/user |


---

## 5. LOCATION

### province_v2


| Column        | Type         | Note             |
| ------------- | ------------ | ---------------- |
| id            | bigint       | PK               |
| name          | varchar(255) | NOT NULL         |
| code          | int          | NOT NULL, unique |
| codename      | varchar(100) |                  |
| division_type | varchar(100) |                  |
| phone_code    | int          |                  |


### district_v2


| Column        | Type         | Note             |
| ------------- | ------------ | ---------------- |
| id            | bigint       | PK               |
| name          | varchar(255) | NOT NULL         |
| code          | int          | NOT NULL, unique |
| codename      | varchar(100) |                  |
| division_type | varchar(100) |                  |
| province_code | int          | NOT NULL, index  |
| province_id   | bigint       | NOT NULL, index  |


### ward_v2


| Column         | Type         | Note             |
| -------------- | ------------ | ---------------- |
| id             | bigint       | PK               |
| name           | varchar(255) | NOT NULL         |
| code           | int          | NOT NULL, unique |
| codename       | varchar(100) |                  |
| division_type  | varchar(100) |                  |
| short_codename | varchar(100) |                  |
| province_code  | int          | NOT NULL, index  |
| province_id    | bigint       | NOT NULL, index  |
| district_code  | int          | index            |
| district_id    | bigint       | index            |


### province, district, ward (v1 – có thể deprecated)


| Table    | Note |
| -------- | ---- |
| province |      |
| district |      |
| ward     |      |


### region


| Column | Type   | Note |
| ------ | ------ | ---- |
| id     | bigint | PK   |
| ...    |        |      |


---

## 6. PROJECT

### projects


| Column       | Type         | Note     |
| ------------ | ------------ | -------- |
| id           | bigint       | PK       |
| name         | varchar(255) | NOT NULL |
| developer_id | bigint       |          |
| description  | text         |          |


### project_build


| Column                   | Type         | Note     |
| ------------------------ | ------------ | -------- |
| id                       | bigint       | PK       |
| name                     | varchar(255) | NOT NULL |
| num_floor, num_apartment | int          |          |
| note                     | text         |          |
| project_id               | bigint       | NOT NULL |


### project_inventory


| Column | Type   | Note |
| ------ | ------ | ---- |
| id     | bigint | PK   |


### blocks


| Column     | Type         | Note            |
| ---------- | ------------ | --------------- |
| id         | bigint       | PK              |
| name       | varchar(255) | NOT NULL        |
| project_id | int8         | NOT NULL, index |


### building


| Column   | Type   | Note |
| -------- | ------ | ---- |
| id       | bigint | PK   |
| block_id | bigint |      |


### build_room


| Column                    | Type         | Note     |
| ------------------------- | ------------ | -------- |
| id                        | bigint       | PK       |
| block_id                  | bigint       | NOT NULL |
| name                      | varchar(255) | NOT NULL |
| floor, ord                | int          |          |
| unit_code                 | varchar(20)  |          |
| area                      | int          |          |
| num_bedroom, num_bathroom | int          |          |
| furniture                 | varchar      |          |
| blueprint_url             | varchar      |          |


### apartment


| Column       | Type         | Note     |
| ------------ | ------------ | -------- |
| id           | bigint       | PK       |
| name         | varchar(100) |          |
| note         | text         |          |
| floor        | int          |          |
| ordinal      | int          | NOT NULL |
| attribute_id | int8         |          |
| status       | smallint     |          |
| archived     | smallint     |          |
| build_id     | int8         |          |


### apartment_attribute


| Column | Type | Note |
| ------ | ---- | ---- |
| ...    |      |      |


---

## 7. DEAL

### deals


| Column                     | Type      | Note          |
| -------------------------- | --------- | ------------- |
| id                         | bigint    | PK            |
| owner_id                   | bigint    |               |
| owner_type                 | smallint  |               |
| name                       | varchar   |               |
| target_profit              | float     | default 0     |
| status                     | smallint  | default 10    |
| deal_type                  | smallint  | default 10    |
| note                       | text      |               |
| charge_person_id           | bigint    |               |
| cancel_reason              | text      |               |
| is_unilateral              | bool      | default false |
| is_manual                  | bool      | default false |
| bank_account_id            | bigint    |               |
| from_date, to_date         | timestamp |               |
| allow_sharing              | bool      |               |
| member_can_add_transaction | bool      |               |
| only_owner_get_commission  | bool      |               |
| internal_note              | text      |               |
| allow_manual_input         | bool      |               |
| description                | text      |               |
| share_visibility           | smallint  | default 10    |


### deal_members


| Column                     | Type      | Note          |
| -------------------------- | --------- | ------------- |
| id                         | bigint    | PK            |
| deal_id                    | bigint    | NOT NULL      |
| member_id                  | bigint    | NOT NULL      |
| member_type                | smallint  | default 10    |
| status                     | smallint  | default 10    |
| role_id                    | bigint    | NOT NULL      |
| amount_commit              | float     | default 0     |
| commission_value           | float     | default 0     |
| commission_type            | smallint  | default 10    |
| is_unilateral              | bool      | default false |
| done_investment            | bool      | default false |
| role_key                   | smallint  | default 430   |
| message                    | text      |               |
| invited_at                 | timestamp |               |
| responded_at, withdrawn_at | timestamp |               |
| inviter_id                 | bigint    |               |
| color_id                   | int       | index         |
| is_owner                   | bool      | default false |


### deal_products


| Column                      | Type   | Note     |
| --------------------------- | ------ | -------- |
| id                          | bigint | PK       |
| deal_id                     | bigint | NOT NULL |
| product_id                  | bigint | NOT NULL |
| unique(deal_id, product_id) |        |          |


### deal_of_organization, deal_of_group, deal_of_branch


| Table                | Note |
| -------------------- | ---- |
| deal_of_organization |      |
| deal_of_group        |      |
| deal_of_branch       |      |


### deal_milestone, deal_investment, deal_internal_note


| Table              | Note |
| ------------------ | ---- |
| deal_milestones    |      |
| deal_investment    |      |
| deal_internal_note |      |


---

## 8. TRANSACTION / TX

### tx_transaction


| Column           | Type      | Note        |
| ---------------- | --------- | ----------- |
| id               | bigint    | PK          |
| transaction_name | varchar   |             |
| from_id, to_id   | bigint    |             |
| from_of, to_of   | smallint  | TxOwnerType |
| method           | smallint  | TxMethod    |
| status           | smallint  | TxStatus    |
| timestamp        | timestamp |             |
| last_action_id   | bigint    |             |


### transactions

(TxTransaction dùng bảng `transactions` – có thể trùng với Transaction cũ)


| Column                         | Type     | Note |
| ------------------------------ | -------- | ---- |
| id                             | bigint   | PK   |
| owner_id, owner_type           |          |      |
| product_id                     | bigint   |      |
| deposite_amount                | float    |      |
| deposite_note                  | text     |      |
| contact_phone                  | varchar  |      |
| contact_id                     | bigint   |      |
| amount                         | float    |      |
| currency                       | varchar  |      |
| transaction_name               | varchar  |      |
| description                    | text     |      |
| transaction_type               | smallint |      |
| category_id, payment_method_id | int      |      |


### tx_action, tx_contract_deal, tx_deal_cost, tx_cost_type

### tx_transaction_history, tx_transaction_product


| Table                  | Note |
| ---------------------- | ---- |
| tx_action              |      |
| tx_contract_deal       |      |
| tx_deal_cost           |      |
| tx_cost_type           |      |
| tx_transaction_history |      |
| tx_transaction_product |      |


---

## 9. MASTER / LOOKUP

### property_type


| Column | Type         | Note         |
| ------ | ------------ | ------------ |
| id     | bigint       | PK           |
| name   | varchar(100) | NOT NULL     |
| active | bool         | default true |


### doc_type


| Column | Type         | Note         |
| ------ | ------------ | ------------ |
| id     | bigint       | PK           |
| name   | varchar(100) | NOT NULL     |
| active | bool         | default true |


### amenity


| Column | Type         | Note             |
| ------ | ------------ | ---------------- |
| id     | bigint       | PK               |
| name   | varchar(100) | NOT NULL, unique |
| active | bool         | default true     |


### developer


| Column                | Type         | Note |
| --------------------- | ------------ | ---- |
| id                    | bigint       | PK   |
| name                  | varchar(255) |      |
| slug                  | varchar(15)  |      |
| logo_url              | varchar      |      |
| description, address  |              |      |
| phone, email, website |              |      |
| founded_at            | timestamp    |      |
| status                | int          |      |


### color


| Column | Type   | Note |
| ------ | ------ | ---- |
| id     | bigint | PK   |


### bank_account


| Column | Type   | Note |
| ------ | ------ | ---- |
| id     | bigint | PK   |


### payment_method


| Column | Type   | Note |
| ------ | ------ | ---- |
| id     | bigint | PK   |


---

## 10. SHARING / ACCESS

### sharings


| Column               | Type     | Note |
| -------------------- | -------- | ---- |
| id                   | bigint   | PK   |
| target_id            | bigint   |      |
| target_type          | smallint |      |
| owner_id, owner_type |          |      |
| share_id             | bigint   |      |
| share_type           | varchar  |      |
| permissions          | varchar  |      |


### sharing_access


| Column         | Type     | Note       |
| -------------- | -------- | ---------- |
| id             | bigint   | PK         |
| domain_id      | bigint   |            |
| domain         | smallint | default 10 |
| from_type      | smallint | default 10 |
| from_id, to_id | bigint   |            |
| to_type        | smallint |            |
| fields         | text     |            |
| commission     | float    |            |
| permissions    | varchar  |            |


---

## 11. KHÁC

### intent


| Column | Type   | Note |
| ------ | ------ | ---- |
| id     | bigint | PK   |


### action


| Column | Type   | Note |
| ------ | ------ | ---- |
| id     | bigint | PK   |


### task


| Column | Type   | Note |
| ------ | ------ | ---- |
| id     | bigint | PK   |


### customer


| Column | Type   | Note |
| ------ | ------ | ---- |
| id     | bigint | PK   |


### record_history


| Column | Type   | Note |
| ------ | ------ | ---- |
| id     | bigint | PK   |


### attach_document


| Column   | Type   | Note                   |
| -------- | ------ | ---------------------- |
| id       | bigint | PK                     |
| owner_id | bigint | (có thể dùng cho deal) |


### bds_domain


| Column | Type   | Note |
| ------ | ------ | ---- |
| id     | bigint | PK   |


### profile_transfer, profile_info

### plan, member_plan

### sync_version, sync_version_history (audit)

### commission_stats

---

## Sơ đồ quan hệ nhanh (Property)

```
property_identify (1) ──┬── property_info (1-1)
                        ├── property_location (1-1)
                        ├── property_land_info (1-1)
                        ├── property_building_info (0-1)
                        ├── property_media (1-n)
                        ├── property_external_ref (1-n)
                        └── property_edvidence (0-1)

property_lineage (n) ──┬── property_identify (n-1)
                       ├── property_location (n-1)
                       ├── property_info (n-1)
                       ├── property_land_info (n-1)
                       ├── property_building_info (n-1)
                       └── property_relation (1-n) → products / assets

property_relation: property_id → property_lineage.id, relation_id → product/asset, relation_type (10=Product, 20=Asset)
property_tag_link: property_id → property_lineage.id, tag_id → tag.id
property_product: property_id → property_lineage.id, product_id → products.id
property_user: property_lineage_id → property_lineage.id
```

---

## Sơ đồ quan hệ nhanh (Product)

```
products ──┬── product_price (1-n, last_price_id)
           ├── product_media (1-n)
           ├── product_user (n-n)
           ├── product_organization (n-n)
           ├── product_asset (n-n) → assets
           ├── product_amenity (n-n) → amenity
           ├── house_info (1-1)
           ├── property_id → property_lineage
           ├── project_id → projects
           └── apartment_id → apartment

posts ── product_id → products
deal_products ── product_id → products, deal_id → deals
property_product ── product_id → products, property_id → property_lineage
```

---

*Tài liệu sinh từ domain models. Kiểm tra schema thực tế trong migrate/DB khi có sai lệch.*