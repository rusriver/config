/*
-----------------------------------------------------------------------------
THIS IS UNABRIDGED DOCUMENTATION, PLEASE SEE THE README.md in the project root
for more details.
-----------------------------------------------------------------------------

Package config provides convenient access methods to configuration stored as
JSON or YAML.

Let's start with a simple YAML file config.yml:

	development:
	  database:
	    host: localhost
	  users:
	    - name: calvin
	      password: yukon
	    - name: hobbes
	      password: tuna
	production:
	  database:
	    host: 192.168.1.1

We can parse it using ParseYaml(), which will return a *Config instance on
success:

	c1 := (&config.InitContext{}).FromFile("config.yaml").Load().U()

An equivalent JSON configuration could be built using ParseJson():

	c1 := (&config.InitContext{}).FromFile("config.json").Load().U()

From now, we can retrieve configuration values using a path in dotted notation:

	// "localhost"
	host := c1.DotP("development.database.host").String()

	// or...

	// "192.168.1.1"
	host := c1.DotP("production.database.host").String()

Besides String(), other types can be fetched directly: Bool(), Float64(),
Int(), Map() and List(). All these methods will issue an error if the path
doesn't exist, or the value doesn't match or can't be converted to the
requested type.

A nested configuration can be fetched using DotP(). Here we get a new *Config
instance with a subset of the configuration:

	c2 := c2.DotP("development")

Then the inner values are fetched relatively to the subset:

	// "localhost"
	host := c2.DotP("database.host").String()

For lists, the dotted path must use an index to refer to a specific value.
To retrieve the information from a user stored in the configuration above:

	// map[string]interface{}{ ... }
	user1 := c1.DotP("development.users.0").Map()
	// map[string]interface{}{ ... }
	user2 := c1.DotP("development.users.1").Map()

	// or...

	// "calvin"
	name1 := c1.DotP("development.users.0.name").String()
	// "hobbes"
	name2 := c1.DotP("development.users.1.name").String()
*/
package config
