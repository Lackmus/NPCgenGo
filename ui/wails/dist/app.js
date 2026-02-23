function appBindings() {
  const backend = window?.go?.main?.WailsAPI || window?.go?.main?.App;
  if (!backend) {
    throw new Error("Wails bindings not available. Expected window.go.main.WailsAPI (or legacy App).");
  }
  return backend;
}

function compValue(npc, key) {
  if (!npc || !npc.Components) return "";
  return npc.Components[key] || "";
}

function currentId() {
  return document.getElementById("f_id").value;
}

function readForm() {
  return {
    id: document.getElementById("f_id").value,
    name: document.getElementById("f_name").value,
    type: document.getElementById("f_type").value,
    subtype: document.getElementById("f_subtype").value,
    species: document.getElementById("f_species").value,
    faction: document.getElementById("f_faction").value,
    traits: document.getElementById("f_traits").value.split(",").map((value) => value.trim()).filter(Boolean),
    stats: document.getElementById("f_stats").value,
    items: document.getElementById("f_items").value,
    description: document.getElementById("f_description").value,
    locationID: document.getElementById("f_location").value || "default",
  };
}

function setForm(npc) {
  document.getElementById("f_id").value = npc?.ID || "";
  document.getElementById("f_name").value = compValue(npc, "1");
  document.getElementById("f_type").value = compValue(npc, "2");
  document.getElementById("f_subtype").value = compValue(npc, "3");
  document.getElementById("f_species").value = compValue(npc, "4");
  document.getElementById("f_faction").value = compValue(npc, "5");
  document.getElementById("f_traits").value = compValue(npc, "6");
  document.getElementById("f_stats").value = compValue(npc, "7");
  document.getElementById("f_items").value = compValue(npc, "8");
  document.getElementById("f_description").value = compValue(npc, "9");
  document.getElementById("f_location").value = npc?.LocationID || "default";
}

function clearForm() {
  document.getElementById("npcForm").reset();
  document.getElementById("f_location").value = "default";
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
      alert("No ID present. Use Create New first.");
      return;
    }
    const saved = await backend.SaveNPC(payload);
    setForm(saved);
    await renderList();
  });

  document.getElementById("btnCreate").addEventListener("click", async () => {
    const payload = readForm();
    const created = await backend.CreateNPC(payload);
    setForm(created);
    await renderList();
  });

  document.getElementById("btnClose").addEventListener("click", clearForm);
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
