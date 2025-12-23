package app

type Settings struct {
	EnabledCrawlers []string
}

func NewSettings() Settings {
	return Settings{
		EnabledCrawlers: []string{"goodreads", "libquotes"},
	}
}
