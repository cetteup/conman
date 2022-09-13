package bf2

var voiceOverHelpLines = []string{
	"HUD_HELP_BATTLERECORDER_controlsQ",
	"HUD_HELP_BATTLERECORDER_controlsT",
	"HUD_HELP_COMMANDER_ARTILLERY_description",
	"HUD_HELP_COMMANDER_COMMUNICATION_VOIPTOGGLE",
	"HUD_HELP_COMMANDER_FUNCTIONS_FIREARTILLERY",
	"HUD_HELP_COMMANDER_FUNCTIONS_SATELLITESCAN",
	"HUD_HELP_COMMANDER_FUNCTIONS_SUPPLYDROP",
	"HUD_HELP_COMMANDER_FUNCTIONS_UAV",
	"HUD_HELP_COMMANDER_MAP_FILTERS",
	"HUD_HELP_COMMANDER_MAP_QUICKMENU",
	"HUD_HELP_COMMANDER_MAP_REQUESTS_received",
	"HUD_HELP_COMMANDER_MAP_SATELLITEVIEW_moveKeys",
	"HUD_HELP_COMMANDER_MAP_ZOOM_inOut",
	"HUD_HELP_COMMANDER_MINIMAP_description",
	"HUD_HELP_COMMANDER_SATELLITESCAN_description",
	"HUD_HELP_COMMANDER_SQUADLIST_clickSquadList",
	"HUD_HELP_COMMANDER_SQUADLIST_commoMenu",
	"HUD_HELP_COMMANDER_SQUADLIST_orderIcon",
	"HUD_HELP_COMMANDER_SQUADLIST_squadList",
	"HUD_HELP_COMMANDER_SUPPLIES_description",
	"HUD_HELP_COMMANDER_UAV_description",
	"HUD_HELP_COMMANDER_UAVplaced",
	"HUD_HELP_COMMANDER_VOIP_allSquadLeadersChannel",
	"HUD_HELP_COMMANDER_artilleryDamaged",
	"HUD_HELP_COMMANDER_artilleryFired",
	"HUD_HELP_COMMANDER_commanderApply",
	"HUD_HELP_COMMANDER_placeArtillery",
	"HUD_HELP_COMMANDER_placeSupplies",
	"HUD_HELP_COMMANDER_placeUAV",
	"HUD_HELP_COMMANDER_radarDamaged",
	"HUD_HELP_COMMANDER_uavDamaged",
	"HUD_HELP_GAMEPLAY_captureFlag",
	"HUD_HELP_GAMEPLAY_climbLadder",
	"HUD_HELP_KIT_ENGINEER_NAMETAG_indicatorBar",
	"HUD_HELP_KIT_ENGINEER_inVehicle",
	"HUD_HELP_KIT_MEDIC_NAMETAG_indicatorBar",
	"HUD_HELP_KIT_MEDIC_inVehicle",
	"HUD_HELP_KIT_SUPPORT_NAMETAG_indicatorBar",
	"HUD_HELP_KIT_SUPPORT_inVehicle",
	"HUD_HELP_MENU_COMMOROSE_slCommoRoseFirstUse",
	"HUD_HELP_MENU_CONTROLS_3dMap",
	"HUD_HELP_MENU_CONTROLS_CAPSLOCK_orderMapToggle",
	"HUD_HELP_MENU_CONTROLS_COMMOROSE_SQUADLEADER_spawnScreenShortcut",
	"HUD_HELP_MENU_CONTROLS_COMMOROSE_generalCommoRose",
	"HUD_HELP_MENU_CONTROLS_COMMOROSE_slcommoRoseAvailable",
	"HUD_HELP_MENU_CONTROLS_COMMOROSE_smCommoRoseAvailable",
	"HUD_HELP_MENU_CONTROLS_COMMOROSE_specificSpotted",
	"HUD_HELP_MENU_CONTROLS_callMedicWhenManDown",
	"HUD_HELP_MENU_CONTROLS_changeWeapon",
	"HUD_HELP_MENU_CONTROLS_largeMap",
	"HUD_HELP_MENU_CONTROLS_scoreBoard",
	"HUD_HELP_MENU_CONTROLS_spawnScreen",
	"HUD_HELP_MENU_CONTROLS_zoomMap",
	"HUD_HELP_MENU_applyForCommanderOrSquadHere",
	"HUD_HELP_PLAYER_CONTROLS_crouch",
	"HUD_HELP_PLAYER_CONTROLS_fireLeftMouse",
	"HUD_HELP_PLAYER_CONTROLS_jump",
	"HUD_HELP_PLAYER_CONTROLS_lieDown",
	"HUD_HELP_PLAYER_CONTROLS_moveForwardsBackwards",
	"HUD_HELP_PLAYER_CONTROLS_moveMouseToLook",
	"HUD_HELP_PLAYER_CONTROLS_parachute",
	"HUD_HELP_PLAYER_CONTROLS_radioMessages",
	"HUD_HELP_PLAYER_CONTROLS_sprint",
	"HUD_HELP_PLAYER_CONTROLS_strafeLeftRight",
	"HUD_HELP_PLAYER_CONTROLS_zoomAltFire",
	"HUD_HELP_SPAWNSCREEN_KITTAB_pickArmyAndKit",
	"HUD_HELP_SPAWNSCREEN_KITTAB_pickSpawnPoint",
	"HUD_HELP_SPAWNSCREEN_KITTAB_trySquads",
	"HUD_HELP_SPAWNSCREEN_KITTAB_useUnlock",
	"HUD_HELP_SPAWNSCREEN_SQUADTAB_applyCommander",
	"HUD_HELP_SPAWNSCREEN_SQUADTAB_createSquad",
	"HUD_HELP_SPAWNSCREEN_SQUADTAB_joinSquad",
	"HUD_HELP_SPAWNSCREEN_spawnOnSL",
	"HUD_HELP_SQUAD_LEADER_VOIP_commanderChannel",
	"HUD_HELP_SQUAD_LEADER_mainMapRightMouseMenu",
	"HUD_HELP_SQUAD_callForTarget",
	"HUD_HELP_SQUAD_recievedTarget",
	"HUD_HELP_SQUAD_squadCreated",
	"HUD_HELP_SQUAD_squadJoined",
	"HUD_HELP_VEHICLE_AA_DRIVER_CONTROLS_fireSecondaryWeapon",
	"HUD_HELP_VEHICLE_APC_DRIVER_CONTROLS_fireSecondaryWeapon",
	"HUD_HELP_VEHICLE_APC_PASSENGER_CONTROLS_firePassengerWeapon",
	"HUD_HELP_VEHICLE_CONTROLS_exitVehicle",
	"HUD_HELP_VEHICLE_F35_PILOT_CONTROLS_engageHoverEngines",
	"HUD_HELP_VEHICLE_GENERAL_CONTROLS_enterVehicle",
	"HUD_HELP_VEHICLE_HELO_LOWHEALTH_useFriendlyRepairStation",
	"HUD_HELP_VEHICLE_HELO_PASSENGER_CONTROLS_TVGUIDED_clickToGuide",
	"HUD_HELP_VEHICLE_HELO_PASSENGER_CONTROLS_switchToGuidedMissiles",
	"HUD_HELP_VEHICLE_HELO_PILOT_CONTROLS_deployFlares",
	"HUD_HELP_VEHICLE_JET_GUNNER_CONTROLS_LASERGUIDED_clickToFire",
	"HUD_HELP_VEHICLE_JET_GUNNER_CONTROLS_switchToLaserMissiles",
	"HUD_HELP_VEHICLE_JET_LOWAMMO_flyOverFriendlyAirField",
	"HUD_HELP_VEHICLE_JET_LOWHEALTH_useFriendlyRepairStation",
	"HUD_HELP_VEHICLE_JET_PILOT_CONTROLS_toggleWeaponsShortcut",
	"HUD_HELP_VEHICLE_LCAC_DRIVER_CONTROLS_raiseLowerRamp",
	"HUD_HELP_VEHICLE_TANK_DRIVER_CONTROLS_smoke",
	"HUD_HELP_VEHICLE_TANK_TURRET_CONTROLS_duckInTurret",
	"HUD_HELP_WEAPON_HANDHELD_AMMOBAG_CONTROLS_throwAmmoBagDown",
	"HUD_HELP_WEAPON_HANDHELD_ASSAULTRIFLE_CONTROLS_switchToGrenades",
	"HUD_HELP_WEAPON_HANDHELD_ASSAULTRIFLE_CONTROLS_switchToRifle",
	"HUD_HELP_WEAPON_HANDHELD_ASSAULTRIFLE_CONTROLS_toggleFireModes",
	"HUD_HELP_WEAPON_HANDHELD_ATMINE_vehicles",
	"HUD_HELP_WEAPON_HANDHELD_AT_CONTROLS_guideMissile",
	"HUD_HELP_WEAPON_HANDHELD_C4_CONTROLS_switchToDetonator",
	"HUD_HELP_WEAPON_HANDHELD_CLAYMORE_GENERAL_mineHasDirectionalBlast",
	"HUD_HELP_WEAPON_HANDHELD_HANDGRENADE_CONTROLS_rollGrenade",
	"HUD_HELP_WEAPON_HANDHELD_HEALBAG_GENERAL_holdHealBagToHealLocally",
	"HUD_HELP_WEAPON_HANDHELD_LMG_GENERAL_watchTempGauge",
	"HUD_HELP_WEAPON_HANDHELD_SHOCKPADDLES_CONTROLS_reviveTeamMates",
	"HUD_HELP_WEAPON_HANDHELD_SIMRAD_CONTROLS_transmitTargetCoordinates",
	"HUD_HELP_WEAPON_HANDHELD_SNIPERRIFLE_CONTROLS_switchToScopeView",
	"HUD_HELP_WEAPON_HANDHELD_SUBMACHINEGUN_CONTROLS_toggleFireModes",
	"HUD_HELP_WEAPON_HANDHELD_WRENCH_fixStuff",
	"HUD_HELP_WEAPON_STATIONARY_AA_CONTROLS_lockOnTone",
	"HUD_HELP_WEAPON_STATIONARY_AT_CONTROLS_guideMissile",
	"HUD_HELP_WEAPON_STATIONARY_CONTROLS_enterStationaryWeapon",
	"HUD_HELP_WEPAON_HANDHELD_GENERAL_reloadWeapon",
	"HUD_HELP_WORLD_PLAYER_ARTILLERY_enemy",
	"HUD_HELP_WORLD_PLAYER_ARTILLERY_friendly",
	"HUD_HELP_WORLD_PLAYER_CRATE_description",
	"HUD_HELP_WORLD_PLAYER_RADAR_enemy",
	"HUD_HELP_WORLD_PLAYER_RADAR_friendly",
	"HUD_HELP_WORLD_PLAYER_UAV_TRAILER_enemy",
	"HUD_HELP_WORLD_PLAYER_UAV_TRAILER_friendly",
}
