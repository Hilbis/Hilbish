//go:build darwin && !unix

package main

// String vars that are free to be changed at compile time
var (
	requirePaths = commonRequirePaths + `.. ';'
	.. hilbish.dataDir .. '/libs/?/init.lua;'
	.. hilbish.dataDir .. '/libs/?/?.lua;'` + macosUserPaths
	macosUserPaths = `
	.. hilbish.userDir.data     .. '/hilbish/libs/?/init.lua;'
	.. hilbish.userDir.data	    .. '/hilbish/libs/?/?.lua;'
	.. hilbish.userDir.data	    .. '/hilbish/libs/?.lua;'
	.. hilbish.userDir.config	.. '/hilbish/?/init.lua;'
	.. hilbish.userDir.config	.. '/hilbish/?/?.lua;'
	.. hilbish.userDir.config	.. '/hilbish/?.lua'`
	dataDir = "/usr/local/share/hilbish"
	preloadPath = dataDir + "/nature/init.lua"
	sampleConfPath = dataDir + "/.hilbishrc.lua" // Path to default/sample config
	defaultConfDir = getenv("XDG_CONFIG_HOME", "~/.config")
)
