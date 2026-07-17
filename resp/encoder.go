package main

import (
	"bytes"
	"fmt"
	"strconv"
)

type RespType byte

const (
	SimpleString RespType = '+'
	SimpleError  RespType = '-'
	Integer      RespType = ':'
	BulkString   RespType = '$'
	Array        RespType = '*'
)

type Value struct {
	Type   RespType
	Str    string
	Num    int64
	Array  []Value
	IsNull bool
}

func (v Value) Marshal() []byte {
	switch v.Type {
	case SimpleString:
		return marshalSimpleString(v)
	case SimpleError:
		return marshalSimpleError(v)
	case Integer:
		return marshalInteger(v)
	case BulkString:
		return marshalBulkString(v)
	case Array:
		return marshalArray(v)
	default:
		panic("Invalid type found ")
	}
}

// Simple String -> +<string>\r\n
func marshalSimpleString(v Value) []byte {
	var buf bytes.Buffer

	buf.WriteByte(byte(v.Type))
	buf.WriteString(v.Str)
	buf.WriteString("\r\n")

	return buf.Bytes()
}

// Error         -> -<string>\r\n
func marshalSimpleError(v Value) []byte {
	var buf bytes.Buffer

	buf.WriteByte(byte(v.Type))
	buf.WriteString(v.Str)
	buf.WriteString("\r\n")

	return buf.Bytes()
}

// Integer       -> :<number>\r\n
func marshalInteger(v Value) []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(v.Type))
	buf.WriteString(strconv.FormatInt(v.Num, 10))
	buf.WriteString("\r\n")

	return buf.Bytes()

}

// Bulk String   -> $<len>\r\n<data>\r\n
func marshalBulkString(v Value) []byte {
	var buf bytes.Buffer
	lengthStr := -1
	buf.WriteByte(byte(v.Type))

	if v.IsNull {
		buf.WriteString(strconv.FormatInt(int64(lengthStr), 10))
		buf.WriteString("\r\n")
		return buf.Bytes()
	}
	lengthStr = len(v.Str)
	buf.WriteString(strconv.FormatInt(int64(lengthStr), 10))
	buf.WriteString("\r\n")
	buf.WriteString(v.Str)
	buf.WriteString("\r\n")

	return buf.Bytes()
}

// Array         -> *<count>\r\n<encoded elements...>
func marshalArray(v Value) []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(v.Type))
	if v.IsNull {
		buf.WriteString("-1\r\n")
		return buf.Bytes()
	}
	lengthArr := len(v.Array)
	buf.WriteString(strconv.FormatInt(int64(lengthArr), 10))
	buf.WriteString("\r\n")
	for _, ele := range v.Array {
		buf.Write(ele.Marshal())
	}

	return buf.Bytes()
}

// todo have separate test for encoder
func main() {
	v1 := Value{
		Type: SimpleString,
		Str:  "OK",
	}

	v2 := Value{
		Type: SimpleError,
		Str:  "CTRL",
	}
	v3 := Value{
		Type: Integer,
		Num:  -111111,
	}

	v4 := Value{
		Type: BulkString,
		Str:  "",
	}
	v5 := Value{
		Type:   BulkString,
		IsNull: true,
	}
	v6 := Value{
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
	}

	fmt.Printf("%q\n", v1.Marshal())
	fmt.Printf("%q\n", v2.Marshal())
	fmt.Printf("%q\n", v3.Marshal())
	fmt.Printf("%q\n", v4.Marshal())
	fmt.Printf("%q\n", v5.Marshal())
	fmt.Printf("%q\n", v6.Marshal())
}
