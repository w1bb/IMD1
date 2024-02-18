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
	"reflect"
)

// =====================================
// Generic block structure

type BlockStruct struct {
	Start Pair
	End Pair
	ContentStart Pair
	ContentEnd Pair
}

func (b BlockStruct) String() string {
	return fmt.Sprintf("{S=%v, E=%v, CS=%v, CE=%v}", b.Start, b.End, b.ContentStart, b.ContentEnd)
}

func (b BlockStruct) Empty() bool {
	return b.Start == b.End
}

// =====================================
// Generic block interface

type BlockInterface interface {
	fmt.Stringer
	HTMLInterface
	CheckBlockStarts(line LineStruct) bool
	SeekBufferAfterBlockStarts() int
	ExecuteAfterBlockStarts(line *LineStruct)
	
	CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int)
	CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool
	ExecuteAfterBlockEnds(line *LineStruct)

	SeekBufferAfterBlockEnds() int
	GetBlocksAllowedInside() []BlockInterface
	AcceptBlockInside(other BlockInterface) bool

	IsPartOfParagraph() bool
	DigDeeperForParagraphs() bool

	GetBlockStruct() *BlockStruct
	GetRawContent() *string
}

// =====================================
// Document

type BlockDocument struct {
	BlockStruct
}

func (b BlockDocument) String() string {
	return "BlockDocument" // irrelevant
}

func (b *BlockDocument) CheckBlockStarts(line LineStruct) bool {
	return false // irrelevant
}

func (b BlockDocument) SeekBufferAfterBlockStarts() int {
	return 0 // irrelevant
}

func (b *BlockDocument) ExecuteAfterBlockStarts(line *LineStruct) {
	// irrelevant
}

func (b *BlockDocument) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	return false, nil, 0 // irrelevant
}

func (b BlockDocument) SeekBufferAfterBlockEnds() int {
	return 0 // irrelevant
}

func (b BlockDocument) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return false // irrelevant
}

func (b *BlockDocument) ExecuteAfterBlockEnds(line *LineStruct) {
	// irrelevant
}

func (b BlockDocument) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface {
		&BlockHtml{},
		&BlockCodeListing{},
		&BlockInlineCodeListing{},
		&BlockMath{},
		&BlockInlineMath{},
		&BlockTextbox{},
		&BlockFigure{},
		&BlockUl{},
		&BlockOl{},
		&BlockFootnote{},
		&BlockRef{},
		&BlockHeading{},
		&BlockMeta{},
	}
}

func (b BlockDocument) AcceptBlockInside(other BlockInterface) bool {
	return true
}

func (b BlockDocument) IsPartOfParagraph() bool {
	return false
}

func (b BlockDocument) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockDocument) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct // irrelevant?
}

func (b *BlockDocument) GetRawContent() *string {
	return nil
}

// =====================================
// Paragraph

// Please note that the paragraphs are inserted only once everything else has been
// inserted (with the exception of InlineBlock)

type BlockParagraph struct {
	BlockStruct
}

func (b BlockParagraph) String() string {
	return fmt.Sprintf(
		"BlockParagraph, %v",
		b.BlockStruct.String(),
	)
}

func (b *BlockParagraph) CheckBlockStarts(line LineStruct) bool {
	return false // irrelevant
}

func (b BlockParagraph) SeekBufferAfterBlockStarts() int {
	return 0 // irrelevant
}

func (b *BlockParagraph) ExecuteAfterBlockStarts(line *LineStruct) {
	// irrelevant
}

func (b *BlockParagraph) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	return false, nil, 0// irrelevant
}

func (b BlockParagraph) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return false // irrelevant
}

func (b *BlockParagraph) ExecuteAfterBlockEnds(line *LineStruct) {
	// irrelevant
}

func (b BlockParagraph) SeekBufferAfterBlockEnds() int {
	return 0 // irrelevant
}

func (b BlockParagraph) GetBlocksAllowedInside() []BlockInterface {
	return nil // irrelevant
}

func (b BlockParagraph) AcceptBlockInside(other BlockInterface) bool {
	return false // irrelevant
}

func (b BlockParagraph) IsPartOfParagraph() bool {
	return false // irrelevant
}

func (b BlockParagraph) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockParagraph) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockParagraph) GetRawContent() *string {
	return nil
}

// =====================================
// Heading

type BlockHeading struct {
	BlockStruct
	HeadingLevel int
	Anchor string
}

func (b BlockHeading) String() string {
	return fmt.Sprintf(
		"BlockHeading (level=%v, anchor=%v), %v",
		b.HeadingLevel,
		b.Anchor,
		b.BlockStruct.String(),
	)
}

func (b *BlockHeading) CheckBlockStarts(line LineStruct) bool {
	if !CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "#") {
		return false
	}
	for line.RuneContent[line.RuneJ + b.HeadingLevel] == '#' {
		b.HeadingLevel++
	}
	return true
}

func (b BlockHeading) SeekBufferAfterBlockStarts() int {
	return b.HeadingLevel
}

func (b *BlockHeading) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair{
		i: line.LineIndex,
		j: line.RuneJ - b.HeadingLevel,
	}
	options := GatherBlockOptions(line, []string{"anchor"})
	if value, ok := options["anchor"]; ok {
		b.Anchor = value
	}
	b.ContentStart = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockHeading) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	return (line.LineIndex != b.Start.i), nil, 0
}

func (b BlockHeading) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return false
}

func (b *BlockHeading) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = b.End
}

func (b BlockHeading) SeekBufferAfterBlockEnds() int {
	return 0
}

func (b BlockHeading) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface {
		&BlockHtml{},
		&BlockInlineCodeListing{},
		&BlockInlineMath{},
		&BlockFootnote{},
		&BlockRef{},
	}
}

func (b BlockHeading) AcceptBlockInside(other BlockInterface) bool {
	return true
}

func (b BlockHeading) IsPartOfParagraph() bool {
	return false
}

func (b BlockHeading) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockHeading) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockHeading) GetRawContent() *string {
	return nil
}

// =====================================
// Textbox

type BlockTextbox struct {
	BlockStruct
	Class string
	Title string
}

func (b BlockTextbox) String() string {
	return fmt.Sprintf(
		"BlockTextbox (class=%v, title=\"%v\"), %v",
		b.Class,
		b.Title,
		b.BlockStruct.String(),
	)
}

func (b *BlockTextbox) CheckBlockStarts(line LineStruct) bool {
	if !CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|") {
		return false
	}
	return CheckRunesStartsWithASCII(line.RuneContent[line.RuneJ:], "|textbox>")
}

func (b BlockTextbox) SeekBufferAfterBlockStarts() int {
	return 9
}

func (b *BlockTextbox) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 9,
	}
	options := GatherBlockOptions(line, []string{"class", "title"})
	if value, ok := options["class"]; ok {
		b.Class = value
	}
	if value, ok := options["title"]; ok {
		b.Title = value
	}
	b.ContentStart = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockTextbox) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<textbox|"), nil, 0
}

func (b BlockTextbox) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return false
}

func (b *BlockTextbox) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 9,
	}
}

func (b BlockTextbox) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b BlockTextbox) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface {
		&BlockHtml{},
		&BlockCodeListing{},
		&BlockInlineCodeListing{},
		&BlockMath{},
		&BlockInlineMath{},
		&BlockTextbox{},
		&BlockFigure{},
		&BlockUl{},
		&BlockOl{},
		&BlockFootnote{},
		&BlockRef{},
	}
}

func (b BlockTextbox) AcceptBlockInside(other BlockInterface) bool {
	return true
}

func (b BlockTextbox) IsPartOfParagraph() bool {
	return false
}

func (b BlockTextbox) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockTextbox) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockTextbox) GetRawContent() *string {
	return nil
}

// =====================================
// HTML code

type BlockHtml struct {
	BlockStruct
	RawContent string
}

func (b BlockHtml) String() string {
	return fmt.Sprintf(
		"BlockHtml, %v :: \"%v\"",
		b.BlockStruct.String(),
		b.RawContent,
	)
}

func (b *BlockHtml) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|html>")
}

func (b BlockHtml) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockHtml) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 6,
	}
	b.ContentStart = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockHtml) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<html|"), nil, 0
}

func (b BlockHtml) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return false
}

func (b *BlockHtml) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 6,
	}
}

func (b BlockHtml) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b BlockHtml) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b BlockHtml) AcceptBlockInside(other BlockInterface) bool {
	return false // irrelevant
}

func (b BlockHtml) IsPartOfParagraph() bool {
	return false
}

func (b BlockHtml) DigDeeperForParagraphs() bool {
	return false
}

func (b *BlockHtml) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockHtml) GetRawContent() *string {
	return &b.RawContent
}

// =====================================
// Code listings

type BlockCodeListing struct {
	BlockStruct
	Language string
	Filename string
	TextAlign string
	AllowCopy bool
	RawContent string
}

func (b BlockCodeListing) String() string {
	return fmt.Sprintf(
		"BlockCodeListing (lang=%v, file=%v, align=%v), %v :: \"%v\"",
		b.Language,
		b.Filename,
		b.TextAlign,
		b.BlockStruct.String(),
		b.RawContent,
	)
}

func (b *BlockCodeListing) CheckBlockStarts(line LineStruct) bool {
	if CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "`") {
		if line.RuneJ + 2 >= len(line.RuneContent) {
			return false
		}
		return line.RuneContent[line.RuneJ+1] == rune('`') && line.RuneContent[line.RuneJ+2] == rune('`')
	}
	return false
}

func (b BlockCodeListing) SeekBufferAfterBlockStarts() int {
	return 3
}

func (b *BlockCodeListing) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 3,
	}
	options := GatherBlockOptions(line, []string{"lang", "file", "align", "copy"})
	b.Language = "plaintext"
	b.AllowCopy = true
	if value, ok := options["class"]; ok {
		if value == "text" || value == "txt" {
			b.Language = "plaintext"
		} else {
			b.Language = value
		}
	}
	if value, ok := options["title"]; ok {
		b.Filename = value
	}
	if value, ok := options["align"]; ok {
		b.TextAlign = value
	}
	if value, ok := options["copy"]; ok {
		b.AllowCopy = Contains([]string{"allow", "allowed", "1", "true", "ok", "yes"}, value)
	}
	b.ContentStart = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockCodeListing) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "```"), nil, 0
}

func (b BlockCodeListing) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return false
}

func (b *BlockCodeListing) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 3,
	}
}

func (b BlockCodeListing) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b BlockCodeListing) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b BlockCodeListing) AcceptBlockInside(other BlockInterface) bool {
	return false // irrelevant
}

func (b BlockCodeListing) IsPartOfParagraph() bool {
	return false
}

func (b BlockCodeListing) DigDeeperForParagraphs() bool {
	return false
}

func (b *BlockCodeListing) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockCodeListing) GetRawContent() *string {
	return &b.RawContent
}

// =====================================
// Inline code listings

type BlockInlineCodeListing struct {
	BlockStruct
	RawContent string
}

func (b BlockInlineCodeListing) String() string {
	return fmt.Sprintf(
		"BlockInlineCodeListing, %v :: \"%v\"",
		b.BlockStruct.String(),
		b.RawContent,
	)
}

func (b *BlockInlineCodeListing) CheckBlockStarts(line LineStruct) bool {
	if CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "`") {
		if line.RuneJ + 2 >= len(line.RuneContent) {
			return true
		}
		return line.RuneContent[line.RuneJ+1] != rune('`') || line.RuneContent[line.RuneJ+2] != rune('`')
	}
	return false
}

func (b BlockInlineCodeListing) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockInlineCodeListing) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair{
		i: line.LineIndex,
		j: line.RuneJ-1,
	}
	b.ContentStart = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockInlineCodeListing) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "`"), nil, 0
}

func (b BlockInlineCodeListing) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return false
}

func (b *BlockInlineCodeListing) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair{
		i: line.LineIndex,
		j: line.RuneJ-1,
	}
}

func (b BlockInlineCodeListing) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b BlockInlineCodeListing) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b BlockInlineCodeListing) AcceptBlockInside(other BlockInterface) bool {
	return false // irrelevant
}

func (b BlockInlineCodeListing) IsPartOfParagraph() bool {
	return true
}

func (b BlockInlineCodeListing) DigDeeperForParagraphs() bool {
	return false
}

func (b *BlockInlineCodeListing) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockInlineCodeListing) GetRawContent() *string {
	return &b.RawContent
}

// =====================================
// Math blocks

type BlockMathType uint8

const (
	DoubleDollar BlockMathType = iota
	Brackets
	BeginEquation
	BeginAlign
)

func (t BlockMathType) String() string {
	switch t {
	case DoubleDollar:
		return "DoubleDollar <=> $$...$$"
	case Brackets:
		return "Brackets <=> \\[...\\]"
	case BeginEquation:
		return "BeginEquation <=> \\begin{equation}...\\end{equation}"
	case BeginAlign:
		return "BeginAlign <=> \\begin{align}...\\end{align}"
	default:
		panic(nil) // This should never be reached
	}
}

type BlockMath struct {
	BlockStruct
	TypeOfBlock BlockMathType
	RawContent string
}

func (b BlockMath) String() string {
	return fmt.Sprintf(
		"BlockMath (type: %v), %v :: \"%v\"",
		b.TypeOfBlock.String(),
		b.BlockStruct.String(),
		b.RawContent,
	)
}

func (b *BlockMath) CheckBlockStarts(line LineStruct) bool {
	s := line.RuneContent[:line.RuneJ+1]
	if CheckRunesEndWithUnescapedASCII(s, "\\begin{equation}") {
		b.TypeOfBlock = BeginEquation
		return true
	} else if CheckRunesEndWithUnescapedASCII(s, "\\begin{align}") {
		b.TypeOfBlock = BeginAlign
		return true
	} else if CheckRunesEndWithUnescapedASCII(s, "\\[") {
		b.TypeOfBlock = Brackets
		return true
	} else if CheckRunesEndWithUnescapedASCII(s, "$") {
		if line.RuneJ + 1 < len(line.RuneContent) && line.RuneContent[line.RuneJ+1] == '$' {
			b.TypeOfBlock = DoubleDollar
			return true
		} else {
			return false
		}
	}
	return false
}

func (b BlockMath) SeekBufferAfterBlockStarts() int {
	switch b.TypeOfBlock {
	case BeginEquation, BeginAlign, Brackets:
		return 1
	case DoubleDollar:
		return 2
	}
	panic(nil)
}

func (b *BlockMath) ExecuteAfterBlockStarts(line *LineStruct) {
	switch b.TypeOfBlock {
	case BeginEquation:
		b.Start = Pair{line.LineIndex, line.RuneJ-16}
	case BeginAlign:
		b.Start = Pair{line.LineIndex, line.RuneJ-13}
	case Brackets, DoubleDollar:
		b.Start = Pair{line.LineIndex, line.RuneJ-2}
	}
	b.ContentStart = Pair{line.LineIndex, line.RuneJ}
}

func (b *BlockMath) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	s := line.RuneContent[:line.RuneJ+1]
	switch b.TypeOfBlock {
	case BeginEquation:
		return CheckRunesEndWithUnescapedASCII(s, "\\end{equation}"), nil, 0
	case BeginAlign:
		return CheckRunesEndWithUnescapedASCII(s, "\\end{align}"), nil, 0
	case Brackets:
		return CheckRunesEndWithUnescapedASCII(s, "\\]"), nil, 0
	case DoubleDollar:
		return CheckRunesEndWithUnescapedASCII(s, "$$"), nil, 0
	}
	panic(nil) // This should never be reached
}

func (b BlockMath) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return false
}

func (b *BlockMath) ExecuteAfterBlockEnds(line *LineStruct) {
	switch b.TypeOfBlock {
	case BeginEquation:
		b.ContentEnd = Pair{
			i: line.LineIndex,
			j: line.RuneJ-14,
		}
	case BeginAlign:
		b.ContentEnd = Pair{
			i: line.LineIndex,
			j: line.RuneJ-11,
		}
	case Brackets, DoubleDollar:
		b.ContentEnd = Pair{
			i: line.LineIndex,
			j: line.RuneJ-2,
		}
	}
	b.End = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b BlockMath) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b BlockMath) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b BlockMath) AcceptBlockInside(other BlockInterface) bool {
	return false // irrelevant
}

func (b BlockMath) IsPartOfParagraph() bool {
	return false
}

func (b BlockMath) DigDeeperForParagraphs() bool {
	return false
}

func (b *BlockMath) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockMath) GetRawContent() *string {
	return &b.RawContent
}

// =====================================
// Inline math blocks

type BlockInlineMathType uint8

const (
	SingleDollar BlockInlineMathType = iota
	Parenthesis
)

func (t BlockInlineMathType) String() string {
	switch t {
	case SingleDollar:
		return "SingleDollar <=> $...$"
	case Parenthesis:
		return "Parenthesis <=> \\(...\\)"
	default:
		panic(nil) // This should never be reached
	}
}

type BlockInlineMath struct {
	BlockStruct
	TypeOfBlock BlockInlineMathType
	RawContent string
}

func (b BlockInlineMath) String() string {

	return fmt.Sprintf(
		"BlockInlineMath (type: %v), %v :: \"%v\"",
		b.TypeOfBlock.String(),
		b.BlockStruct.String(),
		b.RawContent,
	)
}

func (b *BlockInlineMath) CheckBlockStarts(line LineStruct) bool {
	s := line.RuneContent[:line.RuneJ+1]
	if CheckRunesEndWithUnescapedASCII(s, "\\(") {
		b.TypeOfBlock = Parenthesis
		return true
	} else if CheckRunesEndWithUnescapedASCII(s, "$") {
		if line.RuneJ + 1 < len(line.RuneContent) && line.RuneContent[line.RuneJ+1] == '$' {
			return false
		} else {
			b.TypeOfBlock = SingleDollar
			return true
		}
	}
	return false
}

func (b BlockInlineMath) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockInlineMath) ExecuteAfterBlockStarts(line *LineStruct) {
	switch b.TypeOfBlock {
	case Parenthesis:
		b.Start = Pair{
			i: line.LineIndex,
			j: line.RuneJ-2,
		}
	case SingleDollar:
		b.Start = Pair{
			i: line.LineIndex,
			j: line.RuneJ-1,
		}
	}
	b.ContentStart = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockInlineMath) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	s := line.RuneContent[:line.RuneJ+1]
	switch b.TypeOfBlock {
	case Parenthesis:
		return CheckRunesEndWithUnescapedASCII(s, "\\)"), nil, 0
	case SingleDollar:
		return CheckRunesEndWithUnescapedASCII(s, "$"), nil, 0
	}
	panic(nil) // This should never be reached
}

func (b BlockInlineMath) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return false
}

func (b *BlockInlineMath) ExecuteAfterBlockEnds(line *LineStruct) {
	switch b.TypeOfBlock {
	case Parenthesis:
		b.ContentEnd = Pair{
			i: line.LineIndex,
			j: line.RuneJ-2,
		}
	case SingleDollar:
		b.ContentEnd = Pair{
			i: line.LineIndex,
			j: line.RuneJ-1,
		}
	}
	b.End = Pair{
		line.LineIndex,
		line.RuneJ,
	}
}

func (b BlockInlineMath) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b BlockInlineMath) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b BlockInlineMath) AcceptBlockInside(other BlockInterface) bool {
	return false // irrelevant
}

func (b BlockInlineMath) IsPartOfParagraph() bool {
	return true
}

func (b BlockInlineMath) DigDeeperForParagraphs() bool {
	return false
}

func (b *BlockInlineMath) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockInlineMath) GetRawContent() *string {
	return &b.RawContent
}

// =====================================
// Unnumbered lists

type BlockUl struct {
	BlockStruct
	Indentation uint16
}

func (b BlockUl) String() string {
	return fmt.Sprintf(
		"BlockUl (indentation=%v), %v",
		b.Indentation,
		b.BlockStruct.String(),
	)
}

func (b *BlockUl) CheckBlockStarts(line LineStruct) bool {
	b.Indentation = line.Indentation
	return line.RuneJ == 0 && line.RuneContent[line.RuneJ] == '-'
}

func (b BlockUl) SeekBufferAfterBlockStarts() int {
	return 0 // This will in fact be a li
}

func (b *BlockUl) ExecuteAfterBlockStarts(line *LineStruct) {
	b.ContentStart = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.Start = b.ContentStart
}

func (b *BlockUl) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	return false, nil, 0
}

func (b BlockUl) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return NewLines >= 1
}

func (b *BlockUl) ExecuteAfterBlockEnds(line *LineStruct) {
	b.ContentEnd = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.End = b.ContentEnd
}

func (b BlockUl) SeekBufferAfterBlockEnds() int {
	return 0
}

func (b BlockUl) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface {
		&BlockUlLi{},
	}
}

func (b BlockUl) AcceptBlockInside(other BlockInterface) bool {
	if reflect.TypeOf(other) != reflect.TypeOf(BlockUlLi{}) {
		return true
	}
	return b.Indentation == other.(*BlockUlLi).Indentation
}

func (b BlockUl) IsPartOfParagraph() bool {
	return false
}

func (b BlockUl) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockUl) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockUl) GetRawContent() *string {
	return nil
}

// =====================================
// Unnumbered list items

type BlockUlLi struct {
	BlockStruct
	Indentation uint16
	LineIndex int
}

func (b BlockUlLi) String() string {
	return fmt.Sprintf(
		"BlockUlLi (indentation=%v, line-index=%v), %v",
		b.Indentation,
		b.LineIndex,
		b.BlockStruct.String(),
	)
}

func (b *BlockUlLi) CheckBlockStarts(line LineStruct) bool {
	b.Indentation = line.Indentation
	b.LineIndex = line.LineIndex
	return line.RuneJ == 0 && line.RuneContent[line.RuneJ] == '-'
}

func (b BlockUlLi) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockUlLi) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair{
		i: line.LineIndex,
		j: line.RuneJ-1,
	}
	b.ContentStart = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockUlLi) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	// Ignore the same li
	if line.LineIndex == b.LineIndex {
		return false, nil, 0
	}
	// Check if a new ordered list could begin instead
	if aux := (BlockOlLi{}); aux.CheckBlockStarts(*line) && aux.Indentation <= b.Indentation {
		return true, nil, 0
	}
	// A li can end when a new one starts. Make sure to accept nested lists
	if aux := (BlockUlLi{}); aux.CheckBlockStarts(*line) {
		return aux.Indentation <= b.Indentation, nil, 0
	}
	// Different indentation
	return line.Indentation != b.Indentation + 2, nil, 0
}

func (b BlockUlLi) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return NewLines >= 1 && Indentation != b.Indentation + 2
}

func (b *BlockUlLi) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = b.End
}

func (b BlockUlLi) SeekBufferAfterBlockEnds() int {
	return 0
}

func (b BlockUlLi) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface {
		&BlockHtml{},
		&BlockCodeListing{},
		&BlockInlineCodeListing{},
		&BlockMath{},
		&BlockInlineMath{},
		&BlockTextbox{},
		&BlockFigure{},
		&BlockUl{},
		&BlockOl{},
		&BlockFootnote{},
		&BlockRef{},
	}
}

func (b BlockUlLi) AcceptBlockInside(other BlockInterface) bool {
	return true
}

func (b BlockUlLi) IsPartOfParagraph() bool {
	return false
}

func (b BlockUlLi) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockUlLi) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockUlLi) GetRawContent() *string {
	return nil
}

// =====================================
// Numbered lists

type BlockOlType uint8

const (
	OlType_1 BlockOlType = iota
	OlType_A
	OlType_a
	OlType_I
	OlType_i
)

func (t BlockOlType) String() string {
	switch t {
	case OlType_1:
		return "OlType_1 <=> 1. ..."
	case OlType_A:
		return "OlType_A <=> A. ..."
	case OlType_a:
		return "OlType_a <=> a. ..."
	case OlType_I:
		return "OlType_I <=> I. ..."
	case OlType_i:
		return "OlType_i <=> i. ..."
	default:
		panic(nil) // This should never be reached
	}
}

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

type BlockOl struct {
	BlockStruct
	Indentation uint16
	TypeOfBlock BlockOlType
}

func (b BlockOl) String() string {
	return fmt.Sprintf(
		"BlockOl (indentation=%v, type=%v), %v",
		b.Indentation,
		b.TypeOfBlock,
		b.BlockStruct.String(),
	)
}

func (b *BlockOl) CheckBlockStarts(line LineStruct) bool {
	b.Indentation = line.Indentation
	if line.RuneJ != 1 {
		return false
	}
	if line.RuneContent[line.RuneJ - 1] >= rune('1') && line.RuneContent[line.RuneJ - 1] <= rune('9') {
		b.TypeOfBlock = OlType_1
	} else if line.RuneContent[line.RuneJ - 1] == rune('I') {
		b.TypeOfBlock = OlType_I
	} else if line.RuneContent[line.RuneJ - 1] == rune('i') {
		b.TypeOfBlock = OlType_i
	} else if line.RuneContent[line.RuneJ - 1] >= rune('A') && line.RuneContent[line.RuneJ - 1] <= rune('Z') { // excepts I
		b.TypeOfBlock = OlType_A
	} else if line.RuneContent[line.RuneJ - 1] >= rune('a') && line.RuneContent[line.RuneJ - 1] <= rune('z') { // excepts i
		b.TypeOfBlock = OlType_a
	} else {
		return false
	}
	return line.RuneContent[line.RuneJ] == '.'
}

func (b BlockOl) SeekBufferAfterBlockStarts() int {
	return 0 // This will in fact be a li
}

func (b *BlockOl) ExecuteAfterBlockStarts(line *LineStruct) {
	b.ContentStart = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 1,
	}
	b.Start = b.ContentStart
}

func (b *BlockOl) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	return false, nil, 0
}

func (b BlockOl) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return NewLines >= 1
}

func (b *BlockOl) ExecuteAfterBlockEnds(line *LineStruct) {
	b.ContentEnd = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.End = b.ContentEnd
}

func (b BlockOl) SeekBufferAfterBlockEnds() int {
	return 0
}

func (b BlockOl) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface {
		&BlockOlLi{},
	}
}

func (b BlockOl) AcceptBlockInside(other BlockInterface) bool {
	if reflect.TypeOf(other) != reflect.TypeOf(&BlockOlLi{}) {
		return true
	}
	if b.Indentation != other.(*BlockOlLi).Indentation {
		return false
	}
	return b.TypeOfBlock == other.(*BlockOlLi).TypeOfBlock
}

func (b BlockOl) IsPartOfParagraph() bool {
	return false
}

func (b BlockOl) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockOl) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockOl) GetRawContent() *string {
	return nil
}

// =====================================
// Unnumbered list items

type BlockOlLi struct {
	BlockStruct
	Indentation uint16
	LineIndex int
	TypeOfBlock BlockOlType
}

func (b BlockOlLi) String() string {
	return fmt.Sprintf(
		"BlockOlLi (indentation=%v, line-index=%v, type=%v), %v",
		b.Indentation,
		b.LineIndex,
		b.TypeOfBlock,
		b.BlockStruct.String(),
	)
}

func (b *BlockOlLi) CheckBlockStarts(line LineStruct) bool {
	b.Indentation = line.Indentation
	b.LineIndex = line.LineIndex
	if line.RuneJ != 1 {
		return false
	}
	if line.RuneContent[line.RuneJ - 1] >= rune('1') && line.RuneContent[line.RuneJ - 1] <= rune('9') {
		b.TypeOfBlock = OlType_1
	} else if line.RuneContent[line.RuneJ - 1] == rune('I') {
		b.TypeOfBlock = OlType_I
	} else if line.RuneContent[line.RuneJ - 1] == rune('i') {
		b.TypeOfBlock = OlType_i
	} else if line.RuneContent[line.RuneJ - 1] >= rune('A') && line.RuneContent[line.RuneJ - 1] <= rune('Z') { // excepts I
		b.TypeOfBlock = OlType_A
	} else if line.RuneContent[line.RuneJ - 1] >= rune('a') && line.RuneContent[line.RuneJ - 1] <= rune('z') { // excepts i
		b.TypeOfBlock = OlType_a
	} else {
		return false
	}
	return line.RuneContent[line.RuneJ] == '.'
}

func (b BlockOlLi) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockOlLi) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair{
		i: line.LineIndex,
		j: line.RuneJ-2,
	}
	b.ContentStart = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockOlLi) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	// Ignore the same li
	if line.LineIndex == b.LineIndex {
		return false, nil, 0
	}
	// Check if a new unordered list could begin instead
	if aux := (BlockUlLi{}); aux.CheckBlockStarts(*line) && aux.Indentation <= b.Indentation {
		return true, nil, 0
	}
	// A li can end when a new one starts. Make sure to accept nested lists
	if aux := (BlockOlLi{}); aux.CheckBlockStarts(*line) {
		return aux.Indentation <= b.Indentation || aux.TypeOfBlock != b.TypeOfBlock, nil, 0
	}
	// Different indentation
	return line.Indentation != b.Indentation + 3, nil, 0
}

func (b BlockOlLi) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return NewLines >= 1 && Indentation != b.Indentation + 3
}

func (b *BlockOlLi) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = b.End
}

func (b BlockOlLi) SeekBufferAfterBlockEnds() int {
	return 0
}

func (b BlockOlLi) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface {
		&BlockHtml{},
		&BlockCodeListing{},
		&BlockInlineCodeListing{},
		&BlockMath{},
		&BlockInlineMath{},
		&BlockTextbox{},
		&BlockFigure{},
		&BlockUl{},
		&BlockOl{},
		&BlockFootnote{},
		&BlockRef{},
	}
}

func (b BlockOlLi) AcceptBlockInside(other BlockInterface) bool {
	return true
}

func (b BlockOlLi) IsPartOfParagraph() bool {
	return false
}

func (b BlockOlLi) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockOlLi) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockOlLi) GetRawContent() *string {
	return nil
}

// =====================================
// Figures

type BlockFigure struct {
	BlockStruct
	MaxWidth string
	Dock string
	Padding string
}

func (b BlockFigure) String() string {
	return fmt.Sprintf(
		"BlockFigure (max-width=%v, dock=%v, padding=%v), %v",
		b.MaxWidth,
		b.Dock,
		b.Padding,
		b.BlockStruct.String(),
	)
}

func (b *BlockFigure) CheckBlockStarts(line LineStruct) bool {
	if !CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|") {
		return false
	}
	return CheckRunesStartsWithASCII(line.RuneContent[line.RuneJ:], "|figure>")
}

func (b BlockFigure) SeekBufferAfterBlockStarts() int {
	return 8
}

func (b *BlockFigure) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 8,
	}
	options := GatherBlockOptions(line, []string{"max-width", "dock", "padding"})
	b.Dock = "center"
	if value, ok := options["max-width"]; ok {
		b.MaxWidth = value
	}
	if value, ok := options["dock"]; ok {
		switch value {
		case "top", "dock-top":
			b.Dock = "dock-top"
		case "bot", "bottom", "dock-bot", "dock-bottom":
			b.Dock = "dock-bottom"
		case "center":
		default:
			// TODO - trigger warning
		}
		b.Dock = value
	}
	if value, ok := options["padding"]; ok {
		b.Padding = value
	}
	b.ContentStart = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockFigure) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<figure|"), nil, 0
}

func (b BlockFigure) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return false
}

func (b *BlockFigure) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 8,
	}
}

func (b BlockFigure) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b BlockFigure) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface {
		&BlockHtml{},
		&BlockSubfigure{},
	}
}

func (b BlockFigure) AcceptBlockInside(other BlockInterface) bool {
	return true
}

func (b BlockFigure) IsPartOfParagraph() bool {
	return false
}

func (b BlockFigure) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockFigure) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockFigure) GetRawContent() *string {
	return nil
}

// =====================================
// Subfigures

type BlockSubfigure struct {
	BlockStruct
	Source string
	Padding string
}

func (b BlockSubfigure) String() string {
	return fmt.Sprintf(
		"BlockSubfigure (source=%v, padding=%v), %v",
		b.Source,
		b.Padding,
		b.BlockStruct.String(),
	)
}

func (b *BlockSubfigure) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|subfigure>")
}

func (b BlockSubfigure) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockSubfigure) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 11,
	}
	options := GatherBlockOptions(line, []string{"src", "padding"})
	if value, ok := options["src"]; ok {
		b.Source = value
	}
	if value, ok := options["padding"]; ok {
		b.Padding = value
	}
	b.ContentStart = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockSubfigure) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<subfigure|"), nil, 0
}

func (b *BlockSubfigure) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 11,
	}
}

func (b BlockSubfigure) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return false
}

func (b BlockSubfigure) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b BlockSubfigure) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface {
		&BlockHtml{},
		&BlockCodeListing{},
		&BlockInlineCodeListing{},
		&BlockMath{},
		&BlockInlineMath{},
		&BlockFootnote{},
		&BlockRef{},
	}
}

func (b BlockSubfigure) AcceptBlockInside(other BlockInterface) bool {
	return true
}

func (b BlockSubfigure) IsPartOfParagraph() bool {
	return false
}

func (b BlockSubfigure) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockSubfigure) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockSubfigure) GetRawContent() *string {
	return nil
}

// =====================================
// Footnotes

type BlockFootnote struct {
	BlockStruct
	FootnoteIndex int
}

func (b BlockFootnote) String() string {
	return fmt.Sprintf(
		"BlockFootnote (index=%v), %v",
		b.FootnoteIndex,
		b.BlockStruct.String(),
	)
}

func (b *BlockFootnote) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|footnote>")
}

func (b BlockFootnote) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockFootnote) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 10,
	}
	b.ContentStart = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockFootnote) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<footnote|"), nil, 0
}

func (b BlockFootnote) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return false
}

func (b *BlockFootnote) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 10,
	}
}

func (b BlockFootnote) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b BlockFootnote) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface {
		&BlockHtml{},
		&BlockCodeListing{},
		&BlockInlineCodeListing{},
		&BlockMath{},
		&BlockInlineMath{},
		&BlockUl{},
		&BlockOl{},
		&BlockRef{},
	}
}

func (b BlockFootnote) AcceptBlockInside(other BlockInterface) bool {
	return true
}

func (b BlockFootnote) IsPartOfParagraph() bool {
	return true
}

func (b BlockFootnote) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockFootnote) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockFootnote) GetRawContent() *string {
	return nil
}

// =====================================
// References

type BlockRef struct {
	BlockStruct
	File string
	RawContent string
	ReferenceIndex int
}

func (b BlockRef) String() string {
	return fmt.Sprintf(
		"BlockRef (file=%v), %v :: \"%v\"",
		b.File,
		b.BlockStruct.String(),
		b.RawContent,
	)
}

func (b *BlockRef) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|ref>")
}

func (b BlockRef) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockRef) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 5,
	}
	options := GatherBlockOptions(line, []string{"file"})
	if value, ok := options["file"]; ok {
		b.File = value
	}
	b.ContentStart = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockRef) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<ref|"), nil, 0
}

func (b BlockRef) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return false
}

func (b *BlockRef) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 5,
	}
}

func (b BlockRef) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b BlockRef) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b BlockRef) AcceptBlockInside(other BlockInterface) bool {
	return false // irrelevant
}

func (b BlockRef) IsPartOfParagraph() bool {
	return true
}

func (b BlockRef) DigDeeperForParagraphs() bool {
	return false
}

func (b *BlockRef) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockRef) GetRawContent() *string {
	return &b.RawContent
}

// =====================================
// Meta info

type BlockMeta struct {
	BlockStruct
}

func (b BlockMeta) String() string {
	return fmt.Sprintf(
		"BlockMeta, %v",
		b.BlockStruct.String(),
	)
}

func (b *BlockMeta) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|meta>")
}

func (b BlockMeta) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockMeta) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 6,
	}
	b.ContentStart = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockMeta) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<meta|"), nil, 0
}

func (b BlockMeta) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return false
}

func (b *BlockMeta) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 6,
	}
}

func (b BlockMeta) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b BlockMeta) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface{
		&BlockMetaAuthor{},
		&BlockMetaCopyright{},
	}
}

func (b BlockMeta) AcceptBlockInside(other BlockInterface) bool {
	return true
}

func (b BlockMeta) IsPartOfParagraph() bool {
	return false
}

func (b BlockMeta) DigDeeperForParagraphs() bool {
	return false
}

func (b *BlockMeta) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockMeta) GetRawContent() *string {
	return nil
}

// =====================================
// Meta info - author

type BlockMetaAuthor struct {
	BlockStruct
	RawContent string
}

func (b BlockMetaAuthor) String() string {
	return fmt.Sprintf(
		"BlockMetaAuthor, %v :: \"%v\"",
		b.BlockStruct.String(),
		b.RawContent,
	)
}

func (b *BlockMetaAuthor) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|author>")
}

func (b BlockMetaAuthor) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockMetaAuthor) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 8,
	}
	b.ContentStart = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockMetaAuthor) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<author|"), nil, 0
}

func (b BlockMetaAuthor) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return false
}

func (b *BlockMetaAuthor) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 8,
	}
}

func (b BlockMetaAuthor) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b BlockMetaAuthor) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b BlockMetaAuthor) AcceptBlockInside(other BlockInterface) bool {
	return false // irrelevant
}

func (b BlockMetaAuthor) IsPartOfParagraph() bool {
	return false
}

func (b BlockMetaAuthor) DigDeeperForParagraphs() bool {
	return false
}

func (b *BlockMetaAuthor) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockMetaAuthor) GetRawContent() *string {
	return &b.RawContent
}

// =====================================
// Meta info - author

type BlockMetaCopyright struct {
	BlockStruct
	RawContent string
}

func (b BlockMetaCopyright) String() string {
	return fmt.Sprintf(
		"BlockMetaCopyright, %v :: \"%v\"",
		b.BlockStruct.String(),
		b.RawContent,
	)
}

func (b *BlockMetaCopyright) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|copyright>")
}

func (b BlockMetaCopyright) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockMetaCopyright) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 11,
	}
	b.ContentStart = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockMetaCopyright) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<copyright|"), nil, 0
}

func (b BlockMetaCopyright) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return false
}

func (b *BlockMetaCopyright) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair{
		i: line.LineIndex,
		j: line.RuneJ - 11,
	}
}

func (b BlockMetaCopyright) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b BlockMetaCopyright) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b BlockMetaCopyright) AcceptBlockInside(other BlockInterface) bool {
	return false // irrelevant
}

func (b BlockMetaCopyright) IsPartOfParagraph() bool {
	return false
}

func (b BlockMetaCopyright) DigDeeperForParagraphs() bool {
	return false
}

func (b *BlockMetaCopyright) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockMetaCopyright) GetRawContent() *string {
	return &b.RawContent
}

// =====================================
// Generic inline block

type BlockInline struct {
	Content InlineInterface
}

func (b BlockInline) String() string {
	return fmt.Sprintf(
		"(BlockInline->)%v",
		b.Content.String(),
	)
}

func (b BlockInline) CheckBlockStarts(line LineStruct) bool {
	return false // irrelevant
}

func (b BlockInline) SeekBufferAfterBlockStarts() int {
	return 0 // irrelevant
}

func (b BlockInline) ExecuteAfterBlockStarts(line *LineStruct) {
	// irrelevant
}

func (b BlockInline) CheckBlockEndsNormally(line *LineStruct, parsing_stack ParsingStack) (bool, BlockInterface, int) {
	return false, nil, 0 // irrelevant
}

func (b BlockInline) CheckBlockEndsViaNewLinesAndIndentation(NewLines int, Indentation uint16) bool {
	return false // irrelevant
}

func (b BlockInline) ExecuteAfterBlockEnds(line *LineStruct) {
	// irrelevant
}

func (b BlockInline) SeekBufferAfterBlockEnds() int {
	return 1 // irrelevant
}

func (b BlockInline) GetBlocksAllowedInside() []BlockInterface {
	return nil // irrelevant
}

func (b BlockInline) AcceptBlockInside(other BlockInterface) bool {
	return false // irrelevant
}

func (b BlockInline) IsPartOfParagraph() bool {
	return true // irrelevant
}

func (b BlockInline) DigDeeperForParagraphs() bool {
	return true // irrelevant
}

func (b BlockInline) GetBlockStruct() *BlockStruct {
	return new(BlockStruct) // irrelevant
}

func (b *BlockInline) GetRawContent() *string {
	return b.Content.GetRawContent() // This might always be nil
}
