# This file is part of the IMD1 project.
# Copyright (c) 2024 Valentin-Ioan Vintilă

# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, version 3.

# This program is distributed in the hope that it will be useful, but
# WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
# General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program. If not, see <http://www.gnu.org/licenses/>.

.PHONY: build run

# Use gcc -Xlinker --verbose for linker-related debugging
mkfile_path := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

build:
	@(cd ../../src && make --no-print-directory)
	@cp ../../src/libimd1{.h,.so} $(mkfile_path)
	gcc -L$(mkfile_path:/=) -limd1 main.c -Wall -o hello -Wl,-rpath .

run:
	./hello
