package userservice

import (
	"context"

	pb "github.com/sushilman/userservice/grpc/proto"
	"github.com/sushilman/userservice/models"
	"github.com/sushilman/userservice/services"
	"github.com/sushilman/userservice/usererrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServicegRPCServer struct {
	pb.UnimplementedUserServiceServer
	US services.IUserService
}

func (s *UserServicegRPCServer) CreateUser(ctx context.Context, userCreation *pb.UserCreation) (*pb.UserCreationResponse, error) {
	uc := mapUserCreation(userCreation)
	id, err := s.US.CreateUser(uc)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Something went wrong")
	}

	response := &pb.UserCreationResponse{Id: id}

	return response, nil
}

func (s *UserServicegRPCServer) GetUsers(filterParams *pb.UserFilter, stream pb.UserService_GetUsersServer) error {
	q := mapFilterParamQueryParams(filterParams)
	users, err := s.US.GetUsers(q)
	if err != nil {
		return status.Errorf(codes.Internal, "Something went wrong")
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
		switch err.(type) {
		case *usererrors.NotFoundError:
			return nil, status.Errorf(codes.NotFound, "User Not found")
		}

		return nil, status.Errorf(codes.Internal, "Something went wrong")
	}

	return mapUser(user), nil
}

func (s *UserServicegRPCServer) UpdateUser(ctx context.Context, userUpdate *pb.UserUpdate) (*pb.Empty, error) {
	uc := mapUserUpdate(userUpdate)
	err := s.US.UpdateUser(userUpdate.Id, uc)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Something went wrong")
	}

	return &pb.Empty{}, nil
}

func (s *UserServicegRPCServer) DeleteUser(ctx context.Context, id *pb.UserId) (*pb.Empty, error) {
	err := s.US.DeleteUserById(id.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Something went wrong")
	}

	return &pb.Empty{}, nil
}

func mapUserCreation(u *pb.UserCreation) models.UserCreation {
	return models.UserCreation{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Nickname:  u.Nickname,
		Password:  u.Password,
		Email:     u.Email,
		Country:   u.Country,
	}
}

func mapUserUpdate(u *pb.UserUpdate) models.UserCreation {
	return models.UserCreation{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Nickname:  u.Nickname,
		Password:  u.Password,
		Email:     u.Email,
		Country:   u.Country,
	}
}

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
