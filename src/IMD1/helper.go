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
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// =====================================
// Pairs

type Pair struct {
	i int
	j int
}

func (b Pair) String() string {
	return fmt.Sprintf("(%v, %v)", b.i, b.j)
}

// =====================================
// Setup log

func SetupLog(lvl log.Level) {
    customFormatter := new(log.TextFormatter)
    customFormatter.TimestampFormat = "2006-01-02 15:04:05"
    customFormatter.FullTimestamp = true
    log.SetFormatter(customFormatter)
	log.SetLevel(lvl)
}

// =====================================
// String manipulation

const INF_BLANKS int = 16

func RemoveExcessSpaces(s string) string {
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(s, " "))
}

// - - - - -

func CheckRunesEndWithUnescapedASCII(r []rune, ending string) bool {
	if len(r) < len(ending) {
		return false
	}
	for ei,ri := (len(ending)-1), (len(r)-1); ei >= 0; ei,ri = ei-1, ri-1 {
		if r[ri] != rune(ending[ei]) {
			return false
		}
	}
	not_escaped := true
	for i := len(r) - len(ending) - 1; i >= 0 && r[i] == '\\'; i-- {
		not_escaped = !not_escaped
	}
	return not_escaped
}

// - - - - -

func CheckRunesStartsWithASCII(r []rune, prefix string) bool {
	if len(r) < len(prefix) {
		return false
	}
	for ei,ri := 0, 0; ei < len(prefix); ei,ri = ei+1, ri+1 {
		if r[ri] != rune(prefix[ei]) {
			return false
		}
	}
	return true
}

// - - - - -

func CountCharSpaces(c rune) uint16 {
	if c == '\t' {
		return 4
	} else if (c == ' ') {
		return 1
	}
	return 0
}

func CountPrefixSpaces(s string) (uint16, uint16) {
	var simple, total uint16
	for _, b := range s {
		add := CountCharSpaces(b)
		if add == 0 {
			return simple, total
		}
		total += add
		simple++
	}
	return simple, total
}

// - - - - -

func Contains[T comparable](s []T, e T) bool {
    for _, v := range s {
        if v == e {
            return true
        }
    }
    return false
}

// - - - - -

func GatherBlockOptions(line *LineStruct, pool []string) map[string]string {
	left := line.RuneJ
	result := make(map[string]string)
	for left < len(line.RuneContent) && line.RuneContent[left] == '[' {
		right := left + 1
		equal_position := left - 1
		for right < len(line.RuneContent) && line.RuneContent[right] != ']' {
			if line.RuneContent[right] == '\\' {
				right++
			} else if line.RuneContent[right] == '=' && equal_position == left-1 {
				equal_position = right
			}
			right++
		}
		if right >= len(line.RuneContent) {
			break
		}
		if equal_position == left-1 {
			// TODO - trigger warning (equal sign not found)
			break
		}
		option := string(line.RuneContent[left+1:equal_position])
		if !Contains(pool, option) {
			// TODO - trigger warning (option not found in pool of options)
		} else {
			result[option] = string(line.RuneContent[equal_position+1:right])
		}
		left = right + 1
	}
	line.RuneJ = left
	return result
}

func StringToHTMLSafe(s string) string {
	var sb strings.Builder
	for _, c := range s {
		switch c {
		case '<':
			sb.WriteString("&lt;")
		case '>':
			sb.WriteString("&gt;")
		case '$':
			sb.WriteString("&#36;")
		default:
			sb.WriteRune(c)
		}
	}
	return sb.String()
}
