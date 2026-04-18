#!/usr/bin/env python3

from sys import argv
import os
import subprocess

_=argv.pop(0)

if not argv:exit(not print("exactly 1 argument expected"))

shll = os.getenv("SHELL",".")

exts={
	"py":"python3.14 '{file}'",
	"user.js":"librewolf '{file}'",
	"mp4":"mpv '{file}'",
	"sh":shll+" '{file}'",
	"c":"gcc '{file}' -o '{file}.out'; '{file}.out'; echo '\n'",
	"html":"librewolf '{file}'",
	"go":"go run '{file}'",
	"hs":"runhaskell '{file}'",
	"lua":"luajit '{file}'",
	"java":"java '{file}'",
	"rkt":"racket '{file}'",
}

def runner(path:str):
	p = os.path.abspath(path)
	with open(p,"r")as f:
		for i in f:
			if not i.strip(): continue
			if(sh:=i.split("!",1))[0]in("#","//","// ")and(r:=sh[-1].strip()):
				return f"{r} '{p}'"
			break
	ext = "sh"
	_,nm = path.rsplit("/",1)
	if"."in path:nm,ext=nm.split(".",1)
	return exts.get(ext,shll).format(file=path,name=nm)

fil=argv.pop(0)
aa=subprocess.call(runner(fil), shell=True)
print() #because zed overwrites the last line sometimes
