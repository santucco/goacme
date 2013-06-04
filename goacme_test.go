

/*3:*/


//line goacme.w:97

package goacme

import(
"os"
"os/exec"
"testing"


/*11:*/


//line goacme.w:184

"fmt"



/*:11*/



/*12:*/


//line goacme.w:190

"time"



/*:12*/



/*17:*/


//line goacme.w:240

"bytes"
"errors"



/*:17*/


//line goacme.w:104

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



/*14:*/


//line goacme.w:198

func TestNewOpen(t*testing.T){
prepare(t)


/*13:*/


//line goacme.w:194

time.Sleep(time.Second)



/*:13*/


//line goacme.w:201

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



/*:14*/



/*18:*/


//line goacme.w:245

func TestReadWrite(t*testing.T){


/*13:*/


//line goacme.w:194

time.Sleep(time.Second)



/*:13*/


//line goacme.w:247

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



/*27:*/


//line goacme.w:356

func TestPipeTo(t*testing.T){


/*13:*/


//line goacme.w:194

time.Sleep(time.Second)



/*:13*/


//line goacme.w:358

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



/*:27*/



/*29:*/


//line goacme.w:416

func TestPipeFrom(t*testing.T){


/*13:*/


//line goacme.w:194

time.Sleep(time.Second)



/*:13*/


//line goacme.w:418

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


/*13:*/


//line goacme.w:194

time.Sleep(time.Second)



/*:13*/


//line goacme.w:435

p,err:=w.PipeFrom("body",f,"cat")
if err!=nil{
t.Fatal(err)
}


/*13:*/


//line goacme.w:194

time.Sleep(time.Second)



/*:13*/


//line goacme.w:440

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



/*:29*/



/*31:*/


//line goacme.w:486

func TestSysRun(t*testing.T){


/*13:*/


//line goacme.w:194

time.Sleep(time.Second)



/*:13*/


//line goacme.w:488

s:="test"
f,p,err:=SysRun("echo","-n",s)
if err!=nil{
t.Fatal(err)
}
p.Wait()
p.Release()


/*13:*/


//line goacme.w:194

time.Sleep(time.Second)



/*:13*/


//line goacme.w:496

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




/*:31*/



/*33:*/


//line goacme.w:528

func TestDel(t*testing.T){


/*13:*/


//line goacme.w:194

time.Sleep(time.Second)



/*:13*/


//line goacme.w:530

w,err:=New()
if err!=nil{
t.Fatal(err)
}


/*13:*/


//line goacme.w:194

time.Sleep(time.Second)



/*:13*/


//line goacme.w:535

w.Del(true)
w.Close()


/*13:*/


//line goacme.w:194

time.Sleep(time.Second)



/*:13*/


//line goacme.w:538

if _,err:=Open(w.id);err==nil{
t.Fatal(errors.New(fmt.Sprintf("window %d is still opened",w.id)))
}
}




/*:33*/



/*40:*/


//line goacme.w:597

func TestDeleteAll(t*testing.T){


/*13:*/


//line goacme.w:194

time.Sleep(time.Second)



/*:13*/


//line goacme.w:599

var l[10]int
for i:=0;i<len(l);i++{


/*13:*/


//line goacme.w:194

time.Sleep(time.Second)



/*:13*/


//line goacme.w:602

w,err:=New()
if err!=nil{
t.Fatal(err)
}
l[i]= w.id
}
DeleteAll()


/*13:*/


//line goacme.w:194

time.Sleep(time.Second)



/*:13*/


//line goacme.w:610

for _,v:=range l{
_,err:=Open(v)
if err==nil{
t.Fatal(errors.New(fmt.Sprintf("window %d is still opened",v)))
}
}
}



/*:40*/



/*65:*/


//line goacme.w:901

func TestEvent(t*testing.T){


/*13:*/


//line goacme.w:194

time.Sleep(time.Second)



/*:13*/


//line goacme.w:903

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


//line goacme.w:122




/*:3*/


