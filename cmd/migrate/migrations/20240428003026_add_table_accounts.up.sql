create table if not exists accounts (
    id serial primary key,
    user_id uuid not null references users(id),
    username text not null unique,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
)