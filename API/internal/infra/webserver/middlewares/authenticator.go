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

	isSuperuserMaybe, ok := token.Get("isSuperuser")
	if !ok {
		return nil, errors.New("isSuperuser param not found")
	}
	isSuperuser, ok := isSuperuserMaybe.(bool)
	if !ok {
		return nil, errors.New("isSuperuser param not found")
	}
	kycLevelMaybe, ok := token.Get("KYCLevel")
	if !ok {
		return nil, errors.New("[1] KYCLevel param not found")
	}
	kycLevel, ok := kycLevelMaybe.(float64)
	if !ok {
		return nil, errors.New("[2] KYCLevel param not found")
	}
	firstNameMaybe, ok := token.Get("firstName")
	if !ok {
		return nil, errors.New("firstName param not found")
	}
	firstName, ok := firstNameMaybe.(string)
	if !ok {
		return nil, errors.New("firstName param not found")
	}
	lastNameMaybe, ok := token.Get("lastName")
	if !ok {
		return nil, errors.New("lastName param not found")
	}
	lastName, ok := lastNameMaybe.(string)
	if !ok {
		return nil, errors.New("lastName param not found")
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
	isPjMaybe, ok := token.Get("isPj")
	if !ok {
		return nil, errors.New("isPj param not found")
	}
	isPj, ok := isPjMaybe.(bool)
	if !ok {
		return nil, errors.New("isPj param not found")
	}
	ctx = context.WithValue(ctx, "subject", token.Subject())
	ctx = context.WithValue(ctx, "isSuperuser", isSuperuser)
	ctx = context.WithValue(ctx, "KYCLevel", int64(kycLevel))
	ctx = context.WithValue(ctx, "firstName", firstName)
	ctx = context.WithValue(ctx, "lastName", lastName)
	ctx = context.WithValue(ctx, "taxId", taxId)
	ctx = context.WithValue(ctx, "email", email)
	ctx = context.WithValue(ctx, "isPj", isPj)

	return ctx, nil
}

func (m *Authenticator) tryJwtAuthentication(r *http.Request) jwt.Token {

	token, err := jwtauth.VerifyRequest(configs.Cfg.TokenAuthUser, r, jwtauth.TokenFromHeader, jwtauth.TokenFromCookie)
	if err == nil && token != nil && jwt.Validate(token) == nil {
		return token
	}
	return nil
}
