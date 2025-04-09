// Package internal of gildedsergigodiff.
package internal //revive:disable-line:comments-density Not a public package

import (
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type SubstitutionPair [2]string

type TextTemplateSubstitutionState struct {
	differ *diffmatchpatch.DiffMatchPatch
	subs   []SubstitutionPair
	from   string
	to     string
	cont   bool
}

func NewTextTemplateSubstitutionState(d *diffmatchpatch.DiffMatchPatch) *TextTemplateSubstitutionState {
	return &TextTemplateSubstitutionState{
		differ: d,
		subs:   nil,
		from:   "",
		to:     "",
		cont:   false,
	}
}

func (p SubstitutionPair) From() string {
	return p[0]
}

func (p SubstitutionPair) To() string {
	return p[1]
}

func (s *TextTemplateSubstitutionState) Update(cur diffmatchpatch.Diff, tail ...diffmatchpatch.Diff) int {
	if !s.cont && s.isActionStart(cur) {
		s.handleActionStart(cur)
	}

	if !s.cont {
		return 0
	}

	if s.isActionClose(cur) {
		return s.handleActionClose(cur, tail...)
	}

	if !s.isActionStart(cur) {
		s.handleActionContinue(cur)
	}

	return 0
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

func (s *TextTemplateSubstitutionState) handleActionClose(cur diffmatchpatch.Diff, tail ...diffmatchpatch.Diff) int {
	if !s.isActionStart(cur) {
		s.to += cur.Text
	}

	var shift int
	if len(tail) > 0 && tail[0].Type == diffmatchpatch.DiffInsert {
		shift++
		s.from += tail[0].Text
	}

	rest := s.differ.DiffText2(tail[shift:])
	s.subs = append(s.subs, SubstitutionPair{s.from + rest, s.to + rest})
	s.from, s.to, s.cont = "", "", false

	return shift
}

func (s *TextTemplateSubstitutionState) Subs() []SubstitutionPair {
	return s.subs
}
