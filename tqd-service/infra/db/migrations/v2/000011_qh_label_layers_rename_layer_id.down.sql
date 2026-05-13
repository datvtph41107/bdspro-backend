DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'qh_label_layers' AND column_name = 'layer_id'
  ) AND NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'qh_label_layers' AND column_name = 'layer_id'
  ) THEN
    ALTER TABLE qh_label_layers RENAME COLUMN layer_id TO layer_id;
  END IF;
END $$;
