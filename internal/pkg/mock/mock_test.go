package mock

import "testing"

func TestNewLogger(t *testing.T) {
	log := NewLogger()

	if log == nil {
		t.Fatal("unexpected nil Logger")
	}
}

func TestLogger_WrappedLogger(t *testing.T) {
	log := NewLogger().WrappedLogger()

	if log == nil {
		t.Fatal("unexpected nil wrapped *logrus.Logger")
	}
}

func TestLogger_UpdateLevel(t *testing.T) {
	NewLogger().UpdateLevel()
}

func TestNewSession(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatal(err)
	}

	_, err = session.User("testUser")
	if err != nil {
		t.Error(err)
	}

	_, err = session.Channel("testChannel")
	if err != nil {
		t.Error(err)
	}

	_, err = session.GuildRoleCreate("testGuild")
	if err != nil {
		t.Error(err)
	}

	defer SessionClose(t, session)
}