package main

import "core:fmt"
import l "core:log"
import "core:os"

import "core:time"
import rl "vendor:raylib"

_screenWidth: i32 = 1000
_screenHeight: i32 = 700

main :: proc() {
	rl.InitWindow(_screenWidth, _screenHeight, "dungeons")
	rl.SetTargetFPS(60)

	rl.BeginDrawing()
	rl.ClearBackground(rl.DARKGRAY)
	run()
	rl.EndDrawing()
	for !rl.WindowShouldClose() {
	}
}


run :: proc() {
	context.logger = l.create_console_logger()
	defer l.destroy_console_logger(context.logger)

	args := os.args
	filename := "input.txt"
	if len(args) == 2 {
		filename = args[1]
	}

	wireset := load_instructions(filename)
	program := create_and_run_program(wireset)
	closest, dist, err := find_closest(program.intersections[:])
	if err != nil {
		l.error(err)
	}

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
