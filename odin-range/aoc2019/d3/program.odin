package main

import l "core:log"
import "core:math"
import rl "vendor:raylib"

Coordinates :: distinct [2]int
Mark :: struct {
	wire_id: int,
}

Program :: struct {
	coord:         Coordinates,
	wireset:       WireSet,
	world:         map[Coordinates]Mark,
	intersections: [dynamic]Coordinates,
}

colors := [?]rl.Color{rl.RED, rl.BLUE, rl.GREEN}

/*
    O------- +x
    |
    |
    |
   +y
*/
create_and_run_program :: proc(wireset_: WireSet) -> Program {
	p := Program {
		coord         = {0, 0},
		wireset       = wireset_,
		world         = make(map[Coordinates]Mark),
		intersections = make([dynamic]Coordinates),
	}
	using p

	l.info("starting program")
	for wire, wire_id in wireset {
		l.infof("starting wire_id: %d, wire len: %d", wire_id, len(wire))
		coord = {0, 0}
		for inst in wire {
			offset: Coordinates
			switch inst.direction {
			case .UP:
				offset = {0, -1}
			case .DOWN:
				offset = {0, 1}
			case .LEFT:
				offset = {-1, 0}
			case .RIGHT:
				offset = {1, 0}
			}

			for step in 0 ..< inst.steps {
				coord += offset // step
				rl.DrawPixel(
					_screenWidth / 2 + i32(coord.x) / 25,
					_screenHeight / 2 + i32(coord.y) / 25,
					colors[wire_id % len(colors)], // choose color
				)
				l.debugf("step to %v", coord)

				mark, marked := world[coord]
				if marked && mark.wire_id != wire_id { 	// check
					l.warnf(
						"got already marked at %v by wire_id %dm (i'm mark id %d)",
						coord,
						mark.wire_id,
						wire_id,
					)
					append(&intersections, coord) // save intersection if already was marked
				}

				world[coord] = {wire_id} // mark
			}
		}
	}

	return p
}

distance :: proc(src, dst: Coordinates) -> int {
	diff := dst - src
	return math.abs(diff.x) + math.abs(diff.y)
}

// returns distance as 2nd value
find_closest :: proc(dsts: []Coordinates, src: Coordinates = {0, 0}) -> (Coordinates, int, Error) {
	if len(dsts) == 0 {
		return {}, 0, "no dsts"
	}

	closest := dsts[0]
	closest_dist := distance(src, closest)
	for dst in dsts {
		cur_dist := distance(src, dst)
		if cur_dist < closest_dist {
			closest = dst
			closest_dist = cur_dist
		}
	}

	return closest, closest_dist, nil
}
