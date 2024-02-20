<div align="center">
<h1>IMD1</h1>

![version](https://img.shields.io/badge/version-0.6-orange?style=flat)
[![License](https://img.shields.io/badge/license-GPL--3.0-red?style=flat)](/LICENSE)
![docs](https://img.shields.io/badge/docs-50%25-blueviolet?style=flat)
![stars](https://img.shields.io/github/stars/w1bb/IMD1)

</div>
<br />

**IMD1** (to be read as _"I am The One"_) is both a very customizable Markdown-like specification and a Go implementation custom tailored for it that can convert from **IMD1** to both **HTML** and **LaTeX**.

The compiler **avoids Regex** at all costs! [Regex is **the worst!**](https://blog.cloudflare.com/details-of-the-cloudflare-outage-on-july-2-2019)

This project is **by no means complete** - it is a work in progress, kept alive only by my passion for great article-writing software.

## License

This project is distribuited under the [**GNU GPLv3 License**](./LICENSE). If you want to read more about other projects I've brought to life, check out my [personal website](https://v-vintila.com) or my [GitHub page](https://github.com/w1bb).

Table of contents
-----------------

- [Usage](#usage)
    - [Use as a Golang import](#use-as-a-golang-import)
    - [Use as a standalone library](#use-as-a-standalone-library)
- [The specification](#the-specification)
    - [Headings and paragraphs](#headings-and-paragraphs)
    - [Inline text modifiers](#inline-text-modifiers)
    - [Links](#links)
    - [Math support](#math-support)
    - [Code listings](#code-listings)
    - [HTML and LaTeX](#html-and-latex)
    - [Lists](#lists)
    - [Textboxes](#textboxes)
    - [Figures and subfigures](#figures-and-subfigures)
    - [Tabs](#tabs)
    - [Footnotes](#footnotes)
    - [Metadata](#metadata)
    - [Bibliography](#bibliography)
    - [Comments](#comments)
- [API](#api)
    - [Golang API](#golang-api)
    - [Exposed API](#exposed-api)
    - [API examples](#api-examples)
- [Q&A](#qa)

## Usage

You should first download/clone this repository:

```bash
git clone https://github.com/w1bb/IMD1.git
cd IMD1/src
```

Depending on what you're trying to accomplish, you should follow the steps below.

### Use as a Golang import

Once you've cloned the project, you are set and ready to import it.

### Use as a standalone library

Once you've cloned the project and `cd`-ed into it, build the `.so` library (make sure you have already installed `make` and `go`):

```bash
make build
```

This will create two files, `libimd1.so` and `libimd1.h`, which can then be used to call some of the _exposed_ functions provided by the library. A list with the API can be found [here](#TODO)

For example, here is a Python snippet that will load the library and call the `C_IMD1_MDFileToHTML` function:

```python
import ctypes
lib = ctypes.cdll.LoadLibrary("./libimd1.so")
lib.C_IMD1_MDFileToHTML.argtypes = [ctypes.c_char_p]
lib.C_IMD1_MDFileToHTML.restype = ctypes.c_char_p
html = ctypes.string_at(lib.C_IMD1_MDFileToHTML("input.md".encode('utf-8'))).decode('utf-8')
```

As you can tell, strings have to be converted to C "strings" (char pointers) before passing them to the API. The return value also needs to be converted from a C char pointer to a Python string.

A file containing the Python variants [exists](./src/py_imd1/imd1.py) within the repository and shall be used to aid development.

Obviously, using C instead of Python would simplify this "translation" process a lot. However, linking errors may arise. Check out the [examples](#api-examples) for more details.

## The specification

This document **will NOT** be detail-oriented. The only reason the specification is here is to illustrate the features that have been included in this Markdown-like language.

### Headings and paragraphs

**Headings** and **paragraphs** work similarly to the original Markdown format. However, headings can only be written using the `#` character, such as `## Heading 2`. This is by design, since the other format is not only confusing for the reader, but it is also complicated to implement.

### Inline text modifiers

The OG Markdown **text modifiers** have also suffered some changes. Whilst both `*` and `_` can be used in almost the same way, the _exceptions_ which were ambiguous without looking up the Commonmark documentation have been removed. Now, these operators work using a simple stack where only the last inserted element can be read. It is now obvious that the infamous `_*hello_*` will **not** produce "_*hello*_", but will generate "\_\*hello\_\*" instead.

There is also no difference between `_` and `*`, they both behave in a similar way to how the `*` operator behaved in vanilla Markdown. However, as I've already examplified, they cannot be interchanged.

A new modifier has also been introduced, the strikeout operator, `~`. For example, `~hello~` would yield <del>hello</del>.

### Links

**Links** are supported as well and their syntax follows the Commonmark guidelines as close as possible. To make it work, I've used a modified pushdown automaton - this ensures both speed and correctness. As I've stated above, I **never** used Regex in this project due to its _many_ flaws!

The syntax is `[text](link/to/resource)`. For example, `[hey](https://archlinux.org)` would produce [hey](https://archlinux.org).

### Math support

The IMD1 format supports **math equations** by design, without altering their content during compilation. The implementation expects a 3rd party library, such as [KaTeX](https://github.com/KaTeX/KaTeX), to handle the way math equations are displayed on the screen.

Equations can be written in various ways, including:

- Inline math equations:

    - `\(..\)` equations, e.g. `\(i=\sqrt{-1}\)`
    - `$..$` equations, e.g. `$f(x)=e^x$` (when outputed to HTML, they will get converted to the `\(..\)` format instead; this is not the case for LaTeX)

- Display math equations

    - `\[..\]` equations, e.g. `\[\frac{x}{23}=\int_0^1\sin t\,dt\]`
    - `$$..$$` equations, e.g. `$$x^{x^x}=\pi$$` (when outputed to HTML, they will get converted to the `\[..\]` format instead; this is not the case for LaTeX)
    - `\begin{equation}..\end{equation}`
    - `\begin{align}..\end{align}`

### Code listings

**Inline code** can easily be insered using the \` operator. The difference between Markdown and IMD1 is visible when handling **multiple lines of code**. For this purpose, the \`\`\` operator has been modified to allow for special options. These have to be written using the `[option=value]` format, right after the \`\`\` operator, as in:

~~~
```[option1=value1][option2=value2]
line of code #1
...
line of code #N
```
~~~

Make sure you include **no** spaces between the options, as well as between the first option and the \`\`\` operator. All of the options have to be specified on the same line.

The code listing allow for the following special options:

- Specify the programming language to be displayed: `[lang=..]`

  To write in plain text, this can be set to `text`, `txt`, `plaintext` or nothing at all.

- Specify the file the snippet of code is part of: `[file=..]`
- Specify the way the code should be aligned `[align=..]`. This will be directly coverted to a tag similar to `<div style="text-align: ..;">`
- Specify if you want to enable or not a copy button on top of the listing. This has to be implemented separately by the website on which the HTML is rendered. The option is: `[copy=..]`. Use `allow`, `allowed`, `1`, `true`, `ok` or `yes` to enable this option - any other value is considered `disabled`. This option has no effect in LaTeX.

### HTML and LaTeX

You can write **HTML** and **LaTeX** directly inside the document. When compiling into HTML, the LaTeX tags will be ignored and vice-versa.

To insert some HTML, use the `|html>..<html|` tag. To insert LaTeX, use `|latex>..<latex|`.

### Lists

**Unordered lists** can be created using the `-` character, found at the beginning of the line (possibly prefixed by an indentation). To continue using the same list item, the following non-empty lines have to be prefixed by the same indentation plus two characters. For example:

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

**Ordered lists** can be used as well. These can be created using the `1.`, `A.`, `a.`, `I.` and `i.` characters at the beginning of the line (possibly prefixed by an indentation). To keep using the same list item, the following non-empty lines have to be preffixed by the same indentation plus three characters. For example:

```
1. Element

   Same element

1. Different element
```

### Textboxes

A **textbox** is meant to be a colored box that contains both a title and some content. These can be used to differentiate between normal text and some crucial piece of information, a warning or an otherwise important message.

In IMD1, you can generate a textbox using the `|textbox>..<textbox|` tag. You can then specify what kind of textbox you want using the `[class=..]` option. The standard (recommended) values are `tip`, `note`, `warning`, `error`, and `admonition`.

You will now be able to specify a title (`|title>..<title|`) and some content (`|content>..<content|`).

_Note: Previous versions did not use the `|title>` and `|content>` tags. Instead, you had to specify the title using the `[title=..]` option. The reason this was changed is because you could only insert unformatable text for the title due to the way options are parsed by the compiler._

### Figures and subfigures

The original Markdown format allows for images to be included using the `![alt](link "title")` format. However, IMD1 opts for both an easier to parse and more powerful approach.

You can create **figures** using the `|figure>..<figure|` syntax; these represent collections of **subfigures** (`|subfigure>..<subfigure|`), defined as images with advanced captions. To illustrate this, check out the example below:

```
|figure>[dock=center]
|subfigure>[src=/path/to/image]
This is the caption for this subfigure. It might contain $\LaTeX$ or even

- Lists
  - And nested lists
<subfigure|
<figure|
```

It is now obvious that both figures and subfigures are customizable using the same syntax described above for code listings.

Figures allow these options:

- Specify if the subfigures (images+captions) will be anchored to the top, the center or the bottom of the structure (figure) as a whole. Use `[dock=..]` with one of the following values: `dock-top` (`top`), `dock-bottom` (`dock-bot`, `bottom`, `bot`) or `center`.
- Specify a maximum width for the whole structure (figure). Use `[max-width=..]`. This will be directly coverted to a tag similar to `<div style="max-width: ..;">`
- Specify a global padding for each subfigure that will be applied if the subfigures themselves don't override this option - `[padding=..]`

Subfigures allow these options:

- Specify the source of the image to be rendered using `[src=..]`. This is optional - figures and subfigures could be used in other creative ways, but this is outside of the scope of our discussion.
- Specify a padding that will override the global padding set by the parent figure. Use `[padding=..]`

Please note that, in order for this to get rendered correctly, some CSS has to be written.

### Tabs

TODO - incomplete

### Footnotes

TODO - the code is complete, docs need to be written

### Metadata

TODO - the code is complete, docs need to be written

### Bibliography

You can write a **bibliography** for your files based on a JSON file. It's structure should be similar to the following example:

```json
{
    "bibliography": [
        {
            "tag": "test-tag-1",
            "type": "article",
            "data": {
                "title": "Article title 1",
                "author": "Article author 2"
            }
        },
        {
            "tag": "test-tag-2",
            "type": "article",
            "data": {
                "title": "Article title 2",
                "author": "Article author 2"
            }
        }
    ]
}
```

As you can see, every entry is contained in this file. Each entry must contain a **tag** that will later be used when referencing your work. Currently, the following types of entries are allowed:

- `article`
- `book`
- `other` (`unknown`)

The "data" section can contain any number of:

- `title` (if missing, the tag will be considered the title)
- `author`
- `journal`
- `volume`
- `number`
- `pages`
- `year`
- `publisher`
- `url`

You can include a bibliography JSON file using `|bibinfo>/path/to/json<bibinfo|` - please note that `|bibinfo>` can only be part of the `|meta>..<meta|` tag (see [Metadata](#metadata)).

You can also include an inline bibliography file (meaning you can directly paste the contents of the JSON file in the tag itself) by specifing the option `[inline=true]`. Multiple `|bibinfo>` tags are allowed in the same document and they will be processed in the order they are written.

To **reference** the bibliography, you can use the `|ref>..<ref|` tag. Make sure to specify a valid tag. The "ref" can precede the the "bibinfo" and it will still be rendered correctly (the `|bibinfo>` tags get processed before the `|ref>` ones).

Here is an example of valid bibliography

```
|meta>|bibinfo>[inline=true]
{
    "bibliography": [
        {
            "tag": "tag-1",
            "type": "article",
            "data": {
                "title": "Article title 1",
                "author": "Article author 1"
            }
        }
    ]
}
<bibinfo||bibinfo>[inline=true]
{
    "bibliography": [
        {
            "tag": "tag-2",
            "type": "article",
            "data": {
                "title": "Article title 2",
                "author": "Article author 2"
            }
        },
        {
            "tag": "tag-2",
            "type": "article",
            "data": {
                "title": "Article title 2",
                "author": "Article author 2"
            }
        }
    ]
}
<bibinfo|<meta|

The reference |ref>tag-1<ref| is valid, and so is |ref>tag-2<ref| (however, the compiler will trigger a warning). But |ref>hello<ref| is invalid and will be converted to [?].
```

### Comments

Finally, there are **comments**. These work in the same way HTML comments would, meaning you can use the `<!--..-->` syntax to write them. Currently, these comments **will get copied** to the final HTML/LaTeX file.

## API

### Golang API

TODO

### Exposed API

As I've already [stated above](#use-as-a-standalone-library), some of the functions have been exposed as C variants that can be called outside of Golang.

As a quick note, they all start with `C_IMD1_`. Generally, a function called `C_IMD1_XYZ` will represent the exposed variant of the Golang function `IMD1_XYZ`.

There is also a Python wrapper that calls these exposed functions. Check out the [py_imd1/imd1.py](./src/py_imd1/imd1.py) file for additional information.

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

#### C_IMD1_MDFileToLaTeXFile

It is the exposed variant of `IMD1_MDFileToLaTeXFile`.

Parameters:

- `c_md_filename` - C string (`char*`)
- `c_latex_filename` - C string (`char*`)

Returns: nothing (`void`)

The function converts an IMD1 file (located at `c_md_filename`) into a LaTeX file (located at `c_latex_filename`).

#### C_IMD1_MDToLaTeXFile

It is the exposed variant of `IMD1_MDToLaTeXFile`.

Parameters:

- `c_s` - C string (`char*`)
- `c_latex_filename` - C string (`char*`)

Returns: nothing (`void`)

The function converts an IMD1 string (`c_s`) into a LaTeX file (located at `c_latex_filename`).

#### C_IMD1_MDFileToLaTeX

It is the exposed variant of `IMD1_MDFileToLaTeX`.

Parameters:

- `c_md_filename` - C string (`char*`)

Returns: C string (`char*`)

The function converts an IMD1 file (located at `c_md_filename`) into a valid LaTeX file, which is not saved on the disk, but is returned as a C string instead.

#### C_IMD1_MDToLaTeX

It is the exposed variant of `IMD1_MDToLaTeX`.

Parameters:

- `c_s` - C string (`char*`)

Returns: C string (`char*`)

The function converts an IMD1 string (`c_s`) into a valid LaTeX file, which is not saved on the disk, but is returned as a C string instead.

### API Examples

You can find a few examples in the [examples](./examples/) folder _(TODO: add more examples)_. Let's explore some of them:

#### 01-hello-world (C example)

The [01-hello-world](./examples/01-hello-world/) gives an example of how a C program can call the IMD1 API using dynamic linking. The most important file in this example is the [Makefile](./examples/01-hello-world/Makefile) - the arguments used when calling `gcc` ensure that the `libimd1.so` library gets loaded correctly at runtime.

If you are experiencing linking issues, you can either use `ldd` or `gcc -Xlinker --verbose` to debug them.

## Q&A

**Q:** Is this a subset/superset of the [Commonmark Markdown specification](https://spec.commonmark.org)?

**A:** (Un)fortunately, it is not. That specification is rather complex and it contains a lot of unnecessary features. Thus, I've decided to keep what I considered important and to build upon that foundation. For example, unlike vanilla markdown, IMD1 can deal with metadata (i.e. author and copyright details), bibliography, math equations, figures and subfigures by design.

**Q:** Why didn't you just modify an existing Markdown engine?

**A:** Because I wanted some first-hand experience with writing a compiler from scratch, as well as to learn some basic Go. The whole project, as of _version 0.6_, took me exactly one week to complete. I am not trying to brag, but every once in a while, it is nice to be satisfied with your results! 