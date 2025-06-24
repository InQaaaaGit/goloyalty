package handlers

import (
	"net/http"

	"gophermart/internal/config"
)

// setAuthCookie устанавливает cookie для аутентификации
func setAuthCookie(w http.ResponseWriter, token string, cfg *config.Config) {
	http.SetCookie(w, &http.Cookie{
		Name:     cfg.CookieName,
		Value:    token,
		Path:     cfg.CookiePath,
		HttpOnly: cfg.CookieHttpOnly,
		Secure:   cfg.CookieSecure,
		MaxAge:   int(cfg.JWTExpiration.Seconds()),
		SameSite: http.SameSiteLaxMode,
	})
}

// clearAuthCookie удаляет cookie аутентификации
// nolint:unused // Эта функция используется в будущих версиях для logout
func clearAuthCookie(w http.ResponseWriter, cfg *config.Config) {
	http.SetCookie(w, &http.Cookie{
		Name:     cfg.CookieName,
		Value:    "",
		Path:     cfg.CookiePath,
		HttpOnly: cfg.CookieHttpOnly,
		Secure:   cfg.CookieSecure,
		MaxAge:   -1, // Удаляем cookie
		SameSite: http.SameSiteLaxMode,
	})
}
