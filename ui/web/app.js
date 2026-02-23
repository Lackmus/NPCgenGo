// ui client — interacts with Go backend API
const API_BASE = '/api';
let creationOptions = null;

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

async function getCreationOptions() {
  const r = await fetch(`${API_BASE}/options`);
  if (!r.ok) {
    return null;
  }
  return await r.json();
}

async function rollSubtypeFields(subtype) {
  if (!subtype) {
	return { stats: '', items: '' };
  }
  const r = await fetch(`${API_BASE}/subtypes/${encodeURIComponent(subtype)}/roll`);
  if (!r.ok) {
    throw new Error(`Failed to roll subtype fields: ${r.statusText}`);
  }
  return await r.json();
}

async function rollSpeciesName(species) {
  if (!species) {
    return { name: '' };
  }
  const r = await fetch(`${API_BASE}/species/${encodeURIComponent(species)}/name`);
  if (!r.ok) {
    throw new Error(`Failed to roll species name: ${r.statusText}`);
  }
  return await r.json();
}

async function applySubtypeRoll(subtype) {
  if (!subtype) {
    document.getElementById('f_stats').textContent = '—';
    document.getElementById('f_items').textContent = '—';
    return;
  }
  const rolled = await rollSubtypeFields(subtype);
  document.getElementById('f_stats').textContent = rolled.stats || '—';
  document.getElementById('f_items').textContent = rolled.items || '—';
}

async function applySpeciesNameRoll(species) {
	if (!species) {
		document.getElementById('f_name').value = '';
		return;
	}
	const rolled = await rollSpeciesName(species);
	document.getElementById('f_name').value = rolled.name || '';
}

function setSelectOptions(selectId, values, includeEmpty = true) {
  const select = document.getElementById(selectId);
  if (!select) return;
  const currentValue = select.value;
  select.innerHTML = '';
  if (includeEmpty) {
    const emptyOption = document.createElement('option');
    emptyOption.value = '';
    emptyOption.textContent = '';
    select.appendChild(emptyOption);
  }
  for (const value of values || []) {
    const option = document.createElement('option');
    option.value = value;
    option.textContent = value;
    select.appendChild(option);
  }
  if (currentValue && Array.from(select.options).some((option) => option.value === currentValue)) {
    select.value = currentValue;
  }
}

function setSelectValue(selectId, value) {
  const select = document.getElementById(selectId);
  if (!select) return;
  const desired = value || '';
  if (desired && !Array.from(select.options).some((option) => option.value === desired)) {
    const option = document.createElement('option');
    option.value = desired;
    option.textContent = desired;
    select.appendChild(option);
  }
  select.value = desired;
}

function setRerollEnabled(enabled) {
  const button = document.getElementById('btnReroll');
  if (!button) return;
  button.disabled = !enabled;
}

function setSubtypeEnabled(enabled) {
  const select = document.getElementById('f_subtype');
  if (!select) return;
  select.disabled = !enabled;
}

function setSpeciesEnabled(enabled) {
  const select = document.getElementById('f_species');
  if (!select) return;
  select.disabled = !enabled;
}

function setNameRerollEnabled(enabled) {
  const button = document.getElementById('btnRerollName');
  if (!button) return;
  button.disabled = !enabled;
}

function updateSubtypeDropdown(selectedType, selectedSubtype = '') {
  if (!creationOptions) return;
  setSubtypeEnabled(Boolean((selectedType || '').trim()));
  const subtypeMap = creationOptions.NpcSubtypeForTypeMap || creationOptions.npcSubtypeForTypeMap || {};
  const subtypes = subtypeMap[selectedType] || [];
  setSelectOptions('f_subtype', subtypes, true);
  setSelectValue('f_subtype', selectedSubtype);
	setRerollEnabled(Boolean((selectedSubtype || '').trim()));
}

function updateSpeciesDropdown(selectedFaction, selectedSpecies = '') {
  if (!creationOptions) return;
  setSpeciesEnabled(Boolean((selectedFaction || '').trim()));
  const speciesMap = creationOptions.NpcSpeciesForFactionMap || creationOptions.npcSpeciesForFactionMap || {};
  const species = speciesMap[selectedFaction] || [];
  setSelectOptions('f_species', species, true);
  setSelectValue('f_species', selectedSpecies);
  setNameRerollEnabled(Boolean((selectedSpecies || '').trim()));
}

function populateOptionDropdowns() {
  if (!creationOptions) return;
  setSelectOptions('f_type', creationOptions.NpcTypes || creationOptions.npcTypes || [], true);
  setSelectOptions('f_faction', creationOptions.Factions || creationOptions.factions || [], true);
  setSelectOptions('f_traits', creationOptions.Traits || creationOptions.traits || [], true);
  updateSubtypeDropdown(document.getElementById('f_type').value);
  updateSpeciesDropdown(document.getElementById('f_faction').value);
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
  creationOptions = await getCreationOptions();
  populateOptionDropdowns();
  document.getElementById('f_type').addEventListener('change', (event) => {
    updateSubtypeDropdown(event.target.value);
    document.getElementById('f_stats').textContent = '—';
    document.getElementById('f_items').textContent = '—';
		setRerollEnabled(false);
  });
  document.getElementById('f_faction').addEventListener('change', (event) => {
    updateSpeciesDropdown(event.target.value);
    document.getElementById('f_name').value = '';
  });
  document.getElementById('f_subtype').addEventListener('change', async (event) => {
    setRerollEnabled(Boolean((event.target.value || '').trim()));
    try {
      await applySubtypeRoll(event.target.value);
    } catch (error) {
      console.error(error);
      alert(error.message || 'Failed to generate subtype fields.');
    }
  });
  document.getElementById('f_species').addEventListener('change', async (event) => {
    setNameRerollEnabled(Boolean((event.target.value || '').trim()));
    try {
      await applySpeciesNameRoll(event.target.value);
    } catch (error) {
      console.error(error);
      alert(error.message || 'Failed to generate species name.');
    }
  });
  document.getElementById('btnGenerate').addEventListener('click', generateNPC);
  document.getElementById('btnRefresh').addEventListener('click', renderList);
  document.getElementById('btnClear').addEventListener('click', async () => {
    if (confirm('Delete all stored NPCs?')) await clearAll();
  });
  document.getElementById('btnSave').addEventListener('click', saveDetails);
  document.getElementById('btnReroll').addEventListener('click', async () => {
    const subtype = document.getElementById('f_subtype').value;
    if (!subtype) {
      alert('Select a subtype first.');
      return;
    }
    try {
      await applySubtypeRoll(subtype);
    } catch (error) {
      console.error(error);
      alert(error.message || 'Failed to reroll subtype fields.');
    }
  });
  document.getElementById('btnRerollName').addEventListener('click', async () => {
    const species = document.getElementById('f_species').value;
    if (!species) {
      alert('Select a species first.');
      return;
    }
    try {
      await applySpeciesNameRoll(species);
    } catch (error) {
      console.error(error);
      alert(error.message || 'Failed to reroll name.');
    }
  });
  document.getElementById('btnClose').addEventListener('click', () => {
    document.getElementById('npcForm').reset();
    setSubtypeEnabled(false);
    setSpeciesEnabled(false);
    setRerollEnabled(false);
    setNameRerollEnabled(false);
  });
  setSubtypeEnabled(Boolean((document.getElementById('f_type').value || '').trim()));
	setSpeciesEnabled(Boolean((document.getElementById('f_faction').value || '').trim()));
	setRerollEnabled(Boolean((document.getElementById('f_subtype').value || '').trim()));
	setNameRerollEnabled(Boolean((document.getElementById('f_species').value || '').trim()));
  await renderList();
}

main().catch(e => { console.error(e); alert('Demo error: '+e.message); });

async function showDetails(id) {
  const npc = await getNpc(id);
  if (!npc) { alert('Failed to load NPC'); return; }
  document.getElementById('f_id').value = npc.ID || npc.id || '';
  document.getElementById('f_name').value = (npc.Components && npc.Components['1']) || npc.name || '';
  document.getElementById('f_type').value = (npc.Components && npc.Components['2']) || '';
  const selectedType = (npc.Components && npc.Components['2']) || '';
  const selectedSubtype = (npc.Components && npc.Components['3']) || '';
  setSelectValue('f_type', selectedType);
  updateSubtypeDropdown(selectedType, selectedSubtype);
  setSelectValue('f_species', (npc.Components && npc.Components['4']) || npc.species || '');
  setSelectValue('f_faction', (npc.Components && npc.Components['5']) || '');
  updateSpeciesDropdown((npc.Components && npc.Components['5']) || '', (npc.Components && npc.Components['4']) || npc.species || '');
  const traitValue = (npc.Components && npc.Components['6']) || '';
  setSelectValue('f_traits', traitValue.split(',')[0]?.trim() || '');
  document.getElementById('f_stats').textContent = (npc.Components && npc.Components['7']) || '—';
  document.getElementById('f_items').textContent = (npc.Components && npc.Components['8']) || '—';
	setRerollEnabled(Boolean((selectedSubtype || '').trim()));
	setNameRerollEnabled(Boolean(((npc.Components && npc.Components['4']) || npc.species || '').trim()));
}

async function gatherForm() {
	const traitValue = document.getElementById('f_traits').value;
	const statsValue = (document.getElementById('f_stats')?.textContent || '').trim();
	const itemsValue = (document.getElementById('f_items')?.textContent || '').trim();
  return {
    id: document.getElementById('f_id').value,
    name: document.getElementById('f_name').value,
    type: document.getElementById('f_type').value,
    subtype: document.getElementById('f_subtype').value,
    species: document.getElementById('f_species').value,
    faction: document.getElementById('f_faction').value,
    trait: traitValue,
    stats: statsValue === '—' ? '' : statsValue,
    items: itemsValue === '—' ? '' : itemsValue
  };
}

async function saveDetails() {
  const payload = await gatherForm();
  if (!payload.id) { alert('No ID present. Generate an NPC first.'); return; }
  const r = await fetch(`${API_BASE}/npcs/${encodeURIComponent(payload.id)}`, {
    method: 'PUT',
    headers: {'Content-Type':'application/json'},
    body: JSON.stringify(payload)
  });
  if (!r.ok) { alert('Save failed: '+r.statusText); return; }
  await renderList();
  alert('Saved');
}

