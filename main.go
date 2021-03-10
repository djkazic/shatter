package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/corvus-ch/shamir"
)

func main() {
	reconFlag := flag.Bool("reconstruct", false, "reconstruction mode")
	flag.Parse()
	if *reconFlag {
		var input string
		reader := bufio.NewReader(os.Stdin)
		inputBytes, err := ioutil.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		input = string(inputBytes)
		inputSplit := strings.Split(input, "\n")
		inputSplit = inputSplit[0 : len(inputSplit)-1]
		validMap := make(map[byte][]byte, 5)
		if strings.Contains(inputSplit[0], "BEGIN SHAMIR SET") && strings.Contains(inputSplit[len(inputSplit)-1], "END SHAMIR SET") {
			var validEntry []string
			validRows := inputSplit[1 : len(inputSplit)-1]
			if len(validRows) >= 3 {
				for _, row := range validRows {
					validEntry = strings.Split(row, "_")
					data, err := hex.DecodeString(validEntry[0])
					if err != nil {
						panic(err)
					}
					shard, err := hex.DecodeString(validEntry[1])
					if err != nil {
						panic(err)
					}
					validMap[data[0]] = shard
				}
				comb, err := shamir.Combine(validMap)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%s", comb)
			} else {
				fmt.Println("Error: insufficient number of shards to reconstruct")
			}
		} else {
			fmt.Println("INVALID INPUT")
			fmt.Println(inputSplit)
		}
	} else {
		var input string
		reader := bufio.NewReader(os.Stdin)
		inputBytes, err := ioutil.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		input = string(inputBytes)
		secret, err := shamir.Split([]byte(input), 5, 3)
		if err != nil {
			panic(err)
		}
		fmt.Println("===== BEGIN SHAMIR SET =====")
		for key, entry := range secret {
			fmt.Printf("%.2x_%x\n", key, entry)
		}
		fmt.Println("===== END SHAMIR SET =====")
	}
}
