package msginfo

type MessageID int

func (m MessageID) Int() int {
	return int(m)
}

func MessageIDFromInt(id int) MessageID {
	return MessageID(id)
}

type ChatID int64

func (id ChatID) Int64() int64 {
	return int64(id)
}

func ChatIDFromInt64(id int64) ChatID {
	return ChatID(id)
}
