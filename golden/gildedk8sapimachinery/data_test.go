package gildedk8sapimachinery_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Serjick/gon-gild-on/golden"
	"github.com/Serjick/gon-gild-on/golden/gildedk8sapimachinery"
)

func ExampleDataJSONMergePatch_Format() {
	d := gildedk8sapimachinery.NewDataJSONMergePatch(
		[]byte(`{"f1":{},"f2":1,"f9":true}`),
		[]byte(`{"f2":2,"f3":{},"f9":true}`),
	)
	b, err := d.Format(golden.NewJSONFormatter())
	fmt.Printf("%s", b)
	fmt.Print(err)
	// Output:
	// {
	//     "f1": null,
	//     "f2": 2,
	//     "f3": {}
	// }
	// <nil>
}

func ExampleDataJSONMergePatch_String() {
	d := gildedk8sapimachinery.NewDataJSONMergePatch(
		[]byte(`{"f1":{},"f2":1,"f9":true}`),
		[]byte(`{"f2":2,"f3":{},"f9":true}`),
	)
	fmt.Printf("%s", d)
	// Output:
	// {"f1":null,"f2":2,"f3":{}}
}

func TestDataJSONMergePatch_TmplVars(t *testing.T) {
	t.Parallel()

	type fields struct {
		before []byte
		after  []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    any
		wantErr bool
	}{
		{
			name: "Equal",
			fields: fields{
				before: []byte(`{"f":1}`),
				after:  []byte(`{"f":1}`),
			},
			want:    map[string]any{},
			wantErr: false,
		},
		{
			name: "AddField",
			fields: fields{
				before: []byte(`{}`),
				after:  []byte(`{"f":1}`),
			},
			want:    map[string]any{"f": float64(1)},
			wantErr: false,
		},
		{
			name: "ChangeField",
			fields: fields{
				before: []byte(`{"f":1}`),
				after:  []byte(`{"f":2}`),
			},
			want:    map[string]any{"f": float64(2)},
			wantErr: false,
		},
		{
			name: "RemoveField",
			fields: fields{
				before: []byte(`{"f":1}`),
				after:  []byte(`{}`),
			},
			want:    map[string]any{"f": nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := gildedk8sapimachinery.NewDataJSONMergePatch(tt.fields.before, tt.fields.after)
			got, err := d.TmplVars()
			if (err != nil) != tt.wantErr {
				t.Errorf("DataJSONMergePatch.TmplVars() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataJSONMergePatch.TmplVars() = %v, want %v", got, tt.want)
			}
		})
	}
}
