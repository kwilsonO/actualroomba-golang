package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"errors"
)


var allDirt map[string]int
var allPos  []string
var Rx, Ry, Dx, Dy int
var source *Vertex
var dirs string


type Vertex struct{
	X, Y int
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Send a file path noob")
	}

	lines, err := ParseLines(os.Args[1], func(s string)(string,bool){
		return s, true
	})

	if err != nil {
		fmt.Println("you're file format sucks", err);
		return
	}
	var s []string
	allDirt = make(map[string]int)
	for i, l := range lines {

		if i == (len(lines) - 1) {
			dirs = l
			continue
		}

		//Room Dimensions
		if i <= 1 {
			s = strings.Split(l, " ")

			x, err := strconv.Atoi(s[0])
			y, err := strconv.Atoi(s[1])

			if err != nil {
				fmt.Println("really can't even type integers")
				return
			}
			if i == 0 {
				Dx = x;
				Dy = y;
				continue
			}
			if i == 1 {
				Rx = x
				Ry = y;
				source = &Vertex{X: x, Y: y}
				continue
			}
		}

		allDirt[l]= 1;
	}

	ParseDirs(source, dirs)

	finalPos := allPos[len(allPos) - 1]
	numDirtCleaned := GetCleanedCount()

	fmt.Printf("%s\n", finalPos)
	fmt.Printf("Cleaned: %d\n", numDirtCleaned)

}

func GetCleanedCount() int {

	count := 0
	for _, v := range allPos {
		if allDirt[v] == 1 {
			count++
		}
	}

	return count

}

func IsSamePoint(v, u *Vertex) bool{

	return (v.X == u.X) && (v.Y == u.Y)
}

func ParseDirs(source *Vertex, dirs string){

	last := source
	for i := 0; i < len(dirs); i++ {
		last = MoveRoomba(last, dirs[i])
		//if out of bounds don't add no use to waste time checkign
		if _, err := CheckPos(last); err == nil {
			allPos = append(allPos, fmt.Sprintf("%d %d", last.X, last.Y))
		}
	}
}

func CheckPos(v *Vertex) (*Vertex, error) {

		if (v.X >= 0 && v.X < Dx) && (v.Y >= 0 && v.Y < Dy) {
			return v, nil
		}

		return nil, errors.New("out of bounds")
}

func MoveRoomba(v *Vertex, b byte) *Vertex {

		switch b {
		case 'N':
			return &Vertex{X: v.X, Y: v.Y + 1}
		case 'E':
			return &Vertex{X: v.X + 1, Y: v.Y}
		case 'S':
			return &Vertex{X: v.X, Y: v.Y - 1}
		case 'W':
			return &Vertex{X: v.X - 1, Y: v.Y}

		}

		return nil
}

func ParseLines(filePath string, parse func(string) (string, bool)) ([]string, error){
	inputFile, err := os.Open(filePath);
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile);
	var results []string
	for scanner.Scan() {
		if output, add := parse(scanner.Text()); add {
			results = append(results, output)
		}
	}

	if err  := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil

}
