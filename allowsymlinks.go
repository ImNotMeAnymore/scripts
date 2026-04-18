//!go run
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func allowsym(path string)(bool,error){
	_,err:=os.Stat(path+"/allowed_symlinks.txt")
	if os.IsNotExist(err) {
		f,err := os.Create(path+"/allowed_symlinks.txt")
		if err != nil {return false,err}
		defer f.Close()
		fmt.Fprintln(f, "[regex].*")
		return true,nil
	} else if err != nil {return false, err}
	return false, nil
}

func find(root string, name string)(string,error){
	entries,err := os.ReadDir(root)//reads the contents of root
	if err!=nil{return "Failed reading file",err}//if it failed propagates the error
	for _,entry:=range entries {//for each file on root
		if entry.Name()==name{return root,nil}//Returns it root
		if !entry.IsDir(){continue}//if it's a regular file ignore it
		chn,err:=find(root+"/"+entry.Name(), name)
		if err==nil{return chn,nil}//if it did find it
		if err!=os.ErrNotExist{return "",err}//any other error
	}
	return "File not found",os.ErrNotExist
}

func lastidtxt(path string)(int,error){
	data,err := os.ReadFile(path)
	if err!=nil{return -1,err}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}

func lastidjson(path string)(int,error){
	//here the number is on a json that looks like this: {
	//  "computer": 2
	//}
	data,err := os.ReadFile(path)
	if err!=nil{return -1,err}
	var last struct {//to unmarshall json you need to create a struct
		ID int `json:"computer"`//then each key is defined with a struct tag
		// the struct tag lets know the struct what to do with what values
		// fields starting with uppercase are exported, and can be accessed by other packages like json
		// so the unmarshall struct needs exported fields
	}
	err=json.Unmarshal(data,&last)//unmarshall takes a pointer to the variable
	if err!=nil{return -1,err}
	return last.ID,nil
}


func addShared(path string, n int)([]error,error){
	//path is something like computercraft/computer/ with folders enumerated from 0 to n
	// the folders don't nesesarly exist yet
	//iterates from 0 to n inclusive, if the folder doesn't exist creates it
	// then it tries to make a link (hard probably)
	// /home/yo/Projects/LUA/CC/shared -> path/{n}/shared
	// until the end
	// if it managed to do so in at leastt one returns
	errs := make([]error,n+1)//make with 2 args creates a slice of len n
	src := "/home/yo/Projects/LUA/CC/shared"
	num:=0
	println(path)
	for i:=0;i<=n;i++{
		com:=fmt.Sprintf("%s/%d/",path,i)//the computer folder
		//this is path+"/"+str(i)+"/"
		err:=os.MkdirAll(com,0o755)//makes a dir with some perms idk
		if err!=nil{goto post}//if error goto the end
		err=os.Symlink(src,com+"/shared")//tries to link shared to the computer folder
		if err!=nil{goto post}//if error goto the end
		num++
		post:
		errs[i]=err//this is nil if it's alright and the error otherwise
	}
	rc := "/home/yo/.local/share/LuaLS/cc-definitions/.luarc.json"
	if _,err:=os.Stat(rc); err==nil { // if it statd with no problems
		return errs,os.ErrExist
	} else if os.IsNotExist(err) {
		dat,err := os.Open(src)
		if err!=nil{return errs,err}
		defer dat.Close()
		f, err := os.Create(rc)
		if err!=nil{return errs,err}
		defer f.Close()
		_,err = io.Copy(f, dat)
		return errs,err //nil if it was created no problem
	} else {
		return errs,err
	}
}


func CCshareds(path string)([]error,error){
	root,err:=find(path,"lastid.txt")//1.12.2 and below
	if err==nil{
		n,err:=lastidtxt(root+"/lastid.txt")
		if err==nil{return addShared(root,n)}
	}
	root,err=find(path,"ids.json")//1.16 and above
	if err==nil{
		n,err:=lastidjson(root+"/ids.json")
		if err==nil{return addShared(root,n)}
	}
	return []error{},os.ErrNotExist
}


func main() {
	//file,err := os.Create("/home/yo/go.log")
	//if err!=nil{fmt.Println("Error opening log file:", err);return}
	//defer file.Close()
	//fmt.Fprintln(file, "test")
	if len(os.Args)<=1{fmt.Println("Not enough arguments, needs an instance dir");return}
	root,err:=find(os.Args[1], "saves")
	if err!=nil{fmt.Println("Error finding instance root:", err);return}
	did,err:=allowsym(root)
	if err!=nil{fmt.Println("Error creating allowed_symlinks.txt:", err);return}
	fmt.Print("allowed_symlinks.txt ")
	if !did{fmt.Print("NOT ")}
	fmt.Println("Created!")
	sh,err := CCshareds(root+"/saves")
	
	if len(sh)==0{
		fmt.Println("No computers found!")
	}else{
		for n,e := range sh {
			if e==nil{
				fmt.Printf("[%d] shared Created!\n",n)
			}else{
				fmt.Printf("[%d] shared not Created: %s\n",n,e.Error())
			}
		}
		if os.IsExist(err){
			//fmt.Println(".luarc.json already exists!")
		} else if err == nil {
			fmt.Println(".luarc.json Created!")
		}
	}
}
