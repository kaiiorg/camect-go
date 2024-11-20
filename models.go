package camect_go

import "encoding/json"

type Info struct {
	// Name is the hub's currently configured name
	Name string `json:"name"`
	// Id is the ID this hub uses to identify itself to Camect's cloud
	Id string `json:"id"`
	// Mode is the current operating mode
	Mode string `json:"mode"`

	// CloudUrl is the url to connect to this hub via Camect cloud
	CloudUrl string `json:"cloud_url"`
	//LocalHttpsUrl is the url needed to connect directly to this hub's local web interface
	LocalHttpsUrl string `json:"local_https_url"`

	ObjectName      []string `json:"object_name"`
	DependentObject []string `json:"dependent_object"`

	TimestampMs uint64 `json:"timestamp_ms"`

	// DiskPartitions is a slice of all of the hub's configured disks, not counting its OS disk
	DiskPartitions []DiskPartition `json:"disk_partitions"`
}

func infoFromJson(j []byte) (*Info, error) {
	i := &Info{}
	err := json.Unmarshal(j, i)
	if err != nil {
		return nil, err
	}
	return i, nil
}

type DiskPartition struct {
	// DevPath is which physical device and partition the hub's operating system reports itself as
	// Example: `/dev/sdb1`
	DevPath string `json:"dev_path"`
	// MountPoint is where in the filesystem this device is mounted
	// Example: `/var/camect/home/mnt/dev_sdb1`
	MountPoint string `json:"mount_point"`
	// Model is what the model this disk reports itself as
	// Example: `ST1000LM048-2E7172`
	Model string `json:"model"`
	// NumCapacityBytes is the total number of bytes this partition can hold
	// Example: `984373075968`
	NumCapacityBytes uint64 `json:"num_capacity_bytes"`
	// NumFreeBytes is the currently free space currently available
	// Example: `934212239360`
	NumFreeBytes uint64 `json:"num_free_bytes"`
}
