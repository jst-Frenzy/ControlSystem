create table users(
    id serial primary key,
    user_name varchar(255) not null,
    email varchar(255) unique not null,
    password_hash varchar(255) not null,
    role varchar(20) default 'user',
    seller_id varchar(100) default null,
    created_at timestamp default now(),
    updated_at timestamp default now()
);

create table refresh_tokens(
    id serial not null,
    user_id integer not null references users(id) on delete cascade,
    token_hash varchar(255) not null,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    expires_at timestamp not null
);