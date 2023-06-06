// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.6.1
// source: photos.proto

package photos

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	protoauth "sara_updated/backend/grpc/proto/protoauth"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// List all available albums for a given account
type AlbumListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageSize  int32                      `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	PageToken string                     `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	OauthInfo *protoauth.OauthConfigInfo `protobuf:"bytes,3,opt,name=oauth_info,json=oauthInfo,proto3" json:"oauth_info,omitempty"`
}

func (x *AlbumListRequest) Reset() {
	*x = AlbumListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_photos_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AlbumListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AlbumListRequest) ProtoMessage() {}

func (x *AlbumListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_photos_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AlbumListRequest.ProtoReflect.Descriptor instead.
func (*AlbumListRequest) Descriptor() ([]byte, []int) {
	return file_photos_proto_rawDescGZIP(), []int{0}
}

func (x *AlbumListRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *AlbumListRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

func (x *AlbumListRequest) GetOauthInfo() *protoauth.OauthConfigInfo {
	if x != nil {
		return x.OauthInfo
	}
	return nil
}

// Gets media from a specific album
type FromAlbumRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AlbumId   string `protobuf:"bytes,1,opt,name=album_id,json=albumId,proto3" json:"album_id,omitempty"`
	PageSize  string `protobuf:"bytes,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	PageToken int32  `protobuf:"varint,3,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
}

func (x *FromAlbumRequest) Reset() {
	*x = FromAlbumRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_photos_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FromAlbumRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FromAlbumRequest) ProtoMessage() {}

func (x *FromAlbumRequest) ProtoReflect() protoreflect.Message {
	mi := &file_photos_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FromAlbumRequest.ProtoReflect.Descriptor instead.
func (*FromAlbumRequest) Descriptor() ([]byte, []int) {
	return file_photos_proto_rawDescGZIP(), []int{1}
}

func (x *FromAlbumRequest) GetAlbumId() string {
	if x != nil {
		return x.AlbumId
	}
	return ""
}

func (x *FromAlbumRequest) GetPageSize() string {
	if x != nil {
		return x.PageSize
	}
	return ""
}

func (x *FromAlbumRequest) GetPageToken() int32 {
	if x != nil {
		return x.PageToken
	}
	return 0
}

type AlbumInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                    string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title                 string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	ProductUrl            string `protobuf:"bytes,3,opt,name=product_url,json=productUrl,proto3" json:"product_url,omitempty"`
	MediaItemsCount       int32  `protobuf:"varint,4,opt,name=media_items_count,json=mediaItemsCount,proto3" json:"media_items_count,omitempty"`
	CoverPhotoBaseUrl     string `protobuf:"bytes,5,opt,name=cover_photo_base_url,json=coverPhotoBaseUrl,proto3" json:"cover_photo_base_url,omitempty"`
	CoverPhotoMediaItemId string `protobuf:"bytes,6,opt,name=cover_photo_media_item_id,json=coverPhotoMediaItemId,proto3" json:"cover_photo_media_item_id,omitempty"`
}

func (x *AlbumInfo) Reset() {
	*x = AlbumInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_photos_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AlbumInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AlbumInfo) ProtoMessage() {}

func (x *AlbumInfo) ProtoReflect() protoreflect.Message {
	mi := &file_photos_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AlbumInfo.ProtoReflect.Descriptor instead.
func (*AlbumInfo) Descriptor() ([]byte, []int) {
	return file_photos_proto_rawDescGZIP(), []int{2}
}

func (x *AlbumInfo) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AlbumInfo) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *AlbumInfo) GetProductUrl() string {
	if x != nil {
		return x.ProductUrl
	}
	return ""
}

func (x *AlbumInfo) GetMediaItemsCount() int32 {
	if x != nil {
		return x.MediaItemsCount
	}
	return 0
}

func (x *AlbumInfo) GetCoverPhotoBaseUrl() string {
	if x != nil {
		return x.CoverPhotoBaseUrl
	}
	return ""
}

func (x *AlbumInfo) GetCoverPhotoMediaItemId() string {
	if x != nil {
		return x.CoverPhotoMediaItemId
	}
	return ""
}

type AlbumsInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Albums        []*AlbumInfo `protobuf:"bytes,1,rep,name=albums,proto3" json:"albums,omitempty"`
	NextPageToken string       `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *AlbumsInfo) Reset() {
	*x = AlbumsInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_photos_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AlbumsInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AlbumsInfo) ProtoMessage() {}

func (x *AlbumsInfo) ProtoReflect() protoreflect.Message {
	mi := &file_photos_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AlbumsInfo.ProtoReflect.Descriptor instead.
func (*AlbumsInfo) Descriptor() ([]byte, []int) {
	return file_photos_proto_rawDescGZIP(), []int{3}
}

func (x *AlbumsInfo) GetAlbums() []*AlbumInfo {
	if x != nil {
		return x.Albums
	}
	return nil
}

func (x *AlbumsInfo) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

type MediaInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MediaItems []*Media `protobuf:"bytes,1,rep,name=media_items,json=mediaItems,proto3" json:"media_items,omitempty"`
	PageToken  string   `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
}

func (x *MediaInfo) Reset() {
	*x = MediaInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_photos_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MediaInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaInfo) ProtoMessage() {}

func (x *MediaInfo) ProtoReflect() protoreflect.Message {
	mi := &file_photos_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaInfo.ProtoReflect.Descriptor instead.
func (*MediaInfo) Descriptor() ([]byte, []int) {
	return file_photos_proto_rawDescGZIP(), []int{4}
}

func (x *MediaInfo) GetMediaItems() []*Media {
	if x != nil {
		return x.MediaItems
	}
	return nil
}

func (x *MediaInfo) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

type Media struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ProductUrl string `protobuf:"bytes,2,opt,name=product_url,json=productUrl,proto3" json:"product_url,omitempty"`
	MimeType   string `protobuf:"bytes,3,opt,name=mime_type,json=mimeType,proto3" json:"mime_type,omitempty"`
}

func (x *Media) Reset() {
	*x = Media{}
	if protoimpl.UnsafeEnabled {
		mi := &file_photos_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Media) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Media) ProtoMessage() {}

func (x *Media) ProtoReflect() protoreflect.Message {
	mi := &file_photos_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Media.ProtoReflect.Descriptor instead.
func (*Media) Descriptor() ([]byte, []int) {
	return file_photos_proto_rawDescGZIP(), []int{5}
}

func (x *Media) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Media) GetProductUrl() string {
	if x != nil {
		return x.ProductUrl
	}
	return ""
}

func (x *Media) GetMimeType() string {
	if x != nil {
		return x.MimeType
	}
	return ""
}

var File_photos_proto protoreflect.FileDescriptor

var file_photos_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b,
	0x6f, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7f, 0x0a, 0x10, 0x41,
	0x6c, 0x62, 0x75, 0x6d, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a,
	0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x2f, 0x0a, 0x0a, 0x6f,
	0x61, 0x75, 0x74, 0x68, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x09, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x69, 0x0a, 0x10,
	0x46, 0x72, 0x6f, 0x6d, 0x41, 0x6c, 0x62, 0x75, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x19, 0x0a, 0x08, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x70,
	0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65,
	0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x70, 0x61,
	0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0xe9, 0x01, 0x0a, 0x09, 0x41, 0x6c, 0x62, 0x75,
	0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x55, 0x72, 0x6c, 0x12, 0x2a, 0x0a, 0x11,
	0x6d, 0x65, 0x64, 0x69, 0x61, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0f, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74,
	0x65, 0x6d, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2f, 0x0a, 0x14, 0x63, 0x6f, 0x76, 0x65,
	0x72, 0x5f, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x5f, 0x62, 0x61, 0x73, 0x65, 0x5f, 0x75, 0x72, 0x6c,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x50, 0x68, 0x6f,
	0x74, 0x6f, 0x42, 0x61, 0x73, 0x65, 0x55, 0x72, 0x6c, 0x12, 0x38, 0x0a, 0x19, 0x63, 0x6f, 0x76,
	0x65, 0x72, 0x5f, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x5f, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x5f, 0x69,
	0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74, 0x65,
	0x6d, 0x49, 0x64, 0x22, 0x58, 0x0a, 0x0a, 0x41, 0x6c, 0x62, 0x75, 0x6d, 0x73, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x22, 0x0a, 0x06, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0a, 0x2e, 0x41, 0x6c, 0x62, 0x75, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x06, 0x61,
	0x6c, 0x62, 0x75, 0x6d, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61,
	0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x6e, 0x65, 0x78, 0x74, 0x50, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x53, 0x0a,
	0x09, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x27, 0x0a, 0x0b, 0x6d, 0x65,
	0x64, 0x69, 0x61, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x06, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x52, 0x0a, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x74,
	0x65, 0x6d, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x22, 0x55, 0x0a, 0x05, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x55, 0x72, 0x6c, 0x12, 0x1b, 0x0a, 0x09,
	0x6d, 0x69, 0x6d, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x6d, 0x69, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x32, 0x72, 0x0a, 0x12, 0x47, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x2c, 0x0a, 0x0a, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x6c, 0x62, 0x75, 0x6d, 0x73, 0x12, 0x11, 0x2e,
	0x41, 0x6c, 0x62, 0x75, 0x6d, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x0b, 0x2e, 0x41, 0x6c, 0x62, 0x75, 0x6d, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x2e, 0x0a,
	0x0d, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x62, 0x75, 0x6d, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x12, 0x11,
	0x2e, 0x46, 0x72, 0x6f, 0x6d, 0x41, 0x6c, 0x62, 0x75, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x0a, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x42, 0x28, 0x5a,
	0x26, 0x73, 0x61, 0x72, 0x61, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x2f, 0x62, 0x61,
	0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_photos_proto_rawDescOnce sync.Once
	file_photos_proto_rawDescData = file_photos_proto_rawDesc
)

func file_photos_proto_rawDescGZIP() []byte {
	file_photos_proto_rawDescOnce.Do(func() {
		file_photos_proto_rawDescData = protoimpl.X.CompressGZIP(file_photos_proto_rawDescData)
	})
	return file_photos_proto_rawDescData
}

var file_photos_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_photos_proto_goTypes = []interface{}{
	(*AlbumListRequest)(nil),          // 0: AlbumListRequest
	(*FromAlbumRequest)(nil),          // 1: FromAlbumRequest
	(*AlbumInfo)(nil),                 // 2: AlbumInfo
	(*AlbumsInfo)(nil),                // 3: AlbumsInfo
	(*MediaInfo)(nil),                 // 4: MediaInfo
	(*Media)(nil),                     // 5: Media
	(*protoauth.OauthConfigInfo)(nil), // 6: OauthConfigInfo
}
var file_photos_proto_depIdxs = []int32{
	6, // 0: AlbumListRequest.oauth_info:type_name -> OauthConfigInfo
	2, // 1: AlbumsInfo.albums:type_name -> AlbumInfo
	5, // 2: MediaInfo.media_items:type_name -> Media
	0, // 3: GooglePhotoService.ListAlbums:input_type -> AlbumListRequest
	1, // 4: GooglePhotoService.GetAlbumMedia:input_type -> FromAlbumRequest
	3, // 5: GooglePhotoService.ListAlbums:output_type -> AlbumsInfo
	4, // 6: GooglePhotoService.GetAlbumMedia:output_type -> MediaInfo
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_photos_proto_init() }
func file_photos_proto_init() {
	if File_photos_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_photos_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AlbumListRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_photos_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FromAlbumRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_photos_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AlbumInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_photos_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AlbumsInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_photos_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MediaInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_photos_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Media); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_photos_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_photos_proto_goTypes,
		DependencyIndexes: file_photos_proto_depIdxs,
		MessageInfos:      file_photos_proto_msgTypes,
	}.Build()
	File_photos_proto = out.File
	file_photos_proto_rawDesc = nil
	file_photos_proto_goTypes = nil
	file_photos_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// GooglePhotoServiceClient is the client API for GooglePhotoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GooglePhotoServiceClient interface {
	ListAlbums(ctx context.Context, in *AlbumListRequest, opts ...grpc.CallOption) (*AlbumsInfo, error)
	GetAlbumMedia(ctx context.Context, in *FromAlbumRequest, opts ...grpc.CallOption) (*MediaInfo, error)
}

type googlePhotoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGooglePhotoServiceClient(cc grpc.ClientConnInterface) GooglePhotoServiceClient {
	return &googlePhotoServiceClient{cc}
}

func (c *googlePhotoServiceClient) ListAlbums(ctx context.Context, in *AlbumListRequest, opts ...grpc.CallOption) (*AlbumsInfo, error) {
	out := new(AlbumsInfo)
	err := c.cc.Invoke(ctx, "/GooglePhotoService/ListAlbums", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *googlePhotoServiceClient) GetAlbumMedia(ctx context.Context, in *FromAlbumRequest, opts ...grpc.CallOption) (*MediaInfo, error) {
	out := new(MediaInfo)
	err := c.cc.Invoke(ctx, "/GooglePhotoService/GetAlbumMedia", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GooglePhotoServiceServer is the server API for GooglePhotoService service.
type GooglePhotoServiceServer interface {
	ListAlbums(context.Context, *AlbumListRequest) (*AlbumsInfo, error)
	GetAlbumMedia(context.Context, *FromAlbumRequest) (*MediaInfo, error)
}

// UnimplementedGooglePhotoServiceServer can be embedded to have forward compatible implementations.
type UnimplementedGooglePhotoServiceServer struct {
}

func (*UnimplementedGooglePhotoServiceServer) ListAlbums(context.Context, *AlbumListRequest) (*AlbumsInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAlbums not implemented")
}
func (*UnimplementedGooglePhotoServiceServer) GetAlbumMedia(context.Context, *FromAlbumRequest) (*MediaInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAlbumMedia not implemented")
}

func RegisterGooglePhotoServiceServer(s *grpc.Server, srv GooglePhotoServiceServer) {
	s.RegisterService(&_GooglePhotoService_serviceDesc, srv)
}

func _GooglePhotoService_ListAlbums_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AlbumListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GooglePhotoServiceServer).ListAlbums(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/GooglePhotoService/ListAlbums",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GooglePhotoServiceServer).ListAlbums(ctx, req.(*AlbumListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GooglePhotoService_GetAlbumMedia_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FromAlbumRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GooglePhotoServiceServer).GetAlbumMedia(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/GooglePhotoService/GetAlbumMedia",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GooglePhotoServiceServer).GetAlbumMedia(ctx, req.(*FromAlbumRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _GooglePhotoService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "GooglePhotoService",
	HandlerType: (*GooglePhotoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListAlbums",
			Handler:    _GooglePhotoService_ListAlbums_Handler,
		},
		{
			MethodName: "GetAlbumMedia",
			Handler:    _GooglePhotoService_GetAlbumMedia_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "photos.proto",
}
