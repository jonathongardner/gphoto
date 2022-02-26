package iom

import (
  "crypto/sha256"
)


func (m *MemoryReadWriteSeeker) SHA256() ([]byte) {
  sum := sha256.Sum256(m.buf)
  return sum[:]
}
