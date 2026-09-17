-- Chỉ rollback khi đã dừng read traffic: xóa bảng sẽ mất binding chống replay.
DROP TABLE quota_read_commands;
