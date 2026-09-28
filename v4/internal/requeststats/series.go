package requeststats

import "time"

func BuildDailySeries(fromTime time.Time, lastDates uint32, rows []DailyItem) []DailyItem {
	byDate := make(map[string]uint64, len(rows))
	for _, row := range rows {
		byDate[row.Date] = row.Count
	}
	data := make([]DailyItem, 0, lastDates)
	for i := 0; i < int(lastDates); i++ {
		day := fromTime.AddDate(0, 0, i)
		date := day.Format("2006-01-02")
		data = append(data, DailyItem{Date: date, Count: byDate[date]})
	}
	return data
}

func BuildHourlySeries(fromTime time.Time, rows []DailyItem) []DailyItem {
	byHour := make(map[string]uint64, len(rows))
	for _, row := range rows {
		byHour[row.Date] = row.Count
	}
	loc := fromTime.Location()
	now := time.Now().In(loc)
	endHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, loc)
	startHour := time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), 0, 0, 0, 0, loc)
	data := make([]DailyItem, 0)
	for cursor := startHour; !cursor.After(endHour); cursor = cursor.Add(time.Hour) {
		key := cursor.Format("2006-01-02 15:00")
		data = append(data, DailyItem{Date: key, Count: byHour[key]})
	}
	return data
}

func BuildMinuteSeries(fromTime time.Time, rows []DailyItem) []DailyItem {
	byMinute := make(map[string]uint64, len(rows))
	for _, row := range rows {
		byMinute[row.Date] = row.Count
	}
	loc := fromTime.Location()
	now := time.Now().In(loc)
	endMinute := now.Truncate(bucketDuration)
	startMinute := time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), 0, 0, 0, 0, loc)
	data := make([]DailyItem, 0)
	for cursor := startMinute; !cursor.After(endMinute); cursor = cursor.Add(bucketDuration) {
		key := cursor.Format("2006-01-02 15:04")
		data = append(data, DailyItem{Date: key, Count: byMinute[key]})
	}
	return data
}

func BuildUserDailySeries(fromTime time.Time, lastDates uint32, rows []UserDailyItem) []UserDailyItem {
	byDate := make(map[string]UserDailyItem, len(rows))
	for _, row := range rows {
		byDate[row.Date] = row
	}
	data := make([]UserDailyItem, 0, lastDates)
	for i := 0; i < int(lastDates); i++ {
		day := fromTime.AddDate(0, 0, i)
		date := day.Format("2006-01-02")
		item := byDate[date]
		data = append(data, UserDailyItem{
			Date: date, Interactions: item.Interactions, ActiveUsers: item.ActiveUsers,
		})
	}
	return data
}

func BuildUserHourlySeries(fromTime time.Time, rows []UserDailyItem) []UserDailyItem {
	byHour := make(map[string]UserDailyItem, len(rows))
	for _, row := range rows {
		byHour[row.Date] = row
	}
	loc := fromTime.Location()
	now := time.Now().In(loc)
	endHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, loc)
	startHour := time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), 0, 0, 0, 0, loc)
	data := make([]UserDailyItem, 0)
	for cursor := startHour; !cursor.After(endHour); cursor = cursor.Add(time.Hour) {
		key := cursor.Format("2006-01-02 15:00")
		item := byHour[key]
		data = append(data, UserDailyItem{
			Date: key, Interactions: item.Interactions, ActiveUsers: item.ActiveUsers,
		})
	}
	return data
}

func BuildUserMinuteSeries(fromTime time.Time, rows []UserDailyItem) []UserDailyItem {
	byMinute := make(map[string]UserDailyItem, len(rows))
	for _, row := range rows {
		byMinute[row.Date] = row
	}
	loc := fromTime.Location()
	now := time.Now().In(loc)
	endMinute := now.Truncate(bucketDuration)
	startMinute := time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), 0, 0, 0, 0, loc)
	data := make([]UserDailyItem, 0)
	for cursor := startMinute; !cursor.After(endMinute); cursor = cursor.Add(bucketDuration) {
		key := cursor.Format("2006-01-02 15:04")
		item := byMinute[key]
		data = append(data, UserDailyItem{
			Date: key, Interactions: item.Interactions, ActiveUsers: item.ActiveUsers,
		})
	}
	return data
}
