/*
Copyright © 2019 Ian Soboroff (ian.soboroff@nist.gov)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	homedir "github.com/mitchellh/go-homedir"
)


var cfgFile string


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dedupe",
	Short: "Compute fingerprints for documents for deduplication.",
	Long: `Dedupe is a library for computing hashes on documents, to enable
clustering near duplicates.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dedupe.yaml)")

	// lsh.buckets is the number of hash buckets for LSH, see lib/lsh.go
	rootCmd.PersistentFlags().Int32P("lsh.buckets", "l", 256, "number of buckets for locality-sensitive hashing (default is 256)")
	viper.BindPFlag("lsh.buckets", rootCmd.PersistentFlags().Lookup("lsh.buckets"))
	viper.SetDefault("lsh.buckets", "256")
	
	// lsh.threshold is the target Jaccard similarity threshold, see lib/lsh.go
	rootCmd.PersistentFlags().Float32P("lsh.threshold", "t", 0.9, "target Jaccard similarity threshold (default 0.9)")
	viper.BindPFlag("lsh.threshold", rootCmd.PersistentFlags().Lookup("lsh.threshold"))
	viper.SetDefault("lsh.threshold", "0.9")

	// minhash.size is the dimensionality of the minhash fingerprint algorithm.  See lib/minhash.go
	rootCmd.PersistentFlags().Int32P("minhash.size", "m", 256, "number of minhash function (default 256)")
	viper.BindPFlag("minhash.size", rootCmd.PersistentFlags().Lookup("minhash.size"))
	viper.SetDefault("minhash.size", "256")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in current or home directory with name ".dedupe" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigName(".dedupe")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

