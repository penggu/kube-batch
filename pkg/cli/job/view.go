package job


import (
	"fmt"
	"io"
	"os"
	"volcano.sh/volcano/pkg/apis/batch/v1alpha1"

	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"volcano.sh/volcano/pkg/client/clientset/versioned"
)

type viewFlags struct {
	commonFlags

	Namespace string
	JobName   string
}


var viewJobFlags = &viewFlags{}

func InitViewFlags(cmd *cobra.Command) {
	initFlags(cmd, &viewJobFlags.commonFlags)

	cmd.Flags().StringVarP(&viewJobFlags.Namespace, "namespace", "", "default", "the namespace of job")
	cmd.Flags().StringVarP(&viewJobFlags.JobName, "name", "n", "", "the name of job")
}

func ViewJob() error {
	config, err := buildConfig(viewJobFlags.Master, viewJobFlags.Kubeconfig)
	if err != nil {
		return err
	}

	jobClient := versioned.NewForConfigOrDie(config)
	job, err := jobClient.BatchV1alpha1().Jobs(viewJobFlags.Namespace).Get(viewJobFlags.JobName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if job == nil {
		fmt.Printf("No resources found\n")
		return nil
	}
	PrintJob(job, os.Stdout)

	return nil
}

func PrintJob(job *v1alpha1.Job, writer io.Writer) {
	replicas := int32(0)
	for _, ts := range job.Spec.Tasks {
		replicas += ts.Replicas
	}
	_, err := fmt.Fprintf(writer, "%-12s %s\n%-12s %s\n%-12s %s\n%-12s %d\n%-12s %d\n%-12s %d\n%-12s %d\n%-12s %d\n%-12s %d\n%-12s %d\n",
		Name, job.Name, Creation, job.CreationTimestamp.Format("2006-01-02 15:04:05"), Phase, job.Status.State.Phase, Replicas, replicas,
		Min, job.Status.MinAvailable, Pending, job.Status.Pending, Running, job.Status.Running, Succeeded, job.Status.Succeeded, Failed, job.Status.Failed, RetryCount, job.Status.RetryCount)
	if err != nil {
		fmt.Printf("Failed to print view command result: %s.\n", err)
	}
}