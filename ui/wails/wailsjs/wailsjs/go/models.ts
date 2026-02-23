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
	    description: string;
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
	        this.description = source["description"];
	        this.locationID = source["locationID"];
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

