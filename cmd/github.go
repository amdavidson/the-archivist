package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// githubCmd represents the github command
var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "Backup your Github account",
	Long:  `Backup your Github account to local storage`,
	Run: func(cmd *cobra.Command, args []string) {
		err := ghList()
		if err != nil {
			fmt.Println(err)
		}
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

type repo struct {
	FullName    string `json:"full_name"`
	CloneHTTPS  string `json:"clone_url"`
	UpdatedDate string `json:"updated_at"`
	CreatedDate string `json:"created_at"`
	ArchiveURL  string `json:"archive_url"`
}

func ghQuery(url string) ([]byte, string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	req.Header.Add("Authorization", "Bearer "+viper.GetString("github.token"))
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	next := ""
	//links := strings.Split(resp.Header["Link"], ",")
	for _, l := range strings.Split(resp.Header["Link"][0], ",") {
		if strings.Contains(l, `rel="next"`) {
			str := strings.Split(l, ";")[0]
			next = strings.Trim(str, " <>")
		}
	}

	return body, next, nil
}

func ghGetRepos() ([]repo, error) {
	next := "https://api.github.com/user/repos"
	var repos []repo

	for {
		var body []byte
		var err error
		body, next, err = ghQuery(next)
		if err != nil {
			return nil, err
		}
		var temp []repo
		err = json.Unmarshal(body, &temp)
		if err != nil {
			return nil, err
		}
		repos = append(repos, temp...)
		if next == "" {
			break
		}
	}

	return repos, nil

}

func ghList() error {
	if Verbose {
		fmt.Println("Listing", C.Green("Github"), "Repos for", C.Magenta(viper.Get("github.user")))
	}

	repos, err := ghGetRepos()
	if err != nil {
		return err
	}

	fmt.Println(viper.GetString("github.user"), "has", len(repos), "repositories.")

	for _, r := range repos {

		fmt.Println(r.FullName)
		fmt.Println("\tUpdated:", r.UpdatedDate)
		fmt.Println("\tCreated:", r.CreatedDate)

	}

	return nil
}
