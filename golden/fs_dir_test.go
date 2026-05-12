package golden_test

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"testing/fstest"

	"github.com/Serjick/gon-gild-on/golden"
)

func ExampleFS_RenderDir() {
	got := fstest.MapFS{
		"file.txt": &fstest.MapFile{
			Data: []byte(`plaintext`),
		},
		"subdir/file.json": &fstest.MapFile{
			Data: []byte(`{"key":"value"}`),
		},
	}

	fsys := golden.NewFS(
		golden.WithFSFormatter(golden.NewStrFormatter()),
		golden.WithFSLocator(func(golden.LocationVars) fmt.Stringer {
			return bytes.NewBufferString(filepath.Join("testdata", "golden", "dir-example", "example.tmpl"))
		}),
		golden.WithFSDirEntryOverrides(
			filepath.Join("subdir", "file.json"), golden.WithFSFormatter(golden.NewJSONFormatter()),
		),
	)
	want, err := fsys.RenderDir(new(testing.T), got)
	_ = fs.WalkDir(want, ".", func(path string, d fs.DirEntry, _ error) error {
		if !d.IsDir() {
			f, er := fs.ReadFile(want, path)
			fmt.Printf("%s: %s, %v\n", path, f, er)
		}

		return nil
	})
	fmt.Println(err)
	// Output:
	// file.txt: plaintext, <nil>
	// subdir/file.json: {
	//     "key": "value"
	// }
	// , <nil>
	// <nil>
}

func TestFS_RenderDir(t *testing.T) {
	t.Parallel()

	type fields struct {
		src  fs.FS
		opts []golden.FSOption
	}
	type args struct {
		actual fs.FS
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    fs.FS
		wantSrc fs.FS
		wantErr bool
	}{
		{
			name: "TwoJsonFilesNoSubdirs",
			fields: fields{
				src: os.DirFS(t.TempDir()),
			},
			args: args{
				actual: fstest.MapFS{
					"lower.json": &fstest.MapFile{
						Data: []byte(`{"key":"value"}`),
						Mode: golden.DefaultFilePerm,
					},
					"UPPER.json": &fstest.MapFile{
						Data: []byte(`{"KEY":"VALUE"}`),
						Mode: golden.DefaultFilePerm,
					},
				},
			},
			want: fstest.MapFS{
				"lower.json": &fstest.MapFile{
					Data: []byte("{\n    \"key\": \"value\"\n}\n"),
					Mode: golden.DefaultFilePerm,
				},
				"UPPER.json": &fstest.MapFile{
					Data: []byte("{\n    \"KEY\": \"VALUE\"\n}\n"),
					Mode: golden.DefaultFilePerm,
				},
			},
			wantErr: false,
		},
		{
			name: "JsonFileInSubdir",
			fields: fields{
				src: os.DirFS(t.TempDir()),
			},
			args: args{
				actual: fstest.MapFS{
					filepath.Join("subdir", "file.json"): &fstest.MapFile{
						Data: []byte(`{"key":"value"}`),
						Mode: golden.DefaultFilePerm,
					},
				},
			},
			want: fstest.MapFS{
				filepath.Join("subdir", "file.json"): &fstest.MapFile{
					Data: []byte("{\n    \"key\": \"value\"\n}\n"),
					Mode: golden.DefaultFilePerm,
				},
			},
			wantErr: false,
		},
		{
			name: "SubdirUpdated",
			fields: fields{
				src: func() fs.FS {
					dir := t.TempDir()
					path := filepath.Join(dir, "testdata", "golden", "fs", t.Name(), "SubdirUpdated", "subdir")
					if err := os.MkdirAll(path, golden.DefaultDirPerm); err != nil {
						t.Fatalf("%q fixture write failed: %s", path, err)
					}
					file := filepath.Join(path, "file.json")
					if err := os.WriteFile(file, []byte("{}\n"), golden.DefaultFilePerm); err != nil {
						t.Fatalf("%q fixture write failed: %s", file, err)
					}

					if err := os.WriteFile(
						filepath.Join(path, "REDUNDANT.json"), []byte(`"string"`), golden.DefaultFilePerm,
					); err != nil {
						t.Fatalf("%q fixture write failed: %s", "REDUNDANT.json", err)
					}

					return os.DirFS(dir)
				}(),
				opts: []golden.FSOption{
					golden.WithFSLocator(golden.NewLocatorSubDir("fs")),
					golden.WithFSForceUpdate(),
				},
			},
			args: args{
				actual: fstest.MapFS{
					filepath.Join("subdir", "file.json"): &fstest.MapFile{
						Data: []byte(`{"key":"value"}`),
						Mode: golden.DefaultFilePerm,
					},
				},
			},
			want: fstest.MapFS{
				filepath.Join("subdir", "file.json"): &fstest.MapFile{
					Data: []byte("{\n    \"key\": \"value\"\n}\n"),
					Mode: golden.DefaultFilePerm,
				},
			},
			wantSrc: fstest.MapFS{
				filepath.Join(
					"testdata", "golden", "fs", t.Name(), "SubdirUpdated", "subdir", "file.json",
				): &fstest.MapFile{
					Data: []byte("{\n    \"key\": \"value\"\n}\n"),
					Mode: golden.DefaultFilePerm,
				},
			},
			wantErr: false,
		},
		{
			name: "WithOverride",
			fields: fields{
				src: os.DirFS(t.TempDir()),
				opts: []golden.FSOption{
					golden.WithFSDirEntryOverrides("subdir", golden.WithFSFormatter(golden.NewStrFormatter())),
				},
			},
			args: args{
				actual: fstest.MapFS{
					"file.json": &fstest.MapFile{
						Data: []byte(`{"key":"value"}`),
						Mode: golden.DefaultFilePerm,
					},
					filepath.Join("subdir", "file.txt"): &fstest.MapFile{
						Data: []byte("plaintext"),
						Mode: golden.DefaultFilePerm,
					},
				},
			},
			want: fstest.MapFS{
				"file.json": &fstest.MapFile{
					Data: []byte("{\n    \"key\": \"value\"\n}\n"),
					Mode: golden.DefaultFilePerm,
				},
				filepath.Join("subdir", "file.txt"): &fstest.MapFile{
					Data: []byte("plaintext"),
					Mode: golden.DefaultFilePerm,
				},
			},
			wantErr: false,
		},
		{
			name: "EmptySubdir",
			fields: fields{
				src: func() fs.FS {
					dir := t.TempDir()
					path := filepath.Join(dir, "testdata", "golden", t.Name(), "EmptySubdir", "dir", "subdir")
					if err := os.MkdirAll(path, golden.DefaultDirPerm); err != nil {
						t.Fatalf("%q fixture write failed: %s", path, err)
					}
					file := filepath.Join(path, golden.DirKeepMarkerFilename)
					if err := os.WriteFile(file, nil, golden.DefaultFilePerm); err != nil {
						t.Fatalf("%q fixture write failed: %s", file, err)
					}

					return os.DirFS(dir)
				}(),
			},
			args: args{
				actual: fstest.MapFS{
					filepath.Join("dir", "subdir"): &fstest.MapFile{
						Mode: golden.DefaultDirPerm,
					},
				},
			},
			want: fstest.MapFS{
				filepath.Join("dir", "subdir"): &fstest.MapFile{
					Mode: golden.DefaultDirPerm,
				},
			},
			wantErr: false,
		},
		{
			name: "PartialUpdate",
			fields: fields{
				src: func() fs.FS {
					dir := t.TempDir()
					path := filepath.Join(dir, "testdata", "golden", t.Name(), "PartialUpdate")
					if err := os.MkdirAll(path, golden.DefaultDirPerm); err != nil {
						t.Fatalf("%q fixture write failed: %s", path, err)
					}
					file := filepath.Join(path, "updatable") // NOTE: this is file
					if err := os.WriteFile(file, nil, golden.DefaultFilePerm); err != nil {
						t.Fatalf("%q fixture write failed: %s", file, err)
					}
					roFile := filepath.Join(path, "readonly") // NOTE: this is file
					if err := os.WriteFile(roFile, []byte{'\n'}, golden.DefaultFilePerm); err != nil {
						t.Fatalf("%q fixture write failed: %s", roFile, err)
					}

					return os.DirFS(dir)
				}(),
				opts: []golden.FSOption{
					golden.WithFSDirEntryOverrides("updatable", golden.WithFSForceUpdate()),
				},
			},
			args: args{
				actual: fstest.MapFS{
					"updatable": &fstest.MapFile{
						Mode: golden.DefaultDirPerm, // NOTE: this is NOT file
					},
					"readonly": &fstest.MapFile{
						Mode: golden.DefaultDirPerm, // NOTE: this is NOT file
					},
				},
			},
			want: fstest.MapFS{
				"updatable": &fstest.MapFile{
					Mode: golden.DefaultDirPerm,
				},
				"readonly": &fstest.MapFile{
					Data: []byte{'\n'},
					Mode: golden.DefaultFilePerm,
				},
			},
			wantErr: false,
		},
		{
			name: "CleanUp",
			fields: fields{
				src: func() fs.FS {
					dir := t.TempDir()

					dir1 := filepath.Join(dir, "testdata", "golden", "dir1", t.Name(), "CleanUp")
					subdir1 := filepath.Join(dir1, "subdir1")
					if err := os.MkdirAll(subdir1, golden.DefaultDirPerm); err != nil {
						t.Fatalf("%q fixture write failed: %s", subdir1, err)
					}
					if err := os.WriteFile(filepath.Join(dir1, "file"), nil, golden.DefaultFilePerm); err != nil {
						t.Fatalf("%q file fixture write failed: %s", dir1, err)
					}
					if err := os.WriteFile(filepath.Join(subdir1, "file"), nil, golden.DefaultFilePerm); err != nil {
						t.Fatalf("%q file fixture write failed: %s", subdir1, err)
					}
					subdir1file := filepath.Join(subdir1, golden.DirKeepMarkerFilename)
					if err := os.WriteFile(subdir1file, nil, golden.DefaultFilePerm); err != nil {
						t.Fatalf("%q fixture write failed: %s", subdir1file, err)
					}

					dir2 := filepath.Join(dir, "testdata", "golden", "dir2", t.Name(), "CleanUp")
					subdir2 := filepath.Join(dir2, "subdir2")
					if err := os.MkdirAll(subdir2, golden.DefaultDirPerm); err != nil {
						t.Fatalf("%q fixture write failed: %s", dir1, err)
					}
					subdir2file := filepath.Join(subdir2, golden.DirKeepMarkerFilename)
					if err := os.WriteFile(subdir2file, nil, golden.DefaultFilePerm); err != nil {
						t.Fatalf("%q fixture write failed: %s", subdir2file, err)
					}

					return os.DirFS(dir)
				}(),
				opts: []golden.FSOption{
					golden.WithFSForceUpdate(),
					golden.WithFSLocator(golden.NewLocatorSubDir("dir1")),
					golden.WithFSDirEntryOverrides("subdir2", golden.WithFSLocator(golden.NewLocatorSubDir("dir2"))),
				},
			},
			args: args{
				actual: fstest.MapFS{
					"subdir1": &fstest.MapFile{
						Mode: golden.DefaultDirPerm,
					},
				},
			},
			want: fstest.MapFS{
				"subdir1": &fstest.MapFile{
					Mode: golden.DefaultDirPerm,
				},
			},
			wantSrc: fstest.MapFS{
				filepath.Join("testdata", "golden", "dir1", t.Name(), "CleanUp", "subdir1"): &fstest.MapFile{
					Mode: golden.DefaultDirPerm,
				},
				filepath.Join(
					"testdata", "golden", "dir1", t.Name(), "CleanUp", "subdir1", golden.DirKeepMarkerFilename,
				): &fstest.MapFile{
					Mode: golden.DefaultFilePerm,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opts := append([]golden.FSOption{golden.WithFSSource(golden.NewSourceFS(tt.fields.src))}, tt.fields.opts...)
			f := golden.NewFS(opts...)
			got, err := f.RenderDir(t, tt.args.actual)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FS.RenderDir() error = %v, wantErr %v", err, tt.wantErr)
			}

			if assertLeftFSIsRightFSSubset(t, got, tt.want) {
				assertLeftFSIsRightFSSubset(t, tt.want, got)
			}

			if tt.wantSrc == nil {
				return
			}

			if assertLeftFSIsRightFSSubset(t, tt.fields.src, tt.wantSrc) {
				assertLeftFSIsRightFSSubset(t, tt.wantSrc, tt.fields.src)
			}
		})
	}
}

func assertLeftFSIsRightFSSubset(t *testing.T, left, right fs.FS) bool {
	t.Helper()

	err := fs.WalkDir(left, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rightF, err := right.Open(path)
		if err != nil {
			return fmt.Errorf("open right fs %q entry failure: %w", path, err)
		}
		defer rightF.Close()

		if entry.IsDir() {
			i, er := rightF.Stat()
			if er != nil {
				return fmt.Errorf("stat right fs %q dir failure: %w", path, er)
			}
			if !i.IsDir() {
				return fmt.Errorf("right fs %q is not dir", path)
			}

			return nil
		}

		rightB, err := io.ReadAll(rightF)
		if err != nil {
			return fmt.Errorf("read right fs %q file failure: %w", path, err)
		}

		leftB, err := fs.ReadFile(left, path)
		if err != nil {
			return fmt.Errorf("read left fs %q file failure: %w", path, err)
		}

		if !reflect.DeepEqual(leftB, rightB) {
			return fmt.Errorf("%q file data mismatch\ngot: %s\nwant: %s", path, leftB, rightB)
		}

		return nil
	})
	if err != nil {
		t.Error(err)

		return false
	}

	return true
}
