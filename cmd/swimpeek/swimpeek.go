package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"swimpeek/internal/config"
	"swimpeek/internal/graph"
	"swimpeek/internal/lanedump"
	"swimpeek/internal/picker"
	"swimpeek/internal/tui"
	"swimpeek/pkg/laneclient"

	"github.com/charmbracelet/log"
)

var logger *log.Logger = config.GetLogger("swimpeek")
var version string = "NOTSET" // set by the build system

func printUsage() {
	fmt.Println("Usage: swimpeek <command> [options]")
	fmt.Println("Available commands:")
	fmt.Println("  config   - Create or modify the SwimPeek configuration.")
	fmt.Println("  dump     - Dump the tenant data to a file for analysis.")
	fmt.Println("  analyze  - Analyze the dumped tenant data.")
	fmt.Println("  version  - Show the SwimPeek version.")
	fmt.Println("Run 'swimpeek <command> -help' for more information on a specific command.")
}

func main() {
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {

		case "config":
			cmdConfig()

		case "dump":
			cmdDump(os.Args[2:])

		case "analyze":
			cmdAnalyze(os.Args[2:])

		case "version":
			logger.Info("swimpeek version: " + version)

		case "help":
			printUsage()

		default:
			logger.Error(fmt.Sprintf("Unknown command: %s", cmd))
			printUsage()
			os.Exit(1)
		}
		return
	}
	logger.Fatal("Please specify a command. Available commands: config, dump, analyze, version.")
}

// cmdConfig creates or modifies the SwimPeek configuration.
func cmdConfig() {
	cfg := loadConfig(true)

	// Test the connection to Swimlane
	logger.Info("Testing connection to Swimlane...")
	client := laneclient.NewLaneClient(cfg.FQDN(), cfg.SwimlaneAccountId, cfg.SwimlaneAccessToken, logger)
	ctx := context.Background()
	tenants, err := client.GetTenants(ctx)
	if err != nil {
		logger.Fatal("Failed to connect", "error", err)
	}
	if len(tenants.Tenants) == 0 {
		logger.Fatal("No tenants found for the account. Please check your configuration.")
	}
	logger.Info("Successfully connected to Swimlane!")
	for _, tenant := range tenants.Tenants {
		logger.Print("Tenant found", "name", tenant.Name, "id", tenant.Id, "users", tenant.UserCount)
	}

	logger.Info("Configuration initialized 🎉")
}

// cmdDump dumps the tenant data to a file for use with the analyze command.
func cmdDump(args []string) {
	outfile := ""
	tenantId := ""
	flagSet := flag.NewFlagSet("dump", flag.ExitOnError)
	flagSet.StringVar(&tenantId, "tenant", "", "Tenant ID to dump (if not specified, a picker dialog will be shown)")
	flagSet.StringVar(&outfile, "outfile", "", "Output file for the dump (default: lanedump_{tenant}.json)")
	if err := flagSet.Parse(args); err != nil {
		logger.Fatal("Failed to parse flags", "error", err)
	}

	cfg := loadConfig(false)

	// List the available tenants (implicitly testing the connection)
	client := laneclient.NewLaneClient(cfg.FQDN(), cfg.SwimlaneAccountId, cfg.SwimlaneAccessToken, config.GetLogger("laneclient"))
	ctx := context.Background()

	logger.Info("Fetching tenants...")
	tenants, err := client.GetTenants(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	tenant, err := selectTenant(tenants.Tenants, tenantId)
	if err != nil {
		logger.Fatal("Failed to select tenant", "error", err)
	}

	// If no output file is specified, use the name of the selected tenant.
	if outfile == "" {
		outfile = fmt.Sprintf("lanedump_%s.json", strings.ToLower(strings.ReplaceAll(tenant.Name, " ", "_")))
	}

	// Dump the tenant configuration to a file
	logger.Info("Dumping tenant configuration", "tenant", tenant.Name, "id", tenant.Id)
	tenantClient := laneclient.NewTenantClient(client, tenant)
	laneState, err := lanedump.LoadFromTenant(ctx, &tenantClient)
	if err != nil {
		logger.Fatal("Failed to dump tenant data", "error", err)
	}
	if err := lanedump.WriteToDisk(laneState, outfile); err != nil {
		logger.Fatal(err)
	}
	logger.Info("Tenant dumped successfully", "outfile", outfile)
}

// cmdAnalyze analyzes the dumped tenant data.
func cmdAnalyze(args []string) {
	infile := ""
	flagSet := flag.NewFlagSet("analyze", flag.ExitOnError)
	flagSet.StringVar(&infile, "infile", "", "Input file for the analysis")
	if err := flagSet.Parse(args); err != nil {
		logger.Fatal("Failed to parse flags", "error", err)
	}

	if infile == "" {
		flagSet.Usage()
		os.Exit(1)
	}

	// Load the tenant dump from disk and build the resource graph
	laneState, err := lanedump.LoadFromDisk(infile)
	if err != nil {
		logger.Fatal("Failed to load dump file", "error", err)
	}

	graph, warns, err := graph.FromState(laneState)
	if err != nil {
		logger.Fatal("Failed to create graph from lane state", "error", err)
	}
	for _, warn := range warns {
		logger.Warn(warn)
	}

	// Launch the resource browser
	if err := tui.LaunchExplorer(laneState, graph); err != nil {
		logger.Fatal("Failed to launch resource explorer", "error", err)
	}
}

// loadConfig loads the SwimPeek configuration from the default location.
// If newCfg is true, it initializes a new configuration.
func loadConfig(newCfg bool) *config.Config {
	// Load the configuration
	cfgDir, err := config.GetConfigDir(newCfg)
	if err != nil {
		logger.Error(err)
		logger.Warn("First run? Please run 'swimpeek config' to create a configuration file.")
		os.Exit(1)
	}

	if newCfg {
		// Run the config initialization
		cfg, err := config.InitConfig(cfgDir)
		if err != nil {
			logger.Fatal(err)
		}

		// Write the config
		if err := config.SaveConfig(cfgDir, cfg); err != nil {
			logger.Fatal(err)
		}
		return cfg
	}

	cfg, err := config.ReadConfig(cfgDir)
	if err != nil {
		logger.Fatal(err)

	}

	return cfg
}

// selectTenant shows a tenant picker dialog and returns the selected tenant.
// If tenantId is provided, it will return the tenant with that ID.
func selectTenant(tenants []laneclient.Tenant, tenantId string) (laneclient.Tenant, error) {
	if tenantId != "" {
		for _, t := range tenants {
			if t.Id == tenantId {
				return t, nil
			}
		}
		return laneclient.Tenant{}, fmt.Errorf("tenant with ID %s not found", tenantId)
	}

	return picker.PickTenant(tenants)
}
