CREATE TABLE amenities (
    -- id is an string identifier for the amenity. Example: "wifi", "pool", "free_parking", "free_breakfast".
    -- TODO(@perebaj): Revisit this. Maybe we should use a UUID. As Im dumping the data manually, the string makes it easier to read.
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE property_amenities (
    property_id UUID NOT NULL,
    -- amenity_id is the id of the amenity.
    amenity_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (property_id, amenity_id)
);

-- TODO(@perebaj) As we are talking about a small app, we can hardcode the amenities.
-- Remove this once we have a real source of amenities.
INSERT INTO
    amenities (id, name)
VALUES
    ('wifi', 'WiFi'),
    ('free_parking', 'Free Parking'),
    ('free_breakfast', 'Free Breakfast'),
    ('pool', 'Pool'),
    ('hot_tub', 'Hot tub'),
    ('washer', 'Washer'),
    ('dryer', 'Dryer'),
    ('tv', 'TV'),
    ('coffee_maker', 'Coffee maker'),
    (
        'air_conditioning',
        'Air conditioning'
    ),
    ('heating', 'Heating'),
    (
        'dedicated_workspace',
        'Dedicated workspace'
    )
