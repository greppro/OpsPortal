package probe

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

// ToolInfo 是输出指标时每个工具附带的标签信息。
type ToolInfo struct {
	ID          uint
	Name        string
	Project     string
	Environment string
	Category    string
	URL         string
}

// ContentType 是 Prometheus 文本格式的响应类型。
const ContentType = "text/plain; version=0.0.4; charset=utf-8"

var labelValueEscaper = strings.NewReplacer(`\`, `\\`, "\n", `\n`, `"`, `\"`)

// WriteMetrics 以 Prometheus 文本格式输出检测结果。
// tools 是开启了检测的工具，totalTools 是全部工具数。还没检测过的工具不输出（缺失即未知）。
func (p *Prober) WriteMetrics(w io.Writer, tools []ToolInfo, totalTools int) error {
	var b strings.Builder

	writeFamily(&b, "opsportal_tools", "Number of tool entries configured in OpsPortal.")
	fmt.Fprintf(&b, "opsportal_tools %d\n", totalTools)

	writeFamily(&b, "opsportal_probe_enabled", "Whether the built-in availability probe is enabled (1) or disabled (0).")
	fmt.Fprintf(&b, "opsportal_probe_enabled %d\n", boolToInt(p.Enabled()))

	type sample struct {
		labels string
		r      Result
	}
	var samples []sample
	if p.Enabled() {
		sorted := append([]ToolInfo(nil), tools...)
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].ID < sorted[j].ID })
		for _, t := range sorted {
			r, ok := p.Lookup(t.URL)
			if !ok {
				continue
			}
			samples = append(samples, sample{labels: toolLabels(t), r: r})
		}
	}

	writeFamily(&b, "opsportal_probe_success", "Whether the last probe of the tool URL succeeded (1) or failed (0).")
	for _, s := range samples {
		fmt.Fprintf(&b, "opsportal_probe_success{%s} %d\n", s.labels, boolToInt(s.r.Up))
	}

	writeFamily(&b, "opsportal_probe_duration_seconds", "Duration of the last probe of the tool URL in seconds.")
	for _, s := range samples {
		fmt.Fprintf(&b, "opsportal_probe_duration_seconds{%s} %s\n", s.labels, formatFloat(s.r.Duration.Seconds()))
	}

	writeFamily(&b, "opsportal_probe_http_status_code", "HTTP status code of the last probe response, 0 if no response was received.")
	for _, s := range samples {
		fmt.Fprintf(&b, "opsportal_probe_http_status_code{%s} %d\n", s.labels, s.r.StatusCode)
	}

	writeFamily(&b, "opsportal_probe_timestamp_seconds", "Unix time of the last probe of the tool URL.")
	for _, s := range samples {
		if s.r.CheckedAt.IsZero() {
			continue
		}
		fmt.Fprintf(&b, "opsportal_probe_timestamp_seconds{%s} %s\n", s.labels, formatUnix(s.r.CheckedAt.UnixMilli()))
	}

	if last := p.LastRound(); !last.IsZero() {
		writeFamily(&b, "opsportal_probe_last_round_timestamp_seconds", "Unix time when the last full probe round finished.")
		fmt.Fprintf(&b, "opsportal_probe_last_round_timestamp_seconds %s\n", formatUnix(last.UnixMilli()))
	}

	_, err := io.WriteString(w, b.String())
	return err
}

func writeFamily(b *strings.Builder, name, help string) {
	fmt.Fprintf(b, "# HELP %s %s\n# TYPE %s gauge\n", name, help, name)
}

func toolLabels(t ToolInfo) string {
	pairs := [][2]string{
		{"tool_id", strconv.FormatUint(uint64(t.ID), 10)},
		{"tool", t.Name},
		{"project", t.Project},
		{"environment", t.Environment},
		{"category", t.Category},
		{"url", strings.TrimSpace(t.URL)},
	}
	parts := make([]string, 0, len(pairs))
	for _, kv := range pairs {
		parts = append(parts, kv[0]+`="`+labelValueEscaper.Replace(kv[1])+`"`)
	}
	return strings.Join(parts, ",")
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func formatUnix(ms int64) string {
	return strconv.FormatFloat(float64(ms)/1000, 'f', 3, 64)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
