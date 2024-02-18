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

import ctypes

imd1_lib_so = ctypes.cdll.LoadLibrary("./libimd1.so")

imd1_lib_so.C_IMD1_MDFileToHTML.argtypes = [ctypes.c_char_p]
imd1_lib_so.C_IMD1_MDFileToHTML.restype = ctypes.c_char_p
def Py_IMD1_MDFileToHTML(py_md_filename):
    c_md_filename = py_md_filename.encode('utf-8')
    c_ret = imd1_lib_so.C_IMD1_MDFileToHTML(c_md_filename)
    return ctypes.string_at(c_ret).decode('utf-8')


html = Py_IMD1_MDFileToHTML("input.md")
print(html)
