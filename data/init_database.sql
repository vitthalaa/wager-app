create table if not exists wager (
    id bigserial not null constraint wager_pk primary key,
    total_wager_value integer not null default 0,
    odds integer not null default 0,
    selling_percentage real not null default 1,
    selling_price real not null default 0,
    current_selling_price real not null default 0,
    percentage_sold real default null,
    amount_sold real default null,
    created_at timestamp default now(),
    updated_at timestamp default null
);

create table if not exists purchases (
    id bigserial not null constraint purchases_pk primary key,
    wager_id bigint not null,
    buying_price real not null default 0,
    created_at timestamp default now(),
    updated_at timestamp default null,
    constraint purchases_wager_fk
        foreign key (wager_id)
            references wager (id)
            on update cascade on delete cascade
);