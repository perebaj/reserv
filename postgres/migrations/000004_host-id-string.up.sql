ALTER TABLE
    properties
ALTER COLUMN
    host_id TYPE VARCHAR USING host_id :: VARCHAR;

ALTER TABLE
    property_images
ALTER COLUMN
    host_id TYPE VARCHAR USING host_id :: VARCHAR;
