package metrics

import (
	"time"

	"github.com/signalfx/golib/v3/datapoint"
	k8sutil "github.com/signalfx/signalfx-agent/internal/monitors/kubernetes/utils"
	atypes "github.com/signalfx/signalfx-agent/internal/monitors/types"
	"github.com/signalfx/signalfx-agent/internal/utils"
	batchv1 "k8s.io/api/batch/v1"
)

func datapointsForJob(job *batchv1.Job) []*datapoint.Datapoint {
	dimensions := map[string]string{
		"metric_source":        "kubernetes",
		"kubernetes_namespace": job.Namespace,
		"kubernetes_uid":       string(job.UID),
		"kubernetes_name":      job.Name,
	}

	return []*datapoint.Datapoint{
		datapoint.New(
			"kubernetes.job.completions",
			dimensions,
			datapoint.NewIntValue(int64(*job.Spec.Completions)),
			datapoint.Gauge,
			time.Now()),
		datapoint.New(
			"kubernetes.job.parallelism",
			dimensions,
			datapoint.NewIntValue(int64(*job.Spec.Parallelism)),
			datapoint.Gauge,
			time.Now()),
		datapoint.New(
			"kubernetes.job.active",
			dimensions,
			datapoint.NewIntValue(int64(job.Status.Active)),
			datapoint.Gauge,
			time.Now()),
		datapoint.New(
			"kubernetes.job.failed",
			dimensions,
			datapoint.NewIntValue(int64(job.Status.Failed)),
			datapoint.Counter,
			time.Now()),
		datapoint.New(
			"kubernetes.job.succeeded",
			dimensions,
			datapoint.NewIntValue(int64(job.Status.Succeeded)),
			datapoint.Counter,
			time.Now()),
	}
}

func dimensionForJob(job *batchv1.Job) *atypes.Dimension {
	props, tags := k8sutil.PropsAndTagsFromLabels(job.Labels)

	props["kubernetes_workload"] = "Job"
	props["kubernetes_workload_name"] = job.Name
	props["job_creation_timestamp"] = job.GetCreationTimestamp().Format(time.RFC3339)

	for _, or := range job.OwnerReferences {
		props[utils.LowercaseFirstChar(or.Kind)] = or.Name
		props[utils.LowercaseFirstChar(or.Kind)+"_uid"] = string(or.UID)
	}

	return &atypes.Dimension{
		Name:       "kubernetes_uid",
		Value:      string(job.UID),
		Properties: props,
		Tags:       tags,
	}
}
