// Package rule HTTP handler 集成测试
package rule

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appRule "github.com/euler/mtap/internal/application/rule"
	domain "github.com/euler/mtap/internal/domain/rule"
)

// ─── Mock Repositories ────────────────────────────────────────────────────────

type htMockConflictRuleRepo struct {
	rules map[string]*domain.ConflictRule
}

func newHtConflictRuleRepo() *htMockConflictRuleRepo {
	return &htMockConflictRuleRepo{rules: make(map[string]*domain.ConflictRule)}
}
func (m *htMockConflictRuleRepo) Save(_ context.Context, r *domain.ConflictRule) error {
	m.rules[r.ID] = r
	return nil
}
func (m *htMockConflictRuleRepo) FindByID(_ context.Context, id string) (*domain.ConflictRule, error) {
	return m.rules[id], nil
}
func (m *htMockConflictRuleRepo) FindByItemPair(_ context.Context, a, b string) (*domain.ConflictRule, error) {
	if a > b {
		a, b = b, a
	}
	for _, r := range m.rules {
		if r.ItemAID == a && r.ItemBID == b {
			return r, nil
		}
	}
	return nil, nil
}
func (m *htMockConflictRuleRepo) FindAll(_ context.Context, status domain.RuleStatus) ([]*domain.ConflictRule, error) {
	var out []*domain.ConflictRule
	for _, r := range m.rules {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out, nil
}
func (m *htMockConflictRuleRepo) Update(_ context.Context, r *domain.ConflictRule) error {
	m.rules[r.ID] = r
	return nil
}
func (m *htMockConflictRuleRepo) Delete(_ context.Context, id string) error {
	delete(m.rules, id)
	return nil
}

type htMockConflictPkgRepo struct {
	pkgs map[string]*domain.ConflictPackage
}

func newHtConflictPkgRepo() *htMockConflictPkgRepo {
	return &htMockConflictPkgRepo{pkgs: make(map[string]*domain.ConflictPackage)}
}
func (m *htMockConflictPkgRepo) Save(_ context.Context, p *domain.ConflictPackage) error {
	m.pkgs[p.ID] = p
	return nil
}
func (m *htMockConflictPkgRepo) FindByID(_ context.Context, id string) (*domain.ConflictPackage, error) {
	return m.pkgs[id], nil
}
func (m *htMockConflictPkgRepo) FindByName(_ context.Context, name string) (*domain.ConflictPackage, error) {
	for _, p := range m.pkgs {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, nil
}
func (m *htMockConflictPkgRepo) FindAll(_ context.Context) ([]*domain.ConflictPackage, error) {
	var out []*domain.ConflictPackage
	for _, p := range m.pkgs {
		out = append(out, p)
	}
	return out, nil
}
func (m *htMockConflictPkgRepo) Update(_ context.Context, p *domain.ConflictPackage) error {
	m.pkgs[p.ID] = p
	return nil
}
func (m *htMockConflictPkgRepo) Delete(_ context.Context, id string) error {
	delete(m.pkgs, id)
	return nil
}

type htMockDepRuleRepo struct {
	rules map[string]*domain.DependencyRule
}

func newHtDepRuleRepo() *htMockDepRuleRepo {
	return &htMockDepRuleRepo{rules: make(map[string]*domain.DependencyRule)}
}
func (m *htMockDepRuleRepo) Save(_ context.Context, r *domain.DependencyRule) error {
	m.rules[r.ID] = r
	return nil
}
func (m *htMockDepRuleRepo) FindByID(_ context.Context, id string) (*domain.DependencyRule, error) {
	return m.rules[id], nil
}
func (m *htMockDepRuleRepo) FindByPostItem(_ context.Context, postItemID string) ([]*domain.DependencyRule, error) {
	var out []*domain.DependencyRule
	for _, r := range m.rules {
		if r.PostItemID == postItemID {
			out = append(out, r)
		}
	}
	return out, nil
}
func (m *htMockDepRuleRepo) FindAll(_ context.Context, status domain.RuleStatus) ([]*domain.DependencyRule, error) {
	var out []*domain.DependencyRule
	for _, r := range m.rules {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out, nil
}
func (m *htMockDepRuleRepo) Update(_ context.Context, r *domain.DependencyRule) error {
	m.rules[r.ID] = r
	return nil
}
func (m *htMockDepRuleRepo) Delete(_ context.Context, id string) error {
	delete(m.rules, id)
	return nil
}

type htMockTagRepo struct {
	tags map[string]*domain.PriorityTag
}

func newHtTagRepo() *htMockTagRepo { return &htMockTagRepo{tags: make(map[string]*domain.PriorityTag)} }
func (m *htMockTagRepo) Save(_ context.Context, t *domain.PriorityTag) error {
	m.tags[t.ID] = t
	return nil
}
func (m *htMockTagRepo) FindByID(_ context.Context, id string) (*domain.PriorityTag, error) {
	return m.tags[id], nil
}
func (m *htMockTagRepo) FindByName(_ context.Context, name string) (*domain.PriorityTag, error) {
	for _, t := range m.tags {
		if t.Name == name {
			return t, nil
		}
	}
	return nil, nil
}
func (m *htMockTagRepo) FindAll(_ context.Context) ([]*domain.PriorityTag, error) {
	var out []*domain.PriorityTag
	for _, t := range m.tags {
		out = append(out, t)
	}
	return out, nil
}
func (m *htMockTagRepo) Update(_ context.Context, t *domain.PriorityTag) error {
	m.tags[t.ID] = t
	return nil
}
func (m *htMockTagRepo) Delete(_ context.Context, id string) error { delete(m.tags, id); return nil }

type htMockSortingRepo struct{ items []*domain.SortingStrategy }

func (m *htMockSortingRepo) Save(_ context.Context, s *domain.SortingStrategy) error {
	m.items = append(m.items, s)
	return nil
}
func (m *htMockSortingRepo) FindByID(_ context.Context, _ string) (*domain.SortingStrategy, error) {
	return nil, nil
}
func (m *htMockSortingRepo) FindAll(_ context.Context) ([]*domain.SortingStrategy, error) {
	return m.items, nil
}
func (m *htMockSortingRepo) FindByScope(_ context.Context, _ domain.EffectiveScope) ([]*domain.SortingStrategy, error) {
	return m.items, nil
}
func (m *htMockSortingRepo) Update(_ context.Context, _ *domain.SortingStrategy) error { return nil }
func (m *htMockSortingRepo) Delete(_ context.Context, _ string) error                  { return nil }

type htMockAdaptRepo struct{ rules []*domain.PatientAdaptRule }

func (m *htMockAdaptRepo) SaveAll(_ context.Context, rules []*domain.PatientAdaptRule) error {
	m.rules = rules
	return nil
}
func (m *htMockAdaptRepo) FindAll(_ context.Context) ([]*domain.PatientAdaptRule, error) {
	return m.rules, nil
}
func (m *htMockAdaptRepo) DeleteAll(_ context.Context) error { m.rules = nil; return nil }

type htMockSourceRepo struct{ controls []*domain.SourceControl }

func (m *htMockSourceRepo) SaveAll(_ context.Context, ctrls []*domain.SourceControl) error {
	m.controls = ctrls
	return nil
}
func (m *htMockSourceRepo) FindAll(_ context.Context) ([]*domain.SourceControl, error) {
	return m.controls, nil
}
func (m *htMockSourceRepo) DeleteAll(_ context.Context) error { m.controls = nil; return nil }

type htMockPreItemChecker struct{}

func (m *htMockPreItemChecker) GetCompletedTime(_ context.Context, _, _ string) (*time.Time, error) {
	return nil, nil
}

// ─── 辅助函数 ─────────────────────────────────────────────────────────────────

func newTestRuleHandler() *Handler {
	crRepo := newHtConflictRuleRepo()
	cpRepo := newHtConflictPkgRepo()
	drRepo := newHtDepRuleRepo()
	tRepo := newHtTagRepo()
	sortRepo := &htMockSortingRepo{}
	adaptRepo := &htMockAdaptRepo{}
	sourceRepo := &htMockSourceRepo{}
	conflictSvc := domain.NewConflictDetectionService(crRepo, cpRepo)
	depSvc := domain.NewDependencyValidationService(drRepo, &htMockPreItemChecker{})
	checker := domain.NewCircularDependencyChecker(drRepo)
	svc := appRule.NewRuleAppService(
		crRepo, cpRepo, drRepo, tRepo, sortRepo, adaptRepo, sourceRepo,
		conflictSvc, depSvc, checker,
	)
	return NewHandler(svc)
}

// setupRuleRouter 创建绕过鉴权中间件的 gin 路由
func setupRuleRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-admin")
		c.Set("role", "admin")
		c.Next()
	})
	api := r.Group("/api/v1")
	h.RegisterRoutes(api)
	return r
}

// doRuleRequest 发送 HTTP 测试请求
func doRuleRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// parseRuleBody 解析 JSON 响应体
func parseRuleBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&m))
	return m
}

// ─── 测试用例 ─────────────────────────────────────────────────────────────────

func TestRuleHandler_CreateConflictRule_OK(t *testing.T) {
	h := newTestRuleHandler()
	r := setupRuleRouter(h)

	w := doRuleRequest(r, http.MethodPost, "/api/v1/rules/conflicts", map[string]interface{}{
		"item_a_id":    "E1",
		"item_b_id":    "E2",
		"min_interval": 24,
		"level":        "warning",
	})
	assert.Equal(t, http.StatusCreated, w.Code)
	body := parseRuleBody(t, w)
	assert.Equal(t, float64(0), body["code"])
	data := body["data"].(map[string]interface{})
	assert.Equal(t, "E1", data["item_a_id"])
}

func TestRuleHandler_CreateConflictRule_BadRequest(t *testing.T) {
	h := newTestRuleHandler()
	r := setupRuleRouter(h)

	// item_b_id 缺失 → 触发 binding 验证失败
	w := doRuleRequest(r, http.MethodPost, "/api/v1/rules/conflicts", map[string]interface{}{
		"item_a_id": "E1",
		"level":     "warning",
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRuleHandler_GetConflictRule_OK(t *testing.T) {
	h := newTestRuleHandler()
	r := setupRuleRouter(h)

	w1 := doRuleRequest(r, http.MethodPost, "/api/v1/rules/conflicts", map[string]interface{}{
		"item_a_id": "E3", "item_b_id": "E4", "min_interval": 12, "level": "warning",
	})
	require.Equal(t, http.StatusCreated, w1.Code)
	id := parseRuleBody(t, w1)["data"].(map[string]interface{})["id"].(string)

	w2 := doRuleRequest(r, http.MethodGet, "/api/v1/rules/conflicts/"+id, nil)
	assert.Equal(t, http.StatusOK, w2.Code)
	body := parseRuleBody(t, w2)
	assert.Equal(t, float64(0), body["code"])
}

func TestRuleHandler_ListConflictRules_OK(t *testing.T) {
	h := newTestRuleHandler()
	r := setupRuleRouter(h)

	doRuleRequest(r, http.MethodPost, "/api/v1/rules/conflicts", map[string]interface{}{
		"item_a_id": "E5", "item_b_id": "E6", "min_interval": 6, "level": "warning",
	})

	// page 和 page_size 不传时 binding:"gte=1" 会失败，传上默认值
	w := doRuleRequest(r, http.MethodGet, "/api/v1/rules/conflicts?page=1&page_size=20", nil)
	assert.Equal(t, http.StatusOK, w.Code)
	body := parseRuleBody(t, w)
	assert.Equal(t, float64(0), body["code"])
}

func TestRuleHandler_DeleteConflictRule_OK(t *testing.T) {
	h := newTestRuleHandler()
	r := setupRuleRouter(h)

	w1 := doRuleRequest(r, http.MethodPost, "/api/v1/rules/conflicts", map[string]interface{}{
		"item_a_id": "E7", "item_b_id": "E8", "min_interval": 48, "level": "forbid",
	})
	require.Equal(t, http.StatusCreated, w1.Code)
	id := parseRuleBody(t, w1)["data"].(map[string]interface{})["id"].(string)

	w2 := doRuleRequest(r, http.MethodDelete, "/api/v1/rules/conflicts/"+id, nil)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestRuleHandler_CreateConflictPackage_OK(t *testing.T) {
	h := newTestRuleHandler()
	r := setupRuleRouter(h)

	w := doRuleRequest(r, http.MethodPost, "/api/v1/rules/conflict-packages", map[string]interface{}{
		"name":         "基础生化套餐",
		"item_ids":     []string{"E1", "E2"},
		"min_interval": 24,
		"level":        "forbid",
	})
	assert.Equal(t, http.StatusCreated, w.Code)
	body := parseRuleBody(t, w)
	assert.Equal(t, float64(0), body["code"])
	data := body["data"].(map[string]interface{})
	assert.Equal(t, "基础生化套餐", data["name"])
}

func TestRuleHandler_ListConflictPackages_OK(t *testing.T) {
	h := newTestRuleHandler()
	r := setupRuleRouter(h)

	doRuleRequest(r, http.MethodPost, "/api/v1/rules/conflict-packages", map[string]interface{}{
		"name": "套餐A", "item_ids": []string{"E1", "E2"}, "level": "forbid",
	})

	w := doRuleRequest(r, http.MethodGet, "/api/v1/rules/conflict-packages", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRuleHandler_CreateDependencyRule_OK(t *testing.T) {
	h := newTestRuleHandler()
	r := setupRuleRouter(h)

	w := doRuleRequest(r, http.MethodPost, "/api/v1/rules/dependencies", map[string]interface{}{
		"pre_item_id":    "E1",
		"post_item_id":   "E2",
		"type":           "mandatory",
		"validity_hours": 48,
	})
	assert.Equal(t, http.StatusCreated, w.Code)
	body := parseRuleBody(t, w)
	assert.Equal(t, float64(0), body["code"])
}

func TestRuleHandler_ListDependencyRules_OK(t *testing.T) {
	h := newTestRuleHandler()
	r := setupRuleRouter(h)

	w := doRuleRequest(r, http.MethodGet, "/api/v1/rules/dependencies?page=1&page_size=10", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRuleHandler_CreatePriorityTag_OK(t *testing.T) {
	h := newTestRuleHandler()
	r := setupRuleRouter(h)

	w := doRuleRequest(r, http.MethodPost, "/api/v1/rules/priority-tags", map[string]interface{}{
		"name":   "紧急",
		"weight": 90,
		"color":  "#FF0000",
	})
	assert.Equal(t, http.StatusCreated, w.Code)
	body := parseRuleBody(t, w)
	assert.Equal(t, float64(0), body["code"])
	data := body["data"].(map[string]interface{})
	assert.Equal(t, "紧急", data["name"])
}

func TestRuleHandler_ListPriorityTags_OK(t *testing.T) {
	h := newTestRuleHandler()
	r := setupRuleRouter(h)

	w := doRuleRequest(r, http.MethodGet, "/api/v1/rules/priority-tags", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRuleHandler_DeletePriorityTag_OK(t *testing.T) {
	h := newTestRuleHandler()
	r := setupRuleRouter(h)

	w1 := doRuleRequest(r, http.MethodPost, "/api/v1/rules/priority-tags", map[string]interface{}{
		"name": "VIP客户", "weight": 80, "color": "#00FF00",
	})
	require.Equal(t, http.StatusCreated, w1.Code)
	id := parseRuleBody(t, w1)["data"].(map[string]interface{})["id"].(string)

	w2 := doRuleRequest(r, http.MethodDelete, "/api/v1/rules/priority-tags/"+id, nil)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestRuleHandler_SaveAndListPatientAdaptRules_OK(t *testing.T) {
	h := newTestRuleHandler()
	r := setupRuleRouter(h)

	w := doRuleRequest(r, http.MethodPost, "/api/v1/rules/patient-adapt", []map[string]interface{}{
		{
			"condition_type":  "age",
			"condition_value": "60",
			"action":          "filter_device",
			"action_params":   map[string]string{"device_type": "CT"},
		},
	})
	assert.Equal(t, http.StatusOK, w.Code)

	w2 := doRuleRequest(r, http.MethodGet, "/api/v1/rules/patient-adapt", nil)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestRuleHandler_SaveAndListSourceControls_OK(t *testing.T) {
	h := newTestRuleHandler()
	r := setupRuleRouter(h)

	w := doRuleRequest(r, http.MethodPost, "/api/v1/rules/source-controls", []map[string]interface{}{
		{
			"source_type":      "outpatient",
			"slot_pool_id":     "pool-001",
			"allocation_ratio": 0.8,
		},
	})
	assert.Equal(t, http.StatusOK, w.Code)

	w2 := doRuleRequest(r, http.MethodGet, "/api/v1/rules/source-controls", nil)
	assert.Equal(t, http.StatusOK, w2.Code)
}
