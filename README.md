![A Dog of Flanders][patrasche-logo] <small>*Logo design by `Summer Ham`*</small>  
![Go Version][go-version-image]

## What is Patrasche?

Patrasche is the utility package for Hyperledger Fabric DApp.  
It is suitable for use with [spf13/cobra][github-cobra]

## Hyperledger Fabirc Compatibility

### Current Compatibility

Integration tested Fabric versions:
- v1.4.12
> v2.x is not tested, but maybe it will be compatible too. (｡•̀ᴗ-)

### Fabric SDK version

- [fabric-sdk-go v1.0.1-0.20221020141211-7af45cede6af][github-fabirc-sdk-go]

> `IMPORTANT`  
> Currently *fabric-sdk-go* is not compatible with *google.golang.org/grpc@v1.43.0 or higher* in *allow-insecure* mode.  
> So, if you want to test locally without TLS, (go mod) replace *google.golang.org/grpc* version to *1.42.0*
```go
// go.mod example
module example

go 1.18

replace google.golang.org/grpc => google.golang.org/grpc v1.42.0 // for allow-insecure
```

## Contents

- [Examples](#Examples)
- [Configuration](#Configuration)
- [Usage](#Usage)
- [Test](#Test)

## Examples

> See test codes in `./test/` and example cli programs in `./cmd/` directory.

```go
package dapp

import (
    "github.com/spf13/cobra"
    "github.com/key-inside/patrasche"
    "github.com/key-inside/patrasche/cmd/ccquery"
    "github.com/key-inside/patrasche/cmd/inspect"
)

func main() {
    c := &cobra.Command{
        Use:     "dapp",
        Short:   "Patrasche DApp",
        Long:    "HLF block listener and chaincode query tool",
        Version: "v0.0.0",
    }
    c.AddCommand(inspect.Command(), ccquery.Command())

    // bite 'c' at any time before executing
    // you can add sub-commands to 'c' after being bitten
    patrasche.Bite(c)

    c.Execute()
}
```

## Configuration

* Patrasche configuration set via config file, environment variable, AWS resource or flag
* If the value is [ARN][arn-doc], it is treated as an AWS resource.
* Currently, AWS resource supports only ParameterStore(SSM) and SecretsManager.
* Patrasche supports multiple config files and resources.

```sh
--config="./config.yaml,./config-stg.json"
```

```yaml
patrasche:
  fabric:
    config:
      - "./fixtures/fabric.yaml"
      - "./fixtures/userstore.yaml"
      - "arn:aws:secretsmanager:::secret:stg/org1/userstore.yaml"
```

* Since Patrasche use the [spf13/viper][github-viper], extensions of file path and ARN resource are used as the content type.
* (viper) Supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env" and "dotenv".
* If there is no extension, it is regarded as "json".

> The command name is used as an environment variable prefix.

```sh
export DAPP_PATRASCHE_LOGGING_LEVEL=info
```

> You can change the env prefix.

```go
func WithEnvPrefix(prefix string) Option
// WithEnvPrefix("MARS") => export MARS_PATRASCHE_LOGGING_LEVEL=info
// WithEnvPrefix("") => export PATRASCHE_LOGGING_LEVEL=info
```

> The default config resource path flag is `config`

```sh
% dapp inspect --config=./stg/config.yaml
% export DAPP_CONFIG=./config.yaml
% dapp inspect
```

> You can change the config flag.

```go
func WithConfigFlagName(name string) Option
// WithConfigFlagName("extra") => --extra.config=./config.yaml
// WithConfigFlagName("") => No 'config' flag!
```

> The default prefix for patrasche configuration is `patrasche`

```sh
% dapp inspect --patrasche.fabric.channel=flanders
% export DAPP_PATRASCHE_FABRIC_CHANNEL=flanders
% dapp inspect
```

> You can change the patrasche config prefix.

```go
func WithConfigPrefix(name string) Option
// WithConfigPrefix("extra") => --extra.fabric.channel=flanders
// WithConfigPrefix("") => --fabric.channel=flanders
```

### ConfigProvider

* When bite a command, Patrasche creates and binds a config provider.
* Config provider function runs when a Patrasche instance is first retrieved by Biter or FromContext.
* You can change the config provider when bite a command.

```go
func WithConfigProvider(provider func() *Config) Option
```

## Usage

### Bite cobra!

* Bite transforms the cobra.Command into carrier of Patrasche.
* It adds configration flags and provider.  
* It overrides some xxxRun methods of the command.
* It adds the Patrasche instance to the context.

```go
func Bite(cmd *cobra.Command, options ...Option) (*Patrasche, error)
```

### Biter

* The bitten command and descendant commands can retrieve biter.

```go
func Biter(cmd *cobra.Command) *Patrasche
func FromContext(ctx context.Context) *Patrasche
```

```go
cobra.Command{
     PreRun: func(cmd *cobra.Command, args []string) {
        // some codes...
        p := patrasche.FromContext(cmd.Context())
        p.SetLogLevel("debug")
        // some codes...
    },
    Run: func(cmd *cobra.Command, args []string) {
        // some codes...
        p := patrasche.Biter(cmd)
        err := p.LitenBlock(handler, options...)
        // some codes...
    }
}
```

### Block Listener

* You can create a block event listener with the function below.

```go
func NewListener(ch *channel.Channel, handler Handler, options ...ListenerOption) (*Listener, error)
```

* Or you can easily listen using the method below.

```go
func (p *Patrasche) ListenBlock(handler block.Handler, options ...block.ListenerOption) error 
```

* Also you can write your own listener code.

> Listener options

```go
func WithStartBlock(blockNum uint64) ListenerOption
func WithEndBlock(blockNum uint64) ListenerOption
func WithShutdown(shutdown func(os.Signal)) ListenerOption
```

### Handler

* To handle block events and transactions, you must implement handlers.
* Handlers can be used in a variety of ways through nesting.
* Common handlers are provided as presets.

> Presets for block handler

```go
// standard block handler to handle tx
func NewStdHandler(handler tx.Handler) (Handler, error)
// standard block logging middleware
func NewStdLogger(next Handler, logger *zerolog.Logger) Handler
// filters
func NewHashFilter(next Handler, pattern string, filteredActions ...Action) Handler
// block number writer
func NewBlockNumberFileWriter(next Handler, path string) Handler
func NewBlockNumberDynamoDBWriter(next Handler, awsCfg aws.Config, table string, itemFactory func(*Block) any) Handler 
```

> Presets for tx handler

```go
// standard tx logging middleware
func NewStdLogger(next Handler, logger *zerolog.Logger) Handler
// filters
func NewHashFilter(next Handler, pattern string, filteredActions ...Action) Handler
func NewValidEndorserFilter(next Handler, filteredActions ...Action) Handler
```

### Action

* Action is a special function for handling filtered objects in filter handlers.

> Presets for block filter

```go
func NewHashFilteredLoggingAction(logger *zerolog.Logger) Action
```

> Presets for tx filter

```go
func NewHashFilteredLoggingAction(logger *zerolog.Logger) Action
func NewValidEndorserFilteredLoggingAction(logger *zerolog.Logger) Action
```

### Logging

* Patrasche currently uses [zerolog][github-zerolog].
* Patrasche does not use the global logger.

```go
// methods related to logger
func (p *Patrasche) Logger() zerolog.Logger
func (p *Patrasche) SetLogLevel(lv string)

// Bite options related to logging
func WithLogWriter(w io.Writer) Option
func WithConsoleLogWriter() Option
```

## Test

> You need to create config files in `'test/fixtures/'` before testing.

```sh
% go test -v ./test
% go test -v ./test --run Test_Help
% make test
```

<!-- references /-->
[patrasche-logo]:   ./logo-1.jpg
[go-version-image]: https://img.shields.io/badge/go%20version-%3E=1.18-61CFDD.svg?style=flat-square

[arn-doc]: https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html

[github-fabirc-sdk-go]: https://github.com/hyperledger/fabric-sdk-go/tree/7af45cede6afa3939a9574bc9948cca9fb424257
[github-cobra]:   https://github.com/spf13/cobra
[github-viper]:   https://github.com/spf13/viper
[github-zerolog]: https://github.com/rs/zerolog
