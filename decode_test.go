package jsondsl

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDecodeEmptyArray(t *testing.T) {
	input := `[]`

	d := &Decoder{}
	d.Reset(strings.NewReader(input))
	got, err := d.Decode()

	wantErr := false
	want := []any(nil)

	gotErr := err != nil
	if gotErr != wantErr {
		t.Fatalf("TestParse(): got err = %v, want err = %v", err, wantErr)
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TestParse(): got diff:\n%s", diff)
	}
}

func TestDecodeEmptyObject(t *testing.T) {
	input := `{}`

	d := &Decoder{}
	d.Reset(strings.NewReader(input))
	got, err := d.Decode()

	wantErr := false
	want := map[string]any(nil)

	gotErr := err != nil
	if gotErr != wantErr {
		t.Fatalf("TestParse(): got err = %v, want err = %v", err, wantErr)
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TestParse(): got diff:\n%s", diff)
	}
}

func TestDecodeOperator(t *testing.T) {
	input := `op()`

	d := &Decoder{}
	d.Reset(strings.NewReader(input))
	got, err := d.Decode()

	wantErr := false
	want := Op{Id: Id("op")}

	gotErr := err != nil
	if gotErr != wantErr {
		t.Fatalf("TestParse(): got err = %v, want err = %v", err, wantErr)
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TestParse(): got diff:\n%s", diff)
	}
}

func TestDecodeEverythingArray(t *testing.T) {
	input := `[
		null,
		false,
		true,
		0,
		1.0,
		-1e+75,
		"\"abc\"",
		[],
		{},
		id,
		add(1,2),
	]`

	d := &Decoder{}
	d.Reset(strings.NewReader(input))
	got, err := d.Decode()

	wantErr := false
	want := []any{
		nil,
		false,
		true,
		float64(0),
		1.0,
		-1e+75,
		`"abc"`,
		[]any(nil),
		map[string]any(nil),
		Id("id"),
		Op{
			Id: Id("add"),
			Args: []any{
				float64(1),
				float64(2),
			},
		},
	}

	gotErr := err != nil
	if gotErr != wantErr {
		t.Fatalf("TestParse(): got err = %v, want err = %v", err, wantErr)
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TestParse(): got diff:\n%s", diff)
	}
}
