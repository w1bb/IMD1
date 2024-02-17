

<div align="center">
<h1>IMD1</h1>

![version](https://img.shields.io/badge/version-0.6-orange?style=flat)
[![License](https://img.shields.io/badge/license-GPL--3.0-red?style=flat)](/LICENSE)
![docs](https://img.shields.io/badge/docs-30%25-blueviolet?style=flat)
![stars](https://img.shields.io/github/stars/w1bb/IMD1)

</div>
<br />

**IMD1** (to be read as _"I am The One"_) is both a very customizable Markdown-like specification and a Go implementation custom tailored for it.

This project is **by no means complete** - it is a work in progress, kept alive only by my passion for great article-writing software.

## License

This project is distribuited under the [**GNU GPLv3 License**](./LICENSE). If you want to read more about other projects I've brought to life, check out my [personal website](https://v-vintila.com) or my [GitHub page](https://github.com/w1bb).

Table of contents
-----------------

- [Usage](#usage)
    - [Use as a Golang import](#use-as-a-golang-import)
    - [Use as a standalone library](#use-as-a-standalone-library)
- [The specification](#the-specification)
- [API](#api)
    - [Golang API](#golang-api)
    - [Exposed API](#exposed-api)
- [Q&A](#qa)

## Usage

You should first download/clone this repository:

```bash
git clone https://github.com/w1bb/IMD1.git
cd IMD1
```

Depending on what you're trying to accomplish, you should follow the steps bellow.

### Use as a Golang import

Once you've cloned the project, you are set and ready to import it.

### Use as a standalone library

Once you've cloned the project and `cd`-ed into it, build the `.so` library (make sure you have already installed `make` and `go`):

```bash
make build
```

This will create two files, `imd1-lib.so` and `imd1-lib.h`, which can then be used to call some of the _exposed_ functions provided by the library. A list with the API can be found [here](#TODO)

For example, here is a Python snippet that will load the library and call the `C_IMD1_MDFileToHTML` function:

```python
import ctypes
lib = ctypes.cdll.LoadLibrary("./imd1-lib.so")
lib.C_IMD1_MDFileToHTML.argtypes = [ctypes.c_char_p]
lib.C_IMD1_MDFileToHTML.restype = ctypes.c_char_p
html = ctypes.string_at(lib.C_IMD1_MDFileToHTML("input.md".encode('utf-8'))).decode('utf-8')
```

As you can tell, strings have to be converted to C "strings" (char pointers) before passing them to the API. The return value also needs to be converted from a C char pointer to a Python string.

Obviously, using C instead of Python would simplify this process a lot.

## The specification

This document **will NOT** be detail-oriented. The only reason the specification is here is to illustrate the features that have been included in this Markdown-like language.

### Headings and paragraphs

Headings and paragraphs work similarly to the original Markdown format. However, headings can only be written using the `#` character, such as `## Heading 2`. This is by design, since the other format is not only confusing for the reader, but it is also complicated to implement.

### Inline text modifiers

The OG Markdown text modifiers have also suffered some changes. Whilst both `*` and `_` can be used in almost the same way, the _exceptions_ which were ambiguous without looking up the Commonmark documentation have been removed. Now, these operators work using a simple stack where only the last inserted element can be read. It is now obvious that the infamous `_*hello_*` will **not** produce "_*hello*_", but will generate "\_\*hello\_\*" instead.

There is also no difference between `_` and `*`, they both behave in a similar way to how the `*` operator behaved in vanilla Markdown. However, as I've already examplified, they cannot be interchanged.

A new modifier has also been introduced, the strikeout operator, `~`. For example, `~hello~` would yield <del>hello</del>.

### Math support

The IMD1 format supports math equations by design, without altering their content during compilation. The implementation expects a 3rd party library, such as [KaTeX](https://github.com/KaTeX/KaTeX), to handle the way math equations are displayed on the screen.

Equations can be written in various ways, including:

- Inline math equations:

    - `\(..\)` equations, e.g. `\(i=\sqrt{-1}\)`
    - `$..$` equations, e.g. `$f(x)=e^x$` (when outputed to HTML, they will get converted to the `\(..\)` format instead)

- Display math equations

    - `\[..\]` equations, e.g. `\[\frac{x}{23}=\int_0^1\sin t\,dt\]`
    - `$$..$$` equations, e.g. `$$x^{x^x}=\pi$$` (when outputed to HTML, they will get converted to the `\[..\]` format instead)
    - `\begin{equation}..\end{equation}`
    - `\begin{align}..\end{align}`

### Code listings

Inline code can easily be insered using the \` operator. The difference between Markdown and IMD1 is visible when handling multiple lines of code. For this purpose, the \`\`\` operator has been modified to allow for special options. These have to be written using the `[option=value]` format, right after the \`\`\` operator, as in:

~~~
```[option1=value1][option2=value2]
line of code #1
...
line of code #N
```
~~~

Make sure you include **no** spaces between the options, as well as between the first option and the \`\`\` operator. All of the options have to be specified on the same line.

The code listing allow for the following special options:

- Specify the programming language to be displayed: `[lang=...]`

  To write in plain text, this can be set to `text`, `txt`, `plaintext` or nothing at all.

- Specify the file the snippet of code is part of: `[file=...]`
- Specify the way the code should be aligned `[align=...]`. This will be directly coverted to a tag similar to `<div style="text-align: ...;">`
- Specify if you want to enable or not a copy button on top of the listing. This has to be implemented separately by the website on which the HTML is rendered. The option is: `[copy=...]`. Use `allow`, `allowed`, `1`, `true`, `ok` or `yes` to enable this option - any other value is considered `disabled`. 

### Lists

Unordered lists can be created using the `-` character, found at the beginning of the line (possibly prefixed by an indentation). To continue using the same list item, the following non-empty lines have to be prefixed by the same indentation plus two characters. For example:

```
- First line
  Second line

  Third line
- Different list item
```

The first and second lines will act as a single paragraph, while the third line will be considred a different paragraph part of the same list item.

Here is a more complex example:

```
- Item 1
  - Item 1.1
  - Item 1.2

    $$\frac{1}{2}+\frac{1}{3}=\frac{5}{6}$$

    Still item 1.2

- Item 2
```

The example above will be converted **the way you intended it to be converted**, not in a seemingly random pattern designed for compilers instead of humans.

Ordered lists can be used as well. These can be created using the `1.`, `A.`, `a.`, `I.` and `i.` characters at the beginning of the line (possibly prefixed by an indentation). To keep using the same list item, the following non-empty lines have to be preffixed by the same indentation plus three characters. For example:

```
1. Element

   Same element

1. Different element
```

### Textboxes

TODO - the code is complete, docs need to be written

### Figures and subfigures

TODO - the code is complete, docs need to be written

### Footnotes

TODO - the code is complete, docs need to be written

### Bibliography

TODO - the code is complete, docs need to be written

### Metadata

TODO - the code is complete, docs need to be written

## API

### Golang API

TODO

### Exposed API

As I've already [stated above](#use-as-a-standalone-library), some of the functions have been exposed as C variants that can be called outside of Golang.

As a quick note, they all start with `C_IMD1_`. Generally, a function called `C_IMD1_XYZ` will represent the exposed variant of the Golang function `IMD1_XYZ`.

#### C_IMD1_MDFileToHTMLFile

It is the exposed variant of `IMD1_MDFileToHTMLFile`.

Parameters:

- `c_md_filename` - C string (`char*`)
- `c_html_filename` - C string (`char*`)

Returns: nothing (`void`)

The function converts an IMD1 file (located at `c_md_filename`) into a HTML file (located at `c_html_filename`).

#### C_IMD1_MDToHTMLFile

It is the exposed variant of `IMD1_MDToHTMLFile`.

Parameters:

- `c_s` - C string (`char*`)
- `c_html_filename` - C string (`char*`)

Returns: nothing (`void`)

The function converts an IMD1 string (`c_s`) into a HTML file (located at `c_html_filename`).

#### C_IMD1_MDFileToHTML

It is the exposed variant of `IMD1_MDFileToHTML`.

Parameters:

- `c_md_filename` - C string (`char*`)

Returns: C string (`char*`)

The function converts an IMD1 file (located at `c_md_filename`) into a valid HTML file, which is not saved on the disk, but is returned as a C string instead.

#### C_IMD1_MDToHTML

It is the exposed variant of `IMD1_MDToHTML`.

Parameters:

- `c_s` - C string (`char*`)

Returns: C string (`char*`)

The function converts an IMD1 string (`c_s`) into a valid HTML file, which is not saved on the disk, but is returned as a C string instead.

## Q&A

**Q:** Is this a subset/superset of the [Commonmark Markdown specification](https://spec.commonmark.org)?

**A:** (Un)fortunately, it is not. That specification is rather complex and it contains a lot of unnecessary features. Thus, I've decided to keep what I considered important and to build upon that foundation. For example, unlike vanilla markdown, IMD1 can deal with metadata (i.e. author and copyright details), bibliography, math equations, figures and subfigures by design.

**Q:** Why didn't you just modify an existing Markdown engine?

**A:** Because I wanted some first-hand experience with writing a compiler from scratch, as well as to learn some basic Go. The whole project, as of _version 0.6_, took me exactly one week to complete. I am not trying to brag, but every once in a while, it is nice to be satisfied with your results! 