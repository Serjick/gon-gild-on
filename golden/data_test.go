package golden_test

import (
	"reflect"
	"testing"

	"github.com/Serjick/gon-gild-on/golden"
)

func TestDataJSON_TmplVars(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		d       golden.DataJSON
		want    any
		wantErr bool
	}{
		{
			name:    "Nil",
			d:       golden.DataJSON("null"),
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Incorrect",
			d:       golden.DataJSON("{"),
			wantErr: true,
		},
		{
			name:    "Map",
			d:       golden.DataJSON(`{"f": true}`),
			want:    map[string]any{"f": true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.d.TmplVars()
			if (err != nil) != tt.wantErr {
				t.Errorf("DataJSON.TmplVars() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataJSON.TmplVars() = %v, want %v", got, tt.want)
			}
		})
	}
}
