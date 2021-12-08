package saveJson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//获取当地uploaded.json中的信息并解析
func  GetJsonInfo() ([]string,error) {
	bytes, err := ioutil.ReadFile("./upload/saveJson/uploaded.json")
	//bytes, err := ioutil.ReadFile("uploaded.json")
	 if err !=nil{
		fmt.Println("ioutil.ReadFile err:",err)
		return nil,err
	}
	var files []string
	mPath:=make(map[string][]string)
	_=json.Unmarshal(bytes, &mPath)
	files=mPath["path"]
	return files,nil
}
//
func  SaveToJson(files []string)error  {
	mPath:=make(map[string][]string)
	mPath["path"]=files
	saveJson,_:=json.Marshal( &mPath)
	f,err:=os.OpenFile("upload/saveJson/uploaded.json",os.O_WRONLY|os.O_CREATE,0766)
	if err !=nil{
		fmt.Println("open file err:",err)
		return err
	}
	defer f.Close()
	_,err =f.Write(saveJson)
	if err !=nil{
		return err
	}
	return nil
}
