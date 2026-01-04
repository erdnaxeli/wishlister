-- name: DeleteUserSession :exec
delete from user_sessions
where id = ?;
