package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var goBasicCmd = &cobra.Command{
	Use:   "go-basic",
	Short: "Run golang basic code",
	Run: func(cmd *cobra.Command, args []string) {
		//Go basic code to run functions
		k8s := Kubernetes{
			Name:    "k8s-demo-cluster",
			Version: "1.31",
			Users:   []string{"alex", "den"},
			NodeNumber: func() int {
				return 10
			},
		}

		//print users
		k8s.GetUsers()

		//add new user to struct
		k8s.AddNewUser("anonymous")

		//print users one more time
		k8s.GetUsers()
	},
}

func init() {
	rootCmd.AddCommand(goBasicCmd)

}

// My go basic fucntions here
type Kubernetes struct {
	Name       string     `json:"name"`
	Version    string     `json:"version"`
	Users      []string   `json:"users,omitempty"`
	NodeNumber func() int `json:"-"`
}

func (k8s Kubernetes) GetUsers() {
	for _, user := range k8s.Users {
		fmt.Println(user)
	}
}

func (k8s *Kubernetes) AddNewUser(user string) {
	k8s.Users = append(k8s.Users, user)
}
