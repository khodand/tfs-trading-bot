package exchanges

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeAuth(t *testing.T) {
	postData := "orderType=lmt&symbol=pi_xbtusd&side=buy&size=10000&limitPrice=9400"
	endpointPath := "/api/v3/sendorder"

	expected := "JUU8ZX7kNBU2bO9OKw3sUF6Qp+R7QPRML1JNgdhmeJpyo/LWJEElts2431zkYdklxVK5sbMMRGDySVXoMvPE5w=="
	actual := encodeAuth(postData, endpointPath)

	assert.Equal(t, expected, actual)
}
