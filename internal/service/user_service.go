package service

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/smailic05/TaskInfoblox/internal/model"
	"github.com/smailic05/TaskInfoblox/internal/pb"
)

const (
	ErrUserNotExist = "The user does not exist"
	ErrUserExist    = "Error, user already exists"
	Success         = "Success"
	Deleted         = "Deleted"
)

type Repository interface {
	Load(user model.User) ([]model.User, error)
	LoadOne(user model.User) (model.User, error)
	Store(user model.User) error
	Update(user model.User) error
	DeleteUser(user model.User) error
}

type UserService struct {
	pb.UnimplementedUserServiceServer
	repo Repository
}

func New(repo Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) AddUser(ctx context.Context, addUser *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	user := model.User{
		Address:  addUser.Address,
		Username: addUser.Username,
		Phone:    addUser.Phone}
	user, err := s.repo.LoadOne(user)
	if err != gorm.ErrRecordNotFound {
		return nil, status.Errorf(codes.AlreadyExists, "%s: when adding user", ErrUserExist)
	}
	err = s.repo.Store(user)
	if err != nil {
		return nil, err
	}
	return &pb.AddUserResponse{Response: Success}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, deleteUser *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	user := model.User{
		Address:  deleteUser.Address,
		Username: deleteUser.Username,
		Phone:    deleteUser.Phone}
	err := s.repo.DeleteUser(user)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{}, nil
}

func (s *UserService) FindUser(findUser *pb.FindUserRequest, srv pb.UserService_FindUserServer) error {
	user := model.User{
		Address:  findUser.Address,
		Username: findUser.Username,
		Phone:    findUser.Phone}
	user = formatRequest(user)
	users, err := s.repo.Load(user)
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return status.Errorf(codes.NotFound, ErrUserNotExist)
	}
	for _, user := range users {
		err := srv.Send(&pb.FindUserResponse{
			Username: user.Username,
			Address:  user.Address,
			Phone:    user.Phone})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, updateUser *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user, err := s.repo.LoadOne(model.User{Username: updateUser.OldUsername, Address: updateUser.OldAddress, Phone: updateUser.OldPhone})
	if err != nil {
		return nil, err
	}
	user.Address = updateUser.NewAddress
	user.Username = updateUser.NewUsername
	user.Phone = updateUser.NewPhone
	err = s.repo.Update(user)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateUserResponse{Response: Success}, nil
}

func (s *UserService) ListUser(list *pb.ListUserRequest, srv pb.UserService_ListUserServer) error {
	user := model.User{
		Address:  "%",
		Username: "%",
		Phone:    "%"}
	users, err := s.repo.Load(user)
	if err != nil {
		return err
	}
	for _, value := range users {
		err := srv.Send(&pb.ListUserResponse{
			Username: value.Username,
			Address:  value.Address,
			Phone:    value.Phone})
		if err != nil {
			return status.Errorf(codes.Internal, "%s: when sending response", err.Error())
		}
	}
	return nil
}

func formatRequest(user model.User) model.User {
	if user.Address == "" {
		user.Address = "%"
	}
	if user.Username == "" {
		user.Username = "%"
	}
	if user.Phone == "" {
		user.Phone = "%"
	}
	user.Address = strings.ReplaceAll(user.Address, "*", "%")
	user.Address = strings.ReplaceAll(user.Address, "?", "_")
	user.Username = strings.ReplaceAll(user.Username, "*", "%")
	user.Username = strings.ReplaceAll(user.Username, "?", "_")
	user.Phone = strings.ReplaceAll(user.Phone, "*", "%")
	user.Phone = strings.ReplaceAll(user.Phone, "?", "_")
	return user
}
