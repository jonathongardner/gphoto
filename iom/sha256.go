package iom

import (
  "crypto/sha256"
  "encoding/hex"
)


func (m *MemoryReadWriteSeeker)  SHA256() (string) {
  sum := sha256.Sum256(m.buf)
  return hex.EncodeToString(sum[:])
}
