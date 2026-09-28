package mutate

import (
	"bytes"
	"testing"
)

func TestApplyExactBytes(t *testing.T) {
	src := []byte("if a == b {\n\treturn true\n}\n")
	// "==" is at offset 5..7
	site := Site{Start: 5, End: 7, Replacement: "!=", Desc: "line 1: flip =="}
	out, err := Apply(src, site)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	want := []byte("if a != b {\n\treturn true\n}\n")
	if !bytes.Equal(out, want) {
		t.Fatalf("Apply = %q, want %q", out, want)
	}
	// Everything before Start and after End must be byte-identical.
	if !bytes.Equal(out[:site.Start], src[:site.Start]) {
		t.Fatalf("prefix changed: got %q want %q", out[:site.Start], src[:site.Start])
	}
	suffix := src[site.End:]
	gotSuffix := out[site.Start+len(site.Replacement):]
	if !bytes.Equal(gotSuffix, suffix) {
		t.Fatalf("suffix changed: got %q want %q", gotSuffix, suffix)
	}
}

func TestApplyOutOfRange(t *testing.T) {
	src := []byte("abc")
	cases := []Site{
		{Start: -1, End: 2},
		{Start: 2, End: 1},
		{Start: 0, End: 10},
		{Start: 10, End: 10},
	}
	for _, s := range cases {
		if _, err := Apply(src, s); err == nil {
			t.Errorf("Apply(%+v) = nil error, want error", s)
		}
	}
}

func TestSampleDeterministicAndCapped(t *testing.T) {
	sites := []Site{
		{Start: 50, Desc: "e"},
		{Start: 10, Desc: "a"},
		{Start: 30, Desc: "c"},
		{Start: 20, Desc: "b"},
		{Start: 40, Desc: "d"},
	}
	seed := []byte("some-module-bytes")

	// Sampling with max >= len(sites) returns everything, sorted by Start.
	all := Sample(sites, 10, seed)
	if len(all) != len(sites) {
		t.Fatalf("len(all) = %d, want %d", len(all), len(sites))
	}
	for i := 1; i < len(all); i++ {
		if all[i-1].Start > all[i].Start {
			t.Fatalf("Sample not sorted by Start: %+v", all)
		}
	}

	// Capped sampling: deterministic given the same seed, at most max
	// sites, sorted by Start.
	got1 := Sample(sites, 2, seed)
	got2 := Sample(sites, 2, seed)
	if len(got1) != 2 {
		t.Fatalf("len(got1) = %d, want 2", len(got1))
	}
	if !equalSites(got1, got2) {
		t.Fatalf("Sample not deterministic: %+v vs %+v", got1, got2)
	}
	for i := 1; i < len(got1); i++ {
		if got1[i-1].Start > got1[i].Start {
			t.Fatalf("capped Sample not sorted by Start: %+v", got1)
		}
	}

	// A different seed may pick a different subset (not asserted which,
	// but the function must not panic and must still respect max/sort).
	got3 := Sample(sites, 2, []byte("different-seed"))
	if len(got3) != 2 {
		t.Fatalf("len(got3) = %d, want 2", len(got3))
	}
}

func equalSites(a, b []Site) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestForKnownAndUnknownLangs(t *testing.T) {
	if m := For("go"); m == nil || m.Lang() != "go" {
		t.Fatalf("For(\"go\") = %+v, want a Go Mutator", m)
	}
	if m := For("nonsense"); m != nil {
		t.Fatalf("For(\"nonsense\") = %+v, want nil", m)
	}
}
