-- name: GetWishlists :many
SELECT *
FROM wishlists
WHERE owner_id = $1;

-- name: CreateWishlist :one
INSERT INTO wishlists (owner_id, title, description, is_private)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CreateFriendsRequest :exec
INSERT INTO friends_requests (user_id_from, user_id_to)
VALUES ($1, $2);

-- name: GetIncomingFriendsRequests :many
SELECT *
FROM users
WHERE id IN (SELECT user_id_from FROM friends_requests WHERE user_id_to = $1);

-- name: CreateFriendsRelationship :exec
INSERT INTO friends (user_id, friend_id)
VALUES ($1, $2), ($2, $1);

-- name: GetFriends :many
SELECT *
FROM users
WHERE id IN (SELECT friend_id FROM friends WHERE user_id = $1);

