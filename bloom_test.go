package bloom

import (
	"encoding/binary"
//	"fmt"
	"testing"
    "io"
)

func TestBasic(t *testing.T) {
	f := New(1000, 4)
	n1 := []byte("Bess")
	n2 := []byte("Jane")
	f.Add(n1)
	n1b := f.Test(n1)
	n2b := f.Test(n2)
	if !n1b {
		t.Errorf("%v should be in.", n1)
	}
	if n2b {
		t.Errorf("%v should not be in.", n2)
	}
}

func TestBasicUint32(t *testing.T) {
	f := New(1000, 4)
	n1 := make([]byte, 4)
	n2 := make([]byte, 4)
	n3 := make([]byte, 4)
	binary.BigEndian.PutUint32(n1, 100)
	binary.BigEndian.PutUint32(n2, 101)
	binary.BigEndian.PutUint32(n3, 102)
	f.Add(n1)
	n1b := f.Test(n1)
	n2b := f.Test(n2)
	f.Test(n3)
	if !n1b {
		t.Errorf("%v should be in.", n1)
	}
	if n2b {
		t.Errorf("%v should not be in.", n2)
	}
}

func TestDirect20_5(t *testing.T) {
	n := uint(10000)
	k := uint(5)
	load := uint(20)
	f := New(n*load, k)
	fp_rate := f.EstimateFalsePositiveRate(n)
	if fp_rate > 0.0001 {
		t.Errorf("False positive rate too high: load=%v, k=%v, %f", load, k, fp_rate)
	}
}

func TestDirect15_10(t *testing.T) {
	n := uint(10000)
	k := uint(10)
	load := uint(15)
	f := New(n*load, k)
	fp_rate := f.EstimateFalsePositiveRate(n)
	if fp_rate > 0.0001 {
		t.Errorf("False positive rate too high: load=%v, k=%v, %f", load, k, fp_rate)
	}
}

func TestEstimated10_0001(t *testing.T) {
	n := uint(10000)
	fp := 0.0001
	m, k := EstimateParameters(n, fp)
	f := NewWithEstimates(n, fp)
	fp_rate := f.EstimateFalsePositiveRate(n)
	if fp_rate > fp {
		t.Errorf("False positive rate too high: n: %v, fp: %f, n: %v, k: %v result: %f", n, fp, m, k, fp_rate)
	}
}

func TestEstimated10_001(t *testing.T) {
	n := uint(10000)
	fp := 0.001
	m, k := EstimateParameters(n, fp)
	f := NewWithEstimates(n, fp)
	fp_rate := f.EstimateFalsePositiveRate(n)
	if fp_rate > fp {
		t.Errorf("False positive rate too high: n: %v, fp: %f, n: %v, k: %v result: %f", n, fp, m, k, fp_rate)
	}
}

type rw struct{
    buf []byte
    r int
}

func (r *rw) Write(b []byte) (int, error) {
    r.buf = append(r.buf, b...)

    return len(b), nil
}

func (r *rw) Read(b []byte) (int, error) {
    n := copy(b, r.buf[r.r:])
    r.r += n
    if n < len(b) {//eof
        return n, io.EOF
    }
    return n, nil
}

func TestDumpRestore(t *testing.T) {
	a := NewWithEstimates(20000, 0.01)
	addValues := [][]byte{
		[]byte("ala"),
		[]byte("ma"),
		[]byte("kota"),
		[]byte("a"),
		[]byte("kot"),
		[]byte("nie")}
	for _, v := range addValues {
		a.Add(v)
	}
    wr := &rw{make([]byte, 0, 10), 0}
	Encode(wr, a)
	b := Decode(wr)
	for _, v := range addValues {
		if !b.Test(v) {//no false negatives!
			t.Error("Did not restore properly")
		}
	}
}

func BenchmarkAdd(b *testing.B){
	b.StopTimer()
	//k, m := EstimateParameters(10000,0.01)
	n := 10000000
	f := NewWithEstimates(uint(n), 0.001)
	n1 := make([]byte, 4)
	b.StartTimer()
	for i := 0; i < b.N ; i++ {
		binary.BigEndian.PutUint32(n1, uint32(i % n))
		f.Add(n1)
	}
}

func BenchmarkNegativeTest(b *testing.B){
	b.StopTimer()
	//k, m := EstimateParameters(10000,0.01)
	n := 10000000
	f := NewWithEstimates(uint(n), 0.001)
	n1 := make([]byte, 4)
	b.StartTimer()
	for i := 0; i < b.N ; i++ {
		binary.BigEndian.PutUint32(n1, uint32(i % n))
		f.Test(n1)
	}
}

func BenchmarkPositiveTest(b *testing.B){
	b.StopTimer()
	//k, m := EstimateParameters(10000,0.01)
	n := 10000000
	f := NewWithEstimates(uint(n), 0.001)
	n1 := make([]byte, 4)
	binary.BigEndian.PutUint32(n1, uint32(1))
	f.Add(n1)
	b.StartTimer()
	for i := 0; i < b.N ; i++ {
		f.Test(n1)
	}
}

