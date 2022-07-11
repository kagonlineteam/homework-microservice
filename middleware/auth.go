package middleware

import (
	"os"
	"strings"

	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const CLASS_ATTRIBUTE_GIN_NAME = "jwt_class"
const GRADE_ATTRIBUTE_GIN_NAME = "jwt_grade"
const USERNAME_ATTRIBUTE_GIN_NAME = "jwt_username"
const TEACHER_ATTRIBUTE_GIN_NAME = "jwt_teacher"
const SHOWALL_ATTRIBUTE_GIN_NAME = "jwt_showall"
const ALLOWDELETE_ATTRIBUTE_GIN_NAME = "jwt_allowdelete"

var jwtPublicKey, _ = jwt.ParseRSAPublicKeyFromPEM([]byte(strings.Replace(os.Getenv("HOMEWORK_JWT_PUB_KEY"), `\n`, "\n", -1)))

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Request.Header["Authorization"]) < 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "No token found!",
			})
			return
		}
		tokenString := strings.Replace(c.Request.Header["Authorization"][0], "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtPublicKey, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token is not valid: " + err.Error(),
			})
			return
		}

		// validate the essential claims
		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token not valid",
			})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Do not accept refresh tokens with
			if claims["sub"] != "access_main" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Token not valid",
				})
				return
			}
			c.Set(GRADE_ATTRIBUTE_GIN_NAME, claims["stufe"])
			c.Set(CLASS_ATTRIBUTE_GIN_NAME, claims["klasse"])
			c.Set(USERNAME_ATTRIBUTE_GIN_NAME, claims["user"])
			if claims["roles"] != nil {
				for _, role := range claims["roles"].([]interface{}) {
					if role.(string) == "ROLE_TEACHER" || role.(string) == "ROLE_ADMINISTRATOR" {
						c.Set(TEACHER_ATTRIBUTE_GIN_NAME, true)
					}
				}
			}
			if claims["consent"] != nil {
				for _, role := range claims["consent"].([]interface{}) {
					if role.(string) == "homework-show-all" {
						c.Set(SHOWALL_ATTRIBUTE_GIN_NAME, true)
					} else if role.(string) == "homework-allow-delete" {
						c.Set(ALLOWDELETE_ATTRIBUTE_GIN_NAME, true)
					}
				}
			}
		}
	}
}
