import { initNPCUI } from "/ui-shared/npc-ui-core.js";

const API_BASE = "/api";

const webAPI = {
  async listNPCs() {
    const response = await fetch(`${API_BASE}/npcs`);
    if (!response.ok) return [];
    return response.json();
  },

  async getNPC(id) {
    const response = await fetch(`${API_BASE}/npcs/${encodeURIComponent(id)}`);
    if (!response.ok) {
      throw new Error(`Failed to load NPC: ${response.statusText}`);
    }
    return response.json();
  },

  async getCreationOptions() {
    const response = await fetch(`${API_BASE}/options`);
    if (!response.ok) {
      throw new Error(`Failed to load options: ${response.statusText}`);
    }
    return response.json();
  },

  async rollSubtypeFields(subtype) {
    const response = await fetch(`${API_BASE}/subtypes/${encodeURIComponent(subtype)}/roll`);
    if (!response.ok) {
      throw new Error(`Failed to roll subtype fields: ${response.statusText}`);
    }
    return response.json();
  },

  async rollSpeciesName(species) {
    const response = await fetch(`${API_BASE}/species/${encodeURIComponent(species)}/name`);
    if (!response.ok) {
      throw new Error(`Failed to roll species name: ${response.statusText}`);
    }
    const payload = await response.json();
    return payload?.name || payload?.Name || "";
  },

  async generateNPC() {
    const response = await fetch(`${API_BASE}/generate`, { method: "POST" });
    if (!response.ok) {
      throw new Error(`Generate failed: ${response.statusText}`);
    }
    return response.json();
  },

  async deleteNPC(id) {
    const response = await fetch(`${API_BASE}/npcs/${encodeURIComponent(id)}`, { method: "DELETE" });
    if (!response.ok && response.status !== 204) {
      throw new Error(`Delete failed: ${response.statusText}`);
    }
  },

  async deleteAllNPCs() {
    const npcs = await this.listNPCs();
    await Promise.all(npcs.map((npc) => this.deleteNPC(npc?.ID || npc?.id || "")));
  },

  async saveNPC(payload) {
    const response = await fetch(`${API_BASE}/npcs/${encodeURIComponent(payload.id)}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    if (!response.ok) {
      throw new Error(`Save failed: ${response.statusText}`);
    }
    return this.getNPC(payload.id);
  },
};

async function main() {
  try {
    const app = initNPCUI(webAPI, window);
    await app.start();
  } catch (error) {
    console.error(error);
    alert(error.message || "Failed to start UI.");
  }
}

main();
