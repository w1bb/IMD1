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

import (
	"encoding/binary"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

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

func RemoveExcessSpaces(s string) string {
	re := regexp.MustCompile(`\s+`)
	s = strings.TrimSpace(re.ReplaceAllString(s, " "))
	if s == " " {
		return ""
	}
	return s
}

// - - - - -

func CheckRunesEndWithUnescapedASCII(r []rune, ending string) bool {
	if len(r) < len(ending) {
		return false
	}
	for ei, ri := len(ending)-1, len(r)-1; ei >= 0; ei, ri = ei-1, ri-1 {
		if r[ri] != rune(ending[ei]) {
			return false
		}
	}
	notEscaped := true
	for i := len(r) - len(ending) - 1; i >= 0 && r[i] == '\\'; i-- {
		notEscaped = !notEscaped
	}
	return notEscaped
}

// - - - - -

func CheckRunesStartsWithASCII(r []rune, prefix string) bool {
	if len(r) < len(prefix) {
		return false
	}
	for ei, ri := 0, 0; ei < len(prefix); ei, ri = ei+1, ri+1 {
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
	} else if c == ' ' {
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
		equalPosition := left - 1
		for right < len(line.RuneContent) && line.RuneContent[right] != ']' {
			if line.RuneContent[right] == '\\' {
				right++
			} else if line.RuneContent[right] == '=' && equalPosition == left-1 {
				equalPosition = right
			}
			right++
		}
		if right >= len(line.RuneContent) {
			break
		}
		if equalPosition == left-1 {
			log.Warnf(
				"Non-option %v detected while searching for options. If you intended to write an option, don't forget the equal sign. If this is not an option, consider writing this on a separate line from the IMD1 tag. The search for other options has halted.",
				string(line.RuneContent[left:right+1]),
			)
			break
		}
		option := string(line.RuneContent[left+1 : equalPosition])
		if !Contains(pool, option) {
			log.Warnf(
				"Option %v (=> \"%v\") will be ignored (it is not part of the IMD1 specification)",
				string(line.RuneContent[left:right+1]),
				option,
			)
		} else {
			result[option] = string(line.RuneContent[equalPosition+1 : right])
		}
		left = right + 1
	}
	line.RuneJ = left
	return result
}

func StringToBool(s string) bool {
	return Contains([]string{"allow", "allowed", "1", "true", "ok", "yes"}, strings.ToLower(s))
}

func StringToHTMLSafe(s string) string {
	var sb strings.Builder
	for _, c := range s {
		switch c {
		case '<':
			sb.WriteString("&lt;")
		case '>':
			sb.WriteString("&gt;")
		case '$': // Math should be handled by now
			sb.WriteString("&#36;")
		case '`': // Code should be handled by now
			sb.WriteString("&#96;")
		case '^':
			sb.WriteString("&#94;")
		case '"':
			sb.WriteString("&quot;")
		case '\'':
			sb.WriteString("&apos;")
		case '~':
			sb.WriteString("&#126;")
		case '{':
			sb.WriteString("&#123;")
		case '}':
			sb.WriteString("&#125;")
		case '|':
			sb.WriteString("&#124;")
		case '/':
			sb.WriteString("&#47;")
		case '\\':
			sb.WriteString("&#92;")
		default:
			sb.WriteRune(c)
		}
	}
	return sb.String()
}

func StringToLaTeXSafe(s string) string {
	var sb strings.Builder
	for _, c := range s {
		switch c {
		case '\\':
			sb.WriteString("\\\\")
		case '{':
			sb.WriteString("\\{")
		case '}':
			sb.WriteString("\\}")
		default:
			sb.WriteRune(c)
		}
	}
	return sb.String()
}

func StringToKaTeXSafe(s string) string {
	var sb strings.Builder
	for _, c := range s {
		switch c {
		case '<':
			sb.WriteString("\\lt ")
		case '>':
			sb.WriteString("\\gt ")
		default:
			sb.WriteRune(c)
		}
	}
	return sb.String()
}

// - - - - -

func StringSerialize(s string) []byte {
	r := make([]byte, len(s)+4)
	binary.LittleEndian.PutUint32(r, uint32(len(s)))
	copy(r[4:], s)
	return r
}

func StringDeserialize(b []byte) string {
	var sb strings.Builder
	l := binary.LittleEndian.Uint32(b)
	for i := uint32(0); i < l; i++ {
		sb.WriteByte(b[i+4])
	}
	return sb.String()
}

// - - - - -

func CheckRecognizedEscapeSequence(c rune) (bool, string) {
	switch c {
	case '_', '*', '|', '~', '<', '>', '\\', '$', '`', '#':
		return true, string(c)
	default:
		r := "\\" + string(c)
		log.Warnf(
			"Unrecognized escape sequence \"%v\". Please use \"\\%v\" instead. The sequence will be treated as \"\\%v\"...",
			r, r, r,
		)
		return false, r
	}
}
