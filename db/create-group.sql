-- name: CreateGroup :exec
insert into groups (
    id, name
)
values (
    ?, ?
);
