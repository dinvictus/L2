package pattern

import (
	"errors"
	"time"
)

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern

	Состояние — это поведенческий паттерн, позволяющий динамически изменять поведение объекта при смене его состояния.
	Поведения, зависящие от состояния, переходят в отдельные структуры.
	Первоначальный объект хранит ссылку на один из таких объектов-состояний и делегирует ему работу

	Паттерн должен применяться:

	- когда поведение объекта зависит от его состояния
	- поведение объекта должно изменяться во время выполнения программы
	- состояний достаточно много и использовать для этого условные операторы, разбросанные по коду, достаточно затруднительно
*/

// ServerState интерфейс для состояний сервера
type ServerState interface {
	start() (string, error)
	stop() (string, error)
	restart() (string, error)
}

// Server структура для сервера
type Server struct {
	curState     ServerState
	stopState    ServerState
	startState   ServerState
	restartState ServerState
}

func newServer() *Server {
	server := Server{}
	stopState := ServerStopState{server: &server}
	startState := ServerStartState{server: &server}
	restartState := ServerRestartState{server: &server}
	server.restartState = &restartState
	server.startState = &startState
	server.stopState = &stopState
	server.setState(&stopState)
	return &server
}

func (s *Server) start() (string, error) {
	return s.curState.start()
}

func (s *Server) stop() (string, error) {
	return s.curState.stop()
}

func (s *Server) restart() (string, error) {
	return s.curState.restart()
}

func (s *Server) setState(state ServerState) {
	s.curState = state
}

// ServerStartState Состояние запущенного сервера
type ServerStartState struct {
	server *Server
}

func (sStS *ServerStartState) start() (string, error) {
	return "", errors.New("server already started")
}

func (sStS *ServerStartState) stop() (string, error) {
	sStS.server.setState(sStS.server.stopState)
	return "Server stopped", nil
}

func (sStS *ServerStartState) restart() (string, error) {
	sStS.server.setState(sStS.server.restartState)
	time.Sleep(5 * time.Second)
	sStS.server.setState(sStS.server.startState)
	return "Server restarted", nil
}

// ServerStopState состояние остановленного сервера
type ServerStopState struct {
	server *Server
}

func (sSpS *ServerStopState) start() (string, error) {
	sSpS.server.setState(sSpS.server.startState)
	return "Server started", nil
}

func (sSpS *ServerStopState) stop() (string, error) {
	return "", errors.New("server already stopped")
}

func (sSpS *ServerStopState) restart() (string, error) {
	return "", errors.New("server stopped")
}

// ServerRestartState состояние перезапускающегося сервера
type ServerRestartState struct {
	server *Server
}

func (sRS *ServerRestartState) start() (string, error) {
	return "", errors.New("server restarting")
}

func (sRS *ServerRestartState) stop() (string, error) {
	return "", errors.New("server restarting")
}

func (sRS *ServerRestartState) restart() (string, error) {
	return "", errors.New("server already restarting")
}
