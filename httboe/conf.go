package httboe

import (
	"errors"
	toml "github.com/pelletier/go-toml"
	"log"
)

type ServConf struct {
	Host     string
	Port     int64
	Location []LocConf
}

type LocConf struct {
	Path string
	Type string
	Root string
	Auth bool
}

type Conf struct {
	Daemon bool
	Log    string
	Server ServConf
}

func (this *Conf) parseLocation(tree *toml.Tree, index int) (err error) {

	if !tree.Has("path") {
		return errors.New("No path configued")
	}

	this.Server.Location[index].Path = tree.Get("path").(string)
	log.Printf("Added path %s", this.Server.Location[index].Path)

	if !(tree.Has("static") || tree.Has("webdav")) {
		return errors.New("You must configure almost one type of handler")
	}

	if tree.Has("static") && tree.Has("webdav") {
		return errors.New("You must configure only one type of handler")
	}

	if tree.Has("static") {
		this.Server.Location[index].Type = "static"
		this.Server.Location[index].Root = tree.Get("static").(string)
		log.Printf("Added static handler with root %s", this.Server.Location[index].Root)
	}

	if tree.Has("webdav") {
		this.Server.Location[index].Type = "webdav"
		this.Server.Location[index].Root = tree.Get("webdav").(string)
		log.Printf("Added webdav handler with root %s", this.Server.Location[index].Root)
	}

	if tree.Has("auth") {
		this.Server.Location[index].Auth = tree.Get("auth").(bool)
	} else {
		this.Server.Location[index].Auth = false
	}

	return nil
	
}

// (server).*\{([^}]+)\} --> match all into server { }
func (this *Conf) parseServer(tree *toml.Tree) (err error) {

	//check for minimal prerequisites
	if !tree.Has("port") {
		return errors.New("No port configured")
	}

	this.Server.Port = tree.Get("port").(int64)
	log.Printf("parseServer port: %d", this.Server.Port)

	if !tree.Has("host") {
		this.Server.Host = "localhost"
	} else {
		this.Server.Host = tree.Get("host").(string)
	}
	log.Printf("parseServer host: %s", this.Server.Host)

	if !tree.Has("location") {
		return errors.New("No locations configured")
	}

	locations := tree.Get("location").([]*toml.Tree)
	this.Server.Location = make([]LocConf, len(locations))
	for x, val := range locations {
		err := this.parseLocation(val, x)
		if err != nil {
			return err
		}
	}
	return nil
	
}

func (this *Conf) Load(path string) (err error) {

	config, errs := toml.LoadFile(path)
	if errs != nil {
		return errs
	}

	if !config.Has("server") {
		return errors.New("No server configured")
	}

	srvTree := config.Get("server").(*toml.Tree)

	errs = this.parseServer(srvTree)
	return errs

}
