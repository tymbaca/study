const std = @import("std");
const eql = std.mem.eql;
const print = std.debug.print;
const log = std.log;
const ArrayList = std.ArrayList;

pub const HashSet = struct {
    buckets: Buckets,
    allocator: std.mem.Allocator = std.heap.page_allocator,

    pub fn init() !HashSet {
        var hs = HashSet{ .buckets = undefined };

        hs.buckets = Buckets.init(hs.allocator);
        for (0..4) |_| {
            log.debug("Created bucket", .{});
            var buck = Bucket.init(hs.allocator);
            try hs.buckets.append(buck);
        }

        return hs;
    }

    pub fn add(self: *HashSet, val: String) !void {
        var h = hash(val);
        var bucket_index = h % self.buckets.items.len;
        var bucket = self.buckets.items[bucket_index];

        for (bucket.items) |item| {
            if (eql(u8, item, val)) {
                return Error.AlreadyAdded;
            }
        }
    }

    pub fn get(self: *HashSet, val: []const u8) void {
        _ = self;
        _ = val;
    }

    pub fn remove(self: *HashSet, val: []const u8) void {
        _ = self;
        _ = val;
    }

    fn hash(val: []const u8) u64 {
        _ = val;
        return 65;
    }
};

const Buckets = ArrayList(Bucket);
const Bucket = ArrayList(String);
const String = []const u8;

const Error = error{
    AlreadyAdded,
};

test "init" {
    var h = try HashSet.init();

    for (h.buckets.items) |item| {
        print("bucket: {any}\n", .{item});
    }
}

test "add" {
    var h = try HashSet.init();
    try h.add("val");
    try h.add("val");
}
