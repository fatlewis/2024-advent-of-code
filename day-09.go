package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"time"
)

type diskPointer struct {
	disk string
	diskIndex int
	fileID int
	spaceLeft int
}

type file struct {
	ID int
	size int
}

type diskLocation struct {
	files []*file
	emptySpace int
}

// Move to difference space. Update fileID if pointing to a file
func (p *diskPointer) move(offset int) {
	p.diskIndex += offset
	p.spaceLeft = toInt(string(p.disk[p.diskIndex]))
	if (p.diskIndex % 2 == 0) {
		p.fileID = p.diskIndex / 2
	}
}

func toInt(str string) int {
	result, _ := strconv.Atoi(str)
	return result
}

func main() {
	dat, _ := os.ReadFile("./assets/day09-input.txt")
	disk := strings.Trim(string(dat), "\n")

part1Start := time.Now()

	// read starts at 0 index, write at last index with file
	lastFileIndex := len(disk)-1
	// If the last digit on the disk is a space, substract 1 from the file index
	if(lastFileIndex % 2 != 0) {
		lastFileIndex -= 1
	}
	maxFileID := lastFileIndex/2

	checkSum := 0
	startPointer := &diskPointer{diskIndex: 1, fileID: 0, spaceLeft: toInt(string(disk[1])), disk: disk}
	endPointer := &diskPointer{diskIndex: lastFileIndex, fileID: maxFileID, spaceLeft: toInt(string(disk[lastFileIndex])), disk: disk}
	// Used to track the exact position for calculating the checksum
	exactDiskPosition := toInt(string(disk[0])) 
	for startPointer.diskIndex < endPointer.diskIndex {
		if (startPointer.spaceLeft > endPointer.spaceLeft) {
			// move files from end to start, calc checksum
			for i := 0; i < endPointer.spaceLeft; i += 1 {
				checkSum += exactDiskPosition * endPointer.fileID
				exactDiskPosition += 1
			}
			// update space left at start
			startPointer.spaceLeft -= endPointer.spaceLeft

			// move endPointer down to next file
			endPointer.move(-2)
		} else {
			// move as much as will fit to start, calc checksum
			for i := 0; i < startPointer.spaceLeft; i += 1 {
				checkSum += exactDiskPosition * endPointer.fileID
				exactDiskPosition += 1
				//fmt.Println("Updating checksum: ", exactDiskPosition, startPointer.fileID, checkSum)
			}
			endPointer.spaceLeft -= startPointer.spaceLeft

			// advance startPointer
			startPointer.move(1)
			if (startPointer.diskIndex == endPointer.diskIndex) {
				startPointer.spaceLeft = endPointer.spaceLeft
			}
			// calc checksum for file at startPointer
			for i := 0; i < startPointer.spaceLeft; i += 1 {
				checkSum += exactDiskPosition * startPointer.fileID
				exactDiskPosition += 1
			}
			// move startPointer to empty space
			startPointer.move(1)
		}
	}
	fmt.Println(checkSum)
fmt.Println(time.Since(part1Start))

part2Start := time.Now()
	// part2 algorithm
	// create disk array with a struct containing files & space values
	diskStructure := []*diskLocation{}
	for index, space := range disk {
		files := []*file{}
		emptySpace := toInt(string(space))
		if(index % 2 == 0) {
			files = append(files, &file{ID: index / 2, size: emptySpace})
			emptySpace = 0
		}
		diskStructure = append(diskStructure, &diskLocation{files, emptySpace})
	}

	// endPointer points at last file movement attempt
	endPointer = &diskPointer{diskIndex: lastFileIndex, fileID: maxFileID, spaceLeft: toInt(string(disk[lastFileIndex])), disk: disk}
	// While files remain untested
	for endPointer.diskIndex > 0 {
		// Try each initial empty space
		for i := 1; i < endPointer.diskIndex; i += 2 {
			// Always one file as we are iterating through the initial files
			currentFile := diskStructure[endPointer.diskIndex].files[0]

			// can we move it?
			if (diskStructure[i].emptySpace >= currentFile.size) {
				// move file
				diskStructure[i].files = append(diskStructure[i].files, currentFile)
				diskStructure[i].emptySpace = diskStructure[i].emptySpace - currentFile.size

				diskStructure[endPointer.diskIndex].files = []*file{}
				diskStructure[endPointer.diskIndex].emptySpace = currentFile.size
				break
			}
		}
		// Move endPointer to next file
		endPointer.move(-2)
	}

	// calculate checksum
	checkSum = 0
	exactDiskPosition = 0
	for _, structureAtIndex := range diskStructure {
		// if there are files, calculate these
		for _, fileAtIndex := range structureAtIndex.files {
			for i := 0; i < fileAtIndex.size; i += 1 {
				checkSum += exactDiskPosition * fileAtIndex.ID
				exactDiskPosition += 1
			}
		}
		// move along exactPosition depending on empty space
		exactDiskPosition += structureAtIndex.emptySpace
	}
	fmt.Println(checkSum)
fmt.Println(time.Since(part2Start))
}
