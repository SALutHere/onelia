CREATE SCHEMA onelia;

CREATE TABLE onelia.cities (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE onelia.segments (
    id UUID PRIMARY KEY,

    from_city_id UUID NOT NULL
        REFERENCES onelia.cities(id),

    to_city_id UUID NOT NULL
        REFERENCES onelia.cities(id),

    transport_type TEXT NOT NULL
        CHECK (transport_type IN ('bus', 'train', 'plane')),

    duration_minutes INTEGER NOT NULL
        CHECK (duration_minutes > 0),

    price NUMERIC(10, 2) NOT NULL
        CHECK (price >= 0),

    CHECK (from_city_id <> to_city_id)
);
