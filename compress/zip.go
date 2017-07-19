package compress

import (
	"archive/zip"
	"io"
	"os"
	"path"
)

type ZipCompress struct {
}

func (z *ZipCompress) Compress() error {
	return nil
}

func (z *ZipCompress) Decompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := dest + file.Name
		err = os.MkdirAll(path.Dir(filename), 0755)
		if err != nil {
			continue
		}
		w, err := os.Create(filename)
		if err != nil {
			continue
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			continue
		}
		w.Close()
		rc.Close()
	}
	return nil
}
