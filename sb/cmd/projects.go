package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

type Project struct {
	UUID        string `json:"uuid"`
	Name        string `json:"display_name"`
	Description string `json:"description"`
	User        int    `json:"user"`
	Cluster     string `json:"cluster"`
}

func fetchProjects() ([]Project, error) {

	sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		fmt.Println("User not logged in")
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/projects/", sbUrl), nil)

token, err := GetToken(sbUrl)
if err != nil {
    fmt.Printf("error getting token: %v\n", err)
}

	/*token := viper.GetString("token")
	if token == "" {
		fmt.Println("User not logged in")
	}
*/
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func checkExistingProject(name, sbUrl, token string) error {
	url := fmt.Sprintf("%s/api/projects/", sbUrl)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the necessary headers
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the response
	var projects []ProjectCreate
	if err := json.Unmarshal(body, &projects); err != nil {
		return fmt.Errorf("failed to parse response body: %w", err)
	}

	// Check if the project name already exists
	for _, p := range projects {
		if p.Name == name {
			return fmt.Errorf("project name already exists")
		}
	}

	return nil
}


var projectsCmd = &cobra.Command{
	Use:   "projects",
	Aliases: []string{"project"}, 
	Short: "Projects are loaded namespaces within a cluster.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify an action like list or add.")
	},
}

func init() {
	rootCmd.AddCommand(projectsCmd)
}
