package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"buf.build/gen/go/pocketsign/apis/connectrpc/go/pocketsign/link/v1/linkv1connect"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

// 取得してきたAccessTokenを保存する
var tokenSourceStore = map[string]oauth2.TokenSource{}
var refreshTokenStore = map[string]string{}
var errNotAuthenticated = errors.New("not authenticated")

// カスタムリソース
var customResourceID string

var userAPIConnectClient linkv1connect.UserResourceServiceClient

type handler struct {
	conf     *oauth2.Config
	provider *oidc.Provider
	store    sessions.Store
}

type UserResource struct {
	Value string `json:"value"`
}

func main() {
	ctx := context.Background()

	// 環境変数の読み込み
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	host := os.Getenv("OIDC_DEMO_HOST")

	// ユーザーリソース取得APIのConnect Clientを初期化
	userAPIConnectClient = linkv1connect.NewUserResourceServiceClient(
		http.DefaultClient,
		"https://api."+host,
	)

	provider, err := oidc.NewProvider(ctx, "https://oidc."+host)
	if err != nil {
		log.Fatal(err)
	}

	// OAuth2 Clientの設定
	oauth2Config := oauth2.Config{
		ClientID:     os.Getenv("OIDC_CLIENT_ID"),
		ClientSecret: os.Getenv("OIDC_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OIDC_REDIRECT_URL"),

		Endpoint: provider.Endpoint(),

		Scopes: []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess},
	}

	// Sessionの設定(Cookieに保存)
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

	// 自分で定義したリソースIDの取得
	customResourceID = os.Getenv("CUSTOM_RESOURCE_ID")

	h := &handler{
		conf:     &oauth2Config,
		provider: provider,
		store:    store,
	}

	// ハンドラの設定
	http.HandleFunc("/", h.handleIndex)
	http.HandleFunc("/update", h.handlePut)
	http.HandleFunc("/login", h.handleRedirect)
	http.HandleFunc("/callback", h.handleCallback)
	http.HandleFunc("/refresh", h.handleRefresh)

	// サーバーの起動
	log.Fatal(http.ListenAndServe(":8080", nil))

}
