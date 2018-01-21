package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/labstack/gommon/log"
)

type Seed struct {
	Execution string `json:"execution"`
	Files     []File `json:"files"`
	Hostname  string `json:"hostname"`
	Resultdir string `json:"resultdir"`
	TimeStamp string `json:"timestamp"`
	Type      string `json:"type"`
}

type File struct {
	MD5Sum string `json:"md5sum"`
	Name   string `json:"name"`
	Data   string `json:"data"`
}

var (
	seed     string
	filesDir string
	level    uint
	logger   *log.Logger
)

func init() {
	logger = log.New("createMessage")
	flag.StringVar(&seed, "s", "", "seed file. Metadata json file")
	flag.StringVar(&filesDir, "d", "", "directory where are file.md and file.pdf")
}

func main() {
	flag.Parse()
	if flag.NFlag() != 2 {
		logger.Warn("missing arguments")
		flag.Usage()
		os.Exit(1)
	}

	// Reading seed file
	logger.Infof("Reading seed file %s", seed)
	jSeed, err := ioutil.ReadFile(seed)
	if err != nil {
		logger.Fatalf("Error reading seed file %s: %s", seed, err.Error())
	}
	logger.Debug(string(jSeed))

	// Unmarshalling seed
	var seedOperator Seed
	jUErr := json.Unmarshal(jSeed, &seedOperator)
	if jUErr != nil {
		logger.Fatalf("Error Unmarshalling seed file %s: %s", jSeed, jUErr.Error())
	}
	logger.Debug(seedOperator)

	// Add files to seed file to send it to pusher
	logger.Info("add files to seed file")
	files, err := ioutil.ReadDir(filesDir)
	if err != nil {
		logger.Fatalf("Error readir files dir %s: %s", filesDir, err.Error())
	}

	for _, file := range files {
		dataFile := filesDir + "/" + file.Name()
		logger.Info("Processing file %s", dataFile)
		data, err := ioutil.ReadFile(dataFile)
		if err != nil {
			logger.Fatalf("Error reading file %s: %s", dataFile, err.Error())
		}
		md5 := md5.New()
		io.WriteString(md5, fmt.Sprintf("%s", data))
		md5sum := fmt.Sprintf("%x", md5.Sum(nil))
		logger.Infof("md5sum: %s %s", md5sum, dataFile)

		jF := File{
			Data:   base64.StdEncoding.EncodeToString(data),
			MD5Sum: md5sum,
			Name:   file.Name(),
		}

		logger.Info("file added to json")
		seedOperator.Files = append(seedOperator.Files, jF)
	}

	logger.Infof("Writing file %s", seed)
	jSeed, err = json.Marshal(seedOperator)
	if err != nil {
		logger.Fatalf("Error json marshalling seedOperator %s", err)
	}

	err = ioutil.WriteFile(seed, jSeed, 0400)
	if err != nil {
		logger.Fatalf("Error writing file %s: %s", seed, err.Error())
	}
}
