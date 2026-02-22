package messageprocessor

func (m *MessageProcessor) EscapeMarkdown(s string) string {
	return m.escaper.EscapeMarkdown(s)
}
