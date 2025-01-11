package config_test

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func Test_Base64_01(t *testing.T) {
	v := `LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFOUhLRysxU0RzWUx3SVF0UEJxSXVWZ1VpWWRvVQo4YlJZVDl5TVlITFpQdUJ5Z3VwM0I4T3o4MVN6Ym52UDVTTHBVZ1pTYXd1Q21McTZCQStUbStCZVd3PT0KLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg==`
	bb, err := base64.StdEncoding.DecodeString(v)
	fmt.Println(err, bb)
}
