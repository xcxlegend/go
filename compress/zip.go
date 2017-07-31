package compress

import (
	"archive/zip"
	"io"
	"os"
	"path"
)

type ZipCompress struct {
}

func (z *ZipCompress) Compress(dirpath, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	f, err := os.Open(dirpath)
	err = compress(f, "", w)
	// fmt.Println(err)
	return err
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = path.Join(prefix, info.Name())
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(path.Join(file.Name(), fi.Name()))
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = path.Join(prefix, header.Name)
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
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
