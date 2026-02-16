// web_demo client — interacts with Go backend API
const API_BASE = '/api';

async function generateNPC() {
  const r = await fetch(`${API_BASE}/generate`, { method: 'POST' });
  if (!r.ok) { alert('Generate failed: ' + r.statusText); return; }
  await renderList();
}

async function listNPCs() {
  const r = await fetch(`${API_BASE}/npcs`);
  if (!r.ok) return [];
  return await r.json();
}

async function getNpc(id) {
  const r = await fetch(`${API_BASE}/npcs/${encodeURIComponent(id)}`);
  if (!r.ok) return null;
  return await r.json();
}

async function deleteNPC(id) {
  await fetch(`${API_BASE}/npcs/${encodeURIComponent(id)}`, { method: 'DELETE' });
  await renderList();
}

async function clearAll() {
  const items = await listNPCs();
  await Promise.all(items.map(it => fetch(`${API_BASE}/npcs/${encodeURIComponent(it.ID||it.id)}`, { method: 'DELETE' })));
  await renderList();
}

function makeItemDiv(it) {
  const id = it.ID || it.id;
  const name = (it.Components && it.Components['1']) || it.name || '';
  const species = (it.Components && it.Components['4']) || it.species || '';
  const type = (it.Components && it.Components['2']) || '';
  const subtype = (it.Components && it.Components['3']) || '';
  const faction = (it.Components && it.Components['5']) || '';
  const created = it.createdAt || '';
  const div = document.createElement('div');
  div.className = 'npc';
  div.innerHTML = `
    <strong class="npc-name">${name}</strong> — <em>${species}</em> <button data-id="${id}">Delete</button>
    <div class="type">Type:${type}</div> 
    <div class="subtype">Subtype: ${subtype}</div>
    <div class="faction">Faction: ${faction}</div>
   
  `;
  const btn = div.querySelector('button');
  btn.addEventListener('click', () => deleteNPC(id));
  div.querySelector('.npc-name').style.cursor = 'pointer';
  div.querySelector('.npc-name').addEventListener('click', () => showDetails(id));
  return div;
}

async function renderList() {
  const container = document.getElementById('list');
  container.innerHTML = '';
  const items = await listNPCs();
  if (!items || !items.length) { container.innerHTML = '<em>No NPCs stored yet.</em>'; return; }
  for (const it of items) {
    container.appendChild(makeItemDiv(it));
  }
}

async function main() {
  document.getElementById('btnGenerate').addEventListener('click', generateNPC);
  document.getElementById('btnRefresh').addEventListener('click', renderList);
  document.getElementById('btnClear').addEventListener('click', async () => {
    if (confirm('Delete all stored NPCs?')) await clearAll();
  });
  document.getElementById('btnSave').addEventListener('click', saveDetails);
  document.getElementById('btnCreate').addEventListener('click', createFromForm);
  document.getElementById('btnClose').addEventListener('click', () => document.getElementById('npcForm').reset());
  await renderList();
}

main().catch(e => { console.error(e); alert('Demo error: '+e.message); });

async function showDetails(id) {
  const npc = await getNpc(id);
  if (!npc) { alert('Failed to load NPC'); return; }
  document.getElementById('f_id').value = npc.ID || npc.id || '';
  document.getElementById('f_name').value = (npc.Components && npc.Components['1']) || npc.name || '';
  document.getElementById('f_type').value = (npc.Components && npc.Components['2']) || '';
  document.getElementById('f_subtype').value = (npc.Components && npc.Components['3']) || '';
  document.getElementById('f_species').value = (npc.Components && npc.Components['4']) || npc.species || '';
  document.getElementById('f_faction').value = (npc.Components && npc.Components['5']) || '';
  document.getElementById('f_traits').value = (npc.Components && npc.Components['6']) || (npc.traits||[]).join(', ');
  document.getElementById('f_stats').value = (npc.Components && npc.Components['7']) || '';
  document.getElementById('f_items').value = (npc.Components && npc.Components['8']) || '';
  document.getElementById('f_description').value = (npc.Components && npc.Components['9']) || npc.description || '';
  document.getElementById('f_location').value = npc.LocationID || '';
}

async function gatherForm() {
  return {
    id: document.getElementById('f_id').value,
    name: document.getElementById('f_name').value,
    type: document.getElementById('f_type').value,
    subtype: document.getElementById('f_subtype').value,
    species: document.getElementById('f_species').value,
    faction: document.getElementById('f_faction').value,
    traits: document.getElementById('f_traits').value.split(',').map(s=>s.trim()).filter(Boolean),
    stats: document.getElementById('f_stats').value,
    items: document.getElementById('f_items').value,
    description: document.getElementById('f_description').value,
    locationID: document.getElementById('f_location').value || 'default'
  };
}

async function saveDetails() {
  const payload = await gatherForm();
  if (!payload.id) { alert('No ID present. Use Create New to add.'); return; }
  const r = await fetch(`${API_BASE}/npcs/${encodeURIComponent(payload.id)}`, {
    method: 'PUT',
    headers: {'Content-Type':'application/json'},
    body: JSON.stringify(payload)
  });
  if (!r.ok) { alert('Save failed: '+r.statusText); return; }
  await renderList();
  alert('Saved');
}

async function createFromForm() {
  const payload = await gatherForm();
  if (!payload.id) payload.id = `${Date.now()}-${Math.random()}`;
  const r = await fetch(`${API_BASE}/npcs`, {
    method: 'POST',
    headers: {'Content-Type':'application/json'},
    body: JSON.stringify(payload)
  });
  if (!r.ok) { alert('Create failed: '+r.statusText); return; }
  await renderList();
  alert('Created');
}
