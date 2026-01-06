-- name: GetUserSession :one
select
    id,
    user_id
from user_sessions
where id = ?;
