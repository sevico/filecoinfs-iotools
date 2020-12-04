package statics

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sevico/filecoinfs-iotools/model"

	logs "log"

	"github.com/sevico/filecoinfs-iotools/util"

	"github.com/sevico/tachymeter"
	"github.com/thecodeteam/goodbye"
)

var OpLogToFileChan = make(chan model.OpLog)

func formatString(l []model.OpLog) string {
	sRes := ""

	for i := range l {
		temp := fmt.Sprintf("%v|%v|%v|%v|%v\n", l[i].Ts.Format(time.RFC3339), l[i].OP, l[i].Object, l[i].Content, strconv.Itoa(int(l[i].Duration.Microseconds()))+"μs")
		sRes += temp
	}
	return sRes
}

func WriteToLog(debug bool, sampleCount int) {

	tWriteHDD := tachymeter.New(&tachymeter.Config{Size: sampleCount})
	tWriteSSD := tachymeter.New(&tachymeter.Config{Size: sampleCount})
	tRandomRead := tachymeter.New(&tachymeter.Config{Size: sampleCount})
	tCloseHDD := tachymeter.New(&tachymeter.Config{Size: sampleCount})
	tCloseSSD := tachymeter.New(&tachymeter.Config{Size: sampleCount})

	// var opToTachy = map[string]*tachymeter.Tachymeter{
	// 	"writeHDD": tWriteHDD,
	// 	"writeSSD": tWriteSSD,
	// 	"read":     tRandomRead,
	// 	"closeHDD": tCloseHDD,
	// 	"closeSSD": tCloseSSD,
	// }

	goodbye.Register(func(ctx context.Context, sig os.Signal) {
		tl := tachymeter.Timeline{}

		fmt.Printf("the %v statics:%v\n Histogram:\n %v\n", "writeHDD", tWriteHDD.Calc(), tWriteHDD.Calc().Histogram.String(100))
		fmt.Println("===========================")
		fmt.Printf("the %v statics:%v\n Histogram:\n %v\n", "writeSSD", tWriteSSD.Calc(), tWriteSSD.Calc().Histogram.String(100))
		fmt.Println("===========================")
		fmt.Printf("the %v statics:%v\n Histogram:\n %v\n", "randomRead", tRandomRead.Calc(), tRandomRead.Calc().Histogram.String(100))
		fmt.Println("===========================")
		fmt.Printf("the %v statics:%v\n Histogram:\n %v\n", "closeHDD", tCloseHDD.Calc(), tCloseHDD.Calc().Histogram.String(100))
		fmt.Println("===========================")
		fmt.Printf("the %v statics:%v\n Histogram:\n %v\n", "closeSSD", tCloseSSD.Calc(), tCloseSSD.Calc().Histogram.String(100))
		fmt.Println("===========================")

		if tWriteHDD.Calc().Count > 0 {
			tl.AddEvent(tWriteHDD.Calc(), "WriteHDD")
		}
		if tWriteSSD.Calc().Count > 0 {
			tl.AddEvent(tWriteSSD.Calc(), "WriteSSD")
		}
		if tRandomRead.Calc().Count > 0 {
			tl.AddEvent(tRandomRead.Calc(), "RandomRead")
		}
		if tCloseHDD.Calc().Count > 0 {
			tl.AddEvent(tCloseHDD.Calc(), "CloseHDD")
		}
		if tCloseSSD.Calc().Count > 0 {
			tl.AddEvent(tCloseSSD.Calc(), "CloseSSD")
		}
		tl.WriteHTML(util.GetCwd())
	})

	var batch []model.OpLog

	now := time.Now()

	fileName := "op_" + now.Format("2006-01-02-15-04-05") + ".log"
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)

	// fileNameWriteHDD := "writeHDD_" + now.Format("2006-01-02-15-04-05") + ".log"
	// fileWriteHDD, err := os.OpenFile(fileNameWriteHDD, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)

	if err != nil {
		logs.Fatalf("open file %v err %v", fileName, err)
		return
	}
	defer file.Close()
	// opCount := make(map[string]int)

	for logEntry := range OpLogToFileChan {
		if debug {
			batch = append(batch, logEntry)
			// logs.Info("%+v", logEntry)
			if len(batch) > 20 {
				file.WriteString(formatString(batch))
				file.Sync()
				batch = nil
			}
		}
		switch logEntry.OP {
		case "writeHDD":
			tWriteHDD.AddTime(logEntry.Duration)
		case "writeSSD":
			tWriteSSD.AddTime(logEntry.Duration)
		case "read":
			tRandomRead.AddTime(logEntry.Duration)
		case "closeHDD":
			tCloseHDD.AddTime(logEntry.Duration)
		case "closeSSD":
			tCloseSSD.AddTime(logEntry.Duration)
		}
		// opCount[logEntry.OP]++
		// if opCount[logEntry.OP]%100 == 0 {
		// 	fmt.Printf("the %v Histogram:\n %v\n", logEntry.OP, opToTachy[logEntry.OP].Calc().Histogram.String(25))
		// }
	}

}

func FinishHDDWrite(file *os.File, content string, start time.Time, end time.Time) {
	var temp model.OpLog
	temp.Ts = start
	temp.Content = content
	temp.Duration = end.Sub(start)
	temp.OP = "writeHDD"
	temp.Object = file.Name()
	OpLogToFileChan <- temp

}
func FinishSSDWrite(file *os.File, content string, start time.Time, end time.Time) {
	var temp model.OpLog
	temp.Ts = start
	temp.Content = content
	temp.Duration = end.Sub(start)
	temp.OP = "writeSSD"
	temp.Object = file.Name()
	OpLogToFileChan <- temp

}

func FinishRead(file *os.File, content string, start time.Time, dur time.Duration) {
	var temp model.OpLog
	temp.Ts = start
	temp.Content = content
	temp.Duration = dur
	temp.OP = "read"
	temp.Object = file.Name()
	OpLogToFileChan <- temp

}
