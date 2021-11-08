package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var Conf = viper.New()

var rootCmd = &cobra.Command{
	Use:   "yago",
	Short: "yago tools",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initailizeConfig(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func initailizeConfig(cmd *cobra.Command) error {
	v := Conf

	// v.SetConfigName("yacode")

	// v.AddConfigPath(".")

	// if err := v.ReadInConfig(); err != nil {
	// 	if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
	// 		return err
	// 	}
	// }

	// 设置前缀
	// v.SetEnvPrefix("GCODE")

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// if strings.Contains(f.Name, "-") {
		// 	envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
		// 	v.BindEnv(f.Name, fmt.Sprintf("%s_%s", "GCODE", envVarSuffix))
		// }

		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		} else {
			// 如果有 cmd，则用 cmd 值替换 conf
			if f.Value.Type() == "bool" {
				vv, _ := cmd.Flags().GetBool(f.Name)
				v.Set(f.Name, vv)
			} else if f.Value.Type() == "int64" {
				vv, _ := cmd.Flags().GetInt64(f.Name)
				v.Set(f.Name, vv)
			} else {
				v.Set(f.Name, f.Value)
			}
		}
	})

	return nil
}

func main() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(newCmd())
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(genDaoCmd())
	rootCmd.AddCommand(genModelCmd())
	rootCmd.AddCommand(genServiceCmd())
	rootCmd.AddCommand(genHttpCmd())
	rootCmd.AddCommand(genHttpMethodCmd())
	rootCmd.AddCommand(genServiceMethodCmd())
	rootCmd.AddCommand(genTypeCmd())
	rootCmd.AddCommand(genApiCmd())
	rootCmd.AddCommand(genTaskCmd())
	rootCmd.AddCommand(genCommandCmd())
	rootCmd.AddCommand(genRpcCmd())

	// rootCmd.PersistentFlags().StringP("file", "f", GetGoFile(), "file path,eg. ./ab_c.go")

	if err := rootCmd.Execute(); err != nil {
		log.Println("cmd run error:", err.Error())
		os.Exit(1)
	}
}
