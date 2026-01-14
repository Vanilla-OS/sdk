package conf

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vanilla-os/sdk/pkg/v1/conf"
	"github.com/vanilla-os/sdk/pkg/v1/conf/types"
)

type ConfigStruct struct {
	Place    string `mapstructure:"place"`
	Event    string `mapstructure:"event"`
	Duration int    `mapstructure:"duration"`
}

func TestInitConfig(t *testing.T) {
	dir := t.TempDir()

	opts := types.ConfigOptions{
		Domain: "org.vanillaos.sdk.conf-test",
		Prefix: dir,
		Type:   "json",
	}

	filePath := filepath.Join(dir, "/etc", opts.Domain, "config.json")

	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		t.Errorf("error creating directory: %v", err)
	}

	content := []byte(`{
"place": "Gotham",
"event": "Joker's Robbery",
"duration": 24
}`)
	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		t.Errorf("error writing file: %v", err)
	}

	config, err := conf.InitConfig[ConfigStruct](opts)
	if err != nil {
		t.Errorf("error initializing config: %v", err)
		return
	}

	assert.Equal(t, "Gotham", config.Place)
	assert.Equal(t, "Joker's Robbery", config.Event)
	assert.Equal(t, 24, config.Duration)

	t.Logf("Config parsed and loaded correctly: %v", config)
}
