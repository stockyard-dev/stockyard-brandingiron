package server

import "net/http"

func (s *Server) handleUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Branding Iron — Stockyard</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;600&family=Libre+Baskerville:wght@400;700&display=swap" rel="stylesheet">
<style>*{margin:0;padding:0;box-sizing:border-box}body{background:#1a1410;color:#f0e6d3;font-family:'Libre Baskerville',serif;padding:2rem}
.hdr{font-family:'JetBrains Mono',monospace;font-size:.7rem;color:#a0845c;letter-spacing:3px;text-transform:uppercase;margin-bottom:2rem;border-bottom:2px solid #8b3d1a;padding-bottom:.8rem}
.cards{display:grid;grid-template-columns:repeat(2,1fr);gap:1rem;margin-bottom:2rem;font-family:'JetBrains Mono',monospace}.card{background:#241e18;border:1px solid #2e261e;padding:1rem}.card-val{font-size:1.6rem;font-weight:700;display:block}.card-lbl{font-size:.55rem;letter-spacing:2px;text-transform:uppercase;color:#a0845c;margin-top:.2rem}
.section{margin-bottom:2rem}.section h2{font-family:'JetBrains Mono',monospace;font-size:.65rem;letter-spacing:3px;text-transform:uppercase;color:#e8753a;margin-bottom:.8rem;border-bottom:1px solid #2e261e;padding-bottom:.4rem}
.lbl{font-family:'JetBrains Mono',monospace;font-size:.62rem;letter-spacing:1px;text-transform:uppercase;color:#a0845c}
input{font-family:'JetBrains Mono',monospace;font-size:.78rem;background:#2e261e;border:1px solid #2e261e;color:#f0e6d3;padding:.4rem .7rem;outline:none;width:100%}input:focus{border-color:#a0845c}
.row{display:flex;gap:.8rem;align-items:flex-end;flex-wrap:wrap;margin-bottom:1rem}.field{display:flex;flex-direction:column;gap:.3rem}
.btn{font-family:'JetBrains Mono',monospace;font-size:.7rem;padding:.4rem 1rem;border:1px solid #c45d2c;background:transparent;color:#e8753a;cursor:pointer}.btn:hover{background:#c45d2c;color:#f0e6d3}
.preview{background:#241e18;border:1px solid #2e261e;padding:1rem;margin-top:1rem;text-align:center}
.preview img{max-width:100%;height:auto;border:1px solid #2e261e}
pre{font-family:'JetBrains Mono',monospace;font-size:.7rem;background:#241e18;padding:.8rem;color:#bfb5a3;margin-top:.5rem;overflow-x:auto}
</style></head><body>
<div class="hdr">Stockyard · Branding Iron</div>
<div class="cards"><div class="card"><span class="card-val" id="s-tpl">—</span><span class="card-lbl">Templates</span></div><div class="card"><span class="card-val" id="s-gen">—</span><span class="card-lbl">Generated</span></div></div>
<div class="section"><h2>Generate OG Image</h2>
<div class="row">
  <div class="field" style="flex:2"><span class="lbl">Title</span><input id="g-title" placeholder="My Blog Post Title" value="Hello World"></div>
  <div class="field" style="flex:1"><span class="lbl">Subtitle</span><input id="g-sub" placeholder="Optional subtitle" value="A great article"></div>
  <button class="btn" onclick="generate()">Generate</button>
</div>
<div id="preview" class="preview" style="display:none"></div>
<div id="url-box" style="margin-top:.5rem"></div>
</div>
<div class="section"><h2>Usage</h2>
<pre>
&lt;!-- Add to your HTML head --&gt;
&lt;meta property="og:image" content="http://localhost:9040/api/og?title=My+Post&amp;subtitle=Read+more" /&gt;

# Or fetch via API
curl "http://localhost:9040/api/og?title=Hello+World&amp;subtitle=My+Blog" &gt; og.svg
</pre>
</div>
<script>
async function refresh(){
  try{const s=await(await fetch('/api/status')).json();document.getElementById('s-tpl').textContent=s.templates||0;document.getElementById('s-gen').textContent=s.generations||0;}catch(e){}
}
function generate(){
  const title=encodeURIComponent(document.getElementById('g-title').value);
  const sub=encodeURIComponent(document.getElementById('g-sub').value);
  const url='/api/og?title='+title+'&subtitle='+sub;
  document.getElementById('preview').style.display='block';
  document.getElementById('preview').innerHTML='<img src="'+url+'" alt="OG Preview">';
  document.getElementById('url-box').innerHTML='<pre style="font-size:.65rem">'+location.origin+url+'</pre>';
  refresh();
}
refresh();setInterval(refresh,10000);
</script></body></html>`))
}
