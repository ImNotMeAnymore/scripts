# scripts
Small programs I want in all my computers

# TODO
Installer and sorter



## runner.py
so that zed runs programs on the terminal instead of the nasty debugger it has

1) move to bin as runner or something like that
2) f1 > zed: open tasks
3) paste this into it:
```json
{
		"label": "Runner",
		"command": "runner '$ZED_FILE'; $SHELL",
		"hide": "always",
}
```


## allowsymlinks.go
adds allowsymlinks.txt to minecraft instances + manages computercraft files


1) go build it anywhere, it doesn't really matter
2) idk I'll add the rest later
