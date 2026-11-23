package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"jing-sync/boot/logger"
	"jing-sync/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
	logger.LoggerInit()
}

func generateTestToken(t *testing.T, userId uint, username string) string {
	sk, err := utils.GetSecretKey()
	require.NoError(t, err)

	claims := &utils.Claims{
		UserId:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(sk))
	require.NoError(t, err)
	return tokenStr
}

func TestAuthMiddleware(t *testing.T) {
	t.Cleanup(func() {
		_ = os.RemoveAll("data")
	})

	r := gin.New()
	r.Use(I18nMiddleware())
	r.Use(AuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no claims"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":  "ok",
			"user_id":  claims.(*utils.Claims).UserId,
			"username": claims.(*utils.Claims).Username,
		})
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, false, resp["success"])
		assert.Equal(t, float64(http.StatusUnauthorized), resp["code"])
	})

	t.Run("Invalid Token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "invalid.jwt.token")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, false, resp["success"])
		assert.Equal(t, float64(http.StatusUnauthorized), resp["code"])
	})

	t.Run("Valid Token", func(t *testing.T) {
		token := generateTestToken(t, 100, "testadmin")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", token)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, float64(100), resp["user_id"])
		assert.Equal(t, "testadmin", resp["username"])
	})
}

func TestTimeoutMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(Timeout(50 * time.Millisecond))

	r.GET("/fast", func(c *gin.Context) {
		c.String(http.StatusOK, "fast response")
	})

	r.GET("/slow", func(c *gin.Context) {
		time.Sleep(120 * time.Millisecond)
		c.String(http.StatusOK, "slow response")
	})

	t.Run("Fast Request Completes Within Timeout", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/fast", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "fast response", w.Body.String())
	})

	t.Run("Slow Request Triggers 504 Timeout", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/slow", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, false, resp["success"])
		assert.Equal(t, float64(http.StatusGatewayTimeout), resp["code"])
	})
}

func TestI18nMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(I18nMiddleware())
	r.GET("/i18n-check", func(c *gin.Context) {
		localizer, exists := c.Get("localizer")
		if !exists || localizer == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no localizer"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/i18n-check", nil)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
