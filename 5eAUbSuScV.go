func Search(usr string) *ldap.SearchResult {
//Connect and Bind
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, ldapPort))
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	defer l.Close()
	// l.Debug = true

	err = l.Bind(user, passwd)
	if err != nil {
		log.Printf("ERROR: Cannot bind: %s\n", err.Error())
	}
	
	// Search 
	filter := fmt.Sprintf("(cn=%v)", usr)
	

	search := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		Attributes,
		nil)

	sr, err := l.Search(search)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	return sr
}