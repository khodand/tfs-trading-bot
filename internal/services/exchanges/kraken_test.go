package exchanges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeAuth(t *testing.T) {
	postData := "orderType=lmt&symbol=pi_xbtusd&side=buy&size=10000&limitPrice=9400"
	endpointPath := "/api/v3/sendorder"

	expected := "9sZ3JGYUmDjullEp/qyR034ktzXCmU/oIQmFJWZyNHntTig8zgCo5/uD25LMJeZkSouyupn7NTw++Su11+kUjA=="
	actual := encodeAuth(postData, endpointPath, "U3R1YlNlY3JldEFwaUtleQ==")

	assert.Equal(t, expected, actual)
}
