function byId(id) {
  return document.getElementById(id);
}

function textOrDash(value) {
  const normalized = String(value ?? "").trim();
  return normalized || "—";
}

function isPresent(value) {
  return String(value ?? "").trim().length > 0;
}

function npcId(npc) {
  return npc?.id || npc?.ID || "";
}

function setSelectOptions(selectId, values, includeEmpty = true) {
  const select = byId(selectId);
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
  const select = byId(selectId);
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

function setButtonEnabled(id, enabled) {
  const button = byId(id);
  if (!button) return;
  button.disabled = !enabled;
}

function setFieldEnabled(id, enabled) {
  const field = byId(id);
  if (!field) return;
  field.disabled = !enabled;
}

function showDetailsPanel() {
  const details = byId("npcDetails");
  const form = byId("npcForm");
  if (details) details.style.display = "grid";
  if (form) form.style.display = "none";
}

function showEditPanel() {
  const details = byId("npcDetails");
  const form = byId("npcForm");
  if (details) details.style.display = "none";
  if (form) form.style.display = "grid";
}

function renderDetails(npc) {
  byId("d_name").textContent = textOrDash(npc?.name);
  byId("d_type").textContent = textOrDash(npc?.type);
  byId("d_subtype").textContent = textOrDash(npc?.subtype);
  byId("d_species").textContent = textOrDash(npc?.species);
  byId("d_faction").textContent = textOrDash(npc?.faction);
  byId("d_traits").textContent = textOrDash(npc?.trait);
  byId("d_stats").textContent = textOrDash(npc?.stats);
  byId("d_items").textContent = textOrDash(npc?.items);
  setButtonEnabled("btnEdit", isPresent(npcId(npc)));
}

function readForm() {
  const traitValue = byId("f_traits")?.value || "";
  const statsValue = (byId("f_stats")?.textContent || "").trim();
  const itemsValue = (byId("f_items")?.textContent || "").trim();

  return {
    id: byId("f_id")?.value || "",
    name: byId("f_name")?.value || "",
    type: byId("f_type")?.value || "",
    subtype: byId("f_subtype")?.value || "",
    species: byId("f_species")?.value || "",
    faction: byId("f_faction")?.value || "",
    trait: traitValue,
    stats: statsValue === "—" ? "" : statsValue,
    items: itemsValue === "—" ? "" : itemsValue,
  };
}

function validatePayload(payload) {
  const checks = [
    ["name", payload.name],
    ["type", payload.type],
    ["subtype", payload.subtype],
    ["species", payload.species],
    ["faction", payload.faction],
    ["trait", payload.trait],
  ];

  const missing = checks
    .filter(([, value]) => !isPresent(value))
    .map(([label]) => label);

  return {
    ok: missing.length === 0,
    missing,
  };
}

export function initNPCUI(api, dialog = window) {
  let creationOptions = null;
  let selectedNPC = null;

  function updateSubtypeDropdown(selectedType, selectedSubtype = "") {
    if (!creationOptions) return;
    setFieldEnabled("f_subtype", isPresent(selectedType));
    const subtypeMap = creationOptions.NpcSubtypeForTypeMap || creationOptions.npcSubtypeForTypeMap || {};
    const subtypes = subtypeMap[selectedType] || [];
    setSelectOptions("f_subtype", subtypes, true);
    setSelectValue("f_subtype", selectedSubtype);
    setButtonEnabled("btnReroll", isPresent(selectedSubtype));
  }

  function updateSpeciesDropdown(selectedFaction, selectedSpecies = "") {
    if (!creationOptions) return;
    setFieldEnabled("f_species", isPresent(selectedFaction));
    const speciesMap = creationOptions.NpcSpeciesForFactionMap || creationOptions.npcSpeciesForFactionMap || {};
    const species = speciesMap[selectedFaction] || [];
    setSelectOptions("f_species", species, true);
    setSelectValue("f_species", selectedSpecies);
    setButtonEnabled("btnRerollName", isPresent(selectedSpecies));
  }

  async function applySubtypeRoll(subtype) {
    if (!isPresent(subtype)) {
      byId("f_stats").textContent = "—";
      byId("f_items").textContent = "—";
      return;
    }
    const rolled = await api.rollSubtypeFields(subtype);
    byId("f_stats").textContent = rolled?.stats || rolled?.Stats || "—";
    byId("f_items").textContent = rolled?.items || rolled?.Items || "—";
  }

  async function applySpeciesNameRoll(species) {
    if (!isPresent(species)) {
      byId("f_name").value = "";
      return;
    }
    const rolledName = await api.rollSpeciesName(species);
    byId("f_name").value = rolledName || "";
  }

  function populateOptionDropdowns() {
    if (!creationOptions) return;
    setSelectOptions("f_type", creationOptions.NpcTypes || creationOptions.npcTypes || [], true);
    setSelectOptions("f_faction", creationOptions.Factions || creationOptions.factions || [], true);
    setSelectOptions("f_traits", creationOptions.Traits || creationOptions.traits || [], true);
    updateSubtypeDropdown(byId("f_type").value);
    updateSpeciesDropdown(byId("f_faction").value);
  }

  function setForm(npc) {
    selectedNPC = npc || null;
    byId("f_id").value = npcId(npc);
    byId("f_name").value = npc?.name || "";
    const selectedType = npc?.type || "";
    const selectedSubtype = npc?.subtype || "";
    const selectedSpecies = npc?.species || "";
    const selectedFaction = npc?.faction || "";
    setSelectValue("f_type", selectedType);
    updateSubtypeDropdown(selectedType, selectedSubtype);
    setSelectValue("f_faction", selectedFaction);
    updateSpeciesDropdown(selectedFaction, selectedSpecies);
    const traitValue = npc?.trait || "";
    setSelectValue("f_traits", traitValue.split(",")[0]?.trim() || "");
    byId("f_stats").textContent = npc?.stats || "—";
    byId("f_items").textContent = npc?.items || "—";
    setButtonEnabled("btnReroll", isPresent(selectedSubtype));
    setButtonEnabled("btnRerollName", isPresent(selectedSpecies));
  }

  function clearForm() {
    selectedNPC = null;
    byId("npcForm")?.reset();
    setFieldEnabled("f_subtype", false);
    setFieldEnabled("f_species", false);
    setButtonEnabled("btnReroll", false);
    setButtonEnabled("btnRerollName", false);
    renderDetails(null);
    showDetailsPanel();
  }

  function startCreateNPC() {
    selectedNPC = null;
    byId("f_id").value = "";
    byId("f_name").value = "";
    setSelectValue("f_type", "");
    updateSubtypeDropdown("", "");
    setSelectValue("f_faction", "");
    updateSpeciesDropdown("", "");
    setSelectValue("f_traits", "");
    byId("f_stats").textContent = "—";
    byId("f_items").textContent = "—";
    setButtonEnabled("btnReroll", false);
    setButtonEnabled("btnRerollName", false);
    renderDetails(null);
    showEditPanel();
  }

  function npcCard(npc) {
    const id = npcId(npc);
    const name = npc?.name || "";
    const species = npc?.species || "";
    const type = npc?.type || "";
    const subtype = npc?.subtype || "";
    const faction = npc?.faction || "";

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

    element.querySelector(".npc-name")?.addEventListener("click", async () => {
      const loaded = await api.getNPC(id);
      setForm(loaded);
      renderDetails(loaded);
      showDetailsPanel();
    });

    element.querySelector("button")?.addEventListener("click", async () => {
      await api.deleteNPC(id);
      await renderList();
      if ((byId("f_id")?.value || "") === id) {
        clearForm();
      }
    });

    return element;
  }

  async function renderList() {
    const list = byId("list");
    list.innerHTML = "";

    const npcs = await api.listNPCs();
    if (!npcs || npcs.length === 0) {
      list.innerHTML = "<em>No NPCs stored yet.</em>";
      return;
    }

    npcs.forEach((npc) => {
      list.appendChild(npcCard(npc));
    });
  }

  async function setupActions() {
    creationOptions = await api.getCreationOptions();
    populateOptionDropdowns();

    byId("f_type")?.addEventListener("change", (event) => {
      updateSubtypeDropdown(event.target.value);
      byId("f_stats").textContent = "—";
      byId("f_items").textContent = "—";
      setButtonEnabled("btnReroll", false);
    });

    byId("f_faction")?.addEventListener("change", (event) => {
      updateSpeciesDropdown(event.target.value);
      byId("f_name").value = "";
    });

    byId("f_subtype")?.addEventListener("change", async (event) => {
      setButtonEnabled("btnReroll", isPresent(event.target.value));
      try {
        await applySubtypeRoll(event.target.value);
      } catch (error) {
        console.error(error);
        dialog.alert(error.message || "Failed to generate subtype fields.");
      }
    });

    byId("f_species")?.addEventListener("change", async (event) => {
      setButtonEnabled("btnRerollName", isPresent(event.target.value));
      try {
        await applySpeciesNameRoll(event.target.value);
      } catch (error) {
        console.error(error);
        dialog.alert(error.message || "Failed to generate species name.");
      }
    });

    byId("btnReroll")?.addEventListener("click", async () => {
      const subtype = byId("f_subtype")?.value || "";
      if (!isPresent(subtype)) {
        dialog.alert("Select a subtype first.");
        return;
      }
      try {
        await applySubtypeRoll(subtype);
      } catch (error) {
        console.error(error);
        dialog.alert(error.message || "Failed to reroll subtype fields.");
      }
    });

    byId("btnRerollName")?.addEventListener("click", async () => {
      const species = byId("f_species")?.value || "";
      if (!isPresent(species)) {
        dialog.alert("Select a species first.");
        return;
      }
      try {
        await applySpeciesNameRoll(species);
      } catch (error) {
        console.error(error);
        dialog.alert(error.message || "Failed to reroll name.");
      }
    });

    byId("btnGenerate")?.addEventListener("click", async () => {
      await api.generateNPC();
      await renderList();
    });

    byId("btnRefresh")?.addEventListener("click", renderList);

    byId("btnClear")?.addEventListener("click", async () => {
      if (!dialog.confirm("Delete all stored NPCs?")) return;
      await api.deleteAllNPCs();
      await renderList();
      clearForm();
    });

    byId("btnCreate")?.addEventListener("click", startCreateNPC);

    byId("btnSave")?.addEventListener("click", async () => {
      const payload = readForm();
      const validation = validatePayload(payload);
      if (!validation.ok) {
        dialog.alert(`Please fill all fields before saving. Missing: ${validation.missing.join(", ")}`);
        return;
      }

      let saved;
      if (isPresent(payload.id)) {
        saved = await api.saveNPC(payload);
      } else if (typeof api.createNPC === "function") {
        saved = await api.createNPC(payload);
      } else {
        dialog.alert("No ID present. Generate an NPC first.");
        return;
      }

      setForm(saved || payload);
      renderDetails(saved || payload);
      showDetailsPanel();
      await renderList();
    });

    byId("btnEdit")?.addEventListener("click", () => {
      if (!isPresent(byId("f_id")?.value || "")) {
        dialog.alert("Select an NPC first.");
        return;
      }
      showEditPanel();
    });

    byId("btnCancelEdit")?.addEventListener("click", () => {
      setForm(selectedNPC);
      renderDetails(selectedNPC);
      showDetailsPanel();
    });

    byId("btnClose")?.addEventListener("click", clearForm);

    setFieldEnabled("f_subtype", isPresent(byId("f_type")?.value || ""));
    setFieldEnabled("f_species", isPresent(byId("f_faction")?.value || ""));
    setButtonEnabled("btnReroll", isPresent(byId("f_subtype")?.value || ""));
    setButtonEnabled("btnRerollName", isPresent(byId("f_species")?.value || ""));
    renderDetails(null);
    showDetailsPanel();
  }

  return {
    async start() {
      await setupActions();
      await renderList();
    },
  };
}
