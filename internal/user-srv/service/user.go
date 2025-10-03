package service

import (
	//"context"

	"context"

	v1 "github.com/kyson/e-shop-native/api/protobuf/user/v1"
	"github.com/kyson/e-shop-native/internal/user-srv/biz"
)


type UserService struct {
	uc biz.UserService
	// Add any dependencies or fields here
	v1.UnimplementedUserServiceServer
}

func NewUserService(uc biz.UserService) v1.UserServiceServer {
	return &UserService{
		uc: uc,
	}
}

func (s *UserService) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterReply, error){
	user := &biz.User{
		UserName: req.Username,
		Password: req.Password,
		Phone: req.Phone,
		Email: req.Email,
	}
	_, err := s.uc.RegisterUser(ctx, user)
	if err != nil {
		return nil, err
	}	
	return &v1.RegisterReply{
		User: &v1.User{
			Id: int32(user.ID),
			Username: user.UserName,
			Phone: user.Phone,
			Email: user.Email,
		},
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error){
	user, err := s.uc.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		User: &v1.User{
			Id: int32(user.ID),
			Username: user.UserName,
			Phone: user.Phone,
			Email: user.Email,
		},
	}, nil		
}

func (s *UserService) GetMyProfile(ctx context.Context, req *v1.GetMyProfileRequest) (*v1.GetMyProfileReply, error){
	user, err := s.uc.GetMyProfile(ctx, uint(1))
	if err != nil {
		return nil, err
	}
	return &v1.GetMyProfileReply{
		User: &v1.User{
			Id: int32(user.ID),	
			Username: user.UserName,
			Phone: user.Phone,
			Email: user.Email,
		},
	}, nil
}

