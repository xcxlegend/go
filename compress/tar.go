package compress

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path"
)

type TarCompress struct {
}

func (t *TarCompress) Compress(dirpath, dest string) error {
	fw, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fw.Close()

	// gzip write
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// 打开文件夹

	// fmt.Println("tar.gz ok")
	return t.readDirToTar(tw, dirpath, "", path.Base(dest))
	// return nil
}

func (t *TarCompress) readDirToTar(tw *tar.Writer, dirpath, base string, exclude string) error {

	dirpath = path.Join(path.Dir(dirpath), path.Base(dirpath)) + string(os.PathSeparator)
	// beego.Debug(dirpath)
	dir, err := os.Open(dirpath)
	if err != nil {
		return err
	}
	defer dir.Close()

	// 读取文件列表
	fis, err := dir.Readdir(0)
	if err != nil {
		panic(err)
	}

	// 遍历文件列表
	for _, fi := range fis {
		if fi.Name() == exclude {
			continue
		}
		// 逃过文件夹, 我这里就不递归了
		if fi.IsDir() {
			t.readDirToTar(tw, path.Join(dirpath, fi.Name()), path.Join(base, fi.Name()), "")
			continue
		}

		// 打印文件名称
		// fmt.Println(fi.Name())

		// 打开文件
		fr, err := os.Open(path.Join(dir.Name(), fi.Name()))
		if err != nil {
			return err
		}
		defer fr.Close()

		// 信息头
		h := new(tar.Header)
		h.Name = path.Join(base, fi.Name())
		h.Size = fi.Size()
		h.Mode = int64(fi.Mode())
		h.ModTime = fi.ModTime()

		// 写信息头
		err = tw.WriteHeader(h)
		if err != nil {
			return err
		}

		// 写文件
		_, err = io.Copy(tw, fr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TarCompress) Decompress(zipFile, dest string) error {

	return nil
}
