package mcdb




const (
	keySubChunkData = '/' 
)


const (
	
	keyVersion = ',' 
	
	
	keyVersionOld = 'v' 
	
	
	keyBlockEntities = '1' 
	
	
	keyEntitiesOld = '2' 
	
	
	keyPendingScheduledTicks = '3'
	
	
	
	keyFinalisation = '6' 
	
	key3DData = '+' 
	
	
	key2DData = '-' 
	
	
	keyChecksums = ';' 

	keyEntityIdentifiers = "digp"

	keyEntity = "actorprefix"
)


const (
	keyAutonomousEntities = "AutonomousEntities"
	keyOverworld          = "Overworld"
	keyMobEvents          = "mobevents"
	keyBiomeData          = "BiomeData"
	keyScoreboard         = "scoreboard"
	keyLocalPlayer        = "~local_player"
)

const (
	finalisationGenerated = iota + 1
	finalisationPopulated
)
