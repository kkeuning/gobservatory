package main

// PonzuConnection specifies connection details for the Ponzu server
type PonzuConnection struct {
	Scheme string
	Host   string
	Port   string
	Auth   func(*Auth)
}

// Auth contains authentication details for Ponzu
type Auth struct {
	PonzuSecret string
	PonzuUser   string
	PonzuToken  string
	AuthMethod  string
}

// AuthMethod is initialized with the supported Ponzu authentication methods
var AuthMethod = struct {
	Secret string
	Token  string
	None   string
}{
	"Secret",
	"Token",
	"None",
}

// PonzuSecretAuth sets Auth details for authentication using secret
func PonzuSecretAuth(secret string, user string) func(a *Auth) {
	return func(a *Auth) {
		a.PonzuSecret = secret
		a.PonzuUser = user
		a.AuthMethod = AuthMethod.Secret
	}
}

// PonzuTokenAuth sets Auth details for authentication using token
func PonzuTokenAuth(token string) func(a *Auth) {
	return func(a *Auth) {
		a.PonzuToken = token
		a.AuthMethod = AuthMethod.Token
	}
}

// PonzuNoAuth sets Auth details for skipping authentication
func PonzuNoAuth() func(a *Auth) {
	return func(a *Auth) {
		a.AuthMethod = AuthMethod.None
	}
}
