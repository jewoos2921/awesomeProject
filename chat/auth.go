package main

import (
	"crypto/md5"
	"fmt"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"io"
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	//if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == ""
	if err == http.ErrNoCookie {
		// 인증 불가
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		// 다른 에러
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 성공 - 다음 헨들러
	h.next.ServeHTTP(w, r)
}

func MathAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler는 서드파티 로그인 프로세스를 처리한다.
// 형식: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s",
				provider, err), http.StatusBadRequest)
			return
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to GetBeginAuthURL for %s: %s",
				provider, err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s", provider, err), http.StatusBadRequest)
			return
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to complet auth for %s: %s",
				provider, err), http.StatusInternalServerError)
			return
		}
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("Error when trying to get user from", provider, "-", err)
		}
		chatUser := &chatUser{User: user}

		// 사용자 식별
		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Email()))
		chatUser.uniqueID = fmt.Sprintf("%x", m.Sum(nil))
		avatarURL, err := avatars.GetAvatarURL(chatUser)
		if err != nil {
			log.Fatalln("Error whn trying to GetAvatarURL", "-", err)
		}
		authCookieValue := objx.New(map[string]interface{}{
			"userid":     chatUser.uniqueID,
			"name":       user.Name(),
			"avatar_url": avatarURL,
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{Name: "auth", Value: authCookieValue, Path: "/"})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}
