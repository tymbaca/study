const std = @import("std");
const print = std.debug.print;
const hashmap = @import("hashmap.zig");

pub fn main() !void {
    // Prints to stderr (it's a shortcut based on `std.io.getStdErr()`)
    std.debug.print("All your {s} are belong to us.\n", .{"codebase"});

    // stdout is for the actual output of your application, for example if you
    // are implementing gzip, then only the compressed bytes should be sent to
    // stdout, not any debugging messages.
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    try stdout.print("Run `zig build test` to run the tests.\n", .{});

    try bw.flush(); // don't forget to flush!

    var hm = hashmap.HashSet{};
    stdout.print("%v", hm);
}

// test "simple test" {
//     // var list = std.ArrayList(i32).init(std.testing.allocator);
//     // defer list.deinit(); // try commenting this out and see if zig detects the memory leak!
//     // try list.append(42);
//     // try std.testing.expectEqual(@as(i32, 42), list.pop());
// }

// test "out of bounds, no safety" {
//     // @setRuntimeSafety(false);
//     // const a = [3]u8{ 1, 2, 3 };
//     // var index: u8 = 2;
//     // const b = a[index];
//     // print("{}\n", .{b});
//     // print("Hello\n", .{});
// }

test "arraylist indexing" {
    const String = []const u8;
    const Bucket = std.ArrayList(String);
    const Buckets = std.ArrayList(Bucket);

    var alloc = std.heap.page_allocator;
    var buck = Bucket.init(alloc);
    try buck.appendNTimes("Hello", 16);

    var bucks = Buckets.init(alloc);
    try bucks.appendNTimes(buck, 8);

    for (buck.items) |item| {
        print("{any}\n", .{item});
    }

    for (bucks.items) |ibuck| {
        print("{any}\n", .{ibuck});
    }
}
