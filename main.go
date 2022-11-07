/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import "github.com/qleet/qleetctl/cmd"

//go:generate bash get_version.sh
func main() {
	cmd.Execute()
}
