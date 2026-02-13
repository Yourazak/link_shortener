package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

type MockRow struct {
	ctrl     *gomock.Controller
	recorder *MockRowMockRecorder
}

type MockRowMockRecorder struct {
	mock *MockRow
}

func NewMockRow(ctrl *gomock.Controller) *MockRow {
	mock := &MockRow{ctrl: ctrl}
	mock.recorder = &MockRowMockRecorder{mock}
	return mock
}

func (m *MockRow) EXPECT() *MockRowMockRecorder {
	return m.recorder
}

func (m *MockRow) Scan(arg0 ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Scan", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockRowMockRecorder) Scan(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockRow)(nil).Scan), arg0...)
}
