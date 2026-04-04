package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Notebook</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.6;height:100vh;display:flex;flex-direction:column}
.hdr{padding:.7rem 1.2rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;flex-shrink:0}
.hdr h1{font-family:var(--mono);font-size:.85rem;letter-spacing:2px;text-transform:uppercase}
.hdr-r{display:flex;gap:.4rem;align-items:center}
.wrap{display:flex;flex:1;overflow:hidden}
.side{width:220px;border-right:1px solid var(--bg3);display:flex;flex-direction:column;flex-shrink:0;background:var(--bg2)}
.side-hdr{padding:.6rem .8rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.side-hdr span{font-family:var(--mono);font-size:.6rem;text-transform:uppercase;letter-spacing:1px;color:var(--cm)}
.nb-list{flex:1;overflow-y:auto;padding:.3rem 0}
.nb-item{padding:.4rem .8rem;cursor:pointer;font-size:.78rem;display:flex;align-items:center;gap:.5rem;border-left:3px solid transparent}
.nb-item:hover{background:var(--bg3)}.nb-item.active{background:var(--bg3);border-left-color:var(--rust)}
.nb-dot{width:8px;height:8px;border-radius:50%;flex-shrink:0}
.nb-name{flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.nb-count{font-family:var(--mono);font-size:.55rem;color:var(--cm)}
.nb-special{color:var(--cm);font-style:italic}
.main{flex:1;display:flex;flex-direction:column;overflow:hidden}
.toolbar{padding:.5rem .8rem;border-bottom:1px solid var(--bg3);display:flex;gap:.4rem;align-items:center;flex-wrap:wrap}
.search{background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem;padding:.3rem .5rem;flex:1;min-width:120px}
.search:focus{outline:none;border-color:var(--leather)}
.stats{display:flex;gap:.3rem;padding:.4rem .8rem;border-bottom:1px solid var(--bg3)}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.3rem .6rem;text-align:center;font-family:var(--mono);flex:1}
.st-v{font-size:1rem;color:var(--cream)}.st-l{font-size:.45rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.05rem}
.notes{flex:1;overflow-y:auto;padding:.5rem .8rem}
.note{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem .9rem;margin-bottom:.4rem;cursor:pointer;transition:border-color .15s}
.note:hover{border-color:var(--leather)}
.note.pinned{border-left:3px solid var(--gold)}
.note-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}
.note-title{font-size:.88rem;flex:1}.note-title.untitled{color:var(--cm);font-style:italic}
.note-acts{display:flex;gap:.2rem;opacity:0;transition:opacity .15s}.note:hover .note-acts{opacity:1}
.note-preview{font-size:.72rem;color:var(--cm);margin-top:.2rem;display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical;overflow:hidden}
.note-meta{font-family:var(--mono);font-size:.55rem;color:var(--cm);margin-top:.3rem;display:flex;gap:.5rem;align-items:center;flex-wrap:wrap}
.tag{font-family:var(--mono);font-size:.5rem;padding:.1rem .3rem;background:var(--bg3);color:var(--cd);border-radius:2px}
.pin-icon{color:var(--gold);font-size:.65rem}
.btn{font-family:var(--mono);font-size:.6rem;padding:.2rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);white-space:nowrap}
.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:var(--bg)}.btn-p:hover{background:#d06830}
.btn-icon{background:none;border:none;color:var(--cm);cursor:pointer;font-size:.7rem;padding:.1rem .2rem}
.btn-icon:hover{color:var(--cream)}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.2rem;width:560px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-family:var(--mono);font-size:.75rem;margin-bottom:.8rem;color:var(--rust);text-transform:uppercase;letter-spacing:1px}
.fr{margin-bottom:.5rem}
.fr label{display:block;font-family:var(--mono);font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.1rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.35rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.fr textarea{font-family:var(--mono);line-height:1.5;resize:vertical}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:.8rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.8rem}
.color-row{display:flex;gap:.3rem;margin-top:.2rem}
.color-dot{width:20px;height:20px;border-radius:50%;cursor:pointer;border:2px solid transparent}
.color-dot.sel{border-color:var(--cream)}
.sort-sel{background:var(--bg2);border:1px solid var(--bg3);color:var(--cm);font-family:var(--mono);font-size:.6rem;padding:.25rem .4rem}
.wc{color:var(--leather);font-size:.5rem}
@media(max-width:640px){.side{width:180px}.stats{flex-wrap:wrap}.st{min-width:60px}}
</style></head><body>
<div class="hdr">
<h1>NOTEBOOK</h1>
<div class="hdr-r">
<button class="btn" onclick="exportAll()">Export All</button>
<button class="btn btn-p" onclick="openNoteForm()">+ Note</button>
</div>
</div>
<div class="stats" id="stats"></div>
<div class="wrap">
<div class="side">
<div class="side-hdr"><span>Notebooks</span><button class="btn-icon" onclick="openNbForm()" title="New Notebook">+</button></div>
<div class="nb-list" id="nbList"></div>
</div>
<div class="main">
<div class="toolbar">
<input class="search" id="search" placeholder="Search notes..." oninput="debounceSearch()">
<select class="sort-sel" id="sortBy" onchange="load()">
<option value="updated">Last Updated</option>
<option value="created">Created</option>
<option value="title">Title</option>
</select>
</div>
<div class="notes" id="notes"></div>
</div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)cm()"><div class="modal" id="mdl"></div></div>
<script>
const A='/api';let notes=[],notebooks=[],tags=[],curNb='',curTag='',showArchive=false,searchTimer;
const COLORS=['#e8753a','#d4a843','#4a9e5c','#5b8ed4','#c94444','#a0845c','#8e6cc4','#c45d8a'];

async function load(){
const[nr,nbr,tr,sr]=await Promise.all([
fetch(A+'/notes?'+new URLSearchParams({notebook_id:curNb,tag:curTag,search:document.getElementById('search').value,sort:document.getElementById('sortBy').value,archived:showArchive?'true':'false',limit:'200'})).then(r=>r.json()),
fetch(A+'/notebooks').then(r=>r.json()),
fetch(A+'/tags').then(r=>r.json()),
fetch(A+'/stats').then(r=>r.json())
]);
notes=nr.notes||[];notebooks=nbr.notebooks||[];tags=tr.tags||[];
renderStats(sr);renderSidebar();renderNotes();
}

function renderStats(s){
document.getElementById('stats').innerHTML=
'<div class="st"><div class="st-v">'+s.notes+'</div><div class="st-l">Notes</div></div>'+
'<div class="st"><div class="st-v">'+s.notebooks+'</div><div class="st-l">Notebooks</div></div>'+
'<div class="st"><div class="st-v">'+s.pinned+'</div><div class="st-l">Pinned</div></div>'+
'<div class="st"><div class="st-v">'+(s.words||0).toLocaleString()+'</div><div class="st-l">Words</div></div>'+
'<div class="st"><div class="st-v">'+s.tags+'</div><div class="st-l">Tags</div></div>';
}

function renderSidebar(){
let h='<div class="nb-item'+(curNb===''&&!curTag&&!showArchive?' active':'')+'" onclick="selNb(\\'\\')"><span class="nb-special">All Notes</span></div>';
h+='<div class="nb-item'+(showArchive?' active':'')+'" onclick="toggleArchive()"><span class="nb-special">Archive</span></div>';
notebooks.forEach(nb=>{
h+='<div class="nb-item'+(curNb===nb.id?' active':'')+'" onclick="selNb(\''+nb.id+'\')">';
h+='<span class="nb-dot" style="background:'+esc(nb.color||'#e8753a')+'"></span>';
h+='<span class="nb-name">'+esc(nb.name)+'</span>';
h+='<span class="nb-count">'+nb.note_count+'</span>';
h+='</div>';
});
if(tags.length){
h+='<div style="padding:.4rem .8rem;margin-top:.3rem;border-top:1px solid var(--bg3)"><span style="font-family:var(--mono);font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px">Tags</span></div>';
tags.sort().forEach(t=>{
h+='<div class="nb-item'+(curTag===t?' active':'')+'" onclick="selTag(\''+esc(t)+'\')"><span class="tag">'+esc(t)+'</span></div>';
});
}
document.getElementById('nbList').innerHTML=h;
}

function renderNotes(){
if(!notes.length){document.getElementById('notes').innerHTML='<div class="empty">'+(showArchive?'No archived notes.':'No notes yet. Create one!')+'</div>';return;}
let h='';notes.forEach(n=>{
const preview=n.body?n.body.substring(0,180).replace(/\n/g,' '):'';
const nbObj=n.notebook_id?notebooks.find(nb=>nb.id===n.notebook_id):null;
h+='<div class="note'+(n.pinned?' pinned':'')+'" ondblclick="openEditor(\''+n.id+'\')">';
h+='<div class="note-top"><div class="note-title'+(n.title?'':' untitled')+'">'+(n.pinned?'<span class="pin-icon">&#x1F4CC;</span> ':'')+(n.title?esc(n.title):'Untitled')+'</div>';
h+='<div class="note-acts">';
if(!n.archived){
h+='<button class="btn-icon" onclick="event.stopPropagation();togglePin(\''+n.id+'\','+!n.pinned+')" title="'+(n.pinned?'Unpin':'Pin')+'">'+(n.pinned?'&#x1F4CC;':'&#x1F4CD;')+'</button>';
h+='<button class="btn-icon" onclick="event.stopPropagation();openEditor(\''+n.id+'\')" title="Edit">&#x270E;</button>';
h+='<button class="btn-icon" onclick="event.stopPropagation();archiveNote(\''+n.id+'\',true)" title="Archive">&#x1F4E6;</button>';
} else {
h+='<button class="btn-icon" onclick="event.stopPropagation();archiveNote(\''+n.id+'\',false)" title="Unarchive">&#x21A9;</button>';
}
h+='<button class="btn-icon" onclick="event.stopPropagation();del(\''+n.id+'\')" title="Delete" style="color:var(--red)">&#x2715;</button>';
h+='</div></div>';
if(preview)h+='<div class="note-preview">'+esc(preview)+'</div>';
h+='<div class="note-meta">';
if(nbObj)h+='<span style="color:'+esc(nbObj.color||'var(--leather)')+'">'+esc(nbObj.name)+'</span>';
h+='<span>'+timeAgo(n.updated_at)+'</span>';
if(n.word_count)h+='<span class="wc">'+n.word_count+' words</span>';
if(n.tags)n.tags.forEach(t=>h+='<span class="tag">'+esc(t)+'</span>');
h+='</div></div>';
});
document.getElementById('notes').innerHTML=h;
}

function selNb(id){curNb=id;curTag='';showArchive=false;load();}
function selTag(t){curTag=t;curNb='';showArchive=false;load();}
function toggleArchive(){showArchive=!showArchive;curNb='';curTag='';load();}
function debounceSearch(){clearTimeout(searchTimer);searchTimer=setTimeout(load,300);}

async function togglePin(id,pin){await fetch(A+'/notes/'+id+'/'+(pin?'pin':'unpin'),{method:'POST'});load();}
async function archiveNote(id,archive){await fetch(A+'/notes/'+id+'/'+(archive?'archive':'unarchive'),{method:'POST'});load();}
async function del(id){if(confirm('Delete this note?')){await fetch(A+'/notes/'+id,{method:'DELETE'});load();}}

function openNoteForm(note){
const n=note||{};
let nbOpts='<option value="">None</option>';
notebooks.forEach(nb=>nbOpts+='<option value="'+nb.id+'"'+(n.notebook_id===nb.id?' selected':'')+'>'+esc(nb.name)+'</option>');
document.getElementById('mdl').innerHTML='<h2>'+(n.id?'Edit':'New')+' Note</h2>'+
'<div class="fr"><label>Title</label><input id="f-title" value="'+esc(n.title||'')+'"></div>'+
'<div class="fr"><label>Notebook</label><select id="f-nb">'+nbOpts+'</select></div>'+
'<div class="fr"><label>Content</label><textarea id="f-body" rows="12">'+esc(n.body||'')+'</textarea></div>'+
'<div class="fr"><label>Tags (comma separated)</label><input id="f-tags" value="'+esc((n.tags||[]).join(', '))+'"></div>'+
'<div class="acts">'+(n.id?'<button class="btn" onclick="exportNote(\''+n.id+'\')" style="margin-right:auto">Export</button>':'')+
'<button class="btn" onclick="cm()">Cancel</button><button class="btn btn-p" onclick="saveNote('+(n.id?"'"+n.id+"'":"null")+')">Save</button></div>';
document.getElementById('mbg').classList.add('open');
setTimeout(function(){document.getElementById('f-title').focus();},100);
}

async function openEditor(id){
const n=await fetch(A+'/notes/'+id).then(r=>r.json());
openNoteForm(n);
}

async function saveNote(id){
const tags=document.getElementById('f-tags').value.split(',').map(function(t){return t.trim();}).filter(Boolean);
const body={title:document.getElementById('f-title').value||'Untitled',body:document.getElementById('f-body').value,notebook_id:document.getElementById('f-nb').value,tags:tags};
if(id){await fetch(A+'/notes/'+id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/notes',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
cm();load();
}

function openNbForm(nb){
const n=nb||{};
let colorH='<div class="color-row">';
COLORS.forEach(function(c){colorH+='<div class="color-dot'+(c===(n.color||'#e8753a')?' sel':'')+'" style="background:'+c+'" onclick="pickColor(this,\''+c+'\')"></div>';});
colorH+='</div>';
document.getElementById('mdl').innerHTML='<h2>'+(n.id?'Edit':'New')+' Notebook</h2>'+
'<div class="fr"><label>Name</label><input id="f-nbname" value="'+esc(n.name||'')+'"></div>'+
'<div class="fr"><label>Color</label><input type="hidden" id="f-nbcolor" value="'+(n.color||'#e8753a')+'">'+colorH+'</div>'+
'<div class="acts">'+
(n.id?'<button class="btn" onclick="delNb(\''+n.id+'\')" style="margin-right:auto;color:var(--red)">Delete</button>':'')+
'<button class="btn" onclick="cm()">Cancel</button><button class="btn btn-p" onclick="saveNb('+(n.id?"'"+n.id+"'":"null")+')">Save</button></div>';
document.getElementById('mbg').classList.add('open');
setTimeout(function(){document.getElementById('f-nbname').focus();},100);
}

function pickColor(el,c){document.querySelectorAll('.color-dot').forEach(function(d){d.classList.remove('sel');});el.classList.add('sel');document.getElementById('f-nbcolor').value=c;}

async function saveNb(id){
const body={name:document.getElementById('f-nbname').value,color:document.getElementById('f-nbcolor').value};
if(!body.name){alert('Name required');return;}
if(id){await fetch(A+'/notebooks/'+id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/notebooks',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
cm();load();
}

async function delNb(id){if(confirm('Delete notebook? Notes will be unassigned.')){await fetch(A+'/notebooks/'+id,{method:'DELETE'});cm();curNb='';load();}}

function exportNote(id){window.open(A+'/notes/'+id+'/export','_blank');}
function exportAll(){window.open(A+'/export','_blank');}

function cm(){document.getElementById('mbg').classList.remove('open');}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
function timeAgo(d){if(!d)return'';var s=Math.floor((Date.now()-new Date(d).getTime())/1000);if(s<60)return'just now';if(s<3600)return Math.floor(s/60)+'m ago';if(s<86400)return Math.floor(s/3600)+'h ago';if(s<604800)return Math.floor(s/86400)+'d ago';return new Date(d).toLocaleDateString();}
document.addEventListener('keydown',function(e){if(e.key==='Escape')cm();if(e.ctrlKey&&e.key==='n'){e.preventDefault();openNoteForm();}});
load();
</script></body></html>`
