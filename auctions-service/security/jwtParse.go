package security

import (
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v4"
)

const (
	Secret = "G+KbPeShVmYq3t6w9z$C&F)J@McQfTjW"
)

type JwtParser struct {
	secret        string
	secretAsBytes []byte
}

func NewJwtParser(secret string) *JwtParser {
	return &JwtParser{
		secret:        secret,
		secretAsBytes: []byte(secret),
	}
}

func (jwtParser *JwtParser) Parse(tokenString string) (map[string]interface{}, error) {
	log.Printf("[JwtParser] authenticating bearer token %s\n", tokenString)
	// alertEngine.sendToConsole(msg)
	// see: https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-Parse-Hmac

	claims := jwt.MapClaims{}

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return jwtParser.secretAsBytes, nil
	})

	payload := map[string]interface{}{}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		for key, val := range claims {
			payload[key] = val
		}
		return payload, nil
	} else {
		log.Println("[JwtParser] could not validate token; provided token invalid, or information was parsed incorrectly from valid token")
		return payload, err
	}
}

type Authenticator struct {
	jwtParser *JwtParser
}

func NewAuthenticator(jwtParser *JwtParser) *Authenticator {
	return &Authenticator{
		jwtParser: jwtParser,
	}
}

func (authenticator *Authenticator) extractDetails(tokenString string) (map[string]interface{}, error) {
	payload, err := authenticator.jwtParser.Parse(tokenString)
	return payload, err
}

func (authenticator *Authenticator) IsAdmin(tokenString string) bool {
	payload, err := authenticator.extractDetails(tokenString)
	if err != nil {
		return false // trivial fail
	}
	for key, val := range payload {
		if key == "authorities" {
			roles := val.([]interface{})
			for _, role := range roles {
				if role.(string) == "ROLE_ADMIN" {
					return true
				}
			}
		}
	}
	return false
}

func (authenticator *Authenticator) IsUser(tokenString string) bool {
	payload, err := authenticator.extractDetails(tokenString)
	if err != nil {
		return false // trivial fail
	}
	for key, val := range payload {
		if key == "authorities" {
			roles := val.([]interface{})
			for _, role := range roles {
				if role.(string) == "ROLE_USER" {
					return true
				}
			}
		}
	}
	return false
}

func (authenticator *Authenticator) MeetsRequirements(tokenString string, mustBeUser, mustBeAdmin bool) bool {
	payload, err := authenticator.extractDetails(tokenString)
	if err != nil {
		return false // trivial fail
	}
	isAdmin := false
	isUser := false
	for key, val := range payload {
		if key == "authorities" {
			roles := val.([]interface{})
			for _, role := range roles {
				if role.(string) == "ROLE_ADMIN" {
					isAdmin = true
				}
				if role.(string) == "ROLE_USER" {
					isUser = true
				}
			}
		}
	}
	meetsRequirements := true
	if mustBeAdmin && !isAdmin {
		meetsRequirements = false
	}
	if mustBeUser && !isUser {
		meetsRequirements = false
	}
	return meetsRequirements
}
