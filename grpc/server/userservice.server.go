package userservice

import (
	"context"
	"log"

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

// CreateUser returns the ID of the newly created user
func (s *UserServicegRPCServer) CreateUser(ctx context.Context, userCreation *pb.UserCreation) (*pb.UserCreationResponse, error) {
	uc := mapUserCreation(userCreation)
	id, err := s.US.CreateUser(ctx, uc)
	if err != nil {
		log.Printf("Error when creating the user. Error: %+v", err)
		return nil, status.Errorf(codes.Internal, "Something went wrong")
	}

	response := &pb.UserCreationResponse{Id: *id}

	return response, nil
}

// GetUsers returns the stream of the Users, filtered by parameters in the filterParams
func (s *UserServicegRPCServer) GetUsers(filterParams *pb.UserFilter, stream pb.UserService_GetUsersServer) error {
	q := mapFilterParamQueryParams(filterParams)
	users, err := s.US.GetUsers(context.Background(), q)
	if err != nil {
		log.Printf("Error while fetching users.\nError: %+v", err)
		return status.Errorf(codes.Internal, "Something went wrong")
	}

	for _, user := range users {
		if err := stream.Send(mapUser(&user)); err != nil {
			return err
		}
	}

	return nil
}

// GetUserById returns the user having the provided userId
func (s *UserServicegRPCServer) GetUserById(ctx context.Context, id *pb.UserId) (*pb.User, error) {
	user, err := s.US.GetUserById(ctx, id.Id)
	if err != nil {
		switch err.(type) {
		case *usererrors.NotFoundError:
			return nil, status.Errorf(codes.NotFound, "User Not found")
		}

		log.Printf("Error while fetching user by ID.\nError: %+v", err)
		return nil, status.Errorf(codes.Internal, "Something went wrong")
	}

	return mapUser(user), nil
}

// UpdateUser updates the user having the provided userId
func (s *UserServicegRPCServer) UpdateUser(ctx context.Context, userUpdate *pb.UserUpdate) (*pb.Empty, error) {
	uc := mapUserUpdate(userUpdate)
	err := s.US.UpdateUser(ctx, userUpdate.Id, uc)
	if err != nil {
		switch err.(type) {
		case *usererrors.NotFoundError:
			return nil, status.Errorf(codes.NotFound, "User Not found")
		}
		log.Printf("Error when updating the user. Error: %+v", err)
		return nil, status.Errorf(codes.Internal, "Something went wrong")
	}

	return &pb.Empty{}, nil
}

// DeleteUser deletes the user having the provided userId
func (s *UserServicegRPCServer) DeleteUser(ctx context.Context, id *pb.UserId) (*pb.Empty, error) {
	err := s.US.DeleteUserById(ctx, id.Id)
	if err != nil {
		log.Printf("Error when deleting the user. Error: %+v", err)
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
