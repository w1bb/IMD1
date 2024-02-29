#!/usr/bin/env python3

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

from imd1 import IMD1

def run_test_IMD1():
    # TODO - create tests
    def run_test_md_to_html_1(imd1):
        imd1.reset()
        imd1.to_html("Hello world")
        print(imd1.html)
        return True

    # Run the tests
    print("Testing IMD1 (run_test_IMD1())")
    imd1 = IMD1()
    tests = [run_test_md_to_html_1]
    for test in tests:
        print(f"> Test {test.__name__} <=> ", end='')
        test_ok = test(imd1)
        if test_ok:
            print("OK")
        else:
            print("FAILED")

run_test_IMD1()
