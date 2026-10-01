package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Settings 运行配置，全部来自环境变量，未设置时使用默认值。
type Settings struct {
	Host             string   // HOST：监听地址，默认空表示所有网卡；只想本机访问时设为 127.0.0.1
	Port             string   // PORT：监听端口，默认 8080
	DBPath           string   // DB_PATH：SQLite 文件路径，默认 data/opsportal.db
	JWTSecret        string   // JWT_SECRET：签名密钥；为空时自动生成并保存在数据库里
	AdminPassword    string   // ADMIN_PASSWORD：首次启动创建管理员时使用；为空时随机生成并打印到日志
	CORSAllowOrigins []string // CORS_ALLOW_ORIGINS：允许跨域的来源，逗号分隔；为空时不开启跨域
	MetricsToken     string   // METRICS_TOKEN：设置后 /metrics 和 /prometheus/targets 需要 Bearer Token
	Probe            ProbeSettings
}

// ProbeSettings 内置可用性检测配置。
type ProbeSettings struct {
	Enabled            bool          // PROBE_ENABLED：是否开启，默认 true
	Interval           time.Duration // PROBE_INTERVAL：检测间隔，默认 60s，最小 10s
	Timeout            time.Duration // PROBE_TIMEOUT：单次请求超时，默认 5s
	Concurrency        int           // PROBE_CONCURRENCY：并发数，默认 10
	InsecureSkipVerify bool          // PROBE_INSECURE_SKIP_VERIFY：跳过 TLS 证书校验，默认 false
}

const (
	minProbeInterval = 10 * time.Second
	minProbeTimeout  = time.Second
	maxConcurrency   = 100
)

// LoadSettings 从环境变量读取配置。非法值打印警告后回退到默认值，不会导致启动失败。
func LoadSettings() Settings {
	s := Settings{
		Host:             strings.TrimSpace(os.Getenv("HOST")),
		Port:             envString("PORT", "8080"),
		DBPath:           envString("DB_PATH", "data/opsportal.db"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		AdminPassword:    os.Getenv("ADMIN_PASSWORD"),
		CORSAllowOrigins: envList("CORS_ALLOW_ORIGINS"),
		MetricsToken:     os.Getenv("METRICS_TOKEN"),
		Probe: ProbeSettings{
			Enabled:            envBool("PROBE_ENABLED", true),
			Interval:           envDuration("PROBE_INTERVAL", 60*time.Second),
			Timeout:            envDuration("PROBE_TIMEOUT", 5*time.Second),
			Concurrency:        envInt("PROBE_CONCURRENCY", 10),
			InsecureSkipVerify: envBool("PROBE_INSECURE_SKIP_VERIFY", false),
		},
	}

	if s.Probe.Interval < minProbeInterval {
		log.Printf("配置警告：PROBE_INTERVAL 不能小于 %s，已改为 %s", minProbeInterval, minProbeInterval)
		s.Probe.Interval = minProbeInterval
	}
	if s.Probe.Timeout < minProbeTimeout {
		log.Printf("配置警告：PROBE_TIMEOUT 不能小于 %s，已改为 %s", minProbeTimeout, minProbeTimeout)
		s.Probe.Timeout = minProbeTimeout
	}
	if s.Probe.Timeout > s.Probe.Interval {
		log.Printf("配置警告：PROBE_TIMEOUT 不能大于 PROBE_INTERVAL，已改为 %s", s.Probe.Interval)
		s.Probe.Timeout = s.Probe.Interval
	}
	if s.Probe.Concurrency < 1 || s.Probe.Concurrency > maxConcurrency {
		log.Printf("配置警告：PROBE_CONCURRENCY 需在 1 到 %d 之间，已改为 10", maxConcurrency)
		s.Probe.Concurrency = 10
	}
	if s.JWTSecret != "" && len(s.JWTSecret) < 32 {
		log.Printf("配置警告：JWT_SECRET 少于 32 个字符，建议使用更长的随机字符串")
	}
	return s
}

func envString(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func envList(key string) []string {
	var out []string
	for _, part := range strings.Split(os.Getenv(key), ",") {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func envBool(key string, def bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Printf("配置警告：%s=%q 不是合法的布尔值，使用默认值 %v", key, v, def)
		return def
	}
	return b
}

func envInt(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("配置警告：%s=%q 不是合法的整数，使用默认值 %d", key, v, def)
		return def
	}
	return n
}

// envDuration 支持 Go 时长格式（30s、5m）和纯数字（按秒计）。
func envDuration(key string, def time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	if n, err := strconv.Atoi(v); err == nil {
		return time.Duration(n) * time.Second
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Printf("配置警告：%s=%q 不是合法的时长，使用默认值 %s", key, v, def)
		return def
	}
	return d
}
