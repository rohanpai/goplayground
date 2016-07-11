func Search(usr string) *ldap.SearchResult {
//Connect and Bind
	l, err := ldap.Dial(&#34;tcp&#34;, fmt.Sprintf(&#34;%s:%d&#34;, ldapServer, ldapPort))
	if err != nil {
		log.Fatalf(&#34;ERROR: %s\n&#34;, err.Error())
	}
	defer l.Close()
	// l.Debug = true

	err = l.Bind(user, passwd)
	if err != nil {
		log.Printf(&#34;ERROR: Cannot bind: %s\n&#34;, err.Error())
	}
	
	// Search 
	filter := fmt.Sprintf(&#34;(cn=%v)&#34;, usr)
	

	search := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		Attributes,
		nil)

	sr, err := l.Search(search)
	if err != nil {
		log.Fatalf(&#34;ERROR: %s\n&#34;, err.Error())
	}
	return sr
}