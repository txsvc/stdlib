package batching

import (
	"sync"
	"time"
)

// BatchWriter will accept messages and invoke the Writer when the batch
// requirements have been fulfilled (either batch size or interval have been
// exceeded). BatchWriter should be created with NewBatchWriter().
type BatchWriter struct {
	w        Writer
	size     int
	interval time.Duration
	batch    []interface{}
	lastSent time.Time
	mu       sync.Mutex
}

// Writer is used to submit the completed batch. The batch may be partial if
// the interval lapsed instead of filling the batch.
type Writer interface {
	// Write submits the batch.
	Write(batch []interface{})
}

// WriterFunc is an adapter to allow ordinary functions to be a Writer.
type WriterFunc func(batch []interface{})

// Write implements Writer.
func (f WriterFunc) Write(batch []interface{}) {
	f(batch)
}

// NewBatchWriter creates a new BatchWriter. For clarity, it is recommenended to use a
// wrapper type such as NewByteBatchWriter or NewV2EnvelopeBatchWriter vs using this directly.
func NewBatchWriter(size int, interval time.Duration, writer Writer) *BatchWriter {
	return &BatchWriter{
		size:     size,
		interval: interval,
		w:        writer,
		lastSent: time.Now(),
	}
}

// Write stores data to the batch. It will not submit the batch to the writer
// until either the batch has been filled, or the interval has lapsed.
func (b *BatchWriter) Write(data interface{}) {
	b.batch = append(b.batch, data)

	if b.partialBatch() && b.partialInterval() {
		return
	}

	b.writeBatch()
}

// ForcedFlush bypasses the batch interval and batch size checks and writes immediately.
func (b *BatchWriter) ForcedFlush() {
	b.writeBatch()
}

// Flush will write a partial batch if there is data and the interval has
// lapsed. Otherwise it is a NOP. This method should be called freqently to
// make sure batches do not stick around for long periods of time.
func (b *BatchWriter) Flush() {
	if b.partialInterval() {
		return
	}

	b.writeBatch()
}

// writeBatch writes the batch (if any) to the writer and resets the batch and interval.
func (b *BatchWriter) writeBatch() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.batch) == 0 {
		return
	}

	b.w.Write(b.batch)

	b.batch = nil
	b.lastSent = time.Now()
}

func (b *BatchWriter) partialBatch() bool {
	return len(b.batch) < b.size
}

func (b *BatchWriter) partialInterval() bool {
	return time.Since(b.lastSent) < b.interval
}
