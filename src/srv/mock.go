package srv

import (
	m "../models"
)

var clone1Status = &m.Status{
	Code:    "OK",
	Message: "Clone is ready",
}

var clone1Db = &m.Database{
	ConnStr:  "connstr",
	Host:     "host",
	Port:     "port",
	Username: "username",
}

var clone1 = &m.Clone{
	Id:          "xxx",
	Name:        "demo-clone-1",
	Project:     "demo",
	Snapshot:    "timestamp-latest",
	CloneSize:   1000,
	CloningTime: 10,
	Protected:   true,
	CreatedAt:   "timestamp",
	Status:      clone1Status,
	Db:          clone1Db,
}

var clone2Status = &m.Status{
	Code:    "OK",
	Message: "Clone is ready",
}

var clone2Db = &m.Database{
	ConnStr:  "connstr",
	Host:     "host",
	Port:     "port",
	Username: "username",
}

var clone2 = &m.Clone{
	Id:          "yyy",
	Name:        "demo-clone-2",
	Project:     "demo",
	Snapshot:    "timestamp-latest",
	CloneSize:   1000,
	CloningTime: 10,
	Protected:   true,
	CreatedAt:   "timestamp",
	Status:      clone2Status,
	Db:          clone2Db,
}

var clones = []*m.Clone{
	clone1,
	clone2,
}

var snapshot1 = m.Snapshot{
	Id:        "xxx",
	Timestamp: "123",
}

var snapshots = []*m.Snapshot{
	&snapshot1,
}
