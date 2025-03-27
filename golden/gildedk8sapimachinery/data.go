package gildedk8sapimachinery

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/util/jsonmergepatch"

	"github.com/Serjick/gon-gild-on/golden"
)

var _ golden.Data = DataJSONMergePatch{}

// DataJSONMergePatch is implementation of golden.Data based on k8s.io/apimachinery/pkg/util/jsonmergepatch.
type DataJSONMergePatch struct {
	before []byte
	after  []byte
}

func NewDataJSONMergePatch(before, after []byte) DataJSONMergePatch {
	return DataJSONMergePatch{
		before: before,
		after:  after,
	}
}

// TmplVars calc diff and returns it as map[string]any.
func (d DataJSONMergePatch) TmplVars() (any, error) {
	return d.diffDecode()
}

// Format calc diff and pass it into formatter as json.RawMessage.
func (d DataJSONMergePatch) Format(f golden.Formatter) ([]byte, error) {
	p, err := d.diff()
	if err != nil {
		return nil, err
	}

	b, err := f.Bytes(json.RawMessage(p))
	if err != nil {
		return nil, fmt.Errorf("diff format failure: %w", err)
	}

	return b, nil
}

// Valid calc diff and pass it into filter as map[string]any.
func (d DataJSONMergePatch) Valid(f golden.DataFilter) bool {
	p, err := d.diffDecode()

	return err == nil && !f(p)
}

func (d DataJSONMergePatch) diff() ([]byte, error) {
	p, err := jsonmergepatch.CreateThreeWayJSONMergePatch(d.before, d.after, d.before)
	if err != nil {
		return nil, fmt.Errorf("create json merge patch failure: %w", err)
	}

	return p, nil
}

func (d DataJSONMergePatch) diffDecode() (map[string]any, error) {
	p, err := d.diff()
	if err != nil {
		return nil, err
	}

	var v map[string]any
	if err := json.Unmarshal(p, &v); err != nil {
		return nil, fmt.Errorf("unmarshal failure: %w", err)
	}

	return v, nil
}

func (d DataJSONMergePatch) String() string {
	if p, err := d.diff(); err == nil {
		return string(p)
	}

	return ""
}
