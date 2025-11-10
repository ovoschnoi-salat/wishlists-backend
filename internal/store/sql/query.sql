-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1;

-- name: CreateUser :execrows
INSERT INTO users (id, username, name, photo_url)
VALUES ($1, $2, $3, $4);

-- name: GetWishlists :many
SELECT *
FROM wishlists
WHERE owner_id = $1;

-- name: CreateWishlist :one
INSERT INTO wishlists (owner_id, title, description, is_private)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateWishlist :execrows
UPDATE wishlists
SET updated_at  = now(),
    title       = $1,
    description = $2,
    is_private  = $3
WHERE id = $4
  and owner_id = $5;

-- name: DeleteWishlist :execrows
DELETE
FROM wishlists
where id = $1
  and owner_id = $2;

-- name: CreateFriendsRequest :execrows
INSERT INTO friends_requests (user_id_from, user_id_to)
VALUES ($1, $2);

-- name: GetIncomingFriendsRequests :many
SELECT *
FROM users
WHERE id IN (SELECT user_id_from FROM friends_requests WHERE user_id_to = $1);

-- name: GetOutcomingFriendsRequests :many
SELECT *
FROM users
WHERE id IN (SELECT user_id_to FROM friends_requests WHERE user_id_from = $1);

-- name: GetIncomingFriendsRequestsCount :one
SELECT count(*)
FROM friends_requests
WHERE user_id_to = $1;

-- name: CreateFriendsRelationship :execrows
INSERT INTO friends (user_id, friend_id)
VALUES ($1, $2),
       ($2, $1);

-- name: GetFriends :many
SELECT *
FROM users
WHERE id IN (SELECT friend_id FROM friends WHERE user_id = $1);

-- name: AcceptFriendsRequest :execrows
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

-- name: DenyFriendsRequest :execrows
DELETE
FROM friends_requests
WHERE user_id_to = $1
  AND user_id_from = $2;

-- name: CheckIfFriends :one
SELECT *
FROM friends
WHERE user_id = $1
  AND friend_id = $2;

-- name: GetWishlistItems :many
SELECT *
FROM wishlist_items
WHERE wishlist_id = $1
  and owner_id = $2;

-- name: GetWishlistItem :one
SELECT *
FROM wishlist_items
WHERE id = $1;

-- name: CreateWishlistItem :one
INSERT INTO wishlist_items (wishlist_id, owner_id, title, description, price, links, reservable)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateWishlistItem :execrows
UPDATE wishlist_items
SET updated_at  = now(),
    title       = $1,
    description = $2,
    price       = $3,
    links       = $4,
    reservable  = $5
WHERE id = $6
  AND owner_id = $7;

-- name: CheckUserHasAccessToPrivateWishlist :one
SELECT *
FROM wishlist_access_list
WHERE list_id = $1
  AND user_id = $2;

-- name: ReserveWishlistItem :execrows
UPDATE wishlist_items
SET updated_at  = now(),
    reserved_by = $2
WHERE id = $1
  AND reserved_by IS NULL;

-- name: UnreserveWishlistItem :execrows
UPDATE wishlist_items
SET updated_at  = now(),
    reserved_by = NULL
WHERE id = $1
  AND reserved_by = $2;

-- name: DeleteWishlistItem :execrows
DELETE
FROM wishlist_items
WHERE id = $1
  and owner_id = $2;

-- name: GetFriendWishlists :many
SELECT *
FROM wishlists
WHERE wishlists.owner_id = $1
  AND (is_private = false OR
       id IN (SELECT list_id FROM wishlist_access_list WHERE wishlist_access_list.owner_id = $1 AND user_id = $2));

-- name: GetWishlistByWishId :one
SELECT *
FROM wishlists
WHERE wishlists.id = (SELECT wishlist_id FROM wishlist_items WHERE wishlist_items.id = $1);

-- name: CheckIfUserHasAccessToWishlist :one
SELECT id
FROM wishlists
WHERE id = $1
  AND (is_private = false OR
       id IN (SELECT list_id FROM wishlist_access_list WHERE wishlist_access_list.list_id = $1 AND user_id = $2));

-- name: GetWishlistAccessList :many
SELECT *
FROM wishlist_access_list
WHERE list_id = $1 AND owner_id = $2;

-- name: InsertWishlistAccessItem :execrows
INSERT INTO wishlist_access_list (list_id, owner_id, user_id)
VALUES ($1, $2, $3);

-- name: DeleteWishlistAccessItem :execrows
DELETE from wishlist_access_list WHERE list_id = $1 AND user_id = $2;

-- name: DeleteWishlistAccessItems :exec
DELETE from wishlist_access_list WHERE list_id = $1;