-- name: GetUserSession :one
select
    user_sessions.id,
    user_id,
    users.name as username,
    users.email as user_email
from user_sessions
join users on users.id = user_sessions.user_id
where user_sessions.id = ?;
