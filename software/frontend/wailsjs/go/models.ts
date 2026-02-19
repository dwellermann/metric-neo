export namespace application {
	
	export class ChronoConfigDTO {
	    enabled: boolean;
	    port: string;
	    baudRate: number;
	    autoRecord: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ChronoConfigDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.port = source["port"];
	        this.baudRate = source["baudRate"];
	        this.autoRecord = source["autoRecord"];
	    }
	}
	export class ShotDTO {
	    velocityMPS: number;
	    energyJoules: number;
	    timestamp: string;
	    valid: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ShotDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.velocityMPS = source["velocityMPS"];
	        this.energyJoules = source["energyJoules"];
	        this.timestamp = source["timestamp"];
	        this.valid = source["valid"];
	    }
	}
	export class ProjectileDTO {
	    id: string;
	    name: string;
	    weightGrams: number;
	    bc: number;
	
	    static createFrom(source: any = {}) {
	        return new ProjectileDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.weightGrams = source["weightGrams"];
	        this.bc = source["bc"];
	    }
	}
	export class OpticDTO {
	    type: string;
	    modelName: string;
	    weightG: number;
	    minMagnification: number;
	    maxMagnification: number;
	
	    static createFrom(source: any = {}) {
	        return new OpticDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.modelName = source["modelName"];
	        this.weightG = source["weightG"];
	        this.minMagnification = source["minMagnification"];
	        this.maxMagnification = source["maxMagnification"];
	    }
	}
	export class ProfileDTO {
	    id: string;
	    name: string;
	    category: string;
	    barrelLengthMM: number;
	    triggerWeightG: number;
	    sightHeightMM: number;
	    optic?: OpticDTO;
	    opticID?: string;
	    twistRateMM?: number;
	    defaultAmmoID?: string;
	    totalWeightG: number;
	
	    static createFrom(source: any = {}) {
	        return new ProfileDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.category = source["category"];
	        this.barrelLengthMM = source["barrelLengthMM"];
	        this.triggerWeightG = source["triggerWeightG"];
	        this.sightHeightMM = source["sightHeightMM"];
	        this.optic = this.convertValues(source["optic"], OpticDTO);
	        this.opticID = source["opticID"];
	        this.twistRateMM = source["twistRateMM"];
	        this.defaultAmmoID = source["defaultAmmoID"];
	        this.totalWeightG = source["totalWeightG"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SessionDTO {
	    id: string;
	    profileSnapshot: ProfileDTO;
	    projectileSnapshot: ProjectileDTO;
	    shots: ShotDTO[];
	    temperatureCelsius?: number;
	    note: string;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new SessionDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.profileSnapshot = this.convertValues(source["profileSnapshot"], ProfileDTO);
	        this.projectileSnapshot = this.convertValues(source["projectileSnapshot"], ProjectileDTO);
	        this.shots = this.convertValues(source["shots"], ShotDTO);
	        this.temperatureCelsius = source["temperatureCelsius"];
	        this.note = source["note"];
	        this.createdAt = source["createdAt"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ChronoPollResultDTO {
	    recorded: boolean;
	    velocityMPS?: number;
	    session?: SessionDTO;
	
	    static createFrom(source: any = {}) {
	        return new ChronoPollResultDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.recorded = source["recorded"];
	        this.velocityMPS = source["velocityMPS"];
	        this.session = this.convertValues(source["session"], SessionDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	export class Result___metric_neo_internal_application_ProfileDTO_ {
	    data: ProfileDTO[];
	    error: string;
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Result___metric_neo_internal_application_ProfileDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], ProfileDTO);
	        this.error = source["error"];
	        this.success = source["success"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Result___metric_neo_internal_application_ProjectileDTO_ {
	    data: ProjectileDTO[];
	    error: string;
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Result___metric_neo_internal_application_ProjectileDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], ProjectileDTO);
	        this.error = source["error"];
	        this.success = source["success"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SessionMetaDTO {
	    id: string;
	    profileName: string;
	    projectileName: string;
	    shotCount: number;
	    validShotCount: number;
	    createdAt: string;
	    note: string;
	    avgVelocityMPS?: number;
	    avgEnergyJoules?: number;
	
	    static createFrom(source: any = {}) {
	        return new SessionMetaDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.profileName = source["profileName"];
	        this.projectileName = source["projectileName"];
	        this.shotCount = source["shotCount"];
	        this.validShotCount = source["validShotCount"];
	        this.createdAt = source["createdAt"];
	        this.note = source["note"];
	        this.avgVelocityMPS = source["avgVelocityMPS"];
	        this.avgEnergyJoules = source["avgEnergyJoules"];
	    }
	}
	export class Result___metric_neo_internal_application_SessionMetaDTO_ {
	    data: SessionMetaDTO[];
	    error: string;
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Result___metric_neo_internal_application_SessionMetaDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], SessionMetaDTO);
	        this.error = source["error"];
	        this.success = source["success"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SightDTO {
	    id: string;
	    type: string;
	    modelName: string;
	    weightG: number;
	    minMagnification: number;
	    maxMagnification: number;
	
	    static createFrom(source: any = {}) {
	        return new SightDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.modelName = source["modelName"];
	        this.weightG = source["weightG"];
	        this.minMagnification = source["minMagnification"];
	        this.maxMagnification = source["maxMagnification"];
	    }
	}
	export class Result___metric_neo_internal_application_SightDTO_ {
	    data: SightDTO[];
	    error: string;
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Result___metric_neo_internal_application_SightDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], SightDTO);
	        this.error = source["error"];
	        this.success = source["success"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Result_bool_ {
	    data: boolean;
	    error: string;
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Result_bool_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = source["data"];
	        this.error = source["error"];
	        this.success = source["success"];
	    }
	}
	export class Result_metric_neo_internal_application_ChronoConfigDTO_ {
	    data: ChronoConfigDTO;
	    error: string;
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Result_metric_neo_internal_application_ChronoConfigDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], ChronoConfigDTO);
	        this.error = source["error"];
	        this.success = source["success"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Result_metric_neo_internal_application_ChronoPollResultDTO_ {
	    data: ChronoPollResultDTO;
	    error: string;
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Result_metric_neo_internal_application_ChronoPollResultDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], ChronoPollResultDTO);
	        this.error = source["error"];
	        this.success = source["success"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Result_metric_neo_internal_application_ProfileDTO_ {
	    data: ProfileDTO;
	    error: string;
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Result_metric_neo_internal_application_ProfileDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], ProfileDTO);
	        this.error = source["error"];
	        this.success = source["success"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Result_metric_neo_internal_application_ProjectileDTO_ {
	    data: ProjectileDTO;
	    error: string;
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Result_metric_neo_internal_application_ProjectileDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], ProjectileDTO);
	        this.error = source["error"];
	        this.success = source["success"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Result_metric_neo_internal_application_SessionDTO_ {
	    data: SessionDTO;
	    error: string;
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Result_metric_neo_internal_application_SessionDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], SessionDTO);
	        this.error = source["error"];
	        this.success = source["success"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Result_metric_neo_internal_application_SightDTO_ {
	    data: SightDTO;
	    error: string;
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Result_metric_neo_internal_application_SightDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], SightDTO);
	        this.error = source["error"];
	        this.success = source["success"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StatisticsDTO {
	    avgVelocityMPS: number;
	    standardDeviation: number;
	    minVelocityMPS: number;
	    maxVelocityMPS: number;
	    extremeSpread: number;
	    avgEnergyJoules: number;
	    validShotCount: number;
	    totalShotCount: number;
	
	    static createFrom(source: any = {}) {
	        return new StatisticsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.avgVelocityMPS = source["avgVelocityMPS"];
	        this.standardDeviation = source["standardDeviation"];
	        this.minVelocityMPS = source["minVelocityMPS"];
	        this.maxVelocityMPS = source["maxVelocityMPS"];
	        this.extremeSpread = source["extremeSpread"];
	        this.avgEnergyJoules = source["avgEnergyJoules"];
	        this.validShotCount = source["validShotCount"];
	        this.totalShotCount = source["totalShotCount"];
	    }
	}
	export class Result_metric_neo_internal_application_StatisticsDTO_ {
	    data: StatisticsDTO;
	    error: string;
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Result_metric_neo_internal_application_StatisticsDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], StatisticsDTO);
	        this.error = source["error"];
	        this.success = source["success"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Result_string_ {
	    data: string;
	    error: string;
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Result_string_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = source["data"];
	        this.error = source["error"];
	        this.success = source["success"];
	    }
	}
	
	
	
	

}

