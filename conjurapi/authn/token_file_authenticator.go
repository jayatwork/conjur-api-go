package authn

import (
	"time"
)

type TokenFileAuthenticator struct {
	TokenFile   string `env:"CONJUR_AUTHN_TOKEN_FILE"`
	mTime       time.Time
	MaxWaitTime time.Duration
}

func (a *TokenFileAuthenticator) RefreshToken() ([]byte, error) {
	maxWaitTime := a.MaxWaitTime
	var timeout <-chan time.Time
	if maxWaitTime == -1 {
		timeout = nil
	} else {
		timeout = time.After(a.MaxWaitTime)
	}

	bytes, err := waitForTextFile(a.TokenFile, timeout)
	if err == nil {
		fi, _ := AppFS.Stat(a.TokenFile)
		a.mTime = fi.ModTime()
	}
	return bytes, err
}

func (a *TokenFileAuthenticator) NeedsTokenRefresh() bool {
	fi, _ := AppFS.Stat(a.TokenFile)
	return a.mTime != fi.ModTime()
}
