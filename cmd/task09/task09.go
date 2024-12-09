package main

import (
	"io"
	"os"
)

type BlockKind int8

const (
	BlockEmpty BlockKind = iota
	BlockNode
)

type Block struct {
	Kind BlockKind
	ID   int
}

type EmptySlot struct {
	Index  int
	Length int
}

func main() {
	blocks, err := readInput()
	if err != nil {
		return
	}

	println(defragmentAndGetChecksum(blocks))
	println(defragmentWholeFilesAndGetChecksum(blocks))
}

func defragmentAndGetChecksum(blocks []*Block) int {
	blocks = copyBlocks(blocks)
	emptyPointer := 0
	nodePointer := len(blocks) - 1

	for {
		emptyPointer = findFirst(BlockEmpty, emptyPointer, blocks)
		nodePointer = findLast(BlockNode, nodePointer, blocks)
		if emptyPointer >= nodePointer {
			break
		}

		blocks[emptyPointer].Kind, blocks[nodePointer].Kind = blocks[nodePointer].Kind, blocks[emptyPointer].Kind
		blocks[emptyPointer].ID, blocks[nodePointer].ID = blocks[nodePointer].ID, blocks[emptyPointer].ID
	}

	return checksum(blocks)
}

func defragmentWholeFilesAndGetChecksum(blocks []*Block) int {
	blocks = copyBlocks(blocks)
	nodePointerRight := len(blocks) - 1

	for nodePointerRight >= 0 {
		nodePointerRight = findLast(BlockNode, nodePointerRight, blocks)
		nodePointerLeft := findFileBeginning(blocks[nodePointerRight].ID, nodePointerRight, blocks)
		if nodePointerLeft < 0 {
			break
		}

		fileLength := nodePointerRight - nodePointerLeft + 1

		slots := findEmptySlots(nodePointerLeft, blocks)
		for _, slot := range slots {
			if slot.Length >= fileLength {
				for i := 0; i < fileLength; i++ {
					blocks[nodePointerLeft+i].Kind, blocks[slot.Index+i].Kind = blocks[slot.Index+i].Kind, blocks[nodePointerLeft+i].Kind
					blocks[nodePointerLeft+i].ID, blocks[slot.Index+i].ID = blocks[slot.Index+i].ID, blocks[nodePointerLeft+i].ID
				}
				break
			}
		}

		nodePointerRight = nodePointerLeft - 1
	}

	return checksum(blocks)
}

func checksum(blocks []*Block) int {
	sum := 0
	for i := 0; i < len(blocks); i++ {
		if blocks[i].Kind == BlockNode {
			sum += blocks[i].ID * i
		}
	}
	return sum
}

func findFirst(kind BlockKind, pos int, blocks []*Block) int {
	for i := pos; i < len(blocks); i++ {
		if blocks[i].Kind == kind {
			return i
		}
	}

	return len(blocks)
}

func findLast(kind BlockKind, pos int, blocks []*Block) int {
	for i := pos; i >= 0; i-- {
		if blocks[i].Kind == kind {
			return i
		}
	}

	return -1
}

func findFileBeginning(id int, pos int, blocks []*Block) int {
	for i := pos; i >= 0; i-- {
		if blocks[i].ID != id {
			return i + 1
		}
	}

	return -1
}

func findEmptySlots(until int, blocks []*Block) []*EmptySlot {
	var slots []*EmptySlot

	lastStart := 0
	counter := 0
	for i := 0; i < len(blocks) && i < until; i++ {
		if blocks[i].Kind == BlockEmpty {
			if counter == 0 {
				lastStart = i
			}
			counter++
		} else {
			if counter > 0 {
				slots = append(slots, &EmptySlot{lastStart, counter})
				lastStart = 0
				counter = 0
			}
		}
	}

	if counter > 0 {
		slots = append(slots, &EmptySlot{lastStart, counter})
	}

	return slots
}

func copyBlocks(blocks []*Block) []*Block {
	b2 := make([]*Block, len(blocks))
	for i := range blocks {
		b2[i] = &Block{Kind: blocks[i].Kind, ID: blocks[i].ID}
	}
	return b2
}

func readInput() ([]*Block, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		println("Failed to open input.txt: " + err.Error())
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var (
		blocks []*Block
		kind   = BlockNode
		id     = 0
	)

	for _, symbol := range content {
		digit := int(symbol - '0')

		if kind == BlockEmpty {
			for range digit {
				blocks = append(blocks, &Block{Kind: BlockEmpty})
			}

			kind = BlockNode
		} else {
			for range digit {
				blocks = append(blocks, &Block{Kind: BlockNode, ID: id})
			}

			kind = BlockEmpty
			id++
		}
	}

	return blocks, nil
}
