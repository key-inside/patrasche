/*
SPDX-FileCopyrightText: © 2023 Key Inside Co., Ltd.

SPDX-License-Identifier: BSD-3-Clause
*/

package patrasche

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/key-inside/patrasche/aws"
	"github.com/key-inside/patrasche/block"
	"github.com/key-inside/patrasche/channel"
	"github.com/key-inside/patrasche/listener"
	"github.com/key-inside/patrasche/logger"
)

type Patrasche struct {
	logOut io.Writer
	logger zerolog.Logger

	envPrefix string
	cfgName   string
	cfgPrefix string

	mutex          sync.Mutex
	configViper    *viper.Viper
	configProvider func() *Config
	config         *Config
}

func Biter(cmd *cobra.Command) *Patrasche {
	// sub commands do not inherit context before execute
	for p := cmd; p != nil; p = p.Parent() {
		if ctx := p.Context(); ctx != nil {
			return FromContext(ctx)
		}
	}
	return nil
}

type ctxKey string

const patrascheCtxKey ctxKey = "patrasche"

func FromContext(ctx context.Context) *Patrasche {
	if ctx != nil {
		v := ctx.Value(patrascheCtxKey)
		if v != nil {
			p := v.(*Patrasche)
			p.mutex.Lock()
			defer p.mutex.Unlock()
			if p.config == nil {
				p.config = p.configProvider()
			}
			return p
		}
	}
	return nil
}

func Bite(cmd *cobra.Command, options ...Option) (*Patrasche, error) {
	out := cmd.OutOrStdout()
	p := &Patrasche{
		logOut:    out,
		logger:    logger.New("patrasche", out, zerolog.InfoLevel),
		envPrefix: cmd.Name(),
		cfgName:   "config",
		cfgPrefix: "patrasche",
	}
	for _, option := range options {
		if err := option(p); err != nil {
			return nil, fmt.Errorf("failed to apply patrasche option: %w", err)
		}
	}

	if p.configProvider == nil {
		// sets up env, flags ...
		v := viper.New()
		v.SetEnvPrefix(p.envPrefix)
		v.AutomaticEnv()
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		ns := p.cfgPrefix // namespace
		if ns != "" {
			ns += "."
		}
		persistentFlags := newDefaultFlagSet(ns)
		if p.cfgName != "" {
			persistentFlags.StringSlice(p.cfgName, []string{}, "config sources")
		}
		v.BindPFlags(persistentFlags)
		cmd.PersistentFlags().AddFlagSet(persistentFlags)

		// pre logging level before loading config
		if lv := v.GetString(ns + "logging.level"); lv != "" {
			p.SetLogLevel(lv)
		}

		p.configViper = v

		p.configProvider = func() *Config {
			cfg, err := p.configFromViper()
			if err != nil {
				p.logger.Panic().Err(err).Msg("failed to load config")
			}
			p.SetLogLevel(cfg.Logging.Level)
			return cfg
		}
	}

	if err := overrideMethods(cmd); err != nil {
		return nil, fmt.Errorf("failed to override command's methods: %w", err)
	}

	// adds patrasche to context
	// it can be got via Biter or FromContext functions
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	cmd.SetContext(context.WithValue(ctx, patrascheCtxKey, p))

	return p, nil
}

func (p *Patrasche) Logger() zerolog.Logger {
	return p.logger
}

func (p *Patrasche) SetLogLevel(lv string) {
	zLv, _ := zerolog.ParseLevel(lv)
	p.logger = p.logger.Level(zLv)
}

func (p *Patrasche) ConfigMap() map[string]any {
	if p.configViper != nil {
		return p.configViper.AllSettings()
	}
	return map[string]any{}
}

func newDefaultFlagSet(namespace string) *pflag.FlagSet {
	fset := pflag.NewFlagSet("patrasche", pflag.ContinueOnError)

	// logging
	fset.String(namespace+"logging.level", "info", "fatal, error, warn, info, debug or disable")

	// fabric
	fset.String(namespace+"fabric.envPrefix", "", "prefix for environment variable of Fabric SDK")
	fset.String(namespace+"fabric.channel", "", "channel name")
	fset.String(namespace+"fabric.identity.organization", "", "user organization")
	fset.String(namespace+"fabric.identity.username", "", "username in the user store")
	fset.StringSlice(namespace+"fabric.config", []string{}, "Fabric SDK config sources")

	return fset
}

func overrideMethods(cmd *cobra.Command) error {
	tmpl := template.New("version")
	tmpl, err := tmpl.Parse(cmd.VersionTemplate())
	if err != nil {
		return fmt.Errorf("can't parse command version template: %w", err)
	}
	printVersion := tmpl.Execute

	persistentPreRunE := cmd.PersistentPreRunE
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) (err error) {
		out := cmd.OutOrStdout()
		PrintLogo(out)
		printVersion(out, cmd)
		fmt.Fprintln(out)

		if p := Biter(cmd); p != nil {
			p.logger.Debug().Str("command", cmd.CommandPath()).
				Dict("runtime", zerolog.Dict().
					Str("os", runtime.GOOS).
					Str("arch", runtime.GOARCH).
					Str("ver", runtime.Version()),
				).Msg("")
		}

		if persistentPreRunE != nil {
			err = persistentPreRunE(cmd, args)
		} else if cmd.PersistentPreRun != nil {
			cmd.PersistentPreRun(cmd, args)
		}

		return err
	}

	if !cmd.Runnable() {
		cmd.Run = func(cmd *cobra.Command, args []string) {
			cmd.Help()
			fmt.Fprintln(cmd.OutOrStdout())
		}
	}

	persistentPostRunE := cmd.PersistentPostRunE
	cmd.PersistentPostRunE = func(cmd *cobra.Command, args []string) (err error) {
		if persistentPostRunE != nil {
			err = persistentPostRunE(cmd, args)
		} else if cmd.PersistentPostRun != nil {
			cmd.PersistentPostRun(cmd, args)
		}
		fmt.Fprintln(cmd.OutOrStdout())
		return err
	}

	return nil
}

func (p *Patrasche) configFromViper() (*Config, error) {
	v := p.configViper

	// forces comma separation
	cfgs := []string{}
	if err := v.UnmarshalKey(p.cfgName, &cfgs); err != nil {
		return nil, err
	}

	for _, src := range cfgs {
		if arn.IsARN(src) {
			p.logger.Debug().Str("arn", src).Msg("load config from AWS")
			cfgMap, err := aws.GetConfigMap(src)
			if err != nil {
				return nil, fmt.Errorf("failed to get config: %w", err)
			}
			if err := v.MergeConfigMap(cfgMap); err != nil {
				return nil, fmt.Errorf("can't merge config: %w", err)
			}
		} else {
			p.logger.Debug().Str("filepath", src).Msg("load config file")
			v.SetConfigFile(src)
			if err := v.MergeInConfig(); err != nil {
				return nil, fmt.Errorf("can't merge config: %w", err)
			}
		}
	}

	cfg, err := unmarshalConfig(v, p.cfgPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return cfg, err
}

func unmarshalConfig(v *viper.Viper, name string) (*Config, error) {
	fabricEnvPrefixKey := "fabric.envPrefix"

	var cfg Config
	if name != "" {
		typ := reflect.StructOf([]reflect.StructField{
			{
				Name: "Config",
				Type: reflect.TypeOf(Config{}),
				Tag:  reflect.StructTag(`mapstructure:"` + name + `"`),
			},
		})
		wrapCfg := reflect.New(typ).Interface()
		if err := v.Unmarshal(&wrapCfg); err != nil {
			return nil, err
		}
		cfg = reflect.ValueOf(wrapCfg).Elem().FieldByName("Config").Interface().(Config)

		fabricEnvPrefixKey = name + "." + fabricEnvPrefixKey // add namespace
	} else {
		if err := v.Unmarshal(&cfg); err != nil {
			return nil, err
		}
	}

	// sets default values
	if !v.IsSet(fabricEnvPrefixKey) {
		cfg.Fabric.EnvPrefix = "FABRIC_SDK"
	}

	return &cfg, nil
}

type fabricConfigBackend struct {
	configViper *viper.Viper
}

func (c *fabricConfigBackend) Lookup(key string) (interface{}, bool) {
	v := c.configViper.Get(key)
	return v, (v != nil)
}

func (p *Patrasche) NewChannel(ctxOpts ...fabsdk.ContextOption) (*channel.Channel, error) {
	v := viper.New()
	v.SetEnvPrefix(p.config.Fabric.EnvPrefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	for _, src := range p.config.Fabric.Config {
		if arn.IsARN(src) {
			p.logger.Debug().Str("arn", src).Msg("load fabric config from AWS")
			cfgMap, err := aws.GetConfigMap(src)
			if err != nil {
				return nil, fmt.Errorf("failed to get fabric config: %w", err)
			}
			if err := v.MergeConfigMap(cfgMap); err != nil {
				return nil, fmt.Errorf("can't merge fabric config: %w", err)
			}
		} else {
			p.logger.Debug().Str("filepath", src).Msg("load fabric config from file")
			v.SetConfigFile(src)
			if err := v.MergeInConfig(); err != nil {
				return nil, fmt.Errorf("can't merge fabric config: %w", err)
			}
		}
	}

	backend := &fabricConfigBackend{configViper: v}
	ctx := func() ([]core.ConfigBackend, error) {
		return []core.ConfigBackend{backend}, nil
	}

	if p.config.Fabric.Identity.Username != "" {
		opts := []fabsdk.ContextOption{fabsdk.WithUser(p.config.Fabric.Identity.Username)}
		if p.config.Fabric.Identity.Organization != "" {
			opts = append(opts, fabsdk.WithOrg(p.config.Fabric.Identity.Organization))
		}
		ctxOpts = append(opts, ctxOpts...)
	}

	return channel.New(p.config.Fabric.Channel, ctx, ctxOpts...)
}

func (p *Patrasche) ListenBlock(handler block.Handler, options ...listener.Option) error {
	ch, err := p.NewChannel()
	if err != nil {
		return fmt.Errorf("failed to connect channel: %w", err)
	}
	defer ch.Close()

	l, err := listener.New(ch, handler, options...)
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}
	return l.Listen()
}

type Option func(*Patrasche) error

func WithEnvPrefix(prefix string) Option {
	return func(p *Patrasche) error {
		p.envPrefix = prefix
		return nil
	}
}

// WithConfigFlagName sets the config flag name
// ex, 'app.config.path', default is 'config'
func WithConfigFlagName(name string) Option {
	return func(p *Patrasche) error {
		p.cfgName = name
		return nil
	}
}

// WithConfigPrefix sets the namespace of the patrasche config set
// Default is 'patrasche' and it can be empty string.
func WithConfigPrefix(name string) Option {
	return func(p *Patrasche) error {
		p.cfgPrefix = name
		return nil
	}
}

// WithConfigProvider sets the custom config provider function
// that runs when a Patrasche instance is first retrieved by Biter or FromContext
func WithConfigProvider(provider func() *Config) Option {
	return func(p *Patrasche) error {
		p.configProvider = provider
		return nil
	}
}

// WithLogWriter injects a writer into the global logger
func WithLogWriter(w io.Writer) Option {
	return func(p *Patrasche) error {
		p.logOut = w
		p.logger = p.logger.Output(w)
		return nil
	}
}

// WithConsoleLogWriter injects a zerolog.ConsoleWriter into the global logger
func WithConsoleLogWriter() Option {
	return func(p *Patrasche) error {
		return WithLogWriter(logger.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.Out = p.logOut
		}))(p)
	}
}
