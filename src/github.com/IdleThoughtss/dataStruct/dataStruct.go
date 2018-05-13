package dataStruct

type SyncKey struct {
	Count int `json:"Count"`
	List []SyncKeyItem `json:"List"`
}
