# Packer Parallels Plugin

This is a custom builder plugin for [Packer](http://www.packer.io/) using [Parallels](http://www.parallels.com/).

**Current status: Alpha - BROKEN!**

## Status
This is currently under development. Don't expect it to work. The only reason to run this plugin is if you wont to contribute to its development.

## Documentation
TODO

## Building / Installing
Install packer and add the following to ```~/.packerconfig```

```
{
  "builders": {
    "parallels-iso": "builder-parallels-iso"
   }
}
```

```
cd $GOPATH/src && git clone https://github.com/rickard-von-essen/packer-parallels
cd packer-parallels
go get -u github.com/mitchellh/gox
make
```

## Issues
If you find any bugs please open a issue at [github](https://github.com/rickard-von-essen/packer-parallels/issues). 

## Contributing
If you have any improvments open a pull request at [github](https://github.com/rickard-von-essen/packer-parallels/pulls). 

## License

This code is distributed under the MIT license, see _LICENSE_.