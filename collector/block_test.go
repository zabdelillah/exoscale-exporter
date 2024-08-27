package collector

import (
    "net/http"
    "github.com/exoscale/egoscale/v3"
    "testing"
)

var dummySnapshots = v3.ListSnapshotsResponse {
    Snapshots: []v3.Snapshot{
        {
            ID: "dummySnapshotID",
            Name: "dummySnapshotName",
            Size: 64000,
        },
    },
}

// var dummyTemplates = v3.ListTemplatesResponse {
//     Templates: []v3.Template{
//         {
//             ID: "dummySnapshotID",
//             Name: "dummySnapshotName",
//             Size: 64000,
//         },
//     },
// }

var dummyBlockStorageVolumes = v3.ListBlockStorageVolumesResponse {
    BlockStorageVolumes: []v3.BlockStorageVolume{
        {
            ID: "dummySnapshotID",
            Name: "dummySnapshotName",
            Size: 64000,
        },
    },
}

func SetupBlockStorageTestEndpoints() {
    http.HandleFunc("/snapshot", HandleTestSnapshotsResponse)
    http.HandleFunc("/block-storage", HandleTestBlockStorageResponse)
}

func HandleTestSnapshotsResponse(w http.ResponseWriter, r *http.Request) {
    WriteObjectToResponse(w, r, dummySnapshots)
}

func HandleTestBlockStorageResponse(w http.ResponseWriter, r *http.Request) {
    WriteObjectToResponse(w, r, dummyBlockStorageVolumes)
}

func TestBlockStorageMetricsExist(t *testing.T) {
    metrics := GetTestMetrics(t)

    metricsToCheck := []string {
        "exoscale_snapshot_count",
        "exoscale_snapshot_size",
        "exoscale_volume_count",
        "exoscale_volume_size",
    }

    _, errs := CheckMetricsExist(t, metricsToCheck, metrics)
    for i := range(errs) {
        t.Errorf("Instance Metric Check Failed: %v", errs[i])
    }
}