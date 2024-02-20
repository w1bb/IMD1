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

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// =====================================
// Line structure

type LineStruct struct {
	LineIndex int
	Indentation uint16

	RuneContent []rune
	RuneJ int
}

func (line LineStruct) String() string {
	return fmt.Sprintf(
		"LineStruct (line-index=%v, indent=%v, rune-j=%v) :: %v",
		line.LineIndex,
		line.Indentation,
		line.RuneJ,
		line.RuneContent,
	)
}

func (line LineStruct) Empty() bool {
	return line.RuneContent == nil || len(line.RuneContent) == 0
}

// =====================================
// File structure

type FileStruct struct {
	Lines []LineStruct
}

func (file FileStruct) String() string {
	var aux string
	for _, line := range(file.Lines) {
		aux += "  " + line.String() + "\n"
	}
	return fmt.Sprintf("FileStruct <%v lines> {\n%v}", len(file.Lines), aux)
}

func (file *FileStruct) ReadString(s string) {
	content := strings.Split(
		strings.ReplaceAll(s, "\r\n", "\n"),
		"\n")
	file.Lines = make([]LineStruct, len(content))
	for i := range(file.Lines) {
		simple, total := CountPrefixSpaces(content[i])
		file.Lines[i] = LineStruct{
			Indentation: total,
			RuneContent: []rune(content[i][simple:]),
			RuneJ: 0, // Default already
			LineIndex: i,
		}
	}
	log.Debug(file)
}

func (file *FileStruct) ReadFile(filename string) {
	s, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	file.ReadString(string(s))
}

func (file FileStruct) GetStringBetween(start Pair[int, int], end Pair[int, int]) string {
	if len(file.Lines) <= start.i {
		return ""
	}
	s := ""
	for start.i < end.i || (start.i == end.i && start.j < end.j) {
		if start.j >= len(file.Lines[start.i].RuneContent) {
			s += "\n"
			start.i++
			start.j = 0
		} else {
			s += string(file.Lines[start.i].RuneContent[start.j])
			start.j++
		}
	}
	return s
}
