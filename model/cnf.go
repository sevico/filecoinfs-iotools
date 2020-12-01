package model

type SizeArgs struct {
	HDD        float64 `yaml:"HDD"`
	SSD        float64 `yaml:"SSD"`
	RandomRead int     `yaml:"RandomRead"`
	HDDChunk   int     `yaml:"HDDChunk"`
	SSDChunk   int     `yaml:"SSDChunk"`
}
type Times struct {
	RandomRead      int `yaml:"RandomRead"`
	GenTimes        int `yaml:"GenTimes"`
	Threshold       int `yaml:"Threshold"`
	NameRountines   int `yaml:"NameRountines"`
	GenFileRoutines int `yaml:"GenFileRoutines"`
	Samples         int `yaml:"Samples"`
}
type Names struct {
	HDDDirName string `yaml:"HDDDirName"`
	SSDDirName string `yaml:"SSDDirName"`
	NameFile   string `yaml:"NameFile"`
}

type Config struct {
	SizeArgs      SizeArgs `yaml:"Sizes"`
	TimeArgs      Times    `yaml:"Times"`
	NameArgs      Names    `yaml:"Names"`
	FSPath        string   `yaml:"FSPath"`
	CleanUp       bool     `yaml:"CleanUp"`
	SubPathMode   bool     `yaml:"SubPath"`
	DebugLog      bool     `yaml:"DebugLog"`
	StaticResidue bool     `yaml:"StaticResidue"`
}
