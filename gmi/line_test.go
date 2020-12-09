package gmi_test

import (
	"testing"

	"git.sr.ht/~kiba/gdn/gmi"
)

func TestLineTypeString(t *testing.T) {
	if gmi.LineType(0).String() != "UNKNOWN" {
		t.Errorf("Expected `UNKNOWN` for line type 0, got: `%s`",
			gmi.LineType(0))
	}

	if gmi.LineType(-1).String() != "UNKNOWN" {
		t.Errorf("Expected `UNKNOWN` for line type -1, got: `%s`",
			gmi.LineType(-1))
	}

	if gmi.Head1.String() != "Head1" {
		t.Errorf("Expected `Head1` for line type, got: `%s`", gmi.Head1)
	}

	if gmi.Head2.String() != "Head2" {
		t.Errorf("Expected `Head2` for line type, got: `%s`", gmi.Head2)
	}

	if gmi.Head3.String() != "Head3" {
		t.Errorf("Expected `Head3` for line type, got: `%s`", gmi.Head3)
	}

	if gmi.Text.String() != "Text" {
		t.Errorf("Expected `Text` for line type, got: `%s`", gmi.Text)
	}

	if gmi.Link.String() != "Link" {
		t.Errorf("Expected `Link` for line type, got: `%s`", gmi.Link)
	}

	if gmi.PreStart.String() != "PreStart" {
		t.Errorf("Expected `PreStart` for line type, got: `%s`", gmi.PreStart)
	}

	if gmi.PreBody.String() != "PreBody" {
		t.Errorf("Expected `Head1` for line type, got: `%s`", gmi.PreBody)
	}

	if gmi.PreEnd.String() != "PreEnd" {
		t.Errorf("Expected `Head1` for line type, got: `%s`", gmi.PreEnd)
	}

	if gmi.List.String() != "List" {
		t.Errorf("Expected `List` for line type, got: `%s`", gmi.List)
	}

	if gmi.Quote.String() != "Quote" {
		t.Errorf("Expected `Quote` for line type, got: `%s`", gmi.Quote)
	}
}
