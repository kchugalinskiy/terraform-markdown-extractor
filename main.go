package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/vharitonsky/iniflags"
)

var (
	rootDir = flag.String("dir", ".", "start dir")
	outPath = flag.String("out", "out.json", "output result filepath")
)

type Line struct {
	Name        string
	Optional    bool
	Description string
}

type Resource struct {
	Name       string
	Arguments  []Line
	Attributes []Line
}

func main() {
	iniflags.Parse()
	// log.SetLevel(log.DebugLevel)
	log.Debug("reading directory: ", *rootDir)

	files, err := ioutil.ReadDir(*rootDir)
	if err != nil {
		log.Error(err)
	}

	resources := []Resource{}

	for _, file := range files {
		err, res := parseResourse(filepath.Join(*rootDir, file.Name()))
		log.Debugf("%+v", res)

		if nil != err {
			log.Error("parsing file: ", err)
			continue
		}

		if nil == res {
			log.Info("parsing file: skipped ", file.Name())
			continue
		}

		resources = append(resources, *res)
	}

	jsonResult, err := json.Marshal(resources)
	if nil != err {
		log.Error("convert to json: ", err)
		return
	}

	f, err := os.OpenFile(*outPath, os.O_WRONLY|os.O_CREATE, 0755)

	if err != nil {
		log.Error("open output file: ", err)
		return
	}
	writer := bufio.NewWriter(f)

	if nil != writer {
		writer.Write(jsonResult)
		writer.Flush()
	}
}

func parseMatchLine(words [][]byte) Line {
	result := Line{Name: "", Optional: true, Description: ""}
	if len(words) >= 4 {
		result.Name = string(words[1])
		result.Optional = string(words[2]) == "Optional"
		result.Description = string(words[3])
	}

	return result
}

func parseResourse(path string) (error, *Resource) {
	matched, err := regexp.MatchString(".*.(markdown|html.md)", path)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("resource parsing: ", err), nil
	}

	if !matched {
		return nil, nil
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("resource parsing: ", err), nil
	}

	name := filepath.Base(path)
	re := regexp.MustCompile("[a-zA-Z_]*")
	resourceName := "aws_" + re.FindString(name)

	result := &Resource{Name: resourceName, Arguments: nil, Attributes: nil}

	argsRegex := regexp.MustCompile("## Argument Reference")
	attribRegex := regexp.MustCompile("## Attributes Reference")

	argsLoc := argsRegex.FindIndex(bytes)
	argumentsStart := 0
	argumentsEnd := 0
	if nil != argsLoc {
		argumentsStart = argsLoc[1]
	}

	attributesStart := 0
	attributesLoc := attribRegex.FindIndex(bytes[argumentsStart:])
	if nil != attributesLoc {
		argumentsEnd = attributesLoc[0] + argumentsStart
		attributesStart = attributesLoc[1] + argumentsStart
	}

	var argumentsBytes []byte
	if argumentsStart <= argumentsEnd {
		argumentsBytes = bytes[argumentsStart:argumentsEnd]
	} else {
		argumentsBytes = make([]byte, 0)
	}

	attributesBytes := bytes[attributesStart:]

	// http://regexr.com
	singleLineRegex := regexp.MustCompile("\\* `([a-zA-Z_0-9]*)` -? ?(\\([a-zA-Z]*\\)|) ?([-\\(\\)A-Za-z0-9 \\.\\,\\[\\]\\:\\/`']*)")

	attributesMatched := singleLineRegex.FindAllSubmatch(attributesBytes, -1)
	result.Attributes = make([]Line, len(attributesMatched))
	for index, attributeParsed := range attributesMatched {
		line := parseMatchLine(attributeParsed)
		result.Attributes[index] = line
	}

	argumentsMatched := singleLineRegex.FindAllSubmatch(argumentsBytes, -1)
	result.Arguments = make([]Line, len(argumentsMatched))
	for index, argumentMatched := range argumentsMatched {
		line := parseMatchLine(argumentMatched)
		result.Arguments[index] = line
	}

	return nil, result
}
