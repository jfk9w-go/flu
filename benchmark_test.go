package flu_test

import "testing"

func BenchmarkInterface_PointerCall(b *testing.B) {
	b.ReportAllocs()
	var v Interface = new(Struct)
	sum := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum += v.PointerCall()
	}
}

func BenchmarkInterface_ValueCall(b *testing.B) {
	b.ReportAllocs()
	var v Interface = new(Struct)
	sum := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum += v.ValueCall()
	}
}

func BenchmarkFunctionPointer_PointerCall(b *testing.B) {
	b.ReportAllocs()
	var v Interface = new(Struct)
	f := v.PointerCall
	sum := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum += f()
	}
}

func BenchmarkFunctionPointer_ValueCall(b *testing.B) {
	b.ReportAllocs()
	var v Interface = new(Struct)
	f := v.ValueCall
	sum := 0
	for i := 0; i < b.N; i++ {
		sum += f()
	}
}

func BenchmarkDirect_PointerCall(b *testing.B) {
	b.ReportAllocs()
	var v = new(Struct)
	sum := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum += v.PointerCall()
	}
}

func BenchmarkDirect_ValueCall(b *testing.B) {
	b.ReportAllocs()
	var v = Struct{}
	sum := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum += v.ValueCall()
	}
}

type Interface interface {
	PointerCall() int
	ValueCall() int
}

type Struct struct {
}

func (s *Struct) PointerCall() int {
	return 0
}

func (s Struct) ValueCall() int {
	return 0
}
