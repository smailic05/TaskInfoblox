package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/smailic05/TaskInfoblox/internal/mock"
	"github.com/smailic05/TaskInfoblox/internal/model"
	"github.com/smailic05/TaskInfoblox/internal/pb"
	"github.com/smailic05/TaskInfoblox/internal/service"
)

const (
	testValue       = "TEST"
	newTestValue    = "NEW_TEST"
	testWrongValue  = "&][()%$!?/"
	ErrUserNotExist = "The user does not exist"
	ErrUserExist    = "Error, user already exists"
	Success         = "Success"
	Deleted         = "Deleted"
)

var user = model.User{
	Username: testValue,
	Address:  testValue,
	Phone:    testValue,
}

var newUser = model.User{
	Username: newTestValue,
	Address:  newTestValue,
	Phone:    newTestValue,
}

type serviceTestSuite struct {
	suite.Suite
	srvFind  *mock.UserService_FindUserServer
	srvList  *mock.UserService_ListUserServer
	repoMock *mock.Repository
	service  *service.UserService
}

func (suite *serviceTestSuite) SetupTest() {
	repo := &mock.Repository{}
	s := service.New(repo)
	suite.repoMock = repo
	suite.srvFind = &mock.UserService_FindUserServer{}
	suite.srvList = &mock.UserService_ListUserServer{}
	suite.service = s
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (suite *serviceTestSuite) TestAddUser() {
	addUserRequest := &pb.AddUserRequest{Username: testValue, Address: testValue, Phone: testValue}
	suite.repoMock.On("LoadOne", user).Once().Return(user, gorm.ErrRecordNotFound)
	suite.repoMock.On("Store", user).Once().Return(nil)
	r, err := suite.service.AddUser(context.Background(), addUserRequest)
	suite.NoError(err)
	suite.Equal(r.GetResponse(), Success)
}

func TestAddUserErrorExist(t *testing.T) {
	//TODO
}

func (suite *serviceTestSuite) TestDeleteUser() {
	deleteUserRequest := &pb.DeleteUserRequest{Username: testValue, Address: testValue, Phone: testValue}
	suite.repoMock.On("DeleteUser", user).Once().Return(nil)
	resp, err := suite.service.DeleteUser(context.Background(), deleteUserRequest)
	suite.NotNil(resp)
	suite.NoError(err)
}

func (suite *serviceTestSuite) TestFindUserNotExist() {
	findUserRequest := &pb.FindUserRequest{Username: testValue, Address: testValue, Phone: testValue}
	suite.repoMock.On("Load", user).Once().Return([]model.User{}, nil)
	err := suite.service.FindUser(findUserRequest, suite.srvFind)
	suite.NotNil(err)
}

func (suite *serviceTestSuite) TestFindUser() {
	findUserRequest := &pb.FindUserRequest{Username: testValue, Address: testValue, Phone: testValue}
	FindUserResponse := &pb.FindUserResponse{Username: testValue, Address: testValue, Phone: testValue}
	users := []model.User{user}
	suite.repoMock.On("Load", user).Once().Return(users, nil)
	suite.srvFind.On("Send", FindUserResponse).Return(nil)
	err := suite.service.FindUser(findUserRequest, suite.srvFind)
	suite.NoError(err)
}

func (suite *serviceTestSuite) TestUpdateUserNotExist() {
	updateUserRequest := &pb.UpdateUserRequest{
		OldUsername: user.Username,
		OldAddress:  user.Address,
		OldPhone:    user.Phone,
		NewUsername: newTestValue,
		NewAddress:  newTestValue,
		NewPhone:    newTestValue,
	}
	suite.repoMock.On("LoadOne", user).Once().Return(user, gorm.ErrRecordNotFound)
	_, err := suite.service.UpdateUser(context.Background(), updateUserRequest)
	suite.Equal(err, gorm.ErrRecordNotFound)
}

func (suite *serviceTestSuite) TestUpdateUser() {
	updateUserRequest := &pb.UpdateUserRequest{
		OldUsername: user.Username,
		OldAddress:  user.Address,
		OldPhone:    user.Phone,
		NewUsername: newTestValue,
		NewAddress:  newTestValue,
		NewPhone:    newTestValue,
	}
	suite.repoMock.On("LoadOne", user).Once().Return(user, nil)
	suite.repoMock.On("Update", newUser).Return(nil)
	resp, err := suite.service.UpdateUser(context.Background(), updateUserRequest)
	suite.NotNil(resp)
	suite.Equal(resp.Response, Success)
	suite.NoError(err)
}

func (suite *serviceTestSuite) TestListUser() {
	listUserRequest := &pb.ListUserRequest{}
	listUserResponse := &pb.ListUserResponse{
		Username: testValue,
		Address:  testValue,
		Phone:    testValue}
	users := []model.User{user}
	suite.repoMock.On("Load", model.User{Address: "%", Username: "%", Phone: "%"}).Return(users, nil)
	suite.srvList.On("Send", listUserResponse).Return(nil)
	err := suite.service.ListUser(listUserRequest, suite.srvList)
	suite.NoError(err)
}
