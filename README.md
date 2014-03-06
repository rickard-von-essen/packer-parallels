# Packer Parallels Plugin

This is a custom builder plugin for [Packer](http://www.packer.io/) using [Parallels Desktop for Mac](http://www.parallels.com/products/desktop/).

**Current status: Alpha - BROKEN!**

![Parallels Logo](imgs/parallels_small.png)


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
Then download and install this plugin. This depends on [Parallels Virtualization SDK 9 for Mac](http://download.parallels.com//desktop/v9/pde.hf1/ParallelsVirtualizationSDK-9.0.24172.951362.dmg), download and install it first.

## Issues
If you find any bugs please open a issue at [github](https://github.com/rickard-von-essen/packer-parallels/issues). 

## Contributing
If you have any improvments open a pull request at [github](https://github.com/rickard-von-essen/packer-parallels/pulls). 

## License

This code is distributed under the MIT license, see _LICENSE_.

© 2014 Rickard von Essen