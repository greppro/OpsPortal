package probe

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestNormalizeURL(t *testing.T) {
	cases := []struct {
		in   string
		want string
		ok   bool
	}{
		{"https://grafana.example.com", "https://grafana.example.com", true},
		{"grafana.example.com:3000/d/abc", "https://grafana.example.com:3000/d/abc", true},
		{"  http://10.0.0.1:9090  ", "http://10.0.0.1:9090", true},
		{"HTTPS://Example.com/#/dashboard", "https://Example.com/", true},
		{"https//example.com/path", "https://example.com/path", true},
		{"javascript:alert(1)", "", false},
		{"data:text/html,<script>alert(1)</script>", "", false},
		{"ftp://example.com", "", false},
		{"file:///etc/passwd", "", false},
		{"https://", "", false},
		{"", "", false},
	}
	for _, c := range cases {
		got, err := NormalizeURL(c.in)
		if c.ok && (err != nil || got != c.want) {
			t.Errorf("NormalizeURL(%q) = %q, %v; want %q", c.in, got, err, c.want)
		}
		if !c.ok && err == nil {
			t.Errorf("NormalizeURL(%q) = %q; want error", c.in, got)
		}
	}
}

func TestIsUpStatus(t *testing.T) {
	for code, want := range map[int]bool{200: true, 204: true, 302: true, 401: true, 403: true, 404: false, 429: false, 500: false, 502: false} {
		if got := IsUpStatus(code); got != want {
			t.Errorf("IsUpStatus(%d) = %v, want %v", code, got, want)
		}
	}
}

func newTestProber(timeout time.Duration, insecure bool, targets func() []string) *Prober {
	return New(Config{Enabled: true, Interval: time.Minute, Timeout: timeout, Concurrency: 4, InsecureSkipVerify: insecure},
		func() ([]string, error) { return targets(), nil })
}

func TestCheckClassifiesResults(t *testing.T) {
	release := make(chan struct{})
	mux := http.NewServeMux()
	mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != userAgent {
			w.WriteHeader(http.StatusTeapot)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/ok", http.StatusFound) })
	mux.HandleFunc("/loop", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/loop", http.StatusFound) })
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusUnauthorized) })
	mux.HandleFunc("/missing", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusBadGateway) })
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-release:
		case <-time.After(5 * time.Second):
		}
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	defer close(release)

	// 拿一个确定没人监听的端口，用来测试连接被拒绝
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	refusedURL := "http://" + ln.Addr().String()
	ln.Close()

	p := newTestProber(300*time.Millisecond, false, func() []string { return nil })
	cases := []struct {
		url    string
		up     bool
		code   int
		reason string
	}{
		{srv.URL + "/ok", true, 200, ""},
		{srv.URL + "/redirect", true, 200, ""},
		{srv.URL + "/loop", true, 302, ""},
		{srv.URL + "/login", true, 401, ""},
		{srv.URL + "/missing", false, 404, ReasonHTTPStatus},
		{srv.URL + "/error", false, 502, ReasonHTTPStatus},
		{srv.URL + "/slow", false, 0, ReasonTimeout},
		{refusedURL, false, 0, ReasonRefused},
		{"ftp://example.com", false, 0, ReasonInvalidURL},
	}
	for _, c := range cases {
		r := p.Check(context.Background(), c.url)
		if r.Up != c.up || r.StatusCode != c.code || r.Reason != c.reason {
			t.Errorf("Check(%s) = up=%v code=%d reason=%q; want up=%v code=%d reason=%q",
				c.url, r.Up, r.StatusCode, r.Reason, c.up, c.code, c.reason)
		}
		if r.CheckedAt.IsZero() {
			t.Errorf("Check(%s) did not set CheckedAt", c.url)
		}
	}
}

func TestCheckTLS(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()

	strict := newTestProber(2*time.Second, false, func() []string { return nil })
	if r := strict.Check(context.Background(), srv.URL); r.Up || r.Reason != ReasonTLS {
		t.Errorf("self-signed cert: up=%v reason=%q, want down with %q", r.Up, r.Reason, ReasonTLS)
	}

	insecure := newTestProber(2*time.Second, true, func() []string { return nil })
	if r := insecure.Check(context.Background(), srv.URL); !r.Up {
		t.Errorf("self-signed cert with InsecureSkipVerify: up=false reason=%q", r.Reason)
	}
}

func TestRunOnceDedupesAndPrunes(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
	}))
	defer srv.Close()

	targets := []string{srv.URL + "/a", srv.URL + "/a", srv.URL + "/a#tab", "  " + srv.URL + "/a  ", srv.URL + "/b", "bad url::"}
	p := newTestProber(2*time.Second, false, func() []string { return targets })

	p.RunOnce(context.Background())
	if got := atomic.LoadInt32(&hits); got != 2 {
		t.Fatalf("expected 2 requests after dedupe, got %d", got)
	}
	if r, ok := p.Lookup(srv.URL + "/a#other"); !ok || !r.Up {
		t.Fatalf("lookup /a: ok=%v up=%v", ok, r.Up)
	}
	if r, ok := p.Lookup("bad url::"); !ok || r.Reason != ReasonInvalidURL {
		t.Fatalf("lookup invalid url: ok=%v reason=%q", ok, r.Reason)
	}
	if p.LastRound().IsZero() {
		t.Fatal("LastRound not set")
	}

	// 单独检测、且晚于本轮开始的结果要保留（比如刚新增的工具）
	fresh := srv.URL + "/fresh"
	p.store(Result{URL: fresh, Up: true, CheckedAt: time.Now().Add(time.Hour)})

	targets = []string{srv.URL + "/b"}
	p.RunOnce(context.Background())
	if _, ok := p.Lookup(srv.URL + "/a"); ok {
		t.Fatal("result for removed target /a should be pruned")
	}
	if _, ok := p.Lookup(fresh); !ok {
		t.Fatal("result checked after the round started should be kept")
	}
}

var sampleLine = regexp.MustCompile(`^[a-zA-Z_:][a-zA-Z0-9_:]*(\{.*\})? -?[0-9.]+$`)

func TestWriteMetrics(t *testing.T) {
	p := New(Config{Enabled: true, Interval: time.Minute, Timeout: time.Second, Concurrency: 1}, nil)
	p.store(Result{URL: "https://a.example.com", Up: true, StatusCode: 200, Duration: 150 * time.Millisecond, CheckedAt: time.UnixMilli(1700000000123)})
	p.store(Result{URL: "https://c.example.com", Up: false, Reason: ReasonTimeout, Duration: 5 * time.Second, CheckedAt: time.UnixMilli(1700000000456)})

	tools := []ToolInfo{
		{ID: 3, Name: "C", Project: "p", Environment: "prod", URL: "https://c.example.com"},
		{ID: 2, Name: `Graf"ana\`, Project: "monitor", Environment: "prod", Category: "监控\n", URL: "https://a.example.com"},
		{ID: 1, Name: "pending", Project: "p", Environment: "dev", URL: "https://b.example.com"},
	}
	var buf bytes.Buffer
	if err := p.WriteMetrics(&buf, tools, 5); err != nil {
		t.Fatal(err)
	}
	out := buf.String()

	labelsA := `tool_id="2",tool="Graf\"ana\\",project="monitor",environment="prod",category="监控\n",url="https://a.example.com"`
	for _, want := range []string{
		"opsportal_tools 5\n",
		"opsportal_probe_enabled 1\n",
		"opsportal_probe_success{" + labelsA + "} 1\n",
		"opsportal_probe_duration_seconds{" + labelsA + "} 0.15\n",
		"opsportal_probe_http_status_code{" + labelsA + "} 200\n",
		"opsportal_probe_timestamp_seconds{" + labelsA + "} 1700000000.123\n",
		`opsportal_probe_success{tool_id="3",tool="C",project="p",environment="prod",category="",url="https://c.example.com"} 0` + "\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("metrics output missing %q\n---\n%s", want, out)
		}
	}
	if strings.Contains(out, `tool_id="1"`) {
		t.Error("tool without a probe result should not be exported")
	}
	if strings.Index(out, `tool_id="2"`) > strings.Index(out, `tool_id="3"`) {
		t.Error("samples should be sorted by tool id")
	}

	// 每一行要么是注释，要么是合法的样本行；每个指标族先有 TYPE 再有样本
	typed := map[string]bool{}
	for _, line := range strings.Split(strings.TrimSuffix(out, "\n"), "\n") {
		if strings.HasPrefix(line, "# TYPE ") {
			typed[strings.Fields(line)[2]] = true
			continue
		}
		if strings.HasPrefix(line, "# HELP ") {
			continue
		}
		if !sampleLine.MatchString(line) {
			t.Errorf("invalid sample line: %q", line)
			continue
		}
		name := line
		if i := strings.IndexAny(line, "{ "); i >= 0 {
			name = line[:i]
		}
		if !typed[name] {
			t.Errorf("sample %q appears before its TYPE line", name)
		}
	}
}

func TestWriteMetricsDisabled(t *testing.T) {
	p := New(Config{Enabled: false}, nil)
	p.store(Result{URL: "https://a.example.com", Up: true})
	var buf bytes.Buffer
	if err := p.WriteMetrics(&buf, []ToolInfo{{ID: 1, URL: "https://a.example.com"}}, 1); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "opsportal_probe_enabled 0\n") || strings.Contains(out, "tool_id=") {
		t.Errorf("unexpected output when probe disabled:\n%s", out)
	}
}
