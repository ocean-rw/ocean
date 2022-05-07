package proto

type Stripe struct {
	ID    uint64   `bson:"_id"`
	Mode  Mode     `bson:"mode"`
	Disks []uint32 `bson:"disks"`
}
