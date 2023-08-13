package middlewares

import (
	"context"
	"errors"
	"github.com/EricBastos/ProjetoTG/API/configs"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
	"log"
	"net/http"
)

type Authenticator struct {
	userType string
	userDb   database.UserInterface
}

func NewAuthenticator(userType string, userDb database.UserInterface) *Authenticator {
	return &Authenticator{
		userType: userType,
		userDb:   userDb,
	}
}

func (m *Authenticator) Authenticate(auths ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			for _, auth := range auths {
				switch auth {
				case "jwt":
					if token := m.tryJwtAuthentication(r); token != nil && token.Subject() != "" {
						var err error
						ctx, err = m.injectTokenInfo(token, ctx)
						if err != nil {
							//log.Println(err.Error())
							continue
						}
						ctx = jwtauth.NewContext(ctx, token, nil)
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}

				default:
					log.Println("(Authenticator) Unknown auth specified")
				}
			}

			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return

		}
		return http.HandlerFunc(hfn)
	}
}

func (m *Authenticator) injectTokenInfo(token jwt.Token, ctx context.Context) (context.Context, error) {

	nameMaybe, ok := token.Get("name")
	if !ok {
		return nil, errors.New("name param not found")
	}
	name, ok := nameMaybe.(string)
	if !ok {
		return nil, errors.New("name param not found")
	}
	taxIdMaybe, ok := token.Get("taxId")
	if !ok {
		return nil, errors.New("taxId param not found")
	}
	taxId, ok := taxIdMaybe.(string)
	if !ok {
		return nil, errors.New("taxId param not found")
	}
	emailMaybe, ok := token.Get("email")
	if !ok {
		return nil, errors.New("email param not found")
	}
	email, ok := emailMaybe.(string)
	if !ok {
		return nil, errors.New("email param not found")
	}
	ctx = context.WithValue(ctx, "subject", token.Subject())
	ctx = context.WithValue(ctx, "taxId", taxId)
	ctx = context.WithValue(ctx, "email", email)
	ctx = context.WithValue(ctx, "name", name)

	return ctx, nil
}

func (m *Authenticator) tryJwtAuthentication(r *http.Request) jwt.Token {

	token, err := jwtauth.VerifyRequest(configs.Cfg.TokenAuthUser, r, jwtauth.TokenFromHeader, jwtauth.TokenFromCookie)
	if err == nil && token != nil && jwt.Validate(token) == nil {
		return token
	}
	return nil
}
