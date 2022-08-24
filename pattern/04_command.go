package pattern

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern

	Команда — это поведенческий паттерн, позволяющий заворачивать запросы или простые операции в отдельные объекты.
	Это позволяет откладывать выполнение команд, выстраивать их в очереди, а также хранить историю и делать отмену.
*/

// Command интерфейс для команд
type Command interface {
	execute() string
}

// StartHTTPServerCommand конкретная команда для старта сервера
type StartHTTPServerCommand struct {
	httpServer *ServerHTTP
}

func (c *StartHTTPServerCommand) execute() string {
	return c.httpServer.StartServer()
}

// StopHTTPServerCommand конкретная команда для остановки сервера
type StopHTTPServerCommand struct {
	httpServer *ServerHTTP
}

func (c *StopHTTPServerCommand) execute() string {
	return c.httpServer.StopServer()
}

// RestartHTTPServerCommand конкретная команда для перезапуска сервера
type RestartHTTPServerCommand struct {
	httpServer *ServerHTTP
}

func (c *RestartHTTPServerCommand) execute() string {
	return c.httpServer.RestartServer()
}

// ServerHTTP структура сервера
type ServerHTTP struct {
}

// StopServer метод для остановки сервера
func (s ServerHTTP) StopServer() string {
	return "Server stopped\n"
}

// StartServer метод для запуска сервера
func (s ServerHTTP) StartServer() string {
	return "Server starting...\n"
}

// RestartServer метод для перезапуска сервера
func (s ServerHTTP) RestartServer() string {
	return "Server restarted\n"
}

// CommandStorage для хранения команд
type CommandStorage struct {
	commands []Command
}

// SaveCommand для добавления новой команды в хранилище
func (cS *CommandStorage) SaveCommand(c Command) {
	cS.commands = append(cS.commands, c)
}

// ExecuteCommands для выполнения всех команд в хранилище
func (cS *CommandStorage) ExecuteCommands() string {
	var res string
	for _, com := range cS.commands {
		res += com.execute()
	}
	return res
}
