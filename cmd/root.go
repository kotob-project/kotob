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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kotob-project/kotob/pkg/translate"
)

const (
	defaultModel          = "gemini-2.5-flash-lite"
	defaultToLang         = "Japanese"
	defaultFromLang       = "auto"
	configName            = "kotob"
	configType            = "json"
	configDir             = "/.config/kotob"
	envPrefix             = "KOTOB"
	maxFileSize           = 1024 * 1024 // 1MB
	binaryCheckBufferSize = 1024
)

type Config struct {
	ToLang   string
	FromLang string
	APIKey   string
	Model    string
	System   string
	Filepath string
	AsJSON   bool
	NoStream bool
}

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
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := buildConfig()

		inputText, err := readInputText(cfg.Filepath, args)
		if err != nil {
			return err
		}

		if err := validateConfigAndInput(cfg, inputText); err != nil {
			return err
		}

		ctx := context.Background()
		return executeTranslation(ctx, cfg, inputText)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("to", "t", defaultToLang, "target language")
	rootCmd.Flags().StringP("from", "f", defaultFromLang, "source language")
	rootCmd.Flags().StringP("api-key", "k", "", "Gemini API key for the session")
	rootCmd.Flags().StringP("model", "m", defaultModel, "AI model to use")
	rootCmd.Flags().StringP("system", "s", "", "custom system prompt for the AI")
	rootCmd.Flags().BoolP("json", "j", false, "output result as a JSON object")
	rootCmd.Flags().StringP("file", "F", "", "path to the text file to translate")
	rootCmd.Flags().BoolP("no-stream", "S", false, "Outputs translations in bulk")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(".")
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(home + configDir)
	}

	_ = viper.ReadInConfig()

	_ = viper.BindPFlag("to", rootCmd.Flags().Lookup("to"))
	_ = viper.BindPFlag("from", rootCmd.Flags().Lookup("from"))
	_ = viper.BindPFlag("api-key", rootCmd.Flags().Lookup("api-key"))
	_ = viper.BindPFlag("model", rootCmd.Flags().Lookup("model"))
	_ = viper.BindPFlag("system", rootCmd.Flags().Lookup("system"))
	_ = viper.BindPFlag("json", rootCmd.Flags().Lookup("json"))
	_ = viper.BindPFlag("file", rootCmd.Flags().Lookup("file"))
	_ = viper.BindPFlag("no-stream", rootCmd.Flags().Lookup("no-stream"))
}

func buildConfig() *Config {
	toLang := viper.GetString("to")
	if toLang == "" {
		toLang = defaultToLang
	}

	return &Config{
		ToLang:   toLang,
		FromLang: viper.GetString("from"),
		APIKey:   viper.GetString("api-key"),
		Model:    viper.GetString("model"),
		System:   viper.GetString("system"),
		Filepath: viper.GetString("file"),
		AsJSON:   viper.GetBool("json"),
		NoStream: viper.GetBool("no-stream"),
	}
}

func readInputText(filepath string, args []string) (string, error) {
	if filepath != "" {
		file, err := os.Open(filepath)
		if err != nil {
			return "", fmt.Errorf("failed to open file %s: %w", filepath, err)
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			return "", fmt.Errorf("failed to get file info: %w", err)
		}
		if info.IsDir() {
			return "", fmt.Errorf("%q is a directory", filepath)
		}

		// バイナリチェック
		buffer := make([]byte, binaryCheckBufferSize)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("failed to read file: %w", err)
		}

		if bytes.IndexByte(buffer[:n], 0) != -1 {
			return "", fmt.Errorf("binary files are not supported")
		}

		remaining, err := io.ReadAll(file)
		if err != nil {
			return "", fmt.Errorf("failed to read remaining file content: %w", err)
		}
		data := append(buffer[:n], remaining...)

		if len(data) > maxFileSize {
			fmt.Fprintln(os.Stderr, "WARNING: File is very large. It may take some time or reach API limits.")
		}

		if len(args) > 0 {
			fmt.Fprintln(os.Stderr, "WARNING: Both file and direct text provided. Using file content.")
		}

		return string(data), nil
	}

	if len(args) < 1 {
		return "", fmt.Errorf("please input the text to be translated or specify a file with -F")
	}

	return args[0], nil
}

func validateConfigAndInput(cfg *Config, inputText string) error {
	if strings.Contains(cfg.FromLang, ".") || strings.Contains(cfg.FromLang, "/") || strings.Contains(cfg.FromLang, "\\") {
		fmt.Fprintln(os.Stderr, "WARNING: -f is for source language. Did you mean -F for file path?")
	}

	if cfg.APIKey == "" {
		return fmt.Errorf("API key is not configured")
	}

	if strings.TrimSpace(inputText) == "" {
		return fmt.Errorf("input text is empty")
	}

	return nil
}

func executeTranslation(ctx context.Context, cfg *Config, inputText string) error {
	client, err := translate.NewClient(ctx, cfg.APIKey, cfg.Model)
	if err != nil {
		return fmt.Errorf("failed to create translation client: %w", err)
	}

	if cfg.AsJSON || cfg.NoStream {
		result, err := client.Translate(ctx, inputText, cfg.FromLang, cfg.ToLang, cfg.System)
		if err != nil {
			return fmt.Errorf("translation error: %w", err)
		}
		if cfg.AsJSON {
			resp := TranslationResponse{
				Source:     cfg.FromLang,
				Target:     cfg.ToLang,
				Input:      inputText,
				Translated: result,
				Model:      cfg.Model,
			}
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(resp); err != nil {
				return fmt.Errorf("failed to encode JSON response: %w", err)
			}
		} else {
			fmt.Print(result)
		}
	} else {
		err = client.TranslateStream(ctx, os.Stdout, inputText, cfg.FromLang, cfg.ToLang, cfg.System)
		if err != nil {
			return fmt.Errorf("stream translation error: %w", err)
		}
	}

	fmt.Println()
	return nil
}
