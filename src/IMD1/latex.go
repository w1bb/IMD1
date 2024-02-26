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

package IMD1

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// =====================================
// LaTeX interface

type LaTeXInterface interface {
	fmt.Stringer

	GenerateLaTeXTagPrefix() string
	GenerateLaTeXTagSuffix() string
}

// =====================================
// Document LaTeX interface

func (b *BlockDocument) GenerateLaTeXTagPrefix() string {
	return "\\documentclass{article}\n" +
		"\\usepackage[normalem]{ulem}\n" +
		"\\usepackage{float}\n" +
		"\\usepackage{graphicx}\n" +
		"\\usepackage{caption}\n" +
		"\\usepackage{subcaption}\n" +
		"\\usepackage{hyperref}\n" +
		"\\newcommand{\\code}[1]{\\texttt{#1}}\n" +
		"\\begin{document}\n"
}

func (b *BlockDocument) GenerateLaTeXTagSuffix() string {
	return "\\end{document}\n"
}

// =====================================
// Paragraph LaTeX interface

func (b *BlockParagraph) GenerateLaTeXTagPrefix() string {
	return ""
}

func (b *BlockParagraph) GenerateLaTeXTagSuffix() string {
	return "\n\n"
}

// =====================================
// Heading LaTeX interface

func (b *BlockHeading) GenerateLaTeXTagPrefix() string {
	hl := b.HeadingLevel
	if hl > 5 {
		hl = 5
	} else if hl < 1 {
		hl = 1 // This should never be reached, but just in case
	}
	switch hl {
	case 1:
		return "\\section{"
	case 2:
		return "\\subsection{"
	case 3:
		return "\\subsubsection{"
	case 4:
		return "\\paragraph{"
	case 5:
		return "\\subparagraph{"
	}
	return "" // This can never be reached
}

func (b *BlockHeading) GenerateLaTeXTagSuffix() string {
	return "}"
}

// =====================================
// Textbox LaTeX interface

func (b *BlockTextbox) GenerateLaTeXTagPrefix() string {
	return "{" // TODO
}

func (b *BlockTextbox) GenerateLaTeXTagSuffix() string {
	return "}"
}

// =====================================
// Textbox title LaTeX interface

func (b *BlockTextboxTitle) GenerateLaTeXTagPrefix() string {
	return "{" // TODO
}

func (b *BlockTextboxTitle) GenerateLaTeXTagSuffix() string {
	return "}"
}

// =====================================
// Textbox content LaTeX interface

func (b *BlockTextboxContent) GenerateLaTeXTagPrefix() string {
	return "{" // TODO
}

func (b *BlockTextboxContent) GenerateLaTeXTagSuffix() string {
	return "}"
}

// =====================================
// Code listings LaTeX interface

func (b *BlockCodeListing) GenerateLaTeXTagPrefix() string {
	return "{" // TODO
}

func (b *BlockCodeListing) GenerateLaTeXTagSuffix() string {
	return "}"
}

// =====================================
// Comment LaTeX interface

func (b *BlockComment) GenerateLaTeXTagPrefix() string {
	return "" // TODO
}

func (b *BlockComment) GenerateLaTeXTagSuffix() string {
	return ""
}

// =====================================
// HTML LaTeX interface (ignored)

func (b *BlockHTML) GenerateLaTeXTagPrefix() string {
	return "" // Obviously ignored
}

func (b *BlockHTML) GenerateLaTeXTagSuffix() string {
	return "" // Obviously ignored
}

// =====================================
// LaTeX LaTeX interface

func (b *BlockLaTeX) GenerateLaTeXTagPrefix() string {
	return b.RawContent
}

func (b *BlockLaTeX) GenerateLaTeXTagSuffix() string {
	return ""
}

// =====================================
// Inline code listing LaTeX interface

func (b *BlockInlineCodeListing) GenerateLaTeXTagPrefix() string {
	return "\\code{" + b.RawContent
}

func (b *BlockInlineCodeListing) GenerateLaTeXTagSuffix() string {
	return "}"
}

// =====================================
// Math block LaTeX interface

func (b *BlockMath) GenerateLaTeXTagPrefix() string {
	var s string
	switch b.TypeOfBlock {
	case DoubleDollar:
		s = "$$"
	case Brackets:
		s = "\\["
	case BeginEquation:
		s = "\\begin{eqaution}"
	case BeginAlign:
		s = "\\begin{align}"
	default:
		panic(nil) // This should never be reached
	}
	return s + b.RawContent
}

func (b *BlockMath) GenerateLaTeXTagSuffix() string {
	switch b.TypeOfBlock {
	case DoubleDollar:
		return "$$"
	case Brackets:
		return "\\]\n"
	case BeginEquation:
		return "\\end{eqaution}\n"
	case BeginAlign:
		return "\\end{align}\n"
	default:
		panic(nil) // This should never be reached
	}
}

// =====================================
// Inline math LaTeX interface

func (b *BlockInlineMath) GenerateLaTeXTagPrefix() string {
	var s string
	switch b.TypeOfBlock {
	case SingleDollar:
		s = "$"
	case Parenthesis:
		s = "\\("
	default:
		panic(nil) // This should never be reached
	}
	return s + b.RawContent
}

func (b *BlockInlineMath) GenerateLaTeXTagSuffix() string {
	switch b.TypeOfBlock {
	case SingleDollar:
		return "$"
	case Parenthesis:
		return "\\)"
	default:
		panic(nil) // This should never be reached
	}
}

// =====================================
// Unordered list LaTeX interface

func (b *BlockUl) GenerateLaTeXTagPrefix() string {
	return "\\begin{itemize}\n"
}

func (b *BlockUl) GenerateLaTeXTagSuffix() string {
	return "\\end{itemize}\n"
}

// =====================================
// Unordered list item LaTeX interface

func (b *BlockUlLi) GenerateLaTeXTagPrefix() string {
	return "\\item{"
}

func (b *BlockUlLi) GenerateLaTeXTagSuffix() string {
	return "}"
}

// =====================================
// Ordered list LaTeX interface

func (b *BlockOl) GenerateLaTeXTagPrefix() string {
	return "\\begin{enumerate}\n" // TODO - customize based on TypeOfBlock
}

func (b *BlockOl) GenerateLaTeXTagSuffix() string {
	return "\\end{enumerate}\n"
}

// =====================================
// Ordered list item LaTeX interface

func (b *BlockOlLi) GenerateLaTeXTagPrefix() string {
	return "\\item{"
}

func (b *BlockOlLi) GenerateLaTeXTagSuffix() string {
	return "}"
}

// =====================================
// Figure LaTeX interface

func (b *BlockFigure) GenerateLaTeXTagPrefix() string {
	r := "\\begin{figure}[H]\n"
	// TODO - add dock & max-width
	return r
}

func (b *BlockFigure) GenerateLaTeXTagSuffix() string {
	return "\\end{subfigure}\n"
}

// =====================================
// Subfigure LaTeX interface

func (b *BlockSubfigure) GenerateLaTeXTagPrefix() string {
	r := "\\begin{subfigure}\n" // TODO - specify width based on how many subfigures
	r += fmt.Sprintf(
		"\\includegraphics[width=\\textwidth]{%v}\n",
		b.Source,
	)
	r += "\\caption{"
	return r
}

func (b *BlockSubfigure) GenerateLaTeXTagSuffix() string {
	return "}\n\\end{subfigure}\n"
}

// =====================================
// Tabs LaTeX interface

func (b *BlockTabs) GenerateLaTeXTagPrefix() string {
	return "{" // TODO
}

func (b *BlockTabs) GenerateLaTeXTagSuffix() string {
	return "}"
}

// =====================================
// Tab LaTeX interface

func (b *BlockTabsTab) GenerateLaTeXTagPrefix() string {
	return "{" // TODO
}

func (b *BlockTabsTab) GenerateLaTeXTagSuffix() string {
	return "}"
}

// =====================================
// Footnote LaTeX interface

func (b *BlockFootnote) GenerateLaTeXTagPrefix() string {
	return "\\footnote{"
}

func (b *BlockFootnote) GenerateLaTeXTagSuffix() string {
	return "}"
}

// =====================================
// Reference LaTeX interface

func (b *BlockRef) GenerateLaTeXTagPrefix() string {
	return "\\ref{" // TODO - test!
}

func (b *BlockRef) GenerateLaTeXTagSuffix() string {
	return "}"
}

// =====================================
// Reference LaTeX interface

func (b *BlockBibliography) GenerateLaTeXTagPrefix() string {
	if b.LaTeXContent == nil {
		// panic(nil) // TODO - currently, this would always panic
		return ""
	}
	return *b.LaTeXContent
}

func (b *BlockBibliography) GenerateLaTeXTagSuffix() string {
	return ""
}

// =====================================
// Block inline LaTeX interface

func (b *BlockInline) GenerateLaTeXTagPrefix() string {
	return b.Content.GenerateLaTeXTagPrefix()
}

func (b *BlockInline) GenerateLaTeXTagSuffix() string {
	return b.Content.GenerateLaTeXTagSuffix()
}

// =====================================
// Inline document LaTeX interface

func (b *InlineDocument) GenerateLaTeXTagPrefix() string {
	return ""
}

func (b *InlineDocument) GenerateLaTeXTagSuffix() string {
	return ""
}

// =====================================
// Raw string LaTeX interface

func (b *InlineRawString) GenerateLaTeXTagPrefix() string {
	return StringToLaTeXSafe(b.Content)
}

func (b *InlineRawString) GenerateLaTeXTagSuffix() string {
	return ""
}

// =====================================
// String modifier LaTeX interface

func (b *InlineStringModifier) GenerateLaTeXTagPrefix() string {
	switch b.TypeOfModifier {
	case ItalicText:
		return "\\textit{"
	case BoldText:
		return "\\textbf{"
	case StrikeoutText:
		return "\\sout{"
	default:
		panic(nil) // This should never be reached
	}
}

func (b *InlineStringModifier) GenerateLaTeXTagSuffix() string {
	return "}"
}

// =====================================
// Delimiter LaTeX interface

func (b *InlineStringDelimiter) GenerateLaTeXTagPrefix() string {
	// Warn the user that something kind of went wrong
	log.Warnf("When generating the LaTeX, there should be no leftover InlineStringDelimiter (%v). This is a bug and should be reported!", b)
	return b.String()
}

func (b *InlineStringDelimiter) GenerateLaTeXTagSuffix() string {
	return ""
}

// =====================================
// Href LaTeX interface

func (b *InlineHref) GenerateLaTeXTagPrefix() string {
	return fmt.Sprintf(
		"\\href{%v}{",
		b.Address,
	)
}

func (b *InlineHref) GenerateLaTeXTagSuffix() string {
	return "}"
}

// =====================================
// Meta LaTeX interfaces (ignore them)

func (b *BlockMeta) GenerateLaTeXTagPrefix() string {
	return ""
}

func (b *BlockMeta) GenerateLaTeXTagSuffix() string {
	return ""
}

// =====================================
// Meta (author) LaTeX interfaces (ignore them)

func (b *BlockMetaAuthor) GenerateLaTeXTagPrefix() string {
	return ""
}

func (b *BlockMetaAuthor) GenerateLaTeXTagSuffix() string {
	return ""
}

// =====================================
// Meta (copyright) LaTeX interfaces (ignore them)

func (b *BlockMetaCopyright) GenerateLaTeXTagPrefix() string {
	return ""
}

func (b *BlockMetaCopyright) GenerateLaTeXTagSuffix() string {
	return ""
}

// =====================================
// Meta (bibliography) LaTeX interfaces (ignore them)

func (b *BlockMetaBibinfo) GenerateLaTeXTagPrefix() string {
	return ""
}

func (b *BlockMetaBibinfo) GenerateLaTeXTagSuffix() string {
	return ""
}

// =====================================
// Generate LaTeX

func GenerateLaTeX(tree *Tree[BlockInterface]) string {
	if tree == nil {
		return "" // Just to be sure
	}
	var s strings.Builder
	var GenerateLaTeXHelper func(tree *Tree[BlockInterface], sb *strings.Builder)
	GenerateLaTeXHelper = func(tree *Tree[BlockInterface], sb *strings.Builder) {
		sb.WriteString(tree.Value.GenerateLaTeXTagPrefix())
		for i := 0; i < len(tree.Children); i++ {
			GenerateLaTeXHelper(tree.Children[i], sb)
		}
		sb.WriteString(tree.Value.GenerateLaTeXTagSuffix())
	}
	GenerateLaTeXHelper(tree, &s)
	return s.String()
}
