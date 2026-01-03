package infra

import (
	"bytes"
	"io"
	"testing"
	"time"
)

func TestProgressReader_Read(t *testing.T) {
	data := make([]byte, 100)
	readerFake := bytes.NewReader(data)
	totalSize := int64(len(data))

	progressChan := make(chan float64, 100)

	callback := func(p float64, b int64) {

		defer func() { recover() }()
		progressChan <- p
	}

	pr := NewProgressReader(readerFake, totalSize, callback)

	buffer := make([]byte, 10)

	for {
		_, err := pr.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Erro: %v", err)
		}
	}

	lastPercent := 0.0
	timeout := time.After(1 * time.Second)
	completed := false

Loop:
	for {
		select {
		case p := <-progressChan:
			lastPercent = p
			if p == 100.0 {
				completed = true
				break Loop
			}
		case <-timeout:
			t.Error("Timeout: Não recebeu 100% a tempo")
			break Loop
		}
	}

	if !completed {
		t.Errorf("Falhou: parou em %.2f%%", lastPercent)
	}
}
