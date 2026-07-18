package resp

import "testing"

func TestMarshal(t *testing.T) {

	tests := []struct {
		name     string
		value    Value
		expected string
	}{
		{
			name: "Simple String",
			value: Value{
				Type: SimpleString,
				Str:  "OK",
			},
			expected: "+OK\r\n",
		},
		{
			name: "Simple Error",
			value: Value{
				Type: SimpleError,
				Str:  "ERR unknown command",
			},
			expected: "-ERR unknown command\r\n",
		},
		{
			name: "Positive Integer",
			value: Value{
				Type: Integer,
				Num:  42,
			},
			expected: ":42\r\n",
		},
		{
			name: "Negative Integer",
			value: Value{
				Type: Integer,
				Num:  -111,
			},
			expected: ":-111\r\n",
		},
		{
			name: "Empty Bulk String",
			value: Value{
				Type: BulkString,
				Str:  "",
			},
			expected: "$0\r\n\r\n",
		},
		{
			name: "Bulk String",
			value: Value{
				Type: BulkString,
				Str:  "hello",
			},
			expected: "$5\r\nhello\r\n",
		},
		{
			name: "Null Bulk String",
			value: Value{
				Type:   BulkString,
				IsNull: true,
			},
			expected: "$-1\r\n",
		},
		{
			name: "Empty Array",
			value: Value{
				Type:  Array,
				Array: []Value{},
			},
			expected: "*0\r\n",
		},
		{
			name: "Null Array",
			value: Value{
				Type:   Array,
				IsNull: true,
			},
			expected: "*-1\r\n",
		},
		{
			name: "GET Command",
			value: Value{
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
			expected: "*2\r\n$3\r\nGET\r\n$3\r\nage\r\n",
		},
		{
			name: "Nested Array",
			value: Value{
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
								Type: BulkString,
								Str:  "PING",
							},
						},
					},
				},
			},
			expected: "*2\r\n:1\r\n*1\r\n$4\r\nPING\r\n",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			got := string(tt.value.Marshal())

			if got != tt.expected {
				t.Errorf(
					"\nExpected: %q\nGot:      %q",
					tt.expected,
					got,
				)
			}
		})
	}
}
