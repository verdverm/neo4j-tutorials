DM.02 TV Shows
==============

This tutorial is along the lines of the Getting Started tutorial.

### Information about a Show

Here is a function to get information about a show.

```Go
func getShowInfo(show string) {
	stmt := `
		MATCH (tvShow:TVShow)-[:HAS_SEASON]->(season)-[:HAS_EPISODE]->(episode)
		WHERE tvShow.name = {showSub}
		RETURN season.name, episode.name
	`
	params := neoism.Props{"showSub": show}

	res := []struct {
		Season  string `json:"season.name"`
		Episode string `json:"episode.name"`
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}

	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	if len(res) > 0 {
		fmt.Println("Show: ", show)
		fmt.Println("  ", res[0].Season, res[0].Episode)
	} else {
		fmt.Println("No results found")
	}

}
```

Same as above, extended to include comments.

``` Go
func getShowInfoWithComments(show string) {
	stmt := `
		MATCH (tvShow:TVShow)-[:HAS_SEASON]->(season)-[:HAS_EPISODE]->(episode)
		WHERE tvShow.name = {showSub}
		WITH season, episode
		OPTIONAL MATCH (episode)-[:HAS_REVIEW]->(review)
		RETURN season.name, episode.name, collect(review.content) AS Reviews
	`
	params := neoism.Props{"showSub": show}

	res := []struct {
		Season  string `json:"season.name"`
		Episode string `json:"episode.name"`
		Reviews []string
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}

	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	if len(res) > 0 {

		fmt.Println("Show & Reviews: ", show, len(res[0].Reviews))
		fmt.Println("Show: ", show)
		fmt.Println("  ", res[0].Season, res[0].Episode)
		for i, r := range res[0].Reviews {
			fmt.Printf("     %d:  %s\n", i, r)
		}
	} else {
		fmt.Println("No results found")
	}
}
```

Here is a function to get the charactor list for a show.

```Go
func getCharacterList(show string) {
	stmt := `
		MATCH (tvShow:TVShow)-[:HAS_SEASON]->()-[:HAS_EPISODE]->()-[:FEATURED_CHARACTER]->(character)
		WHERE tvShow.name = {showSub}
		RETURN DISTINCT character.name
	`
	params := neoism.Props{"showSub": show}

	res := []struct {
		Name string `json:"character.name"`
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}

	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	if len(res) > 0 {
		fmt.Println("Show: ", show)
		for _, n := range res {
			fmt.Println(" ", n.Name)
		}
	} else {
		fmt.Println("No results found")
	}

}
```

Same as above, extended to show both actor and character.

```Go
func getActorList(show string) {
	stmt := `
		MATCH (tvShow:TVShow)-[:HAS_SEASON]->()-[:HAS_EPISODE]->()-[:FEATURED_CHARACTER]->(character)<-[:PLAYED_CHARACTER]-(actor)
		WHERE tvShow.name = {showSub}
		RETURN DISTINCT actor.name, character.name
	`
	params := neoism.Props{"showSub": show}

	res := []struct {
		Name string `json:"actor.name"`
		Char string `json:"character.name"`
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}

	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	if len(res) > 0 {
		fmt.Println("Show: ", show)
		for _, n := range res {
			fmt.Printf("  %-24s  %-24s\n", n.Name, n.Char)
		}
	} else {
		fmt.Println("No results found")
	}
}
```


### Information about an Actor

Here is a function to get information about an actor.

```Go
func getActorInfo(name string) {
	stmt := `
		MATCH (actor:Actor)-[:PLAYED_CHARACTER]->(character)<-[:FEATURED_CHARACTER]-(episode), (episode)<-[:HAS_EPISODE]-(season)<-[:HAS_SEASON]-(tvshow)
		WHERE actor.name = {nameSub}
		RETURN tvshow.name AS Show, season.name AS Season, episode.name AS Episode, character.name AS Character
	`
	params := neoism.Props{"nameSub": name}

	res := []struct {
		Show      string
		Season    string
		Episode   string
		Character string
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}

	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	if len(res) > 0 {
		fmt.Println("Actor: ", name)
		for _, n := range res {
			fmt.Printf("  %-24s  %-16s  %-24s  %-16s\n", n.Show, n.Season, n.Episode, n.Character)
		}
	} else {
		fmt.Println("No results found")
	}
}
```

