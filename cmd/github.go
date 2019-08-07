package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// githubCmd represents the github command
var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "Backup your Github account",
	Long:  `Backup your Github account to local storage`,
	Run: func(cmd *cobra.Command, args []string) {
		githubList()
	},
}

var githubUser string
var githubToken string

func init() {
	rootCmd.AddCommand(githubCmd)

	githubCmd.PersistentFlags().StringVar(&githubUser, "github-user", "", "The github user to backup")
	viper.BindPFlag("github.user", githubCmd.PersistentFlags().Lookup("github-user"))
	githubCmd.PersistentFlags().StringVar(&githubToken, "github-token", "", "The github token to use")
	viper.BindPFlag("github.token", githubCmd.PersistentFlags().Lookup("github-token"))

}

func githubList() {
	if Verbose {
		fmt.Println("Listing", C.Green("Github"), "Repos for", C.Magenta(viper.Get("github.user")))
	}

	githubUrl := "https://api.github.com/users/" + viper.GetString("github.user") + "/repos"

	req, err := http.NewRequest("GET", githubUrl, nil)

	req.Header.Add("Authorization", "Bearer "+viper.GetString("github.token"))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(C.Red("[ERROR]"), err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf(string([]byte(body)))

}
