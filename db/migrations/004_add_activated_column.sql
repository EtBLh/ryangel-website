ALTER TABLE client ADD COLUMN activated BOOLEAN DEFAULT FALSE;

-- Mark existing users as activated
UPDATE client SET activated = TRUE;
