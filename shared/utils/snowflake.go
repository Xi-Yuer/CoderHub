package utils

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"time"
)

var Node *snowflake.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}

	snowflake.Epoch = st.UnixNano() / 1e6
	Node, err = snowflake.NewNode(machineID)
	if err != nil {
		return
	}

	return
}

func init() {
	if err := Init("2024-01-01", 1); err != nil {
		fmt.Println("Init() failed, err = ", err)
		return
	}
}

// GenID 生成 64 位的 雪花 ID
func GenID() int64 {
	return Node.Generate().Int64()
}
