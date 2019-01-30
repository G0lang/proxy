package config

// The following variables are injected at compile time, do not edit
var (
	gver  string
	hash  string
	short string
	date  string
	build string
	count string
	port  string
)

// Config variable inject at build time
type (
	content map[string]string
	config  struct {
		file   string
		config content
	}
)

// Set config
func (c *config) Set(key string, value string) {
	if c.config == nil {
		c.config = make(content)
	}
	c.config[key] = value
}

// Get Config with value
func (c *config) Get(key string) string {
	for k, v := range c.config {
		if k == key {
			return v
		}
	}
	return ""
}

// Config published
var Config config

func init() {
	Config.Set("GITV", gver)
	Config.Set("HASH", hash)
	Config.Set("SHORT", short)
	Config.Set("DATE", date)
	Config.Set("BUILD", build)
	Config.Set("COUNT", count)
	Config.Set("PORT", port)
}
