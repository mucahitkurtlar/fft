package model_test

import (
	"testing"

	"github.com/mucahitkurtlar/fft/pkg/model"
)

func TestGetRobotsTxt(t *testing.T) {
	_, err := model.GetRobotsTxt("https://www.mozilla.org")
	if err != nil {
		t.Error(err)
	}
}

func TestGetUrls(t *testing.T) {
	urls, err := model.GetUrls("https://www.mozilla.org")
	if err != nil {
		t.Error(err)
	}

	if len(urls) == 0 {
		t.Error("no url found")
	}
}
