// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ottlfuncs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottltest"
)

func Test_concat(t *testing.T) {
	tests := []struct {
		name      string
		delimiter string
		vals      []ottl.StandardGetSetter
		expected  string
	}{
		{
			name:      "concat strings",
			delimiter: " ",
			vals: []ottl.StandardGetSetter{
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return "hello"
					},
				},
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return "world"
					},
				},
			},
			expected: "hello world",
		},
		{
			name:      "nil",
			delimiter: "",
			vals: []ottl.StandardGetSetter{
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return "hello"
					},
				},
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return nil
					},
				},
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return "world"
					},
				},
			},
			expected: "hello<nil>world",
		},
		{
			name:      "integers",
			delimiter: "",
			vals: []ottl.StandardGetSetter{
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return "hello"
					},
				},
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return int64(1)
					},
				},
			},
			expected: "hello1",
		},
		{
			name:      "floats",
			delimiter: "",
			vals: []ottl.StandardGetSetter{
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return "hello"
					},
				},
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return 3.14159
					},
				},
			},
			expected: "hello3.14159",
		},
		{
			name:      "booleans",
			delimiter: " ",
			vals: []ottl.StandardGetSetter{
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return "hello"
					},
				},
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return true
					},
				},
			},
			expected: "hello true",
		},
		{
			name:      "byte slices",
			delimiter: "",
			vals: []ottl.StandardGetSetter{
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0e, 0xd2, 0xe6, 0x3c, 0xbe, 0x71, 0xf5, 0xa8}
					},
				},
			},
			expected: "00000000000000000ed2e63cbe71f5a8",
		},
		{
			name:      "non-byte slices",
			delimiter: "",
			vals: []ottl.StandardGetSetter{
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
					},
				},
			},
			expected: "",
		},
		{
			name:      "maps",
			delimiter: "",
			vals: []ottl.StandardGetSetter{
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return map[string]string{"key": "value"}
					},
				},
			},
			expected: "",
		},
		{
			name:      "unprintable value in the middle",
			delimiter: "-",
			vals: []ottl.StandardGetSetter{
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return "hello"
					},
				},
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return map[string]string{"key": "value"}
					},
				},
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return "world"
					},
				},
			},
			expected: "hello--world",
		},
		{
			name:      "empty string values",
			delimiter: "__",
			vals: []ottl.StandardGetSetter{
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return ""
					},
				},
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return ""
					},
				},
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return ""
					},
				},
			},
			expected: "____",
		},
		{
			name:      "single argument",
			delimiter: "-",
			vals: []ottl.StandardGetSetter{
				{
					Getter: func(ctx ottl.TransformContext) interface{} {
						return "hello"
					},
				},
			},
			expected: "hello",
		},
		{
			name:      "no arguments",
			delimiter: "-",
			vals:      []ottl.StandardGetSetter{},
			expected:  "",
		},
		{
			name:      "no arguments with an empty delimiter",
			delimiter: "",
			vals:      []ottl.StandardGetSetter{},
			expected:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ottltest.TestTransformContext{}

			getters := make([]ottl.Getter, len(tt.vals))

			for i, val := range tt.vals {
				getters[i] = val
			}

			exprFunc, _ := Concat(tt.delimiter, getters)
			actual := exprFunc(ctx)

			assert.Equal(t, tt.expected, actual)
		})
	}
}
