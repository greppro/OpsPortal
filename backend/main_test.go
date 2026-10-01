package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"ops-portal/config"
	"ops-portal/handlers"
	"ops-portal/models"
	"ops-portal/probe"
	"ops-portal/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const testAdminPassword = "initial-pass-123"

type testApp struct {
	router   *gin.Engine
	settings config.Settings
	prober   *probe.Prober
}

func newTestApp(t *testing.T, mutate func(*config.Settings)) *testApp {
	t.Helper()
	gin.SetMode(gin.TestMode)

	dir := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// 上传目录是相对路径，切到临时目录避免写进源码树
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	s := config.Settings{
		Port:          "0",
		DBPath:        filepath.Join(dir, "data", "opsportal.db"),
		AdminPassword: testAdminPassword,
		Probe:         config.ProbeSettings{Enabled: true, Interval: time.Minute, Timeout: 2 * time.Second, Concurrency: 4},
	}
	if mutate != nil {
		mutate(&s)
	}
	app := &testApp{settings: s}
	app.restart(t)
	return app
}

// restart 模拟进程重启：在同一个数据库上重新执行启动流程。
func (a *testApp) restart(t *testing.T) {
	t.Helper()
	config.InitDB(a.settings)
	if sqlDB, err := config.DB.DB(); err == nil {
		t.Cleanup(func() { _ = sqlDB.Close() })
	}
	secret, err := config.InitJWTSecret(a.settings)
	if err != nil {
		t.Fatal(err)
	}
	utils.SetJWTSecret(secret)

	// 不启动后台循环，测试里需要时手动调用 RunOnce
	a.prober = probe.New(probe.Config{
		Enabled:     a.settings.Probe.Enabled,
		Interval:    a.settings.Probe.Interval,
		Timeout:     a.settings.Probe.Timeout,
		Concurrency: a.settings.Probe.Concurrency,
	}, handlers.ProbeTargets)
	handlers.SetProber(a.prober)
	a.router = setupRouter(a.settings)
}

func (a *testApp) do(t *testing.T, method, path, token string, body interface{}, headers ...string) *httptest.ResponseRecorder {
	t.Helper()
	var reader *bytes.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		reader = bytes.NewReader(raw)
	} else {
		reader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	for i := 0; i+1 < len(headers); i += 2 {
		req.Header.Set(headers[i], headers[i+1])
	}
	rec := httptest.NewRecorder()
	a.router.ServeHTTP(rec, req)
	return rec
}

func decode(t *testing.T, rec *httptest.ResponseRecorder, v interface{}) {
	t.Helper()
	if err := json.Unmarshal(rec.Body.Bytes(), v); err != nil {
		t.Fatalf("decode %q: %v", rec.Body.String(), err)
	}
}

func (a *testApp) login(t *testing.T, username, password string) (string, bool) {
	t.Helper()
	rec := a.do(t, "POST", "/api/auth/login", "", map[string]string{"username": username, "password": password})
	if rec.Code != http.StatusOK {
		t.Fatalf("login %s: status %d body %s", username, rec.Code, rec.Body.String())
	}
	var resp models.LoginResponse
	decode(t, rec, &resp)
	return resp.Token, resp.DefaultPassword
}

func expectStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, want, rec.Body.String())
	}
}

func TestRestartKeepsEnvironments(t *testing.T) {
	app := newTestApp(t, nil)
	token, _ := app.login(t, "admin", testAdminPassword)

	var before int64
	config.DB.Model(&models.Environment{}).Count(&before)

	rec := app.do(t, "POST", "/api/environments", token, map[string]interface{}{"name": "uat", "label": "验收环境", "project_id": 1})
	expectStatus(t, rec, http.StatusOK)

	app.restart(t)

	var envs []models.Environment
	config.DB.Find(&envs)
	if int64(len(envs)) != before+1 {
		t.Fatalf("environment count after restart = %d, want %d", len(envs), before+1)
	}
	found := false
	for _, e := range envs {
		found = found || (e.Name == "uat" && e.ProjectID == 1)
	}
	if !found {
		t.Fatal("custom environment was lost after restart")
	}
}

func TestInitialAdminAndLogin(t *testing.T) {
	app := newTestApp(t, nil)

	var user models.User
	config.DB.Where("username = ?", "admin").First(&user)
	if !utils.IsPasswordHash(user.Password) {
		t.Fatalf("admin password should be stored as bcrypt hash, got %q", user.Password)
	}

	if _, isDefault := app.login(t, "admin", testAdminPassword); isDefault {
		t.Error("default_password should be false for a custom password")
	}
	expectStatus(t, app.do(t, "POST", "/api/auth/login", "", map[string]string{"username": "admin", "password": "admin123"}), http.StatusUnauthorized)
	expectStatus(t, app.do(t, "POST", "/api/auth/login", "", map[string]string{"username": "nobody", "password": "whatever"}), http.StatusUnauthorized)
}

func TestRandomInitialPassword(t *testing.T) {
	newTestApp(t, func(s *config.Settings) { s.AdminPassword = "" })
	var user models.User
	if err := config.DB.Where("username = ?", "admin").First(&user).Error; err != nil {
		t.Fatal(err)
	}
	if !utils.IsPasswordHash(user.Password) || utils.CheckPassword(user.Password, "admin123") {
		t.Fatal("generated admin password should be a random bcrypt-hashed value")
	}
}

func TestLegacyPlaintextPasswordMigrated(t *testing.T) {
	app := newTestApp(t, nil)
	// 模拟旧版本留下的明文密码
	config.DB.Model(&models.User{}).Where("username = ?", "admin").Update("password", "admin123")

	app.restart(t)

	var user models.User
	config.DB.Where("username = ?", "admin").First(&user)
	if !utils.IsPasswordHash(user.Password) {
		t.Fatalf("plaintext password was not migrated: %q", user.Password)
	}
	if _, isDefault := app.login(t, "admin", "admin123"); !isDefault {
		t.Error("login with the legacy default password should report default_password=true")
	}
}

func TestJWTSecretPersistsAndForgedTokensRejected(t *testing.T) {
	app := newTestApp(t, nil)
	token, _ := app.login(t, "admin", testAdminPassword)

	app.restart(t)
	expectStatus(t, app.do(t, "GET", "/api/notices", token, nil), http.StatusOK)

	// 用旧版本公开的密钥签发的 token 必须无效
	claims := &utils.Claims{Username: "admin", RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))}}
	forged, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("your-secret-key"))
	if err != nil {
		t.Fatal(err)
	}
	expectStatus(t, app.do(t, "GET", "/api/notices", forged, nil), http.StatusUnauthorized)

	// alg=none 的 token 也必须无效
	none, err := jwt.NewWithClaims(jwt.SigningMethodNone, claims).SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatal(err)
	}
	expectStatus(t, app.do(t, "GET", "/api/notices", none, nil), http.StatusUnauthorized)

	// JWT_SECRET 优先于数据库里的密钥
	s := app.settings
	s.JWTSecret = strings.Repeat("x", 40)
	secret, err := config.InitJWTSecret(s)
	if err != nil || string(secret) != s.JWTSecret {
		t.Fatalf("JWT_SECRET should take precedence, got %q, %v", secret, err)
	}
}

func TestChangePasswordInvalidatesOldTokens(t *testing.T) {
	app := newTestApp(t, nil)
	oldToken, _ := app.login(t, "admin", testAdminPassword)

	expectStatus(t, app.do(t, "POST", "/api/auth/change-password", oldToken,
		map[string]string{"oldPassword": testAdminPassword, "newPassword": "short"}), http.StatusBadRequest)
	expectStatus(t, app.do(t, "POST", "/api/auth/change-password", oldToken,
		map[string]string{"oldPassword": "wrong-password", "newPassword": "new-password-456"}), http.StatusBadRequest)

	rec := app.do(t, "POST", "/api/auth/change-password", oldToken,
		map[string]string{"oldPassword": testAdminPassword, "newPassword": "new-password-456"})
	expectStatus(t, rec, http.StatusOK)
	var resp models.ChangePasswordResponse
	decode(t, rec, &resp)
	if resp.Token == "" {
		t.Fatal("change-password should return a new token")
	}

	rec = app.do(t, "GET", "/api/notices", oldToken, nil)
	expectStatus(t, rec, http.StatusUnauthorized)
	if !strings.Contains(rec.Body.String(), "token无效或已过期") {
		t.Errorf("frontend relies on this error text to redirect to login, got %s", rec.Body.String())
	}
	expectStatus(t, app.do(t, "GET", "/api/notices", resp.Token, nil), http.StatusOK)

	expectStatus(t, app.do(t, "POST", "/api/auth/login", "", map[string]string{"username": "admin", "password": testAdminPassword}), http.StatusUnauthorized)
	app.login(t, "admin", "new-password-456")
}

func createTool(t *testing.T, app *testApp, token string, body map[string]interface{}) models.Tool {
	t.Helper()
	rec := app.do(t, "POST", "/api/sites", token, body)
	expectStatus(t, rec, http.StatusOK)
	var tool models.Tool
	decode(t, rec, &tool)
	return tool
}

func TestUpdateToolIgnoresIDInBody(t *testing.T) {
	app := newTestApp(t, func(s *config.Settings) { s.Probe.Enabled = false })
	token, _ := app.login(t, "admin", testAdminPassword)

	a := createTool(t, app, token, map[string]interface{}{"name": "A", "url": "https://a.example.com", "project": "cloud", "environment": "dev"})
	b := createTool(t, app, token, map[string]interface{}{"name": "B", "url": "https://b.example.com", "project": "cloud", "environment": "dev"})

	rec := app.do(t, "PUT", "/api/sites/"+itoa(a.ID), token, map[string]interface{}{
		"id": b.ID, "name": "A2", "url": "https://a2.example.com", "project": "cloud", "environment": "dev",
	})
	expectStatus(t, rec, http.StatusOK)

	var gotA, gotB models.Tool
	config.DB.First(&gotA, a.ID)
	config.DB.First(&gotB, b.ID)
	if gotA.Name != "A2" || gotB.Name != "B" || gotB.URL != "https://b.example.com" {
		t.Fatalf("update leaked to another record: A=%+v B=%+v", gotA, gotB)
	}
}

func TestToolValidation(t *testing.T) {
	app := newTestApp(t, func(s *config.Settings) { s.Probe.Enabled = false })
	token, _ := app.login(t, "admin", testAdminPassword)

	for _, bad := range []string{"javascript:alert(1)", "data:text/html,<b>x</b>", "ftp://files.example.com", "file:///etc/passwd"} {
		rec := app.do(t, "POST", "/api/sites", token, map[string]interface{}{"name": "X", "url": bad, "project": "cloud", "environment": "dev"})
		expectStatus(t, rec, http.StatusBadRequest)
	}
	expectStatus(t, app.do(t, "POST", "/api/sites", token, map[string]interface{}{"name": "", "url": "https://x.example.com", "project": "cloud", "environment": "dev"}), http.StatusBadRequest)

	tool := createTool(t, app, token, map[string]interface{}{"name": "NoScheme", "url": "grafana.example.com:3000", "project": "cloud", "environment": "dev"})
	if tool.URL != "grafana.example.com:3000" {
		t.Errorf("url should be stored as entered, got %q", tool.URL)
	}

	// 旧数据里的非法地址：只改别的字段时不应因为地址校验失败而无法保存
	config.DB.Model(&models.Tool{}).Where("id = ?", tool.ID).Update("url", "legacy bad url")
	rec := app.do(t, "PUT", "/api/sites/"+itoa(tool.ID), token, map[string]interface{}{"name": "Renamed", "url": "legacy bad url", "project": "cloud", "environment": "dev"})
	expectStatus(t, rec, http.StatusOK)
}

func TestNonNumericIDRejected(t *testing.T) {
	app := newTestApp(t, nil)
	token, _ := app.login(t, "admin", testAdminPassword)

	expectStatus(t, app.do(t, "POST", "/api/notices", token, map[string]interface{}{"content": "keep me", "active": true}), http.StatusOK)

	rec := app.do(t, "DELETE", "/api/notices/1%20OR%201=1", token, nil)
	expectStatus(t, rec, http.StatusBadRequest)
	var count int64
	config.DB.Model(&models.Notice{}).Count(&count)
	if count != 1 {
		t.Fatalf("notice should not be deleted by an injected id, count=%d", count)
	}
	expectStatus(t, app.do(t, "PUT", "/api/sites/abc", token, map[string]interface{}{}), http.StatusBadRequest)
	expectStatus(t, app.do(t, "DELETE", "/api/projects/1=1", token, nil), http.StatusBadRequest)
}

func TestNoticeDraftsArePrivate(t *testing.T) {
	app := newTestApp(t, nil)
	token, _ := app.login(t, "admin", testAdminPassword)

	expectStatus(t, app.do(t, "GET", "/api/notices", "", nil), http.StatusUnauthorized)

	expectStatus(t, app.do(t, "POST", "/api/notices", token, map[string]interface{}{"content": "draft", "active": false}), http.StatusOK)
	rec := app.do(t, "GET", "/api/notices/active", "", nil)
	expectStatus(t, rec, http.StatusOK)
	if strings.Contains(rec.Body.String(), "draft") {
		t.Fatalf("an inactive notice must not be published: %s", rec.Body.String())
	}

	expectStatus(t, app.do(t, "POST", "/api/notices", token, map[string]interface{}{"content": "live", "active": true}), http.StatusOK)
	rec = app.do(t, "GET", "/api/notices/active", "", nil)
	if !strings.Contains(rec.Body.String(), "live") {
		t.Fatalf("active notice not returned: %s", rec.Body.String())
	}
}

var tinyPNG, _ = base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==")

func (a *testApp) uploadLogo(t *testing.T, token, filename string, content []byte) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	fw, err := w.CreateFormFile("logo", filename)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fw.Write(content); err != nil {
		t.Fatal(err)
	}
	w.Close()
	req := httptest.NewRequest("POST", "/api/upload/logo", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	a.router.ServeHTTP(rec, req)
	return rec
}

func TestLogoUpload(t *testing.T) {
	app := newTestApp(t, nil)
	token, _ := app.login(t, "admin", testAdminPassword)

	expectStatus(t, app.uploadLogo(t, token, "fake.png", []byte("this is not a png")), http.StatusBadRequest)
	expectStatus(t, app.uploadLogo(t, token, "logo.gif", tinyPNG), http.StatusBadRequest)
	big := append(append([]byte{}, tinyPNG...), make([]byte, 3<<20)...)
	expectStatus(t, app.uploadLogo(t, token, "big.png", big), http.StatusRequestEntityTooLarge)

	rec := app.uploadLogo(t, token, "LOGO.PNG", tinyPNG)
	expectStatus(t, rec, http.StatusOK)
	var first models.Logo
	decode(t, rec, &first)

	svg := []byte(`<?xml version="1.0"?><svg xmlns="http://www.w3.org/2000/svg"><script>alert(1)</script></svg>`)
	rec = app.uploadLogo(t, token, "logo.svg", svg)
	expectStatus(t, rec, http.StatusOK)
	var second models.Logo
	decode(t, rec, &second)

	// 旧 Logo 文件和记录在新 Logo 保存成功后被清理
	if _, err := os.Stat(filepath.Join("uploads", "logos", filepath.Base(first.URL))); !os.IsNotExist(err) {
		t.Errorf("old logo file should be removed, stat err = %v", err)
	}
	var count int64
	config.DB.Model(&models.Logo{}).Count(&count)
	if count != 1 {
		t.Errorf("logo records = %d, want 1", count)
	}

	rec = app.do(t, "GET", second.URL, "", nil)
	expectStatus(t, rec, http.StatusOK)
	if csp := rec.Header().Get("Content-Security-Policy"); !strings.Contains(csp, "sandbox") || !strings.Contains(csp, "default-src 'none'") {
		t.Errorf("uploads must be served with a restrictive CSP, got %q", csp)
	}
	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("uploads must be served with X-Content-Type-Options: nosniff")
	}
}

func TestRenameCascadesAndScopedChecks(t *testing.T) {
	app := newTestApp(t, func(s *config.Settings) { s.Probe.Enabled = false })
	token, _ := app.login(t, "admin", testAdminPassword)

	rec := app.do(t, "POST", "/api/projects", token, map[string]interface{}{"name": "alpha", "label": "Alpha"})
	expectStatus(t, rec, http.StatusOK)
	var project models.Project
	decode(t, rec, &project)
	expectStatus(t, app.do(t, "POST", "/api/projects", token, map[string]interface{}{"name": "alpha", "label": "Dup"}), http.StatusBadRequest)
	expectStatus(t, app.do(t, "POST", "/api/projects", token, map[string]interface{}{"name": "Bad Name", "label": "x"}), http.StatusBadRequest)

	rec = app.do(t, "POST", "/api/environments", token, map[string]interface{}{"name": "qa", "label": "QA", "project_id": project.ID})
	expectStatus(t, rec, http.StatusOK)
	var env models.Environment
	decode(t, rec, &env)
	tool := createTool(t, app, token, map[string]interface{}{"name": "T", "url": "https://t.example.com", "project": "alpha", "environment": "qa"})

	// 有工具时不能删除项目
	expectStatus(t, app.do(t, "DELETE", "/api/projects/"+itoa(project.ID), token, nil), http.StatusBadRequest)

	// 环境改标识：本项目下的工具跟着改
	expectStatus(t, app.do(t, "PUT", "/api/environments/"+itoa(env.ID), token, map[string]interface{}{"name": "uat", "label": "UAT", "project_id": project.ID}), http.StatusOK)
	var got models.Tool
	config.DB.First(&got, tool.ID)
	if got.Environment != "uat" {
		t.Fatalf("tool environment = %q, want uat", got.Environment)
	}

	// 有工具的环境不能移动到其他项目
	expectStatus(t, app.do(t, "PUT", "/api/environments/"+itoa(env.ID), token, map[string]interface{}{"name": "uat", "label": "UAT", "project_id": 1}), http.StatusBadRequest)

	// 项目改标识：环境和工具跟着改
	expectStatus(t, app.do(t, "PUT", "/api/projects/"+itoa(project.ID), token, map[string]interface{}{"name": "beta", "label": "Beta"}), http.StatusOK)
	config.DB.First(&got, tool.ID)
	var gotEnv models.Environment
	config.DB.First(&gotEnv, env.ID)
	if got.Project != "beta" || gotEnv.Project != "beta" {
		t.Fatalf("rename did not cascade: tool.project=%q env.project=%q", got.Project, gotEnv.Project)
	}

	// 删除环境只看本项目：别的项目有 dev 工具，不影响删除这个项目的空 dev 环境
	rec = app.do(t, "POST", "/api/environments", token, map[string]interface{}{"name": "dev", "label": "开发", "project_id": project.ID})
	expectStatus(t, rec, http.StatusOK)
	var emptyDev models.Environment
	decode(t, rec, &emptyDev)
	expectStatus(t, app.do(t, "DELETE", "/api/environments/"+itoa(emptyDev.ID), token, nil), http.StatusOK)

	// 有工具的环境不能删
	expectStatus(t, app.do(t, "DELETE", "/api/environments/"+itoa(env.ID), token, nil), http.StatusBadRequest)
}

func TestCORS(t *testing.T) {
	app := newTestApp(t, nil)
	rec := app.do(t, "GET", "/api/sites", "", nil, "Origin", "https://evil.example.com")
	expectStatus(t, rec, http.StatusOK)
	if v := rec.Header().Get("Access-Control-Allow-Origin"); v != "" {
		t.Fatalf("CORS should be disabled by default, got Access-Control-Allow-Origin=%q", v)
	}

	app = newTestApp(t, func(s *config.Settings) { s.CORSAllowOrigins = []string{"https://portal.example.com"} })
	rec = app.do(t, "GET", "/api/sites", "", nil, "Origin", "https://portal.example.com")
	if v := rec.Header().Get("Access-Control-Allow-Origin"); v != "https://portal.example.com" {
		t.Fatalf("allowed origin not echoed, got %q", v)
	}
	if v := rec.Header().Get("Access-Control-Allow-Credentials"); v != "" {
		t.Errorf("credentials should not be allowed, got %q", v)
	}
	rec = app.do(t, "GET", "/api/sites", "", nil, "Origin", "https://evil.example.com")
	if v := rec.Header().Get("Access-Control-Allow-Origin"); v != "" {
		t.Fatalf("other origins must not be allowed, got %q", v)
	}
}

func TestProbeStatusMetricsAndServiceDiscovery(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/broken" {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	app := newTestApp(t, func(s *config.Settings) { s.MetricsToken = "metrics-secret" })
	token, _ := app.login(t, "admin", testAdminPassword)

	// 去掉示例数据，避免测试去请求公网
	config.DB.Where("1 = 1").Delete(&models.Tool{})
	up := createTool(t, app, token, map[string]interface{}{"name": "Up", "url": srv.URL + "/ok", "project": "cloud", "environment": "prod", "category": "监控"})
	down := createTool(t, app, token, map[string]interface{}{"name": "Down", "url": srv.URL + "/broken", "project": "monitor", "environment": "dev"})
	off := createTool(t, app, token, map[string]interface{}{"name": "Off", "url": srv.URL + "/off", "project": "cloud", "environment": "prod", "probe_disabled": true})

	app.prober.RunOnce(context.Background())

	// 首页状态接口
	rec := app.do(t, "GET", "/api/probe/status", "", nil)
	expectStatus(t, rec, http.StatusOK)
	var status struct {
		Enabled bool `json:"enabled"`
		Results map[string]struct {
			State      string `json:"state"`
			StatusCode int    `json:"status_code"`
			Reason     string `json:"reason"`
		} `json:"results"`
	}
	decode(t, rec, &status)
	if !status.Enabled || status.Results[itoa(up.ID)].State != "up" || status.Results[itoa(down.ID)].State != "down" ||
		status.Results[itoa(down.ID)].StatusCode != 500 || status.Results[itoa(down.ID)].Reason != probe.ReasonHTTPStatus {
		t.Fatalf("unexpected probe status: %s", rec.Body.String())
	}
	if _, ok := status.Results[itoa(off.ID)]; ok {
		t.Fatal("tools with probing disabled must not appear in status")
	}

	// 立即检测需要登录
	expectStatus(t, app.do(t, "POST", "/api/probe/check/"+itoa(up.ID), "", nil), http.StatusUnauthorized)
	rec = app.do(t, "POST", "/api/probe/check/"+itoa(up.ID), token, nil)
	expectStatus(t, rec, http.StatusOK)
	if !strings.Contains(rec.Body.String(), `"state":"up"`) {
		t.Fatalf("check result: %s", rec.Body.String())
	}

	// /metrics 需要 METRICS_TOKEN
	expectStatus(t, app.do(t, "GET", "/metrics", "", nil), http.StatusUnauthorized)
	expectStatus(t, app.do(t, "GET", "/metrics", "wrong", nil), http.StatusUnauthorized)
	rec = app.do(t, "GET", "/metrics", "metrics-secret", nil)
	expectStatus(t, rec, http.StatusOK)
	if ct := rec.Header().Get("Content-Type"); ct != probe.ContentType {
		t.Errorf("metrics content type = %q", ct)
	}
	body := rec.Body.String()
	for _, want := range []string{
		"opsportal_tools 3\n",
		`opsportal_probe_success{tool_id="` + itoa(up.ID) + `",tool="Up",project="cloud",environment="prod",category="监控",url="` + srv.URL + `/ok"} 1`,
		`opsportal_probe_success{tool_id="` + itoa(down.ID) + `",tool="Down",project="monitor",environment="dev",category="",url="` + srv.URL + `/broken"} 0`,
		`opsportal_probe_http_status_code{tool_id="` + itoa(down.ID) + `"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("metrics missing %q\n%s", want, body)
		}
	}
	if strings.Contains(body, `tool="Off"`) {
		t.Error("tools with probing disabled must not be exported")
	}

	// http_sd 服务发现
	expectStatus(t, app.do(t, "GET", "/prometheus/targets", "", nil), http.StatusUnauthorized)
	rec = app.do(t, "GET", "/prometheus/targets", "metrics-secret", nil)
	expectStatus(t, rec, http.StatusOK)
	var groups []struct {
		Targets []string          `json:"targets"`
		Labels  map[string]string `json:"labels"`
	}
	decode(t, rec, &groups)
	if len(groups) != 2 || groups[0].Targets[0] != srv.URL+"/ok" || groups[0].Labels["project"] != "cloud" ||
		groups[0].Labels["category"] != "监控" || groups[0].Labels["tool_id"] != itoa(up.ID) {
		t.Fatalf("unexpected targets: %s", rec.Body.String())
	}
	if _, ok := groups[1].Labels["category"]; ok {
		t.Error("empty category should be omitted from labels")
	}

	rec = app.do(t, "GET", "/prometheus/targets?project=monitor&environment=dev", "metrics-secret", nil)
	decode(t, rec, &groups)
	if len(groups) != 1 || groups[0].Labels["tool"] != "Down" {
		t.Fatalf("filtered targets: %s", rec.Body.String())
	}
	rec = app.do(t, "GET", "/prometheus/targets?project=none", "metrics-secret", nil)
	if strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Fatalf("empty result must be a JSON array, got %s", rec.Body.String())
	}
}

func TestProbeDisabled(t *testing.T) {
	app := newTestApp(t, func(s *config.Settings) { s.Probe.Enabled = false })
	rec := app.do(t, "GET", "/api/probe/status", "", nil)
	expectStatus(t, rec, http.StatusOK)
	if !strings.Contains(rec.Body.String(), `"enabled":false`) {
		t.Fatalf("status when disabled: %s", rec.Body.String())
	}
	rec = app.do(t, "GET", "/metrics", "", nil)
	expectStatus(t, rec, http.StatusOK)
	if !strings.Contains(rec.Body.String(), "opsportal_probe_enabled 0\n") {
		t.Fatalf("metrics when disabled: %s", rec.Body.String())
	}
}

func itoa(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}
