package models

import (
	"time"
)

type PR struct {
	Commits       int           `json:"commits"`
	Additions     int           `json:"additions"`
	Deletions     int           `json:"deletions"`
	Number        int           `json:"number"`
	Title         string        `json:"title"`
	Body          string        `json:"body"`
	Assignees     []interface{} `json:"assignees"`
	AssigneesSize int
	User          *User `json:"user"`
	FetchedAt     time.Time
	RequestTime   int64
}

func (this *PR) SetEndTime(initTime time.Time) {
	this.FetchedAt = initTime
	endMillis := time.Now().UnixNano() / int64(time.Millisecond)
	this.RequestTime = endMillis - (initTime.UnixNano() / int64(time.Millisecond))
}

//
// func (this *PR) UnmarshalJSON(bytes []byte) error {
// 	var f interface{}
// 	if err := json.Unmarshal(bytes, &f); err != nil {
// 		return err
// 	}
// 	// cast bytes to Map
// 	bytesAsMap := f.(map[string]interface{})
//
// 	assigneesMapInterface := bytesAsMap["assignees"] // [{"login": "luchojuarez"}]
// 	if assigneesMapInterface == nil {
// 		log.Printf("no tiene assignees")
// 		return nil
// 	}
//
// 	listAssigneesInterface := assigneesMapInterface.([]interface{})
// 	log.Printf("esto tiene list %s", listAssigneesInterface)
//
// 	for _, u := range listAssigneesInterface {
// 		userMap := u.(map[string]interface{})
// 		log.Printf("esto tiene dentro del for %s", userMap["login"])
// 	}
//
// 	return nil
// }
