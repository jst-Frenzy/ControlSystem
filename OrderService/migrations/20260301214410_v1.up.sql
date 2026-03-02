create table carts(
    id serial primary key,
    cart_id integer not null,
    name varchar(255) not null,
    product_id varchar(255) not null,
    quantity integer not null,
    price decimal(10, 2) not null,
    created_at timestamp default now()
)