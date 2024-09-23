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

In your application's simapp folder, integrate the following imports into the app.go file:

1. Imports

```sh

import (

    // ......

   	availblobkeeper "github.com/vitwit/avail-da-module/keeper"
	availblobmodule "github.com/vitwit/avail-da-module/module"
	availblobrelayer "github.com/vitwit/avail-da-module/relayer"
	"github.com/vitwit/avail-da-module/relayer/avail"
	httpclient "github.com/vitwit/avail-da-module/relayer/http"
	availtypes "github.com/vitwit/avail-da-module/types"
)

```

2. Constants configuration

After importing the necessary packages for the avail-da module in your app.go file, the next step is to declare any constant variables that the module will use. These constants are essential for configuring and integrating the CADA module with your application.

```sh
const (
	appName      = "avail-sdk"
	NodeDir      = ".availsdk"
)
```

3. Keeper and Relayer declaration

Here's a step-by-step guide to integrating the avail-da module keeper and relayer into your Cosmos SDK application

Inside of the ChainApp struct, add the required avail-da module runtime fields.

```sh
type SimApp struct {
    // ...

	AvailBlobKeeper  *availblobkeeper.Keeper
	Availblobrelayer *availblobrelayer.Relayer
	//

}
```

4. Initialize the `avail-da-module` Keeper and Relayer

Within the `NewSimApp` method, the constructor for the app, initialize the avail-da module components.

```go
    func NewSimApp(
	//...
    ) *SimApp {

        // ...

         // pre-existing code: remove optimistic execution in baseapp options
        baseAppOptions = append(baseAppOptions, voteExtOp)

        // NOTE: pre-existing code, add parameter.
            keys := storetypes.NewKVStoreKeys(
            // ...

            // Register avail-da module Store
            availblob1.StoreKey,
        )
    httpClient := httpclient.NewHandler()

        // Avail-DA client
        cfg := availtypes.AvailConfigFromAppOpts(appOpts)
        availDAClient := avail.NewLightClient(cfg.LightClientURL, httpClient)

        app.Availblobrelayer, err = availblobrelayer.NewRelayer(
            logger,
            appCodec,
            cfg,
            NodeDir,
            availDAClient,
        )
        if err != nil {
            panic(err)
        }

        app.AvailBlobKeeper = availblobkeeper.NewKeeper(
            appCodec,
            runtime.NewKVStoreService(keys[availtypes.StoreKey]),
            app.UpgradeKeeper,
            keys[availtypes.StoreKey],
            appOpts,
            logger,
            app.Availblobrelayer,
        )

        // must be done after relayer is created
        app.AvailBlobKeeper.SetRelayer(app.Availblobrelayer)

        //...

```

5.  Integrate Cada module\'s vote extensions and abci methods

    ```go
        voteExtensionHandler := availblobkeeper.NewVoteExtHandler(
            logger,
            app.AvailBlobKeeper,
        )

        dph := baseapp.NewDefaultProposalHandler(bApp.Mempool(), bApp)
        availBlobProposalHandler := availblobkeeper.NewProofOfBlobProposalHandler(
            app.AvailBlobKeeper,
            dph.PrepareProposalHandler(),
            dph.ProcessProposalHandler(),
            *voteExtensionHandler,
        )
        bApp.SetPrepareProposal(availBlobProposalHandler.PrepareProposal)
        bApp.SetProcessProposal(availBlobProposalHandler.ProcessProposal)
        bApp.SetExtendVoteHandler(voteExtensionHandler.ExtendVoteHandler())
        bApp.SetVerifyVoteExtensionHandler(voteExtensionHandler.VerifyVoteExtensionHandler())
        
    ```

    6. Module manager

    ```go

            // pre existing comments

            /**** Module Options ****/

            // ......

            // NOTE: pre-existing code, add parameter.
            app.ModuleManager = module.NewManager(
                // ...

                availblobmodule.NewAppModule(appCodec, app.AvailBlobKeeper),
            )

            // NOTE: pre-existing code, add parameter.
            app.ModuleManager.SetOrderBeginBlockers(
                // ...

                // avail-da-module begin blocker can be last
                availblob1.ModuleName,
            )

            // NOTE: pre-existing code, add parameter.
            app.ModuleManager.SetOrderEndBlockers(
                // ...

                // avail-da-module end blocker can be last
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

6. Integrate `avail-da-module` PreBocker

```go

    // PreBlocker application updates every pre block
    func (app *SimApp) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
        err := app.AvailBlobKeeper.PreBlocker(ctx, req)
        if err != nil {
            return nil, err
        }
        return app.ModuleManager.PreBlock(ctx)
    }

```

### Commands.go wiring

In your simapp application commands file, incorporate the following to wire up the avail-da module CLI commands.

1. Imports

Within the imported packages, add the avail-da module

```go
import (
    // ...
	"github.com/vitwit/avail-da-module/simapp/app"
    availblobcli "github.com/vitwit/avail-da-module/client/cli"
	"github.com/vitwit/avail-da-module/relayer"
)
```

2. Init App Config

````go
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

3. Init Root Command

```go

    func initRootCmd(
        rootCmd *cobra.Command,
        txConfig client.TxConfig,
        _ codectypes.InterfaceRegistry,
        _ codec.Codec,
        basicManager module.BasicManager,
    ) {
        // ......


        AddCommands(rootCmd, app.DefaultNodeHome, newApp, appExport, addModuleInitFlags)

        keysCmd := keys.Commands()
        keysCmd.AddCommand(availblobcli.NewKeysCmd())

        // add keybase, RPC, query, genesis, and tx child commands
        rootCmd.AddCommand(
            server.StatusCommand(),
            genesisCommand(txConfig, basicManager),
            queryCommand(),
            txCommand(),
            keysCmd,
            resetCommand(),
        )
    }
```
