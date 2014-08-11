package main

import (
	"fmt"
	"github.com/jmcvetta/neoism"
	"github.com/verdverm/neo4j-tutorials/common/reset"
)

var (
	db *neoism.Database
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	// resetDB()
	var err error
	db, err = neoism.Connect("http://localhost:7474/db/data")
	if err != nil {
		panic(err)
	}
	// initDB()
}

func resetDB() {
	reset.RemoveNeo4jDB()
	reset.StartNeo4jDB()
}

func initDB() {
	stmt := `
		CREATE (himym:TVShow { name: "How I Met Your Mother" })
		CREATE (himym_s1:Season { name: "HIMYM Season 1" })
		CREATE (himym_s1_e1:Episode { name: "Pilot" })
		CREATE (ted:Character { name: "Ted Mosby" })
		CREATE (joshRadnor:Actor { name: "Josh Radnor" })
		CREATE UNIQUE (joshRadnor)-[:PLAYED_CHARACTER]->(ted)
		CREATE UNIQUE (himym)-[:HAS_SEASON]->(himym_s1)
		CREATE UNIQUE (himym_s1)-[:HAS_EPISODE]->(himym_s1_e1)
		CREATE UNIQUE (himym_s1_e1)-[:FEATURED_CHARACTER]->(ted)
		CREATE (himym_s1_e1_review1 { title: "Meet Me At The Bar In 15 Minutes & Suit Up",
		  content: "It was awesome" })
		CREATE (wakenPayne:User { name: "WakenPayne" })
		CREATE (wakenPayne)-[:WROTE_REVIEW]->(himym_s1_e1_review1)<-[:HAS_REVIEW]-(himym_s1_e1)

		CREATE (marshall:Character { name: "Marshall Eriksen" })
		CREATE (robin:Character { name: "Robin Scherbatsky" })
		CREATE (barney:Character { name: "Barney Stinson" })
		CREATE (lily:Character { name: "Lily Aldrin" })
		CREATE (jasonSegel:Actor { name: "Jason Segel" })
		CREATE (cobieSmulders:Actor { name: "Cobie Smulders" })
		CREATE (neilPatrickHarris:Actor { name: "Neil Patrick Harris" })
		CREATE (alysonHannigan:Actor { name: "Alyson Hannigan" })
		CREATE UNIQUE (jasonSegel)-[:PLAYED_CHARACTER]->(marshall)
		CREATE UNIQUE (cobieSmulders)-[:PLAYED_CHARACTER]->(robin)
		CREATE UNIQUE (neilPatrickHarris)-[:PLAYED_CHARACTER]->(barney)
		CREATE UNIQUE (alysonHannigan)-[:PLAYED_CHARACTER]->(lily)
		CREATE UNIQUE (himym_s1_e1)-[:FEATURED_CHARACTER]->(marshall)
		CREATE UNIQUE (himym_s1_e1)-[:FEATURED_CHARACTER]->(robin)
		CREATE UNIQUE (himym_s1_e1)-[:FEATURED_CHARACTER]->(barney)
		CREATE UNIQUE (himym_s1_e1)-[:FEATURED_CHARACTER]->(lily)
		CREATE (himym_s1_e1_review2 { title: "What a great pilot for a show :)",
		  content: "The humour is great." })
		CREATE (atlasredux:User { name: "atlasredux" })
		CREATE (atlasredux)-[:WROTE_REVIEW]->(himym_s1_e1_review2)<-[:HAS_REVIEW]-(himym_s1_e1)

		CREATE (er:TVShow { name: "ER" })
		CREATE (er_s7:Season { name: "ER S7" })
		CREATE (er_s7_e17:Episode { name: "Peter's Progress" })
		CREATE (tedMosby:Character { name: "The Advocate " })
		CREATE UNIQUE (er)-[:HAS_SEASON]->(er_s7)
		CREATE UNIQUE (er_s7)-[:HAS_EPISODE]->(er_s7_e17)
		WITH er_s7_e17
		MATCH (actor:Actor),(episode:Episode)
		WHERE actor.name = "Josh Radnor" AND episode.name = "Peter's Progress"
		WITH actor, episode
		CREATE (keith:Character { name: "Keith" })
		CREATE UNIQUE (actor)-[:PLAYED_CHARACTER]->(keith)
		CREATE UNIQUE (episode)-[:FEATURED_CHARACTER]->(keith)
	`
	cq := neoism.CypherQuery{
		Statement: stmt,
	}
	err := db.Cypher(&cq)
	panicErr(err)
}

func main() {

	show := "How I Met Your Mother"
	getShowInfo(show)
	getShowInfoWithComments(show)

	getCharacterList(show)
	getActorList(show)

	getActorInfo("Josh Radnor")
	// listGraphData()
}

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

func listGraphData() {
	// query results
	res := []struct {
		From neoism.Node
		Rel  neoism.Relationship
		To   neoism.Node
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement: `
			MATCH (n)-[r]->(m)
			RETURN n AS From, r AS Rel, m AS To;
		`,
		Result: &res,
	}
	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	fmt.Println("Graph Data: ", len(res))
	for i, _ := range res {
		n := res[i]
		fmt.Printf("  [%d] %+v -> %+v -> %+v\n", i, n.From.Data, n.Rel.Data, n.To.Data)
	}
}
