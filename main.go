package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, os.ModePerm)
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func downloadLatestPostalCodes() {
	url := "http://download.geonames.org/export/zip/US.zip"

	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	output, err := os.Create(fileName)

	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
}

func processFile(fileName string) {
	f, err := os.Open(fileName)
	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)

	buffer := bytes.NewBufferString("")

	buffer.WriteString(fmt.Sprintln("INSERT INTO postal_code( " +
		"country_code, postal_code, state, state_code, county, county_code," +
		"latitude, longitude) VALUES"))

	processedPostalCodes := make(map[string]bool)

	for scanner.Scan() {
		text := scanner.Text()
		parts := strings.Split(text, "\t")

		country_code := parts[0]
		postal_code := parts[1]
		state := parts[3]
		state_code := parts[4]
		county := strings.Replace(parts[5], "'", "''", -1)
		county_code := parts[6]
		latitude := parts[9]
		longitude := parts[10]
		if !processedPostalCodes[postal_code] {
			processedPostalCodes[postal_code] = true
			buffer.WriteString(fmt.Sprintf("( '%s', "+ //country code
				"'%s', "+ //postal_code
				"'%s', "+ //state
				"'%s', "+ //state_code
				"'%s', "+ //county
				"'%s', "+ //county_code
				"%s, "+ //latitude
				"%s "+ //longitude
				"),\n", country_code, postal_code, state, state_code, county, county_code, latitude, longitude))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	buffer.Truncate(buffer.Len() - 2)
	buffer.WriteString(";\n")
	fmt.Println(buffer.String())
}

func main() {
	downloadLatestPostalCodes()
	unzip("US.zip", "download")
	processFile("download/US.txt")
}
