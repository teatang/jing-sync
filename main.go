
package main

import (
	"jing-sync/boot"
)

func main() {
	boot.InitDB()
	r := boot.SetupRouter()
	r.Run(":8080")
}
