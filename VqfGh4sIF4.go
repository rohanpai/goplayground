package ttype

import (
	"testing"
	"math"
	"github.com/siddontang/mixer/hack"
)

const (
	TpInt64 int = iota + 1
	TpUint64
	TpFloat64
	TpString
	TpBytes
)

// Datum has struct size of 48 bytes.
type Datum struct {
	tp   int // datum type.
	i    int64 // i can hold int64 uint64 float64 values.
	b    []byte // b can hold string and []byte values
	ptr  interface{} // ptr hold all other types.
}

func (d *Datum) Type() int {
	return d.tp
}

func (d *Datum) GetInt64() int64 {
	return d.i
}

func (d *Datum) SetInt64(i int64) {
	d.tp = TpInt64
	d.i = i
}

func (d *Datum) GetUint64() int64 {
	return d.i
}

func (d *Datum) SetUint64(i int64) {
	d.tp = TpInt64
	d.i = i
}

func (d *Datum) GetFloat64() float64 {
	return math.Float64frombits(uint64(d.i))
}

func (d *Datum) SetFloat64(f float64) {
	d.tp = TpFloat64
	d.i = int64(math.Float64bits(f))
}

func (d *Datum) GetString() string {
	return hack.String(d.b)
}

func (d *Datum) SetString(s string) {
	d.tp = TpString
	d.b = hack.Slice(s)
}

func (d *Datum) GetBytes() []byte {
	return d.b
}

func (d *Datum) SetBytes(b []byte) {
	d.tp = TpBytes
	d.b = b
}

func BenchmarkString(b *testing.B) {
	var v Datum
	for i := 0; i < b.N; i++ {
		v.SetString("abc")
		if v.GetString() != "abc" {
			b.Fatal(v.GetInt64())
		}
	}
}

func BenchmarkInt64(b *testing.B) {
	var v Datum
	for i := 0; i < b.N; i++ {
		v.SetInt64(int64(i))
		if v.GetInt64() != int64(i) {
			b.Fatal(v.GetInt64())
		}
	}
}

func BenchmarkAppendDatum(b *testing.B) {
	array := make([]Datum, 0, 100)
	var value Datum
	for i := 0; i < b.N/100; i++ {
		for j := 0; j < 100; j++ {
			array = append(array, value)
		}
		array = array[0:0]
	}
}
