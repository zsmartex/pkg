// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/order.proto

package GrpcOrder

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	utils "github.com/zsmartex/pkg/Grpc/utils"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type OrderKey struct {
	Id                   int64                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Symbol               string               `protobuf:"bytes,2,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Side                 string               `protobuf:"bytes,3,opt,name=side,proto3" json:"side,omitempty"`
	Uuid                 []byte               `protobuf:"bytes,4,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Price                *utils.Decimal       `protobuf:"bytes,5,opt,name=price,proto3" json:"price,omitempty"`
	StopPrice            *utils.Decimal       `protobuf:"bytes,6,opt,name=stop_price,json=stopPrice,proto3" json:"stop_price,omitempty"`
	Fake                 bool                 `protobuf:"varint,7,opt,name=fake,proto3" json:"fake,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,8,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *OrderKey) Reset()         { *m = OrderKey{} }
func (m *OrderKey) String() string { return proto.CompactTextString(m) }
func (*OrderKey) ProtoMessage()    {}
func (*OrderKey) Descriptor() ([]byte, []int) {
	return fileDescriptor_f65b0626cc3aada8, []int{0}
}

func (m *OrderKey) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrderKey.Unmarshal(m, b)
}
func (m *OrderKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrderKey.Marshal(b, m, deterministic)
}
func (m *OrderKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderKey.Merge(m, src)
}
func (m *OrderKey) XXX_Size() int {
	return xxx_messageInfo_OrderKey.Size(m)
}
func (m *OrderKey) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderKey.DiscardUnknown(m)
}

var xxx_messageInfo_OrderKey proto.InternalMessageInfo

func (m *OrderKey) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *OrderKey) GetSymbol() string {
	if m != nil {
		return m.Symbol
	}
	return ""
}

func (m *OrderKey) GetSide() string {
	if m != nil {
		return m.Side
	}
	return ""
}

func (m *OrderKey) GetUuid() []byte {
	if m != nil {
		return m.Uuid
	}
	return nil
}

func (m *OrderKey) GetPrice() *utils.Decimal {
	if m != nil {
		return m.Price
	}
	return nil
}

func (m *OrderKey) GetStopPrice() *utils.Decimal {
	if m != nil {
		return m.StopPrice
	}
	return nil
}

func (m *OrderKey) GetFake() bool {
	if m != nil {
		return m.Fake
	}
	return false
}

func (m *OrderKey) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

type Order struct {
	Id                   int64                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Uuid                 []byte               `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`
	MemberId             int64                `protobuf:"varint,3,opt,name=member_id,json=memberId,proto3" json:"member_id,omitempty"`
	Symbol               string               `protobuf:"bytes,4,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Side                 string               `protobuf:"bytes,5,opt,name=side,proto3" json:"side,omitempty"`
	Type                 string               `protobuf:"bytes,6,opt,name=type,proto3" json:"type,omitempty"`
	Price                *utils.Decimal       `protobuf:"bytes,7,opt,name=price,proto3" json:"price,omitempty"`
	StopPrice            *utils.Decimal       `protobuf:"bytes,8,opt,name=stop_price,json=stopPrice,proto3" json:"stop_price,omitempty"`
	Quantity             *utils.Decimal       `protobuf:"bytes,9,opt,name=quantity,proto3" json:"quantity,omitempty"`
	FilledQuantity       *utils.Decimal       `protobuf:"bytes,10,opt,name=filled_quantity,json=filledQuantity,proto3" json:"filled_quantity,omitempty"`
	Fake                 bool                 `protobuf:"varint,11,opt,name=fake,proto3" json:"fake,omitempty"`
	Cancelled            bool                 `protobuf:"varint,12,opt,name=cancelled,proto3" json:"cancelled,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,13,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Order) Reset()         { *m = Order{} }
func (m *Order) String() string { return proto.CompactTextString(m) }
func (*Order) ProtoMessage()    {}
func (*Order) Descriptor() ([]byte, []int) {
	return fileDescriptor_f65b0626cc3aada8, []int{1}
}

func (m *Order) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Order.Unmarshal(m, b)
}
func (m *Order) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Order.Marshal(b, m, deterministic)
}
func (m *Order) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Order.Merge(m, src)
}
func (m *Order) XXX_Size() int {
	return xxx_messageInfo_Order.Size(m)
}
func (m *Order) XXX_DiscardUnknown() {
	xxx_messageInfo_Order.DiscardUnknown(m)
}

var xxx_messageInfo_Order proto.InternalMessageInfo

func (m *Order) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Order) GetUuid() []byte {
	if m != nil {
		return m.Uuid
	}
	return nil
}

func (m *Order) GetMemberId() int64 {
	if m != nil {
		return m.MemberId
	}
	return 0
}

func (m *Order) GetSymbol() string {
	if m != nil {
		return m.Symbol
	}
	return ""
}

func (m *Order) GetSide() string {
	if m != nil {
		return m.Side
	}
	return ""
}

func (m *Order) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Order) GetPrice() *utils.Decimal {
	if m != nil {
		return m.Price
	}
	return nil
}

func (m *Order) GetStopPrice() *utils.Decimal {
	if m != nil {
		return m.StopPrice
	}
	return nil
}

func (m *Order) GetQuantity() *utils.Decimal {
	if m != nil {
		return m.Quantity
	}
	return nil
}

func (m *Order) GetFilledQuantity() *utils.Decimal {
	if m != nil {
		return m.FilledQuantity
	}
	return nil
}

func (m *Order) GetFake() bool {
	if m != nil {
		return m.Fake
	}
	return false
}

func (m *Order) GetCancelled() bool {
	if m != nil {
		return m.Cancelled
	}
	return false
}

func (m *Order) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func init() {
	proto.RegisterType((*OrderKey)(nil), "GrpcOrder.OrderKey")
	proto.RegisterType((*Order)(nil), "GrpcOrder.Order")
}

func init() { proto.RegisterFile("proto/order.proto", fileDescriptor_f65b0626cc3aada8) }

var fileDescriptor_f65b0626cc3aada8 = []byte{
	// 409 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0x4f, 0x8b, 0xd4, 0x30,
	0x18, 0xc6, 0x69, 0xa7, 0x9d, 0x6d, 0xde, 0x5d, 0x57, 0xcc, 0x41, 0xc2, 0x28, 0x58, 0xf6, 0xd4,
	0x83, 0xa4, 0xa8, 0x27, 0xd9, 0x93, 0x22, 0x88, 0x78, 0x50, 0x83, 0x5e, 0xbc, 0x94, 0xb4, 0xc9,
	0xd4, 0xb0, 0xcd, 0xa6, 0xa6, 0x29, 0x58, 0x3f, 0x80, 0x5f, 0xcb, 0xaf, 0x26, 0x4d, 0xa6, 0xb3,
	0xe3, 0x9f, 0x01, 0xf7, 0xf6, 0xf4, 0xc9, 0xf3, 0xa4, 0xc9, 0xef, 0x25, 0x70, 0xaf, 0xb7, 0xc6,
	0x99, 0xd2, 0x58, 0x21, 0x2d, 0xf5, 0x1a, 0xa3, 0xd7, 0xb6, 0x6f, 0xde, 0xcd, 0xc6, 0xe6, 0x51,
	0x6b, 0x4c, 0xdb, 0xc9, 0xd2, 0x2f, 0xd4, 0xe3, 0xb6, 0x74, 0x4a, 0xcb, 0xc1, 0x71, 0xdd, 0x87,
	0xec, 0x66, 0x57, 0x1f, 0x9d, 0xea, 0x86, 0x60, 0x5d, 0xfc, 0x88, 0x21, 0xf3, 0xed, 0xb7, 0x72,
	0xc2, 0xe7, 0x10, 0x2b, 0x41, 0xa2, 0x3c, 0x2a, 0x56, 0x2c, 0x56, 0x02, 0xdf, 0x87, 0xf5, 0x30,
	0xe9, 0xda, 0x74, 0x24, 0xce, 0xa3, 0x02, 0xb1, 0xdd, 0x17, 0xc6, 0x90, 0x0c, 0x4a, 0x48, 0xb2,
	0xf2, 0xae, 0xd7, 0xb3, 0x37, 0x8e, 0x4a, 0x90, 0x24, 0x8f, 0x8a, 0x33, 0xe6, 0x35, 0x2e, 0x20,
	0xed, 0xad, 0x6a, 0x24, 0x49, 0xf3, 0xa8, 0x38, 0x7d, 0x8a, 0xe9, 0x7c, 0xd6, 0x4f, 0xfe, 0xef,
	0xaf, 0x64, 0xa3, 0x34, 0xef, 0x58, 0x08, 0xe0, 0x27, 0x00, 0x83, 0x33, 0x7d, 0x15, 0xe2, 0xeb,
	0xa3, 0x71, 0x34, 0xa7, 0xde, 0xfb, 0x0a, 0x86, 0x64, 0xcb, 0xaf, 0x24, 0x39, 0xc9, 0xa3, 0x22,
	0x63, 0x5e, 0xe3, 0xe7, 0x00, 0x8d, 0x95, 0xdc, 0x49, 0x51, 0x71, 0x47, 0x32, 0xbf, 0xcd, 0x86,
	0x06, 0x2c, 0x74, 0xc1, 0x42, 0x3f, 0x2e, 0x58, 0x18, 0xda, 0xa5, 0x5f, 0xb8, 0x8b, 0x9f, 0x2b,
	0x48, 0x3d, 0x88, 0xbf, 0x28, 0x2c, 0x37, 0x8b, 0x0f, 0x6e, 0xf6, 0x00, 0x90, 0x96, 0xba, 0x96,
	0xb6, 0x52, 0xc2, 0x63, 0x58, 0xb1, 0x2c, 0x18, 0x6f, 0x0e, 0xb1, 0x25, 0xff, 0xc4, 0x96, 0xfe,
	0x8e, 0xcd, 0x4d, 0x7d, 0xb8, 0x32, 0x62, 0x5e, 0xdf, 0x60, 0x3b, 0xb9, 0x1d, 0xb6, 0xec, 0x7f,
	0xb0, 0x51, 0xc8, 0xbe, 0x8e, 0xfc, 0xda, 0x29, 0x37, 0x11, 0x74, 0xb4, 0xb0, 0xcf, 0xe0, 0x4b,
	0xb8, 0xbb, 0x55, 0x5d, 0x27, 0x45, 0xb5, 0xaf, 0xc1, 0xd1, 0xda, 0x79, 0x88, 0x7e, 0x58, 0xca,
	0xcb, 0x8c, 0x4e, 0x0f, 0x66, 0xf4, 0x10, 0x50, 0xc3, 0xaf, 0x1b, 0x39, 0x07, 0xc9, 0x99, 0x5f,
	0xb8, 0x31, 0xfe, 0x98, 0xe0, 0x9d, 0x5b, 0x4c, 0xf0, 0x25, 0xfd, 0xfc, 0xb8, 0x55, 0xee, 0xcb,
	0x58, 0xd3, 0xc6, 0xe8, 0xf2, 0xfb, 0xa0, 0xb9, 0x75, 0xf2, 0x5b, 0xd9, 0x5f, 0xb5, 0xe5, 0x7c,
	0xd2, 0xf0, 0x6a, 0x2e, 0xf7, 0xcf, 0xa5, 0x5e, 0xfb, 0xed, 0x9e, 0xfd, 0x0a, 0x00, 0x00, 0xff,
	0xff, 0xdb, 0xca, 0xf1, 0x7c, 0x55, 0x03, 0x00, 0x00,
}
