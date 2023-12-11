package main

import (
	"errors"
	"fmt"
	"github.com/kballard/go-shellquote"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

const BRANCH_PREFIX = "tree"

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use: "git-open",
	Run: openCurrentRepo,
}

func openCurrentRepo(cmd *cobra.Command, args []string) {
	var remote string
	if len(args) > 0 {
		remote = args[0]
	}
	_, err := CurrentGitRepo()
	if err != nil {
		Err("not a git repository")
	}
	branch := CurrentBranch()
	if remote == "" {
		remote = CurrentRemote(branch)
	}
	gitURL := RemoteURL(remote)
	url := TransferToURL(gitURL)
	url = TransferWithBranch(url, branch)
	_ = OpenBrowser(url)
}

// TransferWithBranch 添加分支名
func TransferWithBranch(url string, branch string) string {
	url = fmt.Sprintf("%s/%s/%s", url, BRANCH_PREFIX, branch)
	return url
}

func GitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	return string(output), err
}

// CurrentGitRepo 获取当前git项目路径
func CurrentGitRepo() (string, error) {
	output, err := GitCommand("rev-parse", "-q", "--show-toplevel")
	return output, err
}

// CurrentBranch 获取当前分支，如果不存在则使用master
func CurrentBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("generate branch failed")
	}
	outputStr := strings.TrimSuffix(string(output), "\n")
	branchName := strings.TrimSpace(outputStr)
	return branchName
}

// CurrentRemote 通过分支获取当前remote，默认为origin
func CurrentRemote(branch string) string {
	remote, err := GitCommand("config", fmt.Sprintf("branch.%s.remote", branch))
	fmt.Printf("the URL is %s \n", remote)

	remote = firstLine(remote)
	if remote == "" || err != nil {
		remote = "origin"
	}
	return remote
}

// RemoteURL 通过remote获取remote-url，如果使用错误的remote，则报错退出
func RemoteURL(remote string) string {
	gitURL, err := GitCommand("ls-remote", "--get-url", remote)
	if err != nil {
		Err("git remote is not set for", remote)
	}
	gitURL = firstLine(gitURL)
	if gitURL == remote {
		Err(remote, "is a wrong remote")
	}
	return gitURL
}

// 移除末尾换行符
func firstLine(output string) string {
	if i := strings.Index(output, "\n"); i >= 0 {
		return output[0:i]
	}
	return output
}

// TransferToURL remote url转换为web url
func TransferToURL(gitURL string) string {
	var url string
	if strings.HasPrefix(gitURL, "https://") || strings.HasPrefix(gitURL, "http://") {
		url = gitURL[:len(gitURL)-4]
	}
	//处理一下  https://oauth2:4be21afa355573d342c17a4517537bde7c4c8c8d78dc018df89bd6a73700f5ee@git.dfy.definesys.cn/apaas/xdap-app.git 格式的URL
	if strings.HasPrefix(gitURL, "git@") {
		url = gitURL[:len(gitURL)-4]
		url = strings.Replace(url, ":", "/", 1)
		url = strings.Replace(url, "git@", "https://", 1)
	}

	re := regexp.MustCompile(`https://.*@`)
	match := re.FindString(url)
	if match != "" {
		// 将匹配到的部分替换为https://
		url = strings.Replace(url, match, "https://", 1)
	}
	return url
}

// Err 输出报错，结束程序
func Err(msg ...interface{}) {
	fmt.Println(msg...)
	os.Exit(1)
}

// OpenBrowser 唤起浏览器，打开url
func OpenBrowser(url string) error {
	launcher, err := browserLauncher()
	if err != nil {
		return err
	}
	args := append(launcher, url)
	//fmt.Printf("使用浏览器：%s， 打开链接：%s", launcher, url)
	cmd := exec.Command(args[0], args[1:]...)
	return cmd.Run()
}

// 切分命令为数组，方便exec.Command执行
func browserLauncher() ([]string, error) {
	browser := os.Getenv("BROWSER")
	if browser == "" {
		browser = searchBrowserLauncher()
	} else {
		browser = os.ExpandEnv(browser)
	}

	if browser == "" {
		return nil, errors.New("please set $BROWSER to a web launcher")
	}
	return shellquote.Split(browser)
}

// 根据操作系统，返回默认browser command
func searchBrowserLauncher() (browser string) {
	switch runtime.GOOS {
	case "darwin":
		browser = "open"
	case "windows":
		browser = "cmd /c start"
	case "linux":
		browser = "xdg-open"
	default:
		browser = ""
	}
	return browser
}
