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

type Camera struct {
	// Id is the hub's internal ID for this logical camera
	Id string `json:"id"`
	// Name is the user friendly name for this logical camera
	Name string `json:"name"`
	// Make is the type of camera;
	// Example: ONVIF
	Make string `json:"make"`
	// Model is the model of this camera, if reported it
	Model string `json:"model"`
	// Url is the url to the camera, if reported
	Url string `json:"url"`
	// Width is the current width of the stream in pixels
	Width uint `json:"width"`
	// Height is the current height of the stream in pixels
	Height uint `json:"height"`
	// Wifi is if the hub has determined this camera is on wifi
	Wifi bool `json:"is_on_wifi"`
	// Disabled if the hub is ignoring this camera
	Disabled bool `json:"disabled"`
	// MacAddr reported mac address of this camera
	// Example:x-singlemac-e4:30:22:99:e9:64
	MacAddr string `json:"mac_addr"`
	// IPAddr is the IP address of this camera
	// Example: 192.168.1.221
	IPAddr string `json:"ip_addr"`
	// Streaming is if the hub is currently streaming from this camera
	Streaming bool `json:"is_streaming"`
	// StreamingUrl is the RTSP stream the hub is currently using to stream from this camera
	// Example: rtsp://REDACTED_USERNAME:REDACTED_PASSWORD@192.168.1.221:554/0/onvif/profile2/media.smp
	StreamingUrl string `json:"streaming_url"`
	// AlertsDisabled is if the hub is looking for objects on this camera and will send alerts
	AlertsDisabled bool `json:"is_alerts_disabled"`
	SaveEncoded    uint `json:"save_encoded"`
}
