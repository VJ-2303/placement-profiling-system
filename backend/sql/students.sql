CREATE TABLE IF NOT EXISTS students (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    official_email text UNIQUE NOT NULL,
    profile_image_url text,
    version integer NOT NULL DEFAULT 1
);
