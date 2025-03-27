package golden_test

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Serjick/gon-gild-on/golden"
)

//go:embed *
var goldens embed.FS

func ExampleFS_RenderFile_embedfs() {
	/* -- testdata/golden/example.tmpl --
	{
	    "key": "{{ .Actual.key }}"
	}
	*/
	path := filepath.Join("testdata", "golden", "example.tmpl")
	f := golden.NewFS(goldens, golden.WithFSLocator(func(golden.LocationVars) string {
		return path
	}))
	b, err := f.RenderFile(new(testing.T), map[string]string{"key": "value"})
	fmt.Printf("%s", b)
	fmt.Print(err)
	// Output:
	// {
	//     "key": "value"
	// }
	// <nil>
}

func TestFS_RenderFile(t *testing.T) {
	t.Parallel()

	type fields struct {
		src  fs.FS
		opts []golden.FSOption
	}
	type args struct {
		actual any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantF   string
		wantErr bool
	}{
		{
			name: "AutoCreate",
			fields: fields{
				src: os.DirFS(t.TempDir()),
			},
			args: args{
				actual: true,
			},
			want:    []byte("true\n"),
			wantF:   filepath.Join("testdata", "golden", "TestFS_RenderFile", "AutoCreate", "golden.tmpl"),
			wantErr: false,
		},
		{
			name: "Filtered",
			fields: fields{
				src: os.DirFS("/"),
				opts: []golden.FSOption{
					golden.WithFSDataFilter(func(any) bool {
						return true
					}),
					golden.WithFSWriter(nil, nil),
				},
			},
			args: args{
				actual: true,
			},
			want:    []byte("true\n"),
			wantErr: false,
		},
		{
			name: "RawJSON",
			fields: fields{
				src: os.DirFS(t.TempDir()),
			},
			args: args{
				actual: json.RawMessage(`{"key":"value"}`),
			},
			want:    []byte("{\n    \"key\": \"value\"\n}\n"),
			wantErr: false,
		},
		{
			name: "Updated",
			fields: fields{
				src: func() fs.FS {
					dir := t.TempDir()
					path := filepath.Join(dir, "testdata", "golden", "TestFS_RenderFile", "Updated")
					if err := os.MkdirAll(path, golden.DefaultDirPerm); err != nil {
						t.Errorf("%q fixture write failed: %s", path, err)
					}
					file := filepath.Join(path, "golden.tmpl")
					if err := os.WriteFile(file, []byte("1\n"), golden.DefaultFilePerm); err != nil {
						t.Errorf("%q fixture write failed: %s", file, err)
					}

					return os.DirFS(dir)
				}(),
				opts: []golden.FSOption{golden.WithFSForceUpdate()},
			},
			args: args{
				actual: 2,
			},
			want:    []byte("2\n"),
			wantF:   filepath.Join("testdata", "golden", "TestFS_RenderFile", "Updated", "golden.tmpl"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opts := tt.fields.opts
			if fmt.Sprintf("%T", tt.fields.src) == "os.dirFS" {
				opts = append([]golden.FSOption{golden.WithFSRoot(fmt.Sprintf("%s", tt.fields.src))}, opts...)
			}
			f := golden.NewFS(tt.fields.src, opts...)
			got, err := f.RenderFile(t, tt.args.actual)
			if (err != nil) != tt.wantErr {
				t.Errorf("FS.RenderFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FS.RenderFile() = %s, want %s", got, tt.want)
			}
			if tt.wantF == "" {
				return
			}

			path := filepath.Join(fmt.Sprintf("%s", tt.fields.src), tt.wantF)
			gotF, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("%q read failed: %s", path, err.Error())
				return
			}
			if !reflect.DeepEqual(gotF, tt.want) {
				t.Errorf("FS.writeFile() = %s, want %s", gotF, tt.want)
			}
		})
	}
}
