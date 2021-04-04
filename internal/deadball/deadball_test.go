package deadball

var team = Team{
	Players: [9]*Player{
		NewPlayer("Anna Test", Pitcher),
		NewPlayer("Bob Test", Catcher),
		NewPlayer("Clyde Test", First),
		NewPlayer("Doris Test", Second),
		NewPlayer("Elmer Test", Third),
		NewPlayer("Frank Test", Shortstop),
		NewPlayer("Gillian Test", Left),
		NewPlayer("Helen Test", Center),
		NewPlayer("Ian Test", Right),
	},
}
