% This file is part of goacme package version 0.62
% Author Alexander Sychev
\def\ver{0.62}
\def\title{goacme (version \ver)}
\def\topofcontents{\null\vfill
	\centerline{\titlefont The {\ttitlefont goacme} package for manipulating {\ttitlefont plumb} messages}
	\vskip 15pt
	\centerline{(version \ver)}
	\vfill}
\def\botofcontents{\vfill
\noindent
Copyright \copyright\ 2013, 2014 Alexander Sychev. All rights reserved.
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

@ Legal information.
@c
// Copyright (c) 2013, 2014 Alexander Sychev. All rights reserved.
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

@** Implementation.
@c
// Package |goacme| provides interface to |acme| programming environment
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


@ Let's describe a begin of a test for the package. \.{Acme} will be started for the test.

@(goacme_test.go@>=
package goacme

import (
	"os/exec"
	"9fans.net/go/plan9/client"
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
// |Window| is a structure to manipulate a particular |acme|'s window.
Window struct {
	id		int
	@<|Window| struct members@>
}

@* New.
@<Imports@>=
"9fans.net/go/plan9"
"9fans.net/go/plan9/client"
"sync"
"fmt"

@ At first we have to mount \.{Acme} namespace
@<Variables@>=
fsys	*client.Fsys
once	sync.Once

@
@<Mount \.{Acme} namespace@>=
{
	var err error
	once.Do(func(){fsys,err=client.MountService("acme")})
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

@* Open.
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
"9fans.net/go/plan9"


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
// |Read| reads |len(p)| bytes from |"body"| file of the window.
// |Read| returns a count of read bytes or |error|.
func (this *Window) Read(p []byte) (int, error) {
	f,err:=this.File("body")
	if err!=nil {
		return 0, err
	}
	return f.Read(p)
}

@* Write.
@c
// |Write| writes |len(p)| bytes to |"body"| file of the window.
// |Write| returns a count of written bytes or |error|.
func (this *Window) Write(p []byte) (int, error) {
	f,err:=this.File("body")
	if err!=nil {
		return 0, err
	}
	@<Convert |f| to a wrapper@>
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
// |Seek| sets a position for the next Read or Write to |offset|, interpreted
// according to whence: 0 means relative to the origin of the file, 1 means
// relative to the current offset, and 2 means relative to the end.
// |Seek| returns the new offset or |error|
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
// |File| returns |io.ReadWriteSeeker| of corresponding |file| of the windows or |error|
func (this *Window) File(file string) (io.ReadWriteSeeker, error) {
	fid, ok:=this.files[file]
	if !ok {
		var err error
		if fid,err=fsys.Open(fmt.Sprintf("%d/%s", this.id, file), plan9.ORDWR); err!=nil {
			if fid,err=fsys.Open(fmt.Sprintf("%d/%s", this.id, file), plan9.OREAD); err!=nil {
				if fid,err=fsys.Open(fmt.Sprintf("%d/%s", this.id, file), plan9.OWRITE); err!=nil {
					return nil, err
				}
			}
		}
		this.files[file]=fid
	}
	var f io.ReadWriteSeeker = fid
	@<Convert |f| to a wrapper@>
	return f, nil	
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


@* DeleteAll.
|DeleteAll| deletes all windows opened in a session. So all the windows should be stored
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

@ |readFields| reads properties of an event from |r|. Some trick is used here: |r| is supposed not buffered,
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

@ Here we describe variants of |ActionOrigin|
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
// |ErrInvalidOrigin| will be returned if a case of an unexpected origin of action
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

@ Here we describe variants of |ActionType|
@<Constants@>=
const (
	Delete 	ActionType = 1<< iota
	Insert
	Look
	Execute
	// |Tag| is a flag points out the event has occured in the tag of the window
	Tag
	// |TagMask| is a mask points out the event should be masked by tag
	TagMask
)

@
@<Fields of |Event|@>=
// |Type| will be an type of action with type |ActionType|
Type ActionType

@
@<Variables@>=
// |ErrInvalidType| will be returned if a case of an unexpected type of action
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
// |IsBuiltin| is a flag the action is recognised like an \.{Acme}'s builtin
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

// if there is an expansion
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
@<Check if some arguments are in |Text| field@>

@ If some arguments are in |Text|, then let's add them in the begin of |Arg|
@<Imports@>=
"strings"

@
@<Check if some arguments are in |Text| field@>=
if len(ev.Text)>0 {
	f:=strings.Fields(ev.Text)
	if len(f)>1 {
		ev.Text=f[0]
		s:=ev.Arg
		if len(s)>0 {
			s=" "+ev.Arg
		}
		ev.Arg=strings.Join(f[1:], " ")+s
	}
}


@*1 EventChannel.
@<|Window| struct members@>=
ch	chan *Event

@
@c
// |EventChannel| returns a channel of |*Event| with a buffer |size|
// from which events can be read or |error|.
// Only |ActionType|s set in |tmask| are used.
// If |TagMask| is set in |tmask|, the event will be masked by tag. Otherwise |Tag| flag will be ignored.
// First call of |EventChannel| starts a goroutine to read events from |"event"| file
// and put them to the channel. Subsequent calls of |EventChannel| will return the same channel.
func (this *Window) EventChannel(size int, tmask ActionType) (<-chan *Event, error) {
	if this.ch!=nil {
		return this.ch, nil
	}
	@<Trying to restrict events by type@>
	f,err:=this.File("event")
	if err!=nil	{
		return nil, err
	}
	if tmask&TagMask!=TagMask {
		tmask|=Tag
	}
	this.ch=make(chan *Event, size)
	go func() {
		for ev, err:=readEvent(f); err==nil; ev, err=readEvent(f) {
			if old && ev.Type&tmask!=ev.Type {
				if ev.Type&Insert!=Insert && ev.Type&Delete!=Delete {
					this.UnreadEvent(ev)
				}
				continue
			}
			this.ch<-ev
		}
		close(this.ch)
		this.ch=nil
	} ()
	return this.ch, nil
}


@ Two kinds of filtiring of events are implemented. If \.{Acme} has a support of events restriction,
|old| is false and we do not check events because of \.{Acme} does it. Otherwise we check type
of events.
@<Trying to restrict events by type@>=
old:=false
{
	var em string
	if tmask&Delete==Delete {
		em+="D"
	}
	if tmask&Insert==Insert {
		em+="I"
	}
	if tmask&Look==Look {
		em+="L"
	}
	if tmask&Execute==Execute {
		em+="X"
	}
	if tmask&TagMask!=TagMask {
		em+=strings.ToLower(em)
	}
	if err:=this.WriteCtl("events %s\n", em); err!=nil {
		old=true
	}
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
		case Delete|Tag: t='d'
		case Insert: t='I'
		case Insert|Tag: t='i'
		case Look: t='L'
		case Look|Tag: t='l'
		case Execute: t='X'
		case Execute|Tag: t='x'
		default: return ErrInvalidType
	}

	_,err=fmt.Fprintf(f,"%c%c%d %d\n", o, t, ev.begin, ev.end)
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
	ch,err:=w.EventChannel(0, Look|Execute)
	if err!=nil {
		t.Fatal(err)
	}
	e, ok:=<-ch
	if !ok {
		t.Fatal(errors.New("Channel is closed"))
	}
	if e.Origin!=Mouse || e.Type!=Look || e.Begin!=len(msg) || e.End!=len(msg)+len(test) || e.Text!=test {
		t.Fatal(errors.New(fmt.Sprintf("Something wrong with event: %#v",e)))
	}
	if _,err:=w.Write([]byte("\nChording test: select argument, press middle button of mouse on Execute and press left button of mouse without releasing middle button")); err!=nil {
		t.Fatal(err)
	}
	e, ok=<-ch
	if !ok {
		t.Fatal(errors.New("Channel is closed"))
	}
	if e.Origin!=Mouse || e.Type!=(Execute) || e.Text!="Execute" || e.Arg!="argument" {
		t.Fatal(errors.New(fmt.Sprintf("Something wrong with event: %#v",e)))
	}
	if err:=w.UnreadEvent(e); err!=nil {
		t.Fatal(err)
	}
	if _,err:=w.Write([]byte("\nPress middle button of mouse on Del in the window's tag")); err!=nil {
		t.Fatal(err)
	}
	e, ok=<-ch
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

@* WriteAddr.
@c
// |WriteAddr| writes |format| with |args| in |"addr"| file of the window
func (this *Window) WriteAddr(format string, args ...interface{}) error {
	f,err:=this.File("addr")
	if err!=nil{
		return err
	}
	if len(args)>0 {
		format=fmt.Sprintf(format,args...)
	}
	_,err=f.Write([]byte(format))
	return err
}

@* ReadAddr.
@c
// |ReadAddr| reads the address of the next read/write operation from |"addr"| file of the window.
// |ReadAddr| return |begin| and |end| offsets in symbols or |error|
func (this *Window) ReadAddr() (begin int, end int, err error) {
	f,err:=this.File("addr")
	if err!=nil{
		return
	}
	if _,err=f.Seek(0,0); err!=nil {
		return
	}
	_,err=fmt.Fscanf(f, "%d %d", &begin, &end)
	return
}

@ We should have |"addr"| file is opened because \.{Acme} clears internal address range when |"addr"| is being opened.
@<Init of |Window| members@>=
if _,err:=this.File("addr"); err!=nil {
	return nil, err
}


@ Tests for operations with addresses
@<Test routines@>=
func TestWriteReadAddr(t *testing.T) {
	w,err:=New()
	if err!=nil {
		t.Fatal(err)
	}
	defer w.Close()
	defer w.Del(true)
	if b,e,err:=w.ReadAddr(); err!=nil {
		t.Fatal(err)
	} else if b!=0 || e!=0 {
		t.Fatal(errors.New(fmt.Sprintf("Something wrong with address: %v, %v", b, e)))
	}
	if _,err:=w.Write([]byte("test")); err!=nil {
		t.Fatal(err)
	}
	if err:=w.WriteAddr("0,$"); err!=nil {
		t.Fatal(err)
	}
	if b,e,err:=w.ReadAddr(); err!=nil {
		t.Fatal(err)
	} else if b!=0 || e!=4 {
		t.Fatal(errors.New(fmt.Sprintf("Something wrong with address: %v, %v", b, e)))
	}
}

@* WriteCtl.
@c
// |WriteCtl| writes |format| with |args| in |"ctl"| file of the window
// In case |format| is not ended by newline, |'\n'| will be added to the end of |format|
func (this *Window) WriteCtl(format string, args ...interface{}) error {
	f,err:=this.File("ctl")
	if err!=nil{
		return err
	}
	if len(args)>0 {
		format=fmt.Sprintf(format,args...)
	}
	if len(format)>=0 && format[len(format)-1]!='\n' {
		format+="\n"
	}
	if _,err=f.Seek(0,0); err!=nil {
		return err
	}
	_,err=f.Write([]byte(format))
	return err
}


@* ReadCtl.
@c
// |ReadCtl| reads the address of the next read/write operation from |"ctl"| file of the window.
// |ReadCtl| returns:
//    |id| - the window ID
//    |tlen| - number of characters (runes) in the tag;
//    |blen| - number of characters in the body;
//    |isdir| -  |true| if the window is a directory, |false| otherwise;
//    |isdirty| - |true| if the window is modified, |false|otherwise;
//    |wwidth| - the width of the window in pixels;
//    |font| - the name of the font used in the window;
//    |twidth| - the width of a tab character in pixels;
//    |error| - in case of any error.
func (this *Window) ReadCtl() (id int, tlen int, blen int, isdir bool, isdirty bool, wwidth int, font string, twidth int, err error) {
	f,err:=this.File("ctl")
	if err!=nil{
		return
	}
	if _,err=f.Seek(0,0); err!=nil {
		return
	}
	var dir,dirty int
	_,err=fmt.Fscanf(f, "%d %d %d %d %d %d %s %d", &id, &tlen, &blen, &dir, &dirty, &wwidth, &font, &twidth)
	isdir=dir==1
	isdirty=dirty==1
	return
}

@ Tests for operations with |"ctl"| file
@<Test routines@>=
func TestWriteReadCtl(t *testing.T) {
	w,err:=New()
	if err!=nil {
		t.Fatal(err)
	}
	defer w.Close()
	defer w.Del(true)
	if _,err:=w.Write([]byte("test")); err!=nil {
		t.Fatal(err)
	}
	if _,_,_,_,d,_,_,_,err:=w.ReadCtl(); err!=nil {
		t.Fatal(err)
	} else if !d  {
		t.Fatal(errors.New(fmt.Sprintf("The window has to be dirty\n")))
	}
	if err:=w.WriteCtl("clean"); err!=nil {
		t.Fatal(err)
	}
	if _,_,_,_,d,_,_,_,err:=w.ReadCtl(); err!=nil {
		t.Fatal(err)
	} else if d  {
		t.Fatal(errors.New(fmt.Sprintf("The window has to be clean\n")))
	}
}

@ I found \.{Acme} panics when a size of message is more that 8168 bytes.
So I decided to make a wrapper to replace |Write| method.
@<Types@>=
wrapper struct {
	f io.ReadWriteSeeker
}

@ |wrapper| has to support |io.ReadWriteSeeker| interface, so here are the interface functions.
@c
func (this *wrapper) Read(p []byte) (int, error) {
	return this.f.Read(p)
}

func (this *wrapper) Write(p []byte) (int, error) {
	if len(p)<8168 {
		return this.f.Write(p)
	}
	c:=0
	for i:=0; i<len(p); i+=8168 {
		n:=i+8168
		if n>len(p) {
			n=len(p)
		}
		n,e:=this.f.Write(p[i:n])
		c+=n
		if e!=nil {
			return c,e
		}
	}
	return c, nil
}

func (this *wrapper) Seek(offset int64, whence int) (ret int64, err error) {
	return this.f.Seek(offset, whence)
}

@ This is a convertor to |wrapper|
@<Convert |f| to a wrapper@>=
f=&wrapper{f:f}



@** Index.
