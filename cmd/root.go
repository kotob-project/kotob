/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	toLang   string
	fromLang string
	apiKey   string
	model    string
	system   string
	asJson   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kotob [text]",
	Short: "A lightweight CLI translation tool powered by Gemini API",
	Long: `Kotob is a lightweight CLI translation tool built with Go,
leveraging the Google Gemini API for fast and accurate translations.`,

	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// テスト用に引数とフラグを表示してみる
		fmt.Printf("翻訳対象: %s\n", args[0])
		fmt.Printf("出力言語: %s\n", toLang)
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().StringVarP(&toLang, "to", "t", "", "target language (defaults to en ⇔ ja if unspecified)")
	rootCmd.Flags().StringVarP(&fromLang, "from", "f", "auto", "source language")
	rootCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "Gemini API key for the session")
	rootCmd.Flags().StringVarP(&model, "model", "m", "gemini-2.0-flash-lite", "AI model to use")
	rootCmd.Flags().StringVarP(&system, "system", "s", "", "custom system prompt for the AI")
	rootCmd.Flags().BoolVarP(&asJson, "json", "j", false, "output result as a JSON object")
}
