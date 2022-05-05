package osd

type Args struct {
	DiskID uint32 `pos:"query:disk_id"`
	FD     string `pos:"query:fd"`
}
