

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



/*:2*/



/*3:*/


//line goacme.w:78

// Package goacme provides interface to acme programming environment
package goacme

import(


/*6:*/


//line goacme.w:135

"code.google.com/p/goplan9/plan9"
"code.google.com/p/goplan9/plan9/client"
"sync"
"fmt"



/*:6*/



/*20:*/


//line goacme.w:299

"io"



/*:20*/



/*41:*/


//line goacme.w:512

"errors"



/*:41*/



/*53:*/


//line goacme.w:640

"strings"



/*:53*/


//line goacme.w:83

)

type(


/*5:*/


//line goacme.w:127

// Window is a structure to manipulate a particular acme's window.
Window struct{
id int


/*21:*/


//line goacme.w:303

files map[string]*client.Fid



/*:21*/



/*29:*/


//line goacme.w:378

prev*Window
next*Window



/*:29*/



/*55:*/


//line goacme.w:659

ch chan*Event



/*:55*/


//line goacme.w:131

}



/*:5*/



/*35:*/


//line goacme.w:443

Event struct{


/*40:*/


//line goacme.w:507

// Origin will be an origin of action with type ActionOrigin
Origin ActionOrigin



/*:40*/



/*46:*/


//line goacme.w:549

// Type will be an type of action with type ActionType
Type ActionType



/*:46*/



/*49:*/


//line goacme.w:576

begin int
// Begin is a start address of a text of the action
Begin int
end int
// End is an end address of the text of the action
End int



/*:49*/



/*51:*/


//line goacme.w:593

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


//line goacme.w:445

}



/*:35*/



/*38:*/


//line goacme.w:489

// ActionOrigin is a origin of the action
ActionOrigin int



/*:38*/



/*44:*/


//line goacme.w:531

// ActionType is a type of the action
ActionType int



/*:44*/



/*67:*/


//line goacme.w:941

wrapper struct{
f io.ReadWriteSeeker
}



/*:67*/


//line goacme.w:87

)

var(


/*7:*/


//line goacme.w:142

fsys*client.Fsys



/*:7*/



/*28:*/


//line goacme.w:373

fwin*Window
lwin*Window



/*:28*/



/*42:*/


//line goacme.w:516

// ErrInvalidOrigin will be returned if a case of an unexpected origin of action
ErrInvalidOrigin= errors.New("invalid origin of action")



/*:42*/



/*47:*/


//line goacme.w:554

// ErrInvalidType will be returned if a case of an unexpected type of action
ErrInvalidType= errors.New("invalid type of action")



/*:47*/



/*57:*/


//line goacme.w:699

// ErrChannelAlreadyOpened will be returned
// if channel of events is opened by call of EventChannel
ErrChannelAlreadyOpened= errors.New("channel of events is already opened")



/*:57*/


//line goacme.w:91

)



/*39:*/


//line goacme.w:494

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



/*:39*/



/*45:*/


//line goacme.w:536

const(
Delete ActionType= 1<<iota
Insert
Look
Execute
// Tag is a flag points out the event has occured in the tag of the window
Tag
// TagMask is a mask points out the event should be masked by tag
TagMask
)



/*:45*/


//line goacme.w:94





/*:3*/



/*9:*/


//line goacme.w:156

// New creates a new window and returns *Window or error
func New()(*Window,error){


/*8:*/


//line goacme.w:146

{
var err error
new(sync.Once).Do(func(){fsys,err= client.MountService("acme")})
if err!=nil{
return nil,err
}
}



/*:8*/


//line goacme.w:159

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


//line goacme.w:173

// Open opens a window with a specified id and returns *Window or error
func Open(id int)(*Window,error){


/*8:*/


//line goacme.w:146

{
var err error
new(sync.Once).Do(func(){fsys,err= client.MountService("acme")})
if err!=nil{
return nil,err
}
}



/*:8*/


//line goacme.w:176

if err:=fsys.Access(fmt.Sprintf("%d",id),plan9.OREAD);err!=nil{
return nil,err
}
this:=&Window{id:id}


/*22:*/


//line goacme.w:307

this.files= make(map[string]*client.Fid)



/*:22*/



/*30:*/


//line goacme.w:383

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


//line goacme.w:181

return this,nil
}



/*:10*/



/*11:*/


//line goacme.w:186

// Close releases all resources of the window
func(this*Window)Close()error{


/*23:*/


//line goacme.w:311

for _,v:=range this.files{
v.Close()
}



/*:23*/



/*31:*/


//line goacme.w:395

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


//line goacme.w:189

return nil
}



/*:11*/



/*15:*/


//line goacme.w:223

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


//line goacme.w:235

// Write writes len(p) bytes to "body" file of the window.
// Write returns a count of written bytes or error.
func(this*Window)Write(p[]byte)(int,error){
f,err:=this.File("body")
if err!=nil{
return 0,err
}


/*69:*/


//line goacme.w:976

f= &wrapper{f:f}





/*:69*/


//line goacme.w:243

return f.Write(p)
}



/*:16*/



/*19:*/


//line goacme.w:283

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


//line goacme.w:317

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


/*69:*/


//line goacme.w:976

f= &wrapper{f:f}





/*:69*/


//line goacme.w:333

return f,nil
}



/*:24*/



/*25:*/


//line goacme.w:338

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


//line goacme.w:411

// DeleteAll deletes all the windows opened in a session
func DeleteAll(){
for fwin!=nil{
fwin.Del(true)
fwin.Close()
}
}



/*:32*/



/*36:*/


//line goacme.w:451

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


//line goacme.w:474

func readEvent(r io.Reader)(*Event,error){
o,t,b,e,f,s,err:=readFields(r)
if err!=nil{
return nil,err
}
var ev Event


/*43:*/


//line goacme.w:521

switch o{
case'E':ev.Origin= Edit
case'F':ev.Origin= File
case'K':ev.Origin= Keyboard
case'M':ev.Origin= Mouse
default:return nil,ErrInvalidOrigin
}



/*:43*/


//line goacme.w:481



/*48:*/


//line goacme.w:559

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


//line goacme.w:482



/*50:*/


//line goacme.w:585

ev.begin= b
ev.Begin= b
ev.end= e
ev.End= e



/*:50*/


//line goacme.w:483



/*52:*/


//line goacme.w:607

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


//line goacme.w:644

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


//line goacme.w:637




/*:52*/


//line goacme.w:484

return&ev,nil
}



/*:37*/



/*56:*/


//line goacme.w:663

// EventChannel returns a channel of *Event with a buffer size
// from which events can be read or error.
// Only ActionOrigins set in omask and ActionTypes set in tmask are used.
// If TagMask is set in tmask, the event will be masked by tag. Otherwise Tag flag will be ignored.
// First call of EventChannel starts a goroutine to read events from "event" file
// and put them to the channel. Subsequent calls of EventChannel will return the same channel.
func(this*Window)EventChannel(size int,omask ActionOrigin,tmask ActionType)(<-chan*Event,error){
if this.ch!=nil{
return this.ch,nil
}
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
if ev.Origin&omask!=ev.Origin||ev.Type&tmask!=ev.Type{
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



/*58:*/


//line goacme.w:705

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




/*:58*/



/*59:*/


//line goacme.w:720

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
case Delete:t= 'D'
case Delete|Tag:t= 'd'
case Insert:t= 'I'
case Insert|Tag:t= 'i'
case Look:t= 'L'
case Look|Tag:t= 'l'
case Execute:t= 'X'
case Execute|Tag:t= 'x'
default:return ErrInvalidType
}

_,err= fmt.Fprintf(f,"%c%c%d %d\n",o,t,ev.begin,ev.end)
return err
}



/*:59*/



/*61:*/


//line goacme.w:807

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



/*:61*/



/*62:*/


//line goacme.w:822

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



/*:62*/



/*64:*/


//line goacme.w:865

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
_,err= f.Write([]byte(format))
return err
}




/*:64*/



/*65:*/


//line goacme.w:885

// ReadCtl reads the address of the next read/write operation from "ctl" file of the window.
// ReadCtl returns:
//    id - the window ID
//    tlen - number of characters (runes) in the tag;
//    blen - number of characters in the body;
//    isdir -  true if the window is a directory, false otherwise;
//    isdirty - true if the window is modified, falseotherwise;
//    wwidth - the width of the window in pixels;
//    font - the name of the font used in the window;
//    twidth - the width of a tab character in pixels;
//    error - in case of any error.
func(this*Window)ReadCtl()(id int,tlen int,blen int,isdir bool,isdirty bool,wwidth int,font string,twidth int,err error){
f,err:=this.File("ctl")
if err!=nil{
return
}
if _,err= f.Seek(0,0);err!=nil{
return
}
var dir,dirty int
_,err= fmt.Fscanf(f,"%d %d %d %d %d %d %s %d",&id,&tlen,&blen,&dir,&dirty,&wwidth,&font,&twidth)
isdir= dir==1
isdirty= dirty==1
return
}



/*:65*/



/*68:*/


//line goacme.w:947

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



/*:68*/


