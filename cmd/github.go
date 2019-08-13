package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
	bolt "go.etcd.io/bbolt"
)

// githubCmd represents the github command
var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "Github archivist support",
	Long:  `Archivist tools to backup your Github account to local storage`,
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
	githubCmd.AddCommand(backupCmd)

	githubCmd.PersistentFlags().StringVar(&githubUser, "github-user", "", "The github user to backup")
	viper.BindPFlag("github.user", githubCmd.PersistentFlags().Lookup("github-user"))
	githubCmd.PersistentFlags().StringVar(&githubToken, "github-token", "", "The github token to use")
	viper.BindPFlag("github.token", githubCmd.PersistentFlags().Lookup("github-token"))

}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Run a backup on your Github account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := ghBackup()
		if err != nil {
			fmt.Println(err)
		}
	},
}

type repo struct {
	FullName    string `json:"full_name"`
	UpdatedDate time.Time `json:"updated_at"`
	CreatedDate time.Time `json:"created_at"`
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
		if Verbose {
			fmt.Println("\tUpdated:", r.UpdatedDate)
			fmt.Println("\tCreated:", r.CreatedDate)
		}
	}

	return nil
}

func ghBackup() error {
	if Verbose {
		fmt.Println("Backup up", C.Green("Github"), "repos for", C.Magenta(viper.Get("github.user")))
	}

	repos, err := ghGetRepos() 
	if err != nil {
		return err
	}

	var ghBackupDate time.Time

	db.View(func (tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Github"))
		ghBackupDate := b.Get([]byte("last-backup"))
		return nil
	})

	tasks := make(chan repo)

	// ToDo fix below statements completely non-functioning

	for _, r := range repos {
		if r.UpdatedDate > ghBackupDate {
			tasks <- r
		}
	}

	var wg sync.WaitGroup 
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func () {
			for repo := range tasks {
				ghBackupRepo(repo)
			}
		}
	}
	wg.Wait()

	// End ToDo
}

func ghBackupRepo(repo) error {
	if Verbose {
		fmt.Println("Backing up ", repo.FullName)
	}

	cloneURL := "https://" + viper.GetString("github.token") + "@github.com/" + repo.FullName

	if Verbose {
		_, err := git.PlainClone(DataDir+"github", false, &git.CloneOptions{
			URL:      cloneURL,
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}
	} else {
		_, err := git.PlainClone(DataDir+"github", false, &git.CloneOptions{
			URL: cloneURL,
		})
		if err != nil {
			return err
		}
	}

	if Verbose {
		fmt.Println("Done backing up", repo.FullName)
	}
}
