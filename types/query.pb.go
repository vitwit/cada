// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sdk/avail/v1beta1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// query request
type QuerySubmittedBlobStatusRequest struct {
}

func (m *QuerySubmittedBlobStatusRequest) Reset()         { *m = QuerySubmittedBlobStatusRequest{} }
func (m *QuerySubmittedBlobStatusRequest) String() string { return proto.CompactTextString(m) }
func (*QuerySubmittedBlobStatusRequest) ProtoMessage()    {}
func (*QuerySubmittedBlobStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_30ff5d91ce731c68, []int{0}
}
func (m *QuerySubmittedBlobStatusRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QuerySubmittedBlobStatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QuerySubmittedBlobStatusRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QuerySubmittedBlobStatusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QuerySubmittedBlobStatusRequest.Merge(m, src)
}
func (m *QuerySubmittedBlobStatusRequest) XXX_Size() int {
	return m.Size()
}
func (m *QuerySubmittedBlobStatusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QuerySubmittedBlobStatusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QuerySubmittedBlobStatusRequest proto.InternalMessageInfo

// query response
type QuerySubmittedBlobStatusResponse struct {
	Range        *Range `protobuf:"bytes,1,opt,name=range,proto3" json:"range,omitempty"`
	Status       string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	ProvenHeight uint64 `protobuf:"varint,3,opt,name=proven_height,json=provenHeight,proto3" json:"proven_height,omitempty"`
	VotingEndsAt uint64 `protobuf:"varint,4,opt,name=voting_ends_at,json=votingEndsAt,proto3" json:"voting_ends_at,omitempty"`
}

func (m *QuerySubmittedBlobStatusResponse) Reset()         { *m = QuerySubmittedBlobStatusResponse{} }
func (m *QuerySubmittedBlobStatusResponse) String() string { return proto.CompactTextString(m) }
func (*QuerySubmittedBlobStatusResponse) ProtoMessage()    {}
func (*QuerySubmittedBlobStatusResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_30ff5d91ce731c68, []int{1}
}
func (m *QuerySubmittedBlobStatusResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QuerySubmittedBlobStatusResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QuerySubmittedBlobStatusResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QuerySubmittedBlobStatusResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QuerySubmittedBlobStatusResponse.Merge(m, src)
}
func (m *QuerySubmittedBlobStatusResponse) XXX_Size() int {
	return m.Size()
}
func (m *QuerySubmittedBlobStatusResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QuerySubmittedBlobStatusResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QuerySubmittedBlobStatusResponse proto.InternalMessageInfo

func (m *QuerySubmittedBlobStatusResponse) GetRange() *Range {
	if m != nil {
		return m.Range
	}
	return nil
}

func (m *QuerySubmittedBlobStatusResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *QuerySubmittedBlobStatusResponse) GetProvenHeight() uint64 {
	if m != nil {
		return m.ProvenHeight
	}
	return 0
}

func (m *QuerySubmittedBlobStatusResponse) GetVotingEndsAt() uint64 {
	if m != nil {
		return m.VotingEndsAt
	}
	return 0
}

func init() {
	proto.RegisterType((*QuerySubmittedBlobStatusRequest)(nil), "sdk.avail.v1beta1.QuerySubmittedBlobStatusRequest")
	proto.RegisterType((*QuerySubmittedBlobStatusResponse)(nil), "sdk.avail.v1beta1.QuerySubmittedBlobStatusResponse")
}

func init() { proto.RegisterFile("sdk/avail/v1beta1/query.proto", fileDescriptor_30ff5d91ce731c68) }

var fileDescriptor_30ff5d91ce731c68 = []byte{
	// 401 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xcf, 0x6e, 0xd4, 0x30,
	0x10, 0xc6, 0xd7, 0xa5, 0xad, 0x84, 0xf9, 0x23, 0x61, 0x10, 0x8a, 0x22, 0x48, 0xd3, 0x2d, 0x88,
	0x95, 0x50, 0x63, 0x75, 0xfb, 0x04, 0xad, 0x84, 0xc4, 0x95, 0xf4, 0xc6, 0x65, 0x65, 0x63, 0xe3,
	0xb5, 0x9a, 0x78, 0xd2, 0x78, 0x12, 0xe8, 0x95, 0x27, 0x40, 0xe2, 0x31, 0x38, 0xf7, 0x1d, 0x38,
	0x56, 0xe2, 0xc2, 0x11, 0xed, 0xf2, 0x20, 0x68, 0xed, 0x08, 0x0e, 0xbb, 0x08, 0xf5, 0x96, 0x99,
	0xef, 0x37, 0x5f, 0x3e, 0xcf, 0xd0, 0xa7, 0x5e, 0x9d, 0x73, 0xd1, 0x0b, 0x5b, 0xf1, 0xfe, 0x48,
	0x6a, 0x14, 0x47, 0xfc, 0xa2, 0xd3, 0xed, 0x65, 0xd1, 0xb4, 0x80, 0xc0, 0x1e, 0x78, 0x75, 0x5e,
	0x04, 0xb9, 0x18, 0xe4, 0xf4, 0x91, 0x01, 0x03, 0x41, 0xe5, 0xab, 0xaf, 0x08, 0xa6, 0x4f, 0x0c,
	0x80, 0xa9, 0x34, 0x17, 0x8d, 0xe5, 0xc2, 0x39, 0x40, 0x81, 0x16, 0x9c, 0x1f, 0xd4, 0xfd, 0xf5,
	0xbf, 0xf4, 0xa2, 0xb2, 0x4a, 0x20, 0xb4, 0x03, 0xb2, 0x37, 0x18, 0x84, 0x4a, 0x76, 0xef, 0x39,
	0xda, 0x5a, 0x7b, 0x14, 0x75, 0x33, 0x00, 0xe9, 0xba, 0x07, 0x7e, 0x8c, 0xda, 0x78, 0x9f, 0xee,
	0xbd, 0x59, 0xa5, 0x3e, 0xeb, 0x64, 0x6d, 0x11, 0xb5, 0x3a, 0xad, 0x40, 0x9e, 0xa1, 0xc0, 0xce,
	0x97, 0xfa, 0xa2, 0xd3, 0x1e, 0xc7, 0x57, 0x84, 0xe6, 0xff, 0x66, 0x7c, 0x03, 0xce, 0x6b, 0x56,
	0xd0, 0x9d, 0x56, 0x38, 0xa3, 0x13, 0x92, 0x93, 0xc9, 0x9d, 0x69, 0x52, 0xac, 0x3d, 0xbf, 0x28,
	0x57, 0x7a, 0x19, 0x31, 0xf6, 0x98, 0xee, 0xfa, 0xe0, 0x90, 0x6c, 0xe5, 0x64, 0x72, 0xbb, 0x1c,
	0x2a, 0x76, 0x40, 0xef, 0x35, 0x2d, 0xf4, 0xda, 0xcd, 0xe6, 0xda, 0x9a, 0x39, 0x26, 0xb7, 0x72,
	0x32, 0xd9, 0x2e, 0xef, 0xc6, 0xe6, 0xeb, 0xd0, 0x63, 0xcf, 0xe8, 0xfd, 0x1e, 0xd0, 0x3a, 0x33,
	0xd3, 0x4e, 0xf9, 0x99, 0xc0, 0x64, 0x3b, 0x52, 0xb1, 0xfb, 0xca, 0x29, 0x7f, 0x82, 0xd3, 0x2b,
	0x42, 0x77, 0x42, 0x6e, 0xf6, 0x95, 0xd0, 0x87, 0x1b, 0xc2, 0xb3, 0xe9, 0x86, 0x94, 0xff, 0xd9,
	0x46, 0x7a, 0x7c, 0xa3, 0x99, 0xb8, 0x9d, 0xf1, 0xcb, 0x4f, 0xdf, 0x7f, 0x7d, 0xd9, 0x7a, 0xce,
	0x0e, 0xe2, 0x19, 0x64, 0x05, 0xf2, 0xcf, 0x29, 0x7c, 0x98, 0xfb, 0x3b, 0x74, 0x7a, 0xf2, 0x6d,
	0x91, 0x91, 0xeb, 0x45, 0x46, 0x7e, 0x2e, 0x32, 0xf2, 0x79, 0x99, 0x8d, 0xae, 0x97, 0xd9, 0xe8,
	0xc7, 0x32, 0x1b, 0xbd, 0x7d, 0x61, 0x2c, 0xce, 0x3b, 0x59, 0xbc, 0x83, 0x9a, 0xf7, 0x16, 0x3f,
	0x58, 0x8c, 0x7e, 0x87, 0x4a, 0x1c, 0xd6, 0xa0, 0xba, 0x4a, 0x73, 0xbc, 0x6c, 0xb4, 0x97, 0xbb,
	0xe1, 0xb8, 0xc7, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xcd, 0x4b, 0xa4, 0x7d, 0xa4, 0x02, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// submit Blob Status
	SubmittedBlobStatus(ctx context.Context, in *QuerySubmittedBlobStatusRequest, opts ...grpc.CallOption) (*QuerySubmittedBlobStatusResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) SubmittedBlobStatus(ctx context.Context, in *QuerySubmittedBlobStatusRequest, opts ...grpc.CallOption) (*QuerySubmittedBlobStatusResponse, error) {
	out := new(QuerySubmittedBlobStatusResponse)
	err := c.cc.Invoke(ctx, "/sdk.avail.v1beta1.Query/SubmittedBlobStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// submit Blob Status
	SubmittedBlobStatus(context.Context, *QuerySubmittedBlobStatusRequest) (*QuerySubmittedBlobStatusResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) SubmittedBlobStatus(ctx context.Context, req *QuerySubmittedBlobStatusRequest) (*QuerySubmittedBlobStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmittedBlobStatus not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_SubmittedBlobStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySubmittedBlobStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).SubmittedBlobStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk.avail.v1beta1.Query/SubmittedBlobStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).SubmittedBlobStatus(ctx, req.(*QuerySubmittedBlobStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "sdk.avail.v1beta1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmittedBlobStatus",
			Handler:    _Query_SubmittedBlobStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sdk/avail/v1beta1/query.proto",
}

func (m *QuerySubmittedBlobStatusRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QuerySubmittedBlobStatusRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QuerySubmittedBlobStatusRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QuerySubmittedBlobStatusResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QuerySubmittedBlobStatusResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QuerySubmittedBlobStatusResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.VotingEndsAt != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.VotingEndsAt))
		i--
		dAtA[i] = 0x20
	}
	if m.ProvenHeight != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.ProvenHeight))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Status) > 0 {
		i -= len(m.Status)
		copy(dAtA[i:], m.Status)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Status)))
		i--
		dAtA[i] = 0x12
	}
	if m.Range != nil {
		{
			size, err := m.Range.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QuerySubmittedBlobStatusRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QuerySubmittedBlobStatusResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Range != nil {
		l = m.Range.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	l = len(m.Status)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	if m.ProvenHeight != 0 {
		n += 1 + sovQuery(uint64(m.ProvenHeight))
	}
	if m.VotingEndsAt != 0 {
		n += 1 + sovQuery(uint64(m.VotingEndsAt))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QuerySubmittedBlobStatusRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QuerySubmittedBlobStatusRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QuerySubmittedBlobStatusRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QuerySubmittedBlobStatusResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QuerySubmittedBlobStatusResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QuerySubmittedBlobStatusResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Range", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Range == nil {
				m.Range = &Range{}
			}
			if err := m.Range.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Status = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProvenHeight", wireType)
			}
			m.ProvenHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ProvenHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field VotingEndsAt", wireType)
			}
			m.VotingEndsAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.VotingEndsAt |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
