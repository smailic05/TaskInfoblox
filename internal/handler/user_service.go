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
	Users []User
}

func New() *UserService {
	return &UserService{Users: make([]User, 0)}
}

func (s *UserService) AddUser(ctx context.Context, addUser *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	user := User{
		Address:  addUser.Address,
		Username: addUser.Username,
		Phone:    addUser.Phone}
	s.Users = append(s.Users, user)
	return &pb.AddUserResponse{Response: "Success"}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, deleteUser *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	for key, value := range s.Users {
		if value.Username == deleteUser.Username {
			s.Users = remove(s.Users, key)
			return &pb.DeleteUserResponse{Response: "Deleted"}, nil
		}
	}
	return &pb.DeleteUserResponse{Response: "The user does not exist"}, nil
}
func (s *UserService) FindUser(ctx context.Context, findUser *pb.FindUserRequest) (*pb.FindUserResponse, error) {
	g, err := glob.Compile(findUser.Username)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error")
	}
	for _, value := range s.Users {
		if g.Match(value.Username) {
			return &pb.FindUserResponse{
				Username: value.Username,
				Address:  value.Address,
				Phone:    value.Phone}, nil
		}
	}
	return nil, status.Errorf(codes.InvalidArgument, "The user does not exist")
}

func remove(s []User, i int) []User {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
