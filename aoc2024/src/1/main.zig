const std = @import("std");

fn getNumFromLine(line: []u8) !i32 {
    var first_digit: u8 = undefined;
    var second_digit: u8 = undefined;

    for (line) |char| {
        if (std.ascii.isDigit(char)) {
            if (first_digit == undefined) {
                first_digit = char;
            }
            second_digit = char;
        }
    }
    var togather = [_]u8{ first_digit, second_digit };
    var num = try std.fmt.parseInt(i32, togather[0..], 10);
    _ = num;
    return 47;
}

pub fn main() !void {
    var file = try std.fs.cwd().openFile("src/1/input.txt", .{});
    defer file.close();

    // var allocator = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    // _ = allocator;

    // var reader = file.reader();
    var buf: [1024]u8 = undefined;
    // var stream = std.io.fixedBufferStream(&buf);
    // try reader.streamUntilDelimiter(stream, '\n', .{});
    // std.debug.print("{s}", stream.buffer);

    while (try file.reader().readUntilDelimiterOrEof(&buf, '\n')) |line| {
        var num = try getNumFromLine(line);
        std.debug.print("{}\n", .{num});
    }

    // var sum: i32 = 0;

    // while () |line| {
    //     var num = try getNumFromLine(line);
    //     sum += num;
    // }

    // std.debug.print("{}", line);
}
