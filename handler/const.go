package handler

const (
	sessionKeyUser = "oauth2_user"

	// 認証後に最終的に戻るURL
	headerNameInitiateBackTo = "X-Nginx-Auth-Provider-Initiate-Back-To"

	// 認証後に認証対象のアプリのCookieに認証情報を書き込むために、一時的にリダイレクトするURL
	headerNameInitiateCallback = "X-Nginx-Auth-Provider-Initiate-Callback"

	// 認証対象のアプリのCookieが有効になるパス
	headerNameCallbackCookiePath = "X-Nginx-Auth-Provider-Callback-Cookie-Path"
)
