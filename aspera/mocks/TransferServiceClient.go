// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	transfersdk "github.com/IBM/ibmcloud-cos-cli/aspera/transfersdk"
)

// TransferServiceClient is an autogenerated mock type for the TransferServiceClient type
type TransferServiceClient struct {
	mock.Mock
}

// AddTransferPaths provides a mock function with given fields: ctx, in, opts
func (_m *TransferServiceClient) AddTransferPaths(ctx context.Context, in *transfersdk.TransferPathRequest, opts ...grpc.CallOption) (*transfersdk.TransferPathResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *transfersdk.TransferPathResponse
	if rf, ok := ret.Get(0).(func(context.Context, *transfersdk.TransferPathRequest, ...grpc.CallOption) *transfersdk.TransferPathResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transfersdk.TransferPathResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *transfersdk.TransferPathRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAPIVersion provides a mock function with given fields: ctx, in, opts
func (_m *TransferServiceClient) GetAPIVersion(ctx context.Context, in *transfersdk.APIVersionRequest, opts ...grpc.CallOption) (*transfersdk.APIVersionResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *transfersdk.APIVersionResponse
	if rf, ok := ret.Get(0).(func(context.Context, *transfersdk.APIVersionRequest, ...grpc.CallOption) *transfersdk.APIVersionResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transfersdk.APIVersionResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *transfersdk.APIVersionRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetInfo provides a mock function with given fields: ctx, in, opts
func (_m *TransferServiceClient) GetInfo(ctx context.Context, in *transfersdk.InstanceInfoRequest, opts ...grpc.CallOption) (*transfersdk.InstanceInfoResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *transfersdk.InstanceInfoResponse
	if rf, ok := ret.Get(0).(func(context.Context, *transfersdk.InstanceInfoRequest, ...grpc.CallOption) *transfersdk.InstanceInfoResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transfersdk.InstanceInfoResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *transfersdk.InstanceInfoRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsPeerReachable provides a mock function with given fields: ctx, in, opts
func (_m *TransferServiceClient) IsPeerReachable(ctx context.Context, in *transfersdk.PeerCheckRequest, opts ...grpc.CallOption) (*transfersdk.PeerCheckResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *transfersdk.PeerCheckResponse
	if rf, ok := ret.Get(0).(func(context.Context, *transfersdk.PeerCheckRequest, ...grpc.CallOption) *transfersdk.PeerCheckResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transfersdk.PeerCheckResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *transfersdk.PeerCheckRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LockPersistentTransfer provides a mock function with given fields: ctx, in, opts
func (_m *TransferServiceClient) LockPersistentTransfer(ctx context.Context, in *transfersdk.LockPersistentTransferRequest, opts ...grpc.CallOption) (*transfersdk.LockPersistentTransferResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *transfersdk.LockPersistentTransferResponse
	if rf, ok := ret.Get(0).(func(context.Context, *transfersdk.LockPersistentTransferRequest, ...grpc.CallOption) *transfersdk.LockPersistentTransferResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transfersdk.LockPersistentTransferResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *transfersdk.LockPersistentTransferRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ModifyTransfer provides a mock function with given fields: ctx, in, opts
func (_m *TransferServiceClient) ModifyTransfer(ctx context.Context, in *transfersdk.TransferModificationRequest, opts ...grpc.CallOption) (*transfersdk.TransferModificationResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *transfersdk.TransferModificationResponse
	if rf, ok := ret.Get(0).(func(context.Context, *transfersdk.TransferModificationRequest, ...grpc.CallOption) *transfersdk.TransferModificationResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transfersdk.TransferModificationResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *transfersdk.TransferModificationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MonitorTransfers provides a mock function with given fields: ctx, in, opts
func (_m *TransferServiceClient) MonitorTransfers(ctx context.Context, in *transfersdk.RegistrationRequest, opts ...grpc.CallOption) (transfersdk.TransferService_MonitorTransfersClient, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 transfersdk.TransferService_MonitorTransfersClient
	if rf, ok := ret.Get(0).(func(context.Context, *transfersdk.RegistrationRequest, ...grpc.CallOption) transfersdk.TransferService_MonitorTransfersClient); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(transfersdk.TransferService_MonitorTransfersClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *transfersdk.RegistrationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryTransfer provides a mock function with given fields: ctx, in, opts
func (_m *TransferServiceClient) QueryTransfer(ctx context.Context, in *transfersdk.TransferInfoRequest, opts ...grpc.CallOption) (*transfersdk.QueryTransferResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *transfersdk.QueryTransferResponse
	if rf, ok := ret.Get(0).(func(context.Context, *transfersdk.TransferInfoRequest, ...grpc.CallOption) *transfersdk.QueryTransferResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transfersdk.QueryTransferResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *transfersdk.TransferInfoRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadStream provides a mock function with given fields: ctx, in, opts
func (_m *TransferServiceClient) ReadStream(ctx context.Context, in *transfersdk.ReadStreamRequest, opts ...grpc.CallOption) (transfersdk.TransferService_ReadStreamClient, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 transfersdk.TransferService_ReadStreamClient
	if rf, ok := ret.Get(0).(func(context.Context, *transfersdk.ReadStreamRequest, ...grpc.CallOption) transfersdk.TransferService_ReadStreamClient); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(transfersdk.TransferService_ReadStreamClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *transfersdk.ReadStreamRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StartTransfer provides a mock function with given fields: ctx, in, opts
func (_m *TransferServiceClient) StartTransfer(ctx context.Context, in *transfersdk.TransferRequest, opts ...grpc.CallOption) (*transfersdk.StartTransferResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *transfersdk.StartTransferResponse
	if rf, ok := ret.Get(0).(func(context.Context, *transfersdk.TransferRequest, ...grpc.CallOption) *transfersdk.StartTransferResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transfersdk.StartTransferResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *transfersdk.TransferRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StartTransferWithMonitor provides a mock function with given fields: ctx, in, opts
func (_m *TransferServiceClient) StartTransferWithMonitor(ctx context.Context, in *transfersdk.TransferRequest, opts ...grpc.CallOption) (transfersdk.TransferService_StartTransferWithMonitorClient, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 transfersdk.TransferService_StartTransferWithMonitorClient
	if rf, ok := ret.Get(0).(func(context.Context, *transfersdk.TransferRequest, ...grpc.CallOption) transfersdk.TransferService_StartTransferWithMonitorClient); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(transfersdk.TransferService_StartTransferWithMonitorClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *transfersdk.TransferRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StopTransfer provides a mock function with given fields: ctx, in, opts
func (_m *TransferServiceClient) StopTransfer(ctx context.Context, in *transfersdk.StopTransferRequest, opts ...grpc.CallOption) (*transfersdk.StopTransferResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *transfersdk.StopTransferResponse
	if rf, ok := ret.Get(0).(func(context.Context, *transfersdk.StopTransferRequest, ...grpc.CallOption) *transfersdk.StopTransferResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transfersdk.StopTransferResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *transfersdk.StopTransferRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Validate provides a mock function with given fields: ctx, in, opts
func (_m *TransferServiceClient) Validate(ctx context.Context, in *transfersdk.ValidationRequest, opts ...grpc.CallOption) (*transfersdk.ValidationResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *transfersdk.ValidationResponse
	if rf, ok := ret.Get(0).(func(context.Context, *transfersdk.ValidationRequest, ...grpc.CallOption) *transfersdk.ValidationResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transfersdk.ValidationResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *transfersdk.ValidationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteStream provides a mock function with given fields: ctx, opts
func (_m *TransferServiceClient) WriteStream(ctx context.Context, opts ...grpc.CallOption) (transfersdk.TransferService_WriteStreamClient, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 transfersdk.TransferService_WriteStreamClient
	if rf, ok := ret.Get(0).(func(context.Context, ...grpc.CallOption) transfersdk.TransferService_WriteStreamClient); ok {
		r0 = rf(ctx, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(transfersdk.TransferService_WriteStreamClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteStreamChunk provides a mock function with given fields: ctx, opts
func (_m *TransferServiceClient) WriteStreamChunk(ctx context.Context, opts ...grpc.CallOption) (transfersdk.TransferService_WriteStreamChunkClient, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 transfersdk.TransferService_WriteStreamChunkClient
	if rf, ok := ret.Get(0).(func(context.Context, ...grpc.CallOption) transfersdk.TransferService_WriteStreamChunkClient); ok {
		r0 = rf(ctx, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(transfersdk.TransferService_WriteStreamChunkClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
