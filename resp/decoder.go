package resp

import (
	"bufio"

	"fmt"
	"io"
	"strconv"
)

type Decoder struct {
	reader *bufio.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		reader: bufio.NewReader(r),
	}
}

func (d *Decoder) Decode() (Value, error) {
	prefix, err := d.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch RespType(prefix) {
	case SimpleString:
		return d.parseSimpleString()
	case SimpleError:
		return d.parseSimpleError()
	case Integer:
		return d.parseInteger()
	case BulkString:
		return d.parseBulkString()
	case Array:
		return d.parseArray()
	default:
		return Value{}, fmt.Errorf("unknown RESP type: %q", prefix)
	}
}

func (d *Decoder) readLine() (string, error) {
	str, err := d.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return str[:len(str)-2], err
}

func (d *Decoder) parseSimpleString() (Value, error) {
	v := Value{}
	v.Type = SimpleString
	str, err := d.readLine()
	if err != nil {
		return v, err
	}
	v.Str = str
	return v, err

}

func (d *Decoder) parseSimpleError() (Value, error) {
	v := Value{}
	v.Type = SimpleError
	str, err := d.readLine()
	if err != nil {
		return v, err
	}
	v.Str = str
	return v, err

}

func (d *Decoder) parseInteger() (Value, error) {
	v := Value{}
	v.Type = Integer
	str, err := d.readLine()
	if err != nil {
		return v, err
	}
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return v, err
	}
	v.Num = num
	return v, err
}

func (d *Decoder) parseBulkString() (Value, error) {
	v := Value{}
	v.Type = BulkString
	str, err := d.readLine()
	if err != nil {
		return v, err
	}
	num, err := strconv.ParseInt(str, 10, 64)

	if err != nil {
		return v, err
	}
	if num == -1 {
		v.IsNull = true
		return v, nil
	}
	temp := make([]byte, num)
	_, err = io.ReadFull(d.reader, temp)
	if err != nil {
		return Value{}, err
	}
	cr, err := d.reader.ReadByte() //\r
	if err != nil {
		return Value{}, err
	}
	lf, err := d.reader.ReadByte() //\n
	if err != nil {
		return Value{}, err
	}
	if cr != '\r' || lf != '\n' {
		return Value{}, fmt.Errorf("invalid bulk string terminator")
	}

	v.Str = (string)(temp)
	return v, nil

}

func (d *Decoder) parseArray() (Value, error) {
	v := Value{}
	v.Type = Array
	str, err := d.readLine()
	if err != nil {
		return v, err
	}
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return v, err
	}
	if num == -1 {
		v.IsNull = true
		return v, nil
	}
	if num == 0 {
		v.Array = []Value{}
		return v, nil
	}
	for i := int64(0); i < num; i++ {
		val, err := d.Decode()
		if err != nil {
			return v, err
		}
		v.Array = append(v.Array, val)

	}

	return v, nil

}
