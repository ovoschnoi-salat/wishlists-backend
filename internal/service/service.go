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
	g.GET("user/wishlists", s.GetMyWishlists)
	g.POST("user/wishlists", s.CreateWishlist)
	g.GET("user/friends", s.GetFriends)
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
