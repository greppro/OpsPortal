package handlers

import (
	"context"
	"crypto/subtle"
	"net/http"
	"ops-portal/config"
	"ops-portal/models"
	"ops-portal/probe"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var prober *probe.Prober

// SetProber 注入检测器（main 启动时调用）。
func SetProber(p *probe.Prober) {
	prober = p
}

// ProbeTargets 是后台检测的目标来源：所有开启了检测的工具地址。
func ProbeTargets() ([]string, error) {
	var urls []string
	err := config.DB.Model(&models.Tool{}).Where("probe_disabled = ?", false).Pluck("url", &urls).Error
	return urls, err
}

// triggerProbe 保存工具后立即检测一次，不用等下一轮。
func triggerProbe(t models.Tool) {
	if prober == nil || !prober.Enabled() || t.ProbeDisabled {
		return
	}
	go prober.Check(context.Background(), t.URL)
}

type probeStatus struct {
	State      string     `json:"state"` // up / down / pending
	StatusCode int        `json:"status_code,omitempty"`
	DurationMs int64      `json:"duration_ms"`
	Reason     string     `json:"reason,omitempty"`
	CheckedAt  *time.Time `json:"checked_at,omitempty"`
}

func toProbeStatus(r probe.Result, ok bool) probeStatus {
	if !ok {
		return probeStatus{State: "pending"}
	}
	s := probeStatus{
		State:      "down",
		StatusCode: r.StatusCode,
		DurationMs: r.Duration.Milliseconds(),
		Reason:     r.Reason,
	}
	if r.Up {
		s.State = "up"
	}
	if !r.CheckedAt.IsZero() {
		t := r.CheckedAt
		s.CheckedAt = &t
	}
	return s
}

// GetProbeStatus 返回每个工具最近一次的检测结果，key 为工具 ID（公开，首页使用）。
func GetProbeStatus(c *gin.Context) {
	if prober == nil || !prober.Enabled() {
		c.JSON(http.StatusOK, gin.H{"enabled": false, "results": gin.H{}})
		return
	}

	var tools []models.Tool
	if err := config.DB.Select("id", "url").Where("probe_disabled = ?", false).Find(&tools).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取检测结果失败"})
		return
	}
	results := make(map[string]probeStatus, len(tools))
	for _, t := range tools {
		r, ok := prober.Lookup(t.URL)
		results[strconv.FormatUint(uint64(t.ID), 10)] = toProbeStatus(r, ok)
	}

	resp := gin.H{
		"enabled":          true,
		"interval_seconds": int(prober.Interval() / time.Second),
		"results":          results,
	}
	if last := prober.LastRound(); !last.IsZero() {
		resp["last_round_at"] = last
	}
	c.JSON(http.StatusOK, resp)
}

// CheckToolNow 立即检测一个工具并返回结果（需要登录）。
func CheckToolNow(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if prober == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "可用性检测未初始化"})
		return
	}
	var tool models.Tool
	if err := config.DB.First(&tool, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "工具不存在"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), prober.Timeout()+2*time.Second)
	defer cancel()
	c.JSON(http.StatusOK, toProbeStatus(prober.Check(ctx, tool.URL), true))
}

// Metrics 以 Prometheus 文本格式输出检测结果。
func Metrics(c *gin.Context) {
	if prober == nil {
		c.String(http.StatusServiceUnavailable, "probe not initialized\n")
		return
	}
	var tools []models.Tool
	if err := config.DB.Find(&tools).Error; err != nil {
		c.String(http.StatusInternalServerError, "failed to load tools\n")
		return
	}
	infos := make([]probe.ToolInfo, 0, len(tools))
	for _, t := range tools {
		if t.ProbeDisabled {
			continue
		}
		infos = append(infos, probe.ToolInfo{
			ID:          t.ID,
			Name:        t.Name,
			Project:     t.Project,
			Environment: t.Environment,
			Category:    t.Category,
			URL:         t.URL,
		})
	}
	c.Header("Content-Type", probe.ContentType)
	c.Status(http.StatusOK)
	_ = prober.WriteMetrics(c.Writer, infos, len(tools))
}

type sdTargetGroup struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

// PrometheusTargets 是 Prometheus http_sd_configs 服务发现接口，把开启检测的工具地址
// 作为探测目标输出，供 blackbox_exporter 之类的探测器使用。
// 可选查询参数 project / environment / category 用来按项目标识、环境标识、分类筛选。
func PrometheusTargets(c *gin.Context) {
	query := config.DB.Model(&models.Tool{}).Where("probe_disabled = ?", false)
	for _, f := range []string{"project", "environment", "category"} {
		if v := strings.TrimSpace(c.Query(f)); v != "" {
			query = query.Where(f+" = ?", v)
		}
	}
	var tools []models.Tool
	if err := query.Order("id").Find(&tools).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取检测目标失败"})
		return
	}

	groups := make([]sdTargetGroup, 0, len(tools))
	for _, t := range tools {
		target, err := probe.NormalizeURL(t.URL)
		if err != nil {
			continue
		}
		labels := map[string]string{
			"tool_id":     strconv.FormatUint(uint64(t.ID), 10),
			"tool":        t.Name,
			"project":     t.Project,
			"environment": t.Environment,
		}
		if cat := strings.TrimSpace(t.Category); cat != "" {
			labels["category"] = cat
		}
		groups = append(groups, sdTargetGroup{Targets: []string{target}, Labels: labels})
	}
	c.JSON(http.StatusOK, groups)
}

// MetricsAuth 设置了 METRICS_TOKEN 时，要求 Authorization: Bearer <token>。
func MetricsAuth(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if token == "" {
			c.Next()
			return
		}
		got := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if subtle.ConstantTimeCompare([]byte(got), []byte(token)) != 1 {
			c.Header("WWW-Authenticate", `Bearer realm="opsportal"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
