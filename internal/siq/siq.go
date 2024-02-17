package siq

import (
	"archive/zip"
	"encoding/xml"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
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
func CreateSIQPackage() error {
	sourceDir := "./downloads" // Adjust the path as necessary
	targetZipFile := "./quiz_package.siq"
	//processedDir := "./processed"

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

	// Assume a single theme and round for simplicity
	pkg := Package{Name: "Quiz Package", Rounds: []Round{{Name: "Round 1", Themes: []Theme{{Name: "♫"}}}}}

	for _, file := range files {
		filename := file.Name()
		filePath := filepath.Join(sourceDir, filename)

		if !file.IsDir() && strings.HasSuffix(filename, ".mp3") {
			// Directly use URL-encoded filename, without adding "Audio/" prefix
			encodedFilename := url.PathEscape(filename)

			if err := addFileToZip(archive, filePath, encodedFilename); err != nil {
				log.Printf("Failed to add file %s to ZIP: %v", filename, err)
				continue
			}

			// Constructing question using encoded file name directly, without prefix
			question := Question{
				Price: 1,
				Scenario: []Atom{{
					Type: "voice",
					Text: "@" + encodedFilename, // Directly use encoded filename
				}},
				Right: []string{strings.TrimSuffix(filename, filepath.Ext(filename))},
			}
			pkg.Rounds[0].Themes[0].Questions = append(pkg.Rounds[0].Themes[0].Questions, question)
		}
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
