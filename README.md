# ovsdbviz

## How to build and run

```
$ cd $GOPATH/src
$ mkdir -p github.com/davrodpin && cd github.com/davrodpin
$ git clone https://github.com/davrodpin/ovsdbviz.git
$ cd ovsdbviz
$ go install
$ ovsdbviz -schema=./examplesvswitch.ovsschema -out=/tmp/ovsdb.dot
$ brew install graphviz
$ dot -Tpng /tmp/ovsdb.dot -o /tmp/ovsdb.png
$ open /tmp/ovsdb.png
```

![OpenVSwitch Schema](https://github.com/davrodpin/ovsdbviz/blob/master/examples/vswitch.ovsschema.png)
