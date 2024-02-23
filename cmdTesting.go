package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdTesting = &cobra.Command{
	Use:   "testing",
	Short: "Runs testing defined in template",
	Long:  "Runs testing defined in template",
	Run:   RunTesting,
}

func RunTesting(cmd *cobra.Command, args []string) {
	//var flags Config
	fmt.Printf("testrun named %s template %s\n", testingName, testingTemplate)

	default_config := parse_yaml("defaults.yml")
	config_template := parse_yaml("templates/" + testingTemplate + ".yml")

	// range thru globaltestparameters and add value from defaults.yml if !replace is set
	// TODO: skip possible double entries
	// these can happen when default.yml has same value
	for k, gtpa := range config_template.Testing.GlobalTestParameters {
		fmt.Printf("G  parameter: %s :", gtpa.Parameter)
		if !gtpa.Replace {
			//config_template.Testing.GlobalTestParameters[k].Values = nil
			// TODO P1 fix: append from correct defaults index
			config_template.Testing.GlobalTestParameters[k].Values = append(config_template.Testing.GlobalTestParameters[k].Values, default_config.Px_Version)
		}

		for _, tval := range gtpa.Values {
			fmt.Printf(" %s ", tval)
		}
		fmt.Printf("\n")
	}

	for _, gtpa := range config_template.Testing.GlobalTestParameters {
		fmt.Printf("G assembled parameter: %s :", gtpa.Parameter)
		for _, tval := range gtpa.Values {
			fmt.Printf(" %s ", tval)
		}
		fmt.Printf("\n")
	}

	for _, gtpl := range config_template.Testing.GlobalTestPlatform {
		fmt.Printf("G Plat: %s \n", gtpl.Platform)
		for _, tcpa := range gtpl.TestParameters {
			fmt.Printf("G Plat:   parameter: %s :", tcpa.Parameter)
			for _, tval := range tcpa.Values {
				fmt.Printf(" %s ", tval)
			}
			fmt.Printf("\n")
		}
	}

	for _, tcv := range config_template.Testing.TestClouds {
		fmt.Printf("cloud %s\n", tcv.Cloud)
		for _, tcp := range tcv.TestPlatforms {
			fmt.Printf("  platform %s\n", tcp.Platform)
			for _, tcpa := range tcp.TestParameters {
				fmt.Printf("    parameter: %s :", tcpa.Parameter)
				for _, tval := range tcpa.Values {
					fmt.Printf(" %s ", tval)
				}
				fmt.Printf("\n")
			}
		}
	}

	//prep_error := prepare_deployment(&config, &flags, testingName, "", testingTemplate, "")
	//if prep_error != "" {
	//	die(prep_error)
	//}
	//_ = create_deployment(config)

}
