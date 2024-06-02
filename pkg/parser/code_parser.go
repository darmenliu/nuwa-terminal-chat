package parser

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
)

type SourceFile struct {
	FileName     string
	FileContent  string
	MatchContent string
}

// Parse filename and code from the match content
func (s *SourceFile) ParseFileName() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	re := regexp.MustCompile(`@([a-zA-Z0-9_./]+)@`)
	match := re.FindStringSubmatch(s.MatchContent)

	if match != nil {
		s.FileName = match[1]
		logger.Info("Matched:", "match_content ", match[1])
	} else {
		logger.Info("No match found.")
	}
}

func (s *SourceFile) ParseFileContent() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	//```[^\n]*\n([\s\S]*?\n)```
	regstr := "```[^\n]*\n([" + `\s\S` + "]*?)\n```"
	re := regexp.MustCompile(regstr)
	match := re.FindStringSubmatch(s.MatchContent)

	if match != nil {
		s.FileContent = match[1]
		logger.Info("Matched:", "match_content ", match[1])
	} else {
		logger.Info("No match found.")
	}
}

type SourceFileDict struct {
	SourceFiles map[string]SourceFile
}

func (s *SourceFileDict) AddSourceFile(fileName string, fileContent string) {
	s.SourceFiles[fileName] = SourceFile{FileName: fileName, FileContent: fileContent}
}

func (s *SourceFileDict) GetSourceFile(fileName string) (SourceFile, error) {
	file, ok := s.SourceFiles[fileName]
	if !ok {
		return SourceFile{}, fmt.Errorf("file not found")
	}
	return file, nil
}

// PrintSourceFiles function prints the source files in the SourceFileDict
func (s *SourceFileDict) PrintSourceFiles() {
	for key, value := range s.SourceFiles {
		fmt.Println("Key:", key, "Value:", value)
	}
}

func NewSourceFileDict() *SourceFileDict {
	return &SourceFileDict{SourceFiles: make(map[string]SourceFile)}
}

type CodeParser interface {
	ParseCode(text string) (SourceFileDict, error)
}

type GoCodeParser struct {
}

func NewGoCodeParser() *GoCodeParser {
	return &GoCodeParser{}
}

// ParseCode function Parse the code from markdown blocks and return a SourceFileDict
func (g *GoCodeParser) ParseCode(text string) ([]SourceFile, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	var sources []SourceFile

	// Regex to match code blocks

	//```[a-z]*\n[\s\S]*?\n```
	//[a-zA-Z0-9_.]+\n```[^\n]*\n[\s\S]*?\n```
	// @[a-zA-Z0-9_/.]+@\n```[^\n]*\n[\s\S]*?\n```
	regstr := "@[a-zA-Z0-9_/.]+@\n```[^\n]*\n[" + `\s\S` + "]*?\n```"

	regex := regexp.MustCompile(regstr)

	// Find all matches
	matches := regex.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		// Get filename and content
		matchContent := match[0]

		logger.Info("Adding file to source file dict", "content", matchContent)
		// Add to map
		sources = append(sources, SourceFile{FileName: "", FileContent: "", MatchContent: matchContent})
	}

	return sources, nil
}
