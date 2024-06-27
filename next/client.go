package next

//go:generate go run go.uber.org/mock/mockgen@v0.4.0 -source=./public/public.go -destination=./public/mock/public.go
//go:generate go run go.uber.org/mock/mockgen@v0.4.0 -source=./core/core.go -destination=./core/mock/core.go

import (
	"fmt"

	"github.com/krystal/go-katapult/next/core"
	"github.com/krystal/go-katapult/next/public"
)

type Client struct {
	Core   core.ClientInterface
	Public public.ClientInterface
}

func NewClient(
	coreURL string,
	coreToken string,
	coreHTTPClient core.HttpRequestDoer,
	publicURL string,
	publicToken string,
	publicHTTPClient public.HttpRequestDoer,
) (*Client, error) {
	coreClient, err := core.NewClientWithResponses(
		coreURL,
		coreToken,
		core.WithHTTPClient(coreHTTPClient),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating core client: %w", err)
	}

	publicClient, err := public.NewClientWithResponses(
		publicURL,
		publicToken,
		public.WithHTTPClient(publicHTTPClient),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating public client: %w", err)
	}

	return &Client{
		Core:   coreClient,
		Public: publicClient,
	}, nil
}
