

/*2:*/


//line goacme.w:49

// Copyright (c) 2013, 2014, 2020 Alexander Sychev. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * The name of author may not be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.



/*:2*/



/*3:*/


//line goacme.w:78


// Package goacme provides interface to acme programming environment
package goacme

import(


/*6:*/


//line goacme.w:136

"9fans.net/go/plan9"
"9fans.net/go/plan9/client"
"sync"
"fmt"



/*:6*/



/*20:*/


//line goacme.w:301

"io"



/*:20*/



/*41:*/


//line goacme.w:515

"errors"



/*:41*/



/*53:*/


//line goacme.w:644

"strings"



/*:53*/



/*87:*/


//line goacme.w:1204

"bufio"
"sort"



/*:87*/


//line goacme.w:84

)

type(


/*5:*/


//line goacme.w:128

// Window is a structure to manipulate a particular acme's window.
Window struct{
id int


/*21:*/


//line goacme.w:305

files map[string]*client.Fid



/*:21*/



/*29:*/


//line goacme.w:380

prev*Window
next*Window



/*:29*/



/*55:*/


//line goacme.w:663

ch chan*Event



/*:55*/


//line goacme.w:132

}



/*:5*/



/*35:*/


//line goacme.w:445

Event struct{


/*40:*/


//line goacme.w:510

// Origin will be an origin of action with type ActionOrigin
Origin ActionOrigin



/*:40*/



/*46:*/


//line goacme.w:553

// Type will be an type of action with type ActionType
Type ActionType



/*:46*/



/*49:*/


//line goacme.w:580

begin int
// Begin is a start address of a text of the action
Begin int
end int
// End is an end address of the text of the action
End int



/*:49*/



/*51:*/


//line goacme.w:597

flag int
// IsBuiltin is a flag the action is recognised like an \.{Acme}'s builtin
IsBuiltin bool
// NoLoad is a flag of acme can interpret the action without loading a new file
NoLoad bool
// IsFile is a flag the Text is a file or window name
IsFile bool
// Text is a text arguments of the action, perhaps with address
Text string
// Arg is a text of chorded argument if any
Arg string



/*:51*/


//line goacme.w:447

}



/*:35*/



/*38:*/


//line goacme.w:491

// ActionOrigin is a origin of the action
ActionOrigin int



/*:38*/



/*44:*/


//line goacme.w:534

// ActionType is a type of the action
ActionType int



/*:44*/



/*69:*/


//line goacme.w:979

wrapper struct{
f io.ReadWriteSeeker
}



/*:69*/



/*72:*/


//line goacme.w:1020

Log struct{
fid*client.Fid


/*80:*/


//line goacme.w:1109

ch chan*LogEvent



/*:80*/


//line goacme.w:1023

}



/*:72*/



/*75:*/


//line goacme.w:1047

// OperationType is a type of the operation
OperationType int



/*:75*/



/*77:*/


//line goacme.w:1063

LogEvent struct{
Id int
Type OperationType
Name string
}



/*:77*/



/*84:*/


//line goacme.w:1167

Info struct{
Id int
TagSize int
BodySize int
IsDirectory bool
IsDirty bool
Tag[]string
}

Infos[]*Info



/*:84*/


//line goacme.w:88

)

var(


/*7:*/


//line goacme.w:143

fsys*client.Fsys
once sync.Once



/*:7*/



/*28:*/


//line goacme.w:375

fwin*Window
lwin*Window



/*:28*/



/*42:*/


//line goacme.w:519

// ErrInvalidOrigin will be returned if a case of an unexpected origin of action
ErrInvalidOrigin= errors.New("invalid origin of action")



/*:42*/



/*47:*/


//line goacme.w:558

// ErrInvalidType will be returned if a case of an unexpected type of action
ErrInvalidType= errors.New("invalid type of action")



/*:47*/



/*58:*/


//line goacme.w:732

// ErrChannelAlreadyOpened will be returned
// if channel of events is opened by call of EventChannel
ErrChannelAlreadyOpened= errors.New("channel of events is already opened")



/*:58*/



/*78:*/


//line goacme.w:1071

operations= map[string]OperationType{
"new":NewWin,
"zerox":Zerox,
"get":Get,
"put":Put,
"del":DelWin,
"focus":Focus,
}



/*:78*/


//line goacme.w:92

)



/*39:*/


//line goacme.w:496

const(
Unknown ActionOrigin= 0
// Edit is the origin for writes to the body or tag file
Edit ActionOrigin= 1<<iota
// File is the origin for through the window's other files
File
// Keyboard is the origin for keyboard actions
Keyboard
// Mouse is the origin for mouse actions
Mouse
)



/*:39*/



/*45:*/


//line goacme.w:539

const(
Delete ActionType= 1<<iota
Insert
Look
Execute
// Tag is a flag points out the event has occured in the tag of the window
Tag
// TagMask is a mask points out the event should be masked by tag
TagMask
AllTypes= Delete|Insert|Look|Execute
)



/*:45*/



/*76:*/


//line goacme.w:1052

const(
NewWin OperationType= 1<<iota
Zerox
Get
Put
DelWin
Focus
)



/*:76*/


//line goacme.w:95





/*:3*/



/*9:*/


//line goacme.w:158

// New creates a new window and returns *Window or error
func New()(*Window,error){


/*8:*/


//line goacme.w:148

{
var err error
once.Do(func(){fsys,err= client.MountService("acme")})
if err!=nil{
return nil,err
}
}



/*:8*/


//line goacme.w:161

f,err:=fsys.Open("new/ctl",plan9.OREAD)
if err!=nil{
return nil,err
}
defer f.Close()
var id int
if _,err:=fmt.Fscan(f,&id);err!=nil{
return nil,err
}
return Open(id)
}



/*:9*/



/*10:*/


//line goacme.w:175

// Open opens a window with a specified id and returns *Window or error
func Open(id int)(*Window,error){


/*8:*/


//line goacme.w:148

{
var err error
once.Do(func(){fsys,err= client.MountService("acme")})
if err!=nil{
return nil,err
}
}



/*:8*/


//line goacme.w:178

if err:=fsys.Access(fmt.Sprintf("%d",id),plan9.OREAD);err!=nil{
return nil,err
}
this:=&Window{id:id}


/*22:*/


//line goacme.w:309

this.files= make(map[string]*client.Fid)



/*:22*/



/*30:*/


//line goacme.w:385

this.prev= lwin
this.next= nil
if fwin==nil{
fwin= this
}
if lwin!=nil{
lwin.next= this
}
lwin= this



/*:30*/



/*64:*/


//line goacme.w:866

if _,err:=this.File("addr");err!=nil{
return nil,err
}




/*:64*/


//line goacme.w:183

return this,nil
}



/*:10*/



/*11:*/


//line goacme.w:188

// Close releases all resources of the window
func(this*Window)Close()error{


/*23:*/


//line goacme.w:313

for _,v:=range this.files{
v.Close()
}



/*:23*/



/*31:*/


//line goacme.w:397

if this.next!=nil{
this.next.prev= this.prev
}
if this.prev!=nil{
this.prev.next= this.next
}
if fwin==this{
fwin= this.next
}
if lwin==this{
lwin= this.prev
}



/*:31*/


//line goacme.w:191

return nil
}



/*:11*/



/*15:*/


//line goacme.w:225

// Read reads len(p) bytes from "body" file of the window.
// Read returns a count of read bytes or error.
func(this*Window)Read(p[]byte)(int,error){
f,err:=this.File("body")
if err!=nil{
return 0,err
}
return f.Read(p)
}



/*:15*/



/*16:*/


//line goacme.w:237

// Write writes len(p) bytes to "body" file of the window.
// Write returns a count of written bytes or error.
func(this*Window)Write(p[]byte)(int,error){
f,err:=this.File("body")
if err!=nil{
return 0,err
}


/*71:*/


//line goacme.w:1014

f= &wrapper{f:f}



/*:71*/


//line goacme.w:245

return f.Write(p)
}



/*:16*/



/*19:*/


//line goacme.w:285

// Seek sets a position for the next Read or Write to offset, interpreted
// according to whence: 0 means relative to the origin of the file, 1 means
// relative to the current offset, and 2 means relative to the end.
// Seek returns the new offset or error
func(this*Window)Seek(offset int64,whence int)(ret int64,err error){
f,err:=this.File("body")
if err!=nil{
return 0,err
}
return f.Seek(offset,whence)
}





/*:19*/



/*24:*/


//line goacme.w:319

// File returns io.ReadWriteSeeker of corresponding file of the windows or error
func(this*Window)File(file string)(io.ReadWriteSeeker,error){
fid,ok:=this.files[file]
if!ok{
var err error
if fid,err= fsys.Open(fmt.Sprintf("%d/%s",this.id,file),plan9.ORDWR);err!=nil{
if fid,err= fsys.Open(fmt.Sprintf("%d/%s",this.id,file),plan9.OREAD);err!=nil{
if fid,err= fsys.Open(fmt.Sprintf("%d/%s",this.id,file),plan9.OWRITE);err!=nil{
return nil,err
}
}
}
this.files[file]= fid
}
var f io.ReadWriteSeeker= fid


/*71:*/


//line goacme.w:1014

f= &wrapper{f:f}



/*:71*/


//line goacme.w:335

return f,nil
}



/*:24*/



/*25:*/


//line goacme.w:340

// Del deletes the window, without a prompt if sure is true.
func(this*Window)Del(sure bool)error{
f,err:=this.File("ctl")
if err!=nil{
return err
}
s:="del"
if sure{
s= "delete"
}
_,err= f.Write([]byte(s))
return err
}



/*:25*/



/*32:*/


//line goacme.w:413

// DeleteAll deletes all the windows opened in a session
func DeleteAll(){
for fwin!=nil{
fwin.Del(true)
fwin.Close()
}
}



/*:32*/



/*36:*/


//line goacme.w:453

func readFields(r io.Reader)(o rune,t rune,b int,e int,f int,s string,err error){
var l int
if _,err= fmt.Fscanf(r,"%c%c%d %d %d %d",&o,&t,&b,&e,&f,&l);err!=nil{
return
}
if l!=0{
rs:=make([]rune,l)
for i:=0;i<l;i++{
if _,err= fmt.Fscanf(r,"%c",&rs[i]);err!=nil{
return
}
}
s= string(rs)
}
var nl[1]byte
if _,err= r.Read(nl[:]);err!=nil{
return
}
return
}



/*:36*/



/*37:*/


//line goacme.w:476

func readEvent(r io.Reader)(*Event,error){
o,t,b,e,f,s,err:=readFields(r)
if err!=nil{
return nil,err
}
var ev Event


/*43:*/


//line goacme.w:524

switch o{
case'E':ev.Origin= Edit
case'F':ev.Origin= File
case'K':ev.Origin= Keyboard
case'M':ev.Origin= Mouse
default:ev.Origin= Unknown
}



/*:43*/


//line goacme.w:483



/*48:*/


//line goacme.w:563

switch t{
case'D':ev.Type= Delete
case'd':ev.Type= Delete|Tag
case'I':ev.Type= Insert
case'i':ev.Type= Insert|Tag
case'L':ev.Type= Look
case'l':ev.Type= Look|Tag
case'X':ev.Type= Execute
case'x':ev.Type= Execute|Tag
default:return nil,ErrInvalidType
}



/*:48*/


//line goacme.w:484



/*50:*/


//line goacme.w:589

ev.begin= b
ev.Begin= b
ev.end= e
ev.End= e



/*:50*/


//line goacme.w:485



/*52:*/


//line goacme.w:611

ev.flag= f

if ev.Type&Execute==Execute{
ev.IsBuiltin= (ev.flag&1)==1
}else if ev.Type&Look==Look{
ev.NoLoad= (ev.flag&1)==1
ev.IsFile= (ev.flag&4)==4
}

ev.Text= s

// if there is an expansion
if(ev.flag&2)==2{
_,_,ev.Begin,ev.End,_,ev.Text,err= readFields(r)
if err!=nil{
return nil,err
}
}
// if there is a chording
if(ev.flag&8)==8{
_,_,_,_,_,ev.Arg,err= readFields(r)
if err!=nil{
return nil,err
}
_,_,_,_,_,_,err= readFields(r)
if err!=nil{
return nil,err
}
}


/*54:*/


//line goacme.w:648

if len(ev.Text)> 0{
f:=strings.Fields(ev.Text)
if len(f)> 1{
ev.Text= f[0]
s:=ev.Arg
if len(s)> 0{
s= " "+ev.Arg
}
ev.Arg= strings.Join(f[1:]," ")+s
}
}




/*:54*/


//line goacme.w:641




/*:52*/


//line goacme.w:486

return&ev,nil
}



/*:37*/



/*56:*/


//line goacme.w:667

// EventChannel returns a channel of *Event with a buffer size
// from which events can be read or error.
// Only ActionTypes set in tmask are used.
// If TagMask is set in tmask, the event will be masked by tag. Otherwise Tag flag will be ignored.
// First call of EventChannel starts a goroutine to read events from "event" file
// and put them to the channel. Subsequent calls of EventChannel will return the same channel.
func(this*Window)EventChannel(size int,tmask ActionType)(<-chan*Event,error){
if this.ch!=nil{
return this.ch,nil
}


/*57:*/


//line goacme.w:707

old:=false
{
var em string
if tmask&Delete==Delete{
em+= "D"
}
if tmask&Insert==Insert{
em+= "I"
}
if tmask&Look==Look{
em+= "L"
}
if tmask&Execute==Execute{
em+= "X"
}
if tmask&TagMask!=TagMask{
em+= strings.ToLower(em)
}
if err:=this.WriteCtl("events %s\n",em);err!=nil{
old= true
}
}



/*:57*/


//line goacme.w:678

f,err:=this.File("event")
if err!=nil{
return nil,err
}
if tmask&TagMask!=TagMask{
tmask|= Tag
}
this.ch= make(chan*Event,size)
go func(){
for ev,err:=readEvent(f);err==nil;ev,err= readEvent(f){
if old&&ev.Type&tmask!=ev.Type{
if ev.Type&Insert!=Insert&&ev.Type&Delete!=Delete{
this.UnreadEvent(ev)
}
continue
}
this.ch<-ev
}
close(this.ch)
this.ch= nil
}()
return this.ch,nil
}




/*:56*/



/*59:*/


//line goacme.w:738

//  reads an event from "event" file of the window and returns *Event or error
func(this*Window)ReadEvent()(*Event,error){
if this.ch!=nil{
return nil,ErrChannelAlreadyOpened
}
f,err:=this.File("event")
if err!=nil{
return nil,err
}
return readEvent(f)
}




/*:59*/



/*60:*/


//line goacme.w:756

// UnreadEvent writes event ev back to the "event" file,
// indicating to acme that it should be handled internally.
func(this*Window)UnreadEvent(ev*Event)error{
f,err:=this.File("event")
if err!=nil{
return err
}
var o rune
switch ev.Origin{
case Mouse:o= 'M'
default:return ErrInvalidOrigin
}
var t rune
switch ev.Type{
case Look:t= 'L'
case Look|Tag:t= 'l'
case Execute:t= 'X'
case Execute|Tag:t= 'x'
default:return ErrInvalidType
}
_,err= fmt.Fprintf(f,"%c%c%d %d \n",o,t,ev.begin,ev.end)
return err
}



/*:60*/



/*62:*/


//line goacme.w:835

// WriteAddr writes format with args in "addr" file of the window
func(this*Window)WriteAddr(format string,args...interface{})error{
f,err:=this.File("addr")
if err!=nil{
return err
}
if len(args)> 0{
format= fmt.Sprintf(format,args...)
}
_,err= f.Write([]byte(format))
return err
}



/*:62*/



/*63:*/


//line goacme.w:850

// ReadAddr reads the address of the next read/write operation from "addr" file of the window.
// ReadAddr return begin and end offsets in symbols or error
func(this*Window)ReadAddr()(begin int,end int,err error){
f,err:=this.File("addr")
if err!=nil{
return
}
if _,err= f.Seek(0,0);err!=nil{
return
}
_,err= fmt.Fscanf(f,"%d %d",&begin,&end)
return
}



/*:63*/



/*66:*/


//line goacme.w:900

// WriteCtl writes format with args in "ctl" file of the window
// In case format is not ended by newline, '\n' will be added to the end of format
func(this*Window)WriteCtl(format string,args...interface{})error{
f,err:=this.File("ctl")
if err!=nil{
return err
}
if len(args)> 0{
format= fmt.Sprintf(format,args...)
}
if len(format)>=0&&format[len(format)-1]!='\n'{
format+= "\n"
}
if _,err= f.Seek(0,0);err!=nil{
return err
}
_,err= f.Write([]byte(format))
return err
}




/*:66*/



/*67:*/


//line goacme.w:923

// ReadCtl reads the address of the next read/write operation from "ctl" file of the window.
// ReadCtl returns:
//    id - the window ID
//    tlen - number of characters (runes) in the tag;
//    blen - number of characters in the body;
//    isdir -  true if the window is a directory, false otherwise;
//    isdirty - true if the window is modified, false otherwise;
//    wwidth - the width of the window in pixels;
//    font - the name of the font used in the window;
//    twidth - the width of a tab character in pixels;
//    error - in case of any error.
func(this*Window)ReadCtl()(id int,tlen int,blen int,isdir bool,isdirty bool,wwidth int,font string,twidth int,err error){
f,err:=this.File("ctl")
if err!=nil{
return
}
if _,err= f.(io.Seeker).Seek(0,0);err!=nil{
return
}
var dir,dirty int
_,err= fmt.Fscanf(f,"%d %d %d %d %d %d %s %d",&id,&tlen,&blen,&dir,&dirty,&wwidth,&font,&twidth)
isdir= dir==1
isdirty= dirty==1
return
}



/*:67*/



/*70:*/


//line goacme.w:985

func(this*wrapper)Read(p[]byte)(int,error){
return this.f.Read(p)
}

func(this*wrapper)Write(p[]byte)(int,error){
if len(p)<8168{
return this.f.Write(p)
}
c:=0
for i:=0;i<len(p);i+= 8168{
n:=i+8168
if n> len(p){
n= len(p)
}
n,e:=this.f.Write(p[i:n])
c+= n
if e!=nil{
return c,e
}
}
return c,nil
}

func(this*wrapper)Seek(offset int64,whence int)(ret int64,err error){
return this.f.Seek(offset,whence)
}



/*:70*/



/*73:*/


//line goacme.w:1027

// OpenLog opens the log and returns *Log or error
func OpenLog()(*Log,error){


/*8:*/


//line goacme.w:148

{
var err error
once.Do(func(){fsys,err= client.MountService("acme")})
if err!=nil{
return nil,err
}
}



/*:8*/


//line goacme.w:1030

f,err:=fsys.Open("log",plan9.OREAD)
if err!=nil{
return nil,err
}
return&Log{fid:f},nil
}



/*:73*/



/*74:*/


//line goacme.w:1039

// Close close the log
func(this*Log)Close()error{
return this.fid.Close()
}




/*:74*/



/*79:*/


//line goacme.w:1082

// Read reads a log of window operations  of the window from the log.
// Read returns LogEvent or error.
func(this*Log)Read()(*LogEvent,error){
var id int
var op string
var n string
var b[8168]byte
c,err:=this.fid.Read(b[:])
if err!=nil{
return nil,err
}
_,err= fmt.Sscan(string(b[:c]),&id,&op,&n)
if err!=nil{
_,err= fmt.Sscan(string(b[:c]),&id,&op)
}
if err!=nil{
return nil,err
}
t,ok:=operations[op]
if!ok{
return nil,errors.New("unexpected operation code")
}
return&LogEvent{Id:id,Type:t,Name:n},nil
}



/*:79*/



/*81:*/


//line goacme.w:1113

// EventChannel returns a channel of *LogEvent
// from which log events can be read or error.
// Only OperationType set in tmask are used.
// First call of EventChannel starts a goroutine to read events from the log
// and put them to the channel. Subsequent calls of EventChannel will return the same channel.
func(this*Log)EventChannel(tmask OperationType)(<-chan*LogEvent,error){
if this.ch!=nil{
return this.ch,nil
}
this.ch= make(chan*LogEvent)
go func(){
for ev,err:=this.Read();err==nil;ev,err= this.Read(){
if ev.Type&tmask!=ev.Type{
continue
}
this.ch<-ev
}
close(this.ch)
this.ch= nil
}()
return this.ch,nil
}



/*:81*/



/*85:*/


//line goacme.w:1180

func(this Infos)Len()int{
return len(this)
}

func(this Infos)Less(i,j int)bool{
return this[i].Id<this[j].Id
}

func(this Infos)Swap(i,j int){
this[i],this[j]= this[j],this[i]
}



/*:85*/



/*86:*/


//line goacme.w:1194

// Get returns Info by id or an error
func(this Infos)Get(id int)(*Info,error){
i:=sort.Search(this.Len(),func(i int)bool{return this[i].Id==id})
if i<this.Len()&&this[i].Id==id{
return this[i],nil
}
return nil,errors.New(fmt.Sprintf("window with id=%d has not been found",id))
}


/*:86*/



/*88:*/


//line goacme.w:1209

// WindowsInfo returns a list of the existing acme windows.
func WindowsInfo()(res Infos,err error){


/*8:*/


//line goacme.w:148

{
var err error
once.Do(func(){fsys,err= client.MountService("acme")})
if err!=nil{
return nil,err
}
}



/*:8*/


//line goacme.w:1212

f,err:=fsys.Open("index",plan9.OREAD)
if err!=nil{
return nil,err
}
defer f.Close()
r:=bufio.NewReader(f)
if r==nil{
return nil,errors.New("cannot create reader for index file")
}
for s,err:=r.ReadString('\n');err==nil;s,err= r.ReadString('\n'){
var id,ts,bs,d,m int
if _,err:=fmt.Sscanf(s,"%v %v %v %v %v",&id,&ts,&bs,&d,&m);err!=nil{
continue
}
res= append(res,&Info{Id:id,TagSize:ts,BodySize:bs,IsDirectory:d==1,IsDirty:m==1,Tag:strings.Split(s[12*5:]," ")})
}
sort.Sort(res)
return res,nil
}



/*:88*/


