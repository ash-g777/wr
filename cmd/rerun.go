// Copyright © 2018 Genome Research Limited
// Author: Sendu Bala <sb10@sanger.ac.uk>.
//
//  This file is part of wr.
//
//  wr is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Lesser General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  wr is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU Lesser General Public License for more details.
//
//  You should have received a copy of the GNU Lesser General Public License
//  along with wr. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"time"

	"github.com/VertebrateResequencing/wr/jobqueue"
	"github.com/spf13/cobra"
)

// options for this cmd
var cmdAll bool

// rerunCmd represents the rerun command
var rerunCmd = &cobra.Command{
	Use:   "rerun",
	Short: "Rerun completed commands",
	Long: `You can rerun commands you've previously added with "wr add" that
have since completed and become "complete" using this command.

Specify one of the flags -f, -l, -i or -a to choose which commands you want to
rerun. Amongst those, only currently completed jobs will be affected.

-i is the report group (-i) you supplied to "wr add" when you added the job(s)
you want to now rerun. Combining with -z lets you rerun jobs in multiple report
groups, assuming you have arranged that related groups share some substring. Or
-y lets you specify -i as the internal job id reported during "wr status".

The file to provide -f is in the format taken by "wr add".

In -f and -l mode you must provide the cwd the commands were set to run in, if
CwdMatters (and must NOT be provided otherwise). Likewise provide the mounts
options that was used when the command was originally added, if any. You can do this by
using the -c and --mounts/--mounts_json options in -l mode, or by providing the
same file you gave to "wr add" in -f mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		set := countGetJobArgs()
		if set > 1 {
			die("-f, -i, -l and -a are mutually exclusive; only specify one of them")
		}
		if set == 0 {
			die("1 of -f, -i, -l or -a is required")
		}

		timeout := time.Duration(timeoutint) * time.Second
		jq := connect(timeout)
		var err error
		defer func() {
			err = jq.Disconnect()
			if err != nil {
				warn("Disconnecting from the server failed: %s", err)
			}
		}()

		jobs := getJobs(jq, jobqueue.JobStateComplete, cmdAll, 0, false, false)

		if len(jobs) == 0 {
			die("No matching jobs found")
		}

		jes := jobsToJobEssenses(jobs)
		restarted, err := jq.Add(jes)
		if err != nil {
			die("failed to rerun desired jobs: %s", err)
		}
		info("Initiated rerun of %d complete commands (out of %d eligible)", kicked, len(jobs))
	},
}

func init() {
	RootCmd.AddCommand(rerunCmd)

	// flags specific to this sub-command
	rerunCmd.Flags().BoolVarP(&cmdAll, "all", "a", false, "rerun all completed jobs")
	rerunCmd.Flags().StringVarP(&cmdFileStatus, "file", "f", "", "file containing commands you want to rerun; - means read from STDIN")
	rerunCmd.Flags().StringVarP(&cmdIDStatus, "identifier", "i", "", "identifier of the commands you want to rerun")
	rerunCmd.Flags().BoolVarP(&cmdIDIsSubStr, "search", "z", false, "treat -i as a substring to match against all report groups")
	rerunCmd.Flags().BoolVarP(&cmdIDIsInternal, "internal", "y", false, "treat -i as an internal job id")
	rerunCmd.Flags().StringVarP(&cmdLine, "cmdline", "l", "", "a command line you want to rerun")
	rerunCmd.Flags().StringVarP(&cmdCwd, "cwd", "c", "", "working dir that the command(s) specified by -l or -f were set to run in")
	rerunCmd.Flags().StringVarP(&mountJSON, "mount_json", "j", "", "mounts that the command(s) specified by -l or -f were set to use (JSON format)")
	rerunCmd.Flags().StringVar(&mountSimple, "mounts", "", "mounts that the command(s) specified by -l or -f were set to use (simple format)")

	rerunCmd.Flags().IntVar(&timeoutint, "timeout", 120, "how long (seconds) to wait to get a reply from 'wr manager'")
}
