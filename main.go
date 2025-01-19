package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {

	deleteEnabled := flag.Bool("d", false, "Delete detected branches")
	flag.Parse() 

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

	remote, err := repo.Remote("origin")

	if err != nil { panic(err) }

	references, _ := remote.List(&git.ListOptions{})
	// Search through the list of references in that remote for a symbolic reference named HEAD;
	// Its target should be the default branch name.
	defaultBranch := ""
	for _, reference := range references {
		if reference.Name() == "HEAD" && reference.Type() == plumbing.SymbolicReference {
			defaultBranch = reference.Target().Short()
			break
		}
	}

	if defaultBranch == "" {
		panic("Could not find default branch")
	}

	branches.ForEach(func(branch *plumbing.Reference) error {

		branchName := branch.Name().Short()

		if branchName == defaultBranch {
			// Skip default branch
			return nil
		}
		
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

		if hasCommit {
			// This means remote branch includes our latest local commit

			if *deleteEnabled {
				// Delete branch
				err = repo.Storer.RemoveReference(branch.Name())
				if err != nil { panic(err) }
			} else {
				fmt.Println(branch)
			}
		}


		return nil
	})


}
