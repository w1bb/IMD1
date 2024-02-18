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

import "fmt"

// =====================================
// Stack data structure

type Stack[T fmt.Stringer] struct {
	full_content []T
}

func (s Stack[T]) String() string {
	if s.Empty() {
		return "<Stack: []>"
	}
	r := "<Stack: ["
	for i := 0; i < len(s.full_content); i++ {
		r += s.full_content[i].String() + ", "
	}
	return r[:len(r)-2] + "]>"
}

func (s *Stack[T]) Empty() bool {
	return s.full_content == nil || len(s.full_content) == 0
}

func (s *Stack[T]) Size() int {
	if s.full_content == nil {
		return 0
	}
	return len(s.full_content)
}

func (s *Stack[T]) Top() *T {
	if s.Empty() {
		return nil
	}
	return &s.full_content[len(s.full_content) - 1]
}

func (s *Stack[T]) Push(element T) {
	s.full_content = append(s.full_content, element)
}

func (s *Stack[T]) Pop() *T {
	if s.Empty() {
		return nil
	}
	x := &s.full_content[len(s.full_content) - 1]
	s.full_content = s.full_content[:len(s.full_content) - 1]
	return x
}
