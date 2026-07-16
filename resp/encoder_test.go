package main

import (
	"bytes"
	"os"
	"testing"
)

func TestValueMarshal_SimpleString(t *testing.T) {
	tests := []struct {
		name string
		in   Value
		want string
	}{
		{"basic OK", Value{Type: SimpleString, Str: "OK"}, "+OK\r\n"},
		{"empty string", Value{Type: SimpleString, Str: ""}, "+\r\n"},
		{"with spaces", Value{Type: SimpleString, Str: "hello world"}, "+hello world\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.Marshal()
			if string(got) != tt.want {
				t.Errorf("Marshal() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValueMarshal_SimpleError(t *testing.T) {
	tests := []struct {
		name string
		in   Value
		want string
	}{
		{"basic error", Value{Type: SimpleError, Str: "CTRL"}, "-CTRL\r\n"},
		{"empty error", Value{Type: SimpleError, Str: ""}, "-\r\n"},
		{"error with message", Value{Type: SimpleError, Str: "ERR unknown command"}, "-ERR unknown command\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.Marshal()
			if string(got) != tt.want {
				t.Errorf("Marshal() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValueMarshal_Integer(t *testing.T) {
	tests := []struct {
		name string
		in   Value
		want string
	}{
		{"positive", Value{Type: Integer, Num: 1000}, ":1000\r\n"},
		{"negative", Value{Type: Integer, Num: -111111}, ":-111111\r\n"},
		{"zero", Value{Type: Integer, Num: 0}, ":0\r\n"},
		{"max int64", Value{Type: Integer, Num: 9223372036854775807}, ":9223372036854775807\r\n"},
		{"min int64", Value{Type: Integer, Num: -9223372036854775808}, ":-9223372036854775808\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.Marshal()
			if string(got) != tt.want {
				t.Errorf("Marshal() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValueMarshal_BulkString(t *testing.T) {
	tests := []struct {
		name string
		in   Value
		want string
	}{
		{"non-empty string", Value{Type: BulkString, Str: "hello"}, "$5\r\nhello\r\n"},
		{"empty but not null", Value{Type: BulkString, Str: ""}, "$0\r\n\r\n"},
		{"null bulk string", Value{Type: BulkString, IsNull: true}, "$-1\r\n"},
		{
			"null takes precedence over Str contents",
			Value{Type: BulkString, Str: "ignored", IsNull: true},
			"$-1\r\n",
		},
		{"string with CRLF-like content", Value{Type: BulkString, Str: "a\r\nb"}, "$4\r\na\r\nb\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.Marshal()
			if string(got) != tt.want {
				t.Errorf("Marshal() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValueMarshal_Array(t *testing.T) {
	tests := []struct {
		name string
		in   Value
		want string
	}{
		{
			"empty array",
			Value{Type: Array, Array: []Value{}},
			"*0\r\n",
		},
		{
			"nil array behaves like empty",
			Value{Type: Array},
			"*0\r\n",
		},
		{
			"array of bulk strings (GET age)",
			Value{
				Type: Array,
				Array: []Value{
					{Type: BulkString, Str: "GET"},
					{Type: BulkString, Str: "age"},
				},
			},
			"*2\r\n$3\r\nGET\r\n$3\r\nage\r\n",
		},
		{
			"array with mixed types",
			Value{
				Type: Array,
				Array: []Value{
					{Type: SimpleString, Str: "OK"},
					{Type: Integer, Num: 42},
					{Type: BulkString, IsNull: true},
				},
			},
			"*3\r\n+OK\r\n:42\r\n$-1\r\n",
		},
		{
			"nested array",
			Value{
				Type: Array,
				Array: []Value{
					{
						Type: Array,
						Array: []Value{
							{Type: Integer, Num: 1},
							{Type: Integer, Num: 2},
						},
					},
					{Type: BulkString, Str: "x"},
				},
			},
			"*2\r\n*2\r\n:1\r\n:2\r\n$1\r\nx\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.Marshal()
			if string(got) != tt.want {
				t.Errorf("Marshal() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValueMarshal_InvalidTypePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected Marshal() to panic for an invalid RespType, but it did not")
		}
	}()

	v := Value{Type: RespType('!')}
	v.Marshal()
}

func TestValueMarshal_DispatchesOnType(t *testing.T) {
	// Ensures each RespType routes to the correct marshal* helper by
	// checking the leading byte of the output matches the type marker.
	tests := []struct {
		name       string
		in         Value
		wantPrefix byte
	}{
		{"simple string prefix", Value{Type: SimpleString, Str: "a"}, '+'},
		{"simple error prefix", Value{Type: SimpleError, Str: "a"}, '-'},
		{"integer prefix", Value{Type: Integer, Num: 1}, ':'},
		{"bulk string prefix", Value{Type: BulkString, Str: "a"}, '$'},
		{"array prefix", Value{Type: Array, Array: []Value{}}, '*'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.Marshal()
			if len(got) == 0 || got[0] != tt.wantPrefix {
				t.Errorf("Marshal() first byte = %q, want %q", got, tt.wantPrefix)
			}
		})
	}
}

// TestMain_DoesNotPanic exercises the package's demo main() function,
// which builds and marshals a variety of Value instances, to guard
// against regressions that would cause it to panic or otherwise crash.
func TestMain_DoesNotPanic(t *testing.T) {
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = origStdout
	}()

	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("main() panicked: %v", rec)
		}
	}()

	main()

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !bytes.Contains([]byte(output), []byte(`"+OK\r\n"`)) {
		t.Errorf("expected main() output to contain marshaled SimpleString, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte(`"*2\r\n$3\r\nGET\r\n$3\r\nage\r\n"`)) {
		t.Errorf("expected main() output to contain marshaled Array, got: %s", output)
	}
}