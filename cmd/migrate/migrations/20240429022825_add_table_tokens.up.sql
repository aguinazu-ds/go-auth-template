create table if not exists tokens (
    hash bytea primary key,
    user_id uuid not null references users(id) on delete cascade,
    expires_at timestamptz not null,
    scope text not null,
    created_at timestamptz not null default now()
)