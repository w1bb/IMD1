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

package IMD1

//#include <stdlib.h>
import "C"
import (
	"os"
	"time"
	"unsafe"

	log "github.com/sirupsen/logrus"
)

// =====================================
// Unsafe pointer (void*) free (C only)

//export C_FreeUnsafePointer
func C_FreeUnsafePointer(p unsafe.Pointer) {
	C.free(p)
}

// =====================================
// Markdown to HTML (Go)

func ParseString(s string) (Tree[BlockInterface], MDMetaStructure) {
	var file FileStruct
	file.ReadString(s)
	tree, metadata := file.MDParse()
	return tree, metadata
}

func IMD1_MDToHTMLHelper(file FileStruct) (string, []byte) {
	start_time := time.Now()
	tree, metadata := file.MDParse()
	end_time := time.Now()
	log.Infof("Parsing took %v", end_time.Sub(start_time))
	log.Debug(tree)
	log.Debug(metadata)

	start_time = time.Now()
	html := GenerateHTML(&tree, BlockDocumentType_CompleteSpecification)
	end_time = time.Now()
	log.Infof("Generating the HTML took %v", end_time.Sub(start_time))
	return html, metadata.Serialize()
}

func IMD1_MDFileToHTMLFile(md_filename string, html_filename string) []byte {
	var file FileStruct
	file.ReadFile(md_filename)

	html, serial_metadata := IMD1_MDToHTMLHelper(file)

	fout, err := os.Create(html_filename)
	if err != nil {
		panic(err)
	}
	fout.WriteString(html)
	fout.Close()

	return serial_metadata
}

func IMD1_MDToHTMLFile(s string, html_filename string) []byte {
	var file FileStruct
	file.ReadString(s)

	html, serial_metadata := IMD1_MDToHTMLHelper(file)

	fout, err := os.Create(html_filename)
	if err != nil {
		panic(err)
	}
	fout.WriteString(html)
	fout.Close()

	return serial_metadata
}

func IMD1_MDFileToHTML(md_filename string) (string,  []byte) {
	var file FileStruct
	file.ReadFile(md_filename)

	return IMD1_MDToHTMLHelper(file)
}

func IMD1_MDToHTML(s string) (string,  []byte) {
	var file FileStruct
	file.ReadString(s)

	return IMD1_MDToHTMLHelper(file)
}

// =====================================
// Markdown to HTML (C-exported variants)

// //export C_IMD1_MDFileToHTMLFile
// func C_IMD1_MDFileToHTMLFile(c_md_filename *C.char, c_html_filename *C.char) unsafe.Pointer {
// 	md_filename := C.GoString(c_md_filename)
// 	html_filename := C.GoString(c_html_filename)
// 	serial_metadata := IMD1_MDFileToHTMLFile(md_filename, html_filename)
// 	return C.CBytes(serial_metadata)
// }

// //export C_IMD1_MDToHTMLFile
// func C_IMD1_MDToHTMLFile(c_s *C.char, c_html_filename *C.char) unsafe.Pointer {
// 	s := C.GoString(c_s)
// 	html_filename := C.GoString(c_html_filename)
// 	serial_metadata := IMD1_MDToHTMLFile(s, html_filename)
// 	return C.CBytes(serial_metadata)
// }

// //export C_IMD1_MDFileToHTML
// func C_IMD1_MDFileToHTML(c_md_filename *C.char) (*C.char, unsafe.Pointer) {
// 	md_filename := C.GoString(c_md_filename)
// 	html, serial_metadata := IMD1_MDFileToHTML(md_filename)
// 	return C.CString(html), C.CBytes(serial_metadata)
// }

//export C_IMD1_MDToHTML
func C_IMD1_MDToHTML(c_s *C.char) (*C.char, unsafe.Pointer) {
	s := C.GoString(c_s)
	html, serial_metadata := IMD1_MDToHTML(s)
	return C.CString(html), C.CBytes(serial_metadata)
}

// =====================================
// Markdown to LaTeX (Go)

func IMD1_MDToLaTeXHelper(file FileStruct) string {
	start_time := time.Now()
	tree, metadata := file.MDParse()
	end_time := time.Now()
	log.Infof("Parsing took %v", end_time.Sub(start_time))
	log.Debug(tree)
	log.Debug(metadata)

	start_time = time.Now()
	latex := GenerateLaTeX(&tree)
	end_time = time.Now()
	log.Infof("Generating the LaTeX took %v", end_time.Sub(start_time))
	return latex
}

func IMD1_MDFileToLaTeXFile(md_filename string, latex_filename string) {
	var file FileStruct
	file.ReadFile(md_filename)

	latex := IMD1_MDToLaTeXHelper(file)

	fout, err := os.Create(latex_filename)
	if err != nil {
		panic(err)
	}
	fout.WriteString(latex)
	fout.Close()
}

func IMD1_MDToLaTeXFile(s string, latex_filename string) {
	var file FileStruct
	file.ReadString(s)

	latex := IMD1_MDToLaTeXHelper(file)

	fout, err := os.Create(latex_filename)
	if err != nil {
		panic(err)
	}
	fout.WriteString(latex)
	fout.Close()
}

func IMD1_MDFileToLaTeX(md_filename string) string {
	var file FileStruct
	file.ReadFile(md_filename)

	return IMD1_MDToLaTeXHelper(file)
}

func IMD1_MDToLaTeX(s string) string {
	var file FileStruct
	file.ReadString(s)

	return IMD1_MDToLaTeXHelper(file)
}

// =====================================
// Markdown to LaTeX (C-exported variants)

//export C_IMD1_MDFileToLaTeXFile
func C_IMD1_MDFileToLaTeXFile(c_md_filename *C.char, c_latex_filename *C.char) {
	md_filename := C.GoString(c_md_filename)
	latex_filename := C.GoString(c_latex_filename)
	IMD1_MDFileToLaTeXFile(md_filename, latex_filename)
}

//export C_IMD1_MDToLaTeXFile
func C_IMD1_MDToLaTeXFile(c_s *C.char, c_latex_filename *C.char) {
	s := C.GoString(c_s)
	latex_filename := C.GoString(c_latex_filename)
	IMD1_MDToLaTeXFile(s, latex_filename)
}

//export C_IMD1_MDFileToLaTeX
func C_IMD1_MDFileToLaTeX(c_md_filename *C.char) *C.char {
	md_filename := C.GoString(c_md_filename)
	latex := IMD1_MDFileToLaTeX(md_filename)
	return C.CString(latex)
}

//export C_IMD1_MDToLaTeX
func C_IMD1_MDToLaTeX(c_s *C.char) *C.char {
	s := C.GoString(c_s)
	latex := IMD1_MDToLaTeX(s)
	return C.CString(latex)
}
