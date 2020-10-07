package katapult

import "fmt"

type service struct {
	client     *Client
	apiVersion string
}

func (s *service) RequestPath(urlStr string) string {
	return fmt.Sprintf("%s/%s", s.apiVersion, urlStr)
}
