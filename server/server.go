package server

import (
	"context"
	"fmt"

	pb "github.com/casmelad/bootcamp-gateway/server/proto"
	domain "github.com/casmelad/bootcamp-gateway/users"
	mappers "github.com/casmelad/bootcamp-gateway/users/mappers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	appService domain.Service
	pb.UsersServer
}

func NewUserServer(s domain.Service) *UserServer {
	return &UserServer{
		appService: s,
	}
}

//Get a user by the email
func (s UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {

	fmt.Println(req)

	email := req.GetEmail()

	result, err := s.appService.GetByEmail(ctx, email)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user with Email %s could not be found", req.GetEmail())
	}

	mappedUser, err := mappers.ToGrpcUser(result)

	return &pb.GetUserResponse{User: &mappedUser}, nil
}

//Creates a nw user record
func (s UserServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {

	fmt.Println(req)

	err := req.Validate()

	if err != nil {
		return nil, err
	}

	user := domain.User{
		Email:    req.GetEmail(),
		Name:     req.GetName(),
		LastName: req.GetLastName(),
	}

	result, err := s.appService.Create(ctx, user)

	if err != nil {
		switch err {
		case domain.ErrInvalidData:
			return nil, status.Errorf(codes.InvalidArgument, "Invalid data")
		case domain.ErrUserAlreadyExists:
			return nil, status.Errorf(codes.AlreadyExists, "User already exists")
		case domain.ErrInternalError:
			return nil, status.Errorf(codes.Internal, "Internal error")
		}

		return nil, err
	}

	return &pb.CreateResponse{Code: pb.CodeResult_OK, UserId: int32(result)}, nil
}

//Gets all users7
func (s UserServer) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	result, err := s.appService.GetAll(ctx)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error")
	}

	response := pb.GetAllUsersResponse{Users: []*pb.User{}}

	for _, u := range result {
		usr, _ := mappers.ToGrpcUser(u)
		response.Users = append(response.Users, &usr)
	}

	return &response, nil
}

//Updates the user information
func (s UserServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	usr, _ := mappers.ToDomainUser(*req.GetUser())

	err := s.appService.Update(ctx, usr)

	if err != nil {
		switch err {
		case domain.ErrInvalidData:
			return nil, status.Errorf(codes.InvalidArgument, "Invalid data")
		case domain.ErrNotFound:
			return nil, status.Errorf(codes.InvalidArgument, "User does not exist")
		case domain.ErrInternalError:
			return nil, status.Errorf(codes.Internal, "Internal error")
		}

		return nil, err
	}

	return &pb.UpdateResponse{Code: pb.CodeResult_OK}, nil
}

//Deletes a user
func (s UserServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {

	id := req.GetId()

	err := s.appService.Delete(ctx, int(id))

	if err != nil {
		switch err {
		case domain.ErrInvalidData:
			return nil, status.Errorf(codes.InvalidArgument, "invalid id")
		case domain.ErrNotFound:
			return nil, status.Errorf(codes.InvalidArgument, "User does not exist")
		case domain.ErrInternalError:
			return nil, status.Errorf(codes.Internal, "Internal error")
		}

		return nil, err
	}

	return &pb.DeleteResponse{Code: pb.CodeResult_OK}, nil
}
