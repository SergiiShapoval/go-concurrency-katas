package main

import (
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"
)

var benchmarkSink []string

func formatterCases() []struct {
	name         string
	newFormatter func() formatter
	checkCreated func(t *testing.T, got int)
} {
	return []struct {
		name         string
		newFormatter func() formatter
		checkCreated func(t *testing.T, got int)
	}{
		{
			name:         "alloc",
			newFormatter: func() formatter { return NewAllocFormatter() },
			checkCreated: func(t *testing.T, got int) {
				t.Helper()
				if got != 2 {
					t.Fatalf("buffers created mismatch: got %d want 2", got)
				}
			},
		},
		{
			name:         "free_list",
			newFormatter: func() formatter { return NewFreeListFormatter() },
			checkCreated: func(t *testing.T, got int) {
				t.Helper()
				if got != 1 {
					t.Fatalf("buffers created mismatch: got %d want 1", got)
				}
			},
		},
		{
			name:         "sync_pool",
			newFormatter: func() formatter { return NewPoolFormatter() },
			checkCreated: func(t *testing.T, got int) {
				t.Helper()
				if got < 1 || got > 2 {
					t.Fatalf("buffers created mismatch: got %d want 1 or 2", got)
				}
			},
		},
	}
}

func TestFormattersFormatInputsInOrder(t *testing.T) {
	records := []Record{
		{User: "alice", Action: "login", Tags: []string{"auth", "web"}, Comment: "ok"},
		{User: "bob", Action: "deploy", Tags: []string{"ops"}, Comment: "done"},
	}
	want := []string{
		`user=alice action=login tags=[AUTH,WEB] comment="ok"`,
		`user=bob action=deploy tags=[OPS] comment="done"`,
	}

	for _, tt := range formatterCases() {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.newFormatter().FormatAll(records, 2)
			if !slices.Equal(got, want) {
				t.Fatalf("outputs mismatch: got %v want %v", got, want)
			}
		})
	}
}

func TestFormattersCopyOutputBeforeBufferReuse(t *testing.T) {
	for _, tt := range formatterCases() {
		t.Run(tt.name, func(t *testing.T) {
			formatter := tt.newFormatter()

			first := formatter.FormatAll([]Record{
				{User: "alice", Action: "upload", Tags: []string{"blob"}, Comment: strings.Repeat("A", 256)},
			}, 1)
			second := formatter.FormatAll([]Record{
				{User: "bob", Action: "delete", Tags: []string{"cleanup"}, Comment: "short"},
			}, 1)

			wantFirst := []string{
				`user=alice action=upload tags=[BLOB] comment="` + strings.Repeat("A", 256) + `"`,
			}
			wantSecond := []string{
				`user=bob action=delete tags=[CLEANUP] comment="short"`,
			}

			if !slices.Equal(first, wantFirst) {
				t.Fatalf("first outputs mismatch: got %v want %v", first, wantFirst)
			}
			if !slices.Equal(second, wantSecond) {
				t.Fatalf("second outputs mismatch: got %v want %v", second, wantSecond)
			}
		})
	}
}

func TestFormattersTrackCreatedBuffersAcrossBatches(t *testing.T) {
	record := []Record{{User: "alice", Action: "login", Comment: "ok"}}

	for _, tt := range formatterCases() {
		t.Run(tt.name, func(t *testing.T) {
			formatter := tt.newFormatter()

			_ = formatter.FormatAll(record, 1)
			_ = formatter.FormatAll(record, 1)

			tt.checkCreated(t, formatter.BuffersCreated())
		})
	}
}

func TestFormattersHandleEmptyInput(t *testing.T) {
	for _, tt := range formatterCases() {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.newFormatter().FormatAll(nil, 2)
			if len(got) != 0 {
				t.Fatalf("expected empty output, got %v", got)
			}
		})
	}
}

func BenchmarkFormattersSequential(b *testing.B) {
	benchmarkFormatters(b, 1)
}

func BenchmarkFormattersConcurrent(b *testing.B) {
	benchmarkFormatters(b, runtime.NumCPU())
}

func benchmarkFormatters(b *testing.B, limit int) {
	workloads := []struct {
		name    string
		records []Record
	}{
		{name: "retainable", records: benchmarkRecords(80)},
		{name: "oversized", records: benchmarkRecords(280)},
	}

	for _, workload := range workloads {
		b.Run(workload.name, func(b *testing.B) {
			for _, tt := range formatterCases() {
				b.Run(tt.name, func(b *testing.B) {
					formatter := tt.newFormatter()
					b.ReportAllocs()
					for b.Loop() {
						benchmarkSink = formatter.FormatAll(workload.records, limit)
					}
				})
			}
		})
	}
}

func benchmarkRecords(commentRepeats int) []Record {
	records := make([]Record, 0, runtime.NumCPU()*8)
	for i := 0; i < cap(records); i++ {
		records = append(records, Record{
			User:    "user-" + strconv.Itoa(i),
			Action:  "render",
			Tags:    []string{"report", "batch", "export"},
			Comment: strings.Repeat("payload-block-"+strconv.Itoa(i%10)+" ", commentRepeats),
		})
	}
	return records
}
