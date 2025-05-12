CREATE TABLE amenities (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE property_amenities (
    property_id VARCHAR(255) NOT NULL,
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
