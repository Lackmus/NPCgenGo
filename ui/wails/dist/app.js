function appBindings() {
  const backend = window?.go?.main?.WailsAPI || window?.go?.main?.App;
  if (!backend) {
    throw new Error("Wails bindings not available. Expected window.go.main.WailsAPI (or legacy App).");
  }
  return backend;
}

let creationOptions = null;
let selectedNPC = null;

function compValue(npc, key) {
  if (!npc || !npc.Components) return "";
  return npc.Components[key] || "";
}

function currentId() {
  return document.getElementById("f_id").value;
}

function readForm() {
	const traitValue = document.getElementById("f_traits").value;
	const statsValue = (document.getElementById("f_stats")?.textContent || "").trim();
	const itemsValue = (document.getElementById("f_items")?.textContent || "").trim();
  return {
    id: document.getElementById("f_id").value,
    name: document.getElementById("f_name").value,
    type: document.getElementById("f_type").value,
    subtype: document.getElementById("f_subtype").value,
    species: document.getElementById("f_species").value,
    faction: document.getElementById("f_faction").value,
    traits: traitValue ? [traitValue] : [],
    stats: statsValue === "—" ? "" : statsValue,
    items: itemsValue === "—" ? "" : itemsValue,
    locationID: document.getElementById("f_location").value || "default",
  };
}

function validatePayload(payload) {
  const checks = [
    ["name", payload.name],
    ["type", payload.type],
    ["subtype", payload.subtype],
    ["species", payload.species],
    ["faction", payload.faction],
    ["traits", Array.isArray(payload.traits) ? payload.traits.join(",") : payload.traits],
    ["locationID", payload.locationID],
  ];

  const missing = checks
    .filter(([, value]) => !String(value || "").trim())
    .map(([label]) => label);

  return {
    ok: missing.length === 0,
    missing,
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

function setDetailValue(id, value) {
  const element = document.getElementById(id);
  if (!element) return;
  element.textContent = (value || "").trim() || "—";
}

function renderDetails(npc) {
  setDetailValue("d_name", compValue(npc, "1"));
  setDetailValue("d_type", compValue(npc, "2"));
  setDetailValue("d_subtype", compValue(npc, "3"));
  setDetailValue("d_species", compValue(npc, "4"));
  setDetailValue("d_faction", compValue(npc, "5"));
  setDetailValue("d_traits", compValue(npc, "6"));
  setDetailValue("d_stats", compValue(npc, "7"));
  setDetailValue("d_items", compValue(npc, "8"));
  setDetailValue("d_location", npc?.LocationID || "");

  const editButton = document.getElementById("btnEdit");
  if (editButton) {
    editButton.disabled = !(npc?.ID || "");
  }
}

function showDetailsPanel() {
  const details = document.getElementById("npcDetails");
  const form = document.getElementById("npcForm");
  if (details) details.style.display = "grid";
  if (form) form.style.display = "none";
}

function showEditPanel() {
  const details = document.getElementById("npcDetails");
  const form = document.getElementById("npcForm");
  if (details) details.style.display = "none";
  if (form) form.style.display = "grid";
}

function enterEditMode() {
  if (!currentId()) {
    alert("Select an NPC first.");
    return;
  }
  showEditPanel();
}

function exitEditMode() {
  showDetailsPanel();
}

function startCreateNPC() {
  selectedNPC = null;
  document.getElementById("f_id").value = "";
  document.getElementById("f_name").value = "";
  setSelectValue("f_type", "");
  updateSubtypeDropdown("", "");
  setSelectValue("f_faction", "");
  updateSpeciesDropdown("", "");
  setSelectValue("f_traits", "");
  document.getElementById("f_stats").textContent = "—";
  document.getElementById("f_items").textContent = "—";
  document.getElementById("f_location").value = "default";
  setRerollEnabled(false);
  renderDetails(null);
  showEditPanel();
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
    document.getElementById("f_stats").textContent = "—";
    document.getElementById("f_items").textContent = "—";
    return;
  }
  const backend = appBindings();
  const rolled = await backend.RollSubtypeFields(subtype);
  document.getElementById("f_stats").textContent = rolled?.stats || rolled?.Stats || "—";
  document.getElementById("f_items").textContent = rolled?.items || rolled?.Items || "—";
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
  selectedNPC = npc || null;
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
  document.getElementById("f_stats").textContent = compValue(npc, "7") || "—";
  document.getElementById("f_items").textContent = compValue(npc, "8") || "—";
  document.getElementById("f_location").value = npc?.LocationID || "default";

  const id = npc?.ID || "";
  const editButton = document.getElementById("btnEdit");
  if (editButton) {
    editButton.disabled = !id;
  }
}

function clearForm() {
  selectedNPC = null;
  document.getElementById("npcForm").reset();
  document.getElementById("f_location").value = "default";
	setSubtypeEnabled(false);
	setSpeciesEnabled(false);
	setRerollEnabled(false);
  renderDetails(null);
  exitEditMode();
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
    renderDetails(loaded);
    exitEditMode();
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
    document.getElementById("f_stats").textContent = "—";
    document.getElementById("f_items").textContent = "—";

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

  const createButton = document.getElementById("btnCreate");
  if (createButton) {
    createButton.addEventListener("click", startCreateNPC);
  }

  document.getElementById("btnRefresh").addEventListener("click", renderList);

  document.getElementById("btnClear").addEventListener("click", async () => {
    if (!confirm("Delete all stored NPCs?")) return;
    await backend.DeleteAllNPCs();
    await renderList();
    clearForm();
  });

  document.getElementById("btnSave").addEventListener("click", async () => {
    const payload = readForm();
    const validation = validatePayload(payload);
    if (!validation.ok) {
      alert(`Please fill all fields before saving. Missing: ${validation.missing.join(", ")}`);
      return;
    }
    const saved = payload.id ? await backend.SaveNPC(payload) : await backend.CreateNPC(payload);
    setForm(saved);
    renderDetails(saved);
    exitEditMode();
    await renderList();
  });
  const editButton = document.getElementById("btnEdit");
  if (editButton) {
    editButton.addEventListener("click", enterEditMode);
  }
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

  document.getElementById("btnCancelEdit").addEventListener("click", () => {
    setForm(selectedNPC);
    renderDetails(selectedNPC);
    exitEditMode();
  });
  document.getElementById("btnClose").addEventListener("click", clearForm);

  const btn = document.getElementById("btnEdit");
  if (btn) {
    btn.disabled = true;
  }
  setSubtypeEnabled(Boolean((document.getElementById("f_type").value || "").trim()));
	setSpeciesEnabled(Boolean((document.getElementById("f_faction").value || "").trim()));
	setRerollEnabled(Boolean((document.getElementById("f_subtype").value || "").trim()));
  renderDetails(null);
  exitEditMode();
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
