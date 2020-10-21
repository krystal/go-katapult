package katapult

import (
	"net/url"
)

type service struct {
	client *Client
}

type pathHelper struct {
	basePath *url.URL
}

func newPathHelper(path string) (*pathHelper, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	return &pathHelper{basePath: u}, nil
}

func (s *pathHelper) RequestPath(path string) (string, error) {
	u, err := s.basePath.Parse(path)
	if err != nil {
		return "", err
	}

	return u.Path, nil
}
