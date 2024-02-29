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
	InsideTextbox   uint32
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

func (file FileStruct) MDParse() (Tree[BlockInterface], MDMetaStructure) {
	parsing_stack := ParsingStack{}
	current_block := &(Tree[BlockInterface]{})
	current_block.Value = &BlockDocument{}
	allowed_inside_block := current_block.Value.GetBlocksAllowedInside()

	just_skipped := INF_BLANKS
	for i := range file.Lines {
		if file.Lines[i].Empty() {
			just_skipped++
			continue
		}
		for file.Lines[i].RuneJ < len(file.Lines[i].RuneContent) {
			if file.Lines[i].RuneJ > 0 {
				just_skipped = 0
			}
			// Check if we can end this block
			for file.Lines[i].RuneJ < len(file.Lines[i].RuneContent) {
				should_end_block, should_discard_block, should_discard_seek := current_block.Value.CheckBlockEndsNormally(&file.Lines[i], parsing_stack)
				if !should_end_block {
					should_end_block = current_block.Value.CheckBlockEndsViaNewLinesAndIndentation(just_skipped, file.Lines[i].Indentation)
				}
				if !should_end_block {
					break
				}
				log.Debugf(
					"<<< End of %v (line-index=%v, line-j=%v)",
					current_block.Value,
					i, file.Lines[i].RuneJ,
				)

				switch reflect.TypeOf(current_block.Value) {
				case reflect.TypeOf(BlockParagraph{}):
					parsing_stack.InsideParagraph--
				case reflect.TypeOf(BlockTextbox{}):
					parsing_stack.InsideTextbox--
				}

				// Discard if needed
				if should_discard_block != nil {
					file.Lines[i].RuneJ += should_discard_seek
					for reflect.TypeOf(current_block.Value) != reflect.TypeOf(should_discard_block) {
						fmt.Printf("SOFT Discarding %v (%T != %T)\n", current_block.Value, current_block.Value, should_discard_block)
						// TODO - remove
						current_block.Value.ExecuteAfterBlockEnds(&file.Lines[i])
						current_block = current_block.Parent
					}
					file.Lines[i].RuneJ -= should_discard_seek
				}

				file.Lines[i].RuneJ += current_block.Value.SeekBufferAfterBlockEnds()
				current_block.Value.ExecuteAfterBlockEnds(&file.Lines[i])
				current_block = current_block.Parent

				allowed_inside_block = current_block.Value.GetBlocksAllowedInside()
			}
			if file.Lines[i].RuneJ >= len(file.Lines[i].RuneContent) {
				break
			}
			// Check if we can open blocks
			for found := true; found && file.Lines[i].RuneJ < len(file.Lines[i].RuneContent); {
				found = false
				for _, allowed := range allowed_inside_block {
					if allowed.CheckBlockStarts(file.Lines[i]) && current_block.Value.AcceptBlockInside(allowed) {
						file.Lines[i].RuneJ += allowed.SeekBufferAfterBlockStarts()
						next_block := &(Tree[BlockInterface]{})
						next_block.Value = allowed
						next_block.Parent = current_block
						current_block.Children = append(current_block.Children, next_block)
						current_block = next_block
						allowed_inside_block = allowed.GetBlocksAllowedInside()
						found = true
						break
					}
				}
				if !found {
					file.Lines[i].RuneJ++
				} else {
					switch reflect.TypeOf(current_block.Value) {
					case reflect.TypeOf(&BlockParagraph{}):
						parsing_stack.InsideParagraph++
					case reflect.TypeOf(&BlockTextbox{}):
						parsing_stack.InsideTextbox++
					}
					current_block.Value.ExecuteAfterBlockStarts(&file.Lines[i])
					log.Debugf(
						">>> Beginning of %v (line-index=%v, line-j=%v)",
						current_block.Value,
						i, file.Lines[i].RuneJ,
					)
				}
			}
		}
		just_skipped = 0
	}

	// Get back to root
	for current_block.Parent != nil {
		switch reflect.TypeOf(current_block.Value) {
		case reflect.TypeOf(BlockParagraph{}):
			parsing_stack.InsideParagraph--
		case reflect.TypeOf(BlockTextbox{}):
			parsing_stack.InsideTextbox--
		}
		current_block.Value.ExecuteAfterBlockEnds(&file.Lines[len(file.Lines)-1])
		current_block = current_block.Parent
	}
	current_block.Value.GetBlockStruct().ContentEnd = Pair[int, int]{
		i: len(file.Lines) - 1,
		j: len(file.Lines[len(file.Lines)-1].RuneContent),
	}

	// Find first non-paragraph
	FirstNonParagraph := func(tree *Tree[BlockInterface], min_i int) int {
		for i := min_i; i < len(tree.Children); i++ {
			if !tree.Children[i].Value.IsPartOfParagraph() {
				return i
			}
		}
		return -1
	}

	// Create the number of detected paragraphs
	DetectParagraphs := func(before, after Pair[int, int], file *FileStruct) []BlockParagraph {
		var paragraphs []BlockParagraph
		var current_paragraph string

		for i := before.i; i <= after.i; i++ {
			var j = 0
			var max_j = len(file.Lines[i].RuneContent)
			if i == before.i {
				j = before.j
			}
			if i == after.i {
				max_j = after.j
			}
			current_paragraph += string(file.Lines[i].RuneContent[j:max_j]) + " "
			if file.Lines[i].Empty() || i == after.i {
				current_paragraph = RemoveExcessSpaces(current_paragraph)
				if current_paragraph != "" {
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
				current_paragraph = ""
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

		starting_i := 0
		for i := FirstNonParagraph(tree, 0); i != -1; i = FirstNonParagraph(tree, starting_i) {
			after = tree.Children[i].Value.GetBlockStruct().Start
			next_before := tree.Children[i].Value.GetBlockStruct().End
			if p := DetectParagraphs(before, after, &file); p != nil {
				log.Debug("Detected paragraph between ", before, " and ", after, ": ", p)
				generated_p := make([]*Tree[BlockInterface], len(p))
				for j := range p {
					generated_p[j] = new(Tree[BlockInterface])
					generated_p[j].Parent = tree
					generated_p[j].Value = &p[j]
				}
				for ci, j := starting_i, 0; ci < i && j < len(p); ci++ {
					for cbs := tree.Children[ci].Value.GetBlockStruct(); cbs.Start.i > p[j].ContentEnd.i; j++ {
					}
					tree.Children[ci].Parent = generated_p[j]
					generated_p[j].Children = append(generated_p[j].Children, tree.Children[ci])
				}
				// Remove old blocks
				aux := make([]*Tree[BlockInterface], len(p)+len(tree.Children)-(i-starting_i))
				copy(aux[:starting_i], tree.Children[:starting_i])
				copy(aux[starting_i:starting_i+len(p)], generated_p)
				copy(aux[starting_i+len(p):], tree.Children[i:])
				tree.Children = aux
				starting_i += len(p)
			}
			starting_i++
			before = next_before
		}

		after = tree.Value.GetBlockStruct().ContentEnd
		if p := DetectParagraphs(before, after, &file); p != nil {
			log.Debug("Detected paragraph between ", before, " and ", after, ": ", p)
			generated_p := make([]*Tree[BlockInterface], len(p))
			for j := range p {
				generated_p[j] = new(Tree[BlockInterface])
				generated_p[j].Parent = tree
				generated_p[j].Value = &p[j]
			}
			for ci, j := starting_i, 0; ci < len(tree.Children) && j < len(p); ci++ {
				for cbs := tree.Children[ci].Value.GetBlockStruct(); cbs.Start.i > p[j].ContentEnd.i; j++ {
				}
				tree.Children[ci].Parent = generated_p[j]
				generated_p[j].Children = append(generated_p[j].Children, tree.Children[ci])
			}
			// Remove old blocks
			aux := make([]*Tree[BlockInterface], len(p)+len(tree.Children)-(len(tree.Children)-starting_i))
			copy(aux[:starting_i], tree.Children[:starting_i])
			copy(aux[starting_i:starting_i+len(p)], generated_p)
			copy(aux[starting_i+len(p):], tree.Children[len(tree.Children):])
			tree.Children = aux
			starting_i += len(p)
		}
	}
	InsertParagraphs(current_block)

	var md_meta MDMetaStructure

	// Parse paragraphs, raw content, subfigures and meta info
	CompleteBibinfo := make(map[string]BibliographyEntry)
	CompleteBibinfoNextRefIndex := 1

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
		case reflect.TypeOf(&BlockSubfigure{}):
			// Sanity check
			if reflect.TypeOf(tree.Parent.Value) != reflect.TypeOf(&BlockFigure{}) {
				panic(nil)
			}
			if tree.Value.(*BlockSubfigure).Padding == "" {
				tree.Value.(*BlockSubfigure).Padding = tree.Parent.Value.(*BlockFigure).Padding
			}
		case reflect.TypeOf(&BlockMetaAuthor{}):
			// Sanity check
			if reflect.TypeOf(tree.Parent.Value) != reflect.TypeOf(&BlockMeta{}) {
				panic(nil)
			}
			md_meta.Author = *tree.Value.GetRawContent()
		case reflect.TypeOf(&BlockMetaCopyright{}):
			// Sanity check
			if reflect.TypeOf(tree.Parent.Value) != reflect.TypeOf(&BlockMeta{}) {
				panic(nil)
			}
			md_meta.Copyright = *tree.Value.GetRawContent()
		case reflect.TypeOf(&BlockMetaBibinfo{}):
			// Sanity check
			if reflect.TypeOf(tree.Parent.Value) != reflect.TypeOf(&BlockMeta{}) {
				panic(nil)
			}
			new_bibinfo := ParseBibinfo(tree.Value.(*BlockMetaBibinfo))
			ordered_new_bibinfo := make([]string, 0)
			for k := range new_bibinfo {
				ordered_new_bibinfo = append(ordered_new_bibinfo, k)
			}
			sort.Strings(ordered_new_bibinfo)
			for _, new_bibinfo_key := range ordered_new_bibinfo {
				new_bibinfo_value := new_bibinfo[new_bibinfo_key]
				if old_value, ok := (CompleteBibinfo[new_bibinfo_key]); ok {
					log.Warnf("A bibliography entry for tag \"%v\" already exists (%v). Keeping the old value...", new_bibinfo_key, old_value)
					continue
				}
				CompleteBibinfo[new_bibinfo_key] = BibliographyEntry{
					ParentBlock:    new_bibinfo_value.ParentBlock,
					Type:           new_bibinfo_value.Type,
					Fields:         new_bibinfo_value.Fields,
					ReferenceIndex: CompleteBibinfoNextRefIndex,
					File:           tree.Value.(*BlockMetaBibinfo).RefFile,
				}
				CompleteBibinfoNextRefIndex++
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
			tabs_tab := tree.Value.(*BlockTabsTab)
			tabs.Tabs = append(tabs.Tabs, tabs_tab)
		case reflect.TypeOf(&BlockTabs{}):
			tabs := tree.Value.(*BlockTabs)
			TabsBlocks = append(TabsBlocks, tabs)
		}

		for i := 0; i < len(tree.Children); i++ {
			DFSNodeParse(tree.Children[i])
		}
	}
	DFSNodeParse(current_block)

	// Convert references
	for i := 0; i < len(RefBlocks); i++ {
		var entry BibliographyEntry
		var ok bool
		if entry, ok = CompleteBibinfo[RefBlocks[i].RawContent]; !ok {
			log.Warnf("Reference to \"%v\" cannot be resolved (will be rendered as \"[?]\")", RefBlocks[i].RawContent)
			continue
		}
		RefBlocks[i].ReferenceIndex = entry.ReferenceIndex
		if RefBlocks[i].File == "" {
			RefBlocks[i].File = entry.File
		}
	}

	// Generate bibliographies
	bibliography_html_text := GenerateBibliography(CompleteBibinfo)
	for i := 0; i < len(BibliographyBlocks); i++ {
		BibliographyBlocks[i].HTMLContent = &bibliography_html_text
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

	return *current_block, md_meta
}

// =====================================
// Bibliography parse

type BibliographyEntryType uint8

const (
	BibliographyEntryType_Article BibliographyEntryType = iota
	BibliographyEntryType_Book
	BibliographyEntryType_Other
)

func (t BibliographyEntryType) String() string {
	switch t {
	case BibliographyEntryType_Article:
		return "Article"
	case BibliographyEntryType_Book:
		return "Book"
	case BibliographyEntryType_Other:
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
	ParentBlock    *BlockMetaBibinfo
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

func ParseBibinfo(b *BlockMetaBibinfo) map[string]BibliographyEntry {
	entries := make(map[string]BibliographyEntry)
	bibinfo_type_str := ""

	// Types - don't call reflect.TypeOf multiple times!
	reflect_type_string := reflect.TypeOf("a")
	reflect_type_map_string_interface := reflect.TypeOf(make(map[string]interface{}))
	reflect_type_slice_interface := reflect.TypeOf(make([]interface{}, 0))

	var json_sb []byte
	if b.JSONInline {
		bibinfo_type_str = "inline bibinfo \"" + b.RawContent + "\""
		json_sb = []byte(b.RawContent)
	} else {
		bibinfo_type_str = "bibinfo file (" + b.RawContent + ")"
		var err error
		json_sb, err = os.ReadFile(b.RawContent)
		if err != nil {
			log.Errorf("Could not read %v. Skipping nil...", bibinfo_type_str)
			return nil
		}
	}

	// Unmarshall the JSON
	var full_json_interface interface{}
	err := json.Unmarshal(json_sb, &full_json_interface)
	if err != nil {
		log.Errorf("Could not unmarshal %v. Skipping nil...", bibinfo_type_str)
		return nil
	}
	full_json := full_json_interface.(map[string]interface{})

	// Search for "bibliography"
	if _, ok := (full_json["bibliography"]); !ok {
		log.Warnf("Could not find \"bibliography\" entry in %v. Skipping empty...", bibinfo_type_str)
		return entries
	}
	if reflect.TypeOf(full_json["bibliography"]) != reflect_type_slice_interface {
		log.Warnf("The \"bibliography\" in %v is NOT an array. Skipping empty...", bibinfo_type_str)
		return entries
	}
	bibliography_json := full_json["bibliography"].([]interface{})

	// Search for entries
	if len(bibliography_json) == 0 {
		log.Warnf("The \"bibliography\" in %v is empty.", bibinfo_type_str)
		return entries
	}
	for entry_i := 0; entry_i < len(bibliography_json); entry_i++ {
		if reflect.TypeOf(bibliography_json[entry_i]) != reflect_type_map_string_interface {
			log.Warnf("\"bibliography\" entry %v in %v will be ignored (not a valid map[string]...)", entry_i, bibinfo_type_str)
			continue
		}
		bibliography_entry := bibliography_json[entry_i].(map[string]interface{})

		var entry BibliographyEntry
		entry.ParentBlock = b

		entry_tag := ""
		if bibliography_entry_tag, ok := (bibliography_entry["tag"]); !ok {
			log.Warnf("\"bibliography\" entry %v in %v will be ignored (missing tag)", entry_i, bibinfo_type_str)
			continue
		} else if reflect.TypeOf(bibliography_entry_tag) != reflect_type_string {
			log.Warnf("\"bibliography\" entry %v in %v will be ignored (tag exists, but is not a string)", entry_i, bibinfo_type_str)
			continue
		} else {
			entry_tag = bibliography_entry_tag.(string)
		}

		entry_type := ""
		entry.Type = BibliographyEntryType_Other
		if bibliography_entry_type, ok := (bibliography_entry["type"]); !ok {
			log.Warnf("Missing \"bibliography\" entry type in entry %v, %v. Considering type \"other\"...", entry_i, bibinfo_type_str)
		} else if reflect.TypeOf(bibliography_entry_type) != reflect_type_string {
			log.Warnf("\"bibliography\" entry type is not a string in entry %v, %v. Considering type \"other\"...", entry_i, bibinfo_type_str)
		} else {
			entry_type = bibliography_entry_type.(string)
			switch entry_type {
			case "article", "Article":
				entry.Type = BibliographyEntryType_Article
			case "book", "Book":
				entry.Type = BibliographyEntryType_Book
			case "other", "Other", "unknown", "Unknown":
			default:
				log.Warnf("Unrecognized \"bibliography\" entry type \"%v\" in entry %v, %v. Considering type \"other\"...", entry_type, entry_i, bibinfo_type_str)
			}
		}

		if bibliography_entry_data, ok := (bibliography_entry["data"]); ok && reflect.TypeOf(bibliography_entry_data) != reflect_type_map_string_interface {
			log.Warnf("\"bibliography\" entry %v's \"data\" in %v will be ignored (not valid map[string]...)", entry_i, bibinfo_type_str)
		} else if ok {
			entry_data := bibliography_entry_data.(map[string]interface{})
			for entry_data_key, entry_data_value := range entry_data {
				if reflect.TypeOf(entry_data_value) != reflect_type_string {
					log.Warnf("\"bibliography\" entry data value must be string in entry %v, %v. Ignoring only data key %v...", entry_i, bibinfo_type_str, entry_data_key)
					continue
				}
				data_value := entry_data_value.(string)
				switch entry_data_key {
				case "title":
					entry.Fields.Title = new(string)
					*entry.Fields.Title = data_value
				case "author":
					entry.Fields.Author = new(string)
					*entry.Fields.Author = data_value
				case "journal":
					entry.Fields.Journal = new(string)
					*entry.Fields.Journal = data_value
				case "volume":
					entry.Fields.Volume = new(string)
					*entry.Fields.Volume = data_value
				case "number":
					entry.Fields.Number = new(string)
					*entry.Fields.Number = data_value
				case "pages":
					entry.Fields.Pages = new(string)
					*entry.Fields.Pages = data_value
				case "year":
					entry.Fields.Year = new(string)
					*entry.Fields.Year = data_value
				case "publisher":
					entry.Fields.Publisher = new(string)
					*entry.Fields.Publisher = data_value
				case "url":
					entry.Fields.URL = new(string)
					*entry.Fields.URL = data_value
				default:
					log.Warnf("Unrecognized \"bibliography\" entry data key \"%v\" in entry %v, %v. Ignoring only that data key...", entry_data_key, entry_i, bibinfo_type_str)
				}
			}
		}
		entries[entry_tag] = entry
	}
	return entries
}

func GenerateBibliography(mp map[string]BibliographyEntry) string {
	var sb strings.Builder
	sb.WriteString("<div class=\"bibliography\">\n")
	ordered_bibinfo := make([]string, 0)
	for k := range mp {
		ordered_bibinfo = append(ordered_bibinfo, k)
	}
	sort.Strings(ordered_bibinfo)
	for _, key := range ordered_bibinfo {
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

func ParagraphInsertRawHelper(tree *Tree[BlockInterface], s *string) {
	if *s != "" {
		tree.Children = append(tree.Children,
			&Tree[BlockInterface]{
				Parent: tree,
				Value: &BlockInline{
					Content: &InlineRawString{
						Content: *s,
					},
				},
			})
		*s = ""
	}
}

func ParseSingleParagraph(tree *Tree[BlockInterface], file FileStruct) {
	// Insert links
	link_parsed_tree := ParseSingleParagraphLinks(tree, file)
	// Insert emphasis
	emphasis_parsed_tree := ParseSingleBlockInlineEmphasis(link_parsed_tree)
	// Cleanup
	cleaned_tree := CleanupSingleBlockInline(emphasis_parsed_tree)

	// Convert heading-paragraph -> normal heading
	//         textbox title-paragraph into normal textbox title
	tree.Children = []*Tree[BlockInterface]{cleaned_tree}
	cleaned_tree.Parent = tree
	if tree.Parent != nil && (reflect.TypeOf(tree.Parent.Value) == reflect.TypeOf(&BlockHeading{}) || reflect.TypeOf(tree.Parent.Value) == reflect.TypeOf(&BlockTextboxTitle{})) {
		h := tree.Parent
		cleaned_tree.Parent = h
		h.Children = []*Tree[BlockInterface]{cleaned_tree}
	}
}

func CleanupSingleBlockInline(content_tree *Tree[BlockInterface]) *Tree[BlockInterface] {
	cleaned_content_tree := &Tree[BlockInterface]{
		Parent: content_tree.Parent,
		Value: &BlockInline{
			Content: &InlineDocument{},
		},
	}
	for i := 0; i < len(content_tree.Children); i++ {
		var sr string = ""
		j := i
	digging:
		for ; j < len(content_tree.Children); j++ {
			if reflect.TypeOf(content_tree.Children[j].Value) != reflect.TypeOf(&BlockInline{}) {
				break
			}
			bjc := content_tree.Children[j].Value.(*BlockInline).Content
			switch reflect.TypeOf(bjc) {
			case reflect.TypeOf(&InlineStringDelimiter{}):
				switch bjc.(*InlineStringDelimiter).TypeOfDelimiter {
				case UnderlineDelimiter:
					sr += "_"
				case AsteriskDelimiter:
					sr += "*"
				case TildeDelimiter:
					sr += "~"
				case OpenBracketDelimiter:
					sr += "["
				case CloseBracketDelimiter:
					sr += "]"
				case OpenParantDelimiter:
					sr += "("
				case CloseParantDelimiter:
					sr += ")" // Might never be reached
				default:
					panic(nil) // This should never be reached
				}
			case reflect.TypeOf(&InlineRawString{}):
				sr += bjc.(*InlineRawString).Content
			default:
				break digging
			}
		}
		if j != i {
			cleaned_content_tree.Children = append(
				cleaned_content_tree.Children,
				&Tree[BlockInterface]{
					Parent: cleaned_content_tree,
					Value: &BlockInline{
						Content: &InlineRawString{
							Content: sr,
						},
					},
				},
			)
		}
		if j < len(content_tree.Children) {
			content_tree.Children[j].Parent = cleaned_content_tree
			cleaned_content_tree.Children = append(
				cleaned_content_tree.Children,
				content_tree.Children[j],
			)
		}
		i = j
	}

	return cleaned_content_tree
}

func ParseSingleBlockInlineEmphasis(tree *Tree[BlockInterface]) *Tree[BlockInterface] {
	content_tree := &Tree[BlockInterface]{
		Parent: tree,
		Value: &BlockInline{
			Content: &InlineDocument{},
		},
	}

	var delimiter_stack Stack[InlineDelimiter]
	current_string := ""
	for child_i := 0; ; child_i++ {
		if child_i == len(tree.Children) {
			ParagraphInsertRawHelper(content_tree, &current_string)
			break
		}
		tree_child := tree.Children[child_i]

		// Only work with raw strings
		is_raw := false
		is_href := false
		if reflect.TypeOf(tree_child.Value) == reflect.TypeOf(&BlockInline{}) {
			bic := tree_child.Value.(*BlockInline).Content
			switch reflect.TypeOf(bic) {
			case reflect.TypeOf(&InlineRawString{}):
				is_raw = true
			case reflect.TypeOf(&InlineHref{}):
				is_href = true
			}
		}

		if !is_raw {
			ParagraphInsertRawHelper(content_tree, &current_string)
			tree_child.Parent = content_tree
			if is_href {
				aux := ParseSingleBlockInlineEmphasis(tree_child)
				for i := 0; i < len(aux.Children); i++ {
					aux.Children[i].Parent = tree_child
				}
				tree_child.Children = aux.Children
			}
			content_tree.Children = append(content_tree.Children, tree_child)
		} else {
			parsing_now := []rune(tree_child.Value.(*BlockInline).Content.(*InlineRawString).Content)
			is_escaped := false
			for c_i := 0; c_i < len(parsing_now); c_i++ {
				c := parsing_now[c_i]
				if is_escaped {
					switch c {
					case '_', '*', '|', '~', '<', '>', '\\':
						current_string += string(c)
					default:
						log.Warnf(
							"Unrecognized escape sequence \"\\%v\". Please use \"\\\\%v\" instead. The sequence will be treated as \"\\\\%v\"...",
							c, c, c,
						)
						current_string += "\\" + string(c)
					}
					is_escaped = false
				} else {
					switch c {
					case '_', '*', '~':
						ParagraphInsertRawHelper(content_tree, &current_string)
						delim_count := 0
						for parsing_now[c_i] == c {
							delim_count++
							c_i++
							if c_i >= len(parsing_now) {
								break
							}
						}
						c_i--
						var tt InlineDelimiterType
						switch c {
						case '_':
							tt = UnderlineDelimiter
						case '*':
							tt = AsteriskDelimiter
						case '~':
							tt = TildeDelimiter
						}
						for top := delimiter_stack.Top(); top != nil && top.Type == tt && delim_count > 0; top = delimiter_stack.Top() {
							to_extract := 1
							if top.Count >= 2 && delim_count >= 2 {
								to_extract = 2
							}
							i := len(content_tree.Children) - 1
							for ; i >= 0; i-- {
								ctc := content_tree.Children[i].Value
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
							if tt == TildeDelimiter {
								mod.TypeOfModifier = StrikeoutText
							} else if to_extract == 2 {
								mod.TypeOfModifier = BoldText
							} else {
								mod.TypeOfModifier = ItalicText
							}
							mod_tree := &Tree[BlockInterface]{
								Parent: content_tree,
								Value: &BlockInline{
									Content: mod,
								},
								Children: make([]*Tree[BlockInterface], len(content_tree.Children)-(i+1)),
							}
							for j := 0; j < len(mod_tree.Children); j++ {
								mod_tree.Children[j] = content_tree.Children[j+i+1]
								mod_tree.Children[j].Parent = mod_tree
							}
							// Remove temp delimiters
							content_tree.Children = content_tree.Children[:i-(to_extract-1)]
							// Insert string modifier
							content_tree.Children = append(content_tree.Children, mod_tree)

							delim_count -= to_extract
							top.Count -= to_extract
							if top.Count == 0 {
								delimiter_stack.Pop()
							}
						}
						if delim_count > 0 {
							delimiter_stack.Push(InlineDelimiter{
								Type:  tt,
								Count: delim_count,
							})
							for ; delim_count > 0; delim_count-- {
								content_tree.Children = append(
									content_tree.Children,
									&Tree[BlockInterface]{
										Parent: content_tree,
										Value: &BlockInline{
											Content: &InlineStringDelimiter{
												TypeOfDelimiter: tt,
											},
										},
									})
							}
						}
					case '\\':
						is_escaped = true
					default:
						current_string += string(c)
					}
				}
			}
		}
	}

	return content_tree
}

const (
	PDM_LinkState_Start                = 1
	PDM_LinkState_Text_Last_Sq_Bracket = 2
	PDM_LinkState_Text_Last_Space      = 3
	PDM_LinkState_Link                 = 4
)

type PDM_LinkState_Char int

const (
	PDM_LinkState_Char_Bracket PDM_LinkState_Char = iota
	PDM_LinkState_Char_Parant
)

func (c PDM_LinkState_Char) String() string {
	switch c {
	case PDM_LinkState_Char_Bracket:
		return "["
	case PDM_LinkState_Char_Parant:
		return "]"
	default:
		panic(nil) // This should never be reached
	}
}

func ParseSingleParagraphLinks(tree *Tree[BlockInterface], file FileStruct) *Tree[BlockInterface] {
	p := tree.Value.(*BlockParagraph)

	content_tree := &Tree[BlockInterface]{
		Parent: tree,
		Value: &BlockInline{
			Content: &InlineDocument{},
		},
	}

	current_string := ""
	is_escaped := false

	PDM_State := PDM_LinkState_Start
	var PDM_Stack Stack[Pair[PDM_LinkState_Char, int]]
	PDM_text_begin := -1

	for current, expected_child_i := (Pair[int, int]{i: p.ContentStart.i, j: p.ContentStart.j - 1}), 0; ; {
		current.j++
		if current.j >= len(file.Lines[current.i].RuneContent) {
			if current_string != "" && current_string[len(current_string)-1] != ' ' {
				current_string += " "
			}
			current.j = 0
			for current.i++; current.i < len(file.Lines) && file.Lines[current.i].Empty(); {
				current.i++
			}
		}
		if current.i > p.ContentEnd.i || (current.i == p.ContentEnd.i && current.j >= p.ContentEnd.j) {
			ParagraphInsertRawHelper(content_tree, &current_string)
			break
		}

		is_part_of_other_child := false
		if expected_child_i < len(tree.Children) {
			cbs := tree.Children[expected_child_i].Value.GetBlockStruct()
			if current.i == cbs.Start.i && current.j == cbs.Start.j {
				is_part_of_other_child = true
			}
		}

		if is_part_of_other_child {
			ParagraphInsertRawHelper(content_tree, &current_string)
			tree_child := tree.Children[expected_child_i]
			tree_child.Parent = content_tree
			content_tree.Children = append(content_tree.Children, tree_child)
			current = tree.Children[expected_child_i].Value.GetBlockStruct().End
			current.j--
			expected_child_i++

			switch PDM_State {
			case PDM_LinkState_Start, PDM_LinkState_Text_Last_Sq_Bracket, PDM_LinkState_Text_Last_Space:
				PDM_State = PDM_LinkState_Text_Last_Space
			case PDM_LinkState_Link:
				PDM_State = PDM_LinkState_Start
				PDM_Stack.Clear()
			default:
				panic(nil) // This should never be reached!
			}
		} else {
			c := file.Lines[current.i].RuneContent[current.j]
			if is_escaped {
				switch c {
				case '_', '*', '|', '~', '<', '>', '\\':
					current_string += string(c)
				default:
					log.Warnf(
						"Unrecognized escape sequence \"\\%v\". Please use \"\\\\%v\" instead. The sequence will be treated as \"\\\\%v\"...",
						c, c, c,
					)
					current_string += "\\" + string(c)
				}
				is_escaped = false
			} else {
				switch c {
				case '[':
					switch PDM_State {
					case PDM_LinkState_Start, PDM_LinkState_Text_Last_Sq_Bracket, PDM_LinkState_Text_Last_Space:
						ParagraphInsertRawHelper(content_tree, &current_string)
						content_tree.Children = append(
							content_tree.Children,
							&Tree[BlockInterface]{
								Parent: content_tree,
								Value: &BlockInline{
									Content: &InlineStringDelimiter{
										TypeOfDelimiter: OpenBracketDelimiter,
									},
								},
							})
						PDM_State = PDM_LinkState_Start
						PDM_Stack.Push(Pair[PDM_LinkState_Char, int]{
							i: PDM_LinkState_Char_Bracket,
							j: len(content_tree.Children) - 1,
						})
					case PDM_LinkState_Link:
						// Do nothing
						current_string += "["
					default:
						panic(nil) // This should never be reached!
					}
				case ']':
					switch PDM_State {
					case PDM_LinkState_Start, PDM_LinkState_Text_Last_Sq_Bracket, PDM_LinkState_Text_Last_Space:
						ParagraphInsertRawHelper(content_tree, &current_string)
						if !PDM_Stack.Empty() && PDM_Stack.Top().i == PDM_LinkState_Char_Bracket {
							PDM_State = PDM_LinkState_Text_Last_Sq_Bracket
							PDM_text_begin = PDM_Stack.Top().j
							PDM_Stack.Pop()
						} else {
							PDM_State = PDM_LinkState_Start
							PDM_Stack.Clear()
						}
					case PDM_LinkState_Link:
						// Do nothing
					default:
						panic(nil) // This should never be reached!
					}
					current_string += "]"
				case ' ':
					switch PDM_State {
					case PDM_LinkState_Start, PDM_LinkState_Text_Last_Sq_Bracket, PDM_LinkState_Text_Last_Space, PDM_LinkState_Link:
						PDM_State = PDM_LinkState_Text_Last_Space
					default:
						panic(nil) // This should never be reached!
					}
					current_string += " "

				case '(':
					switch PDM_State {
					case PDM_LinkState_Start, PDM_LinkState_Text_Last_Space:
						PDM_State = PDM_LinkState_Text_Last_Space
					case PDM_LinkState_Text_Last_Sq_Bracket, PDM_LinkState_Link:
						PDM_State = PDM_LinkState_Link
						PDM_Stack.Push(Pair[PDM_LinkState_Char, int]{
							i: PDM_LinkState_Char_Parant,
							j: 0, // unused
						})
					default:
						panic(nil) // This should never be reached!
					}
					current_string += "("
				case ')':
					switch PDM_State {
					case PDM_LinkState_Start, PDM_LinkState_Text_Last_Space, PDM_LinkState_Text_Last_Sq_Bracket:
						PDM_State = PDM_LinkState_Text_Last_Space
					case PDM_LinkState_Link:
						if !PDM_Stack.Empty() && PDM_Stack.Top().i == PDM_LinkState_Char_Parant {
							PDM_State = PDM_LinkState_Link
							PDM_Stack.Pop()
						} else {
							PDM_State = PDM_LinkState_Start
						}
					default:
						panic(nil) // This should never be reached!
					}

					if (PDM_Stack.Empty() || PDM_Stack.Size() > 0 && PDM_Stack.Top().i != PDM_LinkState_Char_Parant) && PDM_State == PDM_LinkState_Link {
						PDM_Stack.Clear()

						a := new(InlineHref)
						a.Address = current_string[2:]
						a_tree := &Tree[BlockInterface]{
							Parent: content_tree,
							Value: &BlockInline{
								Content: a,
							},
							Children: make([]*Tree[BlockInterface], len(content_tree.Children)-(PDM_text_begin+1)),
						}
						for j := 0; j < len(a_tree.Children); j++ { // TODO - could use cleanup instead of this
							was_inline_string_delimiter := false
							if reflect.TypeOf(content_tree.Children[j+PDM_text_begin+1].Value) == reflect.TypeOf(&BlockInline{}) {
								bic := content_tree.Children[j+PDM_text_begin+1].Value.(*BlockInline).Content
								if reflect.TypeOf(bic) == reflect.TypeOf(&InlineStringDelimiter{}) {
									t := bic.(*InlineStringDelimiter).TypeOfDelimiter
									// Sanity check
									if t != OpenBracketDelimiter {
										panic(nil)
									}
									was_inline_string_delimiter = true
									a_tree.Children[j] = &Tree[BlockInterface]{
										Parent: content_tree,
										Value: &BlockInline{
											Content: &InlineRawString{
												Content: "[",
											},
										},
										Children: nil,
									}
								}
							}
							if !was_inline_string_delimiter {
								a_tree.Children[j] = content_tree.Children[j+PDM_text_begin+1]
							}
							a_tree.Children[j].Parent = a_tree
						}
						// Remove temp delimiters
						content_tree.Children = content_tree.Children[:PDM_text_begin]
						// Insert string modifier
						content_tree.Children = append(content_tree.Children, a_tree)
						// Clear current string
						current_string = ""
						PDM_State = PDM_LinkState_Start
					} else {
						current_string += ")"
					}
				default:
					switch PDM_State {
					case PDM_LinkState_Start, PDM_LinkState_Text_Last_Space, PDM_LinkState_Text_Last_Sq_Bracket:
						PDM_State = PDM_LinkState_Text_Last_Space
					case PDM_LinkState_Link:
						PDM_State = PDM_LinkState_Link
					default:
						panic(nil) // This should never be reached!
					}
					current_string += string(c)
				}
			}
		}
	}

	return content_tree
}
