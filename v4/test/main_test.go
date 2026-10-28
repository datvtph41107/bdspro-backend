package main

import (
	"encoding/json"
	"testing"
)

var (
	objectResult Player
	arrayResult  PlayerArray
)

func BenchmarkJSONObject(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(int64(len(objectJSON)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var p Player
		if err := json.Unmarshal(objectJSON, &p); err != nil {
			b.Fatal(err)
		}
		objectResult = p
	}
}

func BenchmarkJSONArray(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(int64(len(arrayJSON)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var p PlayerArray
		if err := json.Unmarshal(arrayJSON, &p); err != nil {
			b.Fatal(err)
		}
		arrayResult = p
	}
}

func BenchmarkJSONArrayManual(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(int64(len(arrayJSON)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p, err := ParsePlayerArray(arrayJSON)
		if err != nil {
			b.Fatal(err)
		}
		arrayResult = p
	}
}

func BenchmarkBinaryDecode(b *testing.B) {
	var source PlayerArray
	if err := json.Unmarshal(arrayJSON, &source); err != nil {
		b.Fatal(err)
	}

	payload := EncodePlayerBinary(source)

	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p, err := DecodePlayerBinary(payload)
		if err != nil {
			b.Fatal(err)
		}
		arrayResult = p
	}
}

func BenchmarkBinaryEncodeDecode(b *testing.B) {
	var source PlayerArray
	if err := json.Unmarshal(arrayJSON, &source); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		payload := EncodePlayerBinary(source)

		p, err := DecodePlayerBinary(payload)
		if err != nil {
			b.Fatal(err)
		}
		arrayResult = p
	}
}
