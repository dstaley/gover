# Gover
Gover is a monitoring service for cgminer instances.

## Installation

### Pre-compiled executables
Coming Soon...

### From Source
```
go get github.com/dstaley/gover
go install github.com/dstaley/gover
```

## Setup
After installation, running `gover setup` begins the setup process.
```
~  gover setup
Looks like you don't have a database yet? Shall I create one for you?
y
Okay! Creating a database!
How often would you like to store data in the database (in seconds)? [120]

Now let's setup your configuration.
What's the hostname for your CGMiner instance? [localhost]

What port is the API listening in on? [4028]

Would you like to setup a Mobileminerapp connection?
y
Please enter your email address.
example@email.com
Please enter your application key.
apikey8888
What is the name of your rig?
gover
How often would you like to update Mobileminer (in seconds)? [120]

```
`gover setup` will then create a SQLite database and a `config.json` file in the directory from which you ran the command.

## Usage
Running `gover server` (or `gover s`) starts the server on port 8080. Gover will then log your rig's stats to the created database at the interval you specified.

## License
MIT