export namespace main {
	
	export class NPCInput {
	    id: string;
	    name: string;
	    type: string;
	    subtype: string;
	    species: string;
	    faction: string;
	    trait: string;
	    stats: string;
	    items: string;
	
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
	        this.trait = source["trait"];
	        this.stats = source["stats"];
	        this.items = source["items"];
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
	    Components: Record<number, string>;
	
	    static createFrom(source: any = {}) {
	        return new NPC(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
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

