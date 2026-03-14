/*
Copyright © 2026 kotob-project contributors

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
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kotob-project/kotob/pkg/translate"
)

var (
	toLang   string
	fromLang string
	apiKey   string
	model    string
	system   string
	asJson   bool
	noStream bool
)

type TranslationResponse struct {
	Source     string `json:"source"`
	Target     string `json:"target"`
	Input      string `json:"input"`
	Translated string `json:"translated"`
	Model      string `json:"model"`
}

var rootCmd = &cobra.Command{
	Use:   "kotob [flags] [text]",
	Short: "A lightweight CLI translation tool powered by Gemini API",
	Long: `Kotob is a lightweight CLI translation tool built with Go,
leveraging the Google Gemini API for fast and accurate translations.`,

	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// チェック
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "ERROR: Please input the text to be translated.")
			os.Exit(1)
		}

		if apiKey == "" {
			fmt.Fprintln(os.Stderr, "ERROR: API key is not configured.")
			os.Exit(1)
		}

		//翻訳準備
		ctx := context.Background()
		client, err := translate.NewClient(ctx, apiKey, model)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}

		//翻訳開始

		if asJson || noStream {
			result, err := client.Translate(ctx, args[0], fromLang, toLang, system)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				os.Exit(1)
			}
			if asJson {
				resp := TranslationResponse{
					Source:     fromLang,
					Target:     toLang,
					Input:      args[0],
					Translated: result,
					Model:      model,
				}
				encoder := json.NewEncoder(os.Stdout)
				encoder.SetIndent("", "  ")
				encoder.Encode(resp)
			} else {
				fmt.Print(result)
			}
		} else {
			err = client.TranslateStream(ctx, os.Stdout, args[0], fromLang, toLang, system)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				os.Exit(1)
			}
		}

		fmt.Println()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().StringVarP(&toLang, "to", "t", "", "target language (defaults to en ⇔ ja if unspecified)")
	rootCmd.Flags().StringVarP(&fromLang, "from", "f", "auto", "source language (default auto)")
	rootCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "Gemini API key for the session")
	rootCmd.Flags().StringVarP(&model, "model", "m", "gemini-2.5-flash-lite", "AI model to use")
	rootCmd.Flags().StringVarP(&system, "system", "s", "", "custom system prompt for the AI")
	rootCmd.Flags().BoolVarP(&asJson, "json", "j", false, "output result as a JSON object")
	rootCmd.Flags().BoolVarP(&noStream, "no-stream", "S", false, "Outputs translations in bulk")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetEnvPrefix("KOTOB")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// 設定ファイル関連
	viper.SetConfigName("kotob")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(home + "/.config/kotob")
	}

	_ = viper.ReadInConfig()

	// 補完処理
	if apiKey == "" {
		apiKey = viper.GetString("api-key")
	}
	if model == "" || model == "gemini-2.5-flash-lite" {
		vModel := viper.GetString("model")
		if vModel != "" {
			model = vModel
		}
	}
	if toLang == "" {
		vtoLang := viper.GetString("to")
		if vtoLang != "" {
			toLang = vtoLang
		} else {
			toLang = "Japanese"
		}
	}
	if fromLang == "" || fromLang == "auto" {
		vfromLang := viper.GetString("from")
		if vfromLang != "" {
			fromLang = vfromLang
		}
	}
	if system == "" {
		system = viper.GetString("system")
	}
	if !asJson {
		asJson = viper.GetBool("json")
	}
}
