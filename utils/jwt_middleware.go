package utils

import (
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var JwtSecret = GetConfig("JWT_SECRET")

type JWTClaims struct {
	ID        int64
	Name      string
	Module    uint8
	AccRole   uint8
	OwnerUUID *uuid.UUID
	jwt.RegisteredClaims
}
type JWTResponse struct {
	Token   string
	Refresh string
	Error   string
}

func setAuthContext(c *gin.Context, moduleOrigin uint8) bool {
	jwtToken := GetJWTTokenFromSession(c)
	if jwtToken != "" {
		claims, err := CheckJWTToken(jwtToken, c)
		if err == nil {
			if claims.Module == moduleOrigin {
				log.Printf("- Logged in (ID: %d, Name: %s, Module: %d, Role: %d, Owner UUID: %v, Exp. Time: %s)",
					claims.ID, claims.Name, claims.Module, claims.AccRole, claims.OwnerUUID, claims.RegisteredClaims.ExpiresAt)
				c.Set("LoggedIn", true)
				c.Set("ID", claims.ID)
				c.Set("Name", claims.Name)
				c.Set("Module", claims.Module)
				c.Set("Role", claims.AccRole)
				c.Set("OwnerUUID", claims.OwnerUUID)
				c.Set("ExpiresAt", claims.RegisteredClaims.ExpiresAt)
				return true
			}
			log.Printf("!! JWT Token is for a different module")
		}
	}
	SetJWTTokenFromSession(c, nil, nil)
	c.Set("LoggedIn", false)
	c.Set("ID", nil)
	c.Set("Name", nil)
	c.Set("Module", nil)
	c.Set("Role", nil)
	c.Set("OwnerUUID", nil)
	return false
}

func AuthMiddleware(moduleOrigin uint8) gin.HandlerFunc {
	return func(c *gin.Context) {
		setAuthContext(c, moduleOrigin)
		c.Next()
	}
}

func RequireAuthMiddleware(moduleOrigin uint8) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !setAuthContext(c, moduleOrigin) {
			c.Redirect(http.StatusFound, "/loginForm")
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireRoleMiddleware(moduleOrigin uint8, allowedRoles ...uint8) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !setAuthContext(c, moduleOrigin) {
			c.Redirect(http.StatusFound, "/loginForm")
			c.Abort()
			return
		}
		role, _ := c.Get("Role")
		accRole, ok := role.(uint8)
		if !ok || !slices.Contains(allowedRoles, accRole) {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetJWTTokenFromSession(c *gin.Context) string {
	session := sessions.Default(c)
	jwtToken := session.Get("jwt_token")
	if jwtToken != nil {
		log.Printf("JWT Token found")
		return jwtToken.(string)
	} else {
		log.Printf("!! No JWT Token found")
		return ""
	}
}
func SetJWTTokenFromSession(c *gin.Context, jwtToken *string, refreshToken *string) {
	session := sessions.Default(c)
	session.Set("jwt_token", jwtToken)
	session.Set("refresh_token", refreshToken)
	session.Save()
}

func GetJWTTokenFromHeader(c *gin.Context) string {
	tokenString := c.GetHeader("Authorization")
	log.Printf("Token: %s", tokenString)
	if tokenString == "" {
		return ""
	} else {
		return tokenString
	}
}

func CheckJWTToken(jwtToken string, c *gin.Context) (*JWTClaims, error) {
	token, err := ParseJWTToken(jwtToken)
	if err == nil {
		if claims, ok := token.Claims.(*JWTClaims); ok {
			return claims, nil
		} else {
			log.Println("!! Error obtaining claims")
		}
	} else {
		log.Printf("!! Failed checking JWT Token: %s", err)
	}
	return nil, err
}
func ParseJWTToken(jwtToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(jwtToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JwtSecret), nil
	})
}
func GenerateJWTTokens(accID int64, name string, module uint8, accRole uint8, ownerUUID *uuid.UUID) (signedToken string, signedRefreshToken string) {
	claims := &JWTClaims{
		ID:        accID,
		Name:      name,
		Module:    module,
		AccRole:   accRole,
		OwnerUUID: ownerUUID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	refreshClaims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(JwtSecret))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(JwtSecret))
	if err != nil {
		log.Panic(err)
		return
	}
	return token, refreshToken
}
