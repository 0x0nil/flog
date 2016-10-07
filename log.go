package flog

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	host = "unknownhost"
)

func init() {
	h, err := os.Hostname()
	if err == nil {
		host = shortHostname(h)
	}
}

func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
}

type date struct {
	year  int
	month time.Month
	day   int
}

func (d *date) equal(t time.Time) bool {
	return t.Year() == d.year && t.Month() == d.month && t.Day() == d.day
}

type SyncWriter struct {
	mu sync.Mutex

	prefix string
	file   *os.File
	last   date
}

func New(prefix string) *SyncWriter {
	return &SyncWriter{prefix: prefix}
}

func logName(prefix string, t time.Time) string {
	name := fmt.Sprintf("%s.%s.%04d-%02d-%02d.log",
		prefix,
		host,
		t.Year(),
		t.Month(),
		t.Day())
	return name
}

func create(prefix string, t time.Time) (*os.File, string, error) {
	name := logName(prefix, t)
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, "", fmt.Errorf("cannot create log: %v", err)
	}

	return f, name, nil
}

func (sw *SyncWriter) rotateFile(now time.Time) error {
	if sw.last.equal(now) && sw.file != nil {
		return nil
	}

	if sw.file != nil {
		sw.file.Close()
	}

	var err error
	sw.file, _, err = create(sw.prefix, now)
	if err != nil {
		return err
	}

	sw.last.year, sw.last.month, sw.last.day = now.Date()

	return nil
}

func (sw *SyncWriter) Write(b []byte) (int, error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	now := time.Now()
	err := sw.rotateFile(now)
	if err != nil {
		return 0, err
	}

	n, err := sw.file.Write(b)
	if err != nil {
		return 0, err
	}

	err = sw.file.Sync()
	if err != nil {
		return 0, err
	}

	return n, err
}
