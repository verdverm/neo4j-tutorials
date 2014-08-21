DM.04 Hyperedges
================

In this section, we will explore [hyperedges](http://en.wikipedia.org/wiki/Hypergraph) in the context of users, groups, and roles.

### The setup

![data-viz](http://docs.neo4j.org/chunked/stable/images/cypher-hyperedgecommongroups-graph.svg)


### Listing Users and Groups

We'll start with some simple listing functions for users and groups.

```Go
func listUsers() {
	stmt := `
		MATCH (user:User)
		RETURN user.name
		ORDER BY user.name
	`
	res := []struct {
		Name string `json:"user.name"`
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: nil,
		Result:     &res,
	}

	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	if len(res) > 0 {
		fmt.Println("Users:")
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Name)
		}
	} else {
		fmt.Println("No results found")
	}
}

func listGroups() {
	stmt := `
		MATCH (group:Group)
		RETURN group.name
		ORDER BY group.name
	`
	res := []struct {
		Name string `json:"group.name"`
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: nil,
		Result:     &res,
	}

	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	if len(res) > 0 {
		fmt.Println("Groups:")
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Name)
		}
	} else {
		fmt.Println("No results found")
	}
}
```

### Finding All Roles for a Group

Here is a function for listing all of the Roles associated with a group.
The first MATCH finds the group and associated hyperedges.
The second MATCH finds the roles connected to the hyperedge from the previous step.

```Go
func listGroupRoles(group string) {
	stmt := `
		MATCH (hyperedge)-[:hasGroup]->(group)
		WHERE group.name = {groupSub}
		MATCH (hyperedge)-[:hasRole]->(role)
		RETURN role.name
		ORDER BY role.name
	`

	params := neoism.Props{"groupSub": group}

	res := []struct {
		Name string `json:"role.name"`
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
		fmt.Println(group, " Roles:")
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Name)
		}
	} else {
		fmt.Println("No results found")
	}

}
```

We can also include the `User` who is assigned to the `Role`
by extending the second `MATCH` clause, as follows...

```Go
func listGroupRolesUser(group string) {
	stmt := `
		MATCH (hyperedge)-[:hasGroup]->(group)
		WHERE group.name = {groupSub}
		MATCH (user:User)-[:hasRoleInGroup]->(hyperedge)-[:hasRole]->(role)
		WITH user, role
		ORDER BY user.name
		RETURN role.name, collect(user.name) AS users
		ORDER BY role.name
	`

	params := neoism.Props{"groupSub": group}

	res := []struct {
		Role  string   `json:"role.name"`
		Users []string `json:"users"`
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
		fmt.Println(group, " Role / User:")
		for _, n := range res {
			fmt.Printf("  %-24s  %v\n", n.Role, n.Users)
		}
	} else {
		fmt.Println("No results found")
	}
}
```


### Finding All Groups and Roles for a User

Given a user's name, what groups are they in and
what roles do they have in each group.

``` Go
func listUserGroupRoles(user string) {
	stmt := `
		MATCH (user:User)-[:hasRoleInGroup]->(hyperedge)
		WHERE user.name = {userSub}
		MATCH (hyperedge)-[:hasGroup]->(group), (hyperedge)-[:hasRole]->(role)
		RETURN group.name, collect(role.name) AS Roles
		ORDER BY group.name
	`

	params := neoism.Props{"userSub": user}

	res := []struct {
		Group string   `json:"group.name"`
		Role  []string `json:"Roles"`
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
		fmt.Println(user, " Groups / Roles:")
		for _, n := range res {
			fmt.Printf("  %-24s  %-v\n", n.Group, n.Role)
		}
	} else {
		fmt.Println("No results found")
	}

}
```


### Finding Role of a User in a Group

Given a `User` and a `Group`, list the roles they have in that group.

```Go
func listUserRolesInGroup(user, group string) {
	stmt := `
		MATCH (user:User)-[:hasRoleInGroup]->(hyperedge)-[:hasGroup]->(group:Group)
		WHERE user.name = {userSub} AND group.name = {groupSub}
		MATCH (hyperedge)-[:hasRole]->(role)
		RETURN role.name
		ORDER BY role.name
	`

	params := neoism.Props{"userSub": user, "groupSub": group}

	res := []struct {
		Role string `json:"role.name"`
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
		fmt.Println(user, "-", group, " Roles:")
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Role)
		}
	} else {
		fmt.Println("No results found")
	}

}
```



### Finding Common Groups and Roles shared between Users

Given two user names...

Find groups they are both in and their roles in that group.

```Go
func findCommonGroupsForUsers(user1, user2 string) {
	stmt := `
		MATCH
			(u1:User)-[:hasRoleInGroup]->(hyper1)-[:hasGroup]->(group:Group),
			(hyper1)-[:hasRole]->(role1),
			(u2:User)-[:hasRoleInGroup]->(hyper2)-[:hasGroup]->(group:Group),
			(hyper2)-[:hasRole]->(role2)
		WHERE u1.name = {user1Sub} AND u2.name = {user2Sub}
		RETURN group.name, collect(DISTINCT role1.name) AS role1s, collect(DISTINCT role2.name) AS role2s
		ORDER BY group.name
	`

	params := neoism.Props{"user1Sub": user1, "user2Sub": user2}

	res := []struct {
		Group string   `json:"group.name"`
		Role1 []string `json:"role1s"`
		Role2 []string `json:"role2s"`
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
		fmt.Printf("%-12s  %-14s  %-14s\n", "Group", user1, user2)
		for _, n := range res {
			fmt.Printf("  %-10s  %v  %v\n", n.Group, n.Role1, n.Role2)
		}
	} else {
		fmt.Println("No results found")
	}
}
```

Find the cases where the users have the same role, in different groups.

```Go
func findCommonRolesForUsers(user1, user2 string) {
	stmt := `
		MATCH
			(u1:User)-[:hasRoleInGroup]->(hyper1)-[:hasRole]->(role:Role),
			(hyper1)-[:hasGroup]->(g1),
			(u2:User)-[:hasRoleInGroup]->(hyper2)-[:hasRole]->(role:Role),
			(hyper2)-[:hasGroup]->(g2)
		WHERE u1.name = {user1Sub} AND u2.name = {user2Sub}
		RETURN role.name, collect(DISTINCT g1.name) AS group1s, collect(DISTINCT g2.name) AS group2s
		ORDER BY role.name
	`

	params := neoism.Props{"user1Sub": user1, "user2Sub": user2}

	res := []struct {
		Role   string   `json:"role.name"`
		Group1 []string `json:"group1s"`
		Group2 []string `json:"group2s"`
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
		fmt.Printf("%-12s  %-16s  %-16s\n", "Role", user1, user2)
		for _, n := range res {
			fmt.Printf("  %-10s  %v  %v\n", n.Role, n.Group1, n.Group2)
		}
	} else {
		fmt.Println("No results found")
	}
}
```

Find the cases where multiple users share the same role in the same group.

```Go
func findCommonGroupRoles() {
	stmt := `
		MATCH
			(user:User)-[:hasRoleInGroup]->(hyper)-[:hasRole]->(role:Role),
			(hyper)-[:hasGroup]->(group:Group)
		WITH user
		ORDER BY user.name
		MATCH
			(u1:User)-[:hasRoleInGroup]->(hyper1)-[:hasRole]->(role:Role),
			(hyper1)-[:hasGroup]->(group:Group),
			(u2:User)-[:hasRoleInGroup]->(hyper2)-[:hasRole]->(role:Role),
			(hyper2)-[:hasGroup]->(group:Group)
		WHERE u1.name <> u2.name
		RETURN group.name, role.name, collect(DISTINCT user.name) AS users
		ORDER BY group.name, role.name


	`

	res := []struct {
		Group string   `json:"group.name"`
		Role  string   `json:"role.name"`
		Users []string `json:"users"`
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: nil,
		Result:     &res,
	}

	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	if len(res) > 0 {
		fmt.Printf("%-12s  %-12s  %-12s\n", "Group", "Role", "Users")
		for _, n := range res {
			fmt.Printf("  %-10s  %-12s  %v\n", n.Group, n.Role, n.Users)
		}
	} else {
		fmt.Println("No results found")
	}
}
```

