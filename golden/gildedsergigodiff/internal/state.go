// Package internal of gildedsergigodiff.
package internal //revive:disable-line:comments-density Not a public package

import (
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// SubstitutionPair is a pair of values to substitute first by second.
type SubstitutionPair [2]string

// TextTemplateSubstitutionState is a state of substitutions in text template.
type TextTemplateSubstitutionState struct {
	differ *diffmatchpatch.DiffMatchPatch
	subs   []SubstitutionPair
	from   string
	to     string
	cont   bool
}

// NewTextTemplateSubstitutionState instantiates [TextTemplateSubstitutionState].
func NewTextTemplateSubstitutionState(d *diffmatchpatch.DiffMatchPatch) *TextTemplateSubstitutionState {
	return &TextTemplateSubstitutionState{
		differ: d,
		subs:   nil,
		from:   "",
		to:     "",
		cont:   false,
	}
}

// From is a original value getter.
func (p SubstitutionPair) From() string {
	return p[0]
}

// To is a target value getter.
func (p SubstitutionPair) To() string {
	return p[1]
}

// Update is to add [diffmatchpatch.Diff] to state.
func (s *TextTemplateSubstitutionState) Update(
	cur diffmatchpatch.Diff, tail ...diffmatchpatch.Diff,
) []diffmatchpatch.Diff {
	if !s.cont && s.isActionStart(cur) {
		s.handleActionStart(cur)
	}

	if !s.cont {
		return nil
	}

	if s.isActionClose(cur) {
		return s.handleActionClose(cur, tail...)
	}

	if !s.isActionStart(cur) {
		s.handleActionContinue(cur)
	}

	return nil
}

// Subs is a current state all substitutions getter.
func (s *TextTemplateSubstitutionState) Subs() []SubstitutionPair {
	return s.subs
}

func (*TextTemplateSubstitutionState) isActionStart(cur diffmatchpatch.Diff) bool {
	return cur.Type == diffmatchpatch.DiffDelete &&
		strings.HasPrefix(cur.Text, "{{")
}

func (*TextTemplateSubstitutionState) isActionClose(cur diffmatchpatch.Diff) bool {
	return cur.Type == diffmatchpatch.DiffDelete &&
		strings.Contains(cur.Text, "}}")
}

func (s *TextTemplateSubstitutionState) handleActionStart(cur diffmatchpatch.Diff) {
	s.to, s.cont = cur.Text, true
}

func (s *TextTemplateSubstitutionState) handleActionContinue(cur diffmatchpatch.Diff) {
	switch cur.Type {
	case diffmatchpatch.DiffEqual:
		s.to += cur.Text
		s.from += cur.Text
	case diffmatchpatch.DiffInsert:
		s.from += cur.Text
	case diffmatchpatch.DiffDelete:
		s.to += cur.Text
	}
}

func (s *TextTemplateSubstitutionState) handleActionClose(
	cur diffmatchpatch.Diff, tail ...diffmatchpatch.Diff,
) []diffmatchpatch.Diff {
	if !s.isActionStart(cur) {
		s.to += cur.Text
	}

	rest := tail
	for len(rest) > 0 && rest[0].Type == diffmatchpatch.DiffInsert {
		s.from += rest[0].Text

		rest = rest[1:]
	}

	suffix := s.differ.DiffText2(rest)

	s.subs = append(s.subs, SubstitutionPair{s.from + suffix, s.to + suffix})
	s.from, s.to, s.cont = "", "", false

	return rest
}
