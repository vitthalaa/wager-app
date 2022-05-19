create table if not exists wager (
    id bigserial not null constraint wager_pk primary key,
    total_wager_value integer not null default 0,
    odds integer not null default 0,
    selling_percentage smallint not null default 1,
    selling_price real not null default 0,
    current_selling_price real not null default 0,
    percentage_sold smallint default null,
    amount_sold real default null,
    created_at timestamp default now(),
    updated_at timestamp default null
);

create table if not exists purchases (
    id bigserial not null constraint purchases_pk primary key,
    wager_id bigint not null,
    buying_price real not null default 0,
    created_at timestamp default now(),
    updated_at timestamp default null
);