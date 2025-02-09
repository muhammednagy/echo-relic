package echorelic

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	newrelic "github.com/newrelic/go-agent"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	echoRelic *EchoRelic
	e         *echo.Echo
}

func (suite *TestSuite) SetupTest() {
	echoRelic, err := New("test", "1234567890123456789012345678901234567890")
	if err != nil {
		suite.Fail("Failed to create new EchoRelic")
	}
	suite.echoRelic = echoRelic
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	suite.e = e
}

func (suite *TestSuite) TestUseMiddleware() {
	var t newrelic.Transaction
	suite.e.Use(suite.echoRelic.Transaction)
	req := httptest.NewRequest(echo.GET, "/", nil)
	res := httptest.NewRecorder()
	c := suite.e.NewContext(req, res)
	txn := c.Get("newRelicTransaction")
	suite.IsType(t, txn)
}

func (suite *TestSuite) TestBadConfig() {
	_, err := New("test", "1234567890")
	suite.Error(err, "Should error when a key is passed with an invalid length")
}

func (suite *TestSuite) TestShouldReturnNil() {
	handlerFunc := suite.echoRelic.Transaction(func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	req := httptest.NewRequest(echo.GET, "/", nil)
	res := httptest.NewRecorder()
	c := suite.e.NewContext(req, res)
	err := handlerFunc(c)
	suite.NoError(err)
}

func TestMethodSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
