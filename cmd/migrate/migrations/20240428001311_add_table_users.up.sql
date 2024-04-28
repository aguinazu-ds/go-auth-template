create table if not exists users (
    id uuid primary key default gen_random_uuid(),
    email text not null unique,
    activated boolean not null default false,
    encrypted_password text not null,
    confirmation_token text not null,
    confirmed_at timestamptz,
    recovery_token text not null,
    recovery_sent_at timestamptz,
    raw_app_meta_data jsonb not null default '{}',
    raw_user_meta_data jsonb not null default '{}',
    version integer not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
)