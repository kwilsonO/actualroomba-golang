package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"errors"
)

//Room dimensions
var Dx, Dy int

//quick struct to avoid tons of string manipulation
type Vertex struct{
	X, Y int
}

func main() {

	//double check filepath arg coming in
	if len(os.Args) != 2 {
		fmt.Println("Please send a filepath arg")
		return
	}
	//begin parsing of file
	lines, err := ParseLines(os.Args[1], func(s string)(string,bool){
		return s, true
	})

	//failed parsing file
	if err != nil {
		fmt.Println("Invalid file format please see example", err);
		return
	}
	//check file at least has dimensions, start pos, and directions
	if len(lines) < 3 {
		fmt.Println("Not enough lines in input, try again.")
		return
	}

	//main program logic here
	finalPos, numDirtCleaned := BuildPositionsAndCalculate(lines)

	//print final pos and cleaned count
	fmt.Printf("%s\n", finalPos)
	fmt.Printf("%d\n", numDirtCleaned)

}

func BuildPositionsAndCalculate(lines []string) (string, int){

	//vars
	var s []string
	var dirs string
	var source *Vertex
	allDirt := make(map[string]int)

	//loop through each line of the file
	for i, l := range lines {

		//if the last line of the file then directions
		if i == (len(lines) - 1) {
			dirs = l
			continue
		}

		//if either room dimensions or start pos
		if i <= 1 {
			s = strings.Split(l, " ")

			x, err := strconv.Atoi(s[0])
			y, err := strconv.Atoi(s[1])

			if err != nil {
				fmt.Println("Line format not as expected try 'INT INT'")
				return "", 0
			}
			//if room dimensions
			if i == 0 {
				Dx = x;
				Dy = y;
				continue
			}

			//if start pos
			if i == 1 {
				source = &Vertex{X: x, Y: y}
				continue
			}
		}

		//otherwise it's a dirt spot
		allDirt[l]= 1;
	}

	//do the necessary calcs to find paths
	allPos := ParseDirs(source, dirs)

	var finalPos string
	if len(allPos) > 0 {
		//final pos is always the last pos in all pos
		finalPos = allPos[len(allPos) - 1]
	} else {
		finalPos = fmt.Sprintf("%d %d", source.X, source.Y)
	}

	//go find how many of allPos are present in allDirt map
	numDirtCleaned := GetCleanedCount(allPos, allDirt)

	return finalPos, numDirtCleaned
}


func GetCleanedCount(allPos []string, allDirt map[string]int) int {

	count := 0
	for _, v := range allPos {
		if _, ok := allDirt[v]; ok {
			//make sure we don't double up
			delete(allDirt, v)
			count++
		}
	}

	return count
}

func ParseDirs(source *Vertex, dirs string) []string{

	last := source
	var allPos []string
	for i := 0; i < len(dirs); i++ {
		last = MoveRoomba(last, dirs[i])
		//if out of bounds don't add no use to waste time checkign
		if _, err := CheckPos(last); err == nil {
			allPos = append(allPos, fmt.Sprintf("%d %d", last.X, last.Y))
		}
	}
	return allPos
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
