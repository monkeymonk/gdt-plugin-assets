package report

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/monkeymonk/gdt-assets/internal/asset"
	"github.com/monkeymonk/gdt-assets/internal/diagnostic"
)

func FormatInventory(w io.Writer, assets []asset.Asset, format string) {
	switch format {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		enc.Encode(assets)
	case "csv":
		cw := csv.NewWriter(w)
		cw.Write([]string{"path", "type", "size"})
		for _, a := range assets {
			cw.Write([]string{a.Path, a.Type.String(), fmt.Sprintf("%d", a.Size)})
		}
		cw.Flush()
	case "md":
		fmt.Fprintf(w, "| Path | Type | Size |\n|------|------|------|\n")
		for _, a := range assets {
			fmt.Fprintf(w, "| %s | %s | %s |\n", a.Path, a.Type, asset.HumanSize(a.Size))
		}
	default: // table
		for _, a := range assets {
			fmt.Fprintf(w, "  %-8s %-50s %s\n", a.Type, a.Path, asset.HumanSize(a.Size))
		}
	}
}

func FormatDiagnostics(w io.Writer, diags *diagnostic.Set, format string) {
	switch format {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		enc.Encode(diags.Items)
	case "csv":
		cw := csv.NewWriter(w)
		cw.Write([]string{"severity", "path", "rule", "message"})
		for _, d := range diags.Items {
			cw.Write([]string{d.Severity.String(), d.Path, d.Rule, d.Message})
		}
		cw.Flush()
	default: // table or md
		for _, d := range diags.Items {
			fmt.Fprintf(w, "  %-7s  %s: %s\n", d.Severity, d.Path, d.Message)
		}
	}
}

func FormatSummary(w io.Writer, assets []asset.Asset) {
	byType := make(map[asset.AssetType]int)
	var totalSize int64
	for _, a := range assets {
		byType[a.Type]++
		totalSize += a.Size
	}

	types := make([]string, 0, len(byType))
	for t, n := range byType {
		types = append(types, fmt.Sprintf("%s: %d", t, n))
	}
	sort.Strings(types)

	fmt.Fprintf(w, "Total: %d assets (%s)\n", len(assets), asset.HumanSize(totalSize))
	fmt.Fprintf(w, "  %s\n", strings.Join(types, ", "))
}
