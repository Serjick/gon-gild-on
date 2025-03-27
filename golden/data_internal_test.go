package golden

import (
	"reflect"
	"testing"
)

func TestDataAny_TmplVars(t *testing.T) {
	t.Parallel()

	type fields struct {
		any
	}
	tests := []struct {
		name    string
		fields  fields
		want    any
		wantErr bool
	}{
		{
			name: "Nil",
			fields: fields{
				any: nil,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Struct",
			fields: fields{
				any: struct {
					F bool
				}{F: true},
			},
			want: struct {
				F bool
			}{F: true},
			wantErr: false,
		},
		{
			name: "Map",
			fields: fields{
				any: map[string]any{"f": true},
			},
			want:    map[string]any{"f": true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := DataAny{
				any: tt.fields.any,
			}
			got, err := d.TmplVars()
			if (err != nil) != tt.wantErr {
				t.Errorf("DataAny.TmplVars() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataAny.TmplVars() = %v, want %v", got, tt.want)
			}
		})
	}
}
