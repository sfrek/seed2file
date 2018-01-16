package main

import (
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
)

type SeedFile struct {
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
	jsonSeed        string
	jsonDest        string
	encodedFilePath string
)

func init() {
	flag.StringVar(&jsonSeed, "js", "", "Json seed to add the content file")
	flag.StringVar(&jsonDest, "jd", "", "Json destine file")
	flag.StringVar(&encodedFilePath, "ef", "", "encoded file path")
}

// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	flag.Parse()
	fmt.Println(flag.CommandLine.NFlag())

	F1 := File{
		Name:   "F1.name",
		MD5Sum: "F1.md5sum",
		Data:   "F1.data",
	}

	F2 := File{
		Name:   "F2.name",
		MD5Sum: "F2.md5sum",
		Data:   "F2.data",
	}

	S1 := SeedFile{
		Type:      "S1.periodic",
		TimeStamp: "S1.periodic",
		Files:     []File{F1, F2},
	}

	// Write file
	dj, _ := json.Marshal(&S1)
	fmt.Printf("%s\n", dj)
	wErr := ioutil.WriteFile("temporal.json", dj, 0644)
	check(wErr)

	// Read File
	jSeed, err := ioutil.ReadFile(jsonSeed)
	check(err)
	fmt.Print(string(jSeed))

	var seedOperator SeedFile
	jUErr := json.Unmarshal(jSeed, &seedOperator)
	check(jUErr)
	fmt.Println(seedOperator)

	seedOperator.Files = append(seedOperator.Files, F2)
	fmt.Println(seedOperator)

	// md5Sum
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%s", dj))
	fmt.Printf("temporal.json: %x\n", h.Sum(nil))

	high, err := ioutil.ReadFile(encodedFilePath)
	check(err)
	z := md5.New()
	io.WriteString(z, fmt.Sprintf("%s", high))
	F3 := File{
		Name:   "F3.name",
		MD5Sum: fmt.Sprintf("%x", z.Sum(nil)),
		Data:   fmt.Sprintf("%s", high),
	}

	seedOperator.Files = append(seedOperator.Files, F3)
	fmt.Println(seedOperator)
}
