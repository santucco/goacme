

/*3:*/


//line goacme.w:97

package goacme

import(
"os"
"os/exec"
"testing"
"syscall"


/*11:*/


//line goacme.w:185

"fmt"



/*:11*/



/*16:*/


//line goacme.w:238

"bytes"
"errors"



/*:16*/


//line goacme.w:105

)

func prepare(t*testing.T){
// checking for a running plumber instance
f,err:=os.Open(AcmeDir+"/index")
if err==nil{
f.Close()
t.Log("acme started already")
}else{
cmd:=exec.Command("acme","-m",AcmeDir)
err= cmd.Start()
if err!=nil{
t.Fatal(err)
}
}
}



/*13:*/


//line goacme.w:196

func TestNewOpen(t*testing.T){
prepare(t)


/*12:*/


//line goacme.w:192

syscall.Nanosleep(&syscall.Timespec{Sec:1},nil)



/*:12*/


//line goacme.w:199

w,err:=New()
if err!=nil{
t.Fatal(err)
}
defer w.Close()
defer w.Del(true)
if f,err:=os.Open(fmt.Sprintf("%s/%d",AcmeDir,w.id));err!=nil{
t.Fatal(err)
}else{
f.Close()
}
}



/*:13*/



/*17:*/


//line goacme.w:243

func TestReadWrite(t*testing.T){


/*12:*/


//line goacme.w:192

syscall.Nanosleep(&syscall.Timespec{Sec:1},nil)



/*:12*/


//line goacme.w:245

w,err:=New()
if err!=nil{
t.Fatal(err)
}
defer w.Close()
defer w.Del(true)
b1:=[]byte("test")
_,err= w.Write(b1)
if err!=nil{
t.Fatal(err)
}
w1,err:=Open(w.id)
if err!=nil{
t.Fatal(err)
}
defer w1.Close()
defer w1.Del(true)
b2:=make([]byte,10)
n,err:=w1.Read(b2)
if err!=nil{
t.Fatal(err)
}
if bytes.Compare(b1,b2[:n])!=0{
t.Fatal(errors.New("buffers don't match"))
}
}



/*:17*/



/*26:*/


//line goacme.w:354

func TestPipeTo(t*testing.T){


/*12:*/


//line goacme.w:192

syscall.Nanosleep(&syscall.Timespec{Sec:1},nil)



/*:12*/


//line goacme.w:356

w,err:=New()
if err!=nil{
t.Fatal(err)
}
defer w.Close()
defer w.Del(true)
s:="test"
p,err:=w.PipeTo("body",nil,"echo","-n",s)
if err!=nil{
t.Fatal(err)
}
p.Wait()
p.Release()
w1,err:=Open(w.id)
if err!=nil{
t.Fatal(err)
}
defer w1.Close()
defer w1.Del(true)
b:=make([]byte,10)
n,err:=w1.Read(b)
if err!=nil{
t.Fatal(err)
}
if bytes.Compare([]byte(s),b[:n])!=0{
t.Fatal(errors.New(fmt.Sprintf("buffers don't match: %q and %q",s,string(b))))
}
}



/*:26*/



/*28:*/


//line goacme.w:414

func TestPipeFrom(t*testing.T){


/*12:*/


//line goacme.w:192

syscall.Nanosleep(&syscall.Timespec{Sec:1},nil)



/*:12*/


//line goacme.w:416

w,err:=New()
if err!=nil{
t.Fatal(err)
}
s:="test"
if _,err:=w.Write([]byte(s));err!=nil{
t.Fatal(err)
}
if _,err:=w.Seek(0,0);err!=nil{
t.Fatal(err)
}
f,err:=os.OpenFile("/tmp/goacme.test",os.O_RDWR|os.O_TRUNC|os.O_CREATE,0600)
if err!=nil{
t.Fatal(err)
}
defer f.Close()


/*12:*/


//line goacme.w:192

syscall.Nanosleep(&syscall.Timespec{Sec:1},nil)



/*:12*/


//line goacme.w:433

p,err:=w.PipeFrom("body",f,"cat")
if err!=nil{
t.Fatal(err)
}


/*12:*/


//line goacme.w:192

syscall.Nanosleep(&syscall.Timespec{Sec:1},nil)



/*:12*/


//line goacme.w:438

w.Del(true)
w.Close()
p.Wait()
p.Release()
if _,err:=f.Seek(0,0);err!=nil{
t.Fatal(err)
}
b:=make([]byte,10)
n,err:=f.Read(b)
if err!=nil{
t.Fatal(err)
}
if bytes.Compare([]byte(s),b[:n])!=0{
t.Fatal(errors.New(fmt.Sprintf("buffers don't match: %q and %q",s,string(b))))
}
}



/*:28*/



/*30:*/


//line goacme.w:484

func TestSysRun(t*testing.T){


/*12:*/


//line goacme.w:192

syscall.Nanosleep(&syscall.Timespec{Sec:1},nil)



/*:12*/


//line goacme.w:486

s:="test"
f,p,err:=SysRun("echo","-n",s)
if err!=nil{
t.Fatal(err)
}
p.Wait()
p.Release()


/*12:*/


//line goacme.w:192

syscall.Nanosleep(&syscall.Timespec{Sec:1},nil)



/*:12*/


//line goacme.w:494

b:=make([]byte,10)
if _,err:=f.Seek(0,0);err!=nil{
t.Fatal(err)
}
n,err:=f.Read(b)
if err!=nil{
t.Fatal(err)
}
if bytes.Compare([]byte(s),b[:n])!=0{
t.Fatal(errors.New(fmt.Sprintf("buffers don't match: %q and %q",s,string(b))))
}
}




/*:30*/



/*32:*/


//line goacme.w:526

func TestDel(t*testing.T){


/*12:*/


//line goacme.w:192

syscall.Nanosleep(&syscall.Timespec{Sec:1},nil)



/*:12*/


//line goacme.w:528

w,err:=New()
if err!=nil{
t.Fatal(err)
}


/*12:*/


//line goacme.w:192

syscall.Nanosleep(&syscall.Timespec{Sec:1},nil)



/*:12*/


//line goacme.w:533

w.Del(true)
w.Close()


/*12:*/


//line goacme.w:192

syscall.Nanosleep(&syscall.Timespec{Sec:1},nil)



/*:12*/


//line goacme.w:536

if _,err:=Open(w.id);err==nil{
t.Fatal(errors.New(fmt.Sprintf("window %d is still opened",w.id)))
}
}




/*:32*/



/*39:*/


//line goacme.w:595

func TestDeleteAll(t*testing.T){


/*12:*/


//line goacme.w:192

syscall.Nanosleep(&syscall.Timespec{Sec:1},nil)



/*:12*/


//line goacme.w:597

var l[10]int
for i:=0;i<len(l);i++{


/*12:*/


//line goacme.w:192

syscall.Nanosleep(&syscall.Timespec{Sec:1},nil)



/*:12*/


//line goacme.w:600

w,err:=New()
if err!=nil{
t.Fatal(err)
}
l[i]= w.id
}
DeleteAll()


/*12:*/


//line goacme.w:192

syscall.Nanosleep(&syscall.Timespec{Sec:1},nil)



/*:12*/


//line goacme.w:608

for _,v:=range l{
_,err:=Open(v)
if err==nil{
t.Fatal(errors.New(fmt.Sprintf("window %d is still opened",v)))
}
}
}



/*:39*/



/*65:*/


//line goacme.w:905

func TestEvent(t*testing.T){


/*12:*/


//line goacme.w:192

syscall.Nanosleep(&syscall.Timespec{Sec:1},nil)



/*:12*/


//line goacme.w:907

w,err:=New()
if err!=nil{
t.Fatal(err)
}
defer w.Close()
defer w.Del(true)
msg:="Press left button of mouse on "
test:="Test"
if _,err:=w.Write([]byte(msg+test));err!=nil{
t.Fatal(err)
}
ch,err:=w.EventChannel()
if err!=nil{
t.Fatal(err)
}
ok:=false
var e*Event
for e,ok= <-ch;ok;e,ok= <-ch{
if e.Origin==Mouse{
break
}
}
if!ok{
t.Fatal(errors.New("Channel is closed"))
}
if e.Origin!=Mouse||e.Type!=LookInBody||e.Begin!=len(msg)||e.End!=len(msg)+len(test)||e.Text!=test{
t.Fatal(errors.New(fmt.Sprintf("Something wrong with event: %#v",e)))
}
if _,err:=w.Write([]byte("\nPress middle button of mouse on Del in the window's tag"));err!=nil{
t.Fatal(err)
}
for e,ok= <-ch;ok;e,ok= <-ch{
if e.Origin==Mouse{
break
}
}
if!ok{
t.Fatal(errors.New("Channel is closed"))
}
if e.Origin!=Mouse||e.Type!=ExecuteInTag||e.Text!="Del"{
t.Fatal(errors.New(fmt.Sprintf("Something wrong with event: %#v",e)))
}
if err:=w.UnreadEvent(e);err!=nil{
t.Fatal(err)
}
}



/*:65*/


//line goacme.w:123




/*:3*/


