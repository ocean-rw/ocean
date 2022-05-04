package proto

type ListDisksArgs struct {
	Host  string    `in:"query=host"`
	State DiskState `in:"query=state"`
}
