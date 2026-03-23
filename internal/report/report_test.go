package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/diagnostic"
)

func TestTableFormat(t *testing.T) {
	var buf bytes.Buffer
	assets := []asset.Asset{
		{Path: "hero.png", Type: asset.TypeImage, Size: 1024},
	}
	FormatInventory(&buf, assets, "table")
	out := buf.String()
	if !strings.Contains(out, "hero.png") {
		t.Error("table output missing asset path")
	}
}

func TestDiagnosticsJSON(t *testing.T) {
	var buf bytes.Buffer
	diags := &diagnostic.Set{}
	diags.Add(diagnostic.Diagnostic{
		Path: "hero.png", Severity: diagnostic.Warning,
		Rule: "test.rule", Message: "test msg",
	})
	FormatDiagnostics(&buf, diags, "json")
	out := buf.String()
	if !strings.Contains(out, `"rule"`) {
		t.Error("json output missing rule field")
	}
}
