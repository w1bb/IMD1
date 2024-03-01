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
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// =====================================
// Parsing stack structure

type ParsingStack struct {
	InsideTextBox   uint32
	InsideParagraph uint32
}

// =====================================
// Meta data structure

type MDMetaStructure struct {
	Author    string
	Copyright string
}

func (m MDMetaStructure) String() string {
	return fmt.Sprintf(
		"Metadata:\n| Author: %v\n| Copyright: %v",
		m.Author,
		m.Copyright,
	)
}

func (m MDMetaStructure) Serialize() []byte {
	r := make([]byte, 0)
	r = append(r, StringSerialize(m.Author)...)
	r = append(r, StringSerialize(m.Copyright)...)
	return r
}

// =====================================
// Tree parse

// Parse Please note that this will modify RuneJs
func (file *FileStruct) Parse() (Tree[BlockInterface], MDMetaStructure) {
	parsingStack := ParsingStack{}
	currentBlock := &(Tree[BlockInterface]{})
	currentBlock.Value = &BlockDocument{}
	allowedInsideBlock := currentBlock.Value.GetBlocksAllowedInside()

	justSkipped := 16 // considered infinite
	for i := range file.Lines {
		if file.Lines[i].Empty() {
			justSkipped++
			continue
		}
		for file.Lines[i].RuneJ < len(file.Lines[i].RuneContent) {
			if file.Lines[i].RuneJ > 0 {
				justSkipped = 0
			}
			// Check if we can end this block
			for file.Lines[i].RuneJ < len(file.Lines[i].RuneContent) {
				shouldEndBlock, shouldDiscardBlock, shouldDiscardSeek := currentBlock.Value.CheckBlockEndsNormally(&file.Lines[i], parsingStack)
				if !shouldEndBlock {
					shouldEndBlock = currentBlock.Value.CheckBlockEndsViaNewLinesAndIndentation(justSkipped, file.Lines[i].Indentation)
				}
				if !shouldEndBlock {
					break
				}
				log.Debugf(
					"<<< End of %v (line-index=%v, line-j=%v)",
					currentBlock.Value,
					i, file.Lines[i].RuneJ,
				)

				switch reflect.TypeOf(currentBlock.Value) {
				case reflect.TypeOf(BlockParagraph{}):
					parsingStack.InsideParagraph--
				case reflect.TypeOf(BlockTextBox{}):
					parsingStack.InsideTextBox--
				}

				// Discard if needed
				if shouldDiscardBlock != nil {
					file.Lines[i].RuneJ += shouldDiscardSeek
					for reflect.TypeOf(currentBlock.Value) != reflect.TypeOf(shouldDiscardBlock) {
						fmt.Printf("SOFT Discarding %v (%T != %T)\n", currentBlock.Value, currentBlock.Value, shouldDiscardBlock)
						// TODO - remove
						currentBlock.Value.ExecuteAfterBlockEnds(&file.Lines[i])
						currentBlock = currentBlock.Parent
					}
					file.Lines[i].RuneJ -= shouldDiscardSeek
				}

				file.Lines[i].RuneJ += currentBlock.Value.SeekBufferAfterBlockEnds()
				currentBlock.Value.ExecuteAfterBlockEnds(&file.Lines[i])
				currentBlock = currentBlock.Parent

				allowedInsideBlock = currentBlock.Value.GetBlocksAllowedInside()
			}
			if file.Lines[i].RuneJ >= len(file.Lines[i].RuneContent) {
				break
			}
			// Check if we can open blocks
			for found := true; found && file.Lines[i].RuneJ < len(file.Lines[i].RuneContent); {
				found = false
				for _, allowed := range allowedInsideBlock {
					if allowed.CheckBlockStarts(file.Lines[i]) && currentBlock.Value.AcceptBlockInside(allowed) {
						file.Lines[i].RuneJ += allowed.SeekBufferAfterBlockStarts()
						nextBlock := &(Tree[BlockInterface]{})
						nextBlock.Value = allowed
						nextBlock.Parent = currentBlock
						currentBlock.Children = append(currentBlock.Children, nextBlock)
						currentBlock = nextBlock
						allowedInsideBlock = allowed.GetBlocksAllowedInside()
						found = true
						break
					}
				}
				if !found {
					file.Lines[i].RuneJ++
				} else {
					switch reflect.TypeOf(currentBlock.Value) {
					case reflect.TypeOf(&BlockParagraph{}):
						parsingStack.InsideParagraph++
					case reflect.TypeOf(&BlockTextBox{}):
						parsingStack.InsideTextBox++
					}
					currentBlock.Value.ExecuteAfterBlockStarts(&file.Lines[i])
					log.Debugf(
						">>> Beginning of %v (line-index=%v, line-j=%v)",
						currentBlock.Value,
						i, file.Lines[i].RuneJ,
					)
				}
			}
		}
		justSkipped = 0
	}

	// Get back to root
	for currentBlock.Parent != nil {
		switch reflect.TypeOf(currentBlock.Value) {
		case reflect.TypeOf(BlockParagraph{}):
			parsingStack.InsideParagraph--
		case reflect.TypeOf(BlockTextBox{}):
			parsingStack.InsideTextBox--
		}
		currentBlock.Value.ExecuteAfterBlockEnds(&file.Lines[len(file.Lines)-1])
		currentBlock = currentBlock.Parent
	}
	currentBlock.Value.GetBlockStruct().ContentEnd = Pair[int, int]{
		i: len(file.Lines) - 1,
		j: len(file.Lines[len(file.Lines)-1].RuneContent),
	}

	// Find first non-paragraph
	FirstNonParagraph := func(tree *Tree[BlockInterface], minI int) int {
		for i := minI; i < len(tree.Children); i++ {
			if !tree.Children[i].Value.IsPartOfParagraph() {
				return i
			}
		}
		return -1
	}

	// Create the number of detected paragraphs
	DetectParagraphs := func(before, after Pair[int, int], file *FileStruct) []BlockParagraph {
		var paragraphs []BlockParagraph
		var currentParagraph string

		for i := before.i; i <= after.i; i++ {
			var j = 0
			var maxJ = len(file.Lines[i].RuneContent)
			if i == before.i {
				j = before.j
			}
			if i == after.i {
				maxJ = after.j
			}
			currentParagraph += string(file.Lines[i].RuneContent[j:maxJ]) + " "
			if file.Lines[i].Empty() || i == after.i {
				currentParagraph = RemoveExcessSpaces(currentParagraph)
				if currentParagraph != "" {
					end := Pair[int, int]{
						i: i,
						j: 0, // Default
					}
					if i == after.i {
						end.j = after.j
					}
					paragraphs = append(paragraphs, BlockParagraph{
						BlockStruct: BlockStruct{
							Start:        before,
							End:          end,
							ContentStart: before,
							ContentEnd:   end,
						},
					})
					before = end
				}
				currentParagraph = ""
			}
		}
		return paragraphs
	}

	// Insert paragraphs
	var InsertParagraphs func(*Tree[BlockInterface])
	InsertParagraphs = func(tree *Tree[BlockInterface]) {

		// Go deeper (TODO - multithreading)
		for i := range tree.Children {
			if !tree.Children[i].Value.DigDeeperForParagraphs() {
				continue
			}
			InsertParagraphs(tree.Children[i])
		}

		var before, after Pair[int, int]
		before = tree.Value.GetBlockStruct().ContentStart

		startingI := 0
		for i := FirstNonParagraph(tree, 0); i != -1; i = FirstNonParagraph(tree, startingI) {
			after = tree.Children[i].Value.GetBlockStruct().Start
			nextBefore := tree.Children[i].Value.GetBlockStruct().End
			if p := DetectParagraphs(before, after, file); p != nil {
				log.Debug("Detected paragraph between ", before, " and ", after, ": ", p)
				generatedP := make([]*Tree[BlockInterface], len(p))
				for j := range p {
					generatedP[j] = new(Tree[BlockInterface])
					generatedP[j].Parent = tree
					generatedP[j].Value = &p[j]
				}
				for ci, j := startingI, 0; ci < i && j < len(p); ci++ {
					for cbs := tree.Children[ci].Value.GetBlockStruct(); cbs.Start.i > p[j].ContentEnd.i; j++ {
					}
					tree.Children[ci].Parent = generatedP[j]
					generatedP[j].Children = append(generatedP[j].Children, tree.Children[ci])
				}
				// Remove old blocks
				aux := make([]*Tree[BlockInterface], len(p)+len(tree.Children)-(i-startingI))
				copy(aux[:startingI], tree.Children[:startingI])
				copy(aux[startingI:startingI+len(p)], generatedP)
				copy(aux[startingI+len(p):], tree.Children[i:])
				tree.Children = aux
				startingI += len(p)
			}
			startingI++
			before = nextBefore
		}

		after = tree.Value.GetBlockStruct().ContentEnd
		if p := DetectParagraphs(before, after, file); p != nil {
			log.Debug("Detected paragraph between ", before, " and ", after, ": ", p)
			generatedP := make([]*Tree[BlockInterface], len(p))
			for j := range p {
				generatedP[j] = new(Tree[BlockInterface])
				generatedP[j].Parent = tree
				generatedP[j].Value = &p[j]
			}
			for ci, j := startingI, 0; ci < len(tree.Children) && j < len(p); ci++ {
				for cbs := tree.Children[ci].Value.GetBlockStruct(); cbs.Start.i > p[j].ContentEnd.i; j++ {
				}
				tree.Children[ci].Parent = generatedP[j]
				generatedP[j].Children = append(generatedP[j].Children, tree.Children[ci])
			}
			// Remove old blocks
			aux := make([]*Tree[BlockInterface], len(p)+len(tree.Children)-(len(tree.Children)-startingI))
			copy(aux[:startingI], tree.Children[:startingI])
			copy(aux[startingI:startingI+len(p)], generatedP)
			copy(aux[startingI+len(p):], tree.Children[len(tree.Children):])
			tree.Children = aux
			startingI += len(p)
		}
	}
	InsertParagraphs(currentBlock)

	var mdMeta MDMetaStructure

	// Parse paragraphs, raw content, sub figures and meta info
	CompleteBibInfo := make(map[string]BibliographyEntry)
	CompleteBibInfoNextRefIndex := 1

	var DFSNodeParse func(tree *Tree[BlockInterface])
	RefBlocks := make([]*BlockRef, 0)
	BibliographyBlocks := make([]*BlockBibliography, 0)
	TabsBlocks := make([]*BlockTabs, 0)
	DFSNodeParse = func(tree *Tree[BlockInterface]) {
		if rc := tree.Value.GetRawContent(); rc != nil && *rc == "" {
			bs := tree.Value.GetBlockStruct()
			*rc = file.GetStringBetween(bs.ContentStart, bs.ContentEnd)
		}

		switch reflect.TypeOf(tree.Value) {
		case reflect.TypeOf(&BlockParagraph{}):
			ParseSingleParagraph(tree, file)
		case reflect.TypeOf(&BlockSubFigure{}):
			// Sanity check
			if reflect.TypeOf(tree.Parent.Value) != reflect.TypeOf(&BlockFigure{}) {
				panic(nil)
			}
			if tree.Value.(*BlockSubFigure).Padding == "" {
				tree.Value.(*BlockSubFigure).Padding = tree.Parent.Value.(*BlockFigure).Padding
			}
		case reflect.TypeOf(&BlockMetaAuthor{}):
			// Sanity check
			if reflect.TypeOf(tree.Parent.Value) != reflect.TypeOf(&BlockMeta{}) {
				panic(nil)
			}
			mdMeta.Author = *tree.Value.GetRawContent()
		case reflect.TypeOf(&BlockMetaCopyright{}):
			// Sanity check
			if reflect.TypeOf(tree.Parent.Value) != reflect.TypeOf(&BlockMeta{}) {
				panic(nil)
			}
			mdMeta.Copyright = *tree.Value.GetRawContent()
		case reflect.TypeOf(&BlockMetaBibInfo{}):
			// Sanity check
			if reflect.TypeOf(tree.Parent.Value) != reflect.TypeOf(&BlockMeta{}) {
				panic(nil)
			}
			newBibInfo := ParseBibInfo(tree.Value.(*BlockMetaBibInfo))
			orderedNewBibInfo := make([]string, 0)
			for k := range newBibInfo {
				orderedNewBibInfo = append(orderedNewBibInfo, k)
			}
			sort.Strings(orderedNewBibInfo)
			for _, newBibInfoKey := range orderedNewBibInfo {
				newBibInfoValue := newBibInfo[newBibInfoKey]
				if oldValue, ok := CompleteBibInfo[newBibInfoKey]; ok {
					log.Warnf("A bibliography entry for tag \"%v\" already exists (%v). Keeping the old value...", newBibInfoKey, oldValue)
					continue
				}
				CompleteBibInfo[newBibInfoKey] = BibliographyEntry{
					ParentBlock:    newBibInfoValue.ParentBlock,
					Type:           newBibInfoValue.Type,
					Fields:         newBibInfoValue.Fields,
					ReferenceIndex: CompleteBibInfoNextRefIndex,
					File:           tree.Value.(*BlockMetaBibInfo).RefFile,
				}
				CompleteBibInfoNextRefIndex++
			}
		case reflect.TypeOf(&BlockRef{}):
			RefBlocks = append(RefBlocks, tree.Value.(*BlockRef))
		case reflect.TypeOf(&BlockBibliography{}):
			BibliographyBlocks = append(BibliographyBlocks, tree.Value.(*BlockBibliography))
		case reflect.TypeOf(&BlockTabsTab{}):
			// Sanity check
			if reflect.TypeOf(tree.Parent.Value) != reflect.TypeOf(&BlockTabs{}) {
				panic(nil)
			}
			tabs := tree.Parent.Value.(*BlockTabs)
			tabsTab := tree.Value.(*BlockTabsTab)
			tabs.Tabs = append(tabs.Tabs, tabsTab)
		case reflect.TypeOf(&BlockTabs{}):
			tabs := tree.Value.(*BlockTabs)
			TabsBlocks = append(TabsBlocks, tabs)
		}

		for i := 0; i < len(tree.Children); i++ {
			DFSNodeParse(tree.Children[i])
		}
	}
	DFSNodeParse(currentBlock)

	// Convert references
	for i := 0; i < len(RefBlocks); i++ {
		var entry BibliographyEntry
		var ok bool
		if entry, ok = CompleteBibInfo[RefBlocks[i].RawContent]; !ok {
			log.Warnf("Reference to \"%v\" cannot be resolved (will be rendered as \"[?]\")", RefBlocks[i].RawContent)
			continue
		}
		RefBlocks[i].ReferenceIndex = entry.ReferenceIndex
		if RefBlocks[i].File == "" {
			RefBlocks[i].File = entry.File
		}
	}

	// Generate bibliographies
	bibliographyHtmlText := GenerateBibliography(CompleteBibInfo)
	for i := 0; i < len(BibliographyBlocks); i++ {
		BibliographyBlocks[i].HTMLContent = &bibliographyHtmlText
	}

	// Set selected tabs
	for i := 0; i < len(TabsBlocks); i++ {
		if len(TabsBlocks[i].Tabs) <= TabsBlocks[i].SelectedIndex {
			log.Warnf(
				"The selected index (%v) for a |tabs> element is too large. Will reset to 0...",
				TabsBlocks[i].SelectedIndex,
			)
			TabsBlocks[i].SelectedIndex = 0
		}
		TabsBlocks[i].Tabs[TabsBlocks[i].SelectedIndex].IsSelected = true
	}

	return *currentBlock, mdMeta
}

// =====================================
// Bibliography parse

type BibliographyEntryType uint8

const (
	BibliographyEntryTypeArticle BibliographyEntryType = iota
	BibliographyEntryTypeBook
	BibliographyEntryTypeOther
)

func (t BibliographyEntryType) String() string {
	switch t {
	case BibliographyEntryTypeArticle:
		return "Article"
	case BibliographyEntryTypeBook:
		return "Book"
	case BibliographyEntryTypeOther:
		return "Other"
	default:
		panic(nil) // This should never be reached
	}
}

type BibliographyEntryFields struct {
	Title     *string
	Author    *string
	Journal   *string
	Volume    *string
	Number    *string
	Pages     *string
	Year      *string
	Publisher *string
	URL       *string
}

func (f BibliographyEntryFields) String() string {
	s := ""
	if f.Title != nil {
		s += "title=" + *f.Title + ", "
	}
	if f.Author != nil {
		s += "author=" + *f.Author + ", "
	}
	if f.Journal != nil {
		s += "journal=" + *f.Journal + ", "
	}
	if f.Volume != nil {
		s += "volume=" + *f.Volume + ", "
	}
	if f.Number != nil {
		s += "number=" + *f.Number + ", "
	}
	if f.Pages != nil {
		s += "pages=" + *f.Pages + ", "
	}
	if f.Year != nil {
		s += "Year=" + *f.Year + ", "
	}
	if f.Publisher != nil {
		s += "publisher=" + *f.Publisher + ", "
	}
	if f.URL != nil {
		s += "url=" + *f.URL + ", "
	}
	if s == "" {
		return "Fields (-)"
	}
	return "Fields (" + s[:len(s)-2] + ")"
}

type BibliographyEntry struct {
	ParentBlock    *BlockMetaBibInfo
	Type           BibliographyEntryType
	Fields         BibliographyEntryFields
	ReferenceIndex int
	File           string
}

func (e BibliographyEntry) String() string {
	return fmt.Sprintf(
		"BibliographyEntry (parent=%v, type=%v, fields=%v, index=%v)",
		e.ParentBlock,
		e.Type,
		e.Fields,
		e.ReferenceIndex,
	)
}

func ParseBibInfo(b *BlockMetaBibInfo) map[string]BibliographyEntry {
	entries := make(map[string]BibliographyEntry)
	bibInfoTypeStr := ""

	// Types - don't call reflect.TypeOf multiple times!
	reflectTypeString := reflect.TypeOf("a")
	reflectTypeMapStringInterface := reflect.TypeOf(make(map[string]interface{}))
	reflectTypeSliceInterface := reflect.TypeOf(make([]interface{}, 0))

	var jsonSb []byte
	if b.JSONInline {
		bibInfoTypeStr = "inline bibinfo \"" + b.RawContent + "\""
		jsonSb = []byte(b.RawContent)
	} else {
		bibInfoTypeStr = "bibinfo file (" + b.RawContent + ")"
		var err error
		jsonSb, err = os.ReadFile(b.RawContent)
		if err != nil {
			log.Errorf("Could not read %v. Skipping nil...", bibInfoTypeStr)
			return nil
		}
	}

	// Unmarshall the JSON
	var fullJsonInterface interface{}
	err := json.Unmarshal(jsonSb, &fullJsonInterface)
	if err != nil {
		log.Errorf("Could not unmarshal %v. Skipping nil...", bibInfoTypeStr)
		return nil
	}
	fullJson := fullJsonInterface.(map[string]interface{})

	// Search for "bibliography"
	if _, ok := fullJson["bibliography"]; !ok {
		log.Warnf("Could not find \"bibliography\" entry in %v. Skipping empty...", bibInfoTypeStr)
		return entries
	}
	if reflect.TypeOf(fullJson["bibliography"]) != reflectTypeSliceInterface {
		log.Warnf("The \"bibliography\" in %v is NOT an array. Skipping empty...", bibInfoTypeStr)
		return entries
	}
	bibliographyJson := fullJson["bibliography"].([]interface{})

	// Search for entries
	if len(bibliographyJson) == 0 {
		log.Warnf("The \"bibliography\" in %v is empty.", bibInfoTypeStr)
		return entries
	}
	for entryI := 0; entryI < len(bibliographyJson); entryI++ {
		if reflect.TypeOf(bibliographyJson[entryI]) != reflectTypeMapStringInterface {
			log.Warnf("\"bibliography\" entry %v in %v will be ignored (not a valid map[string]...)", entryI, bibInfoTypeStr)
			continue
		}
		bibliographyEntry := bibliographyJson[entryI].(map[string]interface{})

		var entry BibliographyEntry
		entry.ParentBlock = b

		entryTag := ""
		if bibliographyEntryTag, ok := bibliographyEntry["tag"]; !ok {
			log.Warnf("\"bibliography\" entry %v in %v will be ignored (missing tag)", entryI, bibInfoTypeStr)
			continue
		} else if reflect.TypeOf(bibliographyEntryTag) != reflectTypeString {
			log.Warnf("\"bibliography\" entry %v in %v will be ignored (tag exists, but is not a string)", entryI, bibInfoTypeStr)
			continue
		} else {
			entryTag = bibliographyEntryTag.(string)
		}

		entryType := ""
		entry.Type = BibliographyEntryTypeOther
		if bibliographyEntryType, ok := bibliographyEntry["type"]; !ok {
			log.Warnf("Missing \"bibliography\" entry type in entry %v, %v. Considering type \"other\"...", entryI, bibInfoTypeStr)
		} else if reflect.TypeOf(bibliographyEntryType) != reflectTypeString {
			log.Warnf("\"bibliography\" entry type is not a string in entry %v, %v. Considering type \"other\"...", entryI, bibInfoTypeStr)
		} else {
			entryType = bibliographyEntryType.(string)
			switch entryType {
			case "article", "Article":
				entry.Type = BibliographyEntryTypeArticle
			case "book", "Book":
				entry.Type = BibliographyEntryTypeBook
			case "other", "Other", "unknown", "Unknown":
			default:
				log.Warnf("Unrecognized \"bibliography\" entry type \"%v\" in entry %v, %v. Considering type \"other\"...", entryType, entryI, bibInfoTypeStr)
			}
		}

		if bibliographyEntryData, ok := bibliographyEntry["data"]; ok && reflect.TypeOf(bibliographyEntryData) != reflectTypeMapStringInterface {
			log.Warnf("\"bibliography\" entry %v's \"data\" in %v will be ignored (not valid map[string]...)", entryI, bibInfoTypeStr)
		} else if ok {
			entryData := bibliographyEntryData.(map[string]interface{})
			for entryDataKey, entryDataValue := range entryData {
				if reflect.TypeOf(entryDataValue) != reflectTypeString {
					log.Warnf("\"bibliography\" entry data value must be string in entry %v, %v. Ignoring only data key %v...", entryI, bibInfoTypeStr, entryDataKey)
					continue
				}
				dataValue := entryDataValue.(string)
				switch entryDataKey {
				case "title":
					entry.Fields.Title = new(string)
					*entry.Fields.Title = dataValue
				case "author":
					entry.Fields.Author = new(string)
					*entry.Fields.Author = dataValue
				case "journal":
					entry.Fields.Journal = new(string)
					*entry.Fields.Journal = dataValue
				case "volume":
					entry.Fields.Volume = new(string)
					*entry.Fields.Volume = dataValue
				case "number":
					entry.Fields.Number = new(string)
					*entry.Fields.Number = dataValue
				case "pages":
					entry.Fields.Pages = new(string)
					*entry.Fields.Pages = dataValue
				case "year":
					entry.Fields.Year = new(string)
					*entry.Fields.Year = dataValue
				case "publisher":
					entry.Fields.Publisher = new(string)
					*entry.Fields.Publisher = dataValue
				case "url":
					entry.Fields.URL = new(string)
					*entry.Fields.URL = dataValue
				default:
					log.Warnf("Unrecognized \"bibliography\" entry data key \"%v\" in entry %v, %v. Ignoring only that data key...", entryDataKey, entryI, bibInfoTypeStr)
				}
			}
		}
		entries[entryTag] = entry
	}
	return entries
}

func GenerateBibliography(mp map[string]BibliographyEntry) string {
	var sb strings.Builder
	sb.WriteString("<div class=\"bibliography\">\n")
	orderedBibInfo := make([]string, 0)
	for k := range mp {
		orderedBibInfo = append(orderedBibInfo, k)
	}
	sort.Strings(orderedBibInfo)
	for _, key := range orderedBibInfo {
		value := mp[key]
		sb.WriteString("<div class=\"bib-entry\" id=\"ref-")
		sb.WriteString(strconv.Itoa(value.ReferenceIndex))
		sb.WriteString("\">")

		sb.WriteString("<div class=\"bib-entry-index-wrapper\">")
		sb.WriteString("<div class=\"bib-entry-index\">[")
		sb.WriteString(strconv.Itoa(value.ReferenceIndex))
		sb.WriteString("]</div></div>")

		sb.WriteString("<div class=\"bib-entry-text-wrapper\">")
		sb.WriteString("<div class=\"bib-entry-text\">")

		{
			// Author
			if value.Fields.Author != nil {
				sb.WriteString("<span class=\"author\">")
				sb.WriteString(*value.Fields.Author)
				sb.WriteString("</span>")
				sb.WriteString(" - ")
			}
			// Title
			sb.WriteString("<span class=\"title\">")
			if value.Fields.Title != nil {
				sb.WriteString(*value.Fields.Title)
			} else {
				sb.WriteString(key)
			}
			sb.WriteString("</span>")
			// Year
			if value.Fields.Year != nil {
				sb.WriteString(", <span class=\"year\">")
				sb.WriteString(*value.Fields.Year)
				sb.WriteString("</span>")
			}
			// URL
			if value.Fields.URL != nil {
				sb.WriteString(", <a class=\"url\" href=\"")
				sb.WriteString(*value.Fields.URL)
				sb.WriteString("\">")
				sb.WriteString(*value.Fields.URL)
				sb.WriteString("</a>")
			}
			// TODO - remaining fields
		}

		sb.WriteString("</div></div></div>\n")
	}
	sb.WriteString("</div>\n")
	return sb.String()
}

// =====================================
// Paragraph parse

func ParagraphInsertRawHelperBuilder(tree *Tree[BlockInterface], sb *strings.Builder) {
	if sb.Len() > 0 {
		tree.Children = append(tree.Children,
			&Tree[BlockInterface]{
				Parent: tree,
				Value: &BlockInline{
					Content: &InlineRawString{
						Content: sb.String(),
					},
				},
			})
		sb.Reset()
	}
}

func ParseSingleParagraph(tree *Tree[BlockInterface], file *FileStruct) {
	// Insert links
	linkParsedTree := ParseSingleParagraphLinks(tree, file)
	// Insert emphasis
	emphasisParsedTree := ParseSingleBlockInlineEmphasis(linkParsedTree)
	// Cleanup
	cleanedTree := CleanupSingleBlockInline(emphasisParsedTree)

	// Convert:        heading-paragraph -> normal heading
	//          text box title-paragraph -> normal text box title
	tree.Children = []*Tree[BlockInterface]{cleanedTree}
	cleanedTree.Parent = tree
	if tree.Parent != nil && (reflect.TypeOf(tree.Parent.Value) == reflect.TypeOf(&BlockHeading{}) || reflect.TypeOf(tree.Parent.Value) == reflect.TypeOf(&BlockTextBoxTitle{})) {
		h := tree.Parent
		cleanedTree.Parent = h
		h.Children = []*Tree[BlockInterface]{cleanedTree}
	}
}

func CleanupSingleBlockInline(contentTree *Tree[BlockInterface]) *Tree[BlockInterface] {
	cleanedContentTree := &Tree[BlockInterface]{
		Parent: contentTree.Parent,
		Value: &BlockInline{
			Content: &InlineDocument{},
		},
	}
	for i := 0; i < len(contentTree.Children); i++ {
		var srb strings.Builder
		j := i
	digging:
		for ; j < len(contentTree.Children); j++ {
			if reflect.TypeOf(contentTree.Children[j].Value) != reflect.TypeOf(&BlockInline{}) {
				break
			}
			bjc := contentTree.Children[j].Value.(*BlockInline).Content
			switch reflect.TypeOf(bjc) {
			case reflect.TypeOf(&InlineStringDelimiter{}):
				switch bjc.(*InlineStringDelimiter).TypeOfDelimiter {
				case InlineDelimiterTypeUnderlineDelimiter:
					srb.WriteRune('_')
				case InlineDelimiterTypeAsteriskDelimiter:
					srb.WriteRune('*')
				case InlineDelimiterTypeTildeDelimiter:
					srb.WriteRune('~')
				case InlineDelimiterTypeOpenBracketDelimiter:
					srb.WriteRune('[')
				case InlineDelimiterTypeCloseBracketDelimiter:
					srb.WriteRune(']')
				case InlineDelimiterTypeOpenParenthesesDelimiter:
					srb.WriteRune('(')
				case InlineDelimiterTypeCloseParenthesesDelimiter:
					srb.WriteRune(')') // Might never be reached
				default:
					panic(nil) // This should never be reached
				}
			case reflect.TypeOf(&InlineRawString{}):
				srb.WriteString(bjc.(*InlineRawString).Content)
			default:
				break digging
			}
		}
		if j != i {
			cleanedContentTree.Children = append(
				cleanedContentTree.Children,
				&Tree[BlockInterface]{
					Parent: cleanedContentTree,
					Value: &BlockInline{
						Content: &InlineRawString{
							Content: srb.String(),
						},
					},
				},
			)
		}
		if j < len(contentTree.Children) {
			contentTree.Children[j].Parent = cleanedContentTree
			cleanedContentTree.Children = append(
				cleanedContentTree.Children,
				contentTree.Children[j],
			)
		}
		i = j
	}

	return cleanedContentTree
}

func ParseSingleBlockInlineEmphasis(tree *Tree[BlockInterface]) *Tree[BlockInterface] {
	contentTree := &Tree[BlockInterface]{
		Parent: tree,
		Value: &BlockInline{
			Content: &InlineDocument{},
		},
	}

	var delimiterStack Stack[InlineDelimiter]
	var currentStringBuilder strings.Builder
	for childI := 0; ; childI++ {
		if childI == len(tree.Children) {
			ParagraphInsertRawHelperBuilder(contentTree, &currentStringBuilder)
			break
		}
		treeChild := tree.Children[childI]

		// Only work with raw strings
		isRaw := false
		isHref := false
		if reflect.TypeOf(treeChild.Value) == reflect.TypeOf(&BlockInline{}) {
			bic := treeChild.Value.(*BlockInline).Content
			switch reflect.TypeOf(bic) {
			case reflect.TypeOf(&InlineRawString{}):
				isRaw = true
			case reflect.TypeOf(&InlineHref{}):
				isHref = true
			}
		}

		if !isRaw {
			ParagraphInsertRawHelperBuilder(contentTree, &currentStringBuilder)
			treeChild.Parent = contentTree
			if isHref {
				aux := ParseSingleBlockInlineEmphasis(treeChild)
				for i := 0; i < len(aux.Children); i++ {
					aux.Children[i].Parent = treeChild
				}
				treeChild.Children = aux.Children
			}
			contentTree.Children = append(contentTree.Children, treeChild)
		} else {
			parsingNow := []rune(treeChild.Value.(*BlockInline).Content.(*InlineRawString).Content)
			isEscaped := false
			for cI := 0; cI < len(parsingNow); cI++ {
				c := parsingNow[cI]
				if isEscaped {
					switch c {
					case '_', '*', '|', '~', '<', '>', '\\':
						currentStringBuilder.WriteRune(c)
					default:
						log.Warnf(
							"Unrecognized escape sequence \"\\%v\". Please use \"\\\\%v\" instead. The sequence will be treated as \"\\\\%v\"...",
							c, c, c,
						)
						currentStringBuilder.WriteRune('\\')
						currentStringBuilder.WriteRune(c)
					}
					isEscaped = false
				} else {
					switch c {
					case '_', '*', '~':
						ParagraphInsertRawHelperBuilder(contentTree, &currentStringBuilder)
						delimCount := 0
						for parsingNow[cI] == c {
							delimCount++
							cI++
							if cI >= len(parsingNow) {
								break
							}
						}
						cI--
						var tt InlineDelimiterType
						switch c {
						case '_':
							tt = InlineDelimiterTypeUnderlineDelimiter
						case '*':
							tt = InlineDelimiterTypeAsteriskDelimiter
						case '~':
							tt = InlineDelimiterTypeTildeDelimiter
						}
						for top := delimiterStack.Top(); top != nil && top.Type == tt && delimCount > 0; top = delimiterStack.Top() {
							toExtract := 1
							if top.Count >= 2 && delimCount >= 2 {
								toExtract = 2
							}
							i := len(contentTree.Children) - 1
							for ; i >= 0; i-- {
								ctc := contentTree.Children[i].Value
								if reflect.TypeOf(ctc) != reflect.TypeOf(&BlockInline{}) {
									continue
								}
								ctcbic := ctc.(*BlockInline).Content
								if reflect.TypeOf(ctcbic) != reflect.TypeOf(&InlineStringDelimiter{}) {
									continue
								}
								sm := ctcbic.(*InlineStringDelimiter).TypeOfDelimiter
								if sm != tt {
									continue
								}
								break
							}
							if i == -1 {
								panic(nil) // This should never be reached
							}
							// Create string modifier
							mod := new(InlineStringModifier)
							if tt == InlineDelimiterTypeTildeDelimiter {
								mod.TypeOfModifier = InlineStringModifierTypeStrikeoutText
							} else if toExtract == 2 {
								mod.TypeOfModifier = InlineStringModifierTypeBoldText
							} else {
								mod.TypeOfModifier = InlineStringModifierTypeItalicText
							}
							modTree := &Tree[BlockInterface]{
								Parent: contentTree,
								Value: &BlockInline{
									Content: mod,
								},
								Children: make([]*Tree[BlockInterface], len(contentTree.Children)-(i+1)),
							}
							for j := 0; j < len(modTree.Children); j++ {
								modTree.Children[j] = contentTree.Children[j+i+1]
								modTree.Children[j].Parent = modTree
							}
							// Remove temp delimiters
							contentTree.Children = contentTree.Children[:i-(toExtract-1)]
							// Insert string modifier
							contentTree.Children = append(contentTree.Children, modTree)

							delimCount -= toExtract
							top.Count -= toExtract
							if top.Count == 0 {
								delimiterStack.Pop()
							}
						}
						if delimCount > 0 {
							delimiterStack.Push(InlineDelimiter{
								Type:  tt,
								Count: delimCount,
							})
							for ; delimCount > 0; delimCount-- {
								contentTree.Children = append(
									contentTree.Children,
									&Tree[BlockInterface]{
										Parent: contentTree,
										Value: &BlockInline{
											Content: &InlineStringDelimiter{
												TypeOfDelimiter: tt,
											},
										},
									})
							}
						}
					case '\\':
						isEscaped = true
					default:
						currentStringBuilder.WriteRune(c)
					}
				}
			}
		}
	}
	return contentTree
}

const (
	PDMLinkStateStart             = 1
	PDMLinkStateTextLastSqBracket = 2
	PDMLinkStateTextLastSpace     = 3
	PDMLinkStateLink              = 4
)

type PDMLinkStateChar int

const (
	PDMLinkStateCharBracket PDMLinkStateChar = iota
	PDMLinkStateCharParenthesis
)

func (c PDMLinkStateChar) String() string {
	switch c {
	case PDMLinkStateCharBracket:
		return "["
	case PDMLinkStateCharParenthesis:
		return "]"
	default:
		panic(nil) // This should never be reached
	}
}

func ParseSingleParagraphLinks(tree *Tree[BlockInterface], file *FileStruct) *Tree[BlockInterface] {
	p := tree.Value.(*BlockParagraph)

	contentTree := &Tree[BlockInterface]{
		Parent: tree,
		Value: &BlockInline{
			Content: &InlineDocument{},
		},
	}

	var currentStringBuilder strings.Builder
	var currentStringLastChar rune
	isEscaped := false

	PdmState := PDMLinkStateStart
	var PdmStack Stack[Pair[PDMLinkStateChar, int]]
	PdmTextBegin := -1

	for current, expectedChildI := (Pair[int, int]{i: p.ContentStart.i, j: p.ContentStart.j - 1}), 0; ; {
		current.j++
		if current.j >= len(file.Lines[current.i].RuneContent) {
			if currentStringBuilder.Len() > 0 && currentStringLastChar != ' ' {
				currentStringBuilder.WriteRune(' ')
				currentStringLastChar = ' '
			}
			current.j = 0
			for current.i++; current.i < len(file.Lines) && file.Lines[current.i].Empty(); {
				current.i++
			}
		}
		if current.i > p.ContentEnd.i || (current.i == p.ContentEnd.i && current.j >= p.ContentEnd.j) {
			ParagraphInsertRawHelperBuilder(contentTree, &currentStringBuilder)
			break
		}

		isPartOfOtherChild := false
		if expectedChildI < len(tree.Children) {
			cbs := tree.Children[expectedChildI].Value.GetBlockStruct()
			if current.i == cbs.Start.i && current.j == cbs.Start.j {
				isPartOfOtherChild = true
			}
		}

		if isPartOfOtherChild {
			ParagraphInsertRawHelperBuilder(contentTree, &currentStringBuilder)
			treeChild := tree.Children[expectedChildI]
			treeChild.Parent = contentTree
			contentTree.Children = append(contentTree.Children, treeChild)
			current = tree.Children[expectedChildI].Value.GetBlockStruct().End
			current.j--
			expectedChildI++

			switch PdmState {
			case PDMLinkStateStart, PDMLinkStateTextLastSqBracket, PDMLinkStateTextLastSpace:
				PdmState = PDMLinkStateTextLastSpace
			case PDMLinkStateLink:
				PdmState = PDMLinkStateStart
				PdmStack.Clear()
			default:
				panic(nil) // This should never be reached!
			}
		} else {
			c := file.Lines[current.i].RuneContent[current.j]
			currentStringLastChar = c
			if isEscaped {
				switch c {
				case '_', '*', '|', '~', '<', '>', '\\':
					currentStringBuilder.WriteRune(c)
				default:
					log.Warnf(
						"Unrecognized escape sequence \"\\%v\". Please use \"\\\\%v\" instead. The sequence will be treated as \"\\\\%v\"...",
						c, c, c,
					)
					currentStringBuilder.WriteRune('\\')
					currentStringBuilder.WriteRune(c)
				}
				isEscaped = false
			} else {
				switch c {
				case '[':
					switch PdmState {
					case PDMLinkStateStart, PDMLinkStateTextLastSqBracket, PDMLinkStateTextLastSpace:
						ParagraphInsertRawHelperBuilder(contentTree, &currentStringBuilder)
						contentTree.Children = append(
							contentTree.Children,
							&Tree[BlockInterface]{
								Parent: contentTree,
								Value: &BlockInline{
									Content: &InlineStringDelimiter{
										TypeOfDelimiter: InlineDelimiterTypeOpenBracketDelimiter,
									},
								},
							})
						PdmState = PDMLinkStateStart
						PdmStack.Push(Pair[PDMLinkStateChar, int]{
							i: PDMLinkStateCharBracket,
							j: len(contentTree.Children) - 1,
						})
					case PDMLinkStateLink:
						// Do nothing
						currentStringBuilder.WriteRune('[')
					default:
						panic(nil) // This should never be reached!
					}
				case ']':
					switch PdmState {
					case PDMLinkStateStart, PDMLinkStateTextLastSqBracket, PDMLinkStateTextLastSpace:
						ParagraphInsertRawHelperBuilder(contentTree, &currentStringBuilder)
						if !PdmStack.Empty() && PdmStack.Top().i == PDMLinkStateCharBracket {
							PdmState = PDMLinkStateTextLastSqBracket
							PdmTextBegin = PdmStack.Top().j
							PdmStack.Pop()
						} else {
							PdmState = PDMLinkStateStart
							PdmStack.Clear()
						}
					case PDMLinkStateLink:
						// Do nothing
					default:
						panic(nil) // This should never be reached!
					}
					currentStringBuilder.WriteRune(']')
				case ' ':
					switch PdmState {
					case PDMLinkStateStart, PDMLinkStateTextLastSqBracket, PDMLinkStateTextLastSpace, PDMLinkStateLink:
						PdmState = PDMLinkStateTextLastSpace
					default:
						panic(nil) // This should never be reached!
					}
					currentStringBuilder.WriteRune(' ')
				case '(':
					switch PdmState {
					case PDMLinkStateStart, PDMLinkStateTextLastSpace:
						PdmState = PDMLinkStateTextLastSpace
					case PDMLinkStateTextLastSqBracket, PDMLinkStateLink:
						PdmState = PDMLinkStateLink
						PdmStack.Push(Pair[PDMLinkStateChar, int]{
							i: PDMLinkStateCharParenthesis,
							j: 0, // unused
						})
					default:
						panic(nil) // This should never be reached!
					}
					currentStringBuilder.WriteRune('(')
				case ')':
					switch PdmState {
					case PDMLinkStateStart, PDMLinkStateTextLastSpace, PDMLinkStateTextLastSqBracket:
						PdmState = PDMLinkStateTextLastSpace
					case PDMLinkStateLink:
						if !PdmStack.Empty() && PdmStack.Top().i == PDMLinkStateCharParenthesis {
							PdmState = PDMLinkStateLink
							PdmStack.Pop()
						} else {
							PdmState = PDMLinkStateStart
						}
					default:
						panic(nil) // This should never be reached!
					}

					if (PdmStack.Empty() || PdmStack.Size() > 0 && PdmStack.Top().i != PDMLinkStateCharParenthesis) && PdmState == PDMLinkStateLink {
						PdmStack.Clear()

						a := new(InlineHref)
						a.Address = currentStringBuilder.String()[2:]
						aTree := &Tree[BlockInterface]{
							Parent: contentTree,
							Value: &BlockInline{
								Content: a,
							},
							Children: make([]*Tree[BlockInterface], len(contentTree.Children)-(PdmTextBegin+1)),
						}
						for j := 0; j < len(aTree.Children); j++ { // TODO - could use cleanup instead of this
							wasInlineStringDelimiter := false
							if reflect.TypeOf(contentTree.Children[j+PdmTextBegin+1].Value) == reflect.TypeOf(&BlockInline{}) {
								bic := contentTree.Children[j+PdmTextBegin+1].Value.(*BlockInline).Content
								if reflect.TypeOf(bic) == reflect.TypeOf(&InlineStringDelimiter{}) {
									t := bic.(*InlineStringDelimiter).TypeOfDelimiter
									// Sanity check
									if t != InlineDelimiterTypeOpenBracketDelimiter {
										panic(nil)
									}
									wasInlineStringDelimiter = true
									aTree.Children[j] = &Tree[BlockInterface]{
										Parent: contentTree,
										Value: &BlockInline{
											Content: &InlineRawString{
												Content: "[",
											},
										},
										Children: nil,
									}
								}
							}
							if !wasInlineStringDelimiter {
								aTree.Children[j] = contentTree.Children[j+PdmTextBegin+1]
							}
							aTree.Children[j].Parent = aTree
						}
						// Remove temp delimiters
						contentTree.Children = contentTree.Children[:PdmTextBegin]
						// Insert string modifier
						contentTree.Children = append(contentTree.Children, aTree)
						// Clear current string
						currentStringBuilder.Reset()
						PdmState = PDMLinkStateStart
					} else {
						currentStringBuilder.WriteRune(')')
					}
				default:
					switch PdmState {
					case PDMLinkStateStart, PDMLinkStateTextLastSpace, PDMLinkStateTextLastSqBracket:
						PdmState = PDMLinkStateTextLastSpace
					case PDMLinkStateLink:
						PdmState = PDMLinkStateLink
					default:
						panic(nil) // This should never be reached!
					}
					currentStringBuilder.WriteRune(c)
				}
			}
		}
	}
	return contentTree
}
