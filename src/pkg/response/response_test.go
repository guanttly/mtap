package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	bizErr "github.com/euler/mtap/pkg/errors"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestOK(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	OK(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var r R
	_ = json.Unmarshal(w.Body.Bytes(), &r)
	assert.Equal(t, 0, r.Code)
	assert.Equal(t, "success", r.Message)
}

func TestOKWithData(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	OKWithData(c, map[string]string{"id": "123"})

	assert.Equal(t, http.StatusOK, w.Code)
	var r R
	_ = json.Unmarshal(w.Body.Bytes(), &r)
	assert.Equal(t, 0, r.Code)
	data := r.Data.(map[string]interface{})
	assert.Equal(t, "123", data["id"])
}

func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	Created(c, map[string]string{"id": "new-1"})

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestFail(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	Fail(c, http.StatusBadRequest, bizErr.ErrInvalidParam, "name is required")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var r R
	_ = json.Unmarshal(w.Body.Bytes(), &r)
	assert.Equal(t, int(bizErr.ErrInvalidParam), r.Code)
}

func TestFailWithError_BizError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	err := bizErr.New(bizErr.ErrNotFound)
	FailWithError(c, err)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestFailWithError_StdError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	FailWithError(c, assert.AnError)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestNewPageResult(t *testing.T) {
	items := []string{"a", "b", "c"}
	pr := NewPageResult(items, 25, 1, 10)
	assert.Equal(t, 3, len(pr.Items))
	assert.Equal(t, int64(25), pr.Total)
	assert.Equal(t, 3, pr.TotalPages)
}

func TestNewPageResult_ExactDivision(t *testing.T) {
	pr := NewPageResult([]int{}, 20, 1, 10)
	assert.Equal(t, 2, pr.TotalPages)
}
