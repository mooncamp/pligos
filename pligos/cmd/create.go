// Copyright © 2019 real.digital
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"realcloud.tech/cloud-tools/pkg/pligos"
	"realcloud.tech/cloud-tools/pkg/pligos/helm"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "this command creates a helm chart",
	Run: func(cmd *cobra.Command, args []string) {
		h, err := newHelm()
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "init helm: %v\n", err)
			os.Exit(1)
		}

		if err := h.Create(chartPath); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "create: %v\n", err)
			os.Exit(1)
		}
	},
}

func newHelm() (*helm.Helm, error) {
	config, err := pligos.OpenPligosConfig(configPath)
	if err != nil {
		return nil, err
	}

	c, err := newCompiler()
	if err != nil {
		return nil, err
	}

	context := pligos.FindContext(contextName, config.Contexts)
	return helm.New(context.Flavor, config.DeploymentConfig, config.ChartDependencies, context.Configs, context.Secrets, c), nil
}

var chartPath string

func init() {
	createCmd.Flags().StringVarP(&chartPath, "output", "o", "", "output directory")
	createCmd.MarkFlagRequired("output")

	rootCmd.AddCommand(createCmd)
}
