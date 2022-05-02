package main

import (
	"fmt"
	"github.com/muzin/dpcp"
	"io/ioutil"
	"time"
)

func main() {

	byteurl := "/Users/sirius/bucket/data/__tmp__/bytes/"

	files, err := ioutil.ReadDir(byteurl)
	if err != nil {
		fmt.Println(err)
	}

	var bytes2d [][]byte = make([][]byte, len(files))

	for i := 0; i < len(files); i++ {
		fileInfo := files[i]
		fileurl := byteurl + fileInfo.Name()
		itemBytes, _ := ioutil.ReadFile(fileurl)
		bytes2d[i] = itemBytes
	}

	session := dpcp.NewSession()

	totalLength := 0

	now := time.Now()
	for i := 0; i < len(bytes2d); i++ {
		bytes := bytes2d[i]

		//fmt.Printf("%04d \n", i)
		//if i == 119 {
		//	fmt.Printf("%04d \n", i)
		//}

		processBytes := session.ProcessMessage(bytes, func(msg dpcp.Message) {})
		if processBytes != len(bytes) {
			fmt.Printf("%04d bytes length: %d  process length: %d\n", i, len(bytes), processBytes)
		}

		totalLength += processBytes
	}
	for i := 0; i < len(bytes2d); i++ {
		bytes := bytes2d[i]

		//fmt.Printf("%04d \n", i)
		//if i == 119 {
		//	fmt.Printf("%04d \n", i)
		//}

		processBytes := session.ProcessMessage(bytes, func(msg dpcp.Message) {})
		if processBytes != len(bytes) {
			fmt.Printf("%04d bytes length: %d  process length: %d\n", i, len(bytes), processBytes)
		}

		totalLength += processBytes
	}

	for i := 0; i < len(bytes2d); i++ {
		bytes := bytes2d[i]

		//fmt.Printf("%04d \n", i)
		//if i == 119 {
		//	fmt.Printf("%04d \n", i)
		//}

		processBytes := session.ProcessMessage(bytes, func(msg dpcp.Message) {})
		if processBytes != len(bytes) {
			fmt.Printf("%04d bytes length: %d  process length: %d\n", i, len(bytes), processBytes)
		}

		totalLength += processBytes
	}
	for i := 0; i < len(bytes2d); i++ {
		bytes := bytes2d[i]

		//fmt.Printf("%04d \n", i)
		//if i == 119 {
		//	fmt.Printf("%04d \n", i)
		//}

		processBytes := session.ProcessMessage(bytes, func(msg dpcp.Message) {})
		if processBytes != len(bytes) {
			fmt.Printf("%04d bytes length: %d  process length: %d\n", i, len(bytes), processBytes)
		}

		totalLength += processBytes
	}
	for i := 0; i < len(bytes2d); i++ {
		bytes := bytes2d[i]

		//fmt.Printf("%04d \n", i)
		//if i == 119 {
		//	fmt.Printf("%04d \n", i)
		//}

		processBytes := session.ProcessMessage(bytes, func(msg dpcp.Message) {})
		if processBytes != len(bytes) {
			fmt.Printf("%04d bytes length: %d  process length: %d\n", i, len(bytes), processBytes)
		}

		totalLength += processBytes
	}
	for i := 0; i < len(bytes2d); i++ {
		bytes := bytes2d[i]

		//fmt.Printf("%04d \n", i)
		//if i == 119 {
		//	fmt.Printf("%04d \n", i)
		//}

		processBytes := session.ProcessMessage(bytes, func(msg dpcp.Message) {})
		if processBytes != len(bytes) {
			fmt.Printf("%04d bytes length: %d  process length: %d\n", i, len(bytes), processBytes)
		}

		totalLength += processBytes
	}
	for i := 0; i < len(bytes2d); i++ {
		bytes := bytes2d[i]

		//fmt.Printf("%04d \n", i)
		//if i == 119 {
		//	fmt.Printf("%04d \n", i)
		//}

		processBytes := session.ProcessMessage(bytes, func(msg dpcp.Message) {})
		if processBytes != len(bytes) {
			fmt.Printf("%04d bytes length: %d  process length: %d\n", i, len(bytes), processBytes)
		}

		totalLength += processBytes
	}
	for i := 0; i < len(bytes2d); i++ {
		bytes := bytes2d[i]

		//fmt.Printf("%04d \n", i)
		//if i == 119 {
		//	fmt.Printf("%04d \n", i)
		//}

		processBytes := session.ProcessMessage(bytes, func(msg dpcp.Message) {})
		if processBytes != len(bytes) {
			fmt.Printf("%04d bytes length: %d  process length: %d\n", i, len(bytes), processBytes)
		}

		totalLength += processBytes
	}
	for i := 0; i < len(bytes2d); i++ {
		bytes := bytes2d[i]

		//fmt.Printf("%04d \n", i)
		//if i == 119 {
		//	fmt.Printf("%04d \n", i)
		//}

		processBytes := session.ProcessMessage(bytes, func(msg dpcp.Message) {})
		if processBytes != len(bytes) {
			fmt.Printf("%04d bytes length: %d  process length: %d\n", i, len(bytes), processBytes)
		}

		totalLength += processBytes
	}
	for i := 0; i < len(bytes2d); i++ {
		bytes := bytes2d[i]

		//fmt.Printf("%04d \n", i)
		//if i == 119 {
		//	fmt.Printf("%04d \n", i)
		//}

		processBytes := session.ProcessMessage(bytes, func(msg dpcp.Message) {})
		if processBytes != len(bytes) {
			fmt.Printf("%04d bytes length: %d  process length: %d\n", i, len(bytes), processBytes)
		}

		totalLength += processBytes
	}

	since := time.Since(now)

	fmt.Println(totalLength)
	fmt.Println(since)

}
