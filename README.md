heka-encoder-filter
===================

A filter for [Mozilla Heka](http://hekad.readthedocs.org/) that does nothing but encode the payload.

EncoderFilter
===========


To Build
========

See [Building *hekad* with External Plugins](http://hekad.readthedocs.org/en/latest/installing.html#build-include-externals)
for compiling in plugins.

Edit cmake/plugin_loader.cmake file and add

    add_external_plugin(git https://github.com/michaelgibson/heka-encoder-filter master)

Build Heka:
	. ./build.sh


Config
======
[encoder_filter_zlib]
type = "EncoderFilter"
encoder = "zlib_encoder"
message_matcher = "Fields[StreamAggregatorTag] == 'aggregated'"

[zlib_encoder]
type = "ZlibEncoder"
zlib_tag = "compressed"
