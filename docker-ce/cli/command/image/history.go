/*
Copyright (c) 2014-2020 CGCL Labs
Container_Migrate is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/
package image

import (
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/stringutils"
	"github.com/docker/go-units"
	"github.com/spf13/cobra"
)

type historyOptions struct {
	image string

	human   bool
	quiet   bool
	noTrunc bool
}

// NewHistoryCommand creates a new `docker history` command
func NewHistoryCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts historyOptions

	cmd := &cobra.Command{
		Use:   "history [OPTIONS] IMAGE",
		Short: "Show the history of an image",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.image = args[0]
			return runHistory(dockerCli, opts)
		},
	}

	flags := cmd.Flags()

	flags.BoolVarP(&opts.human, "human", "H", true, "Print sizes and dates in human readable format")
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only show numeric IDs")
	flags.BoolVar(&opts.noTrunc, "no-trunc", false, "Don't truncate output")

	return cmd
}

func runHistory(dockerCli *command.DockerCli, opts historyOptions) error {
	ctx := context.Background()

	history, err := dockerCli.Client().ImageHistory(ctx, opts.image)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(dockerCli.Out(), 20, 1, 3, ' ', 0)

	if opts.quiet {
		for _, entry := range history {
			if opts.noTrunc {
				fmt.Fprintf(w, "%s\n", entry.ID)
			} else {
				fmt.Fprintf(w, "%s\n", stringid.TruncateID(entry.ID))
			}
		}
		w.Flush()
		return nil
	}

	var imageID string
	var createdBy string
	var created string
	var size string

	fmt.Fprintln(w, "IMAGE\tCREATED\tCREATED BY\tSIZE\tCOMMENT")
	for _, entry := range history {
		imageID = entry.ID
		createdBy = strings.Replace(entry.CreatedBy, "\t", " ", -1)
		if !opts.noTrunc {
			createdBy = stringutils.Ellipsis(createdBy, 45)
			imageID = stringid.TruncateID(entry.ID)
		}

		if opts.human {
			created = units.HumanDuration(time.Now().UTC().Sub(time.Unix(entry.Created, 0))) + " ago"
			size = units.HumanSizeWithPrecision(float64(entry.Size), 3)
		} else {
			created = time.Unix(entry.Created, 0).Format(time.RFC3339)
			size = strconv.FormatInt(entry.Size, 10)
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", imageID, created, createdBy, size, entry.Comment)
	}
	w.Flush()
	return nil
}
