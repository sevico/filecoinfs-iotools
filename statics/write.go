package statics

import (
	"os"
	"time"

	"github.com/dustin/go-humanize"
)

//WriteHDDFileWithStatics 带统计的写
func WriteHDDFileWithStatics(file *os.File, content []byte) {
	start := time.Now()
	file.Write(content)
	end := time.Now()
	contentSizeHuman := humanize.Bytes(uint64(len(content)))

	FinishHDDWrite(file, contentSizeHuman, start, end)

}

//WriteSSDFileWithStatics 带统计的写
func WriteSSDFileWithStatics(file *os.File, content []byte) {
	start := time.Now()
	file.Write(content)
	end := time.Now()

	contentSizeHuman := humanize.IBytes(uint64(len(content)))

	FinishSSDWrite(file, contentSizeHuman, start, end)

}
