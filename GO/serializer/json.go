package serializer

import (
	// proto
	proto "github.com/golang/protobuf/proto"
	// jsonpb
	jsonpb "github.com/golang/protobuf/jsonpb"
)

func ProtobufToJSON(message proto.Message) (string, error) {
	// marshaler
	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		Indent:       " ",
		OrigName:     true,
	}
	// marshal to json
	return marshaler.MarshalToString(message)

}
