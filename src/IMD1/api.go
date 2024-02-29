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
	"time"
	"unsafe"

	log "github.com/sirupsen/logrus"
)

// =====================================
// Unsafe pointer (void*) free (only C)

//export CFree
func CFree(p unsafe.Pointer) {
	C.free(p)
}

// =====================================
// Generic parse (only Go)

func ParseIMD1(s string) (Tree[BlockInterface], MDMetaStructure) {
	var file FileStruct
	file.ReadString(s)
	tree, metadata := file.Parse()
	return tree, metadata
}

// =====================================
// Markdown to HTML (Go & C)

func ToHTML(s string) (string, MDMetaStructure) {
	var file FileStruct
	file.ReadString(s)

	startTime := time.Now()
	tree, metadata := file.Parse()
	endTime := time.Now()

	log.Infof("Parsing took %v", endTime.Sub(startTime))
	log.Debug(tree)
	log.Debug(metadata)

	startTime = time.Now()
	html := GenerateHTML(&tree, BlockDocumentTypeCompleteSpecification)
	endTime = time.Now()

	log.Infof("Generating the HTML took %v", endTime.Sub(startTime))
	return html, metadata
}

//export CToHTML
func CToHTML(cS *C.char) (*C.char, unsafe.Pointer) {
	html, serialMetadata := ToHTML(C.GoString(cS))
	return C.CString(html), C.CBytes(serialMetadata.Serialize())
}

// =====================================
// Markdown to LaTeX (Go & C)

func ToLaTeX(s string) (string, MDMetaStructure) {
	var file FileStruct
	file.ReadString(s)

	startTime := time.Now()
	tree, metadata := file.Parse()
	endTime := time.Now()

	log.Infof("Parsing took %v", endTime.Sub(startTime))
	log.Debug(tree)
	log.Debug(metadata)

	startTime = time.Now()
	latex := GenerateLaTeX(&tree)
	endTime = time.Now()

	log.Infof("Generating the LaTeX took %v", endTime.Sub(startTime))
	return latex, metadata
}

//export CToLaTeX
func CToLaTeX(cS *C.char) (*C.char, unsafe.Pointer) {
	html, serialMetadata := ToLaTeX(C.GoString(cS))
	return C.CString(html), C.CBytes(serialMetadata.Serialize())
}
