
ALTER TABLE "crud" ADD COLUMN new_id SERIAL;
UPDATE "crud" SET new_id = id::integer;

ALTER TABLE "crud" DROP CONSTRAINT crud_pkey;

ALTER TABLE "crud" DROP COLUMN id;

ALTER TABLE "crud" CHANGE COLUMN new_id id SERIAL;

ALTER TABLE "crud" ADD PRIMARY KEY (id);