ALTER TABLE notifications
    DROP COLUMN IF EXISTS title_key,
    DROP COLUMN IF EXISTS body_key,
    DROP COLUMN IF EXISTS params;
