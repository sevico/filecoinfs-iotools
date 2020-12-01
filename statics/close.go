package statics

import (
	"infcs/filecoinfs-iotools/model"
	"os"
	"time"
)

func CloseHDD(file *os.File, Content string) {

	t1 := time.Now()
	file.Close()
	t2 := time.Now()

	temp := model.OpLog{
		Ts:       time.Now(),
		OP:       "closeHDD",
		Object:   file.Name(),
		Content:  Content,
		Duration: t2.Sub(t1),
	}
	OpLogToFileChan <- temp
}

func CloseSSD(file *os.File, Content string) {

	t1 := time.Now()
	file.Close()
	t2 := time.Now()

	temp := model.OpLog{
		Ts:       time.Now(),
		OP:       "closeSSD",
		Object:   file.Name(),
		Content:  Content,
		Duration: t2.Sub(t1),
	}
	OpLogToFileChan <- temp
}
