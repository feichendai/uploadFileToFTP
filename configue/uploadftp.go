package configue

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"upload/upload/saveJson"
)
var fs []string
func getListDir(dirPth string) ([]string, error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := GetOSSep()
	for _, fi := range dir {
		if fi.IsDir() {
			//files1 = append(files1, dirPth+PthSep+fi.Name())
			getListDir(dirPth +PthSep+fi.Name())
		}else{
			fs = append(fs, dirPth+PthSep+fi.Name())
		}
	}
	return fs, nil
}
func GetAllFileName(path string) (int, []string ) {
	configPath := GetConfigPath()
	ftpConfig := new(Config)
	confPath:=configPath + "upload"+ GetOSSep()+"configue"+ GetOSSep()+"conf.ini"
	ftpConfig.InitConfig(confPath)
	files, err := getListDir(path)
	if err != nil {
		fmt.Println( "System","get file path err")
	}
	fileLen := len(files)
	return fileLen, files
}
func getLocalIpAddr() string {
	configPath := GetConfigPath()
	ftpConfig := new(Config)
	ftpConfig.InitConfig(configPath + "conf.ini")
	network := ftpConfig.Read("ftp", "comm_way")
	ip := ftpConfig.Read("ftp", "local_ip")
	port := ftpConfig.Read("ftp", "local_port")
	address := ip + ":" + port
	conn, err := net.Dial(network, address)
	if err != nil {
		 return "127.0.0.1"
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}
func ftpUploadFile(ftpserver, ftpuser, pw, localFile, remoteSavePath, saveName string)error {
	configPath := GetConfigPath()
	ftpConfig := new(Config)
	ftpConfig.InitConfig(configPath + "upload"+ GetOSSep()+"configue"+ GetOSSep()+"conf.ini")
	//ftpfile_path := ftpConfig.Read("ftp", "ftpfile_path")
	Rftp, err := ftp.Connect(ftpserver)
	if err !=nil {
		//panic("connect ftp fail")
		fmt.Println( "System", "connect err")
		return err
	}
	err = Rftp.Login(ftpuser, pw)
	if err != nil {
		fmt.Println( "System", "Login err")
		return err
	}
	err =Rftp.ChangeDir(remoteSavePath)
	//dir, err := Rftp.CurrentDir()
	//if err !=nil{
	//	fmt.Println( "System", "Get currentDir err ",dir)
	//	return err
	//}
	//fmt.Println("ftp path: ",dir)
	//fmt.Println("remoter path: ", remoteSavePath)
	if err !=nil {
		err=Rftp.MakeDir(remoteSavePath)
		if err !=nil{
			fmt.Println( "System", "MakeDir1 err:",err)
			Rftp.Logout()
			Rftp.Quit()
			return err
		}
		err=Rftp.ChangeDir(remoteSavePath)
		if err !=nil{
			fmt.Println( "System", "ChangeDir err:",err)
			//fmt.Println( "File path:",saveName)
			Rftp.Logout()
			Rftp.Quit()
			return err
		}
	}
	_, err  = Rftp.CurrentDir()
	if err !=nil{
		fmt.Println( "System", "MakeDir2 err")
		return err
	}
	//fmt.Println( "System", dir)
	file, err := os.Open(localFile)
	if err != nil {
		fmt.Println( "System", "Open err")
		Rftp.Logout()
		Rftp.Quit()
		return err
	}
	defer file.Close()
	err = Rftp.Stor(saveName, file)
	if err != nil {
		fmt.Println( "System", "Stor err")
		Rftp.Logout()
		Rftp.Quit()
		return err
	}
	Rftp.Logout()
	Rftp.Quit()
	return nil
	//logcotent := fmt.Sprintf("%s:%s","success upload file",localFile)
	//fmt.Println( "System", logcotent)
}
func RemoveFile(filePath string, fileName string){
	configPath := GetConfigPath()
	ftpConfig := new(Config)
	ftpConfig.InitConfig(configPath + "upload"+ GetOSSep()+"configue"+ GetOSSep()+"conf.ini")
	err := os.Remove(filePath + fileName)
	if err != nil {
		fmt.Println( "System", "file remove err!")
	} else {
		logcotent := fmt.Sprintf("%s:%s","file remove OK!",fileName)
		fmt.Println( "System", logcotent)
	}
}
func SendAllFileToFtpServer() {
	configPath := GetConfigPath()
	ftpConfig := new(Config)
	ftpConfig.InitConfig(configPath + "upload"+ GetOSSep()+"configue"+ GetOSSep()+"conf.ini")
	local_path := ftpConfig.Read("path", "local_path")
	filePath1:=filepath.ToSlash(local_path)
	flen, fileName := GetAllFileName(filePath1 )
	//fmt.Println("Path file:",fileName)
	ftpserverip := ftpConfig.Read("ftp", "ftp_server_ip")
	ftpPort := ftpConfig.Read("ftp", "ftp_server_port")
	ftpuser := ftpConfig.Read("ftp", "ftp_server_name")
	pw := ftpConfig.Read("ftp", "ftp_server_pwd")

	is_removeFile:=ftpConfig.Read("path", "remove_file")
	ftpserver := ftpserverip + ":" + ftpPort
	files,err:=saveJson.GetJsonInfo()
	//fmt.Println("Path Json:",files)
	if err !=nil{
		fmt.Println("saveJson.GetJsonInfo err:",err)
		return
	}
	for i := 0; i < flen; i++{
		var is_exist bool
		for _, item := range files {
			if item == fileName[i] {
				is_exist=true
			}
		}
		if is_exist{
			continue
		}
		fp:=filepath.ToSlash(fileName[i])
		logcotent,domainFile := filepath.Split(fp)

		domainName := logcotent[len(local_path):len(logcotent)]
		//domainFile:=fp[len(local_path):len(fp)]
		//fmt.Println( "path:", domainName,"file:",domainFile)
		err:=ftpUploadFile(ftpserver, ftpuser, pw, fp, domainName, domainFile)
		if err !=nil {
			fmt.Println("ftp err:",err)
			continue
		}else{
			files= append(files,fileName[i] )
		}
		if is_removeFile=="Yes" {
			RemoveFile(filePath1, fileName[i])
		}



		}
	saveJson.SaveToJson(files)
}
