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

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// =====================================
// HTML interface

type HTMLInterface interface {
	fmt.Stringer

	GenerateHTMLTagPrefix() string
	GenerateHTMLTagSuffix() string
}

// =====================================
// Document HTML interface

func (b *BlockDocument) GenerateHTMLTagPrefix() string {
	return "<!DOCTYPE html>\n<html>\n<head><title></title></head>\n<body>\n"
}

func (b *BlockDocument) GenerateHTMLTagSuffix() string {
	return "</body>\n</html>\n"
}

// =====================================
// Paragraph HTML interface

func (b *BlockParagraph) GenerateHTMLTagPrefix() string {
	return "<p>"
}

func (b *BlockParagraph) GenerateHTMLTagSuffix() string {
	return "</p>\n"
}

// =====================================
// Heading HTML interface

func (b *BlockHeading) GetHTMLHeadingLevel() int {
	hl := b.HeadingLevel
	if hl > 6 {
		hl = 6
	} else if hl < 1 {
		hl = 1 // This should never be reached, but just in case
	}
	return hl
}

func (b *BlockHeading) GenerateHTMLTagPrefix() string {
	return fmt.Sprintf("<h%v>", b.GetHTMLHeadingLevel())
}

func (b *BlockHeading) GenerateHTMLTagSuffix() string {
	return fmt.Sprintf("</h%v>\n", b.GetHTMLHeadingLevel())
}

// =====================================
// Textbox HTML interface

func (b *BlockTextbox) GenerateHTMLTagPrefix() string {
	c := ""
	if b.Class != "" {
		c = " " + b.Class
	}
	return fmt.Sprintf(
		"<div class=\"box%v\">\n",
		c,
	)
}

func (b *BlockTextbox) GenerateHTMLTagSuffix() string {
	return "</div>\n"
}

// =====================================
// Textbox title HTML interface

func (b *BlockTextboxTitle) GenerateHTMLTagPrefix() string {
	return "<div class=\"box-title\">"
}

func (b *BlockTextboxTitle) GenerateHTMLTagSuffix() string {
	return "</div>\n"
}

// =====================================
// Textbox content HTML interface

func (b *BlockTextboxContent) GenerateHTMLTagPrefix() string {
	return "<div class=\"box-content\">"
}

func (b *BlockTextboxContent) GenerateHTMLTagSuffix() string {
	return "</div>\n"
}

// =====================================
// Code listings HTML interface

func (b *BlockCodeListing) GenerateHTMLTagPrefix() string {
	r := "<div class=\"code-listing\">"
	if b.AllowCopy {
		r += "<div class=\"copy-code\"></div>"
	}
	if b.Filename != "" {
		r += fmt.Sprintf(
			"<div class=\"file-name\">%v</div>",
			b.Filename,
		)
	}
	if b.TextAlign != "" {
		r += fmt.Sprintf(
			"<pre style=\"text-align: %v\">",
			b.TextAlign,
		)
	} else {
		r += "<pre>"
	}
	return fmt.Sprintf(
		"%v<code class=\"language-%v\">%v",
		r,
		b.Language,
		b.RawContent,
	)
}

func (b *BlockCodeListing) GenerateHTMLTagSuffix() string {
	return "</code></pre></div>\n"
}

// =====================================
// HTML HTML interface

func (b *BlockHTML) GenerateHTMLTagPrefix() string {
	return b.RawContent
}

func (b *BlockHTML) GenerateHTMLTagSuffix() string {
	return ""
}

// =====================================
// LaTeX HTML interface (ignored)

func (b *BlockLaTeX) GenerateHTMLTagPrefix() string {
	return "" // Obviously ignored
}

func (b *BlockLaTeX) GenerateHTMLTagSuffix() string {
	return "" // Obviously ignored
}

// =====================================
// Inline code listing HTML interface

func (b *BlockInlineCodeListing) GenerateHTMLTagPrefix() string {
	return "<code>" + b.RawContent
}

func (b *BlockInlineCodeListing) GenerateHTMLTagSuffix() string {
	return "</code>"
}

// =====================================
// Math block HTML interface

func (b *BlockMath) GenerateHTMLTagPrefix() string {
	var s string
	switch b.TypeOfBlock {
	case DoubleDollar, Brackets:
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

func (b *BlockMath) GenerateHTMLTagSuffix() string {
	switch b.TypeOfBlock {
	case DoubleDollar, Brackets:
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
// Inline math HTML interface

func (b *BlockInlineMath) GenerateHTMLTagPrefix() string {
	var s string
	switch b.TypeOfBlock {
	case SingleDollar, Parenthesis:
		s = "\\("
	default:
		panic(nil) // This should never be reached
	}
	return s + b.RawContent
}

func (b *BlockInlineMath) GenerateHTMLTagSuffix() string {
	switch b.TypeOfBlock {
	case SingleDollar, Parenthesis:
		return "\\)"
	default:
		panic(nil) // This should never be reached
	}
}

// =====================================
// Unordered list HTML interface

func (b *BlockUl) GenerateHTMLTagPrefix() string {
	return "<ul>\n"
}

func (b *BlockUl) GenerateHTMLTagSuffix() string {
	return "</ul>\n"
}

// =====================================
// Unordered list item HTML interface

func (b *BlockUlLi) GenerateHTMLTagPrefix() string {
	return "<li>\n"
}

func (b *BlockUlLi) GenerateHTMLTagSuffix() string {
	return "</li>\n"
}

// =====================================
// Ordered list HTML interface

func (t BlockOlType) HTMLType() string {
	switch t {
	case OlType_1:
		return "1"
	case OlType_A:
		return "A"
	case OlType_a:
		return "a"
	case OlType_I:
		return "I"
	case OlType_i:
		return "i"
	default:
		panic(nil) // This should never be reached
	}
}

func (b *BlockOl) GenerateHTMLTagPrefix() string {
	return fmt.Sprintf(
		"<ol type=\"%v\">\n",
		b.TypeOfBlock.HTMLType(),
	)
}

func (b *BlockOl) GenerateHTMLTagSuffix() string {
	return "</ol>\n"
}

// =====================================
// Ordered list item HTML interface

func (b *BlockOlLi) GenerateHTMLTagPrefix() string {
	return "<li>\n"
}

func (b *BlockOlLi) GenerateHTMLTagSuffix() string {
	return "</li>\n"
}

// =====================================
// Figure HTML interface

func (b *BlockFigure) GenerateHTMLTagPrefix() string {
	r := "<div class=\"figure"
	switch b.Dock {
	case "dock-top":
		r += " dock-top\""
	case "dock-bottom":
		r += " dock-bottom\""
	default: // warning has already been triggered by b.ExecuteAfterBlockStarts
		r += "\""
	}
	if b.MaxWidth != "" {
		r += fmt.Sprintf(" style=\"max-width: %v;\"", b.MaxWidth)
	}
	r += ">\n"
	return r
}

func (b *BlockFigure) GenerateHTMLTagSuffix() string {
	return "</div>\n"
}

// =====================================
// Subfigure HTML interface

func (b *BlockSubfigure) GenerateHTMLTagPrefix() string {
	img_tag := fmt.Sprintf(
		"<img src=\"%v\"",
		b.Source,
	)
	if b.Padding != "" {
		img_tag += fmt.Sprintf(
			" style=\"padding: %v;\">",
			b.Padding,
		)
	} else {
		img_tag += ">"
	}
	return fmt.Sprintf(
		"<div class=\"subfigure\">%v<div class=\"subcaption\">",
		img_tag,
	)
}

func (b *BlockSubfigure) GenerateHTMLTagSuffix() string {
	return "</div></div>\n"
}

// =====================================
// Footnote HTML interface

func (b *BlockFootnote) GenerateHTMLTagPrefix() string {
	return fmt.Sprintf(
		"<a href=\"#\" class=\"footnote-href\" onclick=\"TODO\"><div class=\"footnote footnote-%v\">",
		b.FootnoteIndex,
	)
}

func (b *BlockFootnote) GenerateHTMLTagSuffix() string {
	return "</div></a>\n"
}

// =====================================
// Reference HTML interface

func (b *BlockRef) GenerateHTMLTagPrefix() string {
	s := "[?]"
	if b.ReferenceIndex >= 1 {
		s = fmt.Sprintf("[%v]", b.ReferenceIndex)
	}
	return fmt.Sprintf(
		"<a href=\"%v#ref-%v\" class=\"reference-href\">%v",
		b.File,
		b.ReferenceIndex,
		s,
	)
}

func (b *BlockRef) GenerateHTMLTagSuffix() string {
	return "</a>\n"
}

// =====================================
// Reference HTML interface

func (b *BlockBibliography) GenerateHTMLTagPrefix() string {
	if b.HTMLContent == nil {
		panic(nil)
	}
	return *b.HTMLContent
}

func (b *BlockBibliography) GenerateHTMLTagSuffix() string {
	return ""
}

// =====================================
// Block inline HTML interface

func (b *BlockInline) GenerateHTMLTagPrefix() string {
	return b.Content.GenerateHTMLTagPrefix()
}

func (b *BlockInline) GenerateHTMLTagSuffix() string {
	return b.Content.GenerateHTMLTagSuffix()
}

// =====================================
// Inline document HTML interface

func (b *InlineDocument) GenerateHTMLTagPrefix() string {
	return ""
}

func (b *InlineDocument) GenerateHTMLTagSuffix() string {
	return ""
}

// =====================================
// Raw string HTML interface

func (b *InlineRawString) GenerateHTMLTagPrefix() string {
	return StringToHTMLSafe(b.Content)
}

func (b *InlineRawString) GenerateHTMLTagSuffix() string {
	return ""
}

// =====================================
// String modifier HTML interface

func (b *InlineStringModifier) GenerateHTMLTagPrefix() string {
	switch b.TypeOfModifier {
	case ItalicText:
		return "<em>"
	case BoldText:
		return "<strong>"
	case StrikeoutText:
		return "<del>"
	default:
		panic(nil) // This should never be reached
	}
}

func (b *InlineStringModifier) GenerateHTMLTagSuffix() string {
	switch b.TypeOfModifier {
	case ItalicText:
		return "</em>"
	case BoldText:
		return "</strong>"
	case StrikeoutText:
		return "</del>"
	default:
		panic(nil) // This should never be reached
	}
}

// =====================================
// Delimiter HTML interface

func (b *InlineStringDelimiter) GenerateHTMLTagPrefix() string {
	// Warn the user that something kind of went wrong
	log.Warnf("When generating the HTML, there should be no leftover InlineStringDelimiter (%v). This is a bug and should be reported!", b)
	return b.String()
}

func (b *InlineStringDelimiter) GenerateHTMLTagSuffix() string {
	return ""
}

// =====================================
// Href HTML interface

func (b *InlineHref) GenerateHTMLTagPrefix() string {
	return fmt.Sprintf("<a href=\"%v\">", b.Address)
}

func (b *InlineHref) GenerateHTMLTagSuffix() string {
	return "</a>"
}

// =====================================
// Meta HTML interfaces (ignore them)

func (b *BlockMeta) GenerateHTMLTagPrefix() string {
	return ""
}

func (b *BlockMeta) GenerateHTMLTagSuffix() string {
	return ""
}

// =====================================
// Meta (author) HTML interfaces (ignore them)

func (b *BlockMetaAuthor) GenerateHTMLTagPrefix() string {
	return ""
}

func (b *BlockMetaAuthor) GenerateHTMLTagSuffix() string {
	return ""
}

// =====================================
// Meta (copyright) HTML interfaces (ignore them)

func (b *BlockMetaCopyright) GenerateHTMLTagPrefix() string {
	return ""
}

func (b *BlockMetaCopyright) GenerateHTMLTagSuffix() string {
	return ""
}

// =====================================
// Meta (bibliography) HTML interfaces (ignore them)

func (b *BlockMetaBibinfo) GenerateHTMLTagPrefix() string {
	return ""
}

func (b *BlockMetaBibinfo) GenerateHTMLTagSuffix() string {
	return ""
}

// =====================================
// Generate HTML

func GenerateHTML(tree *Tree[BlockInterface]) string {
	if tree == nil {
		return "" // Just to be sure
	}
	var s strings.Builder
	var GenerateHTMLHelper func (tree *Tree[BlockInterface], sb *strings.Builder)
	GenerateHTMLHelper = func (tree *Tree[BlockInterface], sb *strings.Builder) {
		sb.WriteString(tree.Value.GenerateHTMLTagPrefix())
		for i := 0; i < len(tree.Children); i++ {
			GenerateHTMLHelper(tree.Children[i], sb)
		}
		sb.WriteString(tree.Value.GenerateHTMLTagSuffix())
	}
	GenerateHTMLHelper(tree, &s)
	return s.String()
}
