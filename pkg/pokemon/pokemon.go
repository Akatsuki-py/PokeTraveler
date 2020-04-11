package pokemon

// Pokemon - ゲーム内のポケモンのデータ
type Pokemon struct {
	*PokeData
	X     int
	Y     int
	Texts []string // 話しかけた時のテキスト
}

// Init - Pokemon packageの初期化処理
func Init() {
	initPokeDex()
}

// NewPokemon - Pokemonのコンストラクタ
func NewPokemon(ID, x, y int, texts []string) *Pokemon {
	p := &Pokemon{
		PokeData: PokeDex[ID],
		X:        x,
		Y:        y,
		Texts:    texts,
	}
	return p
}
