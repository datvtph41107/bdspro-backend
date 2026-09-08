CREATE INDEX idx_parcel_adr_search_trgm ON public.qh_parcel_info USING gin (lower(adr_search) gin_trgm_ops)

-- convert ko chữ ko dấu rồi gán vào adr_search
UPDATE public.qh_parcel_info
SET adr_search = trim(
    regexp_replace(
        lower(
            unaccent(
                coalesce(address_text, '')
            )
        ),
        '[^a-z0-9]+',
        ' ',
        'g'
    )
);


CREATE INDEX idx_search_text_trgm 
ON public.search_index 
USING GIN(search_text gin_trgm_ops);


INSERT INTO public.search_index (
    poi_id,
    title,
    address,
    search_text,
    geom
)
SELECT 
    id AS poi_id,
    name AS title,
    address,
    
    TRIM(
        REGEXP_REPLACE(
            LOWER(
                UNACCENT(
                    COALESCE(name, '') || ' ' || COALESCE(address, '')
                )
            ),
            '[^a-z0-9]+',
            ' ',
            'g'
        )
    ) AS search_text,
    
    geom
FROM public.pois
WHERE deleted_at IS NULL 
  AND is_active = true
ON CONFLICT (poi_id) DO UPDATE 
SET 
    title = EXCLUDED.title,
    address = EXCLUDED.address,
    search_text = EXCLUDED.search_text,
    geom = EXCLUDED.geom,
    updated_at = NOW();
    
    

ALTER TABLE public.search_index 
ADD CONSTRAINT unique_poi_id UNIQUE (poi_id);