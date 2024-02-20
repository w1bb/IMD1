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

import ctypes
import os



# =====================================
# Load the library

current_file_path = os.path.abspath(__file__)
current_folder_path = os.path.dirname(current_file_path)
lib_path = os.path.join(current_folder_path, "libimd1.so")
imd1_lib_so = ctypes.cdll.LoadLibrary(lib_path)



# =====================================
# Markdown to HTML (Python, from C-exported variants)

imd1_lib_so.C_IMD1_MDFileToHTMLFile.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
def Py_IMD1_MDFileToHTMLFile(py_md_filename, py_html_filename):
    c_md_filename = py_md_filename.encode('utf-8')
    c_html_filename = py_html_filename.encode('utf-8')
    imd1_lib_so.C_IMD1_MDFileToHTMLFile(c_md_filename, c_html_filename)

imd1_lib_so.C_IMD1_MDToHTMLFile.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
def Py_IMD1_MDToHTMLFile(py_s, py_html_filename):
    c_s = py_s.encode('utf-8')
    c_html_filename = py_html_filename.encode('utf-8')
    imd1_lib_so.C_IMD1_MDToHTMLFile(c_s, c_html_filename)

imd1_lib_so.C_IMD1_MDFileToHTML.argtypes = [ctypes.c_char_p]
imd1_lib_so.C_IMD1_MDFileToHTML.restype = ctypes.c_char_p
def Py_IMD1_MDFileToHTML(py_md_filename):
    c_md_filename = py_md_filename.encode('utf-8')
    c_ret = imd1_lib_so.C_IMD1_MDFileToHTML(c_md_filename)
    return ctypes.string_at(c_ret).decode('utf-8')

imd1_lib_so.C_IMD1_MDToHTML.argtypes = [ctypes.c_char_p]
imd1_lib_so.C_IMD1_MDToHTML.restype = ctypes.c_char_p
def Py_IMD1_MDToHTML(py_s):
    c_s = py_s.encode('utf-8')
    c_ret = imd1_lib_so.C_IMD1_MDToHTML(c_s)
    return ctypes.string_at(c_ret).decode('utf-8')



# =====================================
# Markdown to LaTeX (Python, from C-exported variants)

imd1_lib_so.C_IMD1_MDFileToLaTeXFile.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
def Py_IMD1_MDFileToLaTeXFile(py_md_filename, py_latex_filename):
    c_md_filename = py_md_filename.encode('utf-8')
    c_latex_filename = py_latex_filename.encode('utf-8')
    imd1_lib_so.C_IMD1_MDFileToLaTeXFile(c_md_filename, c_latex_filename)

imd1_lib_so.C_IMD1_MDToLaTeXFile.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
def Py_IMD1_MDToLaTeXFile(py_s, py_latex_filename):
    c_s = py_s.encode('utf-8')
    c_latex_filename = py_latex_filename.encode('utf-8')
    imd1_lib_so.C_IMD1_MDToLaTeXFile(c_s, c_latex_filename)

imd1_lib_so.C_IMD1_MDFileToLaTeX.argtypes = [ctypes.c_char_p]
imd1_lib_so.C_IMD1_MDFileToLaTeX.restype = ctypes.c_char_p
def Py_IMD1_MDFileToLaTeX(py_md_filename):
    c_md_filename = py_md_filename.encode('utf-8')
    c_ret = imd1_lib_so.C_IMD1_MDFileToLaTeX(c_md_filename)
    return ctypes.string_at(c_ret).decode('utf-8')

imd1_lib_so.C_IMD1_MDToLaTeX.argtypes = [ctypes.c_char_p]
imd1_lib_so.C_IMD1_MDToLaTeX.restype = ctypes.c_char_p
def Py_IMD1_MDToLaTeX(py_s):
    c_s = py_s.encode('utf-8')
    c_ret = imd1_lib_so.C_IMD1_MDToLaTeX(c_s)
    return ctypes.string_at(c_ret).decode('utf-8')
