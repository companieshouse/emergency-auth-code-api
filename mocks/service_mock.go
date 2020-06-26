// Code generated by MockGen. DO NOT EDIT.
// Source: dao/service.go

package mocks

import (
	models "github.com/companieshouse/emergency-auth-code-api/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAuthcodeDAOService is a mock of AuthcodeDAOService interface
type MockAuthcodeDAOService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthcodeDAOServiceMockRecorder
}

// MockAuthcodeDAOServiceMockRecorder is the mock recorder for MockAuthcodeDAOService
type MockAuthcodeDAOServiceMockRecorder struct {
	mock *MockAuthcodeDAOService
}

// NewMockAuthcodeDAOService creates a new mock instance
func NewMockAuthcodeDAOService(ctrl *gomock.Controller) *MockAuthcodeDAOService {
	mock := &MockAuthcodeDAOService{ctrl: ctrl}
	mock.recorder = &MockAuthcodeDAOServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuthcodeDAOService) EXPECT() *MockAuthcodeDAOServiceMockRecorder {
	return m.recorder
}

// CompanyHasAuthCode mocks base method
func (m *MockAuthcodeDAOService) CompanyHasAuthCode(companyNumber string) (bool, error) {
	ret := m.ctrl.Call(m, "CompanyHasAuthCode", companyNumber)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CompanyHasAuthCode indicates an expected call of CompanyHasAuthCode
func (mr *MockAuthcodeDAOServiceMockRecorder) CompanyHasAuthCode(companyNumber interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompanyHasAuthCode", reflect.TypeOf((*MockAuthcodeDAOService)(nil).CompanyHasAuthCode), companyNumber)
}

// UpsertEmptyAuthCode mocks base method
func (m *MockAuthcodeDAOService) UpsertEmptyAuthCode(companyNumber string) error {
	ret := m.ctrl.Call(m, "UpsertEmptyAuthCode", companyNumber)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertEmptyAuthCode indicates an expected call of UpsertEmptyAuthCode
func (mr *MockAuthcodeDAOServiceMockRecorder) UpsertEmptyAuthCode(companyNumber interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertEmptyAuthCode", reflect.TypeOf((*MockAuthcodeDAOService)(nil).UpsertEmptyAuthCode), companyNumber)
}

// MockAuthcodeRequestDAOService is a mock of AuthcodeRequestDAOService interface
type MockAuthcodeRequestDAOService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthcodeRequestDAOServiceMockRecorder
}

// MockAuthcodeRequestDAOServiceMockRecorder is the mock recorder for MockAuthcodeRequestDAOService
type MockAuthcodeRequestDAOServiceMockRecorder struct {
	mock *MockAuthcodeRequestDAOService
}

// NewMockAuthcodeRequestDAOService creates a new mock instance
func NewMockAuthcodeRequestDAOService(ctrl *gomock.Controller) *MockAuthcodeRequestDAOService {
	mock := &MockAuthcodeRequestDAOService{ctrl: ctrl}
	mock.recorder = &MockAuthcodeRequestDAOServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuthcodeRequestDAOService) EXPECT() *MockAuthcodeRequestDAOServiceMockRecorder {
	return m.recorder
}

// InsertAuthCodeRequest mocks base method
func (m *MockAuthcodeRequestDAOService) InsertAuthCodeRequest(dao *models.AuthCodeRequestResourceDao) error {
	ret := m.ctrl.Call(m, "InsertAuthCodeRequest", dao)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertAuthCodeRequest indicates an expected call of InsertAuthCodeRequest
func (mr *MockAuthcodeRequestDAOServiceMockRecorder) InsertAuthCodeRequest(dao interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertAuthCodeRequest", reflect.TypeOf((*MockAuthcodeRequestDAOService)(nil).InsertAuthCodeRequest), dao)
}

// GetAuthCodeRequest mocks base method
func (m *MockAuthcodeRequestDAOService) GetAuthCodeRequest(authCodeRequestID string) (*models.AuthCodeRequestResourceDao, error) {
	ret := m.ctrl.Call(m, "GetAuthCodeRequest", authCodeRequestID)
	ret0, _ := ret[0].(*models.AuthCodeRequestResourceDao)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAuthCodeRequest indicates an expected call of GetAuthCodeRequest
func (mr *MockAuthcodeRequestDAOServiceMockRecorder) GetAuthCodeRequest(authCodeRequestID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthCodeRequest", reflect.TypeOf((*MockAuthcodeRequestDAOService)(nil).GetAuthCodeRequest), authCodeRequestID)
}

// UpdateAuthCodeRequestOfficer mocks base method
func (m *MockAuthcodeRequestDAOService) UpdateAuthCodeRequestOfficer(dao *models.AuthCodeRequestResourceDao) error {
	ret := m.ctrl.Call(m, "UpdateAuthCodeRequestOfficer", dao)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAuthCodeRequestOfficer indicates an expected call of UpdateAuthCodeRequestOfficer
func (mr *MockAuthcodeRequestDAOServiceMockRecorder) UpdateAuthCodeRequestOfficer(dao interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAuthCodeRequestOfficer", reflect.TypeOf((*MockAuthcodeRequestDAOService)(nil).UpdateAuthCodeRequestOfficer), dao)
}

// UpdateAuthCodeRequestStatus mocks base method
func (m *MockAuthcodeRequestDAOService) UpdateAuthCodeRequestStatus(dao *models.AuthCodeRequestResourceDao) error {
	ret := m.ctrl.Call(m, "UpdateAuthCodeRequestStatus", dao)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAuthCodeRequestStatus indicates an expected call of UpdateAuthCodeRequestStatus
func (mr *MockAuthcodeRequestDAOServiceMockRecorder) UpdateAuthCodeRequestStatus(dao interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAuthCodeRequestStatus", reflect.TypeOf((*MockAuthcodeRequestDAOService)(nil).UpdateAuthCodeRequestStatus), dao)
}
