# Packer Parallels Plugin

This is a custom builder plugin for [Packer](http://www.packer.io/) using [Parallels Desktop for Mac](http://www.parallels.com/products/desktop/) and a patched version of the _post-processor-vagrant_.

![Parallels Logo](imgs/parallels_small.png)

**Current status: in Beta!**

[![Build Status](https://travis-ci.org/rickard-von-essen/packer-parallels.png?branch=master)](https://travis-ci.org/rickard-von-essen/packer-parallels)


## Status
This has just been developed and testing is just started. Some feaures are untested or currently broken. If you use this be prepaired to crash, debug and report issues (include debug logs) or submit patches (pull requests).

## Building / Installing

 - Install packer.
 - Install [Parallels Virtualization SDK 9 for Mac](http://download.parallels.com//desktop/v9/pde.hf1/ParallelsVirtualizationSDK-9.0.24172.951362.dmg), download and install it.
 - Install this plugin

```bash
mkdir -p $GOPATH/src/github.com/rickard-von-essen/
cd $GOPATH/src/github.com/rickard-von-essen/
git clone https://github.com/rickard-von-essen/packer-parallels.git
cd packer-parallels
make
```
 - Add the following to ```~/.packerconfig```

```
{
  "builders": {
    "parallels-iso": "builder-parallels-iso",
    "parallels-pvm": "builder-parallels-pvm"
   }
}
```
 - Be sure that ```$GOPATH/bin```is on your path _BEFORE_ packer otherwise the patched _post-processor-vagrant_ plugin won't be used.

## Usage

### Parallels Desktop Builders

 - [parallels-iso](https://github.com/rickard-von-essen/packer-parallels/blob/master/ISO.md)
 - [parallels-pvm](https://github.com/rickard-von-essen/packer-parallels/blob/master/PVM.md)

### Vagrant Post-Processor

 This continas a pacthed _Vagrant Post-Processsor_ that can build _Parallels Desktop Vagrant Boxes_. See the general documentation for [Packer Vagrant Post-Processor](http://www.packer.io/docs/post-processors/vagrant.html)

## Issues
If you find any bugs please open a issue at [github](https://github.com/rickard-von-essen/packer-parallels/issues). 

## Contributing
If you have any improvements open a pull request at [github](https://github.com/rickard-von-essen/packer-parallels/pulls). 

## License

This code is distributed under the MIT license, see _LICENSE_.

© _2014 Rickard von Essen, Yung Sang_

This work is derived from the _Packer VirtualBox builder plugin_ authored by _Mitchell Hashimoto et al._ For more information see [Packer](https://github.com/mitchellh/packer).
