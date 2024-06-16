package main

import (
	"fmt"
	"snowflake"
)

func main() {
	id := snowflake.NewSnowflakeId(31, 31, 100)
	id.Generate()
	fmt.Println(id.Time())
	fmt.Println(id.DataCenter())
	fmt.Println(id.Machine())
	fmt.Println(id.Sequence())
	fmt.Println(id.Int64())
}
