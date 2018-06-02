package main

type Player struct {
	session Session
	username string
	password string
	/* TODO: A lot of player data goes here. */
}

func NewPlayer(session Session, username string, password string) Player {
	return Player{session, username, password}
}