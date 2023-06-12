package base

// 随机获取一条消息
func GetRandomRecord(db *Dbr, record interface{}) error {
	return db.Order("RANDOM()").First(record).Error
}
