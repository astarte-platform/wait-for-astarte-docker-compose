/*
Copyright © 2020-2023 SECO Mind Srl

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
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wait-for-astarte-docker-compose",
	Short: "A utility to wait for a docker-compose Astarte deployment",
	Long: `This program checks the health endpoints of a standard Astarte
docker-compose deployment and returns 0 when all Astarte services are up, or
1 if there is a timeout.

This can be used in scripts where a docker-compose Astarte instance is launched
and the rest of the script has to wait for Astarte to be ready`,
	RunE: rootExecF,
}

var timeoutSeconds int

var serviceToBaseUrl = map[string]string{
	"Realm Management API": "http://api.astarte.localhost/realmmanagement",
	"Housekeeping API":     "http://api.astarte.localhost/housekeeping",
	"AppEngine API":        "http://api.astarte.localhost/appengine",
	"Pairing API":          "http://api.astarte.localhost/pairing",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&timeoutSeconds, "timeout", "t", 300, "Timeout in seconds. Defaults to 300 seconds.")
}

func checkService(wg *sync.WaitGroup, name string, baseUrl string) {
	defer wg.Done()

	url := baseUrl + "/health"
	ok := false
	for !ok {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			fmt.Printf("%s is ready\n", name)
			ok = true
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}

func waitForFinished(wg *sync.WaitGroup, c chan struct{}) {
	wg.Wait()
	c <- struct{}{}
}

func rootExecF(cmd *cobra.Command, args []string) error {
	var wg sync.WaitGroup

	fmt.Println("Waiting for Astarte docker-compose...")

	for k, v := range serviceToBaseUrl {
		wg.Add(1)
		go checkService(&wg, k, v)
	}

	c := make(chan struct{})

	go waitForFinished(&wg, c)

	timeout := time.Duration(timeoutSeconds) * time.Second

	select {
	case <-c:
		fmt.Println("Astarte is ready")
		break
	case <-time.After(timeout):
		fmt.Println("Timed out")
		os.Exit(1)
	}

	return nil
}
