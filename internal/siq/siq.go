package siq

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Package struct {
	XMLName xml.Name `xml:"package"`
	Name    string   `xml:"name,attr"`
	Rounds  []Round  `xml:"rounds>round"`
}

type Round struct {
	Name   string  `xml:"name,attr"`
	Themes []Theme `xml:"themes>theme"`
}

type Theme struct {
	Name      string     `xml:"name,attr"`
	Questions []Question `xml:"questions>question"`
}

type Question struct {
	Price    int      `xml:"price,attr"`
	Scenario []Atom   `xml:"scenario>atom"`
	Right    []string `xml:"right>answer"`
}

type Atom struct {
	Type string `xml:"type,attr"`
	Text string `xml:",chardata"`
}

// Package structure and other types definitions as before

// Adds a file to the ZIP archive and returns the path used inside the ZIP
func addFileToZip(zipWriter *zip.Writer, filePath, baseInZip string) error {
	fileToZip, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Prefix the file name with "Audio/" to place it into an "Audio" directory within the ZIP
	header.Name = filepath.Join("Audio", baseInZip)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

// CreateSIQPackage function now has hardcoded paths
func CreateSIQPackage(siqName string) error {
	sourceDir := "./downloads" // Adjust the path as necessary
	targetZipFile := fmt.Sprintf("./%s.siq", siqName)

	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	zipFile, err := os.Create(targetZipFile)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	pkg := Package{Name: siqName, Rounds: []Round{{Name: "Round 1"}}}
	var questions []Question

	for _, file := range files {
		filename := file.Name()
		filePath := filepath.Join(sourceDir, filename)

		if !file.IsDir() && strings.HasSuffix(filename, ".mp3") {
			encodedFilename := url.PathEscape(filename)
			if err := addFileToZip(archive, filePath, encodedFilename); err != nil {
				log.Printf("Failed to add file %s to ZIP: %v", filename, err)
				continue
			}

			question := Question{
				Price:    1,
				Scenario: []Atom{{Type: "voice", Text: "@" + encodedFilename}},
				Right:    []string{sanitizeRightAnswer(strings.TrimSuffix(filename, filepath.Ext(filename)))},
			}
			questions = append(questions, question)
		}
	}

	// Split questions into groups of 10 for each theme
	for i := 0; i < len(questions); i += 10 {
		end := i + 10
		if end > len(questions) {
			end = len(questions)
		}
		themeName := fmt.Sprintf("%d", i/10+1)
		pkg.Rounds[0].Themes = append(pkg.Rounds[0].Themes, Theme{Name: themeName, Questions: questions[i:end]})
	}

	contentFile, err := archive.Create("content.xml")
	if err != nil {
		return err
	}

	enc := xml.NewEncoder(contentFile)
	enc.Indent("", "  ")
	if err := enc.Encode(pkg); err != nil {
		return err
	}

	return nil
}

func sanitizeRightAnswer(filename string) string {
	// Create a regular expression to find the phrases.
	// This pattern looks for the phrases "official video", "lyrics video", or "lyrics", case-insensitive.
	re := regexp.MustCompile(`\((?i).*?(official|lyrics|lyric video|official video).*?\)|\[(?i).*?(official|lyrics|lyric|video).*?\]`)

	// Replace occurrences of the pattern with an empty string.
	cleanedFilename := re.ReplaceAllString(filename, "")

	// Trim any extra spaces that may have been left as a result of the replacement.
	cleanedFilename = strings.TrimSpace(cleanedFilename)

	return cleanedFilename
}
