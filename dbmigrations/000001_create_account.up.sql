CREATE TABLE accounts (
    id UUID PRIMARY KEY role TEXT NOT NULL CHECK (
        role IN ('client', 'shop', 'rider')
    ),
    email TEXT NOT NULL UNIQUE,
    phone TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    create_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

--Details only clients have

CREATE TABLE client_profiles (
    account_id UUID PRIMARY KEY REFERENCES account (id) ON DELETE CASCADE,
    full_name TEXT NOT NULL,
    state TEXT NOT NULL,
    address TEXT NOT NULL
);

CREATE TABLE shop_profiles (
    account_id UUID PRIMARY KEY REFERENCES account (id) ON DELETE CASCADE,
    business_name TEXT NOT NULL,
    state TEXT NOT NULL,
    address TEXT NOT NULL,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    is_verified BOOLEAN NOT NULL DEFAULT false,
    is_open BOOLEAN NOT NULL DEFAULT false,
    CONSTRAINT shop_lat_range CHECK (
        latitude IS NULL
        OR latitude BETWEEN -90 AND 90
    ),
    CONSTRAINT shop_lng_range CHECK (
        longitude IS NULL
        OR longitude BETWEEN -180 AND 180
    ),

--BOTH CO-ORDINATES ARE SET OR NEITHER IS
CONSTRAINT shop_open_needs_location CHECK (NOT is_open OR latitude IS NOT NULL)
);

--DETAILS ONLY RIDERS HAVE

CREATE TABLE rider_profiles (
    account_id UUID PRIMARY KEY REFERENCES account (id) ON DELETE CASCADE,
    full_name TEXT NOT NULL,
    registered_company TEXT NOT NULL,
    plate_number TEXT,
    has_ride BOOLEAN NOT NULL DEFAULT false,
    is_verified BOLEAN NOT NUL DEFAULT false,
    is_available BOOLEAN NOT NULL DEFAULT false
);