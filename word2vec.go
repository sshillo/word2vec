package word2vec

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
)

const maxSize int = 2000
const N int = 40

type Model struct {
	words, size int
	vocab       []string
	M           []float32
}

type WordData struct {
	Distance float64
	Word     string
}

func (m *Model) Load(fn string) {
	file, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fscanf(file, "%d", &m.words)
	fmt.Fscanf(file, "%d", &m.size)

   // m.words = 1000
	var ch string
	m.vocab = make([]string, m.words)
	m.M = make([]float32, m.size*m.words)
	for b := 0; b < m.words; b++ {
		tmp := make([]float32, m.size)
		fmt.Fscanf(file, "%s%c", &m.vocab[b], &ch)
		binary.Read(file, binary.LittleEndian, tmp)

		length := 0.0
		for _, v := range tmp {
			length += float64(v * v)
		}
		length = math.Sqrt(length)

		for i, _ := range tmp {
			tmp[i] /= float32(length)
		}
		copy(m.M[b*m.size:b*m.size+m.size], tmp)
	}
	file.Close()
}

func (m *Model) MostSimilar(seedWords []string) []WordData {
	inputPosition := make([]int, 100)
	for k, v := range seedWords {
		var b int
		for b = 0; b < m.words; b++ {
			if m.vocab[b] == v {
				break
			}
		}
		inputPosition[k] = b
		//fmt.Printf("Word %v Position %v \n", v, b)
	}
	vec := make([]float32, maxSize)
	for i, _ := range seedWords {
		for j := 0; j < m.size; j++ {
			vec[j] += m.M[j+inputPosition[i]*m.size]
		}
	}

	length := 0.0
	for _, v := range vec {
		length += float64(v * v)
	}
	length = math.Sqrt(length)

	for k, _ := range vec {
		vec[k] /= float32(length)
	}

	bestWords := make([]WordData, N)

	for i := 0; i < m.words; i++ {
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
		for j := 0; j < m.size; j++ {
			dist += float64(vec[j] * m.M[j+i*m.size])
		}

		for j := 0; j < N; j++ {
			if dist > bestWords[j].Distance {
				for d := N - 1; d > j; d-- {
					bestWords[d] = bestWords[d-1]
				}
				bestWords[j] = WordData{dist, m.vocab[i]}
				break
			}
		}
	}
	return bestWords
}
