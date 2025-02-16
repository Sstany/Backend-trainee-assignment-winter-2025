package integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"shop/pkg/client"

	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

func TestBuyItem(t *testing.T) {
	compose, err := tc.NewDockerCompose("../docker-compose.yaml")
	require.NoError(t, err, "NewDockerComposeAPI()")

	t.Cleanup(func() {
		require.NoError(t, compose.Down(
			context.Background(),
			tc.RemoveOrphans(true),
			tc.RemoveImagesLocal,
			tc.RemoveVolumes(true),
		), "compose.Down()")
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	require.NoError(t, compose.Up(ctx, tc.Wait(true)), "compose.Up()")

	cli, err := client.NewClient("http", client.WithBaseURL("http://localhost:8080"))
	require.NoError(t, err)

	//
	// Authenticate user1
	//

	user1 := client.AuthRequest{
		Password: "testtest",
		Username: "user1",
	}

	resp, err := cli.PostApiAuth(
		ctx,
		user1,
	)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var tokenUser1 client.AuthResponse

	err = json.NewDecoder(resp.Body).Decode(&tokenUser1)
	require.NoError(t, err)

	require.NotNil(t, tokenUser1.Token)

	resp.Body.Close()

	//
	// Buy hoody for user1
	//

	authMiddleware := newTokenMiddleware(*tokenUser1.Token)

	resp, err = cli.GetApiBuyItem(ctx, "hoody", authMiddleware)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp.Body.Close()

	//
	// Check hoody in user1 info
	//

	resp, err = cli.GetApiInfo(ctx, authMiddleware)
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var infoUser1 client.InfoResponse

	err = json.NewDecoder(resp.Body).Decode(&infoUser1)
	require.NoError(t, err)

	require.Equal(t, 700, *infoUser1.Coins)

	for _, value := range *infoUser1.Inventory {
		if *value.Type == "hoody" {
			return
		}
	}

	t.Error(`bought item "hoody" not found in info`)
}

func TestSend(t *testing.T) {
	compose, err := tc.NewDockerCompose("../docker-compose.yaml")
	require.NoError(t, err, "NewDockerComposeAPI()")

	t.Cleanup(func() {
		require.NoError(t, compose.Down(
			context.Background(),
			tc.RemoveOrphans(true),
			tc.RemoveImagesLocal,
			tc.RemoveVolumes(true),
		), "compose.Down()")
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	require.NoError(t, compose.Up(ctx, tc.Wait(true)), "compose.Up()")

	cli, err := client.NewClient("http", client.WithBaseURL("http://localhost:8080"))
	require.NoError(t, err)

	//
	// Authenticate user1
	//

	user1 := client.AuthRequest{
		Password: "testtest1",
		Username: "user1",
	}

	resp, err := cli.PostApiAuth(
		ctx,
		user1,
	)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var tokenUser1 client.AuthResponse

	err = json.NewDecoder(resp.Body).Decode(&tokenUser1)
	require.NoError(t, err)

	require.NotNil(t, tokenUser1.Token)

	resp.Body.Close()

	//
	// Authenticate user2
	//

	user2 := client.AuthRequest{
		Password: "testtest2",
		Username: "user2",
	}

	resp2, err := cli.PostApiAuth(
		ctx,
		user2,
	)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp2.StatusCode)

	var tokenUser2 client.AuthResponse

	err = json.NewDecoder(resp2.Body).Decode(&tokenUser2)
	require.NoError(t, err)

	require.NotNil(t, tokenUser2.Token)

	resp2.Body.Close()

	//
	// Send coins from user1 to user2
	//

	authMiddleware1 := newTokenMiddleware(*tokenUser1.Token)
	authMiddleware2 := newTokenMiddleware(*tokenUser2.Token)

	sendCoin := client.SendCoinRequest{
		Amount: 100,
		ToUser: user2.Username,
	}

	resp, err = cli.PostApiSendCoin(ctx, sendCoin, authMiddleware1)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp.Body.Close()

	//
	// Check coins in user1 info
	//

	resp, err = cli.GetApiInfo(ctx, authMiddleware1)
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var infoUser1 client.InfoResponse

	err = json.NewDecoder(resp.Body).Decode(&infoUser1)
	require.NoError(t, err)

	require.Equal(t, 900, *infoUser1.Coins)

	for _, value := range *infoUser1.CoinHistory.Sent {
		if *value.ToUser == user2.Username || *value.Amount == 100 {
			return
		}
	}

	t.Error(`send coins to user2 not found in user1 info`)

	//
	// Check coins in user2 info
	//

	resp, err = cli.GetApiInfo(ctx, authMiddleware2)
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var infoUser2 client.InfoResponse

	err = json.NewDecoder(resp.Body).Decode(&infoUser2)
	require.NoError(t, err)

	require.Equal(t, 1100, *infoUser2.Coins)

	for _, value := range *infoUser2.CoinHistory.Received {
		if *value.FromUser == user1.Username || *value.Amount == 100 {
			return
		}
	}

	t.Error(`get coins from user1 not found in user2 info`)
}

func newTokenMiddleware(token string) func(ctx context.Context, req *http.Request) error {
	return func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}
}
