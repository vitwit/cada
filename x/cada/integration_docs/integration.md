# Integration

Follow these steps to integrate the cada module into your Cosmos SDK-based application.

### app.go wiring

In your application's simapp folder, integrate the following imports into the app.go file:

1. Imports

```sh

import (

    // ......

	cadakeeper "github.com/vitwit/avail-da-module/keeper"
	cadamodule "github.com/vitwit/avail-da-module/module"
	cadarelayer "github.com/vitwit/avail-da-module/relayer"
	"github.com/vitwit/avail-da-module/relayer/avail"
	httpclient "github.com/vitwit/avail-da-module/relayer/http"
	cadatypes "github.com/vitwit/avail-da-module/types"
)

```

2. Constants configuration

After importing the necessary packages for the cada in your app.go file, the next step is to declare any constant variables that the module will use. These constants are essential for configuring and integrating the cada module with your application.

```sh
const (
	appName      = "cada-sdk"
	NodeDir      = ".cada"
)
```

3. Keeper and Relayer declaration

Here's a step-by-step guide to integrating the cada module keeper and relayer into your Cosmos SDK application

Inside of the ChainApp struct, add the required cada module runtime fields.

```sh
type SimApp struct {
    // ...

	CadaKeeper  *cadakeeper.Keeper
	Cadarelayer *cadarelayer.Relayer
	//

}
```

4. Initialize the `Cada` Keeper and Relayer

Within the `NewSimApp` method, the constructor for the app, initialize the cada module components.

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

            // Register cada module Store
            cadatypes.StoreKey,
        )

    httpClient := httpclient.NewHandler()

    // Avail-DA client
        cfg := cadatypes.AvailConfigFromAppOpts(appOpts)
        availDAClient := avail.NewLightClient(cfg.LightClientURL, httpClient)

        app.Cadarelayer, err = cadarelayer.NewRelayer(
            logger,
            appCodec,
            cfg,
            NodeDir,
            availDAClient,
        )
        if err != nil {
            panic(err)
        }

        app.CadaKeeper = cadakeeper.NewKeeper(
            appCodec,
            runtime.NewKVStoreService(keys[cadatypes.StoreKey]),
            app.UpgradeKeeper,
            keys[cadatypes.StoreKey],
            appOpts,
            logger,
            app.Cadarelayer,
        )

        // must be done after relayer is created
        app.CadaKeeper.SetRelayer(app.Cadarelayer)

        //...

```

5.  Integrate Cada module\'s vote extensions and abci methods

    ```go
        voteExtensionHandler := cadakeeper.NewVoteExtHandler(
            logger,
            app.CadaKeeper,
        )

        dph := baseapp.NewDefaultProposalHandler(bApp.Mempool(), bApp)
        cadaProposalHandler := cadakeeper.NewProofOfBlobProposalHandler(
            app.CadaKeeper,
            dph.PrepareProposalHandler(),
            dph.ProcessProposalHandler(),
            *voteExtensionHandler,
        )
        bApp.SetPrepareProposal(cadaProposalHandler.PrepareProposal)
        bApp.SetProcessProposal(cadaProposalHandler.ProcessProposal)
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

                cadamodule.NewAppModule(appCodec, app.CadaKeeper),
            )

            // NOTE: pre-existing code, add parameter.
            app.ModuleManager.SetOrderBeginBlockers(
                // ...

                // cada begin blocker can be last
                cadatypes.ModuleName,
            )

            // NOTE: pre-existing code, add parameter.
            app.ModuleManager.SetOrderEndBlockers(
                // ...

                // cada end blocker can be last
                cadatypes.ModuleName,
            )

                // NOTE: pre-existing code, add parameter.
            genesisModuleOrder := []string{
                // ...

                // cada genesis module order can be last
                cadatypes.ModuleName,
            }

        }

    )
    ```

6. Integrate `cada` PreBocker

```go

    // PreBlocker application updates every pre block
    func (app *SimApp) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
        err := app.CadaKeeper.PreBlocker(ctx, req)
        if err != nil {
            return nil, err
        }
        return app.ModuleManager.PreBlock(ctx)
    }

```

### Commands.go wiring

In your simapp application commands file, incorporate the following to wire up the cada module CLI commands.

1. Imports

Within the imported packages, add the cada module

```go
import (
    // ...
	"github.com/vitwit/avail-da-module/simapp/app"
    cadacli "github.com/vitwit/avail-da-module/client/cli"
	cadatypes "github.com/vitwit/avail-da-module/types"
)
```

2. Init App Config

````go
func initAppConfig() (string, interface{}) {

	type CustomAppConfig struct {
		serverconfig.Config

		Cada *cadatypes.AvailConfiguration `mapstructure:"avail"`
	}

    // ...

	customAppConfig := CustomAppConfig{
		Config: *srvCfg,
		Avail:  &cadatypes.DefaultAvailConfig,
	}

	customAppTemplate := serverconfig.DefaultConfigTemplate + cadatypes.DefaultConfigTemplate

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
        keysCmd.AddCommand(cadacli.NewKeysCmd())

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
