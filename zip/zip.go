package main

import (
	"github.com/xcxlegend/go/compress"
	"os"
)

func main() {
	var zipc = new(compress.ZipCompress)
	const filename = "D:\\LegendXie\\测试文档\\test.zip"
	zipc.Decompress(filename, "D:\\LegendXie\\测试文档\\z\\")
	os.Remove(filename)
	// DeCompressZip()
}
