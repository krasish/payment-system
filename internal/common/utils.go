package common

import (
	"io"

	"github.com/sirupsen/logrus"
)

// CloseWithLogOnError attempts to close the passed io.Closer and logs on error
func CloseWithLogOnError(closer io.Closer) {
	if err := closer.Close(); err != nil {
		logrus.Errorf("while closing db connection: %v", err)
	}
}
