create table wishlists (
    id varchar(21) primary key not null,
    admin_id varchar(21) not null,
    user_id varchar(21) references users (id),
    name varchar(255) not null
);

create table wishlist_elements (
    id varchar(21) primary key not null,
    wishlist_id varchar(21) not null references wishlists (id),
    name varchar(255) not null,
    description text,
    url varchar(2048)
);

create table users (
    id varchar(21) primary key not null,
    name varchar(255) not null,
    email varchar(255)
);
