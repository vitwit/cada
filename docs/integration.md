# Integration

Follow these steps to integrate the avail-da module into your Cosmos SDK-based application.


### app.go wiring 

In your application's simapp folder, integrate the following imports into the app.go file:

1. Imports 

```sh

import (

    // ......

   "github.com/vitwit/avail-da-module"
   availblobkeeper "github.com/vitwit/avail-da-module/keeper"
   availblobmodule "github.com/vitwit/avail-da-module/module"
   availblobrelayer "github.com/vitwit/avail-da-module/relayer" 
)

```

2. Constants configuration

After importing the necessary packages for the avail-da module in your app.go file, the next step is to declare any constant variables that the module will use. These constants are essential for configuring and integrating the avail-da module with your application.

```sh
const (
	// TODO: Change me
	AvailAppID = 1

	// publish blocks to avail every n rollchain blocks.
	publishToAvailBlockInterval = 5 // smaller size == faster testing
)
```

3. Keeper and Relyer declaration

Here's a step-by-step guide to integrating the avail-da module keeper and relayer into your Cosmos SDK application

Inside of the ChainApp struct, add the required avail-da module runtime fields.

```sh
type SimApp struct {
    // ...

	AvailBlobKeeper  *availblobkeeper.Keeper
	Availblobrelayer *availblobrelayer.Relayer
	//

}```

4. Initialize the `avail-da-module` Keeper and Relayer

Within the `NewSimApp` method, the constructor for the app, initialize the avail-da module components.

```sh
    func NewSimApp(
	//
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

        // must be done after relayer is created
        app.AvailBlobKeeper.SetRelayer(app.Availblobrelayer)

        dph := baseapp.NewDefaultProposalHandler(bApp.Mempool(), bApp)
        availBlobProposalHandler := availblobkeeper.NewProofOfBlobProposalHandler(app.AvailBlobKeeper, dph.PrepareProposalHandler(), dph.ProcessProposalHandler())
        bApp.SetPrepareProposal(availBlobProposalHandler.PrepareProposal)
        bApp.SetProcessProposal(availBlobProposalHandler.ProcessProposal)

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

5. Integrate `avail-da-module` PreBocker

```sh

// PreBlocker application updates every pre block
func (app *SimApp) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	err := app.AvailBlobKeeper.PreBlocker(ctx, req)
	if err != nil {
		return nil, err
	}
	return app.ModuleManager.PreBlock(ctx)
}

```

6. Integrate relayer startup

To integrate the relayer startup into your Cosmos SDK application, you will need to query necessary values and initialize the relayer. Here’s how you can do it:

* Modify RegisterNodeService Function :
In your app.go file, locate the RegisterNodeService function. You need to add code to initialize and start the relayer after your application has started.

* Add the Relayer Initialization: Inside the RegisterNodeService function, you will need to query necessary values from the application and initialize the relayer. 

Here’s how you can do it: 

```sh

func (app *SimApp) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)

	app.Availblobrelayer.SetClientContext(clientCtx)

	go app.Availblobrelayer.Start()
}
```

### Commands.go wiring 

In your simapp application commands file, incorporate the following to wire up the avail-da module CLI commands.

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

3. Init Root Command

```sh

func initRootCmd(
	rootCmd *cobra.Command,
	txConfig client.TxConfig,
	interfaceRegistry codectypes.InterfaceRegistry,
	appCodec codec.Codec,
	basicManager module.BasicManager,
) {
        // ......


       server.AddCommands(rootCmd, simapp.DefaultNodeHome, newApp, appExport, addModuleInitFlags)

	    keysCmd := keys.Commands()
	    keysCmd.AddCommand(availblobcli.NewKeysCmd())

        rootCmd.AddCommand(
		server.StatusCommand(),
		genesisCommand(txConfig, basicManager),
		queryCommand(),
		txCommand(),
		keysCmd,
	)
```