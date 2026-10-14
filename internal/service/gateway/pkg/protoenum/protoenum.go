package protoenum

import (
	"reflect"
	"strings"
)

type ProtoEnum interface {
	~int32
	String() string
}

func toUpperSnakeCase(s string) string {
	var res strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			res.WriteByte('_')
		}
		res.WriteRune(r)
	}
	return strings.ToUpper(res.String())
}

func ToCleanString[T ProtoEnum](enumValue T) string {
	str := enumValue.String()

	if (len(str) > 0 && str[0] >= '0' && str[0] <= '9') || strings.HasSuffix(str, "_UNSPECIFIED") {
		return "unspecified"
	}

	typeName := reflect.TypeOf(enumValue).Name()
	prefix := toUpperSnakeCase(typeName) + "_"

	return strings.ToLower(strings.TrimPrefix(str, prefix))
}
