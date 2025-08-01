package ctxdata

import "context"

func GetUId(c context.Context) string {
	if u, ok := c.Value(IdentityKey).(string); ok {
		return u
	}
	return ""
}
