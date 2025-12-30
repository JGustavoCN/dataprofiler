package profiler

import "testing"

var benchRow = []string{"João", "30", "São Paulo", "Brasil", "Developer", "Go"}

func BenchmarkSemPool(b *testing.B) {
    for i := 0; i < b.N; i++ {
        row := make([]string, len(benchRow))
        copy(row, benchRow)
    }
}

func BenchmarkComPool(b *testing.B) {
    for i := 0; i < b.N; i++ {
        row := GetRowSlice()
        row = append(row, benchRow...)
        PutRowSlice(row)
    }
}