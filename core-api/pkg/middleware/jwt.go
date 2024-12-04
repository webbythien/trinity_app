package middleware

import (
	"errors"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	config "github.com/hrshadhin/fiber-go-boilerplate/pkg/config"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/constants"
)

type UserJWTProtect struct {
	RoleID  int    `json:"role_id"`
	UserID  int    `json:"user_id"`
	IsAdmin bool   `json:"admin"`
	Email   string `json:"email"`
	IsExist bool   `json:"is_exist"`
}

// JWTProtected func for specify route group with JWT authentication.
// See: https://github.com/gofiber/jwt
func JWTProtected() func(*fiber.Ctx) error {
	// Create config for JWT authentication middleware.
	jwtwareConfig := jwtware.Config{
		SigningKey:     []byte(config.AppCfg().JWTSecretKey),
		ContextKey:     "user",
		ErrorHandler:   jwtError,
		SuccessHandler: verifyTokenExpiration,
	}

	return jwtware.New(jwtwareConfig)
}

func JWTOptional() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header.
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// No token present, proceed with default user context.
			userLocal := UserJWTProtect{
				IsExist: false,
			}
			c.Locals(constants.UserLocal, userLocal)
			return c.Next()
		}

		// Parse the token.
		token, err := jwt.Parse(authHeader[7:], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(config.AppCfg().JWTSecretKey), nil
		})

		if err != nil {
			// Invalid token, set default user context.
			userLocal := UserJWTProtect{
				IsExist: false,
			}
			c.Locals(constants.UserLocal, userLocal)
			return c.Next()
		}

		// Token is valid, verify claims.
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			expires := int64(claims["exp"].(float64))
			if time.Now().Unix() > expires {
				return jwtError(c, errors.New("token expired"))
			}

			userID, ok := claims["user_id"].(float64)
			if !ok {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"msg": "Invalid user_id claim in token",
				})
			}
			roleID, ok := claims["role_id"].(float64)
			if !ok {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"msg": "Invalid role_id claim in token",
				})
			}
			admin, ok := claims["admin"].(bool)
			if !ok {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"msg": "Invalid admin claim in token",
				})
			}
			email, ok := claims["email"].(string)
			if !ok {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"msg": "Invalid email claim in token",
				})
			}

			// Set user context with token details.
			userLocal := UserJWTProtect{
				RoleID:  int(roleID),
				UserID:  int(userID),
				IsAdmin: admin,
				Email:   email,
				IsExist: true,
			}
			c.Locals(constants.UserLocal, userLocal)
			return c.Next()
		}

		// Token invalid, proceed with default user context.
		userLocal := UserJWTProtect{
			IsExist: false,
		}
		c.Locals(constants.UserLocal, userLocal)
		return c.Next()
	}
}

// "user_id": userID,
//
//	"email":   email,
//	"role_id": roleID,
//	"admin":   isAdmin,
//	"iat":     time.Now().Unix(),
//	"exp":
func verifyTokenExpiration(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*jwt.Token)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "Invalid token",
		})
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "Invalid token claims",
		})
	}

	// Check token expiration
	exp, ok := claims["exp"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "Invalid token expiration",
		})
	}

	expires := int64(exp)
	if time.Now().Unix() > expires {
		return jwtError(c, errors.New("token expired"))
	}

	// Extract and validate claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"msg": "Invalid user_id claim in token",
		})
	}

	roleID, ok := claims["role_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"msg": "Invalid role_id claim in token",
		})
	}

	admin, ok := claims["admin"].(bool)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"msg": "Invalid admin claim in token",
		})
	}
	email, ok := claims["email"].(string)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"msg": "Invalid email claim in token",
		})
	}

	// Populate user local with validated claims
	userLocal := UserJWTProtect{
		RoleID:  int(roleID),
		UserID:  int(userID),
		IsAdmin: admin,
		Email:   email,
		IsExist: true,
	}
	c.Locals(constants.UserLocal, userLocal)

	return c.Next()
}

func jwtError(c *fiber.Ctx, err error) error {
	// Return status 401 and failed authentication error.
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// Return status 401 and failed authentication error.
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"msg": err.Error(),
	})
}

func IsAdmin(c *fiber.Ctx) error {
	user := c.Locals(constants.UserLocal).(UserJWTProtect)

	if !user.IsAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"msg": "Forbidden",
		})
	}

	return c.Next()
}
