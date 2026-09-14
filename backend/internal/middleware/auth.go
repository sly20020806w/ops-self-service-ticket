package middleware

import (
        "net/http"
        "strings"
        "time"

        "github.com/gin-gonic/gin"
        "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
        UserID   uint   `json:"uid"`
        Username string `json:"username"`
        Role     string `json:"role"`
        jwt.RegisteredClaims
}

func Sign(secret string, uid uint, username, role string) (string, error) {
        claims := Claims{
                UserID:   uid,
                Username: username,
                Role:     role,
                RegisteredClaims: jwt.RegisteredClaims{
                        ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
                        IssuedAt:  jwt.NewNumericDate(time.Now()),
                },
        }
        t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        return t.SignedString([]byte(secret))
}

func Auth(secret string) gin.HandlerFunc {
        return func(c *gin.Context) {
                h := c.GetHeader("Authorization")
                if !strings.HasPrefix(h, "Bearer ") {
                        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
                        return
                }
                tokenStr := strings.TrimPrefix(h, "Bearer ")
                claims := &Claims{}
                token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
                        return []byte(secret), nil
                })
                if err != nil || !token.Valid {
                        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录已失效"})
                        return
                }
                c.Set("uid", claims.UserID)
                c.Set("username", claims.Username)
                c.Set("role", claims.Role)
                c.Next()
        }
}

func CORS() gin.HandlerFunc {
        return func(c *gin.Context) {
                c.Header("Access-Control-Allow-Origin", "*")
                c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
                c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
                if c.Request.Method == http.MethodOptions {
                        c.AbortWithStatus(204)
                        return
                }
                c.Next()
        }
}
