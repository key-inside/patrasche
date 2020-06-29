// Copyright Key Inside Co., Ltd. 2020 All Rights Reserved.

package channel

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	fabcfg "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/cryptoutil"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient/seek"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/kataras/golog"
	"github.com/spf13/viper"

	"github.com/key-inside/patrasche/pkg/aws"
	"github.com/key-inside/patrasche/pkg/channel/event"
)

// Channel _
type Channel struct {
	sdk      *fabsdk.FabricSDK
	identity mspctx.SigningIdentity
	ctx      context.ChannelProvider
}

// type _organizationConfig string

// func (c _organizationConfig) Client() *mspctx.ClientConfig {
// 	return &mspctx.ClientConfig{Organization: string(c)}
// }

// New _
func New(ctxOpts ...fabsdk.ContextOption) (*Channel, error) {
	var cfg core.ConfigProvider
	arn, path, err := aws.GetARN("patrasche.network")
	if err != nil {
		cfg = fabcfg.FromFile(path)
	} else { // AWS resource
		in, typ, err := aws.GetReaderWithARN(arn)
		if err != nil {
			return nil, err
		}
		cfg = fabcfg.FromReader(in, typ)
	}

	// sdk
	// opts := []fabsdk.Option{}
	// if viper.IsSet("patrasche.organization") {
	// 	orgCfg := _organizationConfig(viper.GetString("patrasche.organization"))
	// 	opts = append(opts, fabsdk.WithIdentityConfig(orgCfg))
	// }
	// sdk, err := fabsdk.New(cfg, opts...)
	sdk, err := fabsdk.New(cfg)
	if err != nil {
		return nil, err
	}

	// msp
	si, err := getSigningIdentity(sdk.Context())
	if err != nil {
		return nil, err
	}

	// channel provider
	ctxOpts = append([]fabsdk.ContextOption{fabsdk.WithIdentity(si)}, ctxOpts...)
	ctx := sdk.ChannelContext(viper.GetString("patrasche.channel"), ctxOpts...)

	return &Channel{
		sdk:      sdk,
		identity: si,
		ctx:      ctx,
	}, nil
}

// Close _
func (c *Channel) Close() {
	if c.sdk != nil {
		c.sdk.Close()
	}
}

// NewClient returns a channel client
func (c *Channel) NewClient(options ...channel.ClientOption) (*channel.Client, error) {
	return channel.New(c.ctx, options...)
}

// NewEventClient returns an event client
func (c *Channel) NewEventClient(fromBlock uint64, seekType seek.Type) (*event.Client, error) {
	return event.New(c.ctx, fromBlock, seekType)
}

func getSigningIdentity(ctx context.ClientProvider) (mspctx.SigningIdentity, error) {
	client, err := mspclient.New(ctx)
	if err != nil {
		return nil, err
	}

	arn, nameOrPath, err := aws.GetARN("patrasche.identity")
	if err != nil {
		var err error
		si, err := client.GetSigningIdentity(nameOrPath)
		if err != nil {
			if err != mspclient.ErrUserNotFound {
				golog.Debugf("%#v", err)
				return nil, err
			}

			// msp style
			golog.Debug("identity from msp directory")

			mspDir := nameOrPath
			_ctx, err := ctx()
			if err != nil {
				return nil, err
			}
			cert, err := ioutil.ReadFile(filepath.Join(mspDir, "signcerts", "cert.pem"))
			if err != nil {
				return nil, err
			}
			pubKey, err := cryptoutil.GetPublicKeyFromCert(cert, _ctx.CryptoSuite())
			if err != nil {
				return nil, err
			}
			priKey, _ := ioutil.ReadFile(filepath.Join(mspDir, "keystore", fmt.Sprintf("%x_sk", pubKey.SKI())))

			return client.CreateSigningIdentity(mspctx.WithCert(cert), mspctx.WithPrivateKey(priKey))
		}
		// else
		golog.Debug("identity from store")
		return si, nil
	}

	// else AWS resource
	golog.Debug("identity from arn")

	in, typ, err := aws.GetReaderWithARN(arn)
	if err != nil {
		return nil, err
	}
	// uses viper for config consistency
	temp := viper.New()
	temp.SetConfigType(typ)
	if err := temp.ReadConfig(in); err != nil {
		return nil, err
	}
	cert := []byte(temp.GetString("cert"))
	priKey := []byte(temp.GetString("key"))

	return client.CreateSigningIdentity(mspctx.WithCert(cert), mspctx.WithPrivateKey(priKey))
}
