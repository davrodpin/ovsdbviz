# ovsdbviz

## How to build and run

```
$ cd $GOPATH
$ go get github.com/awalterschulze/gographviz
$ cd ./src/github.com/davrodpin/ovsdbviz # you have to clone or go get it
$ go install
$ ovsdbviz -schema=./examplesvswitch.ovsschema -out=/tmp/ovsdb.dot
$ brew install graphviz
$ dot -Tpng /tmp/ovsdb.dot -o /tmp/ovsdb.png
$ open /tmp/ovsdb.png
```

![OpenVSwitch Schema](https://github.com/davrodpin/ovsdbviz/blob/master/examples/vswitch.ovsschema.png)
