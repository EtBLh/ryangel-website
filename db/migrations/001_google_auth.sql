-- Add google_id to client table and make phone nullable
ALTER TABLE client ADD COLUMN google_id VARCHAR(255) UNIQUE;
ALTER TABLE client ALTER COLUMN phone DROP NOT NULL;
ALTER TABLE client ADD CONSTRAINT client_phone_or_google_check CHECK (phone IS NOT NULL OR google_id IS NOT NULL);
