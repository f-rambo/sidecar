// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: api/cloud/v1alpha1/message.proto

package v1alpha1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Cloud struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                       int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	KubeadmInitConfig        string `protobuf:"bytes,2,opt,name=kubeadm_init_config,json=kubeadmInitConfig,proto3" json:"kubeadm_init_config,omitempty"`
	KubeadmConfig            string `protobuf:"bytes,3,opt,name=kubeadm_config,json=kubeadmConfig,proto3" json:"kubeadm_config,omitempty"`
	KubeletService           string `protobuf:"bytes,4,opt,name=kubelet_service,json=kubeletService,proto3" json:"kubelet_service,omitempty"`
	CrioVersion              string `protobuf:"bytes,5,opt,name=crio_version,json=crioVersion,proto3" json:"crio_version,omitempty"`
	Arch                     string `protobuf:"bytes,6,opt,name=arch,proto3" json:"arch,omitempty"`
	Token                    string `protobuf:"bytes,7,opt,name=token,proto3" json:"token,omitempty"`
	DiscoveryTokenCaCertHash string `protobuf:"bytes,8,opt,name=discovery_token_ca_cert_hash,json=discoveryTokenCaCertHash,proto3" json:"discovery_token_ca_cert_hash,omitempty"`
	ControlPlaneEndpoint     string `protobuf:"bytes,9,opt,name=control_plane_endpoint,json=controlPlaneEndpoint,proto3" json:"control_plane_endpoint,omitempty"`
	JoinConfig               string `protobuf:"bytes,10,opt,name=join_config,json=joinConfig,proto3" json:"join_config,omitempty"`
}

func (x *Cloud) Reset() {
	*x = Cloud{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_cloud_v1alpha1_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Cloud) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cloud) ProtoMessage() {}

func (x *Cloud) ProtoReflect() protoreflect.Message {
	mi := &file_api_cloud_v1alpha1_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cloud.ProtoReflect.Descriptor instead.
func (*Cloud) Descriptor() ([]byte, []int) {
	return file_api_cloud_v1alpha1_message_proto_rawDescGZIP(), []int{0}
}

func (x *Cloud) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Cloud) GetKubeadmInitConfig() string {
	if x != nil {
		return x.KubeadmInitConfig
	}
	return ""
}

func (x *Cloud) GetKubeadmConfig() string {
	if x != nil {
		return x.KubeadmConfig
	}
	return ""
}

func (x *Cloud) GetKubeletService() string {
	if x != nil {
		return x.KubeletService
	}
	return ""
}

func (x *Cloud) GetCrioVersion() string {
	if x != nil {
		return x.CrioVersion
	}
	return ""
}

func (x *Cloud) GetArch() string {
	if x != nil {
		return x.Arch
	}
	return ""
}

func (x *Cloud) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *Cloud) GetDiscoveryTokenCaCertHash() string {
	if x != nil {
		return x.DiscoveryTokenCaCertHash
	}
	return ""
}

func (x *Cloud) GetControlPlaneEndpoint() string {
	if x != nil {
		return x.ControlPlaneEndpoint
	}
	return ""
}

func (x *Cloud) GetJoinConfig() string {
	if x != nil {
		return x.JoinConfig
	}
	return ""
}

var File_api_cloud_v1alpha1_message_proto protoreflect.FileDescriptor

var file_api_cloud_v1alpha1_message_proto_rawDesc = []byte{
	0x0a, 0x20, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0x22, 0xfb, 0x02, 0x0a, 0x05, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2e, 0x0a, 0x13,
	0x6b, 0x75, 0x62, 0x65, 0x61, 0x64, 0x6d, 0x5f, 0x69, 0x6e, 0x69, 0x74, 0x5f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x6b, 0x75, 0x62, 0x65, 0x61,
	0x64, 0x6d, 0x49, 0x6e, 0x69, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x25, 0x0a, 0x0e,
	0x6b, 0x75, 0x62, 0x65, 0x61, 0x64, 0x6d, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6b, 0x75, 0x62, 0x65, 0x61, 0x64, 0x6d, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x12, 0x27, 0x0a, 0x0f, 0x6b, 0x75, 0x62, 0x65, 0x6c, 0x65, 0x74, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6b, 0x75,
	0x62, 0x65, 0x6c, 0x65, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x21, 0x0a, 0x0c,
	0x63, 0x72, 0x69, 0x6f, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x63, 0x72, 0x69, 0x6f, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x12, 0x0a, 0x04, 0x61, 0x72, 0x63, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61,
	0x72, 0x63, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x3e, 0x0a, 0x1c, 0x64, 0x69, 0x73,
	0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x63, 0x61, 0x5f,
	0x63, 0x65, 0x72, 0x74, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x18, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x43,
	0x61, 0x43, 0x65, 0x72, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x34, 0x0a, 0x16, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x5f, 0x65, 0x6e, 0x64, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x14, 0x63, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x50, 0x6c, 0x61, 0x6e, 0x65, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12,
	0x1f, 0x0a, 0x0b, 0x6a, 0x6f, 0x69, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6a, 0x6f, 0x69, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x42, 0x1d, 0x5a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x3b, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_cloud_v1alpha1_message_proto_rawDescOnce sync.Once
	file_api_cloud_v1alpha1_message_proto_rawDescData = file_api_cloud_v1alpha1_message_proto_rawDesc
)

func file_api_cloud_v1alpha1_message_proto_rawDescGZIP() []byte {
	file_api_cloud_v1alpha1_message_proto_rawDescOnce.Do(func() {
		file_api_cloud_v1alpha1_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_cloud_v1alpha1_message_proto_rawDescData)
	})
	return file_api_cloud_v1alpha1_message_proto_rawDescData
}

var file_api_cloud_v1alpha1_message_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_cloud_v1alpha1_message_proto_goTypes = []any{
	(*Cloud)(nil), // 0: cloud.v1alpha1.Cloud
}
var file_api_cloud_v1alpha1_message_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_cloud_v1alpha1_message_proto_init() }
func file_api_cloud_v1alpha1_message_proto_init() {
	if File_api_cloud_v1alpha1_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_cloud_v1alpha1_message_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Cloud); i {
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
			RawDescriptor: file_api_cloud_v1alpha1_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_cloud_v1alpha1_message_proto_goTypes,
		DependencyIndexes: file_api_cloud_v1alpha1_message_proto_depIdxs,
		MessageInfos:      file_api_cloud_v1alpha1_message_proto_msgTypes,
	}.Build()
	File_api_cloud_v1alpha1_message_proto = out.File
	file_api_cloud_v1alpha1_message_proto_rawDesc = nil
	file_api_cloud_v1alpha1_message_proto_goTypes = nil
	file_api_cloud_v1alpha1_message_proto_depIdxs = nil
}
