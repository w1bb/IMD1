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
	log "github.com/sirupsen/logrus"
	"github.com/w1bb/IMD1/src/IMD1"
	"os"
)

// =====================================
// Main itself, for sunny days

func main() {
	IMD1.SetupLog(log.DebugLevel /*log.InfoLevel*/)
	imd1, err := os.ReadFile("test.md")
	if err != nil {
		panic(err)
	}
	html, _ := IMD1.ToHTML(string(imd1))

	fout, err := os.Create("test.html")
	if err != nil {
		panic(err)
	}
	_, _ = fout.WriteString(html)
	_ = fout.Close()
}
