package service

import (
	"backend/internal/store"

	"github.com/gin-gonic/gin"
)

type Service struct {
	db *store.Queries
}

func NewService(db *store.Queries) *Service {
	s := &Service{
		db: db,
	}

	return s
}

func (s *Service) RegisterHandlers(g *gin.RouterGroup) {
	g = g.Group("/api")

	g.GET("/user/wishlists", s.GetUserWishlists)
	g.POST("/user/wishlist", s.CreateWishlist)
	g.PATCH("/user/wishlist", s.UpdateWishlist)

	g.GET("/user/wishlist/items", s.GetUserWishlistItems)
	g.POST("/user/wishlist/item", s.CreateUserWishlistItem)
	g.PATCH("/user/wishlist/item", s.UpdateUserWishlistItem)
	g.GET("/user/wishlist/item", s.GetUserWishlistItem)

	g.GET("/user/friends/requests/outcoming", s.GetUserOutcomingFriendsRequests)
	g.GET("/user/friends/requests/incoming", s.GetUserIncomingFriendsRequests)
	g.GET("/user/friends/requests/incoming/count", s.GetUserIncomingFriendsRequestsCount)
	g.POST("/user/friend/request", s.CreateUserFriendsRequest)
	g.POST("/user/friend/request/accept", s.AcceptUserIncomingFriendsRequest)
	g.POST("/user/friend/request/deny", s.DenyUserIncomingFriendsRequest)

	g.GET("/user/friends", s.GetFriends)
	g.GET("/user/friend/wishlists", s.GetUserFriendWishlists)
	g.GET("/user/friend/wishlist/items", s.GetUserFriendWishlistItems)
	g.POST("/user/friend/wishlist/wish/reservation/reserve", s.ReserveFriendWish)
	g.POST("/user/friend/wishlist/wish/reservation/cancel", s.CancelFriendWishReservation)
}

//
//func (s *Service) GetUserInfo(ctx context.Context, _ *emptypb.Empty) (*pb.GetUserInfoResponse, error) {
//	initData, err := middlewares.GetInitDataFromContext(ctx)
//	if err != nil {
//		return nil, err
//	}
//	marshal, err := json.Marshal(initData)
//	if err != nil {
//		return nil, err
//	}
//	return &pb.GetUserInfoResponse{
//		Data: string(marshal),
//	}, nil
//}

//func (s *Service) PatchUserSettings(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
//	//TODO implement me
//	panic("implement me PatchUserSettings")
//}
//
//func (s *Service) AddFriend(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
//	//TODO implement me
//	panic("implement me AddFriend")
//}
//
//func (s *Service) GetFriendsList(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Service) PatchWishlist(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Service) GetMyWishlist(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Service) GetUserWishlist(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Service) GetMyWish(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Service) GetUserWish(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s *Service) mustEmbedUnimplementedWishlistBackendServer() {
//	//TODO implement me
//	panic("implement me")
//}
