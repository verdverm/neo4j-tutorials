DM.08 Find people based on mutual friends and groups
=====================================================

### Listing Functions

Standard practice

```Go
func listUsers() {...}
func listGroups() {...}
```

Sublists

```Go
func listGroupMembers(group string) {
	stmt := `
		MATCH (user:User)-[:member_of_group]->(group:Group)
		WHERE group.name = {groupSub}
		RETURN user.name
		ORDER BY user.name
	`

	params := neoism.Props{"groupSub": group}

	res := []struct {
		User string `json:"user.name"`
	}{}

	...
}

func listUserGroups(user string) {
	stmt := `
		MATCH (user:User)-[:member_of_group]->(group:Group)
		WHERE user.name = {userSub}
		RETURN group.name
		ORDER BY group.name
	`

	params := neoism.Props{"userSub": user}

	res := []struct {
		Group string `json:"group.name"`
	}{}

	...
}

func listUserFriends(user string) {
	stmt := `
		MATCH (user:User)-[:knows]->(friend:User)
		WHERE user.name = {userSub}
		RETURN friend.name
		ORDER BY friend.name
	`

	params := neoism.Props{"userSub": user}

	res := []struct {
		Friend string `json:"friend.name"`
	}{}

	...
}
```

### Finding Mutuals

Here, we are trying to find the number of paths
(through friends or groups) to the `user` to the `other` user.

`OPTIONAL MATCH` is used here so that one of the relationships
can be optional (equal to 0).

```Go
func findMutualsForUser(user string, other string) {
	stmt := `
		MATCH (me:User { name: {userSub} }),(other:User { name: {otherSub} })
		OPTIONAL MATCH pGroups=(me)-[:member_of_group]->(mg)<-[:member_of_group]-(other)
		OPTIONAL MATCH pMutualFriends=(me)-[:knows]->(mf)<-[:knows]-(other)
		RETURN other.name AS name, collect(mg.name) AS mutualGroups,
		  collect(DISTINCT mf.name) AS mutualFriends
		ORDER BY length(collect(DISTINCT mf.name)) DESC
	`

	// othersStr := "['" + others[0]
	// for i, o := range others {
	// 	if i < 1 {
	// 		continue
	// 	}
	// 	othersStr += "', '" + o
	// }
	// othersStr += "']"
	// fmt.Println("others: ", othersStr)

	// ERROR  othersStr doesn't substitute well
	// WHERE other.name IN {othersSub} --> ['Jill', 'Bob']
	// params := neoism.Props{"userSub": user, "othersSub": othersStr}

	params := neoism.Props{"userSub": user, "otherSub": other}

	res := []struct {
		Name    string   `json:"name"`
		Groups  []string `json:"mutualGroups"`
		Friends []string `json:"mutualFriends"`
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

	fmt.Println("Mutuals of", user, ":")
	if len(res) > 0 {
		fmt.Printf("%-12s  %-14s  %-14s\n", "Name", "Groups", "Friends")
		for _, n := range res {
			fmt.Printf("  %-10s  %v  %v\n", n.Name, n.Groups, n.Friends)
		}
	} else {
		fmt.Println("No results found")
	}
}
```
