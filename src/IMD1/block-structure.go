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
	"reflect"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// =====================================
// Generic block structure

type BlockStruct struct {
	Start        Pair[int, int]
	End          Pair[int, int]
	ContentStart Pair[int, int]
	ContentEnd   Pair[int, int]
}

func (b BlockStruct) String() string {
	return fmt.Sprintf(
		"{S=%v, E=%v, CS=%v, CE=%v}",
		b.Start,
		b.End,
		b.ContentStart,
		b.ContentEnd,
	)
}

func (b BlockStruct) Empty() bool {
	return b.Start == b.End
}

type BlockEndDetails struct {
	EndNormally            bool
	EndViaNLI              bool
	DisallowUlLiEndUntilNL bool
	DisallowOlLiEndUntilNL bool
}

// =====================================
// Generic block interface

type BlockInterface interface {
	fmt.Stringer

	HTMLInterface
	LaTeXInterface

	CheckBlockStarts(line LineStruct) bool
	SeekBufferAfterBlockStarts() int
	ExecuteAfterBlockStarts(line *LineStruct)

	CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack)
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

type BlockDocumentType uint8

const (
	BlockDocumentTypeCompleteSpecification BlockDocumentType = iota
	BlockDocumentTypeHTML
	BlockDocumentTypeBody
	BlockDocumentTypeDirect
)

func (t BlockDocumentType) String() string {
	switch t {
	case BlockDocumentTypeCompleteSpecification:
		return "CompleteSpecification"
	case BlockDocumentTypeHTML:
		return "HTML"
	case BlockDocumentTypeBody:
		return "Body"
	case BlockDocumentTypeDirect:
		return "Direct"
	default:
		panic(nil) // This should never be reached
	}
}

type BlockDocument struct {
	BlockStruct
	TypeOfBlock BlockDocumentType
}

func (b *BlockDocument) String() string {
	return fmt.Sprintf(
		"BlockDocument (type=%v)",
		b.TypeOfBlock,
	)
}

func (b *BlockDocument) CheckBlockStarts(_ LineStruct) bool {
	return false // irrelevant
}

func (b *BlockDocument) SeekBufferAfterBlockStarts() int {
	return 0 // irrelevant
}

func (b *BlockDocument) ExecuteAfterBlockStarts(_ *LineStruct) {
	// irrelevant
}

func (b *BlockDocument) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = false
	bed.EndViaNLI = false
}

func (b *BlockDocument) SeekBufferAfterBlockEnds() int {
	return 0 // irrelevant
}

func (b *BlockDocument) ExecuteAfterBlockEnds(_ *LineStruct) {
	// irrelevant
}

func (b *BlockDocument) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface{
		&BlockComment{},
		&BlockHTML{},
		&BlockLaTeX{},
		&BlockTabs{},
		&BlockCodeListing{},
		&BlockInlineCodeListing{},
		&BlockMath{},
		&BlockInlineMath{},
		&BlockTextBox{},
		&BlockFigure{},
		&BlockUl{},
		&BlockOl{},
		&BlockFootnote{},
		&BlockRef{},
		&BlockHeading{},
		&BlockMeta{},
		&BlockBibliography{},
	}
}

func (b *BlockDocument) AcceptBlockInside(_ BlockInterface) bool {
	return true
}

func (b *BlockDocument) IsPartOfParagraph() bool {
	return false
}

func (b *BlockDocument) DigDeeperForParagraphs() bool {
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
// inserted (except InlineBlock)

type BlockParagraph struct {
	BlockStruct
}

func (b *BlockParagraph) String() string {
	return fmt.Sprintf(
		"BlockParagraph, %v",
		b.BlockStruct.String(),
	)
}

func (b *BlockParagraph) CheckBlockStarts(_ LineStruct) bool {
	return false // irrelevant
}

func (b *BlockParagraph) SeekBufferAfterBlockStarts() int {
	return 0 // irrelevant
}

func (b *BlockParagraph) ExecuteAfterBlockStarts(_ *LineStruct) {
	// irrelevant
}

func (b *BlockParagraph) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	// irrelevant
	bed.EndNormally = false
	bed.EndViaNLI = false
}

func (b *BlockParagraph) CheckBlockEndsViaNewLinesAndIndentation(_ int, _ uint16) bool {
	return false // irrelevant
}

func (b *BlockParagraph) ExecuteAfterBlockEnds(_ *LineStruct) {
	// irrelevant
}

func (b *BlockParagraph) SeekBufferAfterBlockEnds() int {
	return 0 // irrelevant
}

func (b *BlockParagraph) GetBlocksAllowedInside() []BlockInterface {
	return nil // irrelevant
}

func (b *BlockParagraph) AcceptBlockInside(_ BlockInterface) bool {
	return false // irrelevant
}

func (b *BlockParagraph) IsPartOfParagraph() bool {
	return false // irrelevant
}

func (b *BlockParagraph) DigDeeperForParagraphs() bool {
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
	Anchor       string
}

func (b *BlockHeading) String() string {
	return fmt.Sprintf(
		"BlockHeading (level=%v, anchor=%v), %v",
		b.HeadingLevel,
		b.Anchor,
		b.BlockStruct.String(),
	)
}

func (b *BlockHeading) CheckBlockStarts(line LineStruct) bool {
	if line.RuneJ != 0 || line.RuneContent[line.RuneJ] != '#' {
		return false
	}
	for (line.RuneJ+b.HeadingLevel) < len(line.RuneContent) && line.RuneContent[line.RuneJ+b.HeadingLevel] == '#' {
		b.HeadingLevel++
	}
	return true
}

func (b *BlockHeading) SeekBufferAfterBlockStarts() int {
	return b.HeadingLevel
}

func (b *BlockHeading) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - b.HeadingLevel,
	}
	options := GatherBlockOptions(line, []string{"anchor"})
	if value, ok := options["anchor"]; ok {
		b.Anchor = value
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockHeading) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = line.LineIndex != b.Start.i
	bed.EndViaNLI = false
}

func (b *BlockHeading) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = b.End
}

func (b *BlockHeading) SeekBufferAfterBlockEnds() int {
	return 0
}

func (b *BlockHeading) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface{
		&BlockComment{},
		&BlockHTML{},
		&BlockLaTeX{},
		&BlockInlineCodeListing{},
		&BlockInlineMath{},
		&BlockFootnote{},
		&BlockRef{},
	}
}

func (b *BlockHeading) AcceptBlockInside(_ BlockInterface) bool {
	return true
}

func (b *BlockHeading) IsPartOfParagraph() bool {
	return false
}

func (b *BlockHeading) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockHeading) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockHeading) GetRawContent() *string {
	return nil
}

// =====================================
// TextBox

type BlockTextBox struct {
	BlockStruct
	Class string
}

func (b *BlockTextBox) String() string {
	return fmt.Sprintf(
		"BlockTextBox (class=%v), %v",
		b.Class,
		b.BlockStruct.String(),
	)
}

func (b *BlockTextBox) CheckBlockStarts(line LineStruct) bool {
	if !CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|") {
		return false
	}
	return CheckRunesStartsWithASCII(line.RuneContent[line.RuneJ:], "|textbox>")
}

func (b *BlockTextBox) SeekBufferAfterBlockStarts() int {
	return 9
}

func (b *BlockTextBox) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 9,
	}
	options := GatherBlockOptions(line, []string{"class"})
	if value, ok := options["class"]; ok {
		b.Class = value
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockTextBox) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<textbox|")
	bed.EndViaNLI = false
}

func (b *BlockTextBox) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 9,
	}
}

func (b *BlockTextBox) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockTextBox) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface{
		&BlockComment{},
		&BlockTextBoxTitle{},
		&BlockTextBoxContent{},
	}
}

func (b *BlockTextBox) AcceptBlockInside(_ BlockInterface) bool {
	return true
}

func (b *BlockTextBox) IsPartOfParagraph() bool {
	return false
}

func (b *BlockTextBox) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockTextBox) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockTextBox) GetRawContent() *string {
	return nil
}

// =====================================
// TextBox title

type BlockTextBoxTitle struct {
	BlockStruct
}

func (b *BlockTextBoxTitle) String() string {
	return fmt.Sprintf(
		"BlockTextBoxTitle, %v",
		b.BlockStruct.String(),
	)
}

func (b *BlockTextBoxTitle) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|title>")
}

func (b *BlockTextBoxTitle) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockTextBoxTitle) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 7,
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockTextBoxTitle) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<title|")
	bed.EndViaNLI = false
}

func (b *BlockTextBoxTitle) CheckBlockEndsViaNewLinesAndIndentation(_ int, _ uint16) bool {
	return false
}

func (b *BlockTextBoxTitle) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 7,
	}
}

func (b *BlockTextBoxTitle) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockTextBoxTitle) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface{
		&BlockComment{},
		&BlockHTML{},
		&BlockLaTeX{},
		&BlockInlineCodeListing{},
		&BlockInlineMath{},
		&BlockFootnote{},
		&BlockRef{},
	}
}

func (b *BlockTextBoxTitle) AcceptBlockInside(_ BlockInterface) bool {
	return true
}

func (b *BlockTextBoxTitle) IsPartOfParagraph() bool {
	return false
}

func (b *BlockTextBoxTitle) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockTextBoxTitle) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockTextBoxTitle) GetRawContent() *string {
	return nil
}

// =====================================
// TextBox content

type BlockTextBoxContent struct {
	BlockStruct
}

func (b *BlockTextBoxContent) String() string {
	return fmt.Sprintf(
		"BlockTextBoxContent, %v",
		b.BlockStruct.String(),
	)
}

func (b *BlockTextBoxContent) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|content>")
}

func (b *BlockTextBoxContent) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockTextBoxContent) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 9,
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockTextBoxContent) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<content|")
	bed.EndViaNLI = false
}

func (b *BlockTextBoxContent) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 9,
	}
}

func (b *BlockTextBoxContent) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockTextBoxContent) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface{
		&BlockComment{},
		&BlockHTML{},
		&BlockLaTeX{},
		&BlockTabs{},
		&BlockCodeListing{},
		&BlockInlineCodeListing{},
		&BlockMath{},
		&BlockInlineMath{},
		&BlockTextBox{},
		&BlockFigure{},
		&BlockUl{},
		&BlockOl{},
		&BlockFootnote{},
		&BlockRef{},
	}
}

func (b *BlockTextBoxContent) AcceptBlockInside(_ BlockInterface) bool {
	return true
}

func (b *BlockTextBoxContent) IsPartOfParagraph() bool {
	return false
}

func (b *BlockTextBoxContent) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockTextBoxContent) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockTextBoxContent) GetRawContent() *string {
	return nil
}

// =====================================
// Toggles

type BlockToggle struct {
	BlockStruct
	Class string
}

func (b *BlockToggle) String() string {
	return fmt.Sprintf(
		"BlockToggle (class=%v), %v",
		b.Class,
		b.BlockStruct.String(),
	)
}

func (b *BlockToggle) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|toggle>")
}

func (b *BlockToggle) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockToggle) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 8,
	}
	options := GatherBlockOptions(line, []string{"class"})
	if value, ok := options["class"]; ok {
		b.Class = value
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockToggle) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<toggle|")
	bed.EndViaNLI = false
}

func (b *BlockToggle) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 8,
	}
}

func (b *BlockToggle) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockToggle) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface{
		&BlockComment{},
		&BlockTextBoxTitle{},
		&BlockTextBoxContent{},
	}
}

func (b *BlockToggle) AcceptBlockInside(_ BlockInterface) bool {
	return true
}

func (b *BlockToggle) IsPartOfParagraph() bool {
	return false
}

func (b *BlockToggle) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockToggle) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockToggle) GetRawContent() *string {
	return nil
}

// =====================================
// HTML code

type BlockComment struct {
	BlockStruct
	RawContent string
}

func (b *BlockComment) String() string {
	return fmt.Sprintf(
		"BlockComment, %v :: \"%v\"",
		b.BlockStruct.String(),
		b.RawContent,
	)
}

func (b *BlockComment) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<!--")
}

func (b *BlockComment) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockComment) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 4,
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockComment) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "-->")
	bed.EndViaNLI = false
}

func (b *BlockComment) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 3,
	}
}

func (b *BlockComment) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockComment) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b *BlockComment) AcceptBlockInside(_ BlockInterface) bool {
	return false // irrelevant
}

func (b *BlockComment) IsPartOfParagraph() bool {
	return false
}

func (b *BlockComment) DigDeeperForParagraphs() bool {
	return false
}

func (b *BlockComment) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockComment) GetRawContent() *string {
	return &b.RawContent
}

// =====================================
// HTML code

type BlockHTML struct {
	BlockStruct
	RawContent string
}

func (b *BlockHTML) String() string {
	return fmt.Sprintf(
		"BlockHTML, %v :: \"%v\"",
		b.BlockStruct.String(),
		b.RawContent,
	)
}

func (b *BlockHTML) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|html>")
}

func (b *BlockHTML) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockHTML) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 6,
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockHTML) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<html|")
	bed.EndViaNLI = false
}

func (b *BlockHTML) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 6,
	}
}

func (b *BlockHTML) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockHTML) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b *BlockHTML) AcceptBlockInside(_ BlockInterface) bool {
	return false // irrelevant
}

func (b *BlockHTML) IsPartOfParagraph() bool {
	return false
}

func (b *BlockHTML) DigDeeperForParagraphs() bool {
	return false
}

func (b *BlockHTML) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockHTML) GetRawContent() *string {
	return &b.RawContent
}

// =====================================
// HTML code

type BlockLaTeX struct {
	BlockStruct
	RawContent string
}

func (b *BlockLaTeX) String() string {
	return fmt.Sprintf(
		"BlockLaTeX, %v :: \"%v\"",
		b.BlockStruct.String(),
		b.RawContent,
	)
}

func (b *BlockLaTeX) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|latex>")
}

func (b *BlockLaTeX) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockLaTeX) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 7,
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockLaTeX) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<latex|")
	bed.EndViaNLI = false
}

func (b *BlockLaTeX) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 7,
	}
}

func (b *BlockLaTeX) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockLaTeX) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b *BlockLaTeX) AcceptBlockInside(_ BlockInterface) bool {
	return false // irrelevant
}

func (b *BlockLaTeX) IsPartOfParagraph() bool {
	return false
}

func (b *BlockLaTeX) DigDeeperForParagraphs() bool {
	return false
}

func (b *BlockLaTeX) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockLaTeX) GetRawContent() *string {
	return &b.RawContent
}

// =====================================
// Code listings

type BlockCodeListing struct {
	BlockStruct
	Language   string
	Filename   string
	TextAlign  string
	AllowCopy  bool
	RawContent string
}

func (b *BlockCodeListing) String() string {
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
		if line.RuneJ+2 >= len(line.RuneContent) {
			return false
		}
		return line.RuneContent[line.RuneJ+1] == '`' && line.RuneContent[line.RuneJ+2] == '`'
	}
	return false
}

func (b *BlockCodeListing) SeekBufferAfterBlockStarts() int {
	return 3
}

func (b *BlockCodeListing) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 3,
	}
	options := GatherBlockOptions(line, []string{"lang", "file", "align", "copy"})
	b.Language = "plaintext"
	b.AllowCopy = true
	if value, ok := options["lang"]; ok {
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
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockCodeListing) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "```")
	bed.EndViaNLI = false
}

func (b *BlockCodeListing) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 3,
	}
}

func (b *BlockCodeListing) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockCodeListing) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b *BlockCodeListing) AcceptBlockInside(_ BlockInterface) bool {
	return false // irrelevant
}

func (b *BlockCodeListing) IsPartOfParagraph() bool {
	return false
}

func (b *BlockCodeListing) DigDeeperForParagraphs() bool {
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

func (b *BlockInlineCodeListing) String() string {
	return fmt.Sprintf(
		"BlockInlineCodeListing, %v :: \"%v\"",
		b.BlockStruct.String(),
		b.RawContent,
	)
}

func (b *BlockInlineCodeListing) CheckBlockStarts(line LineStruct) bool {
	if CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "`") {
		if line.RuneJ+2 >= len(line.RuneContent) {
			return true
		}
		return line.RuneContent[line.RuneJ+1] != '`' || line.RuneContent[line.RuneJ+2] != '`'
	}
	return false
}

func (b *BlockInlineCodeListing) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockInlineCodeListing) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 1,
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockInlineCodeListing) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "`")
	bed.EndViaNLI = false // TODO
}

func (b *BlockInlineCodeListing) CheckBlockEndsViaNewLinesAndIndentation(_ int, _ uint16) bool {
	return false
}

func (b *BlockInlineCodeListing) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 1,
	}
}

func (b *BlockInlineCodeListing) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockInlineCodeListing) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b *BlockInlineCodeListing) AcceptBlockInside(_ BlockInterface) bool {
	return false // irrelevant
}

func (b *BlockInlineCodeListing) IsPartOfParagraph() bool {
	return true
}

func (b *BlockInlineCodeListing) DigDeeperForParagraphs() bool {
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
	BlockMathTypeDoubleDollar BlockMathType = iota
	BlockMathTypeBrackets
	BlockMathTypeBeginEquation
	BlockMathTypeBeginEquationStar
	BlockMathTypeBeginAlign
	BlockMathTypeBeginAlignStar
)

func (t BlockMathType) String() string {
	switch t {
	case BlockMathTypeDoubleDollar:
		return "DoubleDollar <=> $$...$$"
	case BlockMathTypeBrackets:
		return "Brackets <=> \\[...\\]"
	case BlockMathTypeBeginEquation:
		return "BeginEquation <=> \\begin{equation}...\\end{equation}"
	case BlockMathTypeBeginEquationStar:
		return "BeginEquationStar <=> \\begin{equation*}...\\end{equation*}"
	case BlockMathTypeBeginAlign:
		return "BeginAlign <=> \\begin{align}...\\end{align}"
	case BlockMathTypeBeginAlignStar:
		return "BeginAlignStar <=> \\begin{align*}...\\end{align*}"
	default:
		panic(nil) // This should never be reached
	}
}

type BlockMath struct {
	BlockStruct
	TypeOfBlock BlockMathType
	RawContent  string
}

func (b *BlockMath) String() string {
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
		b.TypeOfBlock = BlockMathTypeBeginEquation
		return true
	} else if CheckRunesEndWithUnescapedASCII(s, "\\begin{equation*}") {
		b.TypeOfBlock = BlockMathTypeBeginEquationStar
		return true
	} else if CheckRunesEndWithUnescapedASCII(s, "\\begin{align}") {
		b.TypeOfBlock = BlockMathTypeBeginAlign
		return true
	} else if CheckRunesEndWithUnescapedASCII(s, "\\begin{align*}") {
		b.TypeOfBlock = BlockMathTypeBeginAlignStar
		return true
	} else if CheckRunesEndWithUnescapedASCII(s, "\\[") {
		b.TypeOfBlock = BlockMathTypeBrackets
		return true
	} else if CheckRunesEndWithUnescapedASCII(s, "$") {
		if line.RuneJ+1 < len(line.RuneContent) && line.RuneContent[line.RuneJ+1] == '$' {
			b.TypeOfBlock = BlockMathTypeDoubleDollar
			return true
		} else {
			return false
		}
	}
	return false
}

func (b *BlockMath) SeekBufferAfterBlockStarts() int {
	switch b.TypeOfBlock {
	case BlockMathTypeBeginEquation, BlockMathTypeBeginEquationStar, BlockMathTypeBeginAlign, BlockMathTypeBeginAlignStar, BlockMathTypeBrackets:
		return 1
	case BlockMathTypeDoubleDollar:
		return 2
	}
	panic(nil)
}

func (b *BlockMath) ExecuteAfterBlockStarts(line *LineStruct) {
	switch b.TypeOfBlock {
	case BlockMathTypeBeginEquation:
		b.Start = Pair[int, int]{line.LineIndex, line.RuneJ - 16}
	case BlockMathTypeBeginEquationStar:
		b.Start = Pair[int, int]{line.LineIndex, line.RuneJ - 17}
	case BlockMathTypeBeginAlign:
		b.Start = Pair[int, int]{line.LineIndex, line.RuneJ - 13}
	case BlockMathTypeBeginAlignStar:
		b.Start = Pair[int, int]{line.LineIndex, line.RuneJ - 14}
	case BlockMathTypeBrackets, BlockMathTypeDoubleDollar:
		b.Start = Pair[int, int]{line.LineIndex, line.RuneJ - 2}
	}
	b.ContentStart = Pair[int, int]{line.LineIndex, line.RuneJ}
}

func (b *BlockMath) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = func() bool {
		s := line.RuneContent[:line.RuneJ+1]
		switch b.TypeOfBlock {
		case BlockMathTypeBeginEquation:
			return CheckRunesEndWithUnescapedASCII(s, "\\end{equation}")
		case BlockMathTypeBeginEquationStar:
			return CheckRunesEndWithUnescapedASCII(s, "\\end{equation*}")
		case BlockMathTypeBeginAlign:
			return CheckRunesEndWithUnescapedASCII(s, "\\end{align}")
		case BlockMathTypeBeginAlignStar:
			return CheckRunesEndWithUnescapedASCII(s, "\\end{align*}")
		case BlockMathTypeBrackets:
			return CheckRunesEndWithUnescapedASCII(s, "\\]")
		case BlockMathTypeDoubleDollar:
			return CheckRunesEndWithUnescapedASCII(s, "$$")
		}
		panic(nil) // This should never be reached
	}()
	bed.EndViaNLI = false
}

func (b *BlockMath) CheckBlockEndsNormally(line *LineStruct, _ ParsingStack) (bool, BlockInterface, int) {
	s := line.RuneContent[:line.RuneJ+1]
	switch b.TypeOfBlock {
	case BlockMathTypeBeginEquation:
		return CheckRunesEndWithUnescapedASCII(s, "\\end{equation}"), nil, 0
	case BlockMathTypeBeginEquationStar:
		return CheckRunesEndWithUnescapedASCII(s, "\\end{equation*}"), nil, 0
	case BlockMathTypeBeginAlign:
		return CheckRunesEndWithUnescapedASCII(s, "\\end{align}"), nil, 0
	case BlockMathTypeBeginAlignStar:
		return CheckRunesEndWithUnescapedASCII(s, "\\end{align*}"), nil, 0
	case BlockMathTypeBrackets:
		return CheckRunesEndWithUnescapedASCII(s, "\\]"), nil, 0
	case BlockMathTypeDoubleDollar:
		return CheckRunesEndWithUnescapedASCII(s, "$$"), nil, 0
	}
	panic(nil) // This should never be reached
}

func (b *BlockMath) CheckBlockEndsViaNewLinesAndIndentation(_ int, _ uint16) bool {
	return false
}

func (b *BlockMath) ExecuteAfterBlockEnds(line *LineStruct) {
	switch b.TypeOfBlock {
	case BlockMathTypeBeginEquation:
		b.ContentEnd = Pair[int, int]{
			i: line.LineIndex,
			j: line.RuneJ - 14,
		}
	case BlockMathTypeBeginEquationStar:
		b.ContentEnd = Pair[int, int]{
			i: line.LineIndex,
			j: line.RuneJ - 15,
		}
	case BlockMathTypeBeginAlign:
		b.ContentEnd = Pair[int, int]{
			i: line.LineIndex,
			j: line.RuneJ - 11,
		}
	case BlockMathTypeBeginAlignStar:
		b.ContentEnd = Pair[int, int]{
			i: line.LineIndex,
			j: line.RuneJ - 12,
		}
	case BlockMathTypeBrackets, BlockMathTypeDoubleDollar:
		b.ContentEnd = Pair[int, int]{
			i: line.LineIndex,
			j: line.RuneJ - 2,
		}
	}
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockMath) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockMath) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b *BlockMath) AcceptBlockInside(_ BlockInterface) bool {
	return false // irrelevant
}

func (b *BlockMath) IsPartOfParagraph() bool {
	return false
}

func (b *BlockMath) DigDeeperForParagraphs() bool {
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
	BlockInlineMathTypeSingleDollar BlockInlineMathType = iota
	BlockInlineMathTypeParenthesis
)

func (t BlockInlineMathType) String() string {
	switch t {
	case BlockInlineMathTypeSingleDollar:
		return "SingleDollar <=> $...$"
	case BlockInlineMathTypeParenthesis:
		return "Parenthesis <=> \\(...\\)"
	default:
		panic(nil) // This should never be reached
	}
}

type BlockInlineMath struct {
	BlockStruct
	TypeOfBlock BlockInlineMathType
	RawContent  string
}

func (b *BlockInlineMath) String() string {
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
		b.TypeOfBlock = BlockInlineMathTypeParenthesis
		return true
	} else if CheckRunesEndWithUnescapedASCII(s, "$") {
		if line.RuneJ+1 < len(line.RuneContent) && line.RuneContent[line.RuneJ+1] == '$' {
			return false
		} else {
			b.TypeOfBlock = BlockInlineMathTypeSingleDollar
			return true
		}
	}
	return false
}

func (b *BlockInlineMath) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockInlineMath) ExecuteAfterBlockStarts(line *LineStruct) {
	switch b.TypeOfBlock {
	case BlockInlineMathTypeParenthesis:
		b.Start = Pair[int, int]{
			i: line.LineIndex,
			j: line.RuneJ - 2,
		}
	case BlockInlineMathTypeSingleDollar:
		b.Start = Pair[int, int]{
			i: line.LineIndex,
			j: line.RuneJ - 1,
		}
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockInlineMath) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = func() bool {
		s := line.RuneContent[:line.RuneJ+1]
		switch b.TypeOfBlock {
		case BlockInlineMathTypeParenthesis:
			return CheckRunesEndWithUnescapedASCII(s, "\\)")
		case BlockInlineMathTypeSingleDollar:
			return CheckRunesEndWithUnescapedASCII(s, "$")
		}
		panic(nil) // This should never be reached
	}()
	bed.EndViaNLI = false
	if bed.EndNormally || bed.EndViaNLI {
		bed.DisallowUlLiEndUntilNL = true // experimental
		bed.DisallowOlLiEndUntilNL = true // experimental
	}
}

func (b *BlockInlineMath) ExecuteAfterBlockEnds(line *LineStruct) {
	switch b.TypeOfBlock {
	case BlockInlineMathTypeParenthesis:
		b.ContentEnd = Pair[int, int]{
			i: line.LineIndex,
			j: line.RuneJ - 2,
		}
	case BlockInlineMathTypeSingleDollar:
		b.ContentEnd = Pair[int, int]{
			i: line.LineIndex,
			j: line.RuneJ - 1,
		}
	}
	b.End = Pair[int, int]{
		line.LineIndex,
		line.RuneJ,
	}
}

func (b *BlockInlineMath) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockInlineMath) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b *BlockInlineMath) AcceptBlockInside(_ BlockInterface) bool {
	return false // irrelevant
}

func (b *BlockInlineMath) IsPartOfParagraph() bool {
	return true
}

func (b *BlockInlineMath) DigDeeperForParagraphs() bool {
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

func (b *BlockUl) String() string {
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

func (b *BlockUl) SeekBufferAfterBlockStarts() int {
	return 0 // This will in fact be a li
}

func (b *BlockUl) ExecuteAfterBlockStarts(line *LineStruct) {
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.Start = b.ContentStart
}

func (b *BlockUl) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = false
	bed.EndViaNLI = NewLines >= 1 || Indentation < b.Indentation
}

func (b *BlockUl) ExecuteAfterBlockEnds(line *LineStruct) {
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.End = b.ContentEnd
}

func (b *BlockUl) SeekBufferAfterBlockEnds() int {
	return 0
}

func (b *BlockUl) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface{
		&BlockUlLi{},
	}
}

func (b *BlockUl) AcceptBlockInside(other BlockInterface) bool {
	if reflect.TypeOf(other) != reflect.TypeOf(BlockUlLi{}) {
		return true
	}
	return b.Indentation == other.(*BlockUlLi).Indentation
}

func (b *BlockUl) IsPartOfParagraph() bool {
	return false
}

func (b *BlockUl) DigDeeperForParagraphs() bool {
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
	LineIndex   int
}

func (b *BlockUlLi) String() string {
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

func (b *BlockUlLi) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockUlLi) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 1,
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockUlLi) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	if !bed.DisallowUlLiEndUntilNL {
		if line.LineIndex == b.LineIndex { // Ignore the same li
			bed.EndNormally = false
		} else if aux := (BlockOlLi{}); aux.CheckBlockStarts(*line) && aux.Indentation <= b.Indentation { // Check if a new ordered list could begin instead
			bed.EndNormally = true
		} else if aux := (BlockUlLi{}); aux.CheckBlockStarts(*line) { // A li can end when a new one starts. Make sure to accept nested lis
			bed.EndNormally = aux.Indentation <= b.Indentation
		} else {
			bed.EndNormally = line.Indentation != b.Indentation+2
		}
	} else {
		bed.EndNormally = false
	}
	bed.EndViaNLI = NewLines >= 1 && Indentation != b.Indentation+2
}

func (b *BlockUlLi) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = b.End
}

func (b *BlockUlLi) SeekBufferAfterBlockEnds() int {
	return 0
}

func (b *BlockUlLi) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface{
		&BlockComment{},
		&BlockHTML{},
		&BlockLaTeX{},
		&BlockTabs{},
		&BlockCodeListing{},
		&BlockInlineCodeListing{},
		&BlockMath{},
		&BlockInlineMath{},
		&BlockTextBox{},
		&BlockFigure{},
		&BlockUl{},
		&BlockOl{},
		&BlockFootnote{},
		&BlockRef{},
	}
}

func (b *BlockUlLi) AcceptBlockInside(_ BlockInterface) bool {
	return true
}

func (b *BlockUlLi) IsPartOfParagraph() bool {
	return false
}

func (b *BlockUlLi) DigDeeperForParagraphs() bool {
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
	BlockOlTypeNumber BlockOlType = iota
	BlockOlTypeLetterCapital
	BlockOlTypeLetter
	BlockOlTypeRomanCapital
	BlockOlTypeRoman
)

func (t BlockOlType) String() string {
	switch t {
	case BlockOlTypeNumber:
		return "Number <=> 1. ..."
	case BlockOlTypeLetterCapital:
		return "LetterCapital <=> A. ..."
	case BlockOlTypeLetter:
		return "Letter <=> a. ..."
	case BlockOlTypeRomanCapital:
		return "RomanCapital <=> I. ..."
	case BlockOlTypeRoman:
		return "Roman <=> i. ..."
	default:
		panic(nil) // This should never be reached
	}
}

type BlockOl struct {
	BlockStruct
	Indentation uint16
	TypeOfBlock BlockOlType
}

func (b *BlockOl) String() string {
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
	if line.RuneContent[line.RuneJ-1] >= '1' && line.RuneContent[line.RuneJ-1] <= '9' {
		b.TypeOfBlock = BlockOlTypeNumber
	} else if line.RuneContent[line.RuneJ-1] == 'I' {
		b.TypeOfBlock = BlockOlTypeRomanCapital
	} else if line.RuneContent[line.RuneJ-1] == 'i' {
		b.TypeOfBlock = BlockOlTypeRoman
	} else if line.RuneContent[line.RuneJ-1] >= 'A' && line.RuneContent[line.RuneJ-1] <= 'Z' { // excepts I
		b.TypeOfBlock = BlockOlTypeLetterCapital
	} else if line.RuneContent[line.RuneJ-1] >= 'a' && line.RuneContent[line.RuneJ-1] <= 'z' { // excepts i
		b.TypeOfBlock = BlockOlTypeLetter
	} else {
		return false
	}
	return line.RuneContent[line.RuneJ] == '.'
}

func (b *BlockOl) SeekBufferAfterBlockStarts() int {
	return 0 // This will in fact be a li
}

func (b *BlockOl) ExecuteAfterBlockStarts(line *LineStruct) {
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 1,
	}
	b.Start = b.ContentStart
}

func (b *BlockOl) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = false
	bed.EndViaNLI = NewLines >= 1 || Indentation < b.Indentation
}

func (b *BlockOl) ExecuteAfterBlockEnds(line *LineStruct) {
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.End = b.ContentEnd
}

func (b *BlockOl) SeekBufferAfterBlockEnds() int {
	return 0
}

func (b *BlockOl) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface{
		&BlockOlLi{},
	}
}

func (b *BlockOl) AcceptBlockInside(other BlockInterface) bool {
	if reflect.TypeOf(other) != reflect.TypeOf(&BlockOlLi{}) {
		return true
	}
	if b.Indentation != other.(*BlockOlLi).Indentation {
		return false
	}
	return b.TypeOfBlock == other.(*BlockOlLi).TypeOfBlock
}

func (b *BlockOl) IsPartOfParagraph() bool {
	return false
}

func (b *BlockOl) DigDeeperForParagraphs() bool {
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
	LineIndex   int
	TypeOfBlock BlockOlType
}

func (b *BlockOlLi) String() string {
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
	if line.RuneContent[line.RuneJ-1] >= '1' && line.RuneContent[line.RuneJ-1] <= '9' {
		b.TypeOfBlock = BlockOlTypeNumber
	} else if line.RuneContent[line.RuneJ-1] == 'I' {
		b.TypeOfBlock = BlockOlTypeRomanCapital
	} else if line.RuneContent[line.RuneJ-1] == 'i' {
		b.TypeOfBlock = BlockOlTypeRoman
	} else if line.RuneContent[line.RuneJ-1] >= 'A' && line.RuneContent[line.RuneJ-1] <= 'Z' { // excepts I
		b.TypeOfBlock = BlockOlTypeLetterCapital
	} else if line.RuneContent[line.RuneJ-1] >= 'a' && line.RuneContent[line.RuneJ-1] <= 'z' { // excepts i
		b.TypeOfBlock = BlockOlTypeLetter
	} else {
		return false
	}
	return line.RuneContent[line.RuneJ] == '.'
}

func (b *BlockOlLi) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockOlLi) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 2,
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockOlLi) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	if !bed.DisallowOlLiEndUntilNL {
		if line.LineIndex == b.LineIndex { // Ignore the same li
			bed.EndNormally = false
		} else if aux := (BlockUlLi{}); aux.CheckBlockStarts(*line) && aux.Indentation <= b.Indentation { // Check if a new unordered list could begin instead
			bed.EndNormally = true
		} else if aux := (BlockOlLi{}); aux.CheckBlockStarts(*line) { // A li can end when a new one starts. Make sure to accept nested lists
			bed.EndNormally = aux.Indentation <= b.Indentation || aux.TypeOfBlock != b.TypeOfBlock
		} else { // // Different indentation
			bed.EndNormally = line.Indentation != b.Indentation+3
		}
	} else {
		bed.EndNormally = false
	}
	bed.EndViaNLI = NewLines >= 1 && Indentation != b.Indentation+3
}

func (b *BlockOlLi) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = b.End
}

func (b *BlockOlLi) SeekBufferAfterBlockEnds() int {
	return 0
}

func (b *BlockOlLi) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface{
		&BlockComment{},
		&BlockHTML{},
		&BlockLaTeX{},
		&BlockTabs{},
		&BlockCodeListing{},
		&BlockInlineCodeListing{},
		&BlockMath{},
		&BlockInlineMath{},
		&BlockTextBox{},
		&BlockFigure{},
		&BlockUl{},
		&BlockOl{},
		&BlockFootnote{},
		&BlockRef{},
	}
}

func (b *BlockOlLi) AcceptBlockInside(_ BlockInterface) bool {
	return true
}

func (b *BlockOlLi) IsPartOfParagraph() bool {
	return false
}

func (b *BlockOlLi) DigDeeperForParagraphs() bool {
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
	Dock     string
	Padding  string
}

func (b *BlockFigure) String() string {
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

func (b *BlockFigure) SeekBufferAfterBlockStarts() int {
	return 8
}

func (b *BlockFigure) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
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
			// skip
		default:
			log.Warnf(
				"A figure cannot have \"dock\" set to \"%v\". Resetting it do default (\"%v\")",
				value,
				"center",
			)
		}
	}
	if value, ok := options["padding"]; ok {
		b.Padding = value
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockFigure) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<figure|")
	bed.EndViaNLI = false
}

func (b *BlockFigure) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 8,
	}
}

func (b *BlockFigure) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockFigure) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface{
		&BlockComment{},
		&BlockHTML{},
		&BlockLaTeX{},
		&BlockSubFigure{},
	}
}

func (b *BlockFigure) AcceptBlockInside(_ BlockInterface) bool {
	return true
}

func (b *BlockFigure) IsPartOfParagraph() bool {
	return false
}

func (b *BlockFigure) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockFigure) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockFigure) GetRawContent() *string {
	return nil
}

// =====================================
// SubFigures

type BlockSubFigure struct {
	BlockStruct
	Source  string
	Padding string
	Width   string
}

func (b *BlockSubFigure) String() string {
	return fmt.Sprintf(
		"BlockSubFigure (source=%v, padding=%v), %v",
		b.Source,
		b.Padding,
		b.BlockStruct.String(),
	)
}

func (b *BlockSubFigure) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|subfigure>")
}

func (b *BlockSubFigure) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockSubFigure) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 11,
	}
	options := GatherBlockOptions(line, []string{"src", "padding", "width"})
	if value, ok := options["src"]; ok {
		b.Source = value
	}
	if value, ok := options["padding"]; ok {
		b.Padding = value
	}
	if value, ok := options["width"]; ok {
		b.Width = value
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockSubFigure) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<subfigure|")
	bed.EndViaNLI = false
}

func (b *BlockSubFigure) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 11,
	}
}

func (b *BlockSubFigure) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockSubFigure) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface{
		&BlockComment{},
		&BlockHTML{},
		&BlockLaTeX{},
		&BlockCodeListing{},
		&BlockInlineCodeListing{},
		&BlockMath{},
		&BlockInlineMath{},
		&BlockFootnote{},
		&BlockRef{},
	}
}

func (b *BlockSubFigure) AcceptBlockInside(_ BlockInterface) bool {
	return true
}

func (b *BlockSubFigure) IsPartOfParagraph() bool {
	return false
}

func (b *BlockSubFigure) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockSubFigure) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockSubFigure) GetRawContent() *string {
	return nil
}

// =====================================
// Tabs

type BlockTabs struct {
	BlockStruct
	Tabs          []*BlockTabsTab
	SelectedIndex int
}

func (b *BlockTabs) String() string {
	return fmt.Sprintf(
		"BlockTabs (tabs=%v), %v",
		b.Tabs,
		b.BlockStruct.String(),
	)
}

func (b *BlockTabs) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|tabs>")
}

func (b *BlockTabs) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockTabs) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 6,
	}
	options := GatherBlockOptions(line, []string{"selected"})
	if value, ok := options["selected"]; ok {
		valueInt, err := strconv.Atoi(value)
		if err != nil || valueInt < 0 {
			log.Warnf("Could not use tabs option [selected=%v]. Please use natural numbers", value)
		}
		b.SelectedIndex = valueInt
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockTabs) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<tabs|")
	bed.EndViaNLI = false
}

func (b *BlockTabs) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 6,
	}
}

func (b *BlockTabs) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockTabs) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface{
		&BlockComment{},
		&BlockHTML{},
		&BlockLaTeX{},
		&BlockTabsTab{},
	}
}

func (b *BlockTabs) AcceptBlockInside(_ BlockInterface) bool {
	return true
}

func (b *BlockTabs) IsPartOfParagraph() bool {
	return false
}

func (b *BlockTabs) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockTabs) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockTabs) GetRawContent() *string {
	return nil
}

// =====================================
// Tabs tabs

type BlockTabsTab struct {
	BlockStruct
	Name       string
	IsSelected bool
}

func (b *BlockTabsTab) String() string {
	return fmt.Sprintf(
		"BlockTabsTab (name=%v), %v",
		b.Name,
		b.BlockStruct.String(),
	)
}

func (b *BlockTabsTab) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|tab>")
}

func (b *BlockTabsTab) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockTabsTab) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 5,
	}
	options := GatherBlockOptions(line, []string{"name"})
	if value, ok := options["name"]; ok {
		b.Name = value
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockTabsTab) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<tab|")
	bed.EndViaNLI = false
}

func (b *BlockTabsTab) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 5,
	}
}

func (b *BlockTabsTab) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockTabsTab) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface{
		&BlockComment{},
		&BlockHTML{},
		&BlockLaTeX{},
		&BlockTabs{},
		&BlockTextBox{},
		&BlockCodeListing{},
		&BlockInlineCodeListing{},
		&BlockMath{},
		&BlockInlineMath{},
		&BlockFootnote{},
		&BlockRef{},
	}
}

func (b *BlockTabsTab) AcceptBlockInside(_ BlockInterface) bool {
	return true
}

func (b *BlockTabsTab) IsPartOfParagraph() bool {
	return false
}

func (b *BlockTabsTab) DigDeeperForParagraphs() bool {
	return true
}

func (b *BlockTabsTab) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockTabsTab) GetRawContent() *string {
	return nil
}

// =====================================
// Footnotes

type BlockFootnote struct {
	BlockStruct
	FootnoteIndex int
}

func (b *BlockFootnote) String() string {
	return fmt.Sprintf(
		"BlockFootnote (index=%v), %v",
		b.FootnoteIndex,
		b.BlockStruct.String(),
	)
}

func (b *BlockFootnote) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|footnote>")
}

func (b *BlockFootnote) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockFootnote) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 10,
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockFootnote) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<footnote|")
	bed.EndViaNLI = false
}

func (b *BlockFootnote) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 10,
	}
}

func (b *BlockFootnote) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockFootnote) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface{
		&BlockComment{},
		&BlockHTML{},
		&BlockLaTeX{},
		&BlockCodeListing{},
		&BlockInlineCodeListing{},
		&BlockMath{},
		&BlockInlineMath{},
		&BlockUl{},
		&BlockOl{},
		&BlockRef{},
	}
}

func (b *BlockFootnote) AcceptBlockInside(_ BlockInterface) bool {
	return true
}

func (b *BlockFootnote) IsPartOfParagraph() bool {
	return true
}

func (b *BlockFootnote) DigDeeperForParagraphs() bool {
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
	File           string
	RawContent     string
	ReferenceIndex int
}

func (b *BlockRef) String() string {
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

func (b *BlockRef) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockRef) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 5,
	}
	options := GatherBlockOptions(line, []string{"file"})
	if value, ok := options["file"]; ok {
		b.File = value
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockRef) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<ref|")
	bed.EndViaNLI = false
}

func (b *BlockRef) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 5,
	}
}

func (b *BlockRef) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockRef) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b *BlockRef) AcceptBlockInside(_ BlockInterface) bool {
	return false // irrelevant
}

func (b *BlockRef) IsPartOfParagraph() bool {
	return true
}

func (b *BlockRef) DigDeeperForParagraphs() bool {
	return false
}

func (b *BlockRef) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockRef) GetRawContent() *string {
	return &b.RawContent
}

// =====================================
// Bibliography info

type BlockBibliography struct {
	BlockStruct
	HTMLContent  *string
	LaTeXContent *string
}

func (b *BlockBibliography) String() string {
	hc := "<nil>"
	if b.HTMLContent != nil {
		hc = *b.HTMLContent
	}
	lc := "<nil>"
	if b.LaTeXContent != nil {
		lc = *b.LaTeXContent
	}
	return fmt.Sprintf(
		"BlockBibliography (html-content=%v, latex-content=%v), %v",
		hc,
		lc,
		b.BlockStruct.String(),
	)
}

func (b *BlockBibliography) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|bibliography>")
}

func (b *BlockBibliography) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockBibliography) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 14,
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockBibliography) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<bibliography|")
	bed.EndViaNLI = false
}

func (b *BlockBibliography) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 14,
	}
}

func (b *BlockBibliography) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockBibliography) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b *BlockBibliography) AcceptBlockInside(_ BlockInterface) bool {
	return false // irrelevant
}

func (b *BlockBibliography) IsPartOfParagraph() bool {
	return false
}

func (b *BlockBibliography) DigDeeperForParagraphs() bool {
	return false
}

func (b *BlockBibliography) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockBibliography) GetRawContent() *string {
	return nil
}

// =====================================
// Meta info

type BlockMeta struct {
	BlockStruct
}

func (b *BlockMeta) String() string {
	return fmt.Sprintf(
		"BlockMeta, %v",
		b.BlockStruct.String(),
	)
}

func (b *BlockMeta) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|meta>")
}

func (b *BlockMeta) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockMeta) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 6,
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockMeta) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<meta|")
	bed.EndViaNLI = false
}

func (b *BlockMeta) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 6,
	}
}

func (b *BlockMeta) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockMeta) GetBlocksAllowedInside() []BlockInterface {
	return []BlockInterface{
		// &BlockComment{}, - not needed, will not be ported to HTML
		&BlockMetaAuthor{},
		&BlockMetaCopyright{},
		&BlockMetaBibInfo{},
	}
}

func (b *BlockMeta) AcceptBlockInside(_ BlockInterface) bool {
	return true
}

func (b *BlockMeta) IsPartOfParagraph() bool {
	return false
}

func (b *BlockMeta) DigDeeperForParagraphs() bool {
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

func (b *BlockMetaAuthor) String() string {
	return fmt.Sprintf(
		"BlockMetaAuthor, %v :: \"%v\"",
		b.BlockStruct.String(),
		b.RawContent,
	)
}

func (b *BlockMetaAuthor) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|author>")
}

func (b *BlockMetaAuthor) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockMetaAuthor) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 8,
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockMetaAuthor) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<author|")
	bed.EndViaNLI = false
}

func (b *BlockMetaAuthor) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 8,
	}
}

func (b *BlockMetaAuthor) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockMetaAuthor) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b *BlockMetaAuthor) AcceptBlockInside(_ BlockInterface) bool {
	return false // irrelevant
}

func (b *BlockMetaAuthor) IsPartOfParagraph() bool {
	return false
}

func (b *BlockMetaAuthor) DigDeeperForParagraphs() bool {
	return false
}

func (b *BlockMetaAuthor) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockMetaAuthor) GetRawContent() *string {
	return &b.RawContent
}

// =====================================
// Meta info - copyright

type BlockMetaCopyright struct {
	BlockStruct
	RawContent string
}

func (b *BlockMetaCopyright) String() string {
	return fmt.Sprintf(
		"BlockMetaCopyright, %v :: \"%v\"",
		b.BlockStruct.String(),
		b.RawContent,
	)
}

func (b *BlockMetaCopyright) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|copyright>")
}

func (b *BlockMetaCopyright) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockMetaCopyright) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 11,
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockMetaCopyright) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<copyright|")
	bed.EndViaNLI = false
}

func (b *BlockMetaCopyright) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 11,
	}
}

func (b *BlockMetaCopyright) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockMetaCopyright) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b *BlockMetaCopyright) AcceptBlockInside(_ BlockInterface) bool {
	return false // irrelevant
}

func (b *BlockMetaCopyright) IsPartOfParagraph() bool {
	return false
}

func (b *BlockMetaCopyright) DigDeeperForParagraphs() bool {
	return false
}

func (b *BlockMetaCopyright) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockMetaCopyright) GetRawContent() *string {
	return &b.RawContent
}

// =====================================
// Meta info - bibliography

type BlockMetaBibInfo struct {
	BlockStruct
	RawContent string
	JSONInline bool
	RefFile    string
}

func (b *BlockMetaBibInfo) String() string {
	return fmt.Sprintf(
		"BlockMetaBibInfo (inline=%v), %v :: \"%v\"",
		b.JSONInline,
		b.BlockStruct.String(),
		b.RawContent,
	)
}

func (b *BlockMetaBibInfo) CheckBlockStarts(line LineStruct) bool {
	return CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "|bibinfo>")
}

func (b *BlockMetaBibInfo) SeekBufferAfterBlockStarts() int {
	return 1
}

func (b *BlockMetaBibInfo) ExecuteAfterBlockStarts(line *LineStruct) {
	b.Start = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 9,
	}
	options := GatherBlockOptions(line, []string{"inline", "ref-file"})
	if value, ok := options["inline"]; ok {
		b.JSONInline = Contains([]string{"allow", "allowed", "1", "true", "ok", "yes"}, value)
	}
	if value, ok := options["ref-file"]; ok {
		b.RefFile = value
	}
	b.ContentStart = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
}

func (b *BlockMetaBibInfo) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	bed.EndNormally = CheckRunesEndWithUnescapedASCII(line.RuneContent[:line.RuneJ+1], "<bibinfo|")
	bed.EndViaNLI = false
}

func (b *BlockMetaBibInfo) ExecuteAfterBlockEnds(line *LineStruct) {
	b.End = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ,
	}
	b.ContentEnd = Pair[int, int]{
		i: line.LineIndex,
		j: line.RuneJ - 9,
	}
}

func (b *BlockMetaBibInfo) SeekBufferAfterBlockEnds() int {
	return 1
}

func (b *BlockMetaBibInfo) GetBlocksAllowedInside() []BlockInterface {
	return nil
}

func (b *BlockMetaBibInfo) AcceptBlockInside(_ BlockInterface) bool {
	return false // irrelevant
}

func (b *BlockMetaBibInfo) IsPartOfParagraph() bool {
	return false
}

func (b *BlockMetaBibInfo) DigDeeperForParagraphs() bool {
	return false
}

func (b *BlockMetaBibInfo) GetBlockStruct() *BlockStruct {
	return &b.BlockStruct
}

func (b *BlockMetaBibInfo) GetRawContent() *string {
	return &b.RawContent
}

// =====================================
// Generic inline block

type BlockInline struct {
	Content InlineInterface
}

func (b *BlockInline) String() string {
	return fmt.Sprintf(
		"(BlockInline->)%v",
		b.Content.String(),
	)
}

func (b *BlockInline) CheckBlockStarts(_ LineStruct) bool {
	return false // irrelevant
}

func (b *BlockInline) SeekBufferAfterBlockStarts() int {
	return 0 // irrelevant
}

func (b *BlockInline) ExecuteAfterBlockStarts(_ *LineStruct) {
	// irrelevant
}

func (b *BlockInline) CheckBlockEnds(line *LineStruct, bed *BlockEndDetails, NewLines int, Indentation uint16, parsingStack ParsingStack) {
	// irrelevant
	bed.EndNormally = false
	bed.EndViaNLI = false
}

func (b *BlockInline) ExecuteAfterBlockEnds(_ *LineStruct) {
	// irrelevant
}

func (b *BlockInline) SeekBufferAfterBlockEnds() int {
	return 1 // irrelevant
}

func (b *BlockInline) GetBlocksAllowedInside() []BlockInterface {
	return nil // irrelevant
}

func (b *BlockInline) AcceptBlockInside(_ BlockInterface) bool {
	return false // irrelevant
}

func (b *BlockInline) IsPartOfParagraph() bool {
	return true // irrelevant
}

func (b *BlockInline) DigDeeperForParagraphs() bool {
	return true // irrelevant
}

func (b *BlockInline) GetBlockStruct() *BlockStruct {
	return new(BlockStruct) // irrelevant
}

func (b *BlockInline) GetRawContent() *string {
	return b.Content.GetRawContent() // This might always be nil
}
