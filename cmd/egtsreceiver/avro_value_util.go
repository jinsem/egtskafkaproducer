package main

import (
	egtsschema "github.com/jinsem/egtskafkaproducer/pkg/avro"
)

func setUnionNullDoubleVal(val float64, target *egtsschema.UnionNullDouble) {
	target.UnionType = egtsschema.UnionNullDoubleTypeEnumDouble
	target.Double = val
}

func setUnionNullInt(val int32, target *egtsschema.UnionNullInt) {
	target.UnionType = egtsschema.UnionNullIntTypeEnumInt
	target.Int = val
}

func toNavigationSystem(egtsNavSystemCode uint16) egtsschema.NavigationSystem {
	switch egtsNavSystemCode {
	// Glonass
	case 1:
		return egtsschema.NavigationSystem(egtsschema.NavigationSystemGLONASS)
	// GPS
	case 2:
		return egtsschema.NavigationSystem(egtsschema.NavigationSystemGPS)
	// Galileo
	case 4:
		return egtsschema.NavigationSystem(egtsschema.NavigationSystemGalileo)
	// Compass
	case 8:
		return egtsschema.NavigationSystem(egtsschema.NavigationSystemCompass)
	// Beidou
	case 16:
		return egtsschema.NavigationSystem(egtsschema.NavigationSystemBeidou)
	// DORIS
	case 32:
		return egtsschema.NavigationSystem(egtsschema.NavigationSystemDORIS)
	// unknown
	default: // including 0
		return egtsschema.NavigationSystem(egtsschema.NavigationSystemUknown)
	}
}
