-- DROP FUNCTION public.sync_profile_func();

CREATE OR REPLACE FUNCTION public.sync_profile_func()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    INSERT INTO profile_transfer(profile_id, full_name, phone)
    VALUES (NEW.profile_id, NEW.full_name, NEW.phone)
    ON CONFLICT (profile_id) DO UPDATE 
    SET full_name = EXCLUDED.full_name, phone= EXCLUDED.phone;
    RETURN NEW;
END;
$function$
;

SELECT trigger_name, event_manipulation, event_object_table
FROM information_schema.triggers
WHERE event_object_table = 'profile_transfer';

SELECT * FROM pg_proc WHERE proname = 'sync_profile_func';


SELECT * FROM pg_class WHERE relname = 'profile_transfer';
ALTER TABLE user_profile ENABLE TRIGGER profile_sync_trigger;
