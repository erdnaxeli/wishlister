create table groups (
    id TEXT primary key,
    name TEXT not null
) strict;

create table wishlists (
    id TEXT primary key,
    admin_id TEXT not null,
    user_id TEXT not null references users (id),
    name TEXT not null,
    group_id TEXT references groups (id)
) strict;

create table wishlist_elements (
    id TEXT primary key,
    wishlist_id TEXT not null references wishlists (id),
    name TEXT not null,
    description text,
    url TEXT
) strict;

create table users (
    id TEXT primary key,
    name TEXT not null,
    email TEXT unique not null
) strict;
