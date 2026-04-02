package main

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	defaultBufferCap     = 1 << 10
	maxRetainedBufferCap = 1 << 12
)

type Record struct {
	User    string
	Action  string
	Tags    []string
	Comment string
}

type formatter interface {
	FormatAll(records []Record, limit int) []string
	BuffersCreated() int
}

type bufferLifecycle interface {
	getBuffer() *bytes.Buffer
	putBuffer(buf *bytes.Buffer)
}

type AllocFormatter struct {
	created atomic.Int32
}

func NewAllocFormatter() *AllocFormatter {
	return &AllocFormatter{}
}

func (f *AllocFormatter) BuffersCreated() int {
	return int(f.created.Load())
}

func (f *AllocFormatter) getBuffer() *bytes.Buffer {
	f.created.Add(1)
	return newOwnedBuffer()
}

func (f *AllocFormatter) putBuffer(_ *bytes.Buffer) {}

func (f *AllocFormatter) FormatAll(records []Record, limit int) []string {
	return formatAllWith(f, records, limit)
}

type FreeListFormatter struct {
	mu      sync.Mutex
	free    []*bytes.Buffer
	created atomic.Int32
}

func NewFreeListFormatter() *FreeListFormatter {
	return &FreeListFormatter{}
}

func (f *FreeListFormatter) BuffersCreated() int {
	return int(f.created.Load())
}

func (f *FreeListFormatter) getBuffer() *bytes.Buffer {
	// TODO: reuse a buffer from the free list when available.
	// TODO: allocate and count a new buffer only when the free list is empty.
	f.created.Add(1)
	return newOwnedBuffer()
}

func (f *FreeListFormatter) putBuffer(buf *bytes.Buffer) {
	// TODO: reset the buffer before reuse.
	// TODO: drop oversized buffers instead of retaining them.
	// TODO: return reusable buffers back to the free list under the mutex.
	_ = buf
}

func (f *FreeListFormatter) FormatAll(records []Record, limit int) []string {
	return formatAllWith(f, records, limit)
}

type PoolFormatter struct {
	pool    sync.Pool
	created atomic.Int32
}

func NewPoolFormatter() *PoolFormatter {
	// TODO: configure pool.New so the pool can allocate counted buffers lazily.
	return &PoolFormatter{}
}

func (f *PoolFormatter) BuffersCreated() int {
	return int(f.created.Load())
}

func (f *PoolFormatter) getBuffer() *bytes.Buffer {
	// TODO: borrow a buffer from sync.Pool.
	f.created.Add(1)
	return newOwnedBuffer()
}

func (f *PoolFormatter) putBuffer(buf *bytes.Buffer) {
	// TODO: reset the buffer before reuse.
	// TODO: drop oversized buffers instead of pooling them.
	// TODO: return reusable buffers back to sync.Pool.
	_ = buf
}

func (f *PoolFormatter) FormatAll(records []Record, limit int) []string {
	return formatAllWith(f, records, limit)
}

func formatAllWith(lifecycle bufferLifecycle, records []Record, limit int) []string {
	if len(records) == 0 {
		return nil
	}
	if limit <= 0 {
		limit = 1
	}

	outputs := make([]string, len(records))
	sem := make(chan struct{}, limit)

	var wg sync.WaitGroup
	for i, record := range records {
		wg.Add(1)
		sem <- struct{}{}

		go func(i int, record Record) {
			defer wg.Done()
			defer func() { <-sem }()

			buf := lifecycle.getBuffer()
			defer lifecycle.putBuffer(buf)

			writeRecord(buf, record)
			outputs[i] = string(append([]byte(nil), buf.Bytes()...))
		}(i, record)
	}

	wg.Wait()
	return outputs
}

func newOwnedBuffer() *bytes.Buffer {
	return bytes.NewBuffer(make([]byte, 0, defaultBufferCap))
}

func writeRecord(buf *bytes.Buffer, record Record) {
	buf.WriteString("user=")
	buf.WriteString(record.User)
	buf.WriteString(" action=")
	buf.WriteString(record.Action)
	buf.WriteString(" tags=[")
	for i, tag := range record.Tags {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(strings.ToUpper(tag))
	}
	buf.WriteString(`] comment="`)
	buf.WriteString(record.Comment)
	buf.WriteByte('"')
}

func main() {
	records := []Record{
		{User: "alice", Action: "login", Tags: []string{"auth", "web"}, Comment: strings.Repeat("ok ", 20)},
		{User: "bob", Action: "deploy", Tags: []string{"ops", "release"}, Comment: strings.Repeat("done ", 25)},
		{User: "carol", Action: "export", Tags: []string{"csv", "report"}, Comment: strings.Repeat("ready ", 18)},
	}

	formatters := []struct {
		name      string
		formatter formatter
	}{
		{name: "alloc", formatter: NewAllocFormatter()},
		{name: "free-list", formatter: NewFreeListFormatter()},
		{name: "sync-pool", formatter: NewPoolFormatter()},
	}

	for _, item := range formatters {
		fmt.Println("formatter:", item.name)
		for _, output := range item.formatter.FormatAll(records, 2) {
			fmt.Println(output)
		}
		fmt.Println("buffers created:", item.formatter.BuffersCreated())
	}
}
