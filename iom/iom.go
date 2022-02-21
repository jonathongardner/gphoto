package iom // io in memory

import (
  "io"
  "errors"
)

type MemoryReadWriteSeeker struct {
    buf []byte
    pos int
}

func (m *MemoryReadWriteSeeker) Write(p []byte) (n int, err error) {
    minCap := m.pos + len(p)
    if minCap > cap(m.buf) { // Make sure buf has enough capacity:
        buf2 := make([]byte, len(m.buf), minCap+len(p)) // add some extra
        copy(buf2, m.buf)
        m.buf = buf2
    }
    if minCap > len(m.buf) {
        m.buf = m.buf[:minCap]
    }
    copy(m.buf[m.pos:], p)
    m.pos += len(p)
    return len(p), nil
}

func (m *MemoryReadWriteSeeker) Seek(offset int64, whence int) (int64, error) {
    newPos, offs := 0, int(offset)
    switch whence {
    case io.SeekStart:
      newPos = offs
    case io.SeekCurrent:
      newPos = m.pos + offs
    case io.SeekEnd:
      newPos = len(m.buf) + offs
    }
    if newPos < 0 {
      return 0, errors.New("negative result pos")
    }
    m.pos = newPos
    return int64(newPos), nil
}

func (m *MemoryReadWriteSeeker) Read(p []byte) (n int, err error) {
  if m.pos >= len(m.buf) {
    return 0, io.EOF
  }
  length := copy(p, m.buf[m.pos:])
  m.pos += length
  return length, nil
}
