package main

import &#34;fmt&#34;

// Factory A Factory which defines a GetFactory method to return a factory instance
type Factory interface {
	GetFactory() Factory
}

// Database describes a database which should be created. The client will choose
// the type of database they want.
type Database interface {
	RunQuery(sql string) (result string)
	AddData(sql string, data string) error
}

// Filesystem describes a filesystem with the capabilities to open and read from
// files.
type Filesystem interface {
	CreateFile(fileName string) (bool, error)
	GetFile(fileName string) (File, error)
}

// File type
type File struct {
	content  string
	fileName string
}

// MongoDB a concrete Database implementation
type MongoDB struct {
	database map[string]string
}

// OracleDB a concrete Database implementation
type OracleDB struct {
	database map[string]string
}

// ZFS concrete filesystem
type ZFS struct {
	files map[string]File
}

// NTFS concrete filesystem
type NTFS struct {
	files map[string]File
}

// Databases abstract database groupper which acts like it is a mongoDb or an OracleDB
type Databases struct {
	*MongoDB
	*OracleDB
}

// Filesystems abstract filesystem groupper which acts like a ZFS or an NTFS
type Filesystems struct {
	*ZFS
	*NTFS
}

// RunQuery runs a query on a OracleDB type database
func (odb OracleDB) RunQuery(sql string) (result string) {
	return odb.database[sql]
}

// RunQuery runs a query on a MongoDB type database
func (mdb MongoDB) RunQuery(sql string) (result string) {
	return mdb.database[sql]
}

// AddData add data to the database
func (odb *OracleDB) AddData(sql string, data string) error {
	odb.database[sql] = data
	return nil
}

// AddData add data to the database
func (mdb *MongoDB) AddData(sql string, data string) error {
	mdb.database[sql] = data
	return nil
}

// CreateFile creates file
func (zfs *ZFS) CreateFile(fileName string) (bool, error) {
	file := File{content: &#34;ZFS Content&#34;, fileName: fileName}
	zfs.files[fileName] = file
	if _, ok := zfs.files[fileName]; ok {
		return true, nil
	}

	return false, fmt.Errorf(&#34;Something bad happened.&#34;)
}

// GetFile gets a file
func (zfs *ZFS) GetFile(fileName string) (File, error) {
	if f, ok := zfs.files[fileName]; ok {
		return f, nil
	}

	return File{}, fmt.Errorf(&#34;File is still there.&#34;)
}

// CreateFile creates file
func (ntfs *NTFS) CreateFile(fileName string) (bool, error) {
	file := File{content: &#34;NTFS Content&#34;, fileName: fileName}
	ntfs.files[fileName] = file
	if _, ok := ntfs.files[fileName]; ok {
		return true, nil
	}

	return false, fmt.Errorf(&#34;Something bad happened.&#34;)
}

// GetFile gets a file
func (ntfs *NTFS) GetFile(fileName string) (File, error) {
	if f, ok := ntfs.files[fileName]; ok {
		return f, nil
	}

	return File{}, fmt.Errorf(&#34;File is still there.&#34;)
}

// GetFactory Create a Factory for the databases
func (db Databases) GetFactory() Factory {
	return Databases{&amp;MongoDB{make(map[string]string)}, &amp;OracleDB{make(map[string]string)}}
}

// GetFactory Create a Factory for the filesystems
func (fs Filesystems) GetFactory() Factory {
	return Filesystems{&amp;ZFS{make(map[string]File)}, &amp;NTFS{make(map[string]File)}}
}

// GetFactory is an abstract factory which returns factories
func GetFactory(factoryType string) Factory {
	switch factoryType {
	case &#34;database&#34;:
		return Databases{}.GetFactory()
	case &#34;filesystems&#34;:
		return Filesystems{}.GetFactory()
	}
	return nil
}

// GetDatabase This works like a concrete database factory. It returns a concrete
// database based on databaseType
func GetDatabase(databaseType string) Database {
	f := GetFactory(&#34;database&#34;)
	switch databaseType {
	case &#34;mongo&#34;:
		return f.(Databases).MongoDB
	case &#34;oracle&#34;:
		return f.(Databases).OracleDB
	}
	return nil
}

// GetFilesystems This works like a concrete filesystem factory. Returns a concrete
// filesystem
func GetFilesystems(filesystemType string) Filesystem {
	f := GetFactory(&#34;filesystems&#34;)
	switch filesystemType {
	case &#34;zfs&#34;:
		return f.(Filesystems).ZFS
	case &#34;ntfs&#34;:
		return f.(Filesystems).NTFS
	}
	return nil
}

func main() {
	database := GetDatabase(&#34;mongo&#34;)
	database.AddData(&#34;bla&#34;, &#34;data bla&#34;)
	fmt.Println(&#34;database: &#34;, database.RunQuery(&#34;bla&#34;))

	filesystem := GetFilesystems(&#34;zfs&#34;)
	filesystem.CreateFile(&#34;bla&#34;)
	file, _ := filesystem.GetFile(&#34;bla&#34;)
	fmt.Println(&#34;file content: &#34;, file.content)
}
