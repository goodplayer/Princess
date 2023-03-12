package ginplugin

import (
	"crypto/rand"
	"encoding/base32"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

var sessionUtilsLogger = logrus.New()

func NewSessionId() string {
	k := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		sessionUtilsLogger.Errorf("generate session id failed:%s", err)
		panic(err)
	}
	return strings.TrimRight(base32.StdEncoding.EncodeToString(k), "=")
}
