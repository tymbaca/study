const debug = @import("std").debug;

// Toy example of "writers" that share the same "contract".
const FileWriter = struct {
    pub fn writeAll(_: FileWriter, bytes: []const u8) void {
        debug.print("[FileWriter] {s};", .{bytes});
    }
};
const MultiWriter = struct {
    pub fn writeAll(_: MultiWriter, bytes: []const u8) void {
        debug.print("[MultiWriter] {s};", .{bytes});
    }
};
const NullWriter = struct {
    pub fn writeAll(_: NullWriter, _: []const u8) void {}
};
// This writer differs from the other, so could be said to have a different "contract".
const BadWriter = struct {
    pub fn writeAll(_: BadWriter, val: i32) void {
        debug.print("[BadWriter] {d};", .{val});
    }
};

const example_01_duck_typing = struct {
    fn save(writer: anytype, bytes: []const u8) void {
        writer.writeAll(bytes);
    }

    pub fn main() void {
        var file_writer = FileWriter{};
        save(file_writer, "a");
        var multi_writer = MultiWriter{};
        save(multi_writer, "b");
        var null_writer = NullWriter{};
        save(null_writer, "c");
    }
};
