package handler

import (
	"context"
	"sync"

	"github.com/gobwas/glob"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smailic05/TaskInfoblox/internal/pb"
)

const (
	ErrUserNotExist = "The user does not exist"
	ErrUserExist    = "Error, user already exists"
	Success         = "Success"
	Deleted         = "Deleted"
)

type User struct {
	Address  string
	Username string
	Phone    string
}

type Compiled struct {
	Address  glob.Glob
	Username glob.Glob
	Phone    glob.Glob
}

type UserService struct {
	pb.UnimplementedUserServiceServer
	UserSlice []User
	mtx       sync.RWMutex
}

func New() *UserService {
	return &UserService{UserSlice: make([]User, 0)}
}

func (s *UserService) AddUser(ctx context.Context, addUser *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	user := User{
		Address:  addUser.Address,
		Username: addUser.Username,
		Phone:    addUser.Phone}
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if findExist(addUser, s.UserSlice) {
		return nil, status.Errorf(codes.InvalidArgument, ErrUserExist)
	}
	s.UserSlice = append(s.UserSlice, user)
	return &pb.AddUserResponse{Response: Success}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, deleteUser *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	compiled, err := compileRequest(deleteUser.Username, deleteUser.Address, deleteUser.Phone)
	if err != nil {
		return nil, err
	}
	count := 0
	s.mtx.Lock()
	for key := 0; key < len(s.UserSlice); key++ {
		value := s.UserSlice[key]
		if compiled.Username.Match(value.Username) && compiled.Address.Match(value.Address) && compiled.Phone.Match(value.Phone) {
			s.UserSlice = remove(s.UserSlice, key)
			count++
			key--
		}
	}
	s.mtx.Unlock()
	if count == 0 {
		return nil, status.Errorf(codes.InvalidArgument, ErrUserNotExist)
	}
	return &pb.DeleteUserResponse{Response: Deleted}, nil
}

func (s *UserService) FindUser(findUser *pb.FindUserRequest, srv pb.UserService_FindUserServer) error {
	compiled, err := compileRequest(findUser.Username, findUser.Address, findUser.Phone)
	if err != nil {
		return err
	}
	count := 0
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	for _, value := range s.UserSlice {
		if compiled.Username.Match(value.Username) && compiled.Address.Match(value.Address) && compiled.Phone.Match(value.Phone) {
			count++
			err := srv.Send(&pb.FindUserResponse{
				Username: value.Username,
				Address:  value.Address,
				Phone:    value.Phone})
			if err != nil {
				return status.Errorf(codes.InvalidArgument, "%s: when sending response", err.Error())
			}
		}
	}
	if count == 0 {
		return status.Errorf(codes.InvalidArgument, ErrUserNotExist)
	}
	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, updateUser *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	for key, value := range s.UserSlice {
		if updateUser.GetOldAddress() == value.Address && updateUser.GetOldPhone() == value.Phone && updateUser.GetOldUsername() == value.Username {
			s.UserSlice[key] = User{
				Address:  updateUser.GetNewAddress(),
				Username: updateUser.GetNewUsername(),
				Phone:    updateUser.GetNewPhone()}
			return &pb.UpdateUserResponse{Response: Success}, nil
		}
	}
	return nil, status.Errorf(codes.InvalidArgument, ErrUserNotExist)
}

func (s *UserService) ListUser(list *pb.ListUserRequest, srv pb.UserService_ListUserServer) error {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	for _, value := range s.UserSlice {
		err := srv.Send(&pb.ListUserResponse{
			Username: value.Username,
			Address:  value.Address,
			Phone:    value.Phone})
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "%s: when sending response", err.Error())
		}
	}
	return nil
}

func remove(s []User, i int) []User {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func findExist(addUser *pb.AddUserRequest, user []User) bool {
	for _, v := range user {
		if v.Address == addUser.Address && v.Phone == addUser.Phone && v.Username == addUser.Username {
			return true
		}
	}
	return false
}

func compileRequest(username, address, phone string) (*Compiled, error) {
	if username == "" {
		username = "*"
	}
	nameCompiled, err := glob.Compile(username)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: when parsing username", err.Error())
	}
	if address == "" {
		address = "*"
	}
	addressCompiled, err := glob.Compile(address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: when parsing address", err.Error())
	}
	if phone == "" {
		phone = "*"
	}
	phoneCompiled, err := glob.Compile(phone)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: when parsing phone", err.Error())
	}
	return &Compiled{Username: nameCompiled, Address: addressCompiled, Phone: phoneCompiled}, nil
}
