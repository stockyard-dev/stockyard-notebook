package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Notebook</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#c45d2c;--rl:#e8753a;--leather:#a0845c;--ll:#c4a87a;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c44040;--mono:'JetBrains Mono',Consolas,monospace;--serif:'Libre Baskerville',Georgia,serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);font-size:13px;line-height:1.6;height:100vh;overflow:hidden}
a{color:var(--rl);text-decoration:none}a:hover{color:var(--gold)}
.app{display:flex;height:100vh}
.sidebar{width:220px;background:var(--bg2);border-right:1px solid var(--bg3);display:flex;flex-direction:column;flex-shrink:0}
.sidebar-hdr{padding:.6rem .8rem;border-bottom:1px solid var(--bg3);font-family:var(--serif);font-size:.9rem}
.sidebar-hdr span{color:var(--rl)}
.sidebar-section{padding:.4rem .8rem;font-size:.6rem;text-transform:uppercase;letter-spacing:1.5px;color:var(--rust);margin-top:.5rem}
.sidebar-item{padding:.3rem .8rem;font-size:.75rem;cursor:pointer;display:flex;align-items:center;gap:.4rem;transition:background .1s;color:var(--cd)}
.sidebar-item:hover{background:var(--bg3)}.sidebar-item.active{background:var(--bg3);color:var(--cream)}
.sidebar-dot{width:8px;height:8px;border-radius:50%;flex-shrink:0}
.sidebar-count{margin-left:auto;font-size:.6rem;color:var(--cm)}
.sidebar-bottom{margin-top:auto;padding:.5rem .8rem;border-top:1px solid var(--bg3);font-size:.65rem;color:var(--cm)}

.list-pane{width:280px;border-right:1px solid var(--bg3);display:flex;flex-direction:column;flex-shrink:0}
.list-toolbar{padding:.4rem .6rem;border-bottom:1px solid var(--bg3);display:flex;gap:.4rem;align-items:center}
.list-toolbar input{flex:1;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);padding:.3rem .5rem;font-family:var(--mono);font-size:.72rem;outline:none}
.list-toolbar input:focus{border-color:var(--rust)}
.list-scroll{flex:1;overflow-y:auto}
.note-item{padding:.5rem .7rem;border-bottom:1px solid var(--bg3);cursor:pointer;transition:background .1s}
.note-item:hover{background:var(--bg2)}.note-item.active{background:var(--bg2);border-left:2px solid var(--rl)}
.note-item-title{font-size:.78rem;font-weight:600;display:flex;align-items:center;gap:.3rem}
.note-item-title .pin{color:var(--gold);font-size:.6rem}
.note-item-preview{font-size:.68rem;color:var(--cm);margin-top:.15rem;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.note-item-meta{font-size:.6rem;color:var(--cm);margin-top:.15rem;display:flex;gap:.5rem}
.tag-chip{font-size:.55rem;padding:0 .25rem;background:var(--bg3);color:var(--ll);border-radius:2px}

.editor-pane{flex:1;display:flex;flex-direction:column;min-width:0}
.editor-toolbar{padding:.4rem .8rem;border-bottom:1px solid var(--bg3);display:flex;align-items:center;gap:.5rem}
.editor-toolbar select{background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.72rem;padding:.25rem .4rem}
.btn{font-family:var(--mono);font-size:.68rem;padding:.25rem .6rem;border:1px solid;cursor:pointer;background:transparent;transition:.15s;white-space:nowrap}
.btn-p{border-color:var(--rust);color:var(--rl)}.btn-p:hover{background:var(--rust);color:var(--cream)}
.btn-d{border-color:var(--bg3);color:var(--cm)}.btn-d:hover{border-color:var(--red);color:var(--red)}
.btn-s{border-color:var(--green);color:var(--green)}.btn-s:hover{background:var(--green);color:var(--bg)}
.btn-g{border-color:var(--gold);color:var(--gold)}.btn-g:hover{background:var(--gold);color:var(--bg)}
.editor-title{width:100%;background:transparent;border:none;color:var(--cream);font-family:var(--serif);font-size:1.1rem;padding:.6rem .8rem;outline:none;border-bottom:1px solid var(--bg3)}
.editor-tags{padding:.3rem .8rem;border-bottom:1px solid var(--bg3)}
.editor-tags input{background:transparent;border:none;color:var(--ll);font-family:var(--mono);font-size:.7rem;outline:none;width:100%}
.editor-body{flex:1;display:flex;overflow:hidden}
.editor-body textarea{flex:1;background:transparent;border:none;color:var(--cd);font-family:var(--mono);font-size:.8rem;padding:.8rem;outline:none;resize:none;line-height:1.7}
.preview{flex:1;padding:.8rem;overflow-y:auto;border-left:1px solid var(--bg3);font-size:.8rem;color:var(--cd);line-height:1.7;display:none}
.preview h1,.preview h2,.preview h3{color:var(--cream);font-family:var(--serif);margin:1rem 0 .4rem}
.preview h1{font-size:1.2rem}.preview h2{font-size:1rem}.preview h3{font-size:.9rem}
.preview p{margin:.4rem 0}.preview code{background:var(--bg3);padding:.1rem .3rem;font-size:.75rem;border-radius:2px}
.preview pre{background:var(--bg3);padding:.6rem;margin:.5rem 0;overflow-x:auto;font-size:.75rem;border-radius:2px}
.preview pre code{background:transparent;padding:0}
.preview ul,.preview ol{padding-left:1.2rem;margin:.4rem 0}
.preview blockquote{border-left:3px solid var(--rust);padding-left:.8rem;color:var(--cm);margin:.5rem 0}
.preview a{color:var(--rl)}.preview strong{color:var(--cream)}.preview em{color:var(--ll)}

.empty{text-align:center;padding:3rem 1rem;color:var(--cm);font-style:italic;font-family:var(--serif)}

.modal-bg{position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,.65);display:flex;align-items:center;justify-content:center;z-index:100}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:90%;max-width:400px}
.modal h2{font-family:var(--serif);font-size:.9rem;margin-bottom:.8rem}
label.fl{display:block;font-size:.65rem;color:var(--leather);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem;margin-top:.5rem}
input[type=text],input[type=color]{background:var(--bg);border:1px solid var(--bg3);color:var(--cream);padding:.35rem .5rem;font-family:var(--mono);font-size:.78rem;width:100%;outline:none}
input:focus{border-color:var(--rust)}
</style>
<link href="https://fonts.googleapis.com/css2?family=Libre+Baskerville:ital@0;1&family=JetBrains+Mono:wght@400;600&display=swap" rel="stylesheet">
</head><body>
<div class="app">
<div class="sidebar">
<div class="sidebar-hdr"><span>Notebook</span></div>
<div class="sidebar-section">Notebooks</div>
<div class="sidebar-item active" onclick="filterNotebook('')" data-nb="all">All notes <span class="sidebar-count" id="sAll">-</span></div>
<div id="nbList"></div>
<div class="sidebar-item" style="color:var(--rl)" onclick="showNewNotebook()">+ New notebook</div>
<div class="sidebar-section">Tags</div>
<div id="tagList"></div>
<div class="sidebar-section" style="margin-top:.5rem">Views</div>
<div class="sidebar-item" onclick="filterPinned()">&#x1f4cc; Pinned</div>
<div class="sidebar-item" onclick="filterArchived()">&#x1f4e6; Archive</div>
<div class="sidebar-bottom" id="sWords">-</div>
</div>
<div class="list-pane">
<div class="list-toolbar">
<input type="text" id="searchBox" placeholder="Search notes..." onkeydown="if(event.key==='Enter')loadNotes()">
<button class="btn btn-p" onclick="newNote()">+</button>
</div>
<div class="list-scroll" id="noteList"></div>
</div>
<div class="editor-pane" id="editorPane" style="display:none">
<div class="editor-toolbar">
<select id="edNb" onchange="saveNote()"></select>
<button class="btn btn-g" id="pinBtn" onclick="togglePin()">Pin</button>
<button class="btn btn-d" id="previewBtn" onclick="togglePreview()">Preview</button>
<button class="btn btn-d" onclick="exportNote()">Export</button>
<button class="btn btn-d" onclick="archiveCur()">Archive</button>
<span style="flex:1"></span>
<span style="font-size:.6rem;color:var(--cm)" id="edMeta"></span>
<button class="btn btn-d" onclick="deleteCur()">Del</button>
</div>
<input class="editor-title" id="edTitle" placeholder="Note title..." oninput="autoSave()">
<div class="editor-tags"><input type="text" id="edTags" placeholder="Tags (comma-separated)" onchange="saveNote()"></div>
<div class="editor-body">
<textarea id="edBody" oninput="autoSave()" placeholder="Write in Markdown..."></textarea>
<div class="preview" id="previewArea"></div>
</div>
</div>
<div id="editorEmpty" class="editor-pane" style="display:flex;align-items:center;justify-content:center">
<div class="empty">Select a note or create a new one.</div>
</div>
</div>
<div id="modal"></div>

<script>
let notes=[],notebooks=[],tags=[],curNote=null,curFilter={},saveTimer=null,previewOn=false;

async function api(url,opts){const r=await fetch(url,opts);return r.json()}
function esc(s){return String(s||'').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;')}
function timeAgo(d){if(!d)return'';const s=Math.floor((Date.now()-new Date(d))/1e3);if(s<60)return s+'s ago';if(s<3600)return Math.floor(s/60)+'m ago';if(s<86400)return Math.floor(s/3600)+'h ago';return Math.floor(s/86400)+'d ago'}

async function init(){
  const [nbd,td,sd]=await Promise.all([api('/api/notebooks'),api('/api/tags'),api('/api/stats')]);
  notebooks=nbd.notebooks||[];tags=td.tags||[];
  renderSidebar(sd);
  loadNotes();
}

function renderSidebar(sd){
  document.getElementById('sAll').textContent=sd.notes;
  document.getElementById('sWords').textContent=sd.words.toLocaleString()+' words';
  const nbEl=document.getElementById('nbList');
  nbEl.innerHTML=notebooks.map(nb=>'<div class="sidebar-item" onclick="filterNotebook(\''+nb.id+'\')" data-nb="'+nb.id+'"><div class="sidebar-dot" style="background:'+esc(nb.color)+'"></div>'+esc(nb.name)+'<span class="sidebar-count">'+nb.note_count+'</span></div>').join('');
  const tagEl=document.getElementById('tagList');
  tagEl.innerHTML=(tags||[]).map(t=>'<div class="sidebar-item" onclick="filterTag(\''+esc(t)+'\')"><span style="color:var(--leather)">#</span>'+esc(t)+'</div>').join('');
  // update notebook selector in editor
  const sel=document.getElementById('edNb');
  sel.innerHTML='<option value="">No notebook</option>'+notebooks.map(nb=>'<option value="'+nb.id+'">'+esc(nb.name)+'</option>').join('');
  if(curNote)sel.value=curNote.notebook_id||'';
}

function filterNotebook(id){curFilter={notebook_id:id};highlightSidebar(id?id:'all');loadNotes()}
function filterTag(t){curFilter={tag:t};loadNotes()}
function filterPinned(){curFilter={pinned:'true'};loadNotes()}
function filterArchived(){curFilter={archived:'true'};loadNotes()}

function highlightSidebar(id){
  document.querySelectorAll('.sidebar-item').forEach(el=>{
    el.classList.toggle('active',el.dataset.nb===id)
  })
}

async function loadNotes(){
  const p=new URLSearchParams();
  if(curFilter.notebook_id)p.set('notebook_id',curFilter.notebook_id);
  if(curFilter.tag)p.set('tag',curFilter.tag);
  if(curFilter.pinned)p.set('pinned',curFilter.pinned);
  if(curFilter.archived)p.set('archived',curFilter.archived);
  const q=document.getElementById('searchBox').value;
  if(q)p.set('search',q);
  const d=await api('/api/notes?'+p);
  notes=d.notes||[];
  renderNotes();
}

function renderNotes(){
  const el=document.getElementById('noteList');
  if(!notes.length){el.innerHTML='<div class="empty" style="padding:2rem">No notes found.</div>';return}
  el.innerHTML=notes.map(n=>{
    const tags=(n.tags||[]).map(t=>'<span class="tag-chip">'+esc(t)+'</span>').join(' ');
    const preview=n.body.substring(0,80).replace(/\n/g,' ');
    const active=curNote&&curNote.id===n.id?'active':'';
    return '<div class="note-item '+active+'" onclick="openNote(\''+n.id+'\')">'+
      '<div class="note-item-title">'+(n.pinned?'<span class="pin">&#x1f4cc;</span>':'')+esc(n.title)+'</div>'+
      '<div class="note-item-preview">'+esc(preview)+'</div>'+
      '<div class="note-item-meta"><span>'+timeAgo(n.updated_at)+'</span><span>'+n.word_count+'w</span>'+tags+'</div></div>'
  }).join('')
}

async function openNote(id){
  const n=await api('/api/notes/'+id);
  curNote=n;
  document.getElementById('editorPane').style.display='flex';
  document.getElementById('editorEmpty').style.display='none';
  document.getElementById('edTitle').value=n.title;
  document.getElementById('edBody').value=n.body;
  document.getElementById('edTags').value=(n.tags||[]).join(', ');
  document.getElementById('edNb').value=n.notebook_id||'';
  document.getElementById('pinBtn').textContent=n.pinned?'Unpin':'Pin';
  document.getElementById('edMeta').textContent=n.word_count+'w · '+timeAgo(n.updated_at);
  if(previewOn)updatePreview();
  renderNotes();
}

function newNote(){
  const nb=curFilter.notebook_id||'';
  api('/api/notes',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({title:'Untitled',body:'',notebook_id:nb,tags:[]})}).then(n=>{
    loadNotes().then(()=>openNote(n.id));
    init();
  })
}

function autoSave(){
  if(saveTimer)clearTimeout(saveTimer);
  saveTimer=setTimeout(saveNote,500);
}

async function saveNote(){
  if(!curNote)return;
  const tags=(document.getElementById('edTags').value||'').split(',').map(s=>s.trim()).filter(Boolean);
  const body={
    title:document.getElementById('edTitle').value||'Untitled',
    body:document.getElementById('edBody').value,
    tags:tags,
    notebook_id:document.getElementById('edNb').value
  };
  await api('/api/notes/'+curNote.id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
  curNote.title=body.title;curNote.body=body.body;curNote.tags=tags;
  curNote.word_count=body.body.split(/\s+/).filter(Boolean).length;
  document.getElementById('edMeta').textContent=curNote.word_count+'w · just now';
  if(previewOn)updatePreview();
  renderNotes();
}

async function togglePin(){
  if(!curNote)return;
  const action=curNote.pinned?'unpin':'pin';
  await api('/api/notes/'+curNote.id+'/'+action,{method:'POST'});
  curNote.pinned=!curNote.pinned;
  document.getElementById('pinBtn').textContent=curNote.pinned?'Unpin':'Pin';
  loadNotes();
}

async function archiveCur(){
  if(!curNote)return;
  await api('/api/notes/'+curNote.id+'/archive',{method:'POST'});
  curNote=null;
  document.getElementById('editorPane').style.display='none';
  document.getElementById('editorEmpty').style.display='flex';
  loadNotes();init();
}

async function deleteCur(){
  if(!curNote||!confirm('Delete this note?'))return;
  await api('/api/notes/'+curNote.id,{method:'DELETE'});
  curNote=null;
  document.getElementById('editorPane').style.display='none';
  document.getElementById('editorEmpty').style.display='flex';
  loadNotes();init();
}

function exportNote(){
  if(!curNote)return;
  window.open('/api/notes/'+curNote.id+'/export','_blank');
}

function togglePreview(){
  previewOn=!previewOn;
  document.getElementById('previewArea').style.display=previewOn?'block':'none';
  document.getElementById('previewBtn').textContent=previewOn?'Edit':'Preview';
  if(previewOn)updatePreview();
}

function updatePreview(){
  const md=document.getElementById('edBody').value;
  document.getElementById('previewArea').innerHTML=renderMd(md);
}

function renderMd(md){
  let html=esc(md);
  html=html.replace(/^### (.+)$/gm,'<h3>$1</h3>');
  html=html.replace(/^## (.+)$/gm,'<h2>$1</h2>');
  html=html.replace(/^# (.+)$/gm,'<h1>$1</h1>');
  html=html.replace(/\*\*(.+?)\*\*/g,'<strong>$1</strong>');
  html=html.replace(/\*(.+?)\*/g,'<em>$1</em>');
  html=html.replace(/` + "`" + `([^` + "`" + `]+)` + "`" + `/g,'<code>$1</code>');
  html=html.replace(/^&gt; (.+)$/gm,'<blockquote>$1</blockquote>');
  html=html.replace(/^- (.+)$/gm,'<li>$1</li>');
  html=html.replace(/(<li>.*<\/li>)/s,function(m){return '<ul>'+m+'</ul>'});
  html=html.replace(/\n\n/g,'</p><p>');
  html='<p>'+html+'</p>';
  return html;
}

function showNewNotebook(){
  document.getElementById('modal').innerHTML='<div class="modal-bg" onclick="if(event.target===this)closeModal()"><div class="modal">'+
    '<h2>New Notebook</h2>'+
    '<label class="fl">Name</label><input type="text" id="nn-name">'+
    '<label class="fl">Color</label><input type="color" id="nn-color" value="#c45d2c" style="height:30px;padding:2px">'+
    '<div style="display:flex;gap:.5rem;margin-top:1rem"><button class="btn btn-p" onclick="saveNewNb()">Create</button><button class="btn btn-d" onclick="closeModal()">Cancel</button></div>'+
  '</div></div>'
}
async function saveNewNb(){
  const body={name:document.getElementById('nn-name').value,color:document.getElementById('nn-color').value};
  if(!body.name){alert('Name required');return}
  await api('/api/notebooks',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
  closeModal();init();
}
function closeModal(){document.getElementById('modal').innerHTML=''}

init();
</script></body></html>`
