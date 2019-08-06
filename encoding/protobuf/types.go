// Copyright 2019 CUE Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protobuf

import (
	"text/scanner"

	"github.com/emicklei/proto"
)

func protoToCUE(typ string, options []*proto.Option) (ref string, ok bool) {
	t, ok := scalars[typ]
	return t, ok
}

var scalars = map[string]string{
	// Differing
	"sint32":   "int32",
	"sint64":   "int64",
	"fixed32":  "uint32",
	"fixed64":  "uint64",
	"sfixed32": "int32",
	"sfixed64": "int64",

	// Identical to CUE
	"int32":  "int32",
	"int64":  "int64",
	"uint32": "uint32",
	"uint64": "uint64",

	"double": "float64",
	"float":  "float32",

	"bool":   "bool",
	"string": "string",
	"bytes":  "bytes",
}

func (p *protoConverter) setBuiltin(from, to string, pkg *protoConverter) {
	p.scope[0][from] = mapping{to, "", pkg}
}

var (
	pkgTime   = &protoConverter{goPkgPath: "time"}
	pkgStruct = &protoConverter{goPkgPath: "struct"}
)

func (p *protoConverter) mapBuiltinPackage(pos scanner.Position, file string, required bool) (generate bool) {
	// Map some builtin types to their JSON/CUE mappings.
	switch file {
	case "gogoproto/gogo.proto":

	case "google/protobuf/struct.proto":
		p.setBuiltin("google.protobuf.Struct", "{}", nil)
		p.setBuiltin("google.protobuf.Value", "_", nil)
		p.setBuiltin("google.protobuf.NullValue", "null", nil)
		p.setBuiltin("google.protobuf.ListValue", "[...]", nil)
		p.setBuiltin("google.protobuf.StringValue", "string", nil)
		p.setBuiltin("google.protobuf.BoolValue", "bool", nil)
		p.setBuiltin("google.protobuf.NumberValue", "number", nil)
		return false

	case "google/protobuf/empty.proto":
		p.setBuiltin("google.protobuf.Empty", "struct.MaxFields(0)", pkgStruct)
		return false

	case "google/protobuf/duration.proto":
		p.setBuiltin("google.protobuf.Duration", "time.Duration", pkgTime)
		return false

	case "google/protobuf/timestamp.proto":
		p.setBuiltin("google.protobuf.Timestamp", "time.Time", pkgTime)
		return false

	// case "google/protobuf/field_mask.proto":
	// 	p.setBuiltin("google.protobuf.FieldMask", "protobuf.FieldMask", nil)

	// 	protobuf.Any

	default:
		if required {
			failf(pos, "import %q not found", file)
		}
	}
	return true
}
