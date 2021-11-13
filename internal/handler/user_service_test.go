package handler

import (
	"context"
	"testing"

	"github.com/smailic05/TaskInfoblox/internal/pb"
	"google.golang.org/grpc"
)

const (
	testValue      = "TEST"
	newTestValue   = "NEW_TEST"
	testWrongValue = "&][()%$!?/"
)

type UserService_FindUserServer_Mock struct {
	grpc.ServerStream
}

func (srv UserService_FindUserServer_Mock) Send(*pb.FindUserResponse) error {
	return nil
}

type UserService_ListUserServer_Mock struct {
	grpc.ServerStream
}

func (srv UserService_ListUserServer_Mock) Send(*pb.ListUserResponse) error {
	return nil
}

func TestAddUser(t *testing.T) {
	us := New()
	addUserRequest := pb.AddUserRequest{Username: testValue, Address: testValue, Phone: testValue}
	resp, err := us.AddUser(context.Background(), &addUserRequest)
	if resp.Response != Success || err != nil {
		t.Fatalf("func AddUser returned an error")
	}
}

func TestAddUserErrorExist(t *testing.T) {
	us := New()
	user := User{Username: testValue, Address: testValue, Phone: testValue}
	us.UserSlice = append(us.UserSlice, user)
	addUserRequest := pb.AddUserRequest{Username: testValue, Address: testValue, Phone: testValue}
	resp, err := us.AddUser(context.Background(), &addUserRequest)
	if resp != nil || err == nil {
		t.Fatalf("func AddUser was supposed to return an error")
	}
}

func TestDeleteUser(t *testing.T) {
	us := New()
	user := User{Username: testValue, Address: testValue, Phone: testValue}
	us.UserSlice = append(us.UserSlice, user)
	deleteUserRequest := pb.DeleteUserRequest{Username: testValue, Address: testValue, Phone: testValue}
	resp, err := us.DeleteUser(context.Background(), &deleteUserRequest)
	if resp.Response != Deleted || err != nil {
		t.Fatalf("func DeleteUser returned an error")
	}
}

func TestDeleteUserNotExist(t *testing.T) {
	us := New()
	deleteUserRequest := pb.DeleteUserRequest{Username: testValue, Address: testValue, Phone: testValue}
	resp, err := us.DeleteUser(context.Background(), &deleteUserRequest)
	if resp != nil || err == nil {
		t.Fatalf("func DeleteUser was supposed to return an error")
	}
}

func TestFindUserNotExist(t *testing.T) {
	us := New()
	findUserRequest := pb.FindUserRequest{Username: testValue, Address: testValue, Phone: testValue}
	srv := UserService_FindUserServer_Mock{}
	err := us.FindUser(&findUserRequest, srv)
	if err == nil {
		t.Fatalf("func FindUser was supposed to return an error")
	}
}

func TestFindUser(t *testing.T) {
	us := New()
	user := User{Username: testValue, Address: testValue, Phone: testValue}
	us.UserSlice = append(us.UserSlice, user)
	findUserRequest := pb.FindUserRequest{Username: testValue, Address: testValue, Phone: testValue}
	srv := UserService_FindUserServer_Mock{}
	err := us.FindUser(&findUserRequest, srv)
	if err != nil {
		t.Fatalf("func FindUser returned an error")
	}
}

func TestUpdateUserNotExist(t *testing.T) {
	us := New()
	updateUserRequest := pb.UpdateUserRequest{}
	resp, err := us.UpdateUser(context.Background(), &updateUserRequest)
	if resp != nil || err == nil {
		t.Fatalf("func UpdateUser was supposed to return an error")
	}
}

func TestUpdateUser(t *testing.T) {
	us := New()
	user := User{Username: testValue, Address: testValue, Phone: testValue}
	us.UserSlice = append(us.UserSlice, user)
	updateUserRequest := pb.UpdateUserRequest{
		OldUsername: user.Username,
		OldAddress:  user.Address,
		OldPhone:    user.Phone,
		NewUsername: newTestValue,
		NewAddress:  newTestValue,
		NewPhone:    newTestValue,
	}
	resp, err := us.UpdateUser(context.Background(), &updateUserRequest)
	if resp.Response != Success && err != nil {
		t.Fatalf("func FindUser returned an error")
	}
}

func TestListUser(t *testing.T) {
	us := New()
	user := User{Username: testValue, Address: testValue, Phone: testValue}
	us.UserSlice = append(us.UserSlice, user)
	listUserRequest := pb.ListUserRequest{}
	srv := UserService_ListUserServer_Mock{}
	err := us.ListUser(&listUserRequest, srv)
	if err != nil {
		t.Fatalf("func ListUser returned an error")
	}
}

func TestCompileRequestUsername(t *testing.T) {
	_, err := compileRequest(testWrongValue, testValue, testValue)
	if err == nil {
		t.Fatalf("func CompileRequest was supposed to return an error")
	}
}

func TestCompileRequestAddress(t *testing.T) {
	_, err := compileRequest(testValue, testWrongValue, testValue)
	if err == nil {
		t.Fatalf("func CompileRequest was supposed to return an error")
	}
}

func TestCompileRequestPhone(t *testing.T) {
	_, err := compileRequest(testValue, testValue, testWrongValue)
	if err == nil {
		t.Fatalf("func CompileRequest was supposed to return an error")
	}
}
