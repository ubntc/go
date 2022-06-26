package diff

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type protoFlatter struct {
	dest map[string]string
}

func newProtoFlatter() *protoFlatter {
	return &protoFlatter{make(map[string]string)}
}

func (pf *protoFlatter) walk(msg proto.Message, prefix string) {
	msg.ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		switch {
		case protoreflect.Repeated == fd.Cardinality():
			// TODO: handle repeated fields
		case protoreflect.MessageKind == fd.Kind():
			msg := v.Message().Interface()
			if !IsWellKnown(fd, msg) {
				pf.walk(msg, prefix+fd.JSONName()+".")
				return true
			}
		}

		key := prefix + fd.JSONName()
		pf.dest[key] = v.String()
		return true
	})
}

func allKeys(flatters ...*protoFlatter) map[string]int {
	keys := make(map[string]int)
	for _, pf := range flatters {
		for k := range pf.dest {
			keys[k]++
		}
	}
	return keys
}

func flatten(msg proto.Message) *protoFlatter {
	pf := newProtoFlatter()
	pf.walk(msg, "")
	return pf
}

func Diff(lhs, rhs proto.Message) map[string]string {
	pfLeft := flatten(lhs)
	pfRight := flatten(rhs)
	keys := allKeys(pfLeft, pfRight)
	if len(keys) == 0 {
		return nil
	}

	changes := make(map[string]string)
	for k := range keys {
		l, hasLeft := pfLeft.dest[k]
		r, hasRight := pfRight.dest[k]
		switch {
		case hasLeft && !hasRight:
			changes[k] = "deleted" // store old value ???
		case !hasLeft && hasRight:
			changes[k] = "created" // store new value
		case !hasLeft && !hasRight:
			// cannot happen
		case l != r:
			changes[k] = "updated" // store diff
		default:
			// no change
		}
	}
	return changes
}
