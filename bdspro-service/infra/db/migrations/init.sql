CREATE INDEX idx_rooms_build_floor_ordinal ON apartment(build_id, floor, ordinal);

-- Add sourceContactId and sourceContactNote columns to products table
ALTER TABLE products ADD COLUMN IF NOT EXISTS source_contact_id BIGINT;
ALTER TABLE products ADD COLUMN IF NOT EXISTS source_contact_note TEXT;

CREATE MATERIALIZED VIEW public.product_market
TABLESPACE pg_default
AS SELECT p.id,
    p.created_at,
    p.updated_at,
    p.name,
    p.code,
    p.category_id,
    p.area,
    p.description,
    p.transaction_type,
    p.sale_visibility,
    p.rent_visibility,
    p.owner_id AS owner_user_id,
    p.owner_org_id,
    p.property_type_id,
    p.doc_type_id,
    p.position_url,
    p.address,
    p.google_map_link,
    p.asset_id,
    p.apartment_id,
    p.last_price_id,
    p.transaction_price,
    p.deleted_at,
    p.archived,
    hi.num_bedroom,
    hi.num_bathroom,
    hi.num_floor,
    hi.furniture,
    hi.orientation,
    hi.num_front,
    hi.num_car_park,
    hi.num_toilet,
    string_agg(DISTINCT a.amenity_id::text, ','::text) AS amenity_ids,
    array_remove(array_agg(DISTINCT pm.media_url), NULL::character varying) AS media_urls,
    array_remove(array_agg(DISTINCT pm.media_type), NULL::character varying) AS media_types,
    max(
        CASE
            WHEN pm.is_main THEN pm.media_url
            ELSE NULL::character varying
        END::text) AS main_image,
    province.name AS province_name,
    district.name AS district_name,
    ward.name AS ward_name,
    dt.name AS doc_type_name,
    pt.name AS property_type_name,
    pp.currency,
    pp.price_owner,
    pp.sale_price,
    pp.sale_commission,
    pp.sale_commission_type,
    pp.rent_price,
    pp.rent_commission,
    pp.rent_commission_type,
    pp.rent_payment_cycle,
    pp.import_price,
    pp.operating_cost,
    pp.internal_note,
    pp.target_profit,
    pp.private_docs,
    pp.deposite
   FROM products p
     LEFT JOIN house_info hi ON hi.product_id = p.id AND hi.deleted_at IS NULL
     LEFT JOIN product_amenity a ON a.product_id = p.id
     LEFT JOIN product_media pm ON pm.product_id = p.id AND pm.deleted_at IS NULL
     LEFT JOIN region province ON province.id = p.province_id AND province.deleted_at IS NULL
     LEFT JOIN region district ON district.id = p.district_id AND district.deleted_at IS NULL
     LEFT JOIN region ward ON ward.id = p.ward_id AND ward.deleted_at IS NULL
     LEFT JOIN doc_type dt ON dt.id = p.doc_type_id AND dt.deleted_at IS NULL
     LEFT JOIN property_type pt ON pt.id = p.property_type_id AND pt.deleted_at IS NULL
     LEFT JOIN LATERAL ( SELECT pp2.created_by,
            pp2.updated_by,
            pp2.id,
            pp2.created_at,
            pp2.updated_at,
            pp2.deleted_at,
            pp2.product_id,
            pp2.currency,
            pp2.price_owner,
            pp2.sale_price,
            pp2.sale_commission,
            pp2.sale_commission_type,
            pp2.rent_price,
            pp2.rent_commission,
            pp2.rent_commission_type,
            pp2.rent_payment_cycle,
            pp2.import_price,
            pp2.operating_cost,
            pp2.internal_note,
            pp2.target_profit,
            pp2.private_docs,
            pp2.deposite
           FROM product_price pp2
          WHERE pp2.product_id = p.id AND pp2.deleted_at IS NULL
          ORDER BY pp2.created_at DESC NULLS LAST, pp2.id DESC
         LIMIT 1) pp ON true
  WHERE p.deleted_at IS NULL AND (p.transaction_type = 20 AND p.rent_visibility = 10 OR p.transaction_type = 10 AND p.sale_visibility = 10)
  GROUP BY p.id, pp.id, hi.product_id, province.id, 
  district.id, ward.id, pt.id, dt.id, pp.currency, 
  pp.price_owner, pp.sale_price, pp.sale_commission, 
  pp.sale_commission_type, pp.rent_price, pp.rent_commission, 
  pp.rent_commission_type, pp.rent_payment_cycle, pp.import_price, 
  pp.operating_cost, pp.internal_note, pp.target_profit, pp.private_docs, pp.deposite