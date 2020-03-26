package cmd

import (
	goflag "flag"
	"fmt"
	"github.com/alcideio/rbac-tool/pkg/kube"
	"github.com/alcideio/rbac-tool/pkg/utils"
	"github.com/alcideio/rbac-tool/pkg/visualize"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func NewCommandVisualize() *cobra.Command {

	opts := visualize.Opts{}

	// Support overrides
	cmd := &cobra.Command{
		Use:     "visualize",
		Aliases: []string{"vis", "viz"},
		Short:   "A RBAC visualizer",
		Long: `
A Kubernetes RBAC visualizer - Generate a graph as dot file format.

By default 'rbac-tool viz' will connect to the local cluster (pointed by kubeconfig)
Create a RBAC graph of the actively running workload on all namespaces except kube-system

See run options on how to render specific namespaces, other clusters, etc.

#Render Locally
rbac-tool viz --outformat dot && cat rbac.dot | dot -Tpng > rbac.png  && open rbac.png

# Render Online
https://dreampuf.github.io/GraphvizOnline

Examples:

# Scan the cluster pointed by the kubeconfig context 'myctx'
rbac-tool viz --cluster-context myctx

# Scan and create a PNG image from the graph
rbac-tool viz  --outformat dot --exclude-namespaces=soemns && cat rbac.dot | dot -Tpng > rbac.png && google-chrome rbac.png

`,
		Hidden: false,
		RunE: func(c *cobra.Command, args []string) error {

			utils.ConsolePrinter(fmt.Sprintf("Connecting to cluster '%v'", color.HiBlueString(opts.ClusterContext)))

			kubeClient, err := kube.NewClient(opts.ClusterContext)
			if err != nil {
				return fmt.Errorf("Failed to create kubernetes client - %v", err)
			}

			utils.ConsolePrinter(fmt.Sprintf("Generating RBAC Graph to cluster '%v'", color.HiBlueString(opts.ClusterContext)))

			utils.ConsolePrinter(fmt.Sprintf("Namespaces included %v", color.GreenString("'%v'", opts.IncludedNamespaces)))

			if len(opts.ExcludedNamespaces) > 0 {
				utils.ConsolePrinter(fmt.Sprintf("Namespaces excluded %v", color.HiRedString("'%v'", opts.ExcludedNamespaces)))
			}

			return visualize.CreateRBACGraph(kubeClient, &opts)
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&opts.ClusterContext, "cluster-context", "", "Cluster Context .use 'kubectl config get-contexts' to list available contexts")
	flags.StringVar(&opts.Outfile, "outfile", "rbac.html", "Output file")
	flags.StringVar(&opts.Outformat, "outformat", "html", "Output format: dot or html")
	flags.StringVar(&opts.IncludedNamespaces, "include-namespaces", "*", "Comma-delimited list of namespaces to include in the visualization")
	flags.StringVar(&opts.ExcludedNamespaces, "exclude-namespaces", "kube-system", "Comma-delimited list of namespaces to include in the visualization")

	flags.BoolVar(&opts.ShowLegend, "show-legend", false, "Whether to show the legend or not (for dot format)")
	flags.BoolVar(&opts.ShowRules, "show-rules", true, "Whether to render RBAC access rules (e.g. \"get pods\") or not")

	klog.InitFlags(nil)
	flags.AddGoFlagSet(goflag.CommandLine)
	return cmd
}
