package profiler

import "sync"

var RowPool = sync.Pool{
	New: func() any {
		return make([]string, 0, 10)
	},
}

func GetRowSlice() []string {
	return RowPool.Get().([]string)
}

func PutRowSlice(row []string) {
	RowPool.Put(row[:0])
}