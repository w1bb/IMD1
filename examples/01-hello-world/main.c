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

#include <stdlib.h>
#include <stdio.h>

#include "libimd1.h"

int main() {
    // Read the contents of input.md
    FILE *fin = fopen("input.md", "r+");
    if (!fin) {
        fprintf(stderr, "Could not open input.md");
        return -1;
    }
    fseek(fin, 0, SEEK_END);
    long size = ftell(fin);
    fseek(fin, 0, SEEK_SET);
    char *buffer = malloc(sizeof(char) * (size + 1));
    if (!buffer) {
        fprintf(stderr, "Could not allocate memory for buffer");
        return -1;
    }
    fread(buffer, sizeof(char), size, fin);
    buffer[size] = '\0';
    fclose(fin);

    // IMPORTANT: Call IMD1
    struct C_IMD1_MDToHTML_return r = C_IMD1_MDToHTML(buffer);

    // Output the HTML file
    FILE *fout = fopen("output.html", "w+");
    if (!fout) {
        fprintf(stderr, "Could not open output.html");
        return -1;
    }
    fprintf(fout, r.r0);
    fclose(fout);

    // IMPORTANT: Free the values returned by IMD1
    free(r.r0);
    free(r.r1);

    return 0;
}
