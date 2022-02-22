create table wallet
(
    id      serial primary key,
    user_id integer not null,
    balance integer default 0
);
create table transactions
(
    uuid             uuid      default gen_random_uuid() not null unique primary key,
    wallet_id        integer                             not null,
    amount           integer                             not null,
    transaction_type varchar                             not null,
    created          timestamp default now()
);
