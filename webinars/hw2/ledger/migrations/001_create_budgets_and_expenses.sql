-- +goose Up

create table if not exists budgets (
    id serial primary key,
    category text not null,
    limit_amount numeric(14,2) not null check (limit_amount > 0),
    period text not null,
    unique (category, period)
);

create table if not exists expenses (
    id serial primary key,
    amount numeric(14,2) not null check (amount > 0),
    category text not null,
    description text,
    date date not null
);

create index if not exists idx_expenses_category_date
on expenses(category, date);

-- +goose Down

drop index if exists idx_expenses_category_date;

drop table if exists expenses;

drop table if exists budgets;