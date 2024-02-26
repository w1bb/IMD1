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
	"fmt"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

// =====================================
// Tree data structure

type Tree[T fmt.Stringer] struct {
	Children []*Tree[T]
	Value T
	Parent *Tree[T]
}

func (t Tree[T]) String() string {
	var spaces = 0
	var helper func (t *Tree[T]) string
	helper = func (t *Tree[T]) string {
		var spaces_str, result string
		if spaces == 0 {
			spaces_str = ""
		} else {
			spaces_str = strings.Repeat("  ", spaces - 2) + "- "
		}
		result += spaces_str + t.Value.String() + "\n"
		spaces += 2
		for i := range t.Children {
			result += helper(t.Children[i])
		}
		spaces -= 2
		return result
	}
	return helper(&t)
}

func (t *Tree[T]) Verify(enable_log bool) bool {
	var helper func (t *Tree[T], tail []*Tree[T], enable_log bool) bool
	helper = func (t *Tree[T], tail []*Tree[T], enable_log bool) bool {
		if t.Children == nil || len(t.Children) == 0 {
			return true
		}
		new_tail := append(tail, t)
		ok := true
		for i := 0; i < len(t.Children); i++ {
			if !reflect.DeepEqual(t.Children[i].Parent, t) {
				ok = false
				if enable_log {
					s_tail := ""
					for ti := 0; ti < len(tail); ti++ {
						s_tail += tail[ti].Value.String() + "->"
					}
					s_tail += t.Value.String()
					log.Errorf(
						"Invalid tree:\n  Child %v\n  Correct parent: %v\n  Detected parent: %v\n",
						t.Children[i].Value.String(),
						s_tail,
						t.Children[i].Parent.Value.String(),
					)
				}
			}
			if !helper(t.Children[i], new_tail, enable_log) {
				ok = false
			}
		}
		return ok
	}
	return helper(t, make([]*Tree[T], 0), enable_log)
}

func FindElementsOfType[T fmt.Stringer, U any](tree *Tree[T], x U) []U {
	if tree == nil {
		return nil
	}
	r := make([]U, 0)
	if reflect.TypeOf(tree.Value) == reflect.TypeOf(x) {
		r = append(r, any(tree.Value).(U))
	}
	for i := 0; i < len(tree.Children); i++ {
		r = append(r, FindElementsOfType[T, U](tree.Children[i], x)...)
	}
	return r
}
