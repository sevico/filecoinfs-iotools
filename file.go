package main

import (
	"fmt"
	"infcs/filecoinfs-iotools/model"
	"infcs/filecoinfs-iotools/statics"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	logs "log"

	"github.com/dustin/go-humanize"
	"github.com/hashicorp/go-uuid"
	cmap "github.com/orcaman/concurrent-map"
)

var genedFileName = cmap.New()
var finishedFileName = cmap.New()
var lock sync.Mutex
var fileNameChan = make(chan string)
var c = sync.NewCond(&sync.Mutex{})

var finishCount int32

func UUIDCheckFunc(id string) (bool, error) {

	lock.Lock()
	defer lock.Unlock()

	var exist = true
	if _, ok := genedFileName.Get(id); ok {
		return false, nil
	}

	HDDDirPath := path.Join(cnf.FSPath, cnf.NameArgs.HDDDirName, subDirName)
	// SSDDirPath := path.Join(cnf.BeeFSPath, cnf.NameArgs.SSDDirName)

	path := path.Join(HDDDirPath, id)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return !exist, nil
}
func genFileName() string {
	for {
		var id string
		for {
			id, _ = uuid.GenerateUUID()
			success, _ := UUIDCheckFunc(id)
			if success {
				break
			}
		}
		genedFileName.Set(id, true)

		fileNameChan <- id
	}

}

var HDDChunkSize = 0 //conf
var SSDChunkSize = 0 //conf
var subDirName = ""

func genFile() {
	hddDirPath := path.Join(cnf.FSPath, subDirName, cnf.NameArgs.HDDDirName)
	ssdDirPath := path.Join(cnf.FSPath, subDirName, cnf.NameArgs.SSDDirName)
	hundredMBytes := make([]byte, HDDChunkSize)
	fourKBBytes := make([]byte, SSDChunkSize)
	_, err := rand.Read(hundredMBytes)
	_, err = rand.Read(fourKBBytes)
	for {
		func() {
			fileName := <-fileNameChan

			var SSDDirPaths = make([]string, cnf.TimeArgs.GenTimes)
			HDDFilePath := path.Join(hddDirPath, fileName)
			// fmt.Println(HDDFilePath)
			for i := 0; i < cnf.TimeArgs.GenTimes; i++ {
				SSDDirPaths[i] = path.Join(ssdDirPath, fileName+"_"+strconv.Itoa(i+1))
			}

			hddFile, _ := os.OpenFile(HDDFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
			defer statics.CloseHDD(hddFile, fmt.Sprintf("%.2f", cnf.SizeArgs.HDD)+"GB")

			for i := 0; i < int(int(cnf.SizeArgs.HDD*1024*1024*1024)/HDDChunkSize); i++ {

				if err != nil {
					logs.Fatalf("rand read err %v", err)
				}
				statics.WriteHDDFileWithStatics(hddFile, hundredMBytes)
			}
			residue := int(cnf.SizeArgs.HDD*1024*1024*1024) % HDDChunkSize
			if cnf.StaticResidue {
				statics.WriteHDDFileWithStatics(hddFile, hundredMBytes[:residue])
			} else {
				hddFile.Write(hundredMBytes[:residue])
			}

			for i := 0; i < len(SSDDirPaths); i++ {
				ssdFile, _ := os.OpenFile(SSDDirPaths[i], os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)

				defer statics.CloseSSD(ssdFile, fmt.Sprintf("%.2f", cnf.SizeArgs.SSD)+"MB")
				for j := 0; j < int(int(cnf.SizeArgs.SSD*1024*1024)/SSDChunkSize); j++ {

					if err != nil {
						logs.Fatalf("rand read err %v", err)
					}
					statics.WriteSSDFileWithStatics(ssdFile, fourKBBytes)
				}
				residue := int(cnf.SizeArgs.SSD*1024*1024) % SSDChunkSize
				if cnf.StaticResidue {
					statics.WriteSSDFileWithStatics(ssdFile, fourKBBytes[:residue])
				} else {
					ssdFile.Write(fourKBBytes[:residue])
				}
				// util.WriteSSDFileWithStatics(ssdFile, fourKBBytes[:residue])
			}
			finishedFileName.Set(fileName, true)
			writeToNames(fileName)
			c.L.Lock()
			finishCount += 1
			c.L.Unlock()
			c.Broadcast()
		}()

		// if int(finishCount) > cnf.TimeArgs.Threshold+20 {
		// 	c.Broadcast()
		// }
	}

}

func getNameFiles() []string {
	var nameFiles []string

	files, err := ioutil.ReadDir(cnf.FSPath)
	if err != nil {
		logs.Printf("ioutil.ReadDir the dir:%v", cnf.FSPath)
	}
	for _, item := range files {
		if item.IsDir() == true && strings.Contains(item.Name(), "_"+cnf.NameArgs.NameFile) == false {
			continue
		}
		nameFiles = append(nameFiles, item.Name())
	}
	return nameFiles
}

func getAllFilename() []model.FileName {
	var fileNames []model.FileName
	for _, item := range getNameFiles() {
		// if item.IsDir() == true && strings.Contains(item.Name(), "_"+cnf.NameArgs.NameFile) == false {
		// 	continue
		// }
		subDir := strings.Split(item, "_")[0]

		content, err := ioutil.ReadFile(path.Join(cnf.FSPath, item))
		if err != nil {
			logs.Fatalf("ReadFile %v err %v", path.Join(cnf.FSPath, item), err)
			continue
		}
		fileNamesPerfile := strings.Split(string(content), "\n")
		for _, item := range fileNamesPerfile {
			if item == "" {
				continue
			}
			temp := model.FileName{
				SubDir:   subDir,
				DataFile: item,
			}
			fileNames = append(fileNames, temp)
		}

	}

	return fileNames
}

func getFilename(nameFile string) {

}

func RandomReads() {
	var wg sync.WaitGroup
	for {

		fileNames := getAllFilename()
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(fileNames), func(i, j int) { fileNames[i], fileNames[j] = fileNames[j], fileNames[i] })
		selectedFileNames := fileNames[:cnf.TimeArgs.Threshold]
		for _, item := range selectedFileNames {
			wg.Add(1)
			go Random4KRead(item, &wg)
		}
		wg.Wait()
	}

}

func Random4KRead(fileName model.FileName, wg *sync.WaitGroup) {
	defer wg.Done()
	var filePath string
	if cnf.SubPathMode {
		filePath = path.Join(cnf.FSPath, fileName.SubDir, cnf.NameArgs.HDDDirName, fileName.DataFile)
	} else {
		filePath = path.Join(cnf.FSPath, cnf.NameArgs.HDDDirName, fileName.DataFile)
	}

	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		logs.Fatalf("open file %v err %v", filePath, err)
	}
	stat, err := f.Stat()
	if err != nil {
		logs.Fatalf("stat file %v err %v", filePath, err)
		return
	}
	fileSize := stat.Size()
	readSize := cnf.SizeArgs.RandomRead
	offset := rand.Intn(int(fileSize) - readSize)
	readBuf := make([]byte, readSize)
	t1 := time.Now()
	_, err = f.ReadAt(readBuf, int64(offset))
	dur := time.Since(t1)
	content := humanize.IBytes(uint64(readSize))

	statics.FinishRead(f, content, t1, dur)
	if err != nil {
		logs.Fatalf("read offset %v length %v error %v", offset, readSize, err)
	}

}
func RecoverGenStatus() {
	if cnf.SubPathMode {
		fs, err := ioutil.ReadDir(cnf.FSPath)
		if err != nil {
			logs.Fatalf("ioutil.ReadDir %v err %v", cnf.FSPath, err)
			return
		}
		for _, item := range fs {
			if item.IsDir() == false || item.Name() == cnf.NameArgs.HDDDirName || item.Name() == cnf.NameArgs.SSDDirName {
				continue
			}
			recoverGenStatus(path.Join(cnf.FSPath, item.Name()))
		}
	} else {
		recoverGenStatus(cnf.FSPath)
	}
}

func recoverGenStatus(dir string) {

	HDDDirPath := path.Join(dir, cnf.NameArgs.HDDDirName)

	fmt.Println(HDDDirPath)
	files, err := ioutil.ReadDir(HDDDirPath)
	if err != nil {
		logs.Fatalf("ioutil.ReadDir %v", err)
	}
	for _, item := range files {
		fi, err := os.Stat(path.Join(HDDDirPath, item.Name()))
		if err != nil {
			logs.Fatalf("os.stat %v err %v", path.Join(HDDDirPath, item.Name()), err)
			continue
		}
		if fi.Size() != int64(cnf.SizeArgs.HDD*1024*1024*1024) {
			os.Remove(path.Join(HDDDirPath, item.Name()))
			continue
		}
		ssdEncountErr := false
		for i := 0; i < cnf.TimeArgs.GenTimes; i++ {

			SSDFilePath := path.Join(dir, cnf.NameArgs.SSDDirName, item.Name()+"_"+strconv.Itoa(i+1))
			ssdFi, err := os.Stat(SSDFilePath)
			if err != nil {
				logs.Fatalf("os.Stat(%v) err %v", SSDFilePath, err)
				os.Remove(SSDFilePath)
				ssdEncountErr = true
				continue
			}

			// fmt.Println(ssdFi.Name(), ssdFi.Size(), int64(cnf.SizeArgs.SSD*1024*1024))
			if ssdFi.Size() != int64(cnf.SizeArgs.SSD*1024*1024) {
				ssdEncountErr = true
			}
			if ssdEncountErr {
				os.Remove(SSDFilePath)
				break
			}
		}
		if ssdEncountErr {
			continue
		}
		fmt.Println("recover file:", item.Name())
		genedFileName.Set(item.Name(), true)
		finishedFileName.Set(item.Name(), true)
		writeToNames(item.Name())
		finishCount++
	}

}

func getNamesPath() string {
	host, err := os.Hostname()
	if err != nil {
		logs.Fatalf("get hostname err %v", err)
	}
	namesPath := path.Join(cnf.FSPath, host+"_"+cnf.NameArgs.NameFile)
	return namesPath
}

// writeToNames 将生成的文件名写到文件中
func writeToNames(fileName string) {

	namesFile, err := os.OpenFile(getNamesPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logs.Fatalf("join path %v err %v", getNamesPath(), err)
	}
	defer namesFile.Close()
	_, err = namesFile.WriteString(fileName + "\n")
	if err != nil {
		logs.Fatalf("write string %v err %v", getNamesPath(), err)
	}

}
