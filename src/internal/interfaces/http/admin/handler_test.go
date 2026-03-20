package admin

import (
"bytes"
"encoding/json"
"net/http"
"net/http/httptest"
"testing"
"time"

sqlmock "github.com/DATA-DOG/go-sqlmock"
"github.com/gin-gonic/gin"
"gorm.io/driver/mysql"
"gorm.io/gorm"
"gorm.io/gorm/logger"
)

func init() { gin.SetMode(gin.TestMode) }

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
t.Helper()
db, mock, err := sqlmock.New()
if err != nil {
t.Fatalf("sqlmock.New: %v", err)
}
gormDB, err := gorm.Open(mysql.New(mysql.Config{
Conn:                      db,
SkipInitializeWithVersion: true,
}), &gorm.Config{Logger: logger.Discard})
if err != nil {
t.Fatalf("gorm.Open: %v", err)
}
return gormDB, mock
}

func setupAdminRouter(db *gorm.DB) *gin.Engine {
h := NewHandler(db)
r := gin.New()
r.Use(func(c *gin.Context) {
c.Set("user_id", "admin-001")
c.Set("role", "admin")
c.Next()
})
v1 := r.Group("/api/v1")
h.RegisterRoutes(v1)
return r
}

func doAdminRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
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

func parseAdminBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
t.Helper()
var m map[string]interface{}
if err := json.NewDecoder(w.Body).Decode(&m); err != nil {
t.Fatalf("decode body: %v", err)
}
return m
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestAdminHandler_ListUsers_OK(t *testing.T) {
db, mock := newMockDB(t)
defer func() { _ = mock.ExpectationsWereMet() }()

now := time.Now()
// COUNT query from GORM paginator
mock.ExpectQuery("SELECT count").
WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
// SELECT users
mock.ExpectQuery("SELECT").
WillReturnRows(sqlmock.NewRows([]string{
"id", "username", "password_hash", "real_name", "role_id",
"department_id", "status", "last_login_at", "created_at", "updated_at",
}).AddRow("u-001", "admin", "hash", "Admin", "role-001", "", "active", now, now, now))
// SELECT roles for role name mapping
mock.ExpectQuery("SELECT").
WillReturnRows(sqlmock.NewRows([]string{
"id", "name", "permissions", "is_preset", "created_at", "updated_at",
}).AddRow("role-001", "admin", "[]", true, now, now))

r := setupAdminRouter(db)
w := doAdminRequest(r, http.MethodGet, "/api/v1/admin/users", nil)
if w.Code != http.StatusOK {
t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
}
}

func TestAdminHandler_CreateUser_BadRequest(t *testing.T) {
db, mock := newMockDB(t)
defer func() { _ = mock.ExpectationsWereMet() }()

r := setupAdminRouter(db)
// missing required: password, role_id
w := doAdminRequest(r, http.MethodPost, "/api/v1/admin/users", map[string]interface{}{
"username": "newuser",
})
if w.Code != http.StatusBadRequest {
t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
}
}

func TestAdminHandler_CreateUser_OK(t *testing.T) {
db, mock := newMockDB(t)
defer func() { _ = mock.ExpectationsWereMet() }()

now := time.Now()
// duplicate check count=0
mock.ExpectQuery("SELECT count").
WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
// GORM Create wraps in transaction: BEGIN -> INSERT -> COMMIT
mock.ExpectBegin()
mock.ExpectExec("INSERT INTO").
WillReturnResult(sqlmock.NewResult(1, 1))
mock.ExpectCommit()
// SELECT role for response
mock.ExpectQuery("SELECT").
WillReturnRows(sqlmock.NewRows([]string{
"id", "name", "permissions", "is_preset", "created_at", "updated_at",
}).AddRow("role-001", "operator", "[]", false, now, now))

r := setupAdminRouter(db)
w := doAdminRequest(r, http.MethodPost, "/api/v1/admin/users", map[string]interface{}{
"username":  "operator1",
"password":  "pass123456",
"real_name": "Operator One",
"role_id":   "role-001",
})
if w.Code != http.StatusOK {
t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
}
body := parseAdminBody(t, w)
data := body["data"].(map[string]interface{})
if data["username"] != "operator1" {
t.Errorf("unexpected username: %v", data["username"])
}
}

func TestAdminHandler_ListRoles_OK(t *testing.T) {
db, mock := newMockDB(t)
defer func() { _ = mock.ExpectationsWereMet() }()

now := time.Now()
mock.ExpectQuery("SELECT").
WillReturnRows(sqlmock.NewRows([]string{
"id", "name", "permissions", "is_preset", "created_at", "updated_at",
}).AddRow("role-001", "admin", "[]", true, now, now))

r := setupAdminRouter(db)
w := doAdminRequest(r, http.MethodGet, "/api/v1/admin/roles", nil)
if w.Code != http.StatusOK {
t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
}
body := parseAdminBody(t, w)
// ListRoles returns {"items": [...], "total": N}
data := body["data"].(map[string]interface{})
items, ok := data["items"].([]interface{})
if !ok || len(items) == 0 {
t.Errorf("expected non-empty role items, got %v", data)
}
}

func TestAdminHandler_CreateRole_BadRequest(t *testing.T) {
db, mock := newMockDB(t)
defer func() { _ = mock.ExpectationsWereMet() }()

r := setupAdminRouter(db)
w := doAdminRequest(r, http.MethodPost, "/api/v1/admin/roles", map[string]interface{}{})
if w.Code != http.StatusBadRequest {
t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
}
}
