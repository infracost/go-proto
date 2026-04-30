package storage

type Storage struct {
	Buckets []Bucket `tree:"buckets"`
}
