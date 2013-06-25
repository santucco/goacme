% This file is part of goacme package version 0.2
% Author Alexander Sychev

\def\title{goacme (version 0.2)}
\def\topofcontents{\null\vfill
	\centerline{\titlefont The {\ttitlefont goacme} package for manipulating {\ttitlefont plumb} messages}
	\vskip 15pt
	\centerline{(version 0.2)}
	\vfill}
\def\botofcontents{\vfill
\noindent
Copyright \copyright\ 2013 Alexander Sychev. All rights reserved.
\bigskip\noindent
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

\yskip\item{$\bullet$}Redistributions of source code must retain the 
above copyright
notice, this list of conditions and the following disclaimer.
\yskip\item{$\bullet$}Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
\yskip\item{$\bullet$}The name of author may not be used to endorse 
or promote products derived from
this software without specific prior written permission.

\bigskip\noindent
THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
}

\pageno=\contentspagenumber \advance\pageno by 1
\let\maybe=\iftrue

@** Introduction.
It is a package to manupulate windows of \.{Acme} 

@** The package \.{goacme}.
@c
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
@#
// Package goacme provides interface to acme programming environment
package goacme

import (
	@<Imports@>
)@#

type (
	@<Types@>
)@#

var (
	@<Variables@>
)@#

@<Constants@>@#


@ Let's describe a begin of a test for the package. The \.{Acme} will be be started for the test.

@(goacme_test.go@>=
package goacme

import (
	"os"
	"os/exec"
	"code.google.com/p/goplan9/plan9"
	"code.google.com/p/goplan9/plan9/client"
	"testing"
	@<Test specific imports@>
)@#

func prepare(t *testing.T)  {
	_,err:=client.MountService("acme")
	if err==nil {
		t.Log("acme started already")
	} else {
		cmd:=exec.Command("acme")
		err=cmd.Start()
		if err!=nil {
			t.Fatal(err)
		}
		@<Sleep a bit@>	
	}
}

@<Test routines@>

@ Let's describe |Window| structure. All the fields are unexported. 
For now |Window| contains |id| of a window, but the structure will be extended.
@<Types@>=
// |Window| is a structure to manipulate a particular acme's window.
Window struct {
	id		int
	@<|Window| struct members@>
}

@* New.

@
@<Imports@>=
"code.google.com/p/goplan9/plan9"
"code.google.com/p/goplan9/plan9/client"
"sync"
"fmt"

@ At first we have to mount \.{Acme} namespace
@<Variables@>=
fsys	*client.Fsys

@
@<Mount \.{Acme} namespace@>=
{
	var err error
	new(sync.Once).Do(func(){fsys,err=client.MountService("acme")})
	if err!=nil {
		return nil, err
	}
}

@ 
@c
// |New| creates a new window and returns |*Window| or |error|
func New() (*Window, error) {
	@<Mount \.{Acme} namespace@>
	f,err:=fsys.Open("new/ctl",plan9.OREAD)
	if err!=nil {
		return nil, err
	}
	defer f.Close()
	var id int
	if _,err:=fmt.Fscan(f, &id); err!=nil {
		return nil, err
	}
	return Open(id)
}

@* Open. |Open| opens a window with specified |id| and returns |*Window| or |error|
@c
// |Open| opens a window with a specified |id| and returns |*Window| or |error|
func Open(id int) (*Window, error) {
	@<Mount \.{Acme} namespace@>
	if err:=fsys.Access(fmt.Sprintf("%d", id), plan9.OREAD); err!=nil {
		return nil, err
	}
	this:=&Window{id: id}
	@<Init of |Window| members@>
	return this, nil
}

@* Close.
@c
// |Close| releases all resources of the window
func (this *Window) Close() error {
	@<Releasing of |Window| members@>
	return nil
}

@ Let's test |New| and |Open|

@<Test specific imports@>=
"fmt"
"time"

@	
@<Sleep a bit@>=
time.Sleep(time.Second)

@
@<Test routines@>=
func TestNewOpen(t *testing.T) {
	prepare(t)
	w,err:=New()
	if err!=nil {
		t.Fatal(err)
	}
	defer w.Close()
	defer w.Del(true)
	if f,err:=fsys.Open(fmt.Sprintf("%d", w.id), plan9.OREAD); err!=nil {
		t.Fatal(err)
	} else {
		f.Close()
	}
}

@* Read.
@c
// |Read| reads |len(p)| bytes from |"body"| of the window. 
// It returns count of readed bytes or |error|.
func (this *Window) Read(p []byte) (int, error) {
	f,err:=this.File("body")
	if err!=nil {
		return 0, err
	}
	return f.Read(p)
}

@* Write.
@c
// |Write| writes |len(p)| bytes to |"body"| of the window.
// It returns count of writen bytes or |error|.
func (this *Window) Write(p []byte) (int, error) {
	f,err:=this.File("body")
	if err!=nil {
		return 0, err
	}
	return f.Write(p)
}

@ Test of |Read| and |Write| function
@<Test specific imports@>=
"bytes"
"errors"

@
@<Test routines@>=
func TestReadWrite(t *testing.T) {
	w,err:=New()
	if err!=nil {
		t.Fatal(err)
	}
	defer w.Close()
	defer w.Del(true)
	b1:=[]byte("test")
	_,err=w.Write(b1)
	if err!=nil {
		t.Fatal(err)
	}
	w1,err:=Open(w.id)
	if err!=nil {
		t.Fatal(err)
	}
	defer w1.Close()
	defer w1.Del(true)
	b2:=make([]byte,10)
	n,err:=w1.Read(b2)
	if err!=nil {
		t.Fatal(err)
	}
	if bytes.Compare(b1,b2[:n])!=0 {
		t.Fatal(errors.New("buffers don't match"))
	}
}

@* Seek.
@c
// Seek sets the offset for the next Read or Write to offset, interpreted
// according to whence: 0 means relative to the origin of the file, 1 means
// relative to the current offset, and 2 means relative to the end. Seek
// returns the new offset and an Error, if any.
func (this *Window) Seek(offset int64, whence int) (ret int64, err error) {
	f,err:=this.File("body")
	if err!=nil {
		return 0, err
	}
	return f.Seek(offset, whence)
}



@* File.
@<Imports@>=
"io"

@ I have decided to store open files in |map[string]*os.File|.
@<|Window| struct members@>=
files map[string] *client.Fid

@
@<Init of |Window| members@>=
this.files=make(map[string] *client.Fid)

@ When |Window| is destroyed, all members of |files| have to be closed.
@<Releasing of |Window| members@>=
for _,v:=range this.files {
	v.Close()
}

@
@c
// |File| returns |io.ReadWriteSeeker| of corresponding |file| of the windows or error if |file| doesn't exist.
// Opened file will be cached inside the window, so caller shouldn't close it.
func (this *Window) File(file string) (io.ReadWriteSeeker, error) {
	f:=this.files[file]
	if f!=nil {
		return f, nil
	}
	f,err:=fsys.Open(fmt.Sprintf("%d/%s", this.id, file), plan9.ORDWR)	
	if err!=nil {
		return nil, err
	}
	this.files[file]=f
	return f, nil
}

@* PipeTo. 

@<Imports@>=
"os"
"os/exec"

@
@c
// |PipeTo| runs shell command line |cmd| with /dev/null on standard input and the window's file 
// on standard output.  If |stderr| is non-zero, it is used as standard error.
// Otherwise the command inherits the caller's standard error.
// |PipeTo| returns |*os.Process| of created process or |error|.
// Caller has to wait the running process and read its exit status.
func (this *Window) PipeTo(name string, stderr io.Writer, cmd ...string) (*os.Process, error) {
	if len(cmd)==0 {
		return nil, os.ErrInvalid
	}
	c:=exec.Command(cmd[0], cmd[1:]...)
	if stderr==nil {
		stderr=os.Stderr
	}
	c.Stdin=nil
	c.Stderr=stderr
	var err error
	if c.Stdout,err=this.File(name); err!=nil {
		return nil, err
	}
	if err:=c.Start(); err!=nil {
		return nil, err
	}
	return c.Process, nil	
}

@ Test of |PipeTo| function.
@<Test routines@>=
func TestPipeTo(t *testing.T) {
	w,err:=New()
	if err!=nil {
		t.Fatal(err)
	}
	defer w.Close()
	defer w.Del(true)
	s:="test"
	p,err:=w.PipeTo("body", nil, "echo", "-n", s)
	if err!=nil {
		t.Fatal(err)
	}
	p.Wait()
	p.Release()
	@<Sleep a bit@>
	w1,err:=Open(w.id)
	if err!=nil {
		t.Fatal(err)
	}
	defer w1.Close()
	defer w1.Del(true)
	b:=make([]byte,10)
	n,err:=w1.Read(b)
	if err!=nil {
		t.Fatal(err)
	}
	if bytes.Compare([]byte(s),b[:n])!=0 {
		t.Fatal(errors.New(fmt.Sprintf("buffers don't match: %q and %q", s,string(b))))
	}	
}

@* PipeFrom. 
@c
// |PipeFrom| runs a shell command line |cmd| with the
// window's file on standard input.  The command runs with
// |stdout| as its standard output and standard error.
// |PipeFrom| returns |*os.Process| of created process or |error|.
// Caller has to wait the running process and read its exit status.
func (this *Window) PipeFrom(name string, stdout io.Writer, cmd ...string) (*os.Process, error) {
	if len(cmd)==0 {
		return nil, os.ErrInvalid
	}
	c:=exec.Command(cmd[0], cmd[1:]...)
	if stdout==nil {
		stdout=os.Stdout
	}
	c.Stdout=stdout
	c.Stderr=os.Stderr
	var err error
	if c.Stdin,err=this.File(name); err!=nil {
		return nil, err
	} 
	if err:=c.Start(); err!=nil {
		return nil, err
	}
	return c.Process, nil	
}

@ Test of |PipeFrom| function.
@<Test routines@>=
func TestPipeFrom(t *testing.T) {
	w,err:=New()
	if err!=nil {
		t.Fatal(err)
	}
	s:="test"
	if _,err:=w.Write([]byte(s)); err!=nil {
		t.Fatal(err)
	}
	if _,err:=w.Seek(0,0); err!=nil {
		t.Fatal(err)
	}
	f,err:=os.OpenFile("/tmp/goacme.test",os.O_RDWR|os.O_TRUNC|os.O_CREATE,0600)
	if err!=nil {
		t.Fatal(err)
	}
	defer f.Close()
	p,err:=w.PipeFrom("body", f, "cat")
	if err!=nil {
		t.Fatal(err)
	}
	w.Del(true)
	w.Close()
	p.Wait()
	p.Release()
	@<Sleep a bit@>
	if _,err:=f.Seek(0,0); err!=nil {
		t.Fatal(err)
	}
	b:=make([]byte,10)
	n,err:=f.Read(b)
	if err!=nil {
		t.Fatal(err)
	}
	if bytes.Compare([]byte(s),b[:n])!=0 {
		t.Fatal(errors.New(fmt.Sprintf("buffers don't match: %q and %q", s,string(b))))
	}
}

@* SysRun.
@c
// |SysRun| runs shell command |cmd| and returns |io.ReadSeeker| of a window, 
// |*os.Process| of a running process or |error|.
// Caller has to wait the running process and read its exit status.
func SysRun(cmd ...string) (io.ReadSeeker, *os.Process, error) {
	w,err:=New()
	if err!=nil {
		return nil, nil, err
	}
	if len(cmd)==0 {
		w.Close()
		return nil, nil, os.ErrInvalid
	}
	c:=exec.Command(cmd[0], cmd[1:]...)
	f,err:=w.File("body")
	if err!=nil {
		return nil, nil, err
	}
	c.Stdout=f
	if err:=c.Start(); err!=nil {
		w.Close()
		return nil, nil, err
	}
	return f, c.Process, nil
}

@ Test of |SysRun| function.
@<Test routines@>=
func TestSysRun(t *testing.T) {
	s:="test"
	f,p,err:=SysRun("echo", "-n", s)
	if err!=nil {
		t.Fatal(err)
	}
	p.Wait()
	p.Release()
	@<Sleep a bit@>
	b:=make([]byte,10)
	if _,err:=f.Seek(0,0); err!=nil {
		t.Fatal(err)
	}
	n,err:=f.Read(b)
	if err!=nil {
		t.Fatal(err)
	}
	if bytes.Compare([]byte(s),b[:n])!=0 {
		t.Fatal(errors.New(fmt.Sprintf("buffers don't match: %q and %q", s,string(b))))
	}		
}


@* Del.
@c
// |Del| deletes the window, without a prompt if |sure| is true.
func (this *Window) Del(sure bool) error {
	f,err:=this.File("ctl")
	if err!=nil {
		return err
	}
	s:="del"
	if sure {
		s="delete"
	}
	_,err=f.Write([]byte(s))
	return err	
}

@ Test of |Del| function.
@<Test routines@>=
func TestDel(t *testing.T) {
	w,err:=New()
	if err!=nil {
		t.Fatal(err)
	}
	w.Del(true)
	w.Close()
	if _,err:=Open(w.id); err==nil {
		t.Fatal(errors.New(fmt.Sprintf("window %d is still opened", w.id)))
	}
}


@* DeleteAll. |DeleteAll| deletes all windows opened in a session. So all the windows should be stored
in a list. Some global variables and |Window| members are needed for this purpose.

@ |fwin| is a pointer to a first |Window| and |lwin| is a pointer to a last |Window|
@<Variables@>=
fwin	*Window
lwin	*Window

@ |prev| and |next| are pointer on previous |Window| and next |Window| respectively.
@<|Window| struct members@>=
prev	*Window
next	*Window

@ We need to place the window in the end of list of all windows
@<Init of |Window| members@>=
this.prev=lwin
this.next=nil
if fwin==nil {
	fwin=this
}
if lwin!=nil {
	lwin.next=this
}
lwin=this

@ When |Window| is destroyed, the |Window| has to be excluded from the list of windows
@<Releasing of |Window| members@>=
if this.next!=nil {
	this.next.prev=this.prev
}
if this.prev!=nil {
	this.prev.next=this.next
}
if fwin==this {
	fwin=this.next
}
if lwin==this {
	lwin=this.prev
}

@ Some trick is used to delete all |Window| - when |fwin| is closed, |fwin| is set to |fwin.next|,
so to delete all the windows |fwin| will be closed until |fwin| is not null.
@c
// |DeleteAll| deletes all the windows opened in a session
func DeleteAll() {
	for fwin!=nil {
		fwin.Del(true)
		fwin.Close()
	}
}

@ Test of |DeleteAll| function.
@<Test routines@>=
func TestDeleteAll(t *testing.T) {
	var l [10]int
	for i:=0;i<len(l);i++ {
		w,err:=New()
		if err!=nil {
			t.Fatal(err)
		}
		l[i]=w.id
	}
	DeleteAll()
	for _,v:=range l {
		_,err:=Open(v)
		if err==nil {
			t.Fatal(errors.New(fmt.Sprintf("window %d is still opened", v)))
		}
	}
}

@* Events processing.

@ At first let's describe |Event| structure. Fiels of |Event| will be specified a bit later.
@<Types@>=
Event struct {
	@<Fields of |Event|@>
}

@ |readEvent| reads properties of an event from |r|. Some trick is used here: |r| is supposed not buffered,
so it doesn't implement |RuneScanner| interface. When a length of text is parsing in event, 
a space followed by the length is read by |Fscanf| and we shouldn't read it. 
@c
func readFields(r io.Reader) (o rune, t rune, b int, e int, f int, s string, err error) {
	var l int
	if _,err=fmt.Fscanf(r, "%c%c%d %d %d %d", &o, &t, &b, &e, &f, &l); err!=nil {
		return
	}
	if l!=0 {
		rs:=make([]rune, l)
		for i:=0; i<l; i++ {
			if _,err=fmt.Fscanf(r,"%c", &rs[i]); err!=nil {
				return
			}
		}
		s=string(rs)
	}
	var nl [1]byte
	if _,err=r.Read(nl[:]); err!=nil {
		return
	}
	return
}

@ |readEvent| is unexported function to read |Event| from |f|. 
@c
func readEvent(r io.Reader) (*Event, error) {
	o, t, b, e, f, s, err:=readFields(r)
	if err!=nil {
		return nil, err
	}
	var ev Event
	@<Interpret origin@>
	@<Interpret action@>
	@<Fill addresses@>
	@<Interpret flag@>
	return &ev, nil	
}

@ Let's make a type for origin of an action
@<Types@>=
// |ActionOrigin| is a origin of the action
ActionOrigin	int

@ Here described variants of |ActionOrigin|
@<Constants@>=
const (
	// |Edit| is the origin for writes to the body or tag file
	Edit		ActionOrigin = 1 << iota
	// |File| is the origin for through the window's other files
	File
	// |Keyboard| is the origin for keyboard actions
	Keyboard	
	// |Mouse| is the origin for mouse actions
	Mouse
)

@ 
@<Fields of |Event|@>=
// |Origin| will be an origin of action with type |ActionOrigin|
Origin ActionOrigin

@
@<Imports@>=
"errors"

@ 
@<Variables@>=
// |ErrInvalidOrigin| will be returned if |Origin| is unexptected 
ErrInvalidOrigin=errors.New("invalid origin of action")

@
@<Interpret origin@>=
switch o {
	case 'E': ev.Origin=Edit
	case 'F': ev.Origin=File
	case 'K': ev.Origin=Keyboard
	case 'M': ev.Origin=Mouse
	default: return nil, ErrInvalidOrigin 
}

@ Let's make a type for type of an action
@<Types@>=
// |ActionType| is a type of the action
ActionType	int

@ Here described variants of |ActionType|
@<Constants@>=
const (
	//|Tag| is a flag points out event has occured in tag of window
	Tag		ActionType = 1
	Delete 	ActionType = 2<< iota
	Insert
	Look
	Execute
)

@
@<Fields of |Event|@>=
// |Type| will be an type of action with type |ActionType|
Type ActionType

@
@<Variables@>=
// |ErrInvalidType| will be returned if |Type| is unexpected 
ErrInvalidType=errors.New("invalid type of action")

@
@<Interpret action@>=
switch t {
	case 'D': ev.Type=Delete
	case 'd': ev.Type=Delete|Tag
	case 'I': ev.Type=Insert
	case 'i': ev.Type=Insert|Tag
	case 'L': ev.Type=Look
	case 'l': ev.Type=Look|Tag
	case 'X': ev.Type=Execute
	case 'x': ev.Type=Execute|Tag
	default: return nil, ErrInvalidType
}

@ |Begin| and |End| are addresses of the action.
|begin| and |end| are unexported addresses from an original event - they should be stored, 
but I decided to hide them to avoid collisions.

@<Fields of |Event|@>=
begin	int
// |Begin| is a start address of a text of the action
Begin	int
end		int
// |End| is an end address of the text of the action
End		int

@
@<Fill addresses@>=
ev.begin=b
ev.Begin=b
ev.end=e
ev.End=e

@ |flag| is an unexported copy of flag from an original event

@<Fields of |Event|@>=
flag		int
// |IsBuiltin| is a flag the action is recognised like a \.{Acme}'s builtin
IsBuiltin	bool
// |NoLoad| is a flag of acme can interpret the action without loading a new file
NoLoad		bool
// |IsFile| is a flag the |Text| is a file or window name
IsFile		bool
// |Text| is a text arguments of the action, perhaps with address
Text		string
// |Arg| is a text of chorded argument if any
Arg			string

@
@<Interpret flag@>=
ev.flag=f

if ev.Type&Execute==Execute {
	ev.IsBuiltin=(ev.flag&1)==1
} else if ev.Type&Look==Look {
	ev.NoLoad=(ev.flag&1)==1
	ev.IsFile=(ev.flag&4)==4
}

ev.Text=s

// if there is an expantion
if (ev.flag&2)==2 {
	_, _, ev.Begin, ev.End, _, ev.Text, err=readFields(r)
	if err!=nil {
		return nil, err
	}
}
// if there is a chording
if (ev.flag&8)==8 {
	_, _, _, _, _, ev.Arg, err=readFields(r)
	if err!=nil {
		return nil, err
	}
	_, _, _, _, _, _, err=readFields(r)
	if err!=nil {
		return nil, err
	}		
}

@*1 EventChannel. 
@<|Window| struct members@>=
ch	chan *Event

@
@c
// |EventChannel| returns a channel of |*Event| from which events can be read or |error|.
// First call of |EventChannel| starts a goroutine to read events from |"event"| file 
// and put them to the channel. Subsequent calls of |EventChannel| will return the same channel.
func (this *Window) EventChannel() (<-chan *Event, error) {
	if this.ch!=nil {
		return this.ch, nil
	}
	f,err:=this.File("event")
	if err!=nil	{
		return nil, err
	}
	this.ch=make(chan *Event)
	go func() {
		for ev, err:=readEvent(f); err==nil; ev, err=readEvent(f) {
			this.ch<-ev
		}
		close(this.ch)
		this.ch=nil
	} ()
	return this.ch, nil
}

@*1 ReadEvent.
@<Variables@>=
// |ErrChannelAlreadyOpened| will be returned 
// if channel of events is opened by call of |EventChannel| 
ErrChannelAlreadyOpened=errors.New("channel of events is already opened")

@
@c
// || reads an event from |"event"| file of the window and returns |*Event| or |error|
func (this *Window) ReadEvent() (*Event, error) {
	if this.ch!=nil {
		return nil, ErrChannelAlreadyOpened
	}
	f,err:=this.File("event")
	if err!=nil {
		return nil, err
	}
	return readEvent(f)
}


@*1 UnreadEvent.
@c
// |UnreadEvent| writes event |ev| back to the |"event"| file, 
// indicating to acme that it should be handled internally.
func (this *Window) UnreadEvent(ev *Event) error {
	f,err:=this.File("event")
	if err!=nil{
		return err
	}
	var o rune
	switch ev.Origin {
		case Edit: o='E'
		case File: o='F'
		case Keyboard: o='K'
		case Mouse: o='M'
		default: return ErrInvalidOrigin 
	}
	var t rune
	switch ev.Type {
		case Delete : t='D'
		case Delete&Tag: t='d'
		case Insert: t='I'
		case Insert|Tag: t='i'
		case Look: t='L'
		case Look|Tag: t='l'
		case Execute: t='X'
		case Execute|Tag: t='x'
		default: return ErrInvalidType
	}
	
	_,err=fmt.Fprintf(f,"%c%c%d %d \n", o, t, ev.begin, ev.end)
	return err
}

@ Tests for events
@<Test routines@>=
func TestEvent(t *testing.T) {
	w,err:=New()
	if err!=nil {
		t.Fatal(err)
	}
	defer w.Close()
	defer w.Del(true)
	msg:="Press left button of mouse on " 	
	test:="Test"
	if _,err:=w.Write([]byte(msg+test)); err!=nil {
		t.Fatal(err)
	}
	ch,err:=w.EventChannel() 
	if err!=nil {
		t.Fatal(err)
	}
	ok:=false
	var e *Event
	for e, ok=<-ch; ok; e,ok=<-ch {
		if e.Origin==Mouse {
			break
		}
	}
	if !ok {
		t.Fatal(errors.New("Channel is closed"))
	}
	if e.Origin!=Mouse || e.Type!=Look || e.Begin!=len(msg) || e.End!=len(msg)+len(test) || e.Text!=test {
		t.Fatal(errors.New(fmt.Sprintf("Something wrong with event: %#v",e)))
	}
	if _,err:=w.Write([]byte("\nPress middle button of mouse on Del in the window's tag")); err!=nil {
		t.Fatal(err)
	}
	for e, ok=<-ch; ok; e,ok=<-ch {
		if e.Origin==Mouse {
			break
		}
	}
	if !ok {
		t.Fatal(errors.New("Channel is closed"))
	}
	if e.Origin!=Mouse || e.Type!=(Execute|Tag) || e.Text!="Del" {
		t.Fatal(errors.New(fmt.Sprintf("Something wrong with event: %#v",e)))
	}
	if err:=w.UnreadEvent(e); err!=nil {
		t.Fatal(err)
	}
}

@** Index.
