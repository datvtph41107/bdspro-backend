package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
)

var errInvalidBinary = errors.New("invalid binary payload")

const maxElements = 4096

func appendU16(dst []byte, v uint16) []byte {
	return binary.LittleEndian.AppendUint16(dst, v)
}

func appendI32(dst []byte, v int) []byte {
	return binary.LittleEndian.AppendUint32(dst, uint32(int32(v)))
}

func appendF64(dst []byte, v float64) []byte {
	return binary.LittleEndian.AppendUint64(dst, math.Float64bits(v))
}

// Binary layout:
// int32 ID, X, Y, HP
// uint16 skillsCount + int32 skills[]
// uint16 buffsCount + int32 buffs[]
// uint16 positionCount + float64 position[]
// uint16 nearbyCount + (uint16 innerCount + int32 inner[])[]

func EncodePlayerBinary(p PlayerArray) []byte {
	// Đủ dung lượng cho dữ liệu mẫu; có thể tính chính xác hơn
	// nếu cần tối ưu encode và tránh realloc.
	dst := make([]byte, 0, 128)

	dst = appendI32(dst, p.ID)
	dst = appendI32(dst, p.X)
	dst = appendI32(dst, p.Y)
	dst = appendI32(dst, p.HP)

	dst = appendU16(dst, uint16(len(p.Skills)))
	for _, v := range p.Skills {
		dst = appendI32(dst, v)
	}

	dst = appendU16(dst, uint16(len(p.Buffs)))
	for _, v := range p.Buffs {
		dst = appendI32(dst, v)
	}

	dst = appendU16(dst, uint16(len(p.Position)))
	for _, v := range p.Position {
		dst = appendF64(dst, v)
	}

	dst = appendU16(dst, uint16(len(p.Nearby)))
	for _, group := range p.Nearby {
		dst = appendU16(dst, uint16(len(group)))
		for _, v := range group {
			dst = appendI32(dst, v)
		}
	}

	return dst
}

type binaryReader struct {
	data []byte
	pos  int
}

func (r *binaryReader) take(n int) ([]byte, error) {
	if n < 0 || n > len(r.data)-r.pos {
		return nil, errInvalidBinary
	}
	v := r.data[r.pos : r.pos+n]
	r.pos += n
	return v, nil
}

func (r *binaryReader) u16() (uint16, error) {
	b, err := r.take(2)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b), nil
}

func (r *binaryReader) i32() (int, error) {
	b, err := r.take(4)
	if err != nil {
		return 0, err
	}
	return int(int32(binary.LittleEndian.Uint32(b))), nil
}

func (r *binaryReader) f64() (float64, error) {
	b, err := r.take(8)
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(binary.LittleEndian.Uint64(b)), nil
}

func (r *binaryReader) count() (int, error) {
	n, err := r.u16()
	if err != nil {
		return 0, err
	}
	if n > maxElements {
		return 0, errInvalidBinary
	}
	return int(n), nil
}

func (r *binaryReader) intSlice() ([]int, error) {
	n, err := r.count()
	if err != nil {
		return nil, err
	}

	// Kiểm tra đủ byte trước khi cấp phát.
	if n > (len(r.data)-r.pos)/4 {
		return nil, errInvalidBinary
	}

	values := make([]int, n)
	for i := range values {
		values[i], err = r.i32()
		if err != nil {
			return nil, err
		}
	}
	return values, nil
}

func DecodePlayerBinary(data []byte) (PlayerArray, error) {
	r := &binaryReader{data: data}
	var p PlayerArray
	var err error

	if p.ID, err = r.i32(); err != nil {
		return p, err
	}
	if p.X, err = r.i32(); err != nil {
		return p, err
	}
	if p.Y, err = r.i32(); err != nil {
		return p, err
	}
	if p.HP, err = r.i32(); err != nil {
		return p, err
	}

	if p.Skills, err = r.intSlice(); err != nil {
		return p, err
	}
	if p.Buffs, err = r.intSlice(); err != nil {
		return p, err
	}

	n, err := r.count()
	if err != nil {
		return p, err
	}
	if n > (len(data)-r.pos)/8 {
		return p, errInvalidBinary
	}

	p.Position = make([]float64, n)
	for i := range p.Position {
		if p.Position[i], err = r.f64(); err != nil {
			return p, err
		}
	}

	n, err = r.count()
	if err != nil {
		return p, err
	}
	if n > (len(data)-r.pos)/2 {
		return p, errInvalidBinary
	}

	p.Nearby = make([][]int, n)
	for i := range p.Nearby {
		if p.Nearby[i], err = r.intSlice(); err != nil {
			return p, err
		}
	}

	if r.pos != len(data) {
		return p, errInvalidBinary
	}

	return p, nil
}

type ArrayParser struct {
	data []byte
	pos  int
}

func NewArrayParser(data []byte) *ArrayParser {
	return &ArrayParser{data: data}
}

func (p *ArrayParser) skipWS() {
	for p.pos < len(p.data) {
		switch p.data[p.pos] {
		case ' ', '\n', '\r', '\t':
			p.pos++
		default:
			return
		}
	}
}

func (p *ArrayParser) expect(c byte) error {
	p.skipWS()

	if p.pos >= len(p.data) || p.data[p.pos] != c {
		return fmt.Errorf("expected %q at byte %d", c, p.pos)
	}

	p.pos++
	return nil
}

func (p *ArrayParser) separator(end bool) error {
	p.skipWS()

	if p.pos >= len(p.data) {
		return fmt.Errorf("unexpected end of input")
	}

	if end {
		if p.data[p.pos] == ']' {
			return nil
		}
		if p.data[p.pos] == ',' {
			p.pos++
			return nil
		}
		return fmt.Errorf("expected ',' or ']' at byte %d", p.pos)
	}

	if p.data[p.pos] != ',' {
		return fmt.Errorf("expected ',' at byte %d", p.pos)
	}
	p.pos++
	return nil
}

func (p *ArrayParser) integer() (int, error) {
	p.skipWS()

	start := p.pos
	if p.pos < len(p.data) && p.data[p.pos] == '-' {
		p.pos++
	}

	digitStart := p.pos
	for p.pos < len(p.data) &&
		p.data[p.pos] >= '0' && p.data[p.pos] <= '9' {
		p.pos++
	}

	if p.pos == digitStart {
		return 0, fmt.Errorf("invalid integer at byte %d", start)
	}

	n, err := strconv.ParseInt(
		string(p.data[start:p.pos]), 10, strconv.IntSize,
	)
	return int(n), err
}

func (p *ArrayParser) decimal() (float64, error) {
	p.skipWS()

	start := p.pos
	for p.pos < len(p.data) {
		c := p.data[p.pos]
		if (c >= '0' && c <= '9') ||
			c == '-' || c == '+' ||
			c == '.' || c == 'e' || c == 'E' {
			p.pos++
		} else {
			break
		}
	}

	if start == p.pos {
		return 0, fmt.Errorf("invalid number at byte %d", start)
	}

	return strconv.ParseFloat(string(p.data[start:p.pos]), 64)
}

func (p *ArrayParser) intArray() ([]int, error) {
	if err := p.expect('['); err != nil {
		return nil, err
	}

	var result []int
	p.skipWS()

	if p.pos < len(p.data) && p.data[p.pos] == ']' {
		p.pos++
		return result, nil
	}

	for {
		v, err := p.integer()
		if err != nil {
			return nil, err
		}
		result = append(result, v)

		p.skipWS()
		if p.pos < len(p.data) && p.data[p.pos] == ']' {
			p.pos++
			return result, nil
		}

		if err := p.expect(','); err != nil {
			return nil, err
		}
	}
}

func (p *ArrayParser) floatArray() ([]float64, error) {
	if err := p.expect('['); err != nil {
		return nil, err
	}

	var result []float64
	p.skipWS()

	if p.pos < len(p.data) && p.data[p.pos] == ']' {
		p.pos++
		return result, nil
	}

	for {
		v, err := p.decimal()
		if err != nil {
			return nil, err
		}
		result = append(result, v)

		p.skipWS()
		if p.pos < len(p.data) && p.data[p.pos] == ']' {
			p.pos++
			return result, nil
		}

		if err := p.expect(','); err != nil {
			return nil, err
		}
	}
}

func (p *ArrayParser) nestedIntArray() ([][]int, error) {
	if err := p.expect('['); err != nil {
		return nil, err
	}

	var result [][]int
	p.skipWS()

	if p.pos < len(p.data) && p.data[p.pos] == ']' {
		p.pos++
		return result, nil
	}

	for {
		values, err := p.intArray()
		if err != nil {
			return nil, err
		}
		result = append(result, values)

		p.skipWS()
		if p.pos < len(p.data) && p.data[p.pos] == ']' {
			p.pos++
			return result, nil
		}

		if err := p.expect(','); err != nil {
			return nil, err
		}
	}
}

// Schema:
// [id, x, y, hp, skills, buffs, position, nearby]
func ParsePlayerArray(data []byte) (PlayerArray, error) {
	p := NewArrayParser(data)
	var result PlayerArray

	if err := p.expect('['); err != nil {
		return result, err
	}

	fields := []func() error{
		func() error {
			v, err := p.integer()
			result.ID = v
			return err
		},
		func() error {
			v, err := p.integer()
			result.X = v
			return err
		},
		func() error {
			v, err := p.integer()
			result.Y = v
			return err
		},
		func() error {
			v, err := p.integer()
			result.HP = v
			return err
		},
		func() error {
			v, err := p.intArray()
			result.Skills = v
			return err
		},
		func() error {
			v, err := p.intArray()
			result.Buffs = v
			return err
		},
		func() error {
			v, err := p.floatArray()
			result.Position = v
			return err
		},
		func() error {
			v, err := p.nestedIntArray()
			result.Nearby = v
			return err
		},
	}

	for i, parseField := range fields {
		if i > 0 {
			if err := p.expect(','); err != nil {
				return result, err
			}
		}
		if err := parseField(); err != nil {
			return result, fmt.Errorf("field %d: %w", i, err)
		}
	}

	if err := p.expect(']'); err != nil {
		return result, err
	}

	p.skipWS()
	if p.pos != len(p.data) {
		return result, fmt.Errorf("unexpected trailing data at byte %d", p.pos)
	}

	return result, nil
}

type Player struct {
	ID       int       `json:"id"`
	X        int       `json:"x"`
	Y        int       `json:"y"`
	HP       int       `json:"hp"`
	Skills   []int     `json:"skills"`
	Buffs    []int     `json:"buffs"`
	Position []float64 `json:"position"`
	Nearby   [][]int   `json:"nearby"`
}

// Array schema:
// [id, x, y, hp, skills, buffs, position, nearby]
type PlayerArray struct {
	ID       int
	X        int
	Y        int
	HP       int
	Skills   []int
	Buffs    []int
	Position []float64
	Nearby   [][]int
}

func (p *PlayerArray) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))

	token, err := dec.Token()
	if err != nil {
		return err
	}
	if token != json.Delim('[') {
		return fmt.Errorf("expected JSON array")
	}

	// Decode each positional field into its actual type.
	fields := []any{
		&p.ID,
		&p.X,
		&p.Y,
		&p.HP,
		&p.Skills,
		&p.Buffs,
		&p.Position,
		&p.Nearby,
	}

	for _, field := range fields {
		if !dec.More() {
			return fmt.Errorf("not enough array fields")
		}
		if err := dec.Decode(field); err != nil {
			return err
		}
	}

	if dec.More() {
		return fmt.Errorf("too many array fields")
	}

	token, err = dec.Token()
	if err != nil {
		return err
	}
	if token != json.Delim(']') {
		return fmt.Errorf("expected end of array")
	}

	return nil
}

var objectJSON = []byte(`{
	"id":123,
	"x":10,
	"y":20,
	"hp":100,
	"skills":[1,2,3,4],
	"buffs":[10,20,30],
	"position":[10.5,20.5,30.5],
	"nearby":[[1,10,20],[2,30,40],[3,50,60]]
}`)

var arrayJSON = []byte(
	`[123,10,20,100,[1,2,3,4],[10,20,30],[10.5,20.5,30.5],[[1,10,20],[2,30,40],[3,50,60]]]`,
)

func main() {
	var object Player
	var array PlayerArray

	if err := json.Unmarshal(objectJSON, &object); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(arrayJSON, &array); err != nil {
		panic(err)
	}

	fmt.Printf("Object: %+v\n", object)
	fmt.Printf("Array:  %+v\n", array)

	fmt.Printf("\nObject payload: %d bytes\n", len(objectJSON))
	fmt.Printf("Array payload:  %d bytes\n", len(arrayJSON))
}
