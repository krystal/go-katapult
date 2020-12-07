package katapult

import "net/url"

type queryable interface {
	queryValues() *url.Values
}

func queryValues(objs ...queryable) *url.Values {
	merged := &url.Values{}

	for _, obj := range objs {
		if obj == nil {
			continue
		}

		urlVals := *obj.queryValues()
		for k, vals := range urlVals {
			for _, v := range vals {
				merged.Add(k, v)
			}
		}
	}

	return merged
}
