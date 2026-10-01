// Package probe 提供内置的 URL 可用性检测：按固定间隔并发请求所有工具地址，
// 在内存里保存每个地址最近一次的结果，供首页状态展示和 Prometheus 指标使用。
package probe

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
)

// 失败原因（接口和指标里使用的固定取值）
const (
	ReasonHTTPStatus = "http_status"        // 收到响应，但状态码表示不可用
	ReasonTimeout    = "timeout"            // 请求超时
	ReasonDNS        = "dns"                // 域名解析失败
	ReasonRefused    = "connection_refused" // 连接被拒绝
	ReasonTLS        = "tls"                // TLS 握手或证书校验失败
	ReasonInvalidURL = "invalid_url"        // 地址格式无效，未发起请求
	ReasonNetwork    = "network"            // 其他网络错误
)

const (
	maxRedirects = 5
	userAgent    = "OpsPortal-Probe/1.0 (+https://github.com/greppro/OpsPortal)"
)

// Config 检测配置。
type Config struct {
	Enabled            bool
	Interval           time.Duration
	Timeout            time.Duration
	Concurrency        int
	InsecureSkipVerify bool
}

// Result 是某个地址最近一次检测的结果。
type Result struct {
	URL        string
	Up         bool
	StatusCode int // 未收到 HTTP 响应时为 0
	Duration   time.Duration
	Reason     string // Up 为 false 时说明原因
	CheckedAt  time.Time
}

// TargetSource 返回需要检测的原始地址列表（重复地址只检测一次）。
type TargetSource func() ([]string, error)

// Prober 负责定时检测并保存结果，可并发使用。
type Prober struct {
	cfg    Config
	client *http.Client
	source TargetSource

	mu        sync.RWMutex
	results   map[string]Result // key 为规范化后的地址
	lastRound time.Time
}

// New 创建检测器。source 每轮检测前调用一次。
func New(cfg Config, source TargetSource) *Prober {
	if cfg.Concurrency < 1 {
		cfg.Concurrency = 1
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 5 * time.Second
	}
	if cfg.Interval <= 0 {
		cfg.Interval = time.Minute
	}

	// 沿用默认 Transport：支持 HTTP_PROXY / HTTPS_PROXY / NO_PROXY 环境变量
	transport := http.DefaultTransport.(*http.Transport).Clone()
	// 每次检测都新建连接，测的是真实的连通性，而不是连接池里的旧连接
	transport.DisableKeepAlives = true
	// 内网自签证书场景可通过 PROBE_INSECURE_SKIP_VERIFY 显式关闭校验
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: cfg.InsecureSkipVerify} // #nosec G402

	return &Prober{
		cfg:     cfg,
		source:  source,
		results: make(map[string]Result),
		client: &http.Client{
			Timeout:   cfg.Timeout,
			Transport: transport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= maxRedirects {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
	}
}

// Enabled 返回是否开启了后台定时检测。
func (p *Prober) Enabled() bool { return p.cfg.Enabled }

// Interval 返回检测间隔。
func (p *Prober) Interval() time.Duration { return p.cfg.Interval }

// Timeout 返回单次请求超时。
func (p *Prober) Timeout() time.Duration { return p.cfg.Timeout }

// LastRound 返回最近一轮检测完成的时间，还没跑完第一轮时为零值。
func (p *Prober) LastRound() time.Time {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.lastRound
}

// Run 立即检测一轮，然后按间隔循环，直到 ctx 结束。未开启时直接返回。
func (p *Prober) Run(ctx context.Context) {
	if !p.cfg.Enabled {
		log.Printf("可用性检测未开启（PROBE_ENABLED=false）")
		return
	}
	log.Printf("可用性检测已开启：每 %s 一轮，超时 %s，并发 %d", p.cfg.Interval, p.cfg.Timeout, p.cfg.Concurrency)

	p.RunOnce(ctx)
	ticker := time.NewTicker(p.cfg.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.RunOnce(ctx)
		}
	}
}

// RunOnce 检测一轮：读取全部目标，去重后并发检测，并清理已不存在的地址的旧结果。
func (p *Prober) RunOnce(ctx context.Context) {
	raws, err := p.source()
	if err != nil {
		log.Printf("可用性检测：读取检测目标失败：%v", err)
		return
	}
	targets := uniqueTargets(raws)
	roundStart := time.Now()

	sem := make(chan struct{}, p.cfg.Concurrency)
	var wg sync.WaitGroup
loop:
	for _, target := range targets {
		select {
		case <-ctx.Done():
			break loop
		case sem <- struct{}{}:
		}
		wg.Add(1)
		go func(target string) {
			defer wg.Done()
			defer func() { <-sem }()
			p.store(p.probe(ctx, target))
		}(target)
	}
	wg.Wait()

	keep := make(map[string]struct{}, len(targets))
	for _, t := range targets {
		keep[t] = struct{}{}
	}
	p.mu.Lock()
	for key, r := range p.results {
		// 本轮开始后单独检测过的结果（比如刚新增的工具）先保留
		if _, ok := keep[key]; !ok && r.CheckedAt.Before(roundStart) {
			delete(p.results, key)
		}
	}
	p.lastRound = time.Now()
	p.mu.Unlock()
}

// Check 立即检测一个地址并保存结果（不受后台定时开关影响）。
func (p *Prober) Check(ctx context.Context, raw string) Result {
	target, err := NormalizeURL(raw)
	if err != nil {
		return Result{URL: strings.TrimSpace(raw), Reason: ReasonInvalidURL, CheckedAt: time.Now()}
	}
	r := p.probe(ctx, target)
	p.store(r)
	return r
}

// Lookup 返回地址最近一次的检测结果。地址格式无效时直接返回 invalid_url 结果；
// 还没检测过时 ok 为 false。
func (p *Prober) Lookup(raw string) (Result, bool) {
	target, err := NormalizeURL(raw)
	if err != nil {
		return Result{URL: strings.TrimSpace(raw), Reason: ReasonInvalidURL}, true
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	r, ok := p.results[target]
	return r, ok
}

func (p *Prober) store(r Result) {
	p.mu.Lock()
	p.results[r.URL] = r
	p.mu.Unlock()
}

func (p *Prober) probe(ctx context.Context, target string) Result {
	res := Result{URL: target}
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		res.Reason = ReasonInvalidURL
		res.CheckedAt = time.Now()
		return res
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := p.client.Do(req)
	res.Duration = time.Since(start)
	res.CheckedAt = time.Now()
	if err != nil {
		res.Reason = classifyError(err)
		return res
	}
	resp.Body.Close()

	res.StatusCode = resp.StatusCode
	res.Up = IsUpStatus(resp.StatusCode)
	if !res.Up {
		res.Reason = ReasonHTTPStatus
	}
	return res
}

// IsUpStatus 判断状态码是否算"可用"：小于 400 算可用（重定向会先跟随），
// 401/403 说明服务在线、只是需要登录，也算可用。
func IsUpStatus(code int) bool {
	return code < 400 || code == http.StatusUnauthorized || code == http.StatusForbidden
}

func classifyError(err error) string {
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return ReasonDNS
	}

	var certErr *tls.CertificateVerificationError
	var unknownAuthority x509.UnknownAuthorityError
	var hostnameErr x509.HostnameError
	var invalidCert x509.CertificateInvalidError
	var recordErr tls.RecordHeaderError
	if errors.As(err, &certErr) || errors.As(err, &unknownAuthority) || errors.As(err, &hostnameErr) ||
		errors.As(err, &invalidCert) || errors.As(err, &recordErr) ||
		strings.Contains(err.Error(), "HTTP response to HTTPS client") ||
		strings.Contains(err.Error(), "tls: ") {
		return ReasonTLS
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return ReasonTimeout
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return ReasonTimeout
	}
	if errors.Is(err, syscall.ECONNREFUSED) {
		return ReasonRefused
	}
	return ReasonNetwork
}

var (
	schemePattern       = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.-]*://`)
	missingColonPattern = regexp.MustCompile(`(?i)^https?/+`)
	missingColonPrefix  = regexp.MustCompile(`(?i)^https?/+:?`)
	errEmptyURL         = errors.New("地址为空")
	errUnsupportedURL   = errors.New("只支持 http 和 https 地址")
	errMissingHost      = errors.New("地址缺少主机名")
)

// NormalizeURL 把工具地址整理成可请求的 http(s) 地址，规则与前端打开链接时一致：
// 没写协议时补 https://，"https//host" 这类漏写冒号的地址修正为 https://host。
func NormalizeURL(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", errEmptyURL
	}
	if missingColonPattern.MatchString(s) {
		s = missingColonPrefix.ReplaceAllString(s, "https://")
	}
	if !schemePattern.MatchString(s) {
		s = "https://" + s
	}
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", errUnsupportedURL
	}
	if u.Hostname() == "" {
		return "", errMissingHost
	}
	u.Scheme = scheme
	u.Fragment = ""
	u.RawFragment = ""
	return u.String(), nil
}

func uniqueTargets(raws []string) []string {
	seen := make(map[string]struct{}, len(raws))
	out := make([]string, 0, len(raws))
	for _, raw := range raws {
		target, err := NormalizeURL(raw)
		if err != nil {
			continue
		}
		if _, ok := seen[target]; ok {
			continue
		}
		seen[target] = struct{}{}
		out = append(out, target)
	}
	sort.Strings(out)
	return out
}
