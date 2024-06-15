package main

import "core:fmt"
import "core:io"
import l "core:log"
import "core:strconv"
import "core:strings"

Error :: union {
	string,
}

Instruction :: struct {
	direction: Direction,
	steps:     int,
}

Direction :: enum {
	UP,
	DOWN,
	LEFT,
	RIGHT,
}

WireSet :: distinct [][]Instruction

decode :: proc(data: []byte) -> (wireset: WireSet, err: Error) {
	sdata := string(data)
	lines := strings.split(sdata, "\n")
	wireset = make(WireSet, len(lines))

	for line, li in lines {
		if len(line) == 0 do continue

		raw_insructions := strings.split(line, ",")
		instructions := make([]Instruction, len(raw_insructions))

		for inst, i in raw_insructions {
			dir := get_direction(inst) or_return
			steps := get_steps(inst) or_return
			instructions[i] = Instruction {
				direction = dir,
				steps     = steps,
			}
		}

		wireset[li] = instructions
	}

	return wireset, nil
}

get_direction :: proc(inst: string) -> (Direction, Error) {
	if len(inst) < 2 {
		return nil, fmt.aprintf("incorrect instruction: %s", inst)
	}

	switch inst[0] {
	case 'U':
		return .UP, nil
	case 'D':
		return .DOWN, nil
	case 'L':
		return .LEFT, nil
	case 'R':
		return .RIGHT, nil
	}

	return nil, fmt.aprintf("incorrect instruction: %s", inst)
}

get_steps :: proc(inst: string) -> (int, Error) {
	if len(inst) < 2 {
		return 0, fmt.aprintf("incorrect instruction: %s", inst)
	}

	steps_str := inst[1:]
	i, ok := strconv.parse_int(steps_str)
	if !ok {
		return 0, fmt.aprintf("can't decode steps number instruction: %s", inst)
	}

	return i, nil
}
