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

import "fmt"

// =====================================
// Generic inline structure

type InlineInterface interface {
	fmt.Stringer

	HTMLInterface
	LaTeXInterface

	GetRawContent() *string
}

// =====================================
// Inline document

type InlineDocument struct {
}

func (b InlineDocument) String() string {
	return "InlineDocument"
}

func (b *InlineDocument) GetRawContent() *string {
	return nil
}

// =====================================
// Raw string

type InlineRawString struct {
	Content string
}

func (b InlineRawString) String() string {
	return "InlineRawString (\"" + b.Content + "\")"
}

func (b *InlineRawString) GetRawContent() *string {
	return nil
}

// =====================================
// String modifier

type InlineStringModifierType uint8

const (
	ItalicText InlineStringModifierType = iota
	BoldText
	StrikeoutText
)

func (t InlineStringModifierType) String() string {
	switch t {
	case ItalicText:
		return "ItalicText"
	case BoldText:
		return "BoldText"
	case StrikeoutText:
		return "StrikeoutText"
	default:
		panic(nil)
	}
}

type InlineStringModifier struct {
	TypeOfModifier InlineStringModifierType
}

func (b InlineStringModifier) String() string {
	return "InlineStringModifier (" + b.TypeOfModifier.String() + ")"
}

func (b *InlineStringModifier) GetRawContent() *string {
	return nil
}

// =====================================
// String delimiter

type InlineDelimiterType uint8

const (
	AsteriskDelimiter InlineDelimiterType = iota
	UnderlineDelimiter
	TildeDelimiter
	OpenBracketDelimiter
	CloseBracketDelimiter
	OpenParantDelimiter
	CloseParantDelimiter
)

func (t InlineDelimiterType) String() string {
	switch t {
	case AsteriskDelimiter:
		return "AsteriskDelimiter"
	case UnderlineDelimiter:
		return "UnderlineDelimiter"
	case TildeDelimiter:
		return "TildeDelimiter"
	case OpenBracketDelimiter:
		return "OpenBracketDelimiter"
	case CloseBracketDelimiter:
		return "CloseBracketDelimiter"
	case OpenParantDelimiter:
		return "OpenParantDelimiter"
	case CloseParantDelimiter:
		return "CloseParantDelimiter"
	}
	panic(nil) // This should never be reached
}

type InlineDelimiter struct {
	Type InlineDelimiterType
	Count int
}

func (d InlineDelimiter) String() string {
	return fmt.Sprintf(
		"InlineDelimiter (type=%v, count=%v)",
		d.Type,
		d.Count,
	)
}

type InlineStringDelimiter struct {
	TypeOfDelimiter InlineDelimiterType
}

func (b InlineStringDelimiter) String() string {
	return "InlineStringDelimiter (" + b.TypeOfDelimiter.String() + ")"
}

func (b *InlineStringDelimiter) GetRawContent() *string {
	return nil
}

// =====================================
// Hrefs

type InlineHref struct {
	Address string
}

func (b InlineHref) String() string {
	return "InlineHref (href=\"" + b.Address + "\")"
}

func (b *InlineHref) GetRawContent() *string {
	return nil
}
