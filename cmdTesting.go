package main

import (
	"fmt"
	"reflect"
	"slices"

	"strings"

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
	fmt.Printf("testrun named %s template %s\n\n", testingName, testingTemplate)

	default_config := parse_yaml("defaults.yml")
	config_template := parse_yaml("templates/" + testingTemplate + ".yml")

	// range thru globaltestparameters and add value from defaults.yml if !replace is set
	// TODO: skip possible double entries
	//       this can happen when default.yml has same value
	fmt.Printf("Global Settings for test run: (entries with * being added from defaults.yml, [] indicates double entries)\n")
	for k, gtpa := range config_template.Testing.GlobalTestParameters {
		fmt.Printf("globaltestparameter: %s :", gtpa.Parameter)
		for _, tval := range gtpa.Values {
			fmt.Printf(" %s", tval)
		}
		if !gtpa.Ignoredefault {
			defVal := getDefaultValue(gtpa.Parameter, &default_config)
			if defVal != "" {
				if !slices.Contains(config_template.Testing.GlobalTestParameters[k].Values, defVal) {
					fmt.Printf(" %s*", defVal)
					config_template.Testing.GlobalTestParameters[k].Values = append(config_template.Testing.GlobalTestParameters[k].Values, defVal)
				} else {
					fmt.Printf(" [%s*]", defVal)
				}
			}
		}
		fmt.Printf("\n")
	}

	for k, globalplatforms := range config_template.Testing.GlobalTestPlatform {
		fmt.Printf("globaltestplatform : %s\n", globalplatforms.Platform)
		for j, globalplatformvalues := range globalplatforms.TestParameters {
			fmt.Printf("  %s:", globalplatformvalues.Parameter)
			for _, tval := range globalplatformvalues.Values {
				fmt.Printf(" %s", tval)
			}
			if !globalplatformvalues.Ignoredefault {
				defVal := getDefaultValue(globalplatformvalues.Parameter, &default_config)
				if defVal != "" {
					if !slices.Contains(config_template.Testing.GlobalTestPlatform[k].TestParameters[j].Values, defVal) {
						fmt.Printf(" %s*", defVal)
						config_template.Testing.GlobalTestPlatform[k].TestParameters[j].Values = append(config_template.Testing.GlobalTestPlatform[k].TestParameters[j].Values, defVal)
					} else {
						fmt.Printf(" [%s*]", defVal)
					}
				}
				fmt.Printf("\n")
			}
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

func getDefaultValue(field string, config *Config) string {
	refConf := reflect.ValueOf(*config)
	typeOfC := refConf.Type()
	for i := 0; i < refConf.NumField(); i++ {
		if strings.ToLower(typeOfC.Field(i).Name) == strings.ToLower(field) {
			//fmt.Printf("found field %s in defaults\n", refConf.Field(i).Interface())
			return fmt.Sprintf("%s", refConf.Field(i).Interface())
		}
	}
	return ""
}
