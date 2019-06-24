package protoc

import (
	"reflect"
	"testing"
)

func TestExtractImports(t *testing.T) {
	tests := []struct {
		content string
		results []string
	}{
		{
			`
import "a.proto";
 import     "b.proto"   ; // some comments
// import "c.proto";
			`,
			[]string{
				"a.proto",
				"b.proto",
			},
		},
		{
			`
import "a.proto";
/**
 import     "b.proto"   ; // some comments
*/
// import "c.proto";
			`,
			[]string{
				"a.proto",
				"b.proto",
			},
		},
	}

	for i, test := range tests {
		got := extractImports(test.content)
		if !reflect.DeepEqual(got, test.results) {
			t.Errorf("extractImports(CASE %d) returns %v, but want %v", i, got, test.results)
		}
	}
}

func TestParse(t *testing.T) {
	p := New()
	tests := []struct {
		files       []string
		wantFileNum int
	}{
		{
			[]string{"testdata/test1.proto"},
			1,
		},
		{
			[]string{"testdata/test2.proto"},
			1,
		},
		{
			[]string{"testdata/test1.proto", "testdata/test2.proto"},
			2,
		},
	}
	for i, test := range tests {
		f, err := p.Parse(test.files...)
		if err != nil {
			t.Errorf("Parse(CASE #%d) returns error: %v", i, err)
		}
		if got, want := len(f.GetFile()), test.wantFileNum; got != want {
			t.Errorf("Parse(CASE #%d) returns %d files, but want %d", i, got, want)
		}
	}
}

func TestParse_IncludeImports(t *testing.T) {
	p := New(
		IncludeImports(),
	)
	tests := []struct {
		files       []string
		wantFileNum int
	}{
		{
			[]string{"testdata/test2.proto"},
			2,
		},
	}
	for i, test := range tests {
		f, err := p.Parse(test.files...)
		if err != nil {
			t.Errorf("Parse(CASE #%d) returns error: %v", i, err)
		}
		if got, want := len(f.GetFile()), test.wantFileNum; got != want {
			t.Errorf("Parse(CASE #%d) returns %d files, but want %d", i, got, want)
		}
	}
}

func TestParse_ExcludeCommonProtos(t *testing.T) {
	p := New(ExcludeCommonProtos())
	tests := []struct {
		files  []string
		hasErr bool
	}{
		{
			[]string{"testdata/test1.proto"},
			false,
		},
		{
			[]string{"testdata/test2.proto"},
			true,
		},
	}
	for i, test := range tests {
		_, err := p.Parse(test.files...)
		if (err != nil) != test.hasErr {
			t.Errorf("Parse(CASE #%d) returns error: %v, but want having error %v", i, err, test.hasErr)
		}
	}
}
