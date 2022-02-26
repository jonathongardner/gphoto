package photo

import (
	"io"
  "time"
  "strconv"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/exif"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/xmp"
)

func YearPhotoTaken(read meta.Reader) (string) {
  var err error
	var x xmp.XMP
	var e *exif.Data
  var yr time.Time

	exifDecodeFn := func(r io.Reader, m *meta.Metadata) error {
		e, err = e.ParseExifWithMetadata(read, m)
		return nil
	}
	xmpDecodeFn := func(r io.Reader, m *meta.Metadata) error {
		x, err = xmp.ParseXmp(r)
		return err
	}

	_, err = imagemeta.NewMetadata(read, xmpDecodeFn, exifDecodeFn)
  if e != nil {
    yr, err = e.DateTime()
  }

	if err != nil {
		return "unknown"
	}
  return strconv.Itoa(yr.Year())
}
