package main

import "core:fmt"
import l "core:log"
import "core:os"

main :: proc() {
	context.logger = l.create_console_logger()
	defer l.destroy_console_logger(context.logger)

	args := os.args
	filename := "input.txt"
	if len(args) == 2 {
		filename = args[1]
	}

	wireset := load_instructions(filename)
	program := create_and_run_program(wireset)
	closest, dist := find_closest(program.intersections[:])
	l.info(dist)
}

load_instructions :: proc(filename: string) -> WireSet {
	data, ok := os.read_entire_file(filename)
	if !ok do panic("can't open the file")

	wires, err := decode(data)
	if err != nil {
		panic(fmt.aprint(err))
	}

	return wires
}
