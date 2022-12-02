package proxy

import (
	"context"
	"crypto/subtle"
	"net/http"

	"github.com/mmohamed/loki-auth-proxy/src/pkg"
)

type key int

const (
	OrgIDKey key = iota
	realm        = "Loki Auth Proxy"
)

func BasicAuth(handler http.HandlerFunc, authConfig *pkg.Authn, orgCheck bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		requestOrgID := r.Header.Get("X-Scope-OrgID")
		authorized, orgID := isAuthorized(user, pass, authConfig, orgCheck, requestOrgID)
		if !ok || !authorized {
			writeUnauthorisedResponse(w)
			return
		}
		ctx := context.WithValue(r.Context(), OrgIDKey, orgID)
		handler(w, r.WithContext(ctx))
	}
}

func isAuthorized(user string, pass string, authConfig *pkg.Authn, orgCheck bool, requestOrgID string) (bool, string) {
	for _, v := range authConfig.Users {
		if subtle.ConstantTimeCompare([]byte(user), []byte(v.Username)) == 1 && subtle.ConstantTimeCompare([]byte(pass), []byte(v.Password)) == 1 {
			if orgCheck && (subtle.ConstantTimeCompare([]byte(requestOrgID), []byte(v.OrgID)) == 1 || subtle.ConstantTimeCompare([]byte("*"), []byte(v.OrgID)) == 1) {
				//  is orgCheck is activated and an account with "*" configured as OrgID, the request orgID will be sended to Loki
				return true, requestOrgID
			} else if !orgCheck {
				return true, v.OrgID
			}
		}
	}
	return false, ""
}

func writeUnauthorisedResponse(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
	w.WriteHeader(401)
	w.Write([]byte("Unauthorised\n"))
}
