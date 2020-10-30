-- +goose Up
create table key_values (
    k varchar(255) not null unique,
    v text
);

-- +goose Down
drop table key_values;
