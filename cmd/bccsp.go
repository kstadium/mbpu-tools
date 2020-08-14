package cmd

import (
  "fmt"
  "errors"
  "os"
  "plugin"
  "strconv"

  "github.com/hyperledger/fabric/bccsp"
  "github.com/hyperledger/fabric/bccsp/factory"
  "github.com/spf13/cobra"
)

var (
  cmdBCCSP = &cobra.Command{
    Use:   "bccsp [command]",
    Args: cobra.MinimumNArgs(1),
    Short: "BCCSP tool",
  }

  cmdBCCSPVersion = &cobra.Command{
    Use:   "version [subcommand]",
    Args: cobra.MinimumNArgs(1),
    Short: "Print BCCSP version ",
    Run: func(cmd *cobra.Command, args []string) {
      err := versionPlugin(args[0])
      if err != nil {
        fmt.Println(err)
      }
    },
  }

  cmdBCCSPTest = &cobra.Command{
    Use:   "test [MBPU device index] [BCCSP plugin path]",
    Args: cobra.MinimumNArgs(2),
    Short: "Test BCCSP Plugin",
    Run: func(cmd *cobra.Command, args []string) {
      testPlugin(args[0], args[1])
    },
  }
)

func versionPlugin(filepath string) (err error) {
  	// make sure the library exists
	if _, err = os.Stat(filepath); err != nil {
    return
	}

	// attempt to load the library as a plugin
	plug, err := plugin.Open(filepath)
	if err != nil {
    return
	}

	// lookup the required symbol 'New'
	sym, err := plug.Lookup("PrintVersion")
	if err != nil {
    return
	}

	// check to make sure symbol New meets the required function signature
	version, ok := sym.(func() (string))
	if !ok {
    err = fmt.Errorf("Plugin does not implement the required function signature for 'PrintVersion'")
    return
	}

  version()
  
  return nil
}
func testPlugin(idxStr string, filepath string){
  count := getMBPUCount()

  idx, err := strconv.Atoi(idxStr)
  if err != nil {
     fmt.Printf("Invalid MBPU Index")
     return
  }

  if idx >= count {
    fmt.Printf("Out of index range")
    return
  }
  
  
  opts := &factory.FactoryOpts{
    PluginOpts: &factory.PluginOpts{
      Library: filepath,
      Config: map[string]interface{}{
        "SecLevel":     "256",
        "HashFamily":   "SHA2",
        "KeyStorePath": "./",
        "BoardIndex": idx,
      },
    },
  }
  bsp, err := loadPlugin(opts)
  if err != nil {
    // TODO kind fail message
    fmt.Println(err)
  }

  if bsp != nil {
    fmt.Print("plugin loaded\n")
  }

  bsp = nil

}


func loadPlugin(config *factory.FactoryOpts) (bccsp.BCCSP, error){
	// check for valid config
	if config == nil || config.PluginOpts == nil {
		return nil, errors.New("Invalid config. It must not be nil.")
	}

	// Library is required property
	if config.PluginOpts.Library == "" {
		return nil, errors.New("Invalid config: missing property 'Library'")
	}

	// make sure the library exists
	if _, err := os.Stat(config.PluginOpts.Library); err != nil {
		return nil, fmt.Errorf("Could not find library '%s' [%s]", config.PluginOpts.Library, err)
	}

	// attempt to load the library as a plugin
	plug, err := plugin.Open(config.PluginOpts.Library)
	if err != nil {
		return nil, fmt.Errorf("Failed to load plugin '%s' [%s]", config.PluginOpts.Library, err)
	}

	// lookup the required symbol 'New'
	sym, err := plug.Lookup("New")
	if err != nil {
		return nil, fmt.Errorf("Could not find required symbol 'CryptoServiceProvider' [%s]", err)
	}

	// check to make sure symbol New meets the required function signature
	new, ok := sym.(func(config map[string]interface{}) (bccsp.BCCSP, error))
	if !ok {
		return nil, fmt.Errorf("Plugin does not implement the required function signature for 'New'")
	}

	return new(config.PluginOpts.Config)
}