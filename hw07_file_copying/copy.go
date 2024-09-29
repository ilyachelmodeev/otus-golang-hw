package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	in, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	fi, err := in.Stat()
	if err != nil {
		return err
	}
	fm := fi.Mode()
	if !fm.IsRegular() {
		return ErrUnsupportedFile
	}
	inSize := fi.Size()
	if offset > inSize {
		return ErrOffsetExceedsFileSize
	}
	_, err = in.Seek(offset, 0)
	if err != nil {
		return err
	}
	out, err := os.OpenFile(toPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	var barSize int64
	if limit == 0 {
		barSize = inSize - offset
	} else {
		if limit < inSize-offset {
			barSize = limit
		} else {
			barSize = inSize - offset
		}
	}
	bar := pb.Full.Start64(barSize)

	var done int64
	buf := make([]byte, 1024*4)
	for limit == 0 || done < limit {
		read, err := in.Read(buf)
		if errors.Is(err, io.EOF) {
			bar.Finish()
			break
		}
		if err != nil {
			return err
		}
		bar.Add(read)
		var written int
		if int64(read) < limit-done || limit == 0 {
			written, err = out.Write(buf[0:read])
		} else {
			written, err = out.Write(buf[0 : limit-done])
		}
		if err != nil {
			return err
		}
		done += int64(written)
	}
	return nil
}
