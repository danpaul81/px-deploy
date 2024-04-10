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
	fmt.Printf("global parameter: (* inherit defaults.yml, [] double entries)\n")
	for k, gtpa := range config_template.Testing.GlobalTestParameters {
		fmt.Printf("  %s :", gtpa.Parameter)
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

	fmt.Printf("\nplatform specific settings: (+ inherit global parameter, * inherit from defaults.yml, [] double entries)\n")
	for k, globalplatforms := range config_template.Testing.GlobalTestPlatform {
		fmt.Printf(" %s\n", globalplatforms.Platform)
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
			}

			for _, gtpa := range config_template.Testing.GlobalTestParameters {
				if strings.EqualFold(gtpa.Parameter, globalplatformvalues.Parameter) {
					//fmt.Printf(" match %s", gtpa.Parameter)
					for _, gtpv := range gtpa.Values {
						if !slices.Contains(config_template.Testing.GlobalTestPlatform[k].TestParameters[j].Values, gtpv) {
							fmt.Printf(" %s+", gtpv)
							config_template.Testing.GlobalTestPlatform[k].TestParameters[j].Values = append(config_template.Testing.GlobalTestPlatform[k].TestParameters[j].Values, gtpv)
						} else {
							fmt.Printf(" [%s+]", gtpv)
						}
					}
				}
			}
			fmt.Printf("\n")
		}
	}

	for _, tcv := range config_template.Testing.TestClouds {
		//fmt.Printf("\n cloud %s\n", tcv.Cloud)
		for _, tcp := range tcv.Platforms {
			fmt.Printf("%s\t%s\n", tcv.Cloud, tcp)

			// check if for this platform there is a specific override
			// e.g. px_version set in global test parameters AND in global platform setting
			// if yes, take the values from the global platform settings
			// otherwise take from global setting
			//for _, gtpa := range config_template.Testing.GlobalTestParameters {
			//	fmt.Printf(" glob %s\n", gtpa.Parameter)
			for _, gppa := range config_template.Testing.GlobalTestPlatform {
				if gppa.Platform == tcp {
					for _, gpval := range gppa.TestParameters {

						fmt.Printf("    setting %s\n", gpval.Parameter)
						for _, gtpa := range config_template.Testing.GlobalTestParameters {
							fmt.Printf("   g setting %s\n", gtpa.Parameter)
						}
						//if gpval.Parameter == gtpa.Parameter {
						//	fmt.Printf("      found g %s override for platform %s\n", gtpa.Parameter, gppa.Platform)
						//		platform_override = true
						//}
					}
				}
			}

			//}
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
		if strings.EqualFold(typeOfC.Field(i).Name, field) {
			//fmt.Printf("found field %s in defaults\n", refConf.Field(i).Interface())
			return fmt.Sprintf("%s", refConf.Field(i).Interface())
		}
	}
	return ""
}
