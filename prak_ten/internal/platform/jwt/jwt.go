package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"

	"github.com/CyberGeo335/prak_ten/internal/platform/config"
)

// Validator используется в AuthN-мидлваре
type Validator interface {
	Parse(tokenStr string) (jwtlib.MapClaims, error)
}

// RS256 менеджер токенов (access/refresh)
type RS256 struct {
	issuer     string
	audience   string
	accessTTL  time.Duration
	refreshTTL time.Duration

	signKey   *rsa.PrivateKey
	signKID   string
	verifyKey map[string]*rsa.PublicKey
}

// NewRS256 — ИМЕННО эту функцию сейчас ищет router.go
func NewRS256(cfg config.Config) (*RS256, error) {
	if len(cfg.RSKeys) == 0 {
		return nil, errors.New("no RSA keys configured")
	}

	verify := make(map[string]*rsa.PublicKey)
	var sign *rsa.PrivateKey
	var signKID string

	for i, kc := range cfg.RSKeys {
		priv, err := loadPrivateKey(kc.PrivateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("load private key %s: %w", kc.KID, err)
		}
		pub, err := loadPublicKey(kc.PublicKeyPath)
		if err != nil {
			return nil, fmt.Errorf("load public key %s: %w", kc.KID, err)
		}
		verify[kc.KID] = pub
		if i == 0 { // первый ключ в списке — активный для подписи
			sign = priv
			signKID = kc.KID
		}
	}

	return &RS256{
		issuer:     cfg.Issuer,
		audience:   cfg.Audience,
		accessTTL:  cfg.AccessTTL,
		refreshTTL: cfg.RefreshTTL,
		signKey:    sign,
		signKID:    signKID,
		verifyKey:  verify,
	}, nil
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, errors.New("no PEM block")
	}
	if !strings.Contains(block.Type, "PRIVATE KEY") {
		return nil, fmt.Errorf("unexpected PEM type: %s", block.Type)
	}

	// PKCS#1
	if k, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return k, nil
	}
	// PKCS#8
	pk, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	k, ok := pk.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not RSA private key")
	}
	return k, nil
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, errors.New("no PEM block")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	k, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}
	return k, nil
}

// внутренний метод для генерации токена
func (r *RS256) signToken(userID int64, email, role, typ string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwtlib.MapClaims{
		"sub":   userID,
		"email": email,
		"role":  role,
		"iat":   now.Unix(),
		"exp":   now.Add(ttl).Unix(),
		"iss":   r.issuer,
		"aud":   r.audience,
		"typ":   typ, // "access" или "refresh"
	}

	t := jwtlib.NewWithClaims(jwtlib.SigningMethodRS256, claims)
	// kid в заголовок
	t.Header["kid"] = r.signKID

	return t.SignedString(r.signKey)
}

// Публичные методы — используются в сервисе

func (r *RS256) SignAccess(userID int64, email, role string) (string, error) {
	return r.signToken(userID, email, role, "access", r.accessTTL)
}

func (r *RS256) SignRefresh(userID int64, email, role string) (string, error) {
	return r.signToken(userID, email, role, "refresh", r.refreshTTL)
}

func (r *RS256) Parse(tokenStr string) (jwtlib.MapClaims, error) {
	token, err := jwtlib.Parse(tokenStr, func(t *jwtlib.Token) (any, error) {
		if _, ok := t.Method.(*jwtlib.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %T", t.Method)
		}
		kid, _ := t.Header["kid"].(string)
		if kid == "" {
			// fallback — публичный ключ от активного приватного
			return &r.signKey.PublicKey, nil
		}
		pub, ok := r.verifyKey[kid]
		if !ok {
			return nil, fmt.Errorf("unknown kid: %s", kid)
		}
		return pub, nil
	}, jwtlib.WithIssuer(r.issuer), jwtlib.WithAudience(r.audience))
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwtlib.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}
	return claims, nil
}
