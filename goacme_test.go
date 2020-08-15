

/*3:*/


//line goacme.w:31

package goacme

import(
"os/exec"
"9fans.net/go/plan9/client"
"testing"


/*12:*/


//line goacme.w:130

"fmt"
"time"
"9fans.net/go/plan9"




/*:12*/



/*17:*/


//line goacme.w:183

"bytes"
"errors"



/*:17*/


//line goacme.w:38

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


/*13:*/


//line goacme.w:137

time.Sleep(time.Second)



/*:13*/


//line goacme.w:51

}
}



/*14:*/


//line goacme.w:141

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



/*:14*/



/*18:*/


//line goacme.w:188

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



/*:18*/



/*26:*/


//line goacme.w:289

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





/*:26*/



/*54:*/


//line goacme.w:644

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
ch,err:=w.EventChannel(0,Look|Execute)
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
if _,err:=w.Write([]byte("\nChording test: select 'argument', press middle button of mouse on 'Execute' and press left button of mouse without releasing middle button"));err!=nil{
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



/*:54*/



/*58:*/


//line goacme.w:735

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



/*:58*/



/*61:*/


//line goacme.w:813

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



/*:61*/



/*71:*/


//line goacme.w:932

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





/*:71*/



/*82:*/


//line goacme.w:1074

func TestLog(t*testing.T){
l,err:=OpenLog()
if err!=nil{
t.Fatal(err)
}
defer l.Close()
ch,err:=l.EventChannel(NewWin)
if err!=nil{
t.Fatal(err)
}
w,err:=New()
if err!=nil{
t.Fatal(err)
}
defer w.Del(true)
defer w.Close()
ev,ok:=<-ch
if!ok{
t.Fatal(errors.New("cannot read an event from log"))
}
if w.id!=ev.Id{
t.Fatal(errors.New("unexpected window id"))
}
}



/*:82*/



/*89:*/


//line goacme.w:1170

func TestWindowsInfo(t*testing.T){
l1,err:=WindowsInfo()
if err!=nil{
t.Fatal(err)
}
w,err:=New()
if err!=nil{
t.Fatal(err)
}
defer w.Close()
l2,err:=WindowsInfo()
if err!=nil{
t.Fatal(err)
}
if len(l1)==len(l2)||l2[len(l2)-1].Id!=w.id{
t.Fatal(errors.New(fmt.Sprintf("something wrong with window list: %v, %v",l1,l2)))
}
if _,err:=l1.Get(w.id);err==nil{
t.Fatal(errors.New(fmt.Sprintf(fmt.Sprintf("window with id=%d has been found",w.id))))
}
if i2,err:=l2.Get(w.id);err!=nil||i2.Id!=w.id{
t.Fatal(errors.New(fmt.Sprintf(fmt.Sprintf("window with id=%d has not been found",w.id))))
}

w.Del(true)
l2,err= WindowsInfo()
if err!=nil{
t.Fatal(err)
}
if len(l1)!=len(l2){
t.Fatal(errors.New(fmt.Sprintf("sizes of window lists mismatched: %v, %v",l1,l2)))
}
if _,err:=l2.Get(w.id);err==nil{
t.Fatal(errors.New(fmt.Sprintf(fmt.Sprintf("window with id=%d has been found",w.id))))
}
}



/*:89*/


//line goacme.w:55




/*:3*/


