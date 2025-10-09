DO $$ BEGIN
    CREATE TYPE provider_enum AS ENUM ('GOOGLE', 'ZALO', 'FACEBOOK', 'PHONE', 'ADMIN');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;