-- name: GetWishlists :many
SELECT *
FROM wishlists
WHERE owner_id = $1;

-- name: CreateWishlist :one
INSERT INTO wishlists (owner_id, title, description, is_private)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateWishlist :exec
UPDATE wishlists
SET title       = $1,
    description = $2,
    is_private  = $3
WHERE id = $4
  and owner_id = $5;

-- name: DeleteWishlist :exec
DELETE
FROM wishlists
where id = $1
  and owner_id = $2;

-- name: CreateFriendsRequest :exec
INSERT INTO friends_requests (user_id_from, user_id_to)
VALUES ($1, $2);

-- name: GetIncomingFriendsRequests :many
SELECT *
FROM users
WHERE id IN (SELECT user_id_from FROM friends_requests WHERE user_id_to = $1);

-- name: CreateFriendsRelationship :exec
INSERT INTO friends (user_id, friend_id)
VALUES ($1, $2),
       ($2, $1);

-- name: GetFriends :many
SELECT *
FROM users
WHERE id IN (SELECT friend_id FROM friends WHERE user_id = $1);

-- name: AcceptFriendsRequest :exec
WITH deleted_request AS (
    DELETE FROM friends_requests
        WHERE user_id_to = $1 AND user_id_from = $2
        RETURNING user_id_to, user_id_from)
INSERT
INTO friends (user_id, friend_id)
SELECT user_id_from, user_id_to
FROM deleted_request
UNION ALL
SELECT user_id_to, user_id_from
FROM deleted_request;

-- name: GetWishlistItems :many
SELECT *
FROM wishlist_items
WHERE wishlist_id = $1
  and owner_id = $2;

-- name: CreateWishlistItem :one
INSERT INTO wishlist_items (wishlist_id, owner_id, title, description, price, links, reservable)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateWishlistItem :exec
UPDATE wishlist_items
SET title       = $1,
    description = $2,
    price       = $3,
    links       = $4,
    reservable  = $5
WHERE id = $6
  AND owner_id = $7;


-- name: DeleteWishlistItem :exec
DELETE
FROM wishlist_items
WHERE id = $1
  and owner_id = $2;