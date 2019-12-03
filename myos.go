package golib_os
import(
	"fmt"
	"os/exec"
	"os"
	"context"
	"bytes"
	"time"
	"strings"
)



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

func RunCmd( cmdName string , env []string , stdin_msg string  , timeout_second int ) ( chanFinish , chanCancel chan bool , chanStdoutMsg , chanStderrMsg , chanErr chan string , e error ) {
	var outMsg bytes.Buffer
	var outErr bytes.Buffer
	var ctx context.Context
	var cancel context.CancelFunc
	var cmd *exec.Cmd
	var err error

	chanCancel = make ( chan bool )
	chanFinish = make ( chan bool )
	chanStdoutMsg = make ( chan string ,1 )
	chanStderrMsg = make ( chan string , 1)
	chanErr = make ( chan string ,1 )


	if len(cmdName)==0 {
		return nil , nil , nil , nil ,nil , fmt.Errorf("error, empty cmd")
	}

	log("run cmd=%s , env=%v , stdin_msg=%v \n" , cmdName , env , stdin_msg )

	if timeout_second>0 {
		ctx, cancel = context.WithTimeout( context.Background(), time.Duration(timeout_second) * time.Second )
	}else{
		ctx, cancel = context.WithCancel(context.Background()  )
	}


	rootCmd:="bash"
	if path , _:=SearchExecutable(rootCmd) ; len(path)!=0 {
		cmd = exec.CommandContext( ctx,  rootCmd , "-c" , cmdName )
		goto EXE
	}

	rootCmd="sh"
	if path , _ :=SearchExecutable(rootCmd) ; len(path)!=0 {
		cmd = exec.CommandContext(ctx ,  rootCmd , "-c" , cmdName )
		goto EXE
	}

	return nil , nil , nil , nil ,nil , fmt.Errorf("error, no sh or bash installed")
	

EXE:
	if env!=nil || len(env)!=0 {
		cmd.Env = append(os.Environ(), env... )
	}
	if len(stdin_msg)!=0 {
		cmd.Stdin = strings.NewReader(stdin_msg)
	}


	cmd.Stdout=&outMsg
	cmd.Stderr=&outErr


	log("cmd=%v \n", cmd)
	if err =cmd.Start() ; err!=nil {
		return nil , nil , nil , nil ,nil , err
	}


	go func(){
		log("routine for closing cmd=%v \n" , cmdName )
		<- chanCancel
		cancel()
		log("closing routine eixt for cmd=%v \n" , cmdName )
	}()

	go func(){
		log("routine for executing cmd=%v \n" , cmdName )
		err:=cmd.Wait()
		log("routine ending cmd=%v \n" , cmdName )

		if err!=nil {
			chanErr<-fmt.Sprintf("%v" , err )
		}
		chanStdoutMsg<- outMsg.String()
		chanStderrMsg<- outErr.String()

		close(chanStdoutMsg)
		close(chanStderrMsg)
		close(chanErr)

		select{
		case  <-chanCancel:
			//channel has been closed
		default:
			close(chanCancel)
		}
		close(chanFinish)
		log("executing routine exit for cmd=%v \n" , cmdName )

	}()


	return chanFinish , chanCancel , chanStdoutMsg , chanStderrMsg , chanErr , nil
}








