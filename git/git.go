package git

import (
	"os/exec"
	"strings"
)

type Branch struct {
	Name          string
	AuthorName    string
	CommitterDate string
	IsRemote      bool
}

const format = `branch:%(refname:short)%(HEAD)
authorname:%(authorname)
committerdate:%(committerdate:relative)
`

func getLocalBranches(branches []Branch) ([]Branch, error) {
	cmd := exec.Command(
		"git",
		"for-each-ref",
		"refs/heads",
		"--sort",
		"-committerdate",
		"--sort",
		"-upstream",
		"--format",
		format,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return []Branch{}, err
	}

	s := strings.Split(strings.TrimSpace(string(out)), "\n\n")

	for _, branch := range s {
		fields := strings.Split(branch, "\n")

		branch := strings.TrimPrefix(fields[0], "branch:")
		authorname := strings.TrimPrefix(fields[1], "authorname:")
		committerdate := strings.TrimPrefix(fields[2], "committerdate:")
		branches = append(branches, Branch{
			Name:          strings.TrimSpace(branch),
			AuthorName:    strings.TrimSpace(authorname),
			CommitterDate: strings.TrimSpace(committerdate),
			IsRemote:      false,
		})
	}

	return branches, err
}

func getRemoteBranches(branches []Branch) ([]Branch, error) {
	cmd := exec.Command(
		"git",
		"for-each-ref",
		"refs/remotes",
		"--sort",
		"-committerdate",
		"--sort",
		"-upstream",
		"--format",
		format,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return []Branch{}, err
	}

	s := strings.Split(strings.TrimSpace(string(out)), "\n\n")

	for _, branch := range s {
		fields := strings.Split(branch, "\n")

		branch := strings.TrimPrefix(fields[0], "branch:")
		authorname := strings.TrimPrefix(fields[1], "authorname:")
		committerdate := strings.TrimPrefix(fields[2], "committerdate:")
		branches = append(branches, Branch{
			Name:          strings.TrimSpace(branch),
			AuthorName:    strings.TrimSpace(authorname),
			CommitterDate: strings.TrimSpace(committerdate),
			IsRemote:      true,
		})
	}

	return branches, err
}

func GetAllBranches() (branches []Branch, err error) {
	branches, err = getLocalBranches(branches)
	if err != nil {
		return
	}

	branches, err = getRemoteBranches(branches)
	if err != nil {
		return
	}

	return
}

func CheckoutBranch(branch string) string {
	cmd := exec.Command("git", "checkout", branch)

	out, _ := cmd.CombinedOutput()

	return string(out)
}

func CreateBranch(branch string) string {
	cmd := exec.Command("git", "checkout", "-b", branch)

	out, _ := cmd.CombinedOutput()

	return string(out)
}

func DeleteBranch(branch string) string {
	cmd := exec.Command("git", "branch", "-D", branch)

	out, _ := cmd.CombinedOutput()

	return string(out)
}

func TrackBranch(branch string) string {
	cmd := exec.Command("git", "checkout", "--track", branch)

	out, _ := cmd.CombinedOutput()

	return string(out)
}

func MergeBranch(branch string) string {
	cmd := exec.Command("git", "merge", branch)

	out, _ := cmd.CombinedOutput()

	return string(out)
}

func RebaseBranch(branch string) string {
	cmd := exec.Command("git", "rebase", branch)

	out, _ := cmd.CombinedOutput()

	return string(out)
}

func RenameBranch(oldName, newName string) string {
	cmd := exec.Command("git", "branch", "-m", oldName, newName)

	out, _ := cmd.CombinedOutput()

	return string(out)
}

func RenameRemoteBranch(oldName, newName string) string {
	exec.Command("git", "push", "origin", "--delete", oldName)

	cmd := exec.Command("git", "push", "origin", "-u", newName)
	out, _ := cmd.CombinedOutput()

	return string(out)
}
