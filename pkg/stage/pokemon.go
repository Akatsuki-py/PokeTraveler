package stage

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Akatsuki-py/PokeTraveler/pkg/pokemon"
)

type Pokemons struct {
	List []*JSONPokemon `json:"pokemons"`
}

// JSONPokemon JSONから得られるPokemonの情報
type JSONPokemon struct {
	ID    int      `json:"id"`
	X     int      `json:"x"`
	Y     int      `json:"y"`
	Texts []string `json:"texts"`
}

func (stage *Stage) GetPokemon(x, y int) (target *pokemon.Pokemon) {
	for _, p := range stage.Pokemons {
		if p.X*16 == x && p.Y*16 == y {
			target = p
			break
		}
	}
	return target
}

func (stage *Stage) loadPokemons(filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		stage.Pokemons = []*pokemon.Pokemon{}
		return
	}

	raw := new(Pokemons)
	if err := json.Unmarshal(file, raw); err != nil {
		panic(err)
	}

	pokemons := make([]*pokemon.Pokemon, len(raw.List))
	for i, p := range raw.List {
		pokemons[i] = pokemon.NewPokemon(p.ID, p.X, p.Y, p.Texts)
	}

	stage.Pokemons = pokemons
}
