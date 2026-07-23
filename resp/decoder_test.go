package resp

import (
	"reflect"
	"strings"
	"testing"
)

func TestDecoder(t *testing.T) {
	tests := []struct {
		name     string
		given    string
		expected Value
	}{
		{
			name:  "Simple String",
			given: "+OK\r\n",
			expected: Value{
				Type: SimpleString,
				Str:  "OK",
			},
		},
		{
			name:  "Simple Error",
			given: "-ERR unknown command\r\n",
			expected: Value{
				Type: SimpleError,
				Str:  "ERR unknown command",
			},
		},
		{
			name:  "Positive Integer",
			given: ":42\r\n",
			expected: Value{
				Type: Integer,
				Num:  42,
			},
		}, {
			name:  "Negative Integer",
			given: ":-111\r\n",
			expected: Value{
				Type: Integer,
				Num:  -111,
			},
		},
		{
			name:  "Bulk String",
			given: "$5\r\nhello\r\n",
			expected: Value{
				Type: BulkString,
				Str:  "hello",
			},
		},
		{
			name:  "Empty Bulk String",
			given: "$0\r\n\r\n",
			expected: Value{
				Type: BulkString,
				Str:  "",
			},
		},
		{
			name:  "Null Bulk String",
			given: "$-1\r\n",
			expected: Value{
				Type:   BulkString,
				IsNull: true,
			},
		}, {
			name:  "Empty Array",
			given: "*0\r\n",
			expected: Value{
				Type:  Array,
				Array: []Value{},
			},
		}, {
			name:  "Null Array",
			given: "*-1\r\n",
			expected: Value{
				Type:   Array,
				IsNull: true,
			},
		}, {
			name:  "Get Command",
			given: "*2\r\n$3\r\nGET\r\n$3\r\nage\r\n",
			expected: Value{
				Type: Array,
				Array: []Value{
					{
						Type: BulkString,
						Str:  "GET",
					},
					{
						Type: BulkString,
						Str:  "age",
					},
				},
			},
		},
		{
			name:  "Nested Array",
			given: "*2\r\n:1\r\n*2\r\n+OK\r\n:2\r\n",
			expected: Value{
				Type: Array,
				Array: []Value{
					{
						Type: Integer,
						Num:  1,
					},
					{
						Type: Array,
						Array: []Value{
							{
								Type: SimpleString,
								Str:  "OK",
							},
							{
								Type: Integer,
								Num:  2,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := NewDecoder(strings.NewReader(tt.given))
			got, err := decoder.Decode()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf(
					"\nExpected: %#v\nGot:      %#v",
					tt.expected,
					got,
				)
			}
		})
	}
}
