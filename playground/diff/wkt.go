package diff

import (
	"regexp"

	"google.golang.org/protobuf/reflect/protoreflect"
)

var wktExp = regexp.MustCompile(`^google\.protobuf\.`)

func IsWellKnown(fd protoreflect.FieldDescriptor, msg protoreflect.ProtoMessage) bool {
	isProtobufType := wktExp.Match([]byte(fd.Message().FullName()))
	return isProtobufType

	// TODO: Do we need a more explicit check that would require importing
	//       the actual *pb packages?
	//
	// switch msg.(type) {
	// case *structpb.Struct, *timestamppb.Timestamp, *durationpb.Duration:
	//     return true
	// default:
	// 	   return false
	// }
}
