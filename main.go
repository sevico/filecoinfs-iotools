package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/sevico/filecoinfs-iotools/model"
	"github.com/sevico/filecoinfs-iotools/statics"
	"github.com/sevico/filecoinfs-iotools/util"

	logs "log"
	"os"
	"path"

	"github.com/thecodeteam/goodbye"
	"gopkg.in/yaml.v2"
)

var cnf model.Config

func InitConfig(dir string) bool {

	configFullName := path.Join(util.GetCwd(), fmt.Sprintf("%s/%s", dir, "config.yaml"))

	logs.Printf("【INIT_CONFIG】Use config [%s]", configFullName)

	data, err := ioutil.ReadFile(configFullName)
	if err != nil {
		logs.Fatalf("【INIT_CONFIG】Read config error %v\n", err)
		panic(err)
	}

	if err = yaml.Unmarshal(data, &cnf); err != nil {
		logs.Fatalf("【INIT_CONFIG】Unmarshal conf yaml error %v\n", err)
		panic(err)
	}

	logs.Printf("【INIT_CONFIG】Config InitConfig successfully.")
	return true
}

func init() {
	logs.SetFlags(logs.LstdFlags | logs.Lshortfile)
	InitConfig("conf")
	HDDChunkSize = cnf.SizeArgs.HDDChunk * 1024 * 1024
	SSDChunkSize = cnf.SizeArgs.SSDChunk * 1024

}

func removeNameFiles() {
	nameFiles := getNameFiles()
	for _, fileName := range nameFiles {
		os.Remove(path.Join(cnf.FSPath, fileName))
	}
}

func main() {
	ctx := context.Background()
	defer goodbye.Exit(ctx, -1)
	goodbye.Notify(ctx)

	go statics.WriteToLog(cnf.DebugLog, cnf.TimeArgs.Samples)

	if cnf.SubPathMode {
		subDirName, _ = os.Hostname()
	}

	removeNameFiles()

	//conf
	if cnf.CleanUp {
		os.RemoveAll(cnf.FSPath)
	} else {
		RecoverGenStatus()
	}
	os.MkdirAll(path.Join(cnf.FSPath, subDirName, cnf.NameArgs.HDDDirName), os.ModePerm)
	os.MkdirAll(path.Join(cnf.FSPath, subDirName, cnf.NameArgs.SSDDirName), os.ModePerm)

	for i := 0; i < cnf.TimeArgs.NameRountines; i++ {
		go genFile()
	}

	for i := 0; i < cnf.TimeArgs.GenFileRoutines; i++ {
		go genFileName()
	}

	c.L.Lock()
	for int(finishCount) <= cnf.TimeArgs.Threshold+20 {
		c.Wait()
	}
	c.L.Unlock()
	RandomReads()
	select {}
}
