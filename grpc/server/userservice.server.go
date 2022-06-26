package userservice

import (
	"context"

	pb "github.com/sushilman/userservice/grpc/proto"
	"github.com/sushilman/userservice/models"
	"github.com/sushilman/userservice/services"
)

type UserServicegRPCServer struct {
	pb.UnimplementedUserServiceServer
	US services.IUserService
}

func (s *UserServicegRPCServer) CreateUser(ctx context.Context, userCreation *pb.UserCreation) (*pb.UserCreationResponse, error) {
	return nil, nil
}

func (s *UserServicegRPCServer) GetUsers(filterParams *pb.UserFilter, stream pb.UserService_GetUsersServer) error {
	q := mapFilterParamQueryParams(filterParams)
	users, err := s.US.GetUsers(q)
	if err != nil {
		return err
	}

	for _, user := range users {
		if err := stream.Send(mapUser(&user)); err != nil {
			return err
		}
	}

	return nil
}

func (s *UserServicegRPCServer) GetUserById(ctx context.Context, id *pb.UserId) (*pb.User, error) {
	user, err := s.US.GetUserById(id.Id)
	if err != nil {
		return nil, err
	}
	return mapUser(user), nil
}

func UpdateUser(ctx context.Context, id *pb.UserId) (*pb.UpdateUserResponse, error) { return nil, nil }

func DeleteUser(ctx context.Context, id *pb.UserId) (*pb.DeleteUserResponse, error) { return nil, nil }

func mapUser(u *models.User) *pb.User {
	return &pb.User{
		Id:        u.Id,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Nickname:  u.Nickname,
		Password:  u.Password,
		Email:     u.Email,
		Country:   u.Country,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func mapFilterParamQueryParams(f *pb.UserFilter) models.GetUserQueryParams {
	return models.GetUserQueryParams{
		FirstName: f.FirstName,
		LastName:  f.LastName,
		NickName:  f.Nickname,
		Email:     f.Email,
		Country:   f.Country,
		Offset:    uint(f.Offset),
		Limit:     uint(f.Limit),
	}
}
