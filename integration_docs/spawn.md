# Spawn Integration

## Integration

Follow these steps to integrate the avail-da module into your application. You can use spawn to create a new application, simply include the avail-da module to have it pre-wired in your application.

A fully integrated example application is available in this repository under the simapp directory which can be used as a reference.

### Installing Spawn

Clone this repo and install

```sh
git clone https://github.com/rollchains/spawn.git
cd spawn
git checkout v0.50.4
make install
```

Create your chain using the spawn command and customize it to your needs!

```sh
spawn new rollchain --bech32=cosmos --bin=simd --denom=token 
```

for more details about spawn you can refer to this doc https://github.com/rollchains/spawn

### app.go wiring

In your main application file, typically named `app.go`  incorporate the following to wire up the avail-da module

1. Imports

Within the imported packages, add the `avail-da-module` dependencies.

```sh

import (

    // ......

   "github.com/vitwit/avail-da-module"
   availblobkeeper "github.com/vitwit/avail-da-module/keeper"
   availblobmodule "github.com/vitwit/avail-da-module/module"
   availblobrelayer "github.com/vitwit/avail-da-module/relayer" 
)

```

2.Configuration constants.

After the imports, declare the constant variables required by `avail-da-module`.

```sh
const (
	// TODO: Change me
	AvailAppID = 1

	// publish blocks to avail every n rollchain blocks.
	publishToAvailBlockInterval = 5 // smaller size == faster testing
)
```

3. Keeper and Relyer declaration

Inside of the ChainApp struct, the struct which satisfies the cosmos-sdk runtime.AppI interface, add the required avail-da module runtime fields.

```sh

type ChainApp struct {
	// ...

	AvailBlobKeeper  *availblobkeeper.Keeper
	Availblobrelayer *availblobrelayer.Relayer
	// ....
}
```
4. Initialize the `avail-da-module` Keeper and Relayer

Within the `NewChainApp` method, the constructor for the app, initialize the avail-da module components.

```sh
func NewChainApp(
    // ...
) *ChainApp {
        // ...

        // NOTE: pre-existing code, add parameter.
        keys := storetypes.NewKVStoreKeys(
            // ...

            // Register avail-da module Store
            availblob1.StoreKey,
        )

        app.AvailBlobKeeper = availblobkeeper.NewKeeper(
            appCodec,
            appOpts,
            runtime.NewKVStoreService(keys[availblob1.StoreKey]),
            app.UpgradeKeeper,
            keys[availblob1.StoreKey],
            publishToAvailBlockInterval,
            AvailAppID,
        )

        app.Availblobrelayer, err = availblobrelayer.NewRelayer(
            logger,
            appCodec,
            appOpts,
            homePath,
        )
        if err != nil {
            panic(err)
        }

        // Connect relayer to keeper. Must be done after relayer is created.
        app.AvailBlobKeeper.SetRelayer(app.Availblobrelayer)

        // Rollchains avail-da-module proposal handling
        availBlobProposalHandler := availblobkeeper.NewProofOfBlobProposalHandler(app.AvailBlobKeeper,
            AppSpecificPrepareProposalHandler(), // i.e. baseapp.NoOpPrepareProposal()
            AppSpecificProcessProposalHandler(), // i.e. baseapp.NoOpProcessProposal()
        )
        bApp.SetPrepareProposal(availBlobProposalHandler.PrepareProposal)
        bApp.SetProcessProposal(availBlobProposalHandler.ProcessProposal)

        // ...

        // NOTE: pre-existing code, add parameter.
        app.ModuleManager = module.NewManager(
            // ...

        availblobmodule.NewAppModule(appCodec, app.AvailBlobKeeper),
        )

        // NOTE: pre-existing code, add parameter.
        app.ModuleManager.SetOrderBeginBlockers(
            // ...

            // avail-da module begin blocker can be last
            availblob1.ModuleName,
        )

        // NOTE: pre-existing code, add parameter.
        app.ModuleManager.SetOrderEndBlockers(
            // ...

            // avail-da module end blocker can be last
            availblob1.ModuleName,
        )

        // NOTE: pre-existing code, add parameter.
        genesisModuleOrder := []string{
            // ...

            // avail-da genesis module order can be last
        availblob1.ModuleName,
        }

    }
)
```

5. Integrate Relayer into FinalizeBlock

The `avail-da-module` relayer needs to be notified when the rollchain blocks are committed so that it is aware of the latest height of the chain. In `FinalizeBlock`, rather than returnig app.BaseApp.FinalizeBlock(req), add error handling to that call and notify the relayer afterward.

```sh
func (app *ChainApp) FinalizeBlock(req *abci.RequestFinalizeBlock) (*abci.ResponseFinalizeBlock, error) {
    // ...

	res, err := app.BaseApp.FinalizeBlock(req)
	if err != nil {
		return res, err
	}

	app.Availblobrelayer.NotifyCommitHeight(req.Height)

	return res, nil
}
```

6. Integrate `avail-da-module` PreBocker

The `avail-da-module` PreBlocker must be called in the app's PreBlocker.

```sh
func (app *ChainApp) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	err := app.AvailBlobKeeper.PreBlocker(ctx, req)
	if err != nil {
		return nil, err
	}
	return app.ModuleManager.PreBlock(ctx)
}
```

7. Integrate relayer startup

The relayer needs to query blocks using the client context in order to package them and publish to Avail Light Client. The relayer also needs to be started with some initial values that must be queried from the app after the app has started. Add the following in RegisterNodeService.

```sh

func (app *ChainApp) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
    app.Availblobrelayer.SetClientContext(clientCtx)

	go app.Availblobrelayer.Start()
}

```

### Commands.go wiring

In your application commands file, incorporate the following to wire up the avail-da module CLI commands.

1. Imports

Within the imported packages, add the avail-da module

```sh
import (
    // ...
	"github.com/vitwit/avail-da-module/simapp/app"
    availblobcli "github.com/vitwit/avail-da-module/client/cli"
	"github.com/vitwit/avail-da-module/relayer"
)
```

2. Init App Config


```sh
func initAppConfig() (string, interface{}) {

	type CustomAppConfig struct {
		serverconfig.Config

		Avail *relayer.AvailConfig `mapstructure:"avail"`
	}

    // ...

	customAppConfig := CustomAppConfig{
		Config: *srvCfg,
		Avail:  &relayer.DefaultAvailConfig,
	}

	customAppTemplate := serverconfig.DefaultConfigTemplate + relayer.DefaultConfigTemplate 

	return customAppTemplate, customAppConfig
}
```

3.Init Root Command

```sh

func initRootCmd(
    // ...
) {
	// ...

	keysCmd := keys.Commands()
	keysCmd.AddCommand(availblobcli.NewKeysCmd())


    // Existing code, only modifying one parameter.
	rootCmd.AddCommand(
		server.StatusCommand(),
		genesisCommand(txConfig, basicManager),
		queryCommand(),
		txCommand(),
		keysCmd, // replace keys.Commands() here with this
	)
}

```

