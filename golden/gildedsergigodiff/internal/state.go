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
	if s.isActionStart(cur) {
		s.handleActionStart(cur)

		return 0
	}

	if !s.cont {
		return 0
	}

	next := tail[0]
	if s.isActionClose(cur, next) {
		s.handleActionClose(cur, next, tail[1:]...)

		return 1
	}

	s.handleActionContinue(cur)

	return 0
}

func (s *TextTemplateSubstitutionState) isActionStart(cur diffmatchpatch.Diff) bool {
	return !s.cont && cur.Type == diffmatchpatch.DiffDelete &&
		strings.HasPrefix(cur.Text, "{{")
}

func (s *TextTemplateSubstitutionState) isActionClose(cur, next diffmatchpatch.Diff) bool {
	return s.cont && cur.Type == diffmatchpatch.DiffDelete &&
		next.Type == diffmatchpatch.DiffInsert &&
		strings.HasSuffix(cur.Text, "}}")
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

func (s *TextTemplateSubstitutionState) handleActionClose(cur, next diffmatchpatch.Diff, tail ...diffmatchpatch.Diff) {
	s.to += cur.Text
	s.from += next.Text

	rest := s.differ.DiffText2(tail)
	s.subs = append(s.subs, SubstitutionPair{s.from + rest, s.to + rest})
	s.from, s.to, s.cont = "", "", false
}

func (s *TextTemplateSubstitutionState) Subs() []SubstitutionPair {
	return s.subs
}
