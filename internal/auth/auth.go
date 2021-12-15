package auth

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// KeyLookupFunc is used to map a JWT Key ID (KID) to the corresponding public
// key. It is a requirement for creating an Authenticator.
// ===========================================================================
// * Private keys should be rotated. During the transition period, tokens signed
// with the old and new keys can coexist by looking up the correct public key by
// key id (KID).
// ===========================================================================
// * key-id-to-public-key resolution is usually accomplished via a public JWKS
// endpoints. See https://auth0.com/docs/jwks for more details.
type KeyLookupFunc func(kid string) (*rsa.PublicKey, error)

// NewSimpleKeyLookup is a simple implementation of KeyFunc that only ever
// supports one key. This is easy for development but in production should be
// replaced with a caching layer that calls a JKWS endpoint.
func NewSimpleKeyLookup(activeKID string, publickKey *rsa.PublicKey) KeyLookupFunc {
	f := func(kid string) (*rsa.PublicKey, error) {
		if activeKID != kid {
			return nil, fmt.Errorf("unrecognized key id %q", kid)
		}
		return publickKey, nil
	}

	return f
}

// Authenticator is used to authenticate clients. It can generate a token for a
// set of user claims and receive the claims by parsing the token.
type Authenticator struct {
	privateKey   *rsa.PrivateKey
	activeKID    string
	algorithm    string
	pubKeyLookup KeyLookupFunc
	parser       *jwt.Parser
}

// NewAuthenticator creates an *Authenticator for use. It will error if:
// - The Private key is nil.
// - The Public key func in nil.
// - The public key id is blank.
// - The specified algorithm is unsupported.
func NewAuthenticator(privateKey *rsa.PrivateKey, activeKID, algorithm string,
	pubKeyLookupFunc KeyLookupFunc) (*Authenticator, error) {
	if privateKey == nil {
		return nil, errors.New("Private key cannot be null.")
	}

	if activeKID == "" {
		return nil, errors.New("Active KID cannot be blank.")
	}

	if jwt.GetSigningMethod(algorithm) == nil {
		return nil, errors.Errorf("Unknown algorithm %v", algorithm)
	}

	// Create the token parser to use. The algorithm used to sign the JWT must
	// be validated to avoid a critical vulnerability.
	// https://auth0.com/blog/critical-vulnerability-in-json-web-token-libraries/
	parser := jwt.Parser{
		ValidMethods: []string{algorithm},
	}

	a := Authenticator{
		privateKey:   privateKey,
		activeKID:    activeKID,
		algorithm:    algorithm,
		pubKeyLookup: pubKeyLookupFunc,
		parser:       &parser,
	}

	return &a, nil
}

// GenerateToken generates a signed JWT token string representing the user
// Claims.
func (a *Authenticator) GenerateToken(claims Claims) (string, error) {
	method := jwt.GetSigningMethod(a.algorithm)

	tkn := jwt.NewWithClaims(method, claims)
	tkn.Header["kid"] = a.activeKID

	str, err := tkn.SignedString(a.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "signed token")
	}

	return str, nil
}
