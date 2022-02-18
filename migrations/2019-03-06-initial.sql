-- +migrate Up
create table todos
(
    id    serial primary key,
    title varchar not null,
    done  boolean   default false,
    date  timestamp default now()
);

-- +migrate Down
drop table todos;
