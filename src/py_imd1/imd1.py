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

class IMD1:
    def __init__(self):
        self.html = None
        self.meta = None

        current_file_path = os.path.abspath(__file__)
        current_folder_path = os.path.dirname(current_file_path)
        lib_path = os.path.join(current_folder_path, "libimd1.so")
        self.__imd1_lib_so = ctypes.cdll.LoadLibrary(lib_path)

        self.__imd1_lib_so.C_FreeUnsafePointer.argtypes = [ctypes.c_void_p]
        # self.__imd1_lib_so.C_IMD1_MDFileToHTMLFile.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
        # self.__imd1_lib_so.C_IMD1_MDFileToHTMLFile.restype = ctypes.c_void_p
        # self.__imd1_lib_so.C_IMD1_MDToHTMLFile.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
        # self.__imd1_lib_so.C_IMD1_MDToHTMLFile.restype = ctypes.c_void_p
        # class __C_IMD1_MDFileToHTML_return(ctypes.Structure):
        #     _fields_ = [('r0', ctypes.c_char_p),
        #                 ('r1', ctypes.c_void_p)]
        # self.__imd1_lib_so.C_IMD1_MDFileToHTML.argtypes = [ctypes.c_char_p]
        # self.__imd1_lib_so.C_IMD1_MDFileToHTML.restype = __C_IMD1_MDFileToHTML_return
        class __C_IMD1_MDToHTML_return(ctypes.Structure):
            _fields_ = [('r0', ctypes.c_char_p),
                        ('r1', ctypes.c_void_p)]
        self.__imd1_lib_so.C_IMD1_MDToHTML.argtypes = [ctypes.c_char_p]
        self.__imd1_lib_so.C_IMD1_MDToHTML.restype = __C_IMD1_MDToHTML_return

    def __fill_meta(self, buffer):
        b = ctypes.cast(buffer, ctypes.c_void_p).value
        self.meta = {}
        author_len = ctypes.c_uint32.from_address(b).value
        if author_len > 0:
            self.meta["author"] = ctypes.string_at(b+4, author_len).decode('utf-8')
        b += 4 + author_len
        copyright_len = ctypes.c_uint32.from_address(b).value
        if copyright_len > 0:
            self.meta["copyright"] = ctypes.string_at(b+4, copyright_len).decode('utf-8')
    
    def md_to_html(self, md_string):
        c_s = md_string.encode('utf-8')
        c_ret = self.__imd1_lib_so.C_IMD1_MDToHTML(c_s)
        self.html = ctypes.string_at(c_ret.r0).decode('utf-8')
        self.__fill_meta(c_ret.r1)
        self.__imd1_lib_so.C_FreeUnsafePointer(c_ret.r1)
    
    def reset(self):
        self.html = None
        self.meta = None
