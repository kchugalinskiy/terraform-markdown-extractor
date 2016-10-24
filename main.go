package main

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/vharitonsky/iniflags"
)

var (
	rootDir = flag.String("dir", ".", "start dir")
)

func main() {
	iniflags.Parse()
	log.SetLevel(log.DebugLevel)
	log.Debug("reading directory: ", *rootDir)

	files, err := ioutil.ReadDir(*rootDir)
	if err != nil {
		log.Error(err)
	}

	for _, file := range files {
		parseFile(filepath.Join(*rootDir, file.Name()))
	}

}

func parseFile(path string) {
	matched, err := regexp.MatchString(".*.markdown", path)
	if err != nil {
		log.Error(err)
		return
	}

	if !matched {
		return
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error(err)
		return
	}

	name := filepath.Base(path)
	re := regexp.MustCompile("[a-zA-Z_]*")
	resourceName := re.FindString(name)

	log.Debug("res name: ", resourceName)

	argsRegex := regexp.MustCompile("The following arguments are supported:")
	attribRegex := regexp.MustCompile("The following attributes are exported:")

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

	argumentsBytes := bytes[argumentsStart:argumentsEnd]
	attributesBytes := bytes[attributesStart:]

	singleLineRegex := regexp.MustCompile("\\* `([a-zA-Z_0-9]*)` -? ?(\\([a-zA-Z]*\\)|) ?([-\\(\\)A-Za-z0-9 \\.\\,\\[\\]\\:\\/`']*)")

	attributesMatched := singleLineRegex.FindAllSubmatch(attributesBytes, -1)
	log.Info(attributesMatched)

	argumentsMatched := singleLineRegex.FindAllSubmatch(argumentsBytes, -1)
	log.Info(argumentsMatched)
}
