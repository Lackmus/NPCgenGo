export namespace main {
	
	export class NPCInput {
	    id: string;
	    name: string;
	    type: string;
	    subtype: string;
	    species: string;
	    faction: string;
	    traits: string[];
	    stats: string;
	    items: string;
	    locationID: string;
	
	    static createFrom(source: any = {}) {
	        return new NPCInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.subtype = source["subtype"];
	        this.species = source["species"];
	        this.faction = source["faction"];
	        this.traits = source["traits"];
	        this.stats = source["stats"];
	        this.items = source["items"];
	        this.locationID = source["locationID"];
	    }
	}
	export class SubtypeRoll {
	    stats: string;
	    items: string;
	
	    static createFrom(source: any = {}) {
	        return new SubtypeRoll(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stats = source["stats"];
	        this.items = source["items"];
	    }
	}

}

export namespace model {
	
	export class NPC {
	    ID: string;
	    LocationID: string;
	    Components: Record<number, string>;
	
	    static createFrom(source: any = {}) {
	        return new NPC(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.LocationID = source["LocationID"];
	        this.Components = source["Components"];
	    }
	}

}

export namespace service {
	
	export class NPCCreationOptions {
	    Factions: string[];
	    Species: string[];
	    Traits: string[];
	    NpcTypes: string[];
	    NpcSubtypeForTypeMap: Record<string, Array<string>>;
	    NpcSpeciesForFactionMap: Record<string, Array<string>>;
	
	    static createFrom(source: any = {}) {
	        return new NPCCreationOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Factions = source["Factions"];
	        this.Species = source["Species"];
	        this.Traits = source["Traits"];
	        this.NpcTypes = source["NpcTypes"];
	        this.NpcSubtypeForTypeMap = source["NpcSubtypeForTypeMap"];
	        this.NpcSpeciesForFactionMap = source["NpcSpeciesForFactionMap"];
	    }
	}

}

