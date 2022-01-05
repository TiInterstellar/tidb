package remotestorage

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"sync"
)

// This a hard coded limit of records.
// If per-record is 64B, we have a 64M file.
const recordLimit = 1024 * 1024

var csvPool = sync.Pool{New: func() interface{} {
	return bytes.NewBuffer(make([]byte, 0, 64*recordLimit))
}}

// CsvWriter will write data into memory and then persist into s3.
//
// Safety: CsvWriter is not thread safe.
type CsvWriter struct {
	store *S3

	id    uint64 // unique id for this table.
	count uint   // Count for the current csv.

	buf    *bytes.Buffer
	writer *csv.Writer
}

func NewCsvWriter(store *S3) *CsvWriter {
	buf := csvPool.Get().(*bytes.Buffer)
	buf.Reset()

	writer := csv.NewWriter(buf)

	return &CsvWriter{
		store:  store,
		buf:    buf,
		writer: writer,
	}
}

func (c *CsvWriter) WriteRecord(record []string) (err error) {
	err = c.writer.Write(record)
	if err != nil {
		return fmt.Errorf("write csv record: %w", err)
	}

	c.count++
	if c.count > recordLimit {
		return c.persist()
	}
	return nil
}

func (c *CsvWriter) persist() (err error) {
	// Make sure all data has been flushed.
	c.writer.Flush()

	err = c.store.Write(
		strconv.FormatUint(c.id, 10),
		int64(c.buf.Len()),
		c.buf)
	if err != nil {
		return fmt.Errorf("persist csv: %w", err)
	}

	c.buf.Reset()
	c.count = 0
	c.id++
	return nil
}

func (c *CsvWriter) Close() error {
	csvPool.Put(c.buf)
	c.buf = nil
	return nil
}
