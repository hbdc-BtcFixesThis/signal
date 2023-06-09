package main

import (
	"time"

	"net/http"
)

func VerifyToken(token string, hash string) Response {
	ua := http.StatusUnauthorized
	noTokenErr := "No auth token provided"
	expiredErr := "Your token has expired. Please sign in again"

	if len(token) == 0 {
		return Response{Err: noTokenErr, StatusCode: ua}
	}
	if IsValidToken(token, hash) {
		return Response{Msg: "Welcome!", StatusCode: http.StatusOK}
	}
	if TokenExpired(token, hash) {
		return Response{Err: expiredErr, StatusCode: ua}
	}
	return Response{Err: "Unauthorized!", StatusCode: ua}
}

func IsValidToken(token string, hash string) bool {
	// token create by the client should consist of a hash of
	// the password plus the utc date of the token creation.
	// The token is valid for a day. This seems good enough
	// for now and if needed allowing users to dictate the
	// lifespan of a session should be fairly easy by adding
	// a setting and updating the check here
	return token == SHA256(hash+" "+time.Now().UTC().String()[:10])
}

func TokenExpired(token string, hash string) bool {
	return token == SHA256(hash+" "+time.Now().AddDate(0, 0, -1).UTC().String()[:10])
}

func (ss *SignalServer) CheckAPIKey(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		passHash := ss.SC.PassHash()

		resp := VerifyToken(key, passHash)
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			JSONResponse(w, resp)
		default:
			h.ServeHTTP(w, r)
		}
	})
}

func (ss *SignalServer) verifyHandler(w http.ResponseWriter, r *http.Request) {
	ss.Respond(w, VerifyToken(r.URL.Query().Get("key"), ss.SC.PassHash()))
}
