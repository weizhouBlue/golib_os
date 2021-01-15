package golib_os_test
import (
	"testing"
	myos "github.com/weizhouBlue/golib_os"
	"fmt"
	"time"
	"os"
)

//====================================

func Test_seatch(t *testing.T){

	myos.EnableLog=true


	executable:="ls"
	if path , err:=myos.SearchExecutable( executable ); err!= nil {
		fmt.Println(  "failed to find "+ executable )
		t.FailNow()
	}else{
		fmt.Println( executable + " is @ " + path )

	}


}


func Test_simple_simpleCmd(t *testing.T){

	myos.EnableLog=true

	// exec command
	// 命令其实是直接更 bash -c "..."
	// 命令不能是后台的!不能产生 一直运行的子进程的，否则会一直等待，不能自动超时 
	// 不要跑组合命令  sleep 50 ; echo "hello " , 否则会一直等待直到完成，不能自动超时
	cmd:="echo $WELAN ; echo aaa "
	// addtional environment 
	env:=[]string{
		"WELAN=12345",
		"TOM=uit",
	}
	// stdin for cmd
	stdin_msg:="this is stdin msg"
	// o for no auto timeout
	timeout_second:=5
	if StdoutMsg , StderrMsg  ,exitedCode , e  :=myos.RunCmd( cmd, env , stdin_msg , timeout_second  ); e!=nil {
		fmt.Printf(  "failed to exec %v : %v", cmd , e )
		t.FailNow()

	}else{
		if exitedCode!=0 {
			fmt.Printf(  "error, exitedCode : %v \n" , exitedCode )
		}
		fmt.Println(  "stderrMsg : "+ StderrMsg )
		fmt.Println(  "StdoutMsg : "+ StdoutMsg )
	}

}

func Test_simple_longCmd(t *testing.T){

	myos.EnableLog=true

	// exec command
	// 命令其实是直接更 bash -c "..."
	// 命令不能是后台的!不能产生 一直运行的子进程的，否则会一直等待直到完成，不能自动超时 
	// 不要跑组合命令  sleep 50 && echo "hello " , 否则会一直等待直到完成，不能自动超时
	cmd:=" sleep 50  "  
	// addtional environment 
	env:=[]string{
		"WELAN=12345",
		"TOM=uit",
	}
	// stdin for cmd
	stdin_msg:="this is stdin msg"
	// o for no auto timeout
	timeout_second:=5
	if  StdoutMsg , StderrMsg  ,exitedCode , e  :=myos.RunCmd( cmd, env , stdin_msg , timeout_second ); e!=nil {
		fmt.Printf(  "failed to exec %v : %v", cmd , e )
		t.FailNow()

	}else{
		if exitedCode!=0 {
			fmt.Printf(  "error, exitedCode : %v \n" , exitedCode )
		}
		fmt.Println(  "stderrMsg : "+ StderrMsg )
		fmt.Println(  "StdoutMsg : "+ StdoutMsg )
	}

}



func Test_back(t *testing.T){

	myos.EnableLog=true

	// exec command
	// 命令其实是直接更 bash -c "..."
	// 命令不能是后台的!不能产生 一直运行的子进程的，否则会一直等待，不能自动超时 
	cmd:=` sleep 20 && echo byebye ` 
	// addtional environment 
	env:=[]string{
		"WELAN=12345",
		"TOM=uit",
	}
	// stdin for cmd
	stdin_msg:="this is stdin msg"
	if  process , e  :=myos.RunDaemonCmd( cmd, env , stdin_msg  ); e!=nil {
		fmt.Printf(  "failed to exec %v : %v", cmd , e )
		t.FailNow()

	}else{


		// https://godoc.org/os#Process
		fmt.Printf(  "pid : %v \n" , process.Pid )
		
		time.Sleep(5*time.Second)
		//process.Kill()

		// https://godoc.org/os#Process
		// if _ , e:=process.Wait(); e!=nil {
		// 	fmt.Printf(  "error : %v \n" , e )
		// }

		
	}

}







func Test_json1(t *testing.T){

	data:=map[string] string {
		"k1": "v1" ,
		"k2": "v2" ,

	}
	filePath:="./json_test"

	if e:=myos.WriteJsonToFile( filePath , data ) ; e!=nil {
		fmt.Println(  "failed to WriteJsonToFile " )
		t.FailNow()
	}


	if jsondata , e:=myos.ReadJsonFromFile( filePath ) ; e!=nil {
		fmt.Println(  "failed to ReadJsonFromFile " )
		t.FailNow()
	}else{
		fmt.Printf(  "json data: %v \n" , string(jsondata) )
	}

	os.Remove(filePath)


	
	myos.EmptyFile(filePath)

	size , _ := myos.FileSize( filePath )
	fmt.Printf("size1eddd : %v \n "  , size  ) 




}


func Test_path(t *testing.T){
	fmt.Println( myos.GetMyExecName() )
	fmt.Println( myos.GetMyExecDir() )
	fmt.Println( myos.GetMyRunDir() )
	
	
}


func Test_write( t *testing.T ) {

	data:=[]byte("line1 \n line2 \n\n")
	fmt.Println( myos.WriteFile( "./test"  , data  )  )


}






