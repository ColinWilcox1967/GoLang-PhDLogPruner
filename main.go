package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"

	"flag"
)

const (
	log_pruner_version = "1.1"
	default_output_file = "./OUTPUT.TXT"
	default_keyword_file = "./KEYWORDS.TXT"
	default_log_file = "./LOGFILE.TXT"
)

var (
	outputFile string
	keywordFile string
	logFile string
	eraseOutput bool
	silentRunning bool
	showLineNumbers bool
	caseSensitive bool

	logLines []string
	keywords []string
)

func showBanner() {
	fmt.Printf("Log Pruner (Version %s).\n", log_pruner_version)
	fmt.Println("(c) Colin Wilcox 2024.")
}

func parseCommandLine() {

	flag.StringVar(&keywordFile, "keywords", default_keyword_file, "Name of keyword file.")
	flag.StringVar(&outputFile, "output", default_output_file, "Name of output results file.")
	flag.StringVar(&logFile, "log", default_log_file, "Name of log file to process.")
	flag.BoolVar(&eraseOutput, "erase", false, "Erase output file before writing.")
	flag.BoolVar(&silentRunning, "silent", false, "Run the utility with no echo to console.")
	flag.BoolVar(&showLineNumbers, "linenumbers", true, "Display line numbers of those lines containing a text match.")
	flag.BoolVar(&caseSensitive, "case", false, "Specifies whether the keyword matching be case sensitive.")
	flag.Parse()

	if len(keywordFile) == 0 {
		keywordFile = default_keyword_file
	}

	if !silentRunning {
		fmt.Printf("\nLog            : '%s'.\n", logFile)
		fmt.Printf("Output file    : '%s'.\n", outputFile)
		fmt.Printf("Keyword file   : '%s'.\n", keywordFile)
		fmt.Printf("Erase file     : %v.\n", eraseOutput)
		fmt.Printf("Line numbers   : %v.\n", showLineNumbers)
		fmt.Printf("Case sensitive : %v.\n\n", caseSensitive)
	}
}

func readFile(filePath string) (error, []string) {

	data := make([]string, 0)

	file, err := os.Open(filePath)
    if err != nil {
       return err, []string{}
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
   
    for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			data = append(data, line)
		}
    }

    if err = scanner.Err(); err != nil {
       return err, []string{}
    }

	return nil, data
}

func readLog() (error, []string) {
	return readFile(logFile)
}

func readKeywordFile() (error, []string) {
	return readFile(keywordFile)
}

func writeLineToOutputFile(line string) bool {
		
	f, err := os.OpenFile(outputFile, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0600)
	defer f.Close()

	if err != nil {
		fmt.Println(err)
		return false
	}

	_, err = f.WriteString(line)
	f.WriteString("\n")
	if err != nil {
		fmt.Println(err)
  		return false
	}
	
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func scanLinesForKeywords(keywords []string) int {
	matchingLineCount:=0

	for lineNumber, line := range(logLines) {
		for _, word := range(keywords) {
			matchFound := false
			if caseSensitive {
				matchFound = strings.Contains(line, word)
			} else {
				matchFound = strings.Contains(strings.ToLower(line), strings.ToLower(word))
			}

			if matchFound {

				if !silentRunning {
					str := ""
					if showLineNumbers {
						str += fmt.Sprintf("%04d:%s\n", lineNumber, line)
					} else {
						str += fmt.Sprintf("%s\n", line)
					}
					fmt.Printf(str)
				}

				writeLineToOutputFile(line)
				matchingLineCount++
				
			}
		}
	}

	if !silentRunning {
		fmt.Printf("\nMatching line count = %d.\n", matchingLineCount)
	}

	return matchingLineCount
}


func main () {
	showBanner()

	parseCommandLine()

	if eraseOutput {
		os.Remove(outputFile)
	}

	err, keywords := readKeywordFile()
	if err != nil {
		fmt.Printf("*** Error : Problem reading keywords (%v).\n", err)
		os.Exit(-2)
	}

	err, logLines = readLog()

	if err != nil {
		fmt.Printf("*** Error : Problem reading logfile (%v).\n", err)
		os.Exit(-1)
	}
	
	scanLinesForKeywords(keywords)
}