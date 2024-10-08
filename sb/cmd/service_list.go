package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var servicelistCmd = &cobra.Command{
	Use:   "list",
	Short: "List all services.",
	Run:   listServices,
}

func listServices(cmd *cobra.Command, args []string) {
	// TODO: if cluster context is set, list all projects in cluster
	services, err := fetchServices()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching projects: %v\n", err)
		return
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"UUID", "Name", "Project","Type"})
	for _, service := range services {
		t.AppendRow([]interface{}{service.UUID, service.Name, service.Project.DisplayName,service.Type })
		t.AppendSeparator()
	}
	t.AppendSeparator()
	t.Render()
	if err != nil {
		fmt.Println("Unable to parse response")
	}
}

func init() {
	servicesCmd.AddCommand(servicelistCmd)
}
