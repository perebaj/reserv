// Code generated by MockGen. DO NOT EDIT.
// Source: property.go
//
// Generated by this command:
//
//	mockgen -source property.go -destination ../mock/property.go -package mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	reserv "github.com/perebaj/reserv"
	gomock "go.uber.org/mock/gomock"
)

// MockPropertyRepository is a mock of PropertyRepository interface.
type MockPropertyRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPropertyRepositoryMockRecorder
}

// MockPropertyRepositoryMockRecorder is the mock recorder for MockPropertyRepository.
type MockPropertyRepositoryMockRecorder struct {
	mock *MockPropertyRepository
}

// NewMockPropertyRepository creates a new mock instance.
func NewMockPropertyRepository(ctrl *gomock.Controller) *MockPropertyRepository {
	mock := &MockPropertyRepository{ctrl: ctrl}
	mock.recorder = &MockPropertyRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPropertyRepository) EXPECT() *MockPropertyRepositoryMockRecorder {
	return m.recorder
}

// Amenities mocks base method.
func (m *MockPropertyRepository) Amenities(ctx context.Context) ([]reserv.Amenity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Amenities", ctx)
	ret0, _ := ret[0].([]reserv.Amenity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Amenities indicates an expected call of Amenities.
func (mr *MockPropertyRepositoryMockRecorder) Amenities(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Amenities", reflect.TypeOf((*MockPropertyRepository)(nil).Amenities), ctx)
}

// CreateImage mocks base method.
func (m *MockPropertyRepository) CreateImage(ctx context.Context, image reserv.PropertyImage) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateImage", ctx, image)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateImage indicates an expected call of CreateImage.
func (mr *MockPropertyRepositoryMockRecorder) CreateImage(ctx, image any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateImage", reflect.TypeOf((*MockPropertyRepository)(nil).CreateImage), ctx, image)
}

// CreateProperty mocks base method.
func (m *MockPropertyRepository) CreateProperty(ctx context.Context, property reserv.Property) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProperty", ctx, property)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProperty indicates an expected call of CreateProperty.
func (mr *MockPropertyRepositoryMockRecorder) CreateProperty(ctx, property any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProperty", reflect.TypeOf((*MockPropertyRepository)(nil).CreateProperty), ctx, property)
}

// CreatePropertyAmenities mocks base method.
func (m *MockPropertyRepository) CreatePropertyAmenities(ctx context.Context, propertyID string, amenities []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePropertyAmenities", ctx, propertyID, amenities)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePropertyAmenities indicates an expected call of CreatePropertyAmenities.
func (mr *MockPropertyRepositoryMockRecorder) CreatePropertyAmenities(ctx, propertyID, amenities any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePropertyAmenities", reflect.TypeOf((*MockPropertyRepository)(nil).CreatePropertyAmenities), ctx, propertyID, amenities)
}

// DeleteImage mocks base method.
func (m *MockPropertyRepository) DeleteImage(ctx context.Context, imageID string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteImage", ctx, imageID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteImage indicates an expected call of DeleteImage.
func (mr *MockPropertyRepositoryMockRecorder) DeleteImage(ctx, imageID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteImage", reflect.TypeOf((*MockPropertyRepository)(nil).DeleteImage), ctx, imageID)
}

// DeleteProperty mocks base method.
func (m *MockPropertyRepository) DeleteProperty(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProperty", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProperty indicates an expected call of DeleteProperty.
func (mr *MockPropertyRepositoryMockRecorder) DeleteProperty(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProperty", reflect.TypeOf((*MockPropertyRepository)(nil).DeleteProperty), ctx, id)
}

// GetProperty mocks base method.
func (m *MockPropertyRepository) GetProperty(ctx context.Context, id string) (int, reserv.Property, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProperty", ctx, id)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(reserv.Property)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetProperty indicates an expected call of GetProperty.
func (mr *MockPropertyRepositoryMockRecorder) GetProperty(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProperty", reflect.TypeOf((*MockPropertyRepository)(nil).GetProperty), ctx, id)
}

// GetPropertyAmenities mocks base method.
func (m *MockPropertyRepository) GetPropertyAmenities(ctx context.Context, propertyID string) ([]reserv.Amenity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPropertyAmenities", ctx, propertyID)
	ret0, _ := ret[0].([]reserv.Amenity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPropertyAmenities indicates an expected call of GetPropertyAmenities.
func (mr *MockPropertyRepositoryMockRecorder) GetPropertyAmenities(ctx, propertyID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPropertyAmenities", reflect.TypeOf((*MockPropertyRepository)(nil).GetPropertyAmenities), ctx, propertyID)
}

// Properties mocks base method.
func (m *MockPropertyRepository) Properties(ctx context.Context, filter reserv.PropertyFilter) ([]reserv.Property, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Properties", ctx, filter)
	ret0, _ := ret[0].([]reserv.Property)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Properties indicates an expected call of Properties.
func (mr *MockPropertyRepositoryMockRecorder) Properties(ctx, filter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Properties", reflect.TypeOf((*MockPropertyRepository)(nil).Properties), ctx, filter)
}

// UpdateProperty mocks base method.
func (m *MockPropertyRepository) UpdateProperty(ctx context.Context, property reserv.Property, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProperty", ctx, property, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProperty indicates an expected call of UpdateProperty.
func (mr *MockPropertyRepositoryMockRecorder) UpdateProperty(ctx, property, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProperty", reflect.TypeOf((*MockPropertyRepository)(nil).UpdateProperty), ctx, property, id)
}
