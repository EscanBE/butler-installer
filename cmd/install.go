package cmd

import (
	"bufio"
	"fmt"
	"github.com/EscanBE/butler-installer/constants"
	"github.com/EscanBE/butler-installer/types"
	"github.com/EscanBE/butler-installer/utils"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Escan Butler binary",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := types.UnwrapAppContext(cmd.Context())

		operationUserInfo := ctx.GetOperationUserInfo()
		if !operationUserInfo.OperatingAsSuperUser {
			libutils.PrintlnStdErr("ERR: must be ran as super user")
			libutils.PrintlnStdErr("*** Hint: Try again with 'sudo' or use 'root' ***")
			os.Exit(1)
		}

		workingUser := ctx.GetWorkingUserInfo()

		chownRecursive := func(path string) {
			if workingUser.Username == "root" {
				return
			}
			var groupName string
			if utils.IsDarwin() {
				groupName = "staff"
			} else {
				groupName = workingUser.Username
			}
			owner := fmt.Sprintf("%s:%s", workingUser.Username, groupName)
			fmt.Println("Changing owner of", path, "to", owner)
			err := execCmd("sudo", "", "chown", "-R", owner, path)
			if err != nil {
				libutils.PrintlnStdErr("ERR: failed to change owner of", path, "to", owner, "with error:", err)
				os.Exit(1)
			}
		}

		var cloneButlerRepo bool
		var createNetrcFile bool

		localRepoDirPath := path.Join(workingUser.HomeDir, constants.BUTLER_LOCAL_REPO_DIR)
		_, err := os.Stat(localRepoDirPath)
		if err == nil {
			cloneButlerRepo = false
		} else {
			if !os.IsNotExist(err) {
				libutils.PrintlnStdErr("ERR: failed to get info of local repository directory:", err)
				os.Exit(1)
			}

			cloneButlerRepo = true
		}

		netrcDirPath := path.Join(workingUser.HomeDir, ".netrc")
		const netrcPermission = 0o600
		_, err = os.Stat(netrcDirPath)
		if err != nil {
			if !os.IsNotExist(err) {
				libutils.PrintlnStdErr("ERR: failed to get info of .netrc file:", err)
				os.Exit(1)
			}

			createNetrcFile = true
		} else {
			bz, err := os.ReadFile(netrcDirPath)
			if err != nil {
				libutils.PrintlnStdErr("ERR: failed to read .netrc file:", err)
				os.Exit(1)
			}

			content := string(bz)
			createNetrcFile = libutils.IsBlank(content) || !strings.Contains(content, "machine github.com")
		}

		reader := bufio.NewReader(os.Stdin)

		if createNetrcFile {
			fmt.Println(`About to cloning Butler repository, please prepare credentials for GitHub.
Needed information is a GitHub Fine Grained Access Token with readonly 'repo' scope.`)
			fmt.Println("Please enter your GitHub User:")
			githubUser, err := reader.ReadString('\n')
			if err != nil {
				libutils.PrintlnStdErr("ERR: failed to read GitHub user:", err)
				os.Exit(1)
			}
			githubUser = strings.TrimSpace(githubUser)
			if githubUser == "" {
				libutils.PrintlnStdErr("ERR: GitHub user cannot be empty")
				os.Exit(1)
			}

			fmt.Println("Please enter your GitHub Fine Grained Access Token:")
			bz, err := term.ReadPassword(syscall.Stdin)
			if err != nil {
				libutils.PrintlnStdErr("ERR: failed to read GitHub token:", err)
				os.Exit(1)
			}
			githubToken := strings.TrimSpace(string(bz))
			if githubToken == "" {
				libutils.PrintlnStdErr("ERR: GitHub token cannot be empty")
				os.Exit(1)
			}

			func() {
				f, err := os.OpenFile(netrcDirPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, netrcPermission)
				if err != nil {
					libutils.PrintlnStdErr("ERR: failed to open .netrc file:", err)
					os.Exit(1)
				}
				defer func() {
					_ = f.Close()
				}()

				_, err = f.WriteString(fmt.Sprintf("\nmachine github.com\nlogin %s\npassword %s\n", githubUser, githubToken))
				if err != nil {
					libutils.PrintlnStdErr("ERR: failed to write .netrc file:", err)
					os.Exit(1)
				}

				err = f.Chmod(netrcPermission)
				if err != nil {
					libutils.PrintlnStdErr("ERR: failed to set permission of .netrc file:", err)
					os.Exit(1)
				}

				chownRecursive(netrcDirPath)
			}()
		}

		if cloneButlerRepo {
			fmt.Println("Going to clone Butler repository to:", localRepoDirPath)
			err = execCmd("git", workingUser.HomeDir, "clone", constants.BUTLER_REPO_URL, localRepoDirPath)
			if err != nil {
				libutils.PrintlnStdErr("ERR: failed to clone Butler repository:", err)
				os.Exit(1)
			}

			chownRecursive(localRepoDirPath)
		} else {
			err = execCmd("git", localRepoDirPath, "fetch", "--all")
			if err != nil {
				libutils.PrintlnStdErr("ERR: failed to fetch changes from Butler remote repository:", err)
				os.Exit(1)
			}
		}

		err = execCmd("bash", localRepoDirPath, "-c", "git checkout $(git describe --tags `git rev-list --tags --max-count=1`)")
		if err != nil {
			libutils.PrintlnStdErr("ERR: failed to checkout latest tag of Butler repository:", err)
			os.Exit(1)
		}

		installButler(localRepoDirPath)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func installButler(localRepoDirPath string) {
	err := execCmd("make", localRepoDirPath, "install")
	if err != nil {
		libutils.PrintlnStdErr("ERR: failed to install Butler:", err)
		os.Exit(1)
	}
}

func execCmd(name, dir string, args ...string) error {
	launchCmd := exec.Command(name, args...)
	if len(dir) > 0 {
		launchCmd.Dir = dir
	}
	launchCmd.Stdin = os.Stdin
	launchCmd.Stdout = os.Stdout
	launchCmd.Stderr = os.Stderr
	return launchCmd.Run()
}
