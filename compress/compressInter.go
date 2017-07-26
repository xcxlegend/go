package compress

type CompressTool interface {
	Compress(dirpath, dest string) error
	Decompress(zipFile, dest string) error
}
