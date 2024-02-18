/*
  This file is part of the IMD1 project.
  Copyright (c) 2024 Valentin-Ioan Vintilă

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, version 3.

  This program is distributed in the hope that it will be useful, but
  WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
  General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package main

import "C"
import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// =====================================
// Markdown to HTML (Go)

func IMD1_MDToHTMLHelper(file FileStruct) string {
	start_time := time.Now()
	tree, metadata := file.MDParse()
	end_time := time.Now()
	log.Infof("Parsing took %v", end_time.Sub(start_time))
	log.Debug(tree)
	log.Debug(metadata)

	start_time = time.Now()
	html := GenerateHTML(&tree)
	end_time = time.Now()
	log.Infof("Generating the HTML took %v", end_time.Sub(start_time))
	return html
}

func IMD1_MDFileToHTMLFile(md_filename string, html_filename string) {
	var file FileStruct
	file.ReadFile(md_filename)

	html := IMD1_MDToHTMLHelper(file)

	fout, err := os.Create(html_filename)
	if err != nil {
		panic(err)
	}
	fout.WriteString(html)
	fout.Close()
}

func IMD1_MDToHTMLFile(s string, html_filename string) {
	var file FileStruct
	file.ReadString(s)

	html := IMD1_MDToHTMLHelper(file)

	fout, err := os.Create(html_filename)
	if err != nil {
		panic(err)
	}
	fout.WriteString(html)
	fout.Close()
}

func IMD1_MDFileToHTML(md_filename string) string {
	var file FileStruct
	file.ReadFile(md_filename)

	return IMD1_MDToHTMLHelper(file)
}

func IMD1_MDToHTML(s string) string {
	var file FileStruct
	file.ReadString(s)
	
	return IMD1_MDToHTMLHelper(file)
}

// =====================================
// Markdown to HTML (C-exported variants)

//export C_IMD1_MDFileToHTMLFile
func C_IMD1_MDFileToHTMLFile(c_md_filename *C.char, c_html_filename *C.char) {
	md_filename := C.GoString(c_md_filename)
	html_filename := C.GoString(c_html_filename)
	IMD1_MDFileToHTMLFile(md_filename, html_filename)
}

//export C_IMD1_MDToHTMLFile
func C_IMD1_MDToHTMLFile(c_s *C.char, c_html_filename *C.char) {
	s := C.GoString(c_s)
	html_filename := C.GoString(c_html_filename)
	IMD1_MDToHTMLFile(s, html_filename)
}

//export C_IMD1_MDFileToHTML
func C_IMD1_MDFileToHTML(c_md_filename *C.char) *C.char {
	md_filename := C.GoString(c_md_filename)
	html := IMD1_MDFileToHTML(md_filename)
	return C.CString(html)
}

//export C_IMD1_MDToHTML
func C_IMD1_MDToHTML(c_s *C.char) *C.char {
	s := C.GoString(c_s)
	html := IMD1_MDToHTML(s)
	return C.CString(html)
}
