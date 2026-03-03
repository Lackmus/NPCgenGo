import { initNPCUI } from "./npc-ui-core.js";

function appBindings() {
  const backend = window?.go?.main?.WailsAPI || window?.go?.main?.App;
  if (!backend) {
    throw new Error("Wails bindings not available. Expected window.go.main.WailsAPI (or legacy App).");
  }
  return backend;
}

const wailsAPI = {
  async listNPCs() {
    return appBindings().ListNPCs();
  },

  async getNPC(id) {
    return appBindings().GetNPC(id);
  },

  async getCreationOptions() {
    return appBindings().GetCreationOptions();
  },

  async rollSubtypeFields(subtype) {
    return appBindings().RollSubtypeFields(subtype);
  },

  async rollSpeciesName(species) {
    return appBindings().RollSpeciesName(species);
  },

  async generateNPC() {
    return appBindings().GenerateNPC();
  },

  async deleteNPC(id) {
    await appBindings().DeleteNPC(id);
  },

  async deleteAllNPCs() {
    await appBindings().DeleteAllNPCs();
  },

  async saveNPC(payload) {
    return appBindings().SaveNPC(payload);
  },

  async createNPC(payload) {
    return appBindings().CreateNPC(payload);
  },
};

async function main() {
  try {
    const app = initNPCUI(wailsAPI, window);
    await app.start();
  } catch (error) {
    console.error(error);
    alert(error.message || "Failed to start UI.");
  }
}

main();
