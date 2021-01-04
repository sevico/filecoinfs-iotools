package statics

import (
	logs "log"
	"os"
	"time"

	"github.com/sevico/filecoinfs-iotools/model"
)

func CloseHDD(file *os.File, Content string) {

	t1 := time.Now()
	err := file.Close()
	if err != nil {
		logs.Printf("rand read fourKBBytes err %v", err)
	}
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
	err := file.Close()

	if err != nil {
		logs.Printf("CloseSSD file %v err %v", file, err)
	}
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
