package model_test

import (
	"os"
	"testing"

	"github.com/mucahitkurtlar/fft/pkg/model"
	"github.com/playwright-community/playwright-go"
)

func TestCrawlPage(t *testing.T) {
	err := playwright.Install()
	if err != nil {
		t.Error(err)
	}

	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
	}

	browser, err := pw.Firefox.Launch()
	if err != nil {
		t.Error(err)
	}

	var fontFamilies []string

	err = model.CrawlPage(browser, "https://motherfuckingwebsite.com", &fontFamilies,
		&model.CrawlerOpts{
			MaxPageCount:   1,
			GoRoutineCount: 1,
			GoToTimeout:    30000,
			NetIdleTimeout: 2000,
		},
	)
	if err != nil {
		t.Error(err)
	}

	if len(fontFamilies) == 0 {
		t.Errorf("font families not found")
	}
}

func TestCrawlSite(t *testing.T) {
	err := playwright.Install()
	if err != nil {
		t.Error(err)
	}

	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
	}

	browser, err := pw.Firefox.Launch()
	if err != nil {
		t.Error(err)
	}

	fontFamilies, err := model.CrawlSite(browser, "https://motherfuckingwebsite.com",
		&model.CrawlerOpts{
			MaxPageCount:   1,
			GoRoutineCount: 1,
			GoToTimeout:    30000,
			NetIdleTimeout: 2000,
		},
	)
	if err != nil {
		t.Error(err)
	}

	if len(fontFamilies) == 0 {
		t.Error("font families not found")
	}
}

func TestStartCrawler(t *testing.T) {
	err := playwright.Install()
	if err != nil {
		t.Error(err)
	}

	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
	}

	err = model.StartCrawler(pw, "https://motherfuckingwebsite.com",
		&model.CrawlerOpts{
			MaxPageCount:   1,
			GoRoutineCount: 1,
			GoToTimeout:    30000,
			NetIdleTimeout: 2000,
		},
	)
	if err != nil {
		t.Error(err)
	}

	err = os.Remove("font-families.txt")
	if err != nil {
		t.Error(err)
	}
}

func TestWriteToFile(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		file    string
		content []string
		wantErr bool
	}{
		{
			name:    "write to file",
			file:    "font-families-test.txt",
			content: []string{"foo", "bar", "baz"},
			wantErr: false,
		},
		{
			name:    "write to file error",
			file:    "",
			content: []string{"foo", "bar", "baz"},
			wantErr: true,
		},
		{
			name:    "write to file error",
			file:    "font-families-test.txt",
			content: []string{},
			wantErr: false,
		},
		{
			name:    "write to file error",
			file:    "font-families-test.txt",
			content: []string{"foo", "bar", "baz"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := model.WriteToFile(tt.file, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	os.Remove("font-families-test.txt")
}

func TestContains(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		slice    []string
		contains string
		want     bool
	}{
		{
			name:     "empty slice",
			slice:    []string{},
			contains: "foo",
			want:     false,
		},
		{
			name:     "slice contains",
			slice:    []string{"foo", "bar", "baz"},
			contains: "bar",
			want:     true,
		},
		{
			name:     "slice does not contain",
			slice:    []string{"foo", "bar", "baz"},
			contains: "qux",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := model.Contains(tt.slice, tt.contains); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
