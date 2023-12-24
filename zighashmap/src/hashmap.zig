const std = @import("std");
const eql = std.mem.eql;
const print = std.debug.print;
const ArrayList = std.ArrayList;

pub const HashSet = struct {
    buckets: []Bucket,
    allocator: std.mem.Allocator = std.heap.page_allocator,

    pub fn init() !HashSet {
        var hs = HashSet{ .buckets = undefined };

        hs.buckets = try hs.allocator.alloc(Bucket, 8);
        for (hs.buckets, 0..) |_, index| {
            hs.buckets[index] = try hs.allocator.alloc([]u8, 4);
        }

        return HashSet{ .buckets = undefined };
    }

    pub fn add(self: *HashSet, val: []const u8) !void {
        var h = hash(val);
        var bucket_index = h % self.buckets.len;
        var bucket = self.buckets[bucket_index];

        for (bucket) |item| {
            if (eql([]const u8, item, val)) {
                return .{ .err = Error.AlreadyAdded };
            }
        }

        for (bucket, 0..) |item, index| {
            if (item == undefined) {
                self.buckets[bucket_index][index] = item;
            }
        }
    }

    pub fn get(self: *HashSet, val: []u8) void {
        _ = self;
        _ = val;
    }

    pub fn remove(self: *HashSet, val: []u8) void {
        _ = self;
        _ = val;
    }

    fn hash(val: []const u8) u64 {
        _ = val;
        return 65;
    }
};

const Buckets = []Bucket;
const Bucket = [][]const u8;

const Error = error{
    AlreadyAdded,
};

pub fn main() !void {
    var h = try HashSet.init();
    try h.add("Hello");
    // print("{}\n", h.buckets);
}
