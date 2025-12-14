package store

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

const recreateAccessListRequest = `-- name: RecreateWishlistAccessList :exec
WITH deleted_access AS (
    DELETE FROM wishlist_access_list
        WHERE list_id = $1 AND user_id NOT IN (%s))
INSERT
INTO wishlist_access_list (list_id, owner_id, user_id)
VALUES 
%s
ON CONFLICT DO NOTHING;`

type RecreateAccessListParams struct {
	WishlistId    int64
	OwnerID       int64
	NewFriendsIDs []int64
}

func (q *Queries) RecreateWishlistAccessList(ctx context.Context, arg RecreateAccessListParams) error {
	ids := strings.Join(lo.Map(arg.NewFriendsIDs, func(item int64, _ int) string {
		return strconv.FormatInt(item, 10)
	}), ", ")

	querySB := strings.Builder{}
	for i, id := range arg.NewFriendsIDs {
		if i != 0 {
			querySB.WriteString(", ")
		}
		querySB.WriteString(fmt.Sprintf("(%d, %d, %d)", arg.WishlistId, arg.OwnerID, id))
	}
	query := fmt.Sprintf(recreateAccessListRequest, ids, querySB.String())

	_, err := q.db.Exec(ctx, query, arg.WishlistId)
	if err != nil {
		return err
	}
	return nil
}
