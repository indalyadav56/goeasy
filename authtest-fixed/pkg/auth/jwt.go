package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the JWT claims structure
type JWTClaims struct {
	UserID   uint     `json:"user_id"`
	Email    string   `json:"email"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Type     string   `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// GenerateJWT generates a JWT token with given claims
func GenerateJWT(claims jwt.MapClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateJWT validates a JWT token and returns claims
func ValidateJWT(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Check if token is expired
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, errors.New("token expired")
		}
	}

	return claims, nil
}

// ExtractTokenFromHeader extracts JWT token from Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	// Check if header starts with "Bearer "
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("invalid authorization header format")
	}

	token := authHeader[len(bearerPrefix):]
	if token == "" {
		return "", errors.New("token is required")
	}

	return token, nil
}

// GenerateAccessToken generates an access token for a user
func GenerateAccessToken(userID uint, email, username string, roles []string, secret string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"email":    email,
		"username": username,
		"roles":    roles,
		"type":     "access",
		"exp":      time.Now().Add(expiry).Unix(),
		"iat":      time.Now().Unix(),
		"nbf":      time.Now().Unix(),
	}

	return GenerateJWT(claims, secret)
}

// GenerateRefreshToken generates a refresh token for a user
func GenerateRefreshToken(userID uint, secret string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     time.Now().Add(expiry).Unix(),
		"iat":     time.Now().Unix(),
		"nbf":     time.Now().Unix(),
	}

	return GenerateJWT(claims, secret)
}

// ParseClaims parses JWT claims into JWTClaims struct
func ParseClaims(claims jwt.MapClaims) (*JWTClaims, error) {
	jwtClaims := &JWTClaims{}

	if userID, ok := claims["user_id"].(float64); ok {
		jwtClaims.UserID = uint(userID)
	} else {
		return nil, errors.New("invalid user_id in claims")
	}

	if email, ok := claims["email"].(string); ok {
		jwtClaims.Email = email
	}

	if username, ok := claims["username"].(string); ok {
		jwtClaims.Username = username
	}

	if tokenType, ok := claims["type"].(string); ok {
		jwtClaims.Type = tokenType
	}

	if roles, ok := claims["roles"].([]interface{}); ok {
		jwtClaims.Roles = make([]string, len(roles))
		for i, role := range roles {
			if roleStr, ok := role.(string); ok {
				jwtClaims.Roles[i] = roleStr
			}
		}
	}

	return jwtClaims, nil
}
