-- name: GetUserSessionByMagicLink :one
update user_sessions
set magic_link_token = null
where magic_link_token = ?
    -- safety check to ensure the given token is not null
    and magic_link_token is not null
returning id, user_id;
