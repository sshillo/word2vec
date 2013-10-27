package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
)

const maxSize int = 2000
const N int = 40
const maxW int = 50
const maxRead int = 2000

func main() {
	file, err := os.Open("freebase-vectors-skipgram1000-en.bin")
	if err != nil {
		log.Fatal(err)
	}

	var words, size int
	fmt.Fscanf(file, "%d", &words)
	fmt.Fscanf(file, "%d", &size)

	fmt.Printf("words: %d, size: %d\n", words, size)

	var ch string
	vocab := make([]string, words)
	M := make([]float32, size*words)
	for b := 0; b < words; b++ {
		tmp := make([]float32, size)
		fmt.Fscanf(file, "%s%c", &vocab[b], &ch)
		//fmt.Printf("vocab: %s\n", vocab[b])
		binary.Read(file, binary.LittleEndian, tmp)

		length := 0.0
		for _, v := range tmp {
			length += float64(v * v)
		}
		length = math.Sqrt(length)
		//fmt.Printf("length %f", length)

		for i, _ := range tmp {
			tmp[i] /= float32(length)
			//if (i < 20 && b < 20) {fmt.Printf("M %f\n",tmp[i])}
		}
		copy(M[b*size:b*size+size], tmp)
		//fmt.Printf("M %v\n",M)
	}
	userInput := []string{"/en/united_states", "/en/canada"}
	inputPosition := make([]int, 100)
	for k, v := range userInput {
		var b int
		for b = 0; b < words; b++ {
			if vocab[b] == v {
				break
			}
		}
		inputPosition[k] = b
		fmt.Printf("Word %v Position %v \n", v, b)
	}
	vec := make([]float32, maxSize)
	for i, _ := range userInput {
		for j := 0; j < size; j++ {
			//fmt.Printf("bi %v, M: %v\n", inputPosition[i],M[j + inputPosition[i]*size] )
			vec[j] += M[j+inputPosition[i]*size]
		}
	}
	//fmt.Printf("vec %v\n",vec)

	length := 0.0
	for _, v := range vec {
		length += float64(v * v)
	}
	//fmt.Printf("pre length %v\n", length)
	length = math.Sqrt(length)
	//fmt.Printf("length %v\n", length)

	for k, _ := range vec {
		vec[k] /= float32(length)
	}

	type WordData struct {
		Distance float64
		Word     string
	}

	bestWords := make([]WordData, N)

	for i := 0; i < words; i++ {
		c := 0
		for _, v := range inputPosition {
			if v == i {
				c = 1
			}
		}
		if c == 1 {
			continue
		}
		dist := 0.0
		for j := 0; j < size; j++ {
			dist += float64(vec[j] * M[j+i*size])
		}
		//fmt.Printf("distance %v, word %s\n",dist, vocab[i])
		for j := 0; j < N; j++ {
			if dist > bestWords[j].Distance {
				for d := N - 1; d > j; d-- {
					bestWords[d] = bestWords[d-1]
				}
				bestWords[j] = WordData{dist, vocab[i]}
				break
			}
		}
	}

	fmt.Printf("Best Words %+v\n", bestWords)
}
