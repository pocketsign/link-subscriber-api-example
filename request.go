package main

import (
	"context"
	"log"

	linkv1 "buf.build/gen/go/pocketsign/apis/protocolbuffers/go/pocketsign/link/v1"
	"golang.org/x/oauth2"

	"connectrpc.com/connect"
	"github.com/gorilla/sessions"
)

// リソースの取得を行う
func getUserResource(ctx context.Context, sess *sessions.Session, id string) (UserResource, error) {
	req, err := setupRequest(connect.NewRequest(&linkv1.GetUserResourceRequest{
		ResourceId: id,
	}), sess)
	if err != nil {
		return UserResource{}, err
	}

	resp, err := userAPIConnectClient.GetUserResource(ctx, req)
	if err != nil {
		return UserResource{}, err
	}

	return UserResource{Value: resp.Msg.GetValue()}, nil
}

// リソースの更新を行う
func putUserResource(ctx context.Context, sess *sessions.Session, id string, v UserResource) error {
	req, err := setupRequest(connect.NewRequest(&linkv1.UpdateUserResourceRequest{
		ResourceId: id,
		Value:      v.Value,
	}), sess)

	if err != nil {
		return err
	}

	_, err = userAPIConnectClient.UpdateUserResource(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

// RequestにAuthorizationヘッダを付与する
func setupRequest[T any](r *connect.Request[T], sess *sessions.Session) (*connect.Request[T], error) {
	subject, ok := sess.Values["sub"].(string)
	if !ok {
		return nil, errNotAuthenticated
	}

	ts, ok := tokenSourceStore[subject]
	if !ok {
		return nil, errNotAuthenticated
	}

	token, err := ts.Token()

	if err != nil {
		return nil, err
	}

	r.Header().Add("Authorization", "Bearer "+token.AccessToken)

	return r, nil
}

// トークンの更新を行う
func refresh(ctx context.Context, conf *oauth2.Config, sess *sessions.Session) error {
	subject, ok := sess.Values["sub"].(string)
	if !ok {
		return errNotAuthenticated
	}

	ref, ok := refreshTokenStore[subject]
	if !ok {
		return errNotAuthenticated
	}

	log.Println("refresh token: ", ref)

	// 更新
	tks := conf.TokenSource(ctx, &oauth2.Token{
		RefreshToken: ref,
	})

	tk, err := tks.Token()
	if err != nil {
		log.Println(err)
		return errNotAuthenticated
	}

	tokenSourceStore[subject] = tks
	refreshTokenStore[subject] = tk.RefreshToken

	return nil
}
