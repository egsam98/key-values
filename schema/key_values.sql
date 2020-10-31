create table if not exists key_values (
    k varchar(255) not null unique,
    v text
);

-- name: GetByKey :one
select v from key_values where k = $1;

-- name: PutKey :exec
insert into key_values (k, v) values ($1, $2);