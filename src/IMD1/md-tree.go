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
	"reflect"
)

// =====================================
// Markdown tree manipulation

func MDTreeFindElementsOfType[U any](tree *Tree[BlockInterface], x U) []U {
	if tree == nil {
		return nil
	}
	r := make([]U, 0)
	if reflect.TypeOf(tree.Value) == reflect.TypeOf(x) {
		r = append(r, any(tree.Value).(U))
	} else if reflect.TypeOf(tree.Value) == reflect.TypeOf(&BlockInline{}) &&
		reflect.TypeOf(tree.Value.(*BlockInline).Content) == reflect.TypeOf(x) {
		r = append(r, any(tree.Value).(*BlockInline).Content.(U))
	}
	for i := 0; i < len(tree.Children); i++ {
		r = append(r, MDTreeFindElementsOfType(tree.Children[i], x)...)
	}
	return r
}
