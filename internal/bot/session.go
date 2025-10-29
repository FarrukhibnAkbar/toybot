package bot

func (b *Bot) getSession(userID int64) *Session {
	if s, ok := b.sessions[userID]; ok {
		return s
	}
	// Agar yo‘q bo‘lsa, yangi session yaratamiz
	s := &Session{State: StateNone}
	b.sessions[userID] = s
	return s
}

func (b *Bot) resetSession(userID int64) {
	delete(b.sessions, userID)
}
