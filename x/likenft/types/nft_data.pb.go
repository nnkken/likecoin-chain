// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: likenft/nft_data.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

type NFTData struct {
	Metadata          JsonInput `protobuf:"bytes,1,opt,name=metadata,proto3,customtype=JsonInput" json:"metadata"`
	IscnIdPrefix      string    `protobuf:"bytes,2,opt,name=iscnIdPrefix,proto3" json:"iscnIdPrefix,omitempty"`
	IscnVersionAtMint uint64    `protobuf:"varint,3,opt,name=iscnVersionAtMint,proto3" json:"iscnVersionAtMint,omitempty"`
}

func (m *NFTData) Reset()         { *m = NFTData{} }
func (m *NFTData) String() string { return proto.CompactTextString(m) }
func (*NFTData) ProtoMessage()    {}
func (*NFTData) Descriptor() ([]byte, []int) {
	return fileDescriptor_e4b9b3f4f01f6b92, []int{0}
}
func (m *NFTData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NFTData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NFTData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *NFTData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NFTData.Merge(m, src)
}
func (m *NFTData) XXX_Size() int {
	return m.Size()
}
func (m *NFTData) XXX_DiscardUnknown() {
	xxx_messageInfo_NFTData.DiscardUnknown(m)
}

var xxx_messageInfo_NFTData proto.InternalMessageInfo

func (m *NFTData) GetIscnIdPrefix() string {
	if m != nil {
		return m.IscnIdPrefix
	}
	return ""
}

func (m *NFTData) GetIscnVersionAtMint() uint64 {
	if m != nil {
		return m.IscnVersionAtMint
	}
	return 0
}

func init() {
	proto.RegisterType((*NFTData)(nil), "likecoin.likechain.likenft.NFTData")
}

func init() { proto.RegisterFile("likenft/nft_data.proto", fileDescriptor_e4b9b3f4f01f6b92) }

var fileDescriptor_e4b9b3f4f01f6b92 = []byte{
	// 243 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xcb, 0xc9, 0xcc, 0x4e,
	0xcd, 0x4b, 0x2b, 0xd1, 0xcf, 0x4b, 0x2b, 0x89, 0x4f, 0x49, 0x2c, 0x49, 0xd4, 0x2b, 0x28, 0xca,
	0x2f, 0xc9, 0x17, 0x92, 0x02, 0x89, 0x27, 0xe7, 0x67, 0xe6, 0xe9, 0x81, 0x19, 0x19, 0x89, 0x50,
	0x56, 0x5e, 0x5a, 0x89, 0x94, 0x48, 0x7a, 0x7e, 0x7a, 0x3e, 0x58, 0x99, 0x3e, 0x88, 0x05, 0xd1,
	0xa1, 0xd4, 0xc5, 0xc8, 0xc5, 0xee, 0xe7, 0x16, 0xe2, 0x92, 0x58, 0x92, 0x28, 0xa4, 0xcb, 0xc5,
	0x91, 0x9b, 0x5a, 0x92, 0x08, 0x32, 0x4f, 0x82, 0x51, 0x81, 0x51, 0x83, 0xc7, 0x49, 0xf0, 0xc4,
	0x3d, 0x79, 0x86, 0x5b, 0xf7, 0xe4, 0x39, 0xbd, 0x8a, 0xf3, 0xf3, 0x3c, 0xf3, 0x0a, 0x4a, 0x4b,
	0x82, 0xe0, 0x4a, 0x84, 0x94, 0xb8, 0x78, 0x32, 0x8b, 0x93, 0xf3, 0x3c, 0x53, 0x02, 0x8a, 0x52,
	0xd3, 0x32, 0x2b, 0x24, 0x98, 0x14, 0x18, 0x35, 0x38, 0x83, 0x50, 0xc4, 0x84, 0x74, 0xb8, 0x04,
	0x41, 0xfc, 0xb0, 0xd4, 0xa2, 0xe2, 0xcc, 0xfc, 0x3c, 0xc7, 0x12, 0xdf, 0xcc, 0xbc, 0x12, 0x09,
	0x66, 0x05, 0x46, 0x0d, 0x96, 0x20, 0x4c, 0x09, 0x27, 0xf7, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c,
	0x92, 0x63, 0x7c, 0xf0, 0x48, 0x8e, 0x71, 0xc2, 0x63, 0x39, 0x86, 0x0b, 0x8f, 0xe5, 0x18, 0x6e,
	0x3c, 0x96, 0x63, 0x88, 0xd2, 0x4d, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5,
	0x87, 0xf9, 0x51, 0x1f, 0xee, 0x47, 0xfd, 0x0a, 0x7d, 0x58, 0x80, 0x94, 0x54, 0x16, 0xa4, 0x16,
	0x27, 0xb1, 0x81, 0x3d, 0x67, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xa7, 0x63, 0x03, 0x56, 0x28,
	0x01, 0x00, 0x00,
}

func (m *NFTData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NFTData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *NFTData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.IscnVersionAtMint != 0 {
		i = encodeVarintNftData(dAtA, i, uint64(m.IscnVersionAtMint))
		i--
		dAtA[i] = 0x18
	}
	if len(m.IscnIdPrefix) > 0 {
		i -= len(m.IscnIdPrefix)
		copy(dAtA[i:], m.IscnIdPrefix)
		i = encodeVarintNftData(dAtA, i, uint64(len(m.IscnIdPrefix)))
		i--
		dAtA[i] = 0x12
	}
	{
		size := m.Metadata.Size()
		i -= size
		if _, err := m.Metadata.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintNftData(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintNftData(dAtA []byte, offset int, v uint64) int {
	offset -= sovNftData(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *NFTData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Metadata.Size()
	n += 1 + l + sovNftData(uint64(l))
	l = len(m.IscnIdPrefix)
	if l > 0 {
		n += 1 + l + sovNftData(uint64(l))
	}
	if m.IscnVersionAtMint != 0 {
		n += 1 + sovNftData(uint64(m.IscnVersionAtMint))
	}
	return n
}

func sovNftData(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozNftData(x uint64) (n int) {
	return sovNftData(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *NFTData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNftData
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
			return fmt.Errorf("proto: NFTData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NFTData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Metadata", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNftData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthNftData
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthNftData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Metadata.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IscnIdPrefix", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNftData
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
				return ErrInvalidLengthNftData
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthNftData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IscnIdPrefix = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IscnVersionAtMint", wireType)
			}
			m.IscnVersionAtMint = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNftData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.IscnVersionAtMint |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipNftData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthNftData
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
func skipNftData(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowNftData
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
					return 0, ErrIntOverflowNftData
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
					return 0, ErrIntOverflowNftData
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
				return 0, ErrInvalidLengthNftData
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupNftData
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthNftData
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthNftData        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowNftData          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupNftData = fmt.Errorf("proto: unexpected end of group")
)
