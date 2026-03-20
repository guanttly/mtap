package auth

import (
"bytes"
"encoding/json"
"net/http"
"net/http/httptest"
"testing"
"time"

"github.com/DATA-DOG/go-sqlmock"
"github.com/gin-gonic/gin"
"golang.org/x/crypto/bcrypt"
"gorm.io/driver/mysql"
"gorm.io/gorm"
"gorm.io/gorm/logger"

pkgAuth "github.com/euler/mtap/pkg/auth"
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

func newTestJWT() *pkgAuth.JWTManager {
return pkgAuth.NewJWTManager("test-secret-key", time.Hour, 24*time.Hour)
}

func setupAuthRouter(db *gorm.DB, jwtMgr *pkgAuth.JWTManager) *gin.Engine {
h := NewHandler(db, jwtMgr)
r := gin.New()
v1 := r.Group("/api/v1")
h.RegisterRoutes(v1)
return r
}

func doAuthRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
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

func parseAuthBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
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

func TestAuthHandler_Login_BadRequest(t *testing.T) {
db, mock := newMockDB(t)
defer func() { _ = mock.ExpectationsWereMet() }()

r := setupAuthRouter(db, newTestJWT())
w := doAuthRequest(r, http.MethodPost, "/api/v1/auth/login", map[string]interface{}{
"username": "admin",
// missing password -> ShouldBindJSON returns error before any DB query
})
if w.Code != http.StatusBadRequest {
t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
}
}

func TestAuthHandler_Login_UserNotFound(t *testing.T) {
db, mock := newMockDB(t)
defer func() { _ = mock.ExpectationsWereMet() }()

// SELECT * FROM users WHERE username=? AND status=? LIMIT 1 -> no rows
mock.ExpectQuery(`SELECT \* FROM .users.`).
WillReturnError(gorm.ErrRecordNotFound)

r := setupAuthRouter(db, newTestJWT())
w := doAuthRequest(r, http.MethodPost, "/api/v1/auth/login", map[string]interface{}{
"username": "nobody",
"password": "anything",
})
if w.Code != http.StatusUnauthorized {
t.Fatalf("expected 401, got %d: %s", w.Code, w.Body.String())
}
}

func TestAuthHandler_Login_WrongPassword(t *testing.T) {
db, mock := newMockDB(t)
defer func() { _ = mock.ExpectationsWereMet() }()

hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
now := time.Now()
rows := sqlmock.NewRows([]string{
"id", "username", "password_hash", "real_name", "role_id",
"department_id", "status", "last_login_at", "created_at", "updated_at",
}).AddRow("u-001", "admin", string(hash), "Admin", "role-001", "", "active", now, now, now)

mock.ExpectQuery(`SELECT \* FROM .users.`).WillReturnRows(rows)

r := setupAuthRouter(db, newTestJWT())
w := doAuthRequest(r, http.MethodPost, "/api/v1/auth/login", map[string]interface{}{
"username": "admin",
"password": "wrongpass",
})
if w.Code != http.StatusUnauthorized {
t.Fatalf("expected 401, got %d: %s", w.Code, w.Body.String())
}
}

func TestAuthHandler_Login_OK(t *testing.T) {
db, mock := newMockDB(t)
defer func() { _ = mock.ExpectationsWereMet() }()

hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
now := time.Now()
userRows := sqlmock.NewRows([]string{
"id", "username", "password_hash", "real_name", "role_id",
"department_id", "status", "last_login_at", "created_at", "updated_at",
}).AddRow("u-001", "admin", string(hash), "Admin", "role-001", "", "active", now, now, now)

roleRows := sqlmock.NewRows([]string{
"id", "name", "permissions", "is_preset", "created_at", "updated_at",
}).AddRow("role-001", "admin", "[]", true, now, now)

mock.ExpectQuery(`SELECT \* FROM .users.`).WillReturnRows(userRows)
mock.ExpectQuery(`SELECT \* FROM .roles.`).WillReturnRows(roleRows)
mock.ExpectExec(`UPDATE .users.`).WillReturnResult(sqlmock.NewResult(1, 1))

r := setupAuthRouter(db, newTestJWT())
w := doAuthRequest(r, http.MethodPost, "/api/v1/auth/login", map[string]interface{}{
"username": "admin",
"password": "secret123",
})
if w.Code != http.StatusOK {
t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
}
body := parseAuthBody(t, w)
data := body["data"].(map[string]interface{})
if data["access_token"] == nil || data["access_token"] == "" {
t.Errorf("expected access_token, got %v", data)
}
}

func TestAuthHandler_Logout_OK(t *testing.T) {
db, mock := newMockDB(t)
defer func() { _ = mock.ExpectationsWereMet() }()

r := setupAuthRouter(db, newTestJWT())
w := doAuthRequest(r, http.MethodPost, "/api/v1/auth/logout", nil)
if w.Code != http.StatusOK {
t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
}
}
