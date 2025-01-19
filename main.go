package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {


	cwd, err := os.Getwd()

	if err != nil { panic(err) }

	repo, err := git.PlainOpen(cwd)

	if err != nil { panic(err) }

	err = repo.Fetch(&git.FetchOptions{
		Prune: true,
	})

	if err != nil { panic(err) }

	branches, err := repo.Branches()

	if err != nil { panic(err) }

	// Check if there is uncommited changes in current branch if so stash them and defer pop
	// worktree, err := repo.Worktree()
	// if err != nil { panic(err) }

	// status, err := worktree.Status()
	// if err != nil { panic(err) }


	// if !status.IsClean() {
	// 	panic("Pending changes, stashing not yet supported")
	// }

	// Check for multiple remotes

	remotes, err := repo.Remotes()

	if len(remotes) != 1 {
		panic("Multiple remotes not yet supported")
	}

	remote := remotes[0]

	if err != nil { panic(err) }

	branches.ForEach(func(branch *plumbing.Reference) error {
		fmt.Println(branch.Name().String())

		branchName := branch.Name().Short()
		
		// Check if branch exists in remote
		remoteRef, err := repo.Reference(plumbing.NewRemoteReferenceName(remote.Config().Name, branchName), true)

		if err != nil { 
		
			// This means the remote branch does not exist
			return nil
		}

		iter, err := repo.Log(&git.LogOptions{From: remoteRef.Hash()})

		if err != nil { panic(err) }

		hasCommit := false

		err = iter.ForEach(func(c *object.Commit) error {
			if c.Hash.String() == branch.Hash().String() {
				hasCommit = true
				return nil
				}
			return nil
		})

		if err != nil { panic(err) }

		if !hasCommit {
			fmt.Println("Branch is behind remote")
		}
		

		return nil
	})


}
