package remotestorage

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"sync"
)

// This a hard coded size limit of records, we use 64MiB here.
const sizeLimit = 64 * 1024 * 1024

var csvPool = sync.Pool{New: func() interface{} {
	return bytes.NewBuffer(make([]byte, 0, sizeLimit))
}}

// CsvWriter will write data into memory and then persist into s3.
//
// Safety: CsvWriter is not thread safe.
type CsvWriter struct {
	store *S3
	path  string // path id for this csv.

	buf    *bytes.Buffer
	writer *csv.Writer
}

func NewCsvWriter(store *S3, path string) *CsvWriter {
	buf := csvPool.Get().(*bytes.Buffer)
	buf.Reset()

	writer := csv.NewWriter(buf)

	return &CsvWriter{
		store:  store,
		path:   path,
		buf:    buf,
		writer: writer,
	}
}

// IsFree must be called before calling WriteRecord.
//
// It's OK to write a bit more data than the size limit.
func (c *CsvWriter) IsFree() bool {
	return c.buf.Len() < sizeLimit
}

func (c *CsvWriter) WriteRecord(record []string) (err error) {
	err = c.writer.Write(record)
	if err != nil {
		return fmt.Errorf("write csv record: %w", err)
	}
	return nil
}

// Persist will persist all memory data into s3 and free the writer.
// The writer should not be used after this call.
//
// Caller must decide whether to call Persist. Mostly, Persist should be called in the following cases:
// - IsFree returns false: the writer is full, we should persist it.
// - All records are consumed: we should persist the remaining data.
func (c *CsvWriter) Persist() (err error) {
	// Make sure all data has been flushed.
	c.writer.Flush()

	err = c.store.Write(
		c.path,
		int64(c.buf.Len()),
		c.buf)
	if err != nil {
		return fmt.Errorf("persist csv: %w", err)
	}

	c.close()
	return nil
}

func (c *CsvWriter) close() {
	csvPool.Put(c.buf)
	c.buf = nil
}
