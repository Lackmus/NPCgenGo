function appBindings() {
  const backend = window?.go?.main?.WailsAPI || window?.go?.main?.App;
  if (!backend) {
    throw new Error("Wails bindings not available. Expected window.go.main.WailsAPI (or legacy App).");
  }
  return backend;
}

let creationOptions = null;

function compValue(npc, key) {
  if (!npc || !npc.Components) return "";
  return npc.Components[key] || "";
}

function currentId() {
  return document.getElementById("f_id").value;
}

function readForm() {
	const traitValue = document.getElementById("f_traits").value;
  return {
    id: document.getElementById("f_id").value,
    name: document.getElementById("f_name").value,
    type: document.getElementById("f_type").value,
    subtype: document.getElementById("f_subtype").value,
    species: document.getElementById("f_species").value,
    faction: document.getElementById("f_faction").value,
    traits: traitValue ? [traitValue] : [],
    stats: document.getElementById("f_stats").value,
    items: document.getElementById("f_items").value,
    description: document.getElementById("f_description").value,
    locationID: document.getElementById("f_location").value || "default",
  };
}

function setSelectOptions(selectId, values, includeEmpty = true) {
  const select = document.getElementById(selectId);
  if (!select) return;
  const currentValue = select.value;
  select.innerHTML = "";
  if (includeEmpty) {
    const emptyOption = document.createElement("option");
    emptyOption.value = "";
    emptyOption.textContent = "";
    select.appendChild(emptyOption);
  }
  for (const value of values || []) {
    const option = document.createElement("option");
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
  const desired = value || "";
  if (desired && !Array.from(select.options).some((option) => option.value === desired)) {
    const option = document.createElement("option");
    option.value = desired;
    option.textContent = desired;
    select.appendChild(option);
  }
  select.value = desired;
}

function setRerollEnabled(enabled) {
  const button = document.getElementById("btnReroll");
  if (!button) return;
  button.disabled = !enabled;
}

function setSubtypeEnabled(enabled) {
  const select = document.getElementById("f_subtype");
  if (!select) return;
  select.disabled = !enabled;
}

function setSpeciesEnabled(enabled) {
  const select = document.getElementById("f_species");
  if (!select) return;
  select.disabled = !enabled;
}

function updateSubtypeDropdown(selectedType, selectedSubtype = "") {
  if (!creationOptions) return;
  setSubtypeEnabled(Boolean((selectedType || "").trim()));
  const subtypeMap = creationOptions.NpcSubtypeForTypeMap || creationOptions.npcSubtypeForTypeMap || {};
  const subtypes = subtypeMap[selectedType] || [];
  setSelectOptions("f_subtype", subtypes, true);
  setSelectValue("f_subtype", selectedSubtype);
  setRerollEnabled(Boolean((selectedSubtype || "").trim()));
}

function updateSpeciesDropdown(selectedFaction, selectedSpecies = "") {
  if (!creationOptions) return;
  setSpeciesEnabled(Boolean((selectedFaction || "").trim()));
  const speciesMap = creationOptions.NpcSpeciesForFactionMap || creationOptions.npcSpeciesForFactionMap || {};
  const species = speciesMap[selectedFaction] || [];
  setSelectOptions("f_species", species, true);
  setSelectValue("f_species", selectedSpecies);
}

async function applySubtypeRoll(subtype) {
  if (!subtype) {
    document.getElementById("f_stats").value = "";
    document.getElementById("f_items").value = "";
    document.getElementById("f_description").value = "";
    return;
  }
  const backend = appBindings();
  const rolled = await backend.RollSubtypeFields(subtype);
  document.getElementById("f_stats").value = rolled?.stats || rolled?.Stats || "";
  document.getElementById("f_items").value = rolled?.items || rolled?.Items || "";
  document.getElementById("f_description").value = rolled?.description || rolled?.Description || "";
}

async function applySpeciesNameRoll(species) {
  if (!species) {
    document.getElementById("f_name").value = "";
    return;
  }
  const backend = appBindings();
  const rolledName = await backend.RollSpeciesName(species);
  document.getElementById("f_name").value = rolledName || "";
}

function setForm(npc) {
  document.getElementById("f_id").value = npc?.ID || "";
  document.getElementById("f_name").value = compValue(npc, "1");
  const selectedType = compValue(npc, "2");
  const selectedSubtype = compValue(npc, "3");
  const selectedSpecies = compValue(npc, "4");
  const selectedFaction = compValue(npc, "5");
  setSelectValue("f_type", selectedType);
  updateSubtypeDropdown(selectedType, selectedSubtype);
  setSelectValue("f_faction", selectedFaction);
  updateSpeciesDropdown(selectedFaction, selectedSpecies);
  const traitValue = compValue(npc, "6");
  setSelectValue("f_traits", traitValue.split(",")[0]?.trim() || "");
  document.getElementById("f_stats").value = compValue(npc, "7");
  document.getElementById("f_items").value = compValue(npc, "8");
  document.getElementById("f_description").value = compValue(npc, "9");
  document.getElementById("f_location").value = npc?.LocationID || "default";
}

function clearForm() {
  document.getElementById("npcForm").reset();
  document.getElementById("f_location").value = "default";
	setSubtypeEnabled(false);
	setSpeciesEnabled(false);
	setRerollEnabled(false);
}

function npcCard(npc) {
  const id = npc?.ID || "";
  const name = compValue(npc, "1");
  const species = compValue(npc, "4");
  const type = compValue(npc, "2");
  const subtype = compValue(npc, "3");
  const faction = compValue(npc, "5");

  const element = document.createElement("div");
  element.className = "npc";
  element.innerHTML = `
    <div class="npc-header">
      <span class="npc-name">${name || "Unnamed NPC"}</span>
      <button type="button" data-id="${id}">Delete</button>
    </div>
    <div><strong>Species:</strong> ${species}</div>
    <div><strong>Type:</strong> ${type}</div>
    <div><strong>Subtype:</strong> ${subtype}</div>
    <div><strong>Faction:</strong> ${faction}</div>
  `;

  element.querySelector(".npc-name").addEventListener("click", async () => {
    const backend = appBindings();
    const loaded = await backend.GetNPC(id);
    setForm(loaded);
  });

  element.querySelector("button").addEventListener("click", async () => {
    const backend = appBindings();
    await backend.DeleteNPC(id);
    await renderList();
    if (currentId() === id) {
      clearForm();
    }
  });

  return element;
}

function populateOptionDropdowns() {
  if (!creationOptions) return;
  setSelectOptions("f_type", creationOptions.NpcTypes || creationOptions.npcTypes || [], true);
  setSelectOptions("f_faction", creationOptions.Factions || creationOptions.factions || [], true);
  setSelectOptions("f_traits", creationOptions.Traits || creationOptions.traits || [], true);
  updateSubtypeDropdown(document.getElementById("f_type").value);
  updateSpeciesDropdown(document.getElementById("f_faction").value);
}

async function renderList() {
  const backend = appBindings();
  const list = document.getElementById("list");
  list.innerHTML = "";

  const npcs = await backend.ListNPCs();
  if (!npcs || npcs.length === 0) {
    list.innerHTML = "<em>No NPCs stored yet.</em>";
    return;
  }

  npcs.forEach((npc) => {
    list.appendChild(npcCard(npc));
  });
}

async function setupActions() {
  const backend = appBindings();
  creationOptions = await backend.GetCreationOptions();
  populateOptionDropdowns();
  document.getElementById("f_type").addEventListener("change", (event) => {
    updateSubtypeDropdown(event.target.value);
    document.getElementById("f_stats").value = "";
    document.getElementById("f_items").value = "";
		document.getElementById("f_description").value = "";
		setRerollEnabled(false);
  });
  document.getElementById("f_faction").addEventListener("change", (event) => {
    updateSpeciesDropdown(event.target.value);
    document.getElementById("f_name").value = "";
  });
  document.getElementById("f_subtype").addEventListener("change", async (event) => {
		setRerollEnabled(Boolean((event.target.value || "").trim()));
    try {
      await applySubtypeRoll(event.target.value);
    } catch (error) {
      console.error(error);
      alert(error.message || "Failed to generate subtype fields.");
    }
  });
  document.getElementById("f_species").addEventListener("change", async (event) => {
    try {
      await applySpeciesNameRoll(event.target.value);
    } catch (error) {
      console.error(error);
      alert(error.message || "Failed to generate species name.");
    }
  });

  document.getElementById("btnGenerate").addEventListener("click", async () => {
    await backend.GenerateNPC();
    await renderList();
  });

  document.getElementById("btnRefresh").addEventListener("click", renderList);

  document.getElementById("btnClear").addEventListener("click", async () => {
    if (!confirm("Delete all stored NPCs?")) return;
    await backend.DeleteAllNPCs();
    await renderList();
    clearForm();
  });

  document.getElementById("btnSave").addEventListener("click", async () => {
    const payload = readForm();
    if (!payload.id) {
			alert("No ID present. Generate an NPC first.");
      return;
    }
    const saved = await backend.SaveNPC(payload);
    setForm(saved);
    await renderList();
  });
  document.getElementById("btnReroll").addEventListener("click", async () => {
    const subtype = document.getElementById("f_subtype").value;
    if (!subtype) {
      alert("Select a subtype first.");
      return;
    }
    try {
      await applySubtypeRoll(subtype);
    } catch (error) {
      console.error(error);
      alert(error.message || "Failed to reroll subtype fields.");
    }
  });

  document.getElementById("btnClose").addEventListener("click", clearForm);
  setSubtypeEnabled(Boolean((document.getElementById("f_type").value || "").trim()));
	setRerollEnabled(Boolean((document.getElementById("f_subtype").value || "").trim()));
}

async function main() {
  try {
    await setupActions();
    await renderList();
  } catch (error) {
    console.error(error);
    alert(error.message || "Failed to start UI.");
  }
}

main();
