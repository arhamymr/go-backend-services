-- Migration Down

-- Step 1: Remove the default value for the 'id' column
ALTER TABLE crud
ALTER COLUMN id DROP DEFAULT;

-- Step 2: Revert the column type to its original type (e.g., integer)
-- Warning: This will reset the IDs to sequential values; use a backup strategy if necessary
ALTER TABLE crud
ALTER COLUMN id TYPE integer USING nextval('crud_id_seq'::regclass);

-- Step 3: Optionally, restore the previous default value (if it existed)
-- ALTER TABLE users ALTER COLUMN id SET DEFAULT <original_default_value>;
