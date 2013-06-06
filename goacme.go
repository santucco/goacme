

/*2:*/


//line goacme.w:49

// Copyright (c) 2013 Alexander Sychev. All rights reserved.
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

// Package goacme provides interface to acme programming environment
package goacme

import(


/*6:*/


//line goacme.w:137

"code.google.com/p/goplan9/plan9"
"code.google.com/p/goplan9/plan9/client"
"sync"
"fmt"



/*:6*/



/*20:*/


//line goacme.w:298

"io"



/*:20*/



/*25:*/


//line goacme.w:334

"os"
"os/exec"



/*:25*/



/*48:*/


//line goacme.w:692

"errors"



/*:48*/


//line goacme.w:81

)

type(


/*4:*/


//line goacme.w:127

// Window is a structure to manipulate a particular acme's window.
Window struct{
id int


/*21:*/


//line goacme.w:302

files map[string]*client.Fid



/*:21*/



/*36:*/


//line goacme.w:558

prev*Window
next*Window



/*:36*/



/*60:*/


//line goacme.w:819

ch chan*Event



/*:60*/


//line goacme.w:131

}



/*:4*/



/*42:*/


//line goacme.w:623

Event struct{


/*47:*/


//line goacme.w:687

// Origin will be an origin of action with type ActionOrigin
Origin ActionOrigin



/*:47*/



/*53:*/


//line goacme.w:729

// Type will be an type of action with type ActionType
Type ActionType



/*:53*/



/*56:*/


//line goacme.w:756

begin int
// Begin is a start address of a text of the action
Begin int
end int
// End is an end address of the text of the action
End int



/*:56*/



/*58:*/


//line goacme.w:773

flag int
// IsBuiltin is a flag the action is recognised like a \.{Acme}'s builtin
IsBuiltin bool
// NoLoad is a flag of acme can interpret the action without loading a new file
NoLoad bool
// IsFile is a flag the Text is a file or window name
IsFile bool
// Text is a text arguments of the action, perhaps with address
Text string
// Arg is a text of chorded argument if any
Arg string



/*:58*/


//line goacme.w:625

}



/*:42*/



/*45:*/


//line goacme.w:669

// ActionOrigin is a origin of the action
ActionOrigin int



/*:45*/



/*51:*/


//line goacme.w:711

// ActionType is a type of the action
ActionType int



/*:51*/


//line goacme.w:85

)

var(


/*7:*/


//line goacme.w:144

fsys*client.Fsys



/*:7*/



/*35:*/


//line goacme.w:553

fwin*Window
lwin*Window



/*:35*/



/*49:*/


//line goacme.w:696

// ErrInvalidOrigin will be returned if Origin is unexptected 
ErrInvalidOrigin= errors.New("invalid origin of action")



/*:49*/



/*54:*/


//line goacme.w:734

// ErrInvalidType will be returned if Type is unexpected 
ErrInvalidType= errors.New("invalid type of action")



/*:54*/



/*62:*/


//line goacme.w:847

// ErrChannelAlreadyOpened will be returned 
// if channel of events is opened by call of EventChannel 
ErrChannelAlreadyOpened= errors.New("channel of events is already opened")



/*:62*/


//line goacme.w:89

)



/*46:*/


//line goacme.w:674

const(
// Edit is the origin for writes to the body or tag file
Edit ActionOrigin= 1<<iota
// File is the origin for through the window's other files
File
// Keyboard is the origin for keyboard actions
Keyboard
// Mouse is the origin for mouse actions
Mouse
)



/*:46*/



/*52:*/


//line goacme.w:716

const(
DeletedFromBody ActionType= 1<<iota
DeletedFromTag
InsertInBody
InsertInTag
LookInBody
LookInTag
ExecuteInBody
ExecuteInTag
)



/*:52*/


//line goacme.w:92





/*:2*/



/*9:*/


//line goacme.w:158

// New creates a new window and returns *Window or error
func New()(*Window,error){


/*8:*/


//line goacme.w:148

{
var err error
new(sync.Once).Do(func(){fsys,err= client.MountService("acme")})
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
new(sync.Once).Do(func(){fsys,err= client.MountService("acme")})
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


//line goacme.w:306

this.files= make(map[string]*client.Fid)



/*:22*/



/*37:*/


//line goacme.w:563

this.prev= lwin
this.next= nil
if fwin==nil{
fwin= this
}
if lwin!=nil{
lwin.next= this
}
lwin= this



/*:37*/


//line goacme.w:183

return this,nil
}



/*:10*/



/*11:*/


//line goacme.w:188

// Close releases all resources of the window
func(this*Window)Close()error{


/*23:*/


//line goacme.w:310

for _,v:=range this.files{
v.Close()
}



/*:23*/



/*38:*/


//line goacme.w:575

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



/*:38*/


//line goacme.w:191

return nil
}



/*:11*/



/*15:*/


//line goacme.w:223

// Read reads len(p) bytes from "body" of the window. 
// It returns count of readed bytes or error.
func(this*Window)Read(p[]byte)(int,error){
f,err:=this.File("body")
if err!=nil{
return 0,err
}
return f.Read(p)
}



/*:15*/



/*16:*/


//line goacme.w:235

// Write writes len(p) bytes to "body" of the window.
// It returns count of writen bytes or error.
func(this*Window)Write(p[]byte)(int,error){
f,err:=this.File("body")
if err!=nil{
return 0,err
}
return f.Write(p)
}



/*:16*/



/*19:*/


//line goacme.w:282

// Seek sets the offset for the next Read or Write to offset, interpreted
// according to whence: 0 means relative to the origin of the file, 1 means
// relative to the current offset, and 2 means relative to the end. Seek
// returns the new offset and an Error, if any.
func(this*Window)Seek(offset int64,whence int)(ret int64,err error){
f,err:=this.File("body")
if err!=nil{
return 0,err
}
return f.Seek(offset,whence)
}





/*:19*/



/*24:*/


//line goacme.w:316

// File returns io.ReadWriteSeeker of corresponding file of the windows or error if file doesn't exist.
// Opened file will be cached inside the window, so caller shouldn't close it.
func(this*Window)File(file string)(io.ReadWriteSeeker,error){
f:=this.files[file]
if f!=nil{
return f,nil
}
f,err:=fsys.Open(fmt.Sprintf("%d/%s",this.id,file),plan9.ORDWR)
if err!=nil{
return nil,err
}
this.files[file]= f
return f,nil
}



/*:24*/



/*26:*/


//line goacme.w:339

// PipeTo runs shell command line cmd with /dev/null on standard input and the window's file 
// on standard output.  If stderr is non-zero, it is used as standard error.
// Otherwise the command inherits the caller's standard error.
// PipeTo returns *os.Process of created process or error.
// Caller has to wait the running process and read its exit status.
func(this*Window)PipeTo(name string,stderr io.Writer,cmd...string)(*os.Process,error){
if len(cmd)==0{
return nil,os.ErrInvalid
}
c:=exec.Command(cmd[0],cmd[1:]...)
if stderr==nil{
stderr= os.Stderr
}
c.Stdin= nil
c.Stderr= stderr
var err error
if c.Stdout,err= this.File(name);err!=nil{
return nil,err
}
if err:=c.Start();err!=nil{
return nil,err
}
return c.Process,nil
}



/*:26*/



/*28:*/


//line goacme.w:399

// PipeFrom runs a shell command line cmd with the
// window's file on standard input.  The command runs with
// stdout as its standard output and standard error.
// PipeFrom returns *os.Process of created process or error.
// Caller has to wait the running process and read its exit status.
func(this*Window)PipeFrom(name string,stdout io.Writer,cmd...string)(*os.Process,error){
if len(cmd)==0{
return nil,os.ErrInvalid
}
c:=exec.Command(cmd[0],cmd[1:]...)
if stdout==nil{
stdout= os.Stdout
}
c.Stdout= stdout
c.Stderr= os.Stderr
var err error
if c.Stdin,err= this.File(name);err!=nil{
return nil,err
}
if err:=c.Start();err!=nil{
return nil,err
}
return c.Process,nil
}



/*:28*/



/*30:*/


//line goacme.w:467

// SysRun runs shell command cmd and returns io.ReadSeeker of a window, 
// *os.Process of a running process or error.
// Caller has to wait the running process and read its exit status.
func SysRun(cmd...string)(io.ReadSeeker,*os.Process,error){
w,err:=New()
if err!=nil{
return nil,nil,err
}
if len(cmd)==0{
w.Close()
return nil,nil,os.ErrInvalid
}
c:=exec.Command(cmd[0],cmd[1:]...)
f,err:=w.File("body")
if err!=nil{
return nil,nil,err
}
c.Stdout= f
if err:=c.Start();err!=nil{
w.Close()
return nil,nil,err
}
return f,c.Process,nil
}



/*:30*/



/*32:*/


//line goacme.w:519

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



/*:32*/



/*39:*/


//line goacme.w:591

// DeleteAll deletes all the windows opened in a session
func DeleteAll(){
for fwin!=nil{
fwin.Del(true)
fwin.Close()
}
}



/*:39*/



/*43:*/


//line goacme.w:631

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



/*:43*/



/*44:*/


//line goacme.w:654

func readEvent(r io.Reader)(*Event,error){
o,t,b,e,f,s,err:=readFields(r)
if err!=nil{
return nil,err
}
var ev Event


/*50:*/


//line goacme.w:701

switch o{
case'E':ev.Origin= Edit
case'F':ev.Origin= File
case'K':ev.Origin= Keyboard
case'M':ev.Origin= Mouse
default:return nil,ErrInvalidOrigin
}



/*:50*/


//line goacme.w:661



/*55:*/


//line goacme.w:739

switch t{
case'D':ev.Type= DeletedFromBody
case'd':ev.Type= DeletedFromTag
case'I':ev.Type= InsertInBody
case'i':ev.Type= InsertInTag
case'L':ev.Type= LookInBody
case'l':ev.Type= LookInTag
case'X':ev.Type= ExecuteInBody
case'x':ev.Type= ExecuteInTag
default:return nil,ErrInvalidType
}



/*:55*/


//line goacme.w:662



/*57:*/


//line goacme.w:765

ev.begin= b
ev.Begin= b
ev.end= e
ev.End= e



/*:57*/


//line goacme.w:663



/*59:*/


//line goacme.w:787

ev.flag= f

if ev.Type==ExecuteInBody||ev.Type==ExecuteInTag{
ev.IsBuiltin= (ev.flag&1)==1
}else if ev.Type==LookInBody||ev.Type==LookInTag{
ev.NoLoad= (ev.flag&1)==1
ev.IsFile= (ev.flag&4)==4
}

ev.Text= s

// if there is an expantion
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



/*:59*/


//line goacme.w:664

return&ev,nil
}



/*:44*/



/*61:*/


//line goacme.w:823

// EventChannel returns a channel of *Event from which events can be read or error.
// First call of EventChannel starts a goroutine to read events from "event" file 
// and put them to the channel. Subsequent calls of EventChannel will return the same channel.
func(this*Window)EventChannel()(<-chan*Event,error){
if this.ch!=nil{
return this.ch,nil
}
f,err:=this.File("event")
if err!=nil{
return nil,err
}
this.ch= make(chan*Event)
go func(){
for ev,err:=readEvent(f);err==nil;ev,err= readEvent(f){
this.ch<-ev
}
close(this.ch)
this.ch= nil
}()
return this.ch,nil
}



/*:61*/



/*63:*/


//line goacme.w:853

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




/*:63*/



/*64:*/


//line goacme.w:868

// UnreadEvent writes event ev back to the "event" file, 
// indicating to acme that it should be handled internally.
func(this*Window)UnreadEvent(ev*Event)error{
f,err:=this.File("event")
if err!=nil{
return err
}
var o rune
switch ev.Origin{
case Edit:o= 'E'
case File:o= 'F'
case Keyboard:o= 'K'
case Mouse:o= 'M'
default:return ErrInvalidOrigin
}
var t rune
switch ev.Type{
case DeletedFromBody:t= 'D'
case DeletedFromTag:t= 'd'
case InsertInBody:t= 'I'
case InsertInTag:t= 'i'
case LookInBody:t= 'L'
case LookInTag:t= 'l'
case ExecuteInBody:t= 'X'
case ExecuteInTag:t= 'x'
default:return ErrInvalidType
}

_,err= fmt.Fprintf(f,"%c%c%d %d \n",o,t,ev.begin,ev.end)
return err
}



/*:64*/


