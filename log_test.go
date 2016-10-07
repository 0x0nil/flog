package flog

import (
	"log"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	w := New("./test")
	l := log.New(w, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	l.Println("test")
}

func TestRotateFile(t *testing.T) {
	w := New("./test")
	now := time.Now()
	w.rotateFile(now)

	now = now.AddDate(0, 0, 1)
	w.rotateFile(now)
}
