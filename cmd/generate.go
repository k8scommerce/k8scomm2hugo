/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/k8scommerce/k8scomm2hugo/internal/builder"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var (
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generates Hugo markdown files from K8sCommerce categories & products",
		Long: `Generates Hugo markdown files from K8sCommerce categories & products using 
the API and store-key.`,
		Run: func(cmd *cobra.Command, args []string) {
			storeKey, _ := cmd.Flags().GetString("storekey")
			apiURL, _ := cmd.Flags().GetString("endpoint")
			outputDir, _ := cmd.Flags().GetString("output")
			product, _ := cmd.Flags().GetString("product")
			category, _ := cmd.Flags().GetString("category")
			// fmt.Println("value of the flag apiURL :" + apiURL)
			// fmt.Println("value of the flag storekey :" + storeKey)
			// fmt.Println("value of the flag output :" + outputDir)
			// fmt.Println("value of the flag product :" + product)
			// fmt.Println("value of the flag category :" + category)

			builder.NewCategoryBuilder(apiURL, storeKey, outputDir, category).Build()
			builder.NewProductBuilder(apiURL, storeKey, outputDir, product).Build()
		},
	}
)

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().StringP("endpoint", "e", "", "Admin API endpoint")
	generateCmd.PersistentFlags().StringP("storekey", "k", "", "Store Key")
	generateCmd.PersistentFlags().StringP("output", "o", "./output", "Output directory")
	generateCmd.PersistentFlags().StringP("product", "p", "product", "Product basepath examples: product, products, etc. No trailing slash.")
	generateCmd.PersistentFlags().StringP("category", "c", "category", "Category basepath examples: category, categories, etc. No trailing slash.")

	generateCmd.MarkPersistentFlagRequired("endpoint")
	generateCmd.MarkPersistentFlagRequired("storekey")
}
