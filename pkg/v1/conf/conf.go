package conf

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description:
		Package conf provides a builder for declarative configuration loading with support for cascading overrides.
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/vanilla-os/sdk/pkg/v1/conf/types"
)

// Builder is a builder for configuration loading.
type Builder[T any] struct {
	domain    string
	confType  string
	prefix    string
	cascading bool
	optional  bool
}

// NewBuilder creates a new configuration builder for the given domain.
func NewBuilder[T any](domain string) *Builder[T] {
	return &Builder[T]{
		domain:    domain,
		cascading: true,
	}
}

// WithType sets the configuration file type (e.g. "json").
func (b *Builder[T]) WithType(t string) *Builder[T] {
	b.confType = t
	return b
}

// WithPrefix sets the prefix for the configuration paths (mostly for testing).
func (b *Builder[T]) WithPrefix(p string) *Builder[T] {
	b.prefix = p
	return b
}

// WithCascading configures whether to load configurations in a cascading manner
// (System -> User -> Local), merging them. If false, it stops at the first
// configuration found (Local -> User -> System). Default is true.
func (b *Builder[T]) WithCascading(enable bool) *Builder[T] {
	b.cascading = enable
	return b
}

// WithOptional configures whether the configuration file is optional. If true,
// Build will not return an error if no configuration file is found.
// Default is false.
func (b *Builder[T]) WithOptional(enable bool) *Builder[T] {
	b.optional = enable
	return b
}

// Build loads the configuration and returns it.
func (b *Builder[T]) Build() (*T, error) {
	if b.confType == "" {
		b.confType = "json"
	}
	if b.confType != "json" {
		return nil, fmt.Errorf("unsupported config type: %s", b.confType)
	}

	var config T
	paths := b.getOrderedPaths()
	loaded := false

	if b.cascading {
		for _, dir := range paths {
			path := filepath.Join(dir, "config."+b.confType)
			if err := loadFile(path, &config); err == nil {
				loaded = true
			}
		}
	} else {
		for i := len(paths) - 1; i >= 0; i-- {
			dir := paths[i]
			path := filepath.Join(dir, "config."+b.confType)
			if err := loadFile(path, &config); err == nil {
				loaded = true
				break
			}
		}
	}

	if !loaded {
		if b.optional {
			return &config, nil
		}
		return nil, errors.New("no configuration file found")
	}

	return &config, nil
}

func (b *Builder[T]) getOrderedPaths() []string {
	paths := []string{
		filepath.Join(b.prefix, "/usr/share", b.domain),
		filepath.Join(b.prefix, "/app/share", b.domain),
		filepath.Join(b.prefix, "/etc", b.domain),
	}

	u, err := user.Current()
	if err == nil {
		paths = append(paths, filepath.Join(b.prefix, u.HomeDir, b.domain))

		configHome := os.Getenv("XDG_CONFIG_HOME")
		if configHome == "" {
			configHome = filepath.Join(u.HomeDir, ".config")
		}
		paths = append(paths, filepath.Join(b.prefix, configHome, b.domain))
	}

	paths = append(paths, filepath.Join(".", "conf", b.domain))

	return paths
}

func loadFile(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}

// InitConfig is a compatibility wrapper using the new Builder.
// Deprecated: Use NewBuilder instead.
func InitConfig[T any](opts types.ConfigOptions) (*T, error) {
	return NewBuilder[T](opts.Domain).
		WithType(opts.Type).
		WithPrefix(opts.Prefix).
		Build()
}
