-- name: CreateUserSession :exec
insert into user_sessions (id, user_id)
values (?, ?);
