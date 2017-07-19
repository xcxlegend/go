package compress

type CompressTool interface {
	Compress() error
	Decompress(zipFile, dest string) error
}
