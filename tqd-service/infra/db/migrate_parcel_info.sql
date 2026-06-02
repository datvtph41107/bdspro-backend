
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