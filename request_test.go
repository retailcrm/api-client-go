package retailcrm

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConnectRequestTest struct {
	suite.Suite
}

func TestConnectRequest(t *testing.T) {
	suite.Run(t, new(ConnectRequestTest))
}

func (t *ConnectRequestTest) Test_SystemURL() {
	t.Assert().Equal("", ConnectRequest{}.SystemURL())
	t.Assert().Equal("https://test.retailcrm.pro", ConnectRequest{URL: "https://test.retailcrm.pro"}.SystemURL())
	t.Assert().Equal("https://test.retailcrm.pro", ConnectRequest{URL: "https://test.retailcrm.pro/"}.SystemURL())
}

func (t *ConnectRequestTest) Test_Verify() {
	t.Assert().True(ConnectRequest{
		APIKey: "key",
		Token:  createConnectToken("key", "secret"),
	}.Verify("secret"))
	t.Assert().False(ConnectRequest{
		APIKey: "key",
		Token:  createConnectToken("key", "secret2"),
	}.Verify("secret"))
}

func createConnectToken(apiKey, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	if _, err := mac.Write([]byte(apiKey)); err != nil {
		panic(err)
	}
	return hex.EncodeToString(mac.Sum(nil))
}
