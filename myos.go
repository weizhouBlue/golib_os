package golib_os
import(
	"fmt"
	"os/exec"
	"os"
	"context"
	"bytes"
	"time"
	"strings"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"syscall"
)

/*
func SearchExecutable( name string ) (string , error) 

func RunCmd( cmdName string , env []string , stdin_msg string  , timeout_second int  ) (  stdoutMsg , stderrMsg  string , exitedCode int, e error ) {
func RunDaemonCmd( cmdName string , env []string , stdin_msg string ) (  stdoutMsg , stderrMsg  string , process *os.Process , e error ) {

func  ReadStdin() []byte {

func  ReadArgs() ( []string  ){

func PathExists(path string) (bool, error) {

func DirectoryExists(path string) (bool, error) {

func FileExists(path string) (bool, error) {

func FileSize(path string) ( int64  , error) {

func EmptyFile( sfilePath string ) error {
func WriteFile( sfilePath string , data []byte ) ( err error)   {
func ReadFile( sfilePath string  ) ( []byte , error)   {

func DeleteFile( sfilePath string ) error 
func DeleteDir( dirPath string ) error 

func UniqNumber() string {

func WriteJsonToFile( sfilePath string , jsonData interface{} ) ( err error)   {
func ReadJsonFromFile( sfilePath string  ) ( jsonData []byte , err error) {


func GetMyExecName() string {

func GetMyExecDir() string {
	
func GetMyRunDir() string {


*/

//-----------------------

var (
    EnableLog=false
)

func log( format string, a ...interface{} ) (n int, err error) {
    if EnableLog {
        return fmt.Printf(format , a... )    
    }
    return  0,nil
}



//-----------------------

func SearchExecutable( name string ) (string , error) {
	if len(name)==0 {
		return "" , fmt.Errorf("error, empty name")
	}

	if path, err:=exec.LookPath(name) ; err!=nil {
		return "" , err
	}else{
		return path , nil
	}

}

func RunCmd( cmdName string , env []string , stdin_msg string  , timeout_second int  ) (  stdoutMsg , stderrMsg  string , exitedCode int, e error ) {
	var outMsg bytes.Buffer
	var outErr bytes.Buffer
	var ctx context.Context
	var cancel context.CancelFunc
	var cmd *exec.Cmd
	var err error


	if len(cmdName)==0 {
		e=fmt.Errorf("error, empty cmd")
		return
	}

	log("run cmd=%s , env=%v , stdin_msg=%v \n" , cmdName , env , stdin_msg )

	if timeout_second==0 {
		timeout_second=10
	}

	ctx, cancel = context.WithTimeout( context.Background(), time.Duration(timeout_second) * time.Second )
	defer cancel()

	rootCmd:="bash"
	if path , _:=SearchExecutable(rootCmd) ; len(path)!=0 {
		cmd = exec.CommandContext( ctx,  rootCmd , "-c" ,  cmdName )
		goto EXE
	}


	rootCmd="sh"
	if path , _ :=SearchExecutable(rootCmd) ; len(path)!=0 {
		cmd = exec.CommandContext(ctx ,  rootCmd , "-c" , cmdName )
		goto EXE
	}

	e=fmt.Errorf("error, no sh or bash installed")
	return 

EXE:

	cmd.Env = append(os.Environ(), env... )

	if len(stdin_msg)!=0 {
		cmd.Stdin = strings.NewReader(stdin_msg)
	}

	cmd.Stdout=&outMsg
	cmd.Stderr=&outErr

	log("cmd=%v , timeout=%v \n", cmd , timeout_second )

	if err =cmd.Run() ; err!=nil {
		e=err
		return 
	}

	if a:=strings.TrimSpace( outMsg.String() ) ; len(a)>0{
		stdoutMsg=a			
	}
	if b:=strings.TrimSpace( outErr.String() ) ; len(b)>0{
		stderrMsg=b	
	}
	exitedCode=cmd.ProcessState.ExitCode()

	return 
}

func RunDaemonCmd( cmdName string , env []string , stdin_msg string ) ( process *os.Process , e error ) {
	var cmd *exec.Cmd
	var err error

	if len(cmdName)==0 {
		e=fmt.Errorf("error, empty cmd")
		return
	}

	log("run cmd=%s , env=%v , stdin_msg=%v \n" , cmdName , env , stdin_msg )

	rootCmd:="bash"
	if path , _:=SearchExecutable(rootCmd) ; len(path)!=0 {
		cmd = exec.Command(  rootCmd , "-c" ,  cmdName )
		goto EXE
	}
	rootCmd="sh"
	if path , _ :=SearchExecutable(rootCmd) ; len(path)!=0 {
		cmd = exec.Command(  rootCmd , "-c" , cmdName )
		goto EXE
	}

	e=fmt.Errorf("error, no sh or bash installed")
	return 

EXE:

	cmd.Env = append(os.Environ(), env... )
	cmd.SysProcAttr=&syscall.SysProcAttr{
		Setsid:  true ,
	}

	if len(stdin_msg)!=0 {
		cmd.Stdin = strings.NewReader(stdin_msg)
	}

	if err =cmd.Start() ; err!=nil {
		e=err
		return 
	}

	process=cmd.Process

	return 
}




//-----------------------


//读取标准输入 。 不堵塞
// EXE <<< "string...."
func  ReadStdin() []byte {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if fi.Size() > 0 {
		StdBytes , err:=ioutil.ReadAll(os.Stdin)
		if err!=nil{
			panic("failed to read std " )
		}
		return StdBytes
	} else {
		return nil
	}
}





func  ReadArgs() ( []string  ){
	return os.Args
}





//-----------------------

func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { 
        return true, nil //文件存在
    }
    if os.IsNotExist(err) { //判断错误类型是否是 不存在
        return false, nil
    }
    return false, err //其它原因
}


func DirectoryExists(path string) (bool, error) {
	if exist , e:= PathExists(path)  ; e!=nil || exist==false{
		return false , e
	}

    info , _ := os.Stat(path)
    if info.Mode().IsDir()==true {
    	return true, nil
    }
    return false , nil 

}

func FileExists(path string) (bool, error) {
	if exist , e:= PathExists(path)  ; e!=nil || exist==false{
		return false , e
	}

    info , _ := os.Stat(path)
    if info.Mode().IsRegular()==true {
    	return true, nil
    }
    return false , nil 

}


func FileSize(path string) ( int64  , error) {
	if exist , e:= FileExists(path)  ; e!=nil || exist==false{
		return 0 , e
	}

    info , _ := os.Stat(path)
    return info.Size() , nil 

}




//  creates or truncates the named file
func EmptyFile( sfilePath string ) error {

	// https://godoc.org/os#Create
    if file , err:= os.Create(sfilePath) ; err!=nil {
        return  err
    }else{
        file.Close()
        return nil
    }

}

func UniqNumber() string {

	m :=time.Now()
	return fmt.Sprintf( "%04d%02d%02d%02d%02d%02d%09d" , m.Year() , int(m.Month()) , m.Day() , m.Hour() , m.Minute() , m.Second() , m.Nanosecond()  )
	
}

//-----------------------


func ReadJsonFromFile( sfilePath string  ) ( jsonData []byte , err error) {

	if a, e := FileExists(sfilePath) ; a==false || e!=nil {
		err=fmt.Errorf("no file %v" , sfilePath )
		return
	}

	// https://godoc.org/io/ioutil#ReadAll
	jsonData, err =ioutil.ReadFile(sfilePath)
	if err!=nil {
		return 
	}

	// https://godoc.org/encoding/json#Valid
	if json.Valid( jsonData  )==false {
		err=fmt.Errorf("data is not json format in file %v" , sfilePath )
		return
	}

	return

}

// 覆盖写
func WriteJsonToFile( sfilePath string , jsonData interface{} ) ( err error)   {

	// https://godoc.org/encoding/json#Marshal
	jsonByte , e := json.Marshal( jsonData ) 
	if err != nil {
		err=e
		return
	}

    // https://godoc.org/io/ioutil#WriteFile
    // 覆盖写
    err = ioutil.WriteFile( sfilePath , jsonByte , 0644 ) 
    return 

}



func GetMyExecName() string {
	return filepath.Base(os.Args[0])
}

// 获取 运行的当前的命令所在的目录
// 例如: 在 /a 下执行 /usr/bin/b ， 那么 输出 /usr/bin/b
func GetMyExecDir() string {
	//返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	//将\替换成/
	return strings.Replace(dir, "\\", "/", -1) 


	// path, err := os.Executable()
	// if err != nil {
	//     return ""
	// }
	// dir := filepath.Dir(path)
	// return dir

}

// 获取运行当前命令时 所在的目录
// 例如: 在 /a 下执行 /usr/bin/b ， 那么 输出 /a
func GetMyRunDir() string {

	if dir , e := os.Getwd() ; e!=nil{
		return ""
	}else{
		return dir
	}
}


// 删除一个文件或者一个空目录
func DeleteFile( sfilePath string ) error  {
	// https://godoc.org/os#Remove
    return  os.Remove(sfilePath)
}

// 删除任意，包括目录下的所有东西
func DeleteDirOrFile( dirPath string ) error {
	//https://godoc.org/os#RemoveAll
    return  os.RemoveAll(dirPath)
}


// 覆盖写
func WriteFile( sfilePath string , data []byte ) ( err error)   {

    // https://godoc.org/io/ioutil#WriteFile
    // 覆盖写
    return ioutil.WriteFile( sfilePath , data , 0644 )  

}

func ReadFile( sfilePath string  ) ( []byte , error)   {

	// https://godoc.org/io/ioutil#ReadFile
    return ioutil.ReadFile( sfilePath  )  

}


