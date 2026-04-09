package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Notebook</title>
<link href="https://fonts.googleapis.com/css2?family=Libre+Baskerville:ital,wght@0,400;0,700;1,400&family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.5;font-size:14px;overflow:hidden}
.app{display:grid;grid-template-columns:220px 320px 1fr;height:100vh}
@media(max-width:900px){.app{grid-template-columns:1fr}.sidebar,.notes-col{display:none}.app.show-sidebar .sidebar{display:flex}.app.show-notes .notes-col{display:flex}}

.sidebar{background:var(--bg2);border-right:1px solid var(--bg3);display:flex;flex-direction:column;overflow-y:auto}
.sidebar-hdr{padding:.9rem 1rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.sidebar-hdr h1{font-family:var(--mono);font-size:.75rem;letter-spacing:2px;color:var(--cream)}
.sidebar-hdr h1 span{color:var(--rust)}
.sidebar-section{padding:.6rem 1rem;font-family:var(--mono);font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;display:flex;justify-content:space-between;align-items:center;margin-top:.4rem}
.sidebar-section .add{cursor:pointer;color:var(--cm);background:none;border:none;font-family:var(--mono);font-size:.7rem;padding:0 .2rem}
.sidebar-section .add:hover{color:var(--rust)}
.nb-item{padding:.4rem 1rem;font-family:var(--mono);font-size:.7rem;cursor:pointer;display:flex;justify-content:space-between;align-items:center;gap:.4rem;border-left:2px solid transparent;transition:.15s}
.nb-item:hover{background:var(--bg3)}
.nb-item.active{border-left-color:var(--rust);background:var(--bg3);color:var(--cream)}
.nb-name{display:flex;align-items:center;gap:.4rem;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;flex:1;min-width:0}
.nb-dot{width:8px;height:8px;border-radius:50%;flex-shrink:0}
.nb-count{font-size:.55rem;color:var(--cm);flex-shrink:0}
.nb-edit{display:none;font-size:.55rem;color:var(--cm);background:none;border:none;cursor:pointer}
.nb-item:hover .nb-edit{display:inline}
.nb-edit:hover{color:var(--cream)}

.notes-col{background:var(--bg);border-right:1px solid var(--bg3);display:flex;flex-direction:column;overflow:hidden}
.notes-hdr{padding:.7rem 1rem;border-bottom:1px solid var(--bg3);display:flex;flex-direction:column;gap:.5rem}
.notes-hdr-row{display:flex;justify-content:space-between;align-items:center}
.notes-hdr-title{font-family:var(--mono);font-size:.65rem;color:var(--cd);text-transform:uppercase;letter-spacing:1px}
.notes-search{padding:.35rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.notes-search:focus{outline:none;border-color:var(--leather)}
.notes-list{flex:1;overflow-y:auto}
.note-item{padding:.7rem 1rem;border-bottom:1px solid var(--bg3);cursor:pointer;display:flex;flex-direction:column;gap:.2rem;transition:.15s;border-left:3px solid transparent}
.note-item:hover{background:var(--bg2)}
.note-item.active{background:var(--bg2);border-left-color:var(--rust)}
.note-item.archived{opacity:.55}
.note-item-title{font-family:var(--serif);font-size:.85rem;color:var(--cream);font-weight:700;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;display:flex;align-items:center;gap:.3rem}
.note-pin{color:var(--gold);font-size:.65rem}
.note-item-snippet{font-size:.65rem;color:var(--cm);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.note-item-meta{font-family:var(--mono);font-size:.5rem;color:var(--cm);display:flex;gap:.4rem;flex-wrap:wrap;margin-top:.15rem}
.note-tag{background:var(--bg3);color:var(--cd);padding:.05rem .3rem}

.editor-col{background:var(--bg);display:flex;flex-direction:column;overflow:hidden}
.editor-hdr{padding:.7rem 1rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:.5rem;flex-wrap:wrap}
.editor-actions{display:flex;gap:.3rem;flex-wrap:wrap}
.editor-body{flex:1;display:flex;flex-direction:column;padding:1rem 1.5rem;overflow-y:auto}
.editor-title{width:100%;background:none;border:none;color:var(--cream);font-family:var(--serif);font-size:1.4rem;font-weight:700;margin-bottom:.5rem;padding:.2rem 0}
.editor-title:focus{outline:none;border-bottom:1px solid var(--leather)}
.editor-meta{font-family:var(--mono);font-size:.55rem;color:var(--cm);margin-bottom:.8rem;display:flex;gap:.6rem;flex-wrap:wrap}
.editor-tags-row{display:flex;gap:.4rem;align-items:center;margin-bottom:.7rem;flex-wrap:wrap}
.editor-tags-input{background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.6rem;padding:.3rem .5rem;flex:1;min-width:120px}
.editor-tags-input:focus{outline:none;border-color:var(--leather)}
.editor-text{flex:1;width:100%;background:none;border:none;color:var(--cream);font-family:var(--serif);font-size:.95rem;line-height:1.7;resize:none;min-height:300px}
.editor-text:focus{outline:none}
.editor-extras{margin-top:1rem;padding-top:1rem;border-top:1px solid var(--bg3)}
.editor-extras-label{font-family:var(--mono);font-size:.55rem;color:var(--rust);text-transform:uppercase;letter-spacing:1px;margin-bottom:.5rem}

.btn{font-family:var(--mono);font-size:.6rem;padding:.3rem .55rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:.15s}
.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-p:hover{opacity:.85;color:#fff}
.btn-icon{padding:.3rem .45rem;font-size:.7rem}
.btn-del{color:var(--red);border-color:#3a1a1a}
.btn-del:hover{border-color:var(--red);color:var(--red)}

.empty{padding:3rem 2rem;text-align:center;color:var(--cm);font-style:italic;font-size:.85rem}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:440px;max-width:92vw}
.modal h2{font-family:var(--serif);font-size:1.05rem;margin-bottom:1rem;color:var(--rust)}
.fr{margin-bottom:.7rem}
.fr label{display:block;font-family:var(--mono);font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input[type=color]{height:34px;padding:.1rem;cursor:pointer}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.acts .btn-del{margin-right:auto}

.stats-bar{padding:.5rem 1rem;border-top:1px solid var(--bg3);font-family:var(--mono);font-size:.5rem;color:var(--cm);display:flex;gap:.7rem;flex-wrap:wrap;background:var(--bg2)}
.stats-bar span strong{color:var(--cd);font-weight:700}
.trial-bar{display:none;background:linear-gradient(90deg,#3a2419,#2e1c14);border-bottom:2px solid var(--rust);padding:.7rem 1.5rem;font-family:var(--mono);font-size:.68rem;color:var(--cream);align-items:center;gap:1rem;flex-wrap:wrap;position:sticky;top:0;z-index:100}
.trial-bar.show{display:flex}
.trial-bar-msg{flex:1;min-width:240px;line-height:1.5}
.trial-bar-msg strong{color:var(--rust);text-transform:uppercase;letter-spacing:1px;font-size:.6rem;display:block;margin-bottom:.15rem}
.trial-bar-actions{display:flex;gap:.5rem;align-items:center;flex-wrap:wrap}
.trial-bar a.btn-trial{background:var(--rust);color:#fff;padding:.4rem .8rem;text-decoration:none;font-size:.65rem;text-transform:uppercase;letter-spacing:1px;font-weight:700;border:1px solid var(--rust);transition:all .2s}
.trial-bar a.btn-trial:hover{background:#f08545;border-color:#f08545}
.trial-bar-divider{color:var(--cm);font-size:.6rem}
.trial-bar input.key-input{padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.6rem;width:200px}
.trial-bar input.key-input:focus{outline:none;border-color:var(--rust)}
.trial-bar button.btn-activate{padding:.4rem .7rem;background:var(--bg2);color:var(--cream);border:1px solid var(--leather);font-family:var(--mono);font-size:.6rem;cursor:pointer;text-transform:uppercase;letter-spacing:1px}
.trial-bar button.btn-activate:hover{background:var(--bg3)}
.trial-bar button.btn-activate:disabled{opacity:.5;cursor:wait}
.trial-msg{font-size:.6rem;color:var(--cm);margin-left:.5rem}
.trial-msg.error{color:#e74c3c}
.trial-msg.success{color:#4ade80}
.btn-disabled-trial{opacity:.45;cursor:not-allowed!important}
@media(max-width:600px){.trial-bar{flex-direction:column;align-items:stretch}.trial-bar input.key-input{width:100%}}
</style>
</head>
<body>

<div class="trial-bar" id="trial-bar">
<div class="trial-bar-msg">
<strong>Trial Required</strong>
You can read your existing notes, but creating or editing is locked until you start a 14-day free trial.
</div>
<div class="trial-bar-actions">
<a class="btn-trial" href="https://stockyard.dev/" target="_blank" rel="noopener">Start 14-Day Trial</a>
<span class="trial-bar-divider">or</span>
<input type="text" class="key-input" id="trial-key-input" placeholder="SY-..." autocomplete="off" spellcheck="false">
<button class="btn-activate" id="trial-activate-btn" onclick="activateLicense()">Activate</button>
<span class="trial-msg" id="trial-msg"></span>
</div>
</div>

<div class="app" id="app">

<aside class="sidebar">
<div class="sidebar-hdr">
<h1 id="dash-title"><span>&#9670;</span> NOTEBOOK</h1>
</div>
<div class="sidebar-section">All Notes <button class="add" onclick="selectNotebook('')" title="Show all">All</button></div>
<div id="notebook-list"></div>
<div class="sidebar-section">Notebooks <button class="add" onclick="openNotebookForm()" title="New notebook">+</button></div>
<div id="notebooks"></div>
<div class="sidebar-section">Tags</div>
<div id="tags"></div>
<div class="stats-bar" id="stats"></div>
</aside>

<section class="notes-col">
<div class="notes-hdr">
<div class="notes-hdr-row">
<div class="notes-hdr-title" id="notes-col-title">All Notes</div>
<button class="btn btn-p btn-icon" onclick="newNote()">+ Note</button>
</div>
<input class="notes-search" id="search" placeholder="Search notes..." oninput="debouncedReload()">
<div class="notes-hdr-row">
<select class="notes-search" style="flex:1" id="archived-filter" onchange="reload()">
<option value="">Active</option>
<option value="true">Archived</option>
<option value="all">All</option>
</select>
<select class="notes-search" style="flex:1" id="sort-filter" onchange="reload()">
<option value="updated">Updated</option>
<option value="created">Created</option>
<option value="title">Title</option>
</select>
</div>
</div>
<div class="notes-list" id="notes-list"></div>
</section>

<section class="editor-col">
<div class="editor-hdr">
<div class="editor-actions" id="editor-actions"></div>
</div>
<div class="editor-body" id="editor-body">
<div class="empty">Select a note or click "+ Note" to start writing.</div>
</div>
</section>

</div>

<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()">
<div class="modal" id="mdl"></div>
</div>

<script>
var A='/api';
var notebooks=[],notes=[],tags=[],currentNote=null,currentNotebookID='',searchTimer=null,saveTimer=null;
var notebookExtras={},noteExtras={};
var noteCustomFields=[]; // injected from /api/config note_custom_fields
var notebookCustomFields=[]; // injected from /api/config notebook_custom_fields
var dirty=false;

// ─── Loading ──────────────────────────────────────────────────────

async function loadAll(){
try{
var resps=await Promise.all([
fetch(A+'/notebooks').then(function(r){return r.json()}),
fetch(A+'/tags').then(function(r){return r.json()}),
fetch(A+'/stats').then(function(r){return r.json()}),
fetch(A+'/extras/notebooks').then(function(r){return r.json()}),
fetch(A+'/extras/notes').then(function(r){return r.json()})
]);
notebooks=resps[0].notebooks||[];
tags=resps[1].tags||[];
notebookExtras=resps[3]||{};
noteExtras=resps[4]||{};
renderSidebar(resps[2]||{});
}catch(e){
console.error('loadAll failed',e);
}
await reload();
}

async function reload(){
var q={
search:document.getElementById('search').value,
archived:document.getElementById('archived-filter').value,
sort:document.getElementById('sort-filter').value
};
if(currentNotebookID)q.notebook_id=currentNotebookID;
var qs=Object.keys(q).filter(function(k){return q[k]}).map(function(k){return k+'='+encodeURIComponent(q[k])}).join('&');
try{
var r=await fetch(A+'/notes'+(qs?'?'+qs:'')).then(function(r){return r.json()});
notes=r.notes||[];
notes.forEach(function(n){
var x=noteExtras[n.id];
if(!x)return;
Object.keys(x).forEach(function(k){if(n[k]===undefined)n[k]=x[k]});
});
}catch(e){
notes=[];
}
renderNotesList();
}

function debouncedReload(){
clearTimeout(searchTimer);
searchTimer=setTimeout(reload,200);
}

function renderSidebar(stats){
var sect=document.getElementById('notebooks');
var html='';
notebooks.forEach(function(nb){
var cls='nb-item'+(nb.id===currentNotebookID?' active':'');
html+='<div class="'+cls+'" onclick="selectNotebook(\''+esc(nb.id)+'\')">';
html+='<div class="nb-name"><span class="nb-dot" style="background:'+esc(nb.color||'#c45d2c')+'"></span>'+esc(nb.name)+'</div>';
html+='<span class="nb-count">'+(nb.note_count||0)+'</span>';
html+='<button class="nb-edit" onclick="event.stopPropagation();openNotebookForm(\''+esc(nb.id)+'\')">edit</button>';
html+='</div>';
});
if(!notebooks.length)html='<div style="padding:.5rem 1rem;font-size:.6rem;color:var(--cm);font-style:italic">No notebooks yet</div>';
sect.innerHTML=html;

var allItem=document.getElementById('notebook-list');
allItem.innerHTML='<div class="nb-item'+(currentNotebookID===''?' active':'')+'" onclick="selectNotebook(\'\')"><div class="nb-name">All Notes</div><span class="nb-count">'+(stats.notes||0)+'</span></div>';

var tagsEl=document.getElementById('tags');
if(tags.length){
tagsEl.innerHTML=tags.slice(0,20).map(function(t){return'<div class="nb-item" onclick="filterTag(\''+esc(t)+'\')"><div class="nb-name">#'+esc(t)+'</div></div>'}).join('');
}else{
tagsEl.innerHTML='<div style="padding:.5rem 1rem;font-size:.6rem;color:var(--cm);font-style:italic">No tags yet</div>';
}

document.getElementById('stats').innerHTML=
'<span><strong>'+(stats.notes||0)+'</strong> notes</span>'+
'<span><strong>'+(stats.notebooks||0)+'</strong> notebooks</span>'+
'<span><strong>'+(stats.words||0)+'</strong> words</span>'+
'<span><strong>'+(stats.pinned||0)+'</strong> pinned</span>';
}

function renderNotesList(){
var titleEl=document.getElementById('notes-col-title');
var nb=null;
if(currentNotebookID){
for(var i=0;i<notebooks.length;i++)if(notebooks[i].id===currentNotebookID){nb=notebooks[i];break}
}
titleEl.textContent=nb?nb.name:'All Notes';

var listEl=document.getElementById('notes-list');
if(!notes.length){
listEl.innerHTML='<div class="empty">'+(window._emptyMsg||'No notes here yet')+'</div>';
return;
}
var h='';
notes.forEach(function(n){
var cls='note-item'+(currentNote&&currentNote.id===n.id?' active':'')+(n.archived?' archived':'');
h+='<div class="'+cls+'" onclick="selectNote(\''+esc(n.id)+'\')">';
h+='<div class="note-item-title">';
if(n.pinned)h+='<span class="note-pin">&#9733;</span>';
h+=esc(n.title||'(untitled)');
h+='</div>';
var snip=String(n.body||'').replace(/\s+/g,' ').slice(0,80);
if(snip)h+='<div class="note-item-snippet">'+esc(snip)+'</div>';
h+='<div class="note-item-meta">';
h+='<span>'+fmtDate(n.updated_at)+'</span>';
if(n.word_count)h+='<span>'+n.word_count+'w</span>';
if(n.tags&&n.tags.length){
n.tags.slice(0,3).forEach(function(t){h+='<span class="note-tag">#'+esc(t)+'</span>'});
}
h+='</div>';
h+='</div>';
});
listEl.innerHTML=h;
}

// ─── Note editor ──────────────────────────────────────────────────

function selectNotebook(id){
currentNotebookID=id;
renderSidebar({notes:notes.length,notebooks:notebooks.length});
reload();
}

function filterTag(tag){
document.getElementById('search').value='';
// We can't filter by tag without server support — fallback to client-side
reload().then(function(){
notes=notes.filter(function(n){return n.tags&&n.tags.indexOf(tag)>=0});
renderNotesList();
});
}

async function selectNote(id){
if(dirty&&!confirm('Discard unsaved changes?'))return;
var n=null;
for(var i=0;i<notes.length;i++)if(notes[i].id===id){n=notes[i];break}
if(!n)return;
currentNote=n;
dirty=false;
renderEditor();
renderNotesList();
}

function newNote(){
if(dirty&&!confirm('Discard unsaved changes?'))return;
currentNote={
id:'',
title:'',
body:'',
tags:[],
pinned:false,
archived:false,
notebook_id:currentNotebookID||'',
created_at:'',
updated_at:'',
word_count:0
};
dirty=false;
renderEditor();
setTimeout(function(){var t=document.getElementById('e-title');if(t)t.focus()},50);
}

function renderEditor(){
var n=currentNote;
if(!n){
document.getElementById('editor-body').innerHTML='<div class="empty">Select a note or click "+ Note" to start writing.</div>';
document.getElementById('editor-actions').innerHTML='';
return;
}

// Action buttons
var acts='';
if(n.id){
acts+='<button class="btn btn-icon" onclick="togglePin()" title="Pin">'+(n.pinned?'&#9733; Unpin':'&#9734; Pin')+'</button>';
acts+='<button class="btn btn-icon" onclick="toggleArchive()">'+(n.archived?'Unarchive':'Archive')+'</button>';
acts+='<button class="btn btn-icon" onclick="exportNote()">Export</button>';
acts+='<button class="btn btn-icon btn-del" onclick="deleteCurrent()">Delete</button>';
}
acts+='<button class="btn btn-p btn-icon" onclick="saveNote()">Save</button>';
document.getElementById('editor-actions').innerHTML=acts;

// Body
var nbName='Inbox';
for(var i=0;i<notebooks.length;i++)if(notebooks[i].id===n.notebook_id){nbName=notebooks[i].name;break}

var h='';
h+='<input class="editor-title" id="e-title" placeholder="Untitled" value="'+esc(n.title||'')+'" oninput="markDirty()">';
h+='<div class="editor-meta">';
h+='<span>'+esc(nbName)+'</span>';
if(n.updated_at)h+='<span>Updated '+fmtDate(n.updated_at)+'</span>';
if(n.word_count)h+='<span>'+n.word_count+' words</span>';
h+='</div>';

// Tags input
h+='<div class="editor-tags-row">';
h+='<input class="editor-tags-input" id="e-tags" placeholder="tags, comma separated" value="'+esc((n.tags||[]).join(', '))+'" oninput="markDirty()">';

// Notebook selector
h+='<select class="editor-tags-input" id="e-notebook" oninput="markDirty()" style="flex:0 0 130px">';
h+='<option value=""'+(!n.notebook_id?' selected':'')+'>Inbox</option>';
notebooks.forEach(function(nb){
h+='<option value="'+esc(nb.id)+'"'+(nb.id===n.notebook_id?' selected':'')+'>'+esc(nb.name)+'</option>';
});
h+='</select>';
h+='</div>';

h+='<textarea class="editor-text" id="e-body" placeholder="Start writing..." oninput="markDirty()">'+esc(n.body||'')+'</textarea>';

// Custom fields
if(noteCustomFields.length){
h+='<div class="editor-extras"><div class="editor-extras-label">'+esc(window._customSectionLabel||'Additional Details')+'</div>';
noteCustomFields.forEach(function(f){
var v=n[f.name];
h+=customFieldHTML(f,v);
});
h+='</div>';
}

document.getElementById('editor-body').innerHTML=h;
}

function customFieldHTML(f,value){
var v=value;
if(v===undefined||v===null)v='';
var h='<div class="fr"><label>'+esc(f.label)+'</label>';
if(f.type==='textarea'){
h+='<textarea id="cf-'+f.name+'" rows="2" oninput="markDirty()">'+esc(String(v))+'</textarea>';
}else if(f.type==='select'){
h+='<select id="cf-'+f.name+'" oninput="markDirty()"><option value="">Select...</option>';
(f.options||[]).forEach(function(o){
var sel=String(v)===String(o)?' selected':'';
h+='<option value="'+esc(String(o))+'"'+sel+'>'+esc(String(o))+'</option>';
});
h+='</select>';
}else if(f.type==='number'){
h+='<input type="number" id="cf-'+f.name+'" value="'+esc(String(v))+'" oninput="markDirty()">';
}else{
h+='<input type="text" id="cf-'+f.name+'" value="'+esc(String(v))+'" oninput="markDirty()">';
}
h+='</div>';
return h;
}

function markDirty(){dirty=true}

async function saveNote(){
if(!currentNote)return;
var titleEl=document.getElementById('e-title');
var bodyEl=document.getElementById('e-body');
var tagsEl=document.getElementById('e-tags');
var nbEl=document.getElementById('e-notebook');

var title=titleEl?titleEl.value.trim():'';
if(!title){alert('Title is required');return}

var tagList=(tagsEl?tagsEl.value:'').split(',').map(function(t){return t.trim()}).filter(function(t){return t});
var body={
title:title,
body:bodyEl?bodyEl.value:'',
tags:tagList,
notebook_id:nbEl?nbEl.value:''
};

var extras={};
noteCustomFields.forEach(function(f){
var el=document.getElementById('cf-'+f.name);
if(!el)return;
var v;
if(f.type==='number')v=parseFloat(el.value)||0;
else v=el.value.trim();
extras[f.name]=v;
});

var savedId=currentNote.id;
try{
if(currentNote.id){
var r1=await fetch(A+'/notes/'+currentNote.id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r1.ok){var e1=await r1.json().catch(function(){return{}});alert(e1.error||'Save failed');return}
currentNote=await r1.json();
}else{
var r2=await fetch(A+'/notes',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r2.ok){var e2=await r2.json().catch(function(){return{}});alert(e2.error||'Save failed');return}
currentNote=await r2.json();
savedId=currentNote.id;
}
if(savedId&&Object.keys(extras).length){
await fetch(A+'/extras/notes/'+savedId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(extras)}).catch(function(){});
// Merge extras into currentNote
Object.keys(extras).forEach(function(k){currentNote[k]=extras[k]});
noteExtras[savedId]=extras;
}
}catch(e){
alert('Network error: '+e.message);
return;
}
dirty=false;
await loadAll();
// Re-select the saved note
for(var i=0;i<notes.length;i++)if(notes[i].id===savedId){currentNote=notes[i];break}
renderEditor();
}

async function togglePin(){
if(!currentNote||!currentNote.id)return;
var endpoint=currentNote.pinned?'unpin':'pin';
var r=await fetch(A+'/notes/'+currentNote.id+'/'+endpoint,{method:'POST'});
if(r.ok){
currentNote=await r.json();
await loadAll();
for(var i=0;i<notes.length;i++)if(notes[i].id===currentNote.id){currentNote=notes[i];break}
renderEditor();
}
}

async function toggleArchive(){
if(!currentNote||!currentNote.id)return;
var endpoint=currentNote.archived?'unarchive':'archive';
var r=await fetch(A+'/notes/'+currentNote.id+'/'+endpoint,{method:'POST'});
if(r.ok){
currentNote=await r.json();
await loadAll();
for(var i=0;i<notes.length;i++)if(notes[i].id===currentNote.id){currentNote=notes[i];break}
renderEditor();
}
}

function exportNote(){
if(!currentNote||!currentNote.id)return;
window.open(A+'/notes/'+currentNote.id+'/export','_blank');
}

async function deleteCurrent(){
if(!currentNote||!currentNote.id)return;
if(!confirm('Delete this note?'))return;
await fetch(A+'/notes/'+currentNote.id,{method:'DELETE'});
currentNote=null;
dirty=false;
await loadAll();
renderEditor();
}

// ─── Notebook modal ───────────────────────────────────────────────

function openNotebookForm(id){
var nb=null;
if(id){for(var i=0;i<notebooks.length;i++)if(notebooks[i].id===id){nb=notebooks[i];break}}
var isEdit=!!nb;
var h='<h2>'+(isEdit?'Edit Notebook':'New Notebook')+'</h2>';
h+='<div class="fr"><label>Name *</label><input id="nb-name" value="'+esc(nb?nb.name:'')+'"></div>';
h+='<div class="fr"><label>Slug</label><input id="nb-slug" value="'+esc(nb?nb.slug:'')+'" placeholder="auto-generated"></div>';
h+='<div class="fr"><label>Color</label><input type="color" id="nb-color" value="'+esc(nb?nb.color:'#c45d2c')+'"></div>';

if(notebookCustomFields.length){
var ext=notebookExtras[id]||{};
h+='<div class="editor-extras"><div class="editor-extras-label">'+esc(window._notebookCustomLabel||'Notebook Details')+'</div>';
notebookCustomFields.forEach(function(f){
h+=customFieldHTML(f,ext[f.name]);
});
h+='</div>';
}

h+='<div class="acts">';
if(isEdit)h+='<button class="btn btn-del" onclick="deleteNotebook(\''+esc(id)+'\')">Delete</button>';
h+='<button class="btn" onclick="closeModal()">Cancel</button>';
h+='<button class="btn btn-p" onclick="saveNotebook(\''+(id||'')+'\')">'+(isEdit?'Save':'Create')+'</button>';
h+='</div>';

document.getElementById('mdl').innerHTML=h;
document.getElementById('mbg').classList.add('open');
setTimeout(function(){var n=document.getElementById('nb-name');if(n)n.focus()},50);
}

async function saveNotebook(id){
var name=document.getElementById('nb-name').value.trim();
if(!name){alert('Name is required');return}
var body={
name:name,
slug:document.getElementById('nb-slug').value.trim(),
color:document.getElementById('nb-color').value
};
var extras={};
notebookCustomFields.forEach(function(f){
var el=document.getElementById('cf-'+f.name);
if(!el)return;
var v;
if(f.type==='number')v=parseFloat(el.value)||0;
else v=el.value.trim();
extras[f.name]=v;
});

var savedId=id;
try{
if(id){
var r1=await fetch(A+'/notebooks/'+id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r1.ok){var e1=await r1.json().catch(function(){return{}});alert(e1.error||'Save failed');return}
}else{
var r2=await fetch(A+'/notebooks',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r2.ok){var e2=await r2.json().catch(function(){return{}});alert(e2.error||'Create failed');return}
var created=await r2.json();
savedId=created.id;
}
if(savedId&&Object.keys(extras).length){
await fetch(A+'/extras/notebooks/'+savedId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(extras)}).catch(function(){});
}
}catch(e){
alert('Network error: '+e.message);
return;
}
closeModal();
await loadAll();
}

async function deleteNotebook(id){
if(!confirm('Delete this notebook? Notes will be moved to Inbox.'))return;
await fetch(A+'/notebooks/'+id,{method:'DELETE'});
if(currentNotebookID===id)currentNotebookID='';
closeModal();
await loadAll();
}

function closeModal(){
document.getElementById('mbg').classList.remove('open');
}

// ─── Helpers ──────────────────────────────────────────────────────

function fmtDate(s){
if(!s)return'';
try{
var d=new Date(s);
if(isNaN(d.getTime()))return s;
var now=new Date();
var diffDays=Math.floor((now-d)/(1000*60*60*24));
if(diffDays===0)return'Today';
if(diffDays===1)return'Yesterday';
if(diffDays<7)return diffDays+'d ago';
return d.toLocaleDateString('en-US',{month:'short',day:'numeric',year:'numeric'});
}catch(e){return s}
}

function esc(s){
if(s===undefined||s===null)return'';
var d=document.createElement('div');
d.textContent=String(s);
return d.innerHTML;
}

document.addEventListener('keydown',function(e){
if(e.key==='Escape')closeModal();
if((e.metaKey||e.ctrlKey)&&e.key==='s'){e.preventDefault();saveNote()}
});

// ─── Personalization ──────────────────────────────────────────────

(function loadPersonalization(){
fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
if(!cfg||typeof cfg!=='object')return;

if(cfg.dashboard_title){
var h1=document.getElementById('dash-title');
if(h1)h1.innerHTML='<span>&#9670;</span> '+esc(cfg.dashboard_title);
document.title=cfg.dashboard_title;
}

if(cfg.empty_state_message)window._emptyMsg=cfg.empty_state_message;
if(cfg.note_section_label)window._customSectionLabel=cfg.note_section_label;
if(cfg.notebook_section_label)window._notebookCustomLabel=cfg.notebook_section_label;

if(Array.isArray(cfg.note_custom_fields)){
noteCustomFields=cfg.note_custom_fields.filter(function(f){return f&&f.name&&f.label});
}
if(Array.isArray(cfg.notebook_custom_fields)){
notebookCustomFields=cfg.notebook_custom_fields.filter(function(f){return f&&f.name&&f.label});
}
}).catch(function(){
}).finally(function(){
checkTrialState();
loadAll();
});
})();

// ─── trial-required license gating ───
window._trialRequired=false;

async function checkTrialState(){
try{
var resp=await fetch('/api/tier');
if(!resp.ok)return;
var data=await resp.json();
window._trialRequired=!!data.trial_required;
if(window._trialRequired){
document.getElementById('trial-bar').classList.add('show');
disableWriteControls();
}else{
document.getElementById('trial-bar').classList.remove('show');
}
}catch(e){}
}

function disableWriteControls(){
// Notebook has a three-pane layout with + buttons in the sidebar and
// in the notes column header. Neutralize both, plus the note editor
// save button if it's currently rendered.
var selectors=['.sidebar .add','.notes-hdr .btn-p','.note-editor .btn-p','.acts .btn-p'];
selectors.forEach(function(sel){
document.querySelectorAll(sel).forEach(function(b){
b.classList.add('btn-disabled-trial');
b.title='Locked: trial required';
b.onclick=function(e){
e.preventDefault();
showTrialNudge();
return false;
};
});
});
}

function showTrialNudge(){
var input=document.getElementById('trial-key-input');
if(input){
input.focus();
input.style.borderColor='var(--rust)';
setTimeout(function(){if(input)input.style.borderColor=''},1500);
}
}

async function activateLicense(){
var input=document.getElementById('trial-key-input');
var btn=document.getElementById('trial-activate-btn');
var msg=document.getElementById('trial-msg');
if(!input||!btn||!msg)return;
var key=(input.value||'').trim();
if(!key){
msg.className='trial-msg error';
msg.textContent='Paste your license key first';
input.focus();
return;
}
btn.disabled=true;
msg.className='trial-msg';
msg.textContent='Activating...';
try{
var resp=await fetch('/api/license/activate',{
method:'POST',
headers:{'Content-Type':'application/json'},
body:JSON.stringify({license_key:key})
});
var data=await resp.json();
if(!resp.ok){
msg.className='trial-msg error';
msg.textContent=data.error||'Activation failed';
btn.disabled=false;
return;
}
msg.className='trial-msg success';
msg.textContent='Activated. Reloading...';
setTimeout(function(){location.reload()},800);
}catch(e){
msg.className='trial-msg error';
msg.textContent='Network error: '+e.message;
btn.disabled=false;
}
}

document.addEventListener('DOMContentLoaded',function(){
var input=document.getElementById('trial-key-input');
if(input){
input.addEventListener('keydown',function(e){
if(e.key==='Enter')activateLicense();
});
}
});
</script>
</body>
</html>`
