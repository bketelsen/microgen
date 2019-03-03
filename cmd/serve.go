// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		http.HandleFunc("/", generate)

		fmt.Printf("Starting server...\n")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	},
}

func generate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		return
	case "POST":
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Read Error", http.StatusInternalServerError)
			return
		}
		fname := r.Header.Get("X-Filename")
		if fname == "" {
			fname = "submitted.proto"
		}
		fmt.Println(string(b))
		err = process(b, fname)

		if err != nil {

			fmt.Println(err.Error())
			http.Error(w, "Process Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Not Supported", http.StatusMethodNotAllowed)
	}
}

func process(bb []byte, name string) error {

	var (
		cmdOut []byte
		err    error
	)
	wd, err := os.Getwd()
	fname := filepath.Join(wd, name)
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	_, err = f.Write(bb)

	if err != nil {
		return err
	}
	//protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. path/to/greeter.proto

	cmdName := "/usr/local/bin/protoc"
	cmdArgs := []string{"--proto_path=/home/bketelsen/src:.", "--micro_out=.", "--go_out=.", fname}
	fmt.Println(cmdName, cmdArgs)
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running generation command: ", err)
		fmt.Println(string(cmdOut))
		return err
	}
	return nil
}
func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
