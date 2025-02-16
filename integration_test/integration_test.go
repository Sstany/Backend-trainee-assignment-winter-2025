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

func newTokenMiddleware(token string) func(ctx context.Context, req *http.Request) error {
	return func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}
}
