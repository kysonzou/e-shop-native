package service

import (
	"context"
	v1 "github.com/kyson/e-shop-native/api/protobuf/user/v1"
	"github.com/kyson/e-shop-native/internal/user-srv/auth"
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
	// Token
	token, err := auth.GenerateToken(uint(user.ID), user.UserName)
	if err != nil{
		return nil, err
	}

	return &v1.LoginReply{
		Token: token,
		User: &v1.User{
			Id: int32(user.ID),
			Username: user.UserName,
			Phone: user.Phone,
			Email: user.Email,
		},
	}, nil		
}

func (s *UserService) GetMyProfile(ctx context.Context, req *v1.GetMyProfileRequest) (*v1.GetMyProfileReply, error){
	//这里少了一个验证Token，但是如果每个方法都自己验证的话，就是灾难性的，应该在之前就被验证
	//读取Token
	claims, ok := auth.FromContext(ctx)
	if !ok {
		return nil, auth.ErrTokenInvalid
	}

	user, err := s.uc.GetMyProfile(ctx, claims.Id)
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

