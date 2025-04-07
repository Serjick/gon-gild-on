package gildedsergigodiff

import (
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/Serjick/gon-gild-on/golden"
	"github.com/Serjick/gon-gild-on/golden/gildedsergigodiff/internal"
)

// TextTemplateDiffMatchPatch is a text/template actions transferrer from one string to another.
type TextTemplateDiffMatchPatch struct {
	differ    *diffmatchpatch.DiffMatchPatch
	tmplfuncs template.FuncMap
}

const diffTimeout = time.Minute

func NewTextTemplateDiffMatchPatch(d *diffmatchpatch.DiffMatchPatch, tf template.FuncMap) *TextTemplateDiffMatchPatch {
	return &TextTemplateDiffMatchPatch{
		differ:    d,
		tmplfuncs: tf,
	}
}

// Patch transfers text/template actions from prev to next.
func (p *TextTemplateDiffMatchPatch) Patch(prev, next string) (string, error) {
	patched := next
	for _, pair := range p.calcSubstitions(prev, next) {
		patched = strings.Replace(patched, pair.From(), pair.To(), 1)
	}

	if _, err := template.New("").Funcs(p.tmplfuncs).Parse(patched); err != nil {
		return "", fmt.Errorf("template parse failure: %w", err)
	}

	return patched, nil
}

func (p *TextTemplateDiffMatchPatch) calcSubstitions(prev, next string) []internal.SubstitutionPair {
	state := internal.NewTextTemplateSubstitutionState(p.differ)

	diff := p.differ.DiffBisect(prev, next, time.Now().Add(diffTimeout))
	for i := 0; i < len(diff)-1; i++ {
		i += state.Update(diff[i], diff[i+1:]...)
	}

	return state.Subs()
}

// NewTextTemplateDiffMatchPatchPreSaveHook is a pre save hook to transfer
// text/template actions from present golden file into new.
func NewTextTemplateDiffMatchPatchPreSaveHook() golden.PreSaveHook {
	return func(t golden.TestingT, data []byte, vars golden.PreSaveHookVars) []byte {
		if len(vars.Current) == 0 {
			return data
		}

		patcher := NewTextTemplateDiffMatchPatch(diffmatchpatch.New(), vars.TmplFuncs)

		patched, err := patcher.Patch(string(vars.Current), string(data))
		if err != nil {
			t.Logf("%T patch failed with %q, original content will be used", patcher, err.Error())

			return data
		}

		return []byte(patched)
	}
}
