package infra

import (
	"io"
	"math"
	"sync/atomic"
	"time"
)

type ProgressListener func(percentage float64, bytesRead int64)

type ProgressReader struct {
	Reader     io.Reader
	TotalSize  int64
	OnProgress ProgressListener

	currentBytes int64
	lastPercent  int64
	lastTime     time.Time
}

func (pr *ProgressReader) Read(p []byte) (int, error) {

	n, err := pr.Reader.Read(p)

	newBytes := atomic.AddInt64(&pr.currentBytes, int64(n))

	if pr.TotalSize > 0 && pr.OnProgress != nil {

		percent := int64(math.Floor((float64(newBytes) / float64(pr.TotalSize)) * 80))

		if percent > pr.lastPercent || time.Since(pr.lastTime) > 1*time.Second {
			pr.lastPercent = percent
			pr.lastTime = time.Now()

			go pr.OnProgress(float64(percent), newBytes)
		}
	}

	return n, err
}

func NewProgressReader(r io.Reader, totalSize int64, onProgress ProgressListener) *ProgressReader {
	return &ProgressReader{
		Reader:     r,
		TotalSize:  totalSize,
		OnProgress: onProgress,
		lastTime:   time.Now(),
	}
}
