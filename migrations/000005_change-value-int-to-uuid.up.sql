-- Migration Up

DO $$ 
BEGIN
    -- Create the extension if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp') THEN
        CREATE EXTENSION "uuid-ossp";
    END IF;
END $$;

-- Step 1: Remove the default value from the 'id' column
ALTER TABLE crud
ALTER COLUMN id DROP DEFAULT;

-- Step 2: Change the column type to UUID
ALTER TABLE crud
ALTER COLUMN id TYPE UUID USING uuid_generate_v4();

-- Step 3: Set the new default to automatically generate UUIDs for new records
ALTER TABLE crud
ALTER COLUMN id SET DEFAULT uuid_generate_v4();

-- Step 4: Update the existing records with UUIDs
UPDATE crud
SET id = uuid_generate_v4();
