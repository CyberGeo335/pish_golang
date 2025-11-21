package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Validator interface {
	SignAccess(userID int64, email, role string) (string, error)
	SignRefresh(userID int64, email, role string) (string, error)
	Parse(tokenStr string) (jwtlib.MapClaims, error)
}

// RS256 с двумя действующими ключами и поддержкой kid.
type RS256 struct {
	accessTTL  time.Duration
	refreshTTL time.Duration

	currentKid string
	privKeys   map[string]*rsa.PrivateKey
	pubKeys    map[string]*rsa.PublicKey
}

// NewRS256 — для простоты генерируем 2 ключа при старте процесса.
// key2 считаем "новым", key1 — "старым", оба действуют для проверки.
func NewRS256(_ []byte, accessTTL time.Duration) *RS256 {
	k1, _ := rsa.GenerateKey(rand.Reader, 2048) // в учебном коде ошибки игнорируем
	k2, _ := rsa.GenerateKey(rand.Reader, 2048)

	privs := map[string]*rsa.PrivateKey{
		"key1": k1,
		"key2": k2,
	}
	pubs := map[string]*rsa.PublicKey{
		"key1": &k1.PublicKey,
		"key2": &k2.PublicKey,
	}

	return &RS256{
		accessTTL:  accessTTL,
		refreshTTL: 7 * 24 * time.Hour, // refresh TTL = 7 дней (по заданию)
		currentKid: "key2",             // "новый" ключ для подписи
		privKeys:   privs,
		pubKeys:    pubs,
	}
}

func (r *RS256) SignAccess(userID int64, email, role string) (string, error) {
	return r.sign(userID, email, role, r.accessTTL, "access")
}

func (r *RS256) SignRefresh(userID int64, email, role string) (string, error) {
	return r.sign(userID, email, role, r.refreshTTL, "refresh")
}

func (r *RS256) sign(userID int64, email, role string, ttl time.Duration, typ string) (string, error) {
	now := time.Now()
	claims := jwtlib.MapClaims{
		"sub":   userID,
		"email": email,
		"role":  role,
		"iat":   now.Unix(),
		"exp":   now.Add(ttl).Unix(),
		"iss":   "pz10-auth",
		"aud":   "pz10-clients",
		"typ":   typ, // "access" или "refresh"
	}

	tok := jwtlib.NewWithClaims(jwtlib.SigningMethodRS256, claims)
	tok.Header["kid"] = r.currentKid

	priv := r.privKeys[r.currentKid]
	return tok.SignedString(priv)
}

func (r *RS256) Parse(tokenStr string) (jwtlib.MapClaims, error) {
	token, err := jwtlib.Parse(tokenStr, func(tok *jwtlib.Token) (any, error) {
		if tok.Method != jwtlib.SigningMethodRS256 {
			return nil, errors.New("unexpected signing method")
		}

		kid, _ := tok.Header["kid"].(string)
		if kid == "" {
			return nil, errors.New("missing kid")
		}
		pub, ok := r.pubKeys[kid]
		if !ok {
			return nil, errors.New("unknown kid")
		}
		return pub, nil
	},
		jwtlib.WithAudience("pz10-clients"),
		jwtlib.WithIssuer("pz10-auth"),
	)
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwtlib.MapClaims)
	if !ok {
		return nil, jwtlib.ErrTokenInvalidClaims
	}
	return claims, nil
}
