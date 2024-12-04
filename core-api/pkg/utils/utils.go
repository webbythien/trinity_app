package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/config"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/constants"
)

func NewID() string {
	id := uuid.New().String()
	id = strings.ReplaceAll(id, "-", "")
	return id
}

type TokenMetadata struct {
	UserID    uint
	Email     string
	RoleID    uint
	IssuedAt  time.Time
	ExpiredAt time.Time
}

// GenerateNewTokens creates new JWT tokens (access and refresh)
func GenerateNewTokens(userID uint, email string, roleID uint, isAdmin bool) (string, error) {
	// Create token claims
	secretKey := []byte(config.AppCfg().JWTSecretKey)

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role_id": roleID,
		"admin":   isAdmin,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Duration(config.AppCfg().JWTSecretExpireMinutesCount) * time.Minute).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token using the secret signing key
	encodedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return encodedToken, nil
}

// ExtractTokenMetadata extracts metadata from JWT token
func ExtractTokenMetadata(token *jwt.Token) (*TokenMetadata, error) {
	claims := token.Claims.(jwt.MapClaims)

	userID := uint(claims["user_id"].(float64))
	email := claims["email"].(string)
	roleID := uint(claims["role_id"].(float64))
	issuedAt := time.Unix(int64(claims["iat"].(float64)), 0)
	expiredAt := time.Unix(int64(claims["exp"].(float64)), 0)

	return &TokenMetadata{
		UserID:    userID,
		Email:     email,
		RoleID:    roleID,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}, nil
}

func GenerateHash(campaignID int, platformID int) string {
	source := fmt.Sprintf("%d%d%d", campaignID, platformID, time.Now().UnixNano())
	hash := md5.Sum([]byte(source))
	// Convert to hex and take first 6 characters
	hexHash := strings.ToUpper(hex.EncodeToString(hash[:]))[:6]
	return hexHash
}

func GenerateMockPaymentToken() string {
	timestamp := time.Now().UnixMicro()

	randomBytes := make([]byte, 3)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic("Failed to generate random bytes")
	}

	uniquePart := fmt.Sprintf("%06X", timestamp%1000000)
	randomPart := base64.RawURLEncoding.EncodeToString(randomBytes)[:3]

	result := uniquePart + randomPart
	if len(result) > 6 {
		result = result[:6]
	}

	return strings.ToUpper(result)
}

func GetIPFromContext(ctx *fiber.Ctx) string {
	forwarded := ctx.Get("X-Forwarded-For")
	if forwarded != "" {
		ctx.Locals(constants.UserIP, forwarded)
		return forwarded
	}

	realIP := ctx.Get("X-Real-IP")
	if realIP != "" {
		ctx.Locals(constants.UserIP, realIP)
		return realIP
	}

	ip := ctx.IP()
	ctx.Locals(constants.UserIP, ip)
	return ip
}

func GetUserAgentFromContext(ctx *fiber.Ctx) string {
	userAgent := ctx.Get("User-Agent")
	if userAgent != "" {
		ctx.Locals(constants.UserAgent, userAgent)
	}
	return userAgent
}
