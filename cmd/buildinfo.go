package cmd

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/choria-io/go-choria/build"
	"github.com/choria-io/go-choria/config"
	"github.com/choria-io/go-choria/provtarget"
	"github.com/choria-io/go-protocol/protocol"
	gnatsd "github.com/nats-io/gnatsd/server"
)

type buildinfoCommand struct {
	command
}

func (b *buildinfoCommand) Setup() (err error) {
	b.cmd = cli.app.Command("buildinfo", "View build settings")

	return
}

func (b *buildinfoCommand) Configure() (err error) {
	cfg, err = config.NewDefaultConfig()
	if err != nil {
		return fmt.Errorf("Could not create default configuration: %s", err)
	}

	cfg.DisableSecurityProviderVerify = true
	cfg.Choria.SecurityProvider = "file"

	return
}

func (b *buildinfoCommand) Run(wg *sync.WaitGroup) (err error) {
	defer wg.Done()

	fmt.Println("Choria build settings:")
	fmt.Println()
	fmt.Println("Build Data:")
	fmt.Printf("     Version: %s\n", build.Version)
	fmt.Printf("     Git SHA: %s\n", build.SHA)
	fmt.Printf("  Build Date: %s\n", build.BuildDate)
	fmt.Printf("     License: %s\n", build.License)
	fmt.Printf("  Go Version: %s\n", runtime.Version())
	fmt.Println()
	fmt.Println("Network Broker Settings:")
	fmt.Printf("       Maximum Network Clients: %d\n", build.MaxBrokerClients())
	fmt.Printf("  Embedded NATS Server Version: %s\n", gnatsd.VERSION)

	mutators := config.MutatorNames()
	fmt.Println()
	fmt.Println("Configuration Mutators:")
	if len(mutators) > 0 {
		for _, m := range mutators {
			fmt.Printf("\t%s\n", m)
		}
	} else {
		fmt.Printf("\tnone\n")
	}

	fmt.Println()
	fmt.Println("Server Settings:")
	fmt.Printf("            Provisioning Brokers: %s\n", build.ProvisionBrokerURLs)
	fmt.Printf("            Provisioning Default: %t\n", build.ProvisionDefault())
	fmt.Printf("                Provisioning TLS: %t\n", build.ProvisionSecurity())
	fmt.Printf("      Default Provisioning Agent: %t\n", build.ProvisionAgent == "true")
	fmt.Printf("  Provisioning Registration Data: %s\n", build.ProvisionRegistrationData)
	fmt.Printf("              Provisioning Facts: %s\n", build.ProvisionFacts)
	fmt.Printf("    Provisioning Target Resolver: %s\n", provtarget.Name())
	if build.ProvisionToken != "" {
		fmt.Printf("              Provisioning Token: set\n")
	} else {
		fmt.Printf("              Provisioning Token: not set\n")
	}
	fmt.Println()
	fmt.Println("Agent Providers:")

	for _, p := range build.AgentProviders {
		fmt.Printf("\t%s\n", p)
	}

	fmt.Println()
	fmt.Println("Security Defaults:")
	fmt.Printf("            TLS: %s\n", build.TLS)
	fmt.Printf("  x509 Security: %t\n", protocol.IsSecure())

	if build.TLS != "true" || !protocol.IsSecure() {
		fmt.Println()
		fmt.Println("NOTE: The security of this build is non standard, you might be running without adequate protocol level security.  Please ensure this is the build you intend to be using.")
	}

	return
}

func init() {
	cli.commands = append(cli.commands, &buildinfoCommand{})
}
