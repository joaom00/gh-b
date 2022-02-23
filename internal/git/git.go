package git

import (
	"os/exec"
	"strings"
)

type Branch struct {
	Name          string
	AuthorName    string
	CommitterDate string
	Track         string
	RemoteName    string
}

func (i Branch) FilterValue() string { return "" }

const format = `branch:%(refname:short)%(HEAD)
authorname:%(authorname)
committerdate:%(committerdate:relative)
track:%(upstream:track)
remotename:%(upstream:short)
`

func GetAllBranches() (branches []Branch, err error) {
	cmd := exec.Command(
		"git",
		"for-each-ref",
		"refs/heads",
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
		return
	}

	s := strings.Split(strings.TrimSpace(string(out)), "\n\n")

	for _, branch := range s {
		fields := strings.Split(branch, "\n")

		branch := strings.TrimPrefix(fields[0], "branch:")
		authorname := strings.TrimPrefix(fields[1], "authorname:")
		committerdate := strings.TrimPrefix(fields[2], "committerdate:")
		track := strings.TrimPrefix(fields[3], "track:")
		remotename := strings.TrimPrefix(fields[4], "remotename:")
		branches = append(branches, Branch{
			Name:          strings.TrimSpace(branch),
			AuthorName:    strings.TrimSpace(authorname),
			CommitterDate: strings.TrimSpace(committerdate),
			Track:         strings.TrimSpace(track),
			RemoteName:    strings.TrimSpace(remotename),
		})
	}

	return
}

func CheckoutBranch(branch string) string {
	cmd := exec.Command("git", "checkout", branch)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out)
	}

	return string(out)
}

func CreateBranch(branch string) string {
	cmd := exec.Command("git", "checkout", "-b", branch)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out)
	}

	return string(out)
}

func DeleteBranch(branch string) string {
	cmd := exec.Command("git", "branch", "-D", branch)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out)
	}

	return string(out)
}

func TrackBranch(branch string) string {
	cmd := exec.Command("git", "checkout", "--track", branch)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out)
	}

	return string(out)
}

func MergeBranch(branch string) string {
	cmd := exec.Command("git", "merge", branch)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out)
	}

	return string(out)
}
