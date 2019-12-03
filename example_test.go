package golib_os_test
import (
	"testing"
	myos "github.com/weizhouBlue/golib_os"
	"fmt"
	"time"
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
	cmd:="echo $WELAN"
	// addtional environment 
	env:=[]string{
		"WELAN=12345",
		"TOM=uit",
	}
	// stdin for cmd
	stdin_msg:="this is stdin msg"
	// o for no auto timeout
	timeout_second:=5
	chanFinish , _ , chanStdoutMsg , chanStderrMsg , chanErr , e  :=myos.RunCmd( cmd, env , stdin_msg , timeout_second )
	if e!=nil {
		fmt.Println(  "failed to exec "+ cmd )
		t.FailNow()

	}
	<-chanFinish

	if data , ok := <-chanErr ; ok {
		// return code with no succeed
		fmt.Println(  "err : "+ data )
	}else{
		//fmt.Println("ok for cmd"  )

		if data , ok := <-chanStdoutMsg ; ok {
			fmt.Println(  "stdoutMsg: "+ data )
		}
		if data , ok := <-chanStderrMsg ; ok {
			fmt.Println(  "stderrMsg: "+ data )
		}
	}

}

func Test_simple_longCmd(t *testing.T){

	myos.EnableLog=true

	// exec command
	cmd:="sleep 10d"
	// addtional environment 
	env:=[]string{
		"WELAN=12345",
		"TOM=uit",
	}
	// stdin for cmd
	stdin_msg:="this is stdin msg"
	// o for no auto timeout
	timeout_second:=0
	chanFinish , chanCancel , chanStdoutMsg , chanStderrMsg , chanErr , e  :=myos.RunCmd( cmd, env , stdin_msg , timeout_second )
	if e!=nil {
		fmt.Println(  "failed to exec "+ cmd )
		t.FailNow()

	}

	//wait for cmd
	select{
	case <- chanFinish:
		fmt.Println("cmd finish ")
	case <- time.After(5*time.Second) : 
		fmt.Println("cmd timeout , cancel it")
		close(chanCancel)
	}
	
	//read msg
	if data , ok := <-chanErr ; ok {
		// return code with no succeed
		fmt.Println(  "err : "+ data )
	}else{
		//fmt.Println("ok for cmd"  )

		if data , ok := <-chanStdoutMsg ; ok {
			fmt.Println(  "stdoutMsg: "+ data )
		}
		if data , ok := <-chanStderrMsg ; ok {
			fmt.Println(  "stderrMsg: "+ data )
		}
	}

}

