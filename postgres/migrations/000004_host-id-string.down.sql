ALTER TABLE
    properties
ALTER COLUMN
    host_id TYPE UUID USING host_id :: UUID;

ALTER TABLE
    property_images
ALTER COLUMN
    host_id TYPE UUID USING host_id :: UUID;
