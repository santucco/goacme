

/*4:*/


//line goacme.w:99

package goacme

import(
"os"
"os/exec"
"code.google.com/p/goplan9/plan9/client"
"testing"


/*13:*/


//line goacme.w:198

"fmt"
"time"
"code.google.com/p/goplan9/plan9"




/*:13*/



/*18:*/


//line goacme.w:250

"bytes"
"errors"



/*:18*/


//line goacme.w:107

)

func prepare(t*testing.T){
_,err:=client.MountService("acme")
if err==nil{
t.Log("acme started already")
}else{
cmd:=exec.Command("acme")
err= cmd.Start()
if err!=nil{
t.Fatal(err)
}


/*14:*/


//line goacme.w:205

time.Sleep(time.Second)



/*:14*/


//line goacme.w:120

}
}



/*15:*/


//line goacme.w:209

func TestNewOpen(t*testing.T){
prepare(t)
w,err:=New()
if err!=nil{
t.Fatal(err)
}
defer w.Close()
defer w.Del(true)
if f,err:=fsys.Open(fmt.Sprintf("%d",w.id),plan9.OREAD);err!=nil{
t.Fatal(err)
}else{
f.Close()
}
}



/*:15*/



/*19:*/


//line goacme.w:255

func TestReadWrite(t*testing.T){
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



/*:19*/



/*28:*/


//line goacme.w:368

func TestPipeTo(t*testing.T){
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


/*14:*/


//line goacme.w:205

time.Sleep(time.Second)



/*:14*/


//line goacme.w:383

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



/*:28*/



/*30:*/


//line goacme.w:428

func TestPipeFrom(t*testing.T){
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
p,err:=w.PipeFrom("body",f,"cat")
if err!=nil{
t.Fatal(err)
}
w.Del(true)
w.Close()
p.Wait()
p.Release()


/*14:*/


//line goacme.w:205

time.Sleep(time.Second)



/*:14*/


//line goacme.w:454

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



/*:30*/



/*32:*/


//line goacme.w:496

func TestSysRun(t*testing.T){
s:="test"
f,p,err:=SysRun("echo","-n",s)
if err!=nil{
t.Fatal(err)
}
p.Wait()
p.Release()


/*14:*/


//line goacme.w:205

time.Sleep(time.Second)



/*:14*/


//line goacme.w:505

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



/*:32*/



/*34:*/


//line goacme.w:536

func TestDel(t*testing.T){
w,err:=New()
if err!=nil{
t.Fatal(err)
}
w.Del(true)
w.Close()
if _,err:=Open(w.id);err==nil{
t.Fatal(errors.New(fmt.Sprintf("window %d is still opened",w.id)))
}
}




/*:34*/



/*41:*/


//line goacme.w:603

func TestDeleteAll(t*testing.T){
var l[10]int
for i:=0;i<len(l);i++{
w,err:=New()
if err!=nil{
t.Fatal(err)
}
l[i]= w.id
}
DeleteAll()
for _,v:=range l{
_,err:=Open(v)
if err==nil{
t.Fatal(errors.New(fmt.Sprintf("window %d is still opened",v)))
}
}
}



/*:41*/



/*66:*/


//line goacme.w:916

func TestEvent(t*testing.T){
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
ch,err:=w.EventChannel(0,Mouse,Look|Execute)
if err!=nil{
t.Fatal(err)
}
e,ok:=<-ch
if!ok{
t.Fatal(errors.New("Channel is closed"))
}
if e.Origin!=Mouse||e.Type!=Look||e.Begin!=len(msg)||e.End!=len(msg)+len(test)||e.Text!=test{
t.Fatal(errors.New(fmt.Sprintf("Something wrong with event: %#v",e)))
}
if _,err:=w.Write([]byte("\nChording test: select argument, press middle button of mouse on Execute and press left button of mouse without releasing middle button"));err!=nil{
t.Fatal(err)
}
e,ok= <-ch
if!ok{
t.Fatal(errors.New("Channel is closed"))
}
if e.Origin!=Mouse||e.Type!=(Execute)||e.Text!="Execute"||e.Arg!="argument"{
t.Fatal(errors.New(fmt.Sprintf("Something wrong with event: %#v",e)))
}
if err:=w.UnreadEvent(e);err!=nil{
t.Fatal(err)
}
if _,err:=w.Write([]byte("\nPress middle button of mouse on Del in the window's tag"));err!=nil{
t.Fatal(err)
}
e,ok= <-ch
if!ok{
t.Fatal(errors.New("Channel is closed"))
}
if e.Origin!=Mouse||e.Type!=(Execute|Tag)||e.Text!="Del"{
t.Fatal(errors.New(fmt.Sprintf("Something wrong with event: %#v",e)))
}
if err:=w.UnreadEvent(e);err!=nil{
t.Fatal(err)
}
}



/*:66*/



/*69:*/


//line goacme.w:1000

func TestWriteReadAddr(t*testing.T){
w,err:=New()
if err!=nil{
t.Fatal(err)
}
defer w.Close()
defer w.Del(true)
if b,e,err:=w.ReadAddr();err!=nil{
t.Fatal(err)
}else if b!=0||e!=0{
t.Fatal(errors.New(fmt.Sprintf("Something wrong with address: %v, %v",b,e)))
}
if _,err:=w.Write([]byte("test"));err!=nil{
t.Fatal(err)
}
if err:=w.WriteAddr("0,$");err!=nil{
t.Fatal(err)
}
if b,e,err:=w.ReadAddr();err!=nil{
t.Fatal(err)
}else if b!=0||e!=4{
t.Fatal(errors.New(fmt.Sprintf("Something wrong with address: %v, %v",b,e)))
}
}



/*:69*/



/*72:*/


//line goacme.w:1075

func TestWriteReadCtl(t*testing.T){
w,err:=New()
if err!=nil{
t.Fatal(err)
}
defer w.Close()
defer w.Del(true)
if _,err:=w.Write([]byte("test"));err!=nil{
t.Fatal(err)
}
if _,_,_,_,d,_,_,_,err:=w.ReadCtl();err!=nil{
t.Fatal(err)
}else if!d{
t.Fatal(errors.New(fmt.Sprintf("The window has to be dirty\n")))
}
if err:=w.WriteCtl("clean");err!=nil{
t.Fatal(err)
}
if _,_,_,_,d,_,_,_,err:=w.ReadCtl();err!=nil{
t.Fatal(err)
}else if d{
t.Fatal(errors.New(fmt.Sprintf("The window has to be clean\n")))
}
}




/*:72*/


//line goacme.w:124




/*:4*/


