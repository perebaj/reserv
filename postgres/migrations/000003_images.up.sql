CREATE TABLE property_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id UUID NOT NULL REFERENCES properties(id),
    host_id UUID NOT NULL,
    cloudflare_id UUID NOT NULL,
    filename TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);
