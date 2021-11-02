package handler

import (
	"context"

	"github.com/gobwas/glob"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smailic05/TaskInfoblox/internal/pb"
)

const (
	AddUserPath = "/add"
)

type User struct {
	Address  string
	Username string
	Phone    string
}

type UserService struct {
	pb.UnimplementedUserServiceServer
	UserSlice []User
}

func New() *UserService {
	return &UserService{UserSlice: make([]User, 0)}
}

func (s *UserService) AddUser(ctx context.Context, addUser *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	user := User{
		Address:  addUser.Address,
		Username: addUser.Username,
		Phone:    addUser.Phone}
	s.UserSlice = append(s.UserSlice, user)
	return &pb.AddUserResponse{Response: "Success"}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, deleteUser *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	name, err := glob.Compile(deleteUser.Username)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error")
	}
	address, err := glob.Compile(deleteUser.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error")
	}
	phone, err := glob.Compile(deleteUser.Phone)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error")
	}

	count := 0
	for key := 0; key < len(s.UserSlice); key++ {
		value := s.UserSlice[key]
		if name.Match(value.Username) && address.Match(value.Address) && phone.Match(value.Phone) {
			s.UserSlice = remove(s.UserSlice, key)
			count++
			key--
		}
	}
	if count == 0 {
		return &pb.DeleteUserResponse{Response: "The user does not exist"}, nil
	}
	return &pb.DeleteUserResponse{Response: "Deleted"}, nil
}

func (s *UserService) FindUser(findUser *pb.FindUserRequest, srv pb.UserService_FindUserServer) error {
	name, err := glob.Compile(findUser.Username)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Error")
	}
	address, err := glob.Compile(findUser.Address)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Error")
	}
	phone, err := glob.Compile(findUser.Phone)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Error")
	}
	count := 0
	for _, value := range s.UserSlice {
		if name.Match(value.Username) && address.Match(value.Address) && phone.Match(value.Phone) {
			count++
			err := srv.Send(&pb.FindUserResponse{
				Username: value.Username,
				Address:  value.Address,
				Phone:    value.Phone})
			if err != nil {
				return err
			}
		}
	}
	if count == 0 {
		return status.Errorf(codes.InvalidArgument, "The user does not exist")
	}
	return nil
}

func remove(s []User, i int) []User {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
