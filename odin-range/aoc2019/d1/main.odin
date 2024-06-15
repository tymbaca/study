package main

import "core:fmt"
import "core:os"

LOGGING :: true

main :: proc() {
	wireset := load_instructions()
	program := create_and_run_program(wireset)
	closest, dist := find_closest(program.intersections[:])
	log(closest, dist)
}

load_instructions :: proc() -> WireSet {
	data, ok := os.read_entire_file("input.txt")
	if !ok do panic("can't open the file")

	wires, err := decode(data)
	if err != nil {
		panic(fmt.aprint(err))
	}

	log(wires)
	return wires
}

log :: proc(args: ..any) {
	when LOGGING {
		fmt.println(..args)
	}
}
