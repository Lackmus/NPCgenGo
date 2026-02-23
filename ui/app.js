// Unified app.js — works for both Web and Wails UIs
// Auto-detects environment and uses appropriate API calls

let creationOptions = null;

// Lazy environment detection and API adapter
function getBackend() {
  return window?.go?.main?.WailsAPI || window?.go?.main?.App;
}

function isWailsEnv() {
  return Boolean(getBackend());
}

const API = {
  async generateNPC() {
    if (isWailsEnv()) {
      return await getBackend().GenerateNPC();
    }
    const r = await fetch('/api/generate', { method: 'POST' });
    if (!r.ok) throw new Error(`Generate failed: ${r.statusText}`);
  },
  async listNPCs() {
    if (isWailsEnv()) {
      return await getBackend().ListNPCs();
    }
    const r = await fetch('/api/npcs');
    return r.ok ? await r.json() : [];
  },
  async getNpc(id) {
    if (isWailsEnv()) {
      return await getBackend().GetNPC(id);
    }
    const r = await fetch(`/api/npcs/${encodeURIComponent(id)}`);
    return r.ok ? await r.json() : null;
  },
  async deleteNPC(id) {
    if (isWailsEnv()) {
      return await getBackend().DeleteNPC(id);
    }
    return await fetch(`/api/npcs/${encodeURIComponent(id)}`, { method: 'DELETE' });
  },
  async deleteAllNPCs() {
    if (isWailsEnv()) {
      return await getBackend().DeleteAllNPCs();
    }
    const items = await this.listNPCs();
    await Promise.all(items.map(it => this.deleteNPC(it.ID || it.id)));
  },
  async getCreationOptions() {
    if (isWailsEnv()) {
      return await getBackend().GetCreationOptions();
    }
    const r = await fetch('/api/options');
    return r.ok ? await r.json() : null;
  },
  async rollSubtypeFields(subtype) {
    if (!subtype) return { stats: '', items: '', description: '' };
    if (isWailsEnv()) {
      const rolled = await getBackend().RollSubtypeFields(subtype);
      return {
        stats: rolled?.stats || rolled?.Stats || '',
        items: rolled?.items || rolled?.Items || '',
        description: rolled?.description || rolled?.Description || ''
      };
    }
    const r = await fetch(`/api/subtypes/${encodeURIComponent(subtype)}/roll`);
    if (!r.ok) throw new Error(`Failed to roll subtype fields: ${r.statusText}`);
    return await r.json();
  },
  async rollSpeciesName(species) {
    if (!species) return { name: '' };
    if (isWailsEnv()) {
      const name = await getBackend().RollSpeciesName(species);
      return { name: name || '' };
    }
    const r = await fetch(`/api/species/${encodeURIComponent(species)}/name`);
    if (!r.ok) throw new Error(`Failed to roll species name: ${r.statusText}`);
    return await r.json();
  },
  async saveNPC(payload) {
    if (isWailsEnv()) {
      return await getBackend().SaveNPC(payload);
    }
    const r = await fetch(`/api/npcs/${encodeURIComponent(payload.id)}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    });
    if (!r.ok) throw new Error(`Save failed: ${r.statusText}`);
    return await r.json();
  }
};

// UI helper functions
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
  const div = document.createElement('div');
  div.className = 'npc';
  div.innerHTML = `
    <strong class="npc-name">${name}</strong> — <em>${species}</em> <button data-id="${id}">Delete</button>
    <div class="type">Type: ${type}</div> 
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
  const items = await API.listNPCs();
  if (!items || !items.length) {
    container.innerHTML = '<em>No NPCs stored yet.</em>';
    return;
  }
  for (const it of items) {
    container.appendChild(makeItemDiv(it));
  }
}

async function deleteNPC(id) {
  await API.deleteNPC(id);
  await renderList();
}

async function generateNPC() {
  try {
    await API.generateNPC();
    await renderList();
  } catch (e) {
    alert('Generate failed: ' + e.message);
  }
}

async function showDetails(id) {
  const npc = await API.getNpc(id);
  if (!npc) {
    alert('Failed to load NPC');
    return;
  }
  document.getElementById('f_id').value = npc.ID || npc.id || '';
  document.getElementById('f_name').value = (npc.Components && npc.Components['1']) || npc.name || '';
  const selectedType = (npc.Components && npc.Components['2']) || '';
  const selectedSubtype = (npc.Components && npc.Components['3']) || '';
  const selectedSpecies = (npc.Components && npc.Components['4']) || '';
  const selectedFaction = (npc.Components && npc.Components['5']) || '';
  setSelectValue('f_type', selectedType);
  updateSubtypeDropdown(selectedType, selectedSubtype);
  setSelectValue('f_faction', selectedFaction);
  updateSpeciesDropdown(selectedFaction, selectedSpecies);
  const traitValue = (npc.Components && npc.Components['6']) || '';
  setSelectValue('f_traits', traitValue.split(',')[0]?.trim() || '');
  document.getElementById('f_stats').value = (npc.Components && npc.Components['7']) || '';
  document.getElementById('f_items').value = (npc.Components && npc.Components['8']) || '';
  document.getElementById('f_description').value = (npc.Components && npc.Components['9']) || npc.description || '';
  document.getElementById('f_location').value = npc.LocationID || 'default';
  setRerollEnabled(Boolean((selectedSubtype || '').trim()));
}

async function gatherForm() {
  const traitValue = document.getElementById('f_traits').value;
  return {
    id: document.getElementById('f_id').value,
    name: document.getElementById('f_name').value,
    type: document.getElementById('f_type').value,
    subtype: document.getElementById('f_subtype').value,
    species: document.getElementById('f_species').value,
    faction: document.getElementById('f_faction').value,
    traits: traitValue ? [traitValue] : [],
    stats: document.getElementById('f_stats').value,
    items: document.getElementById('f_items').value,
    description: document.getElementById('f_description').value,
    locationID: document.getElementById('f_location').value || 'default'
  };
}

async function saveDetails() {
  try {
    const payload = await gatherForm();
    if (!payload.id) {
      alert('No ID present. Generate an NPC first.');
      return;
    }
    await API.saveNPC(payload);
    await renderList();
    alert('Saved');
  } catch (e) {
    alert('Save failed: ' + e.message);
  }
}

function clearForm() {
  document.getElementById('npcForm').reset();
  document.getElementById('f_location').value = 'default';
  setSubtypeEnabled(false);
  setSpeciesEnabled(false);
  setRerollEnabled(false);
}

async function main() {
  try {
    creationOptions = await API.getCreationOptions();
    populateOptionDropdowns();
    document.getElementById('f_type').addEventListener('change', (event) => {
      updateSubtypeDropdown(event.target.value);
      document.getElementById('f_stats').value = '';
      document.getElementById('f_items').value = '';
      document.getElementById('f_description').value = '';
      setRerollEnabled(false);
    });
    document.getElementById('f_faction').addEventListener('change', (event) => {
      updateSpeciesDropdown(event.target.value);
      document.getElementById('f_name').value = '';
    });
    document.getElementById('f_subtype').addEventListener('change', async (event) => {
      setRerollEnabled(Boolean((event.target.value || '').trim()));
      try {
        const rolled = await API.rollSubtypeFields(event.target.value);
        document.getElementById('f_stats').value = rolled.stats || '';
        document.getElementById('f_items').value = rolled.items || '';
        document.getElementById('f_description').value = rolled.description || '';
      } catch (error) {
        console.error(error);
        alert(error.message || 'Failed to generate subtype fields.');
      }
    });
    document.getElementById('f_species').addEventListener('change', async (event) => {
      try {
        const rolled = await API.rollSpeciesName(event.target.value);
        document.getElementById('f_name').value = rolled.name || '';
      } catch (error) {
        console.error(error);
        alert(error.message || 'Failed to generate species name.');
      }
    });
    document.getElementById('btnGenerate').addEventListener('click', generateNPC);
    document.getElementById('btnRefresh').addEventListener('click', renderList);
    document.getElementById('btnClear').addEventListener('click', async () => {
      if (confirm('Delete all stored NPCs?')) {
        await API.deleteAllNPCs();
        await renderList();
        clearForm();
      }
    });
    document.getElementById('btnSave').addEventListener('click', saveDetails);
    document.getElementById('btnReroll').addEventListener('click', async () => {
      const subtype = document.getElementById('f_subtype').value;
      if (!subtype) {
        alert('Select a subtype first.');
        return;
      }
      try {
        const rolled = await API.rollSubtypeFields(subtype);
        document.getElementById('f_stats').value = rolled.stats || '';
        document.getElementById('f_items').value = rolled.items || '';
        document.getElementById('f_description').value = rolled.description || '';
      } catch (error) {
        console.error(error);
        alert(error.message || 'Failed to reroll subtype fields.');
      }
    });
    document.getElementById('btnClose').addEventListener('click', clearForm);
    setSubtypeEnabled(Boolean((document.getElementById('f_type').value || '').trim()));
    setSpeciesEnabled(Boolean((document.getElementById('f_faction').value || '').trim()));
    setRerollEnabled(Boolean((document.getElementById('f_subtype').value || '').trim()));
    await renderList();
  } catch (error) {
    console.error(error);
    alert('Failed to initialize app: ' + error.message);
  }
}

main();
