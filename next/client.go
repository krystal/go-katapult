package next

import (
	"fmt"

	"github.com/krystal/go-katapult/next/core"
	"github.com/krystal/go-katapult/next/public"
)

//go:generate ./generate.sh

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
