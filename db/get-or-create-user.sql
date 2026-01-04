-- name: GetOrCreateUser :one
insert into users (id, name, email)
values (?, ?, ?)
on conflict (email) do update set name = excluded.name
returning id, name, email;
