UPDATE public.parcels
SET
    land_number = COALESCE(
        metadata ->> 'thua_dat_so_moi',
        land_number
    ),
    map_number = COALESCE(
        metadata ->> 'to_ban_do_so_moi',
        map_number
    ),
    metadata =
        metadata
        - 'thua_dat_so_moi'
        - 'to_ban_do_so_moi'
where 2000000 < id and id < 4000000 and (metadata ? 'thua_dat_so_moi'
   OR metadata ? 'to_ban_do_so_moi');



CREATE TABLE public.parcel_metadata (
    id BIGSERIAL PRIMARY KEY,
    key TEXT NOT NULL,
    value TEXT NOT NULL,

    created_at timestamptz DEFAULT now(),

    CONSTRAINT uq_parcel_metadata UNIQUE (key, value)
);

CREATE TABLE public.parcel_to_metadata (
    parcel_id BIGINT NOT NULL,
    metadata_id BIGINT NOT NULL,

    PRIMARY KEY (parcel_id, metadata_id),

    CONSTRAINT fk_ptm_parcel
        FOREIGN KEY (parcel_id)
        REFERENCES public.parcels(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_ptm_metadata
        FOREIGN KEY (metadata_id)
        REFERENCES public.parcel_metadata(id)
        ON DELETE CASCADE
);


CREATE INDEX idx_ptm_parcel_id
ON public.parcel_to_metadata(parcel_id);

CREATE INDEX idx_ptm_metadata_id
ON public.parcel_to_metadata(metadata_id);

CREATE INDEX idx_parcel_metadata_key
ON public.parcel_metadata(key);

CREATE INDEX idx_parcel_metadata_value
ON public.parcel_metadata(value);



-- 1. Thêm các giá trị DoAn chưa tồn tại
INSERT INTO parcel_metadata (key, value)
SELECT DISTINCT
    'QuyetDinh',
    p.metadata->>'QuyetDinh'
FROM parcels p
WHERE p.metadata ? 'QuyetDinh'
  AND p.metadata->>'QuyetDinh' <> ''
ON CONFLICT (key, value) DO NOTHING;


-- 2. Tạo quan hệ parcel -> metadata
INSERT INTO parcel_to_metadata (parcel_id, metadata_id)
SELECT
    p.id,
    pm.id
FROM parcels p
JOIN parcel_metadata pm
    ON pm.key = 'QuyetDinh'
   AND pm.value = p.metadata->>'QuyetDinh'
WHERE p.metadata ? 'QuyetDinh'
ON CONFLICT (parcel_id, metadata_id) DO NOTHING;

-- 3. Xóa DoAn khỏi metadata
UPDATE parcels
SET metadata = metadata - 'QuyetDinh'
WHERE metadata ? 'QuyetDinh' and 2000000 < id and id < 4000001;