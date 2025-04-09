package gildedsergigodiff_test

import (
	"fmt"
	"math"
	"testing"
	"text/template"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/Serjick/gon-gild-on/golden/gildedsergigodiff"
)

func TestTextTemplateDiffMatchPatch_Patch(t *testing.T) {
	t.Parallel()

	type fields struct {
		differ *diffmatchpatch.DiffMatchPatch
		tf     template.FuncMap
	}
	type args struct {
		prev string
		next string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "FieldLodge",
			fields: fields{
				differ: diffmatchpatch.New(),
			},
			args: args{
				prev: FieldLodgePrev,
				next: FieldLodgeNext,
			},
			want: FieldLodgeWant,
		},
		{
			name: "FuncMapKey",
			fields: fields{
				differ: diffmatchpatch.New(),
			},
			args: args{
				prev: FuncMapKeyPrev,
				next: FuncMapKeyNext,
			},
			want: FuncMapKeyWant,
		},
		{
			name: "TimeInMap",
			fields: fields{
				differ: diffmatchpatch.New(),
			},
			args: args{
				prev: TimeInMapPrev,
				next: TimeInMapNext,
			},
			want: TimeInMapWant,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := gildedsergigodiff.NewTextTemplateDiffMatchPatch(tt.fields.differ, tt.fields.tf)
			got, err := p.Patch(tt.args.prev, tt.args.next)
			if (err != nil) != tt.wantErr {
				t.Errorf("TextTemplateDiffMatchPatch.Patch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TextTemplateDiffMatchPatch.Patch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func FuzzTextTemplateDiffMatchPatch_Patch_time(f *testing.F) {
	corpus := []struct {
		sec  uint32
		nsec uint32
	}{
		{
			sec:  0,
			nsec: 987654321,
		},
		{
			sec:  1743419509,
			nsec: 123456789,
		},
		{
			sec:  math.MaxUint32,
			nsec: 0,
		},
	}
	for _, c := range corpus {
		f.Add(c.sec, c.nsec)
	}
	f.Fuzz(func(t *testing.T, sec, nsec uint32) {
		ts := time.Unix(int64(sec), int64(nsec))

		differ := diffmatchpatch.New()

		p := gildedsergigodiff.NewTextTemplateDiffMatchPatch(differ, nil)
		got, err := p.Patch(TimeFuzzPrev, fmt.Sprintf(TimeFuzzNext, ts))
		if err != nil {
			t.Errorf("TextTemplateDiffMatchPatch.Patch() error = %v", err)
			return
		}
		if got != TimeFuzzWant {
			t.Errorf("TextTemplateDiffMatchPatch.Patch() = %v, want %v", got, TimeFuzzWant)
		}
	})
}
