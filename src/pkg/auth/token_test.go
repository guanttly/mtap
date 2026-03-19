package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestManager() *JWTManager {
	return NewJWTManager("test-secret-key-32chars-long!!", 2*time.Hour, 7*24*time.Hour)
}

func TestGenerateTokenPair(t *testing.T) {
	m := newTestManager()
	pair, err := m.GenerateTokenPair("user-1", "zhangsan", "admin", "dept-1")

	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
	assert.True(t, pair.ExpiresAt > time.Now().Unix())
}

func TestValidateToken_Valid(t *testing.T) {
	m := newTestManager()
	pair, _ := m.GenerateTokenPair("user-1", "zhangsan", "admin", "dept-1")

	claims, err := m.ValidateToken(pair.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, "user-1", claims.UserID)
	assert.Equal(t, "zhangsan", claims.Username)
	assert.Equal(t, "admin", claims.Role)
	assert.Equal(t, "dept-1", claims.DepartmentID)
}

func TestValidateToken_Expired(t *testing.T) {
	m := NewJWTManager("secret", 1*time.Millisecond, 7*24*time.Hour)
	pair, _ := m.GenerateTokenPair("user-1", "test", "admin", "")

	time.Sleep(10 * time.Millisecond)

	_, err := m.ValidateToken(pair.AccessToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无效或已过期")
}

func TestValidateToken_InvalidString(t *testing.T) {
	m := newTestManager()
	_, err := m.ValidateToken("this-is-not-a-jwt-token")
	assert.Error(t, err)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	m1 := NewJWTManager("secret-1", 2*time.Hour, 168*time.Hour)
	m2 := NewJWTManager("secret-2", 2*time.Hour, 168*time.Hour)

	pair, _ := m1.GenerateTokenPair("user-1", "test", "admin", "")
	_, err := m2.ValidateToken(pair.AccessToken)
	assert.Error(t, err)
}

func TestRefreshAccessToken(t *testing.T) {
	m := newTestManager()
	pair, _ := m.GenerateTokenPair("user-1", "zhangsan", "admin", "dept-1")

	newPair, err := m.RefreshAccessToken(pair.RefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, newPair.AccessToken)
	assert.NotEqual(t, pair.AccessToken, newPair.AccessToken)
}

func TestRefreshAccessToken_InvalidRefresh(t *testing.T) {
	m := newTestManager()
	_, err := m.RefreshAccessToken("invalid-refresh-token")
	assert.Error(t, err)
}

func TestClaims_Issuer(t *testing.T) {
	m := newTestManager()
	pair, _ := m.GenerateTokenPair("user-1", "test", "admin", "")
	claims, _ := m.ValidateToken(pair.AccessToken)
	assert.Equal(t, "mtap", claims.Issuer)
}
