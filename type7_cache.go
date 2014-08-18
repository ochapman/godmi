/*
* File Name:	type7_cache.go
* Description:	
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2014-08-19
*/
package godmi

type OperationalMode byte

const (
	OperationalModeWriteThrough OperationalMode = iota
	OperationalModeWriteBack
	OperationalModeVariesWithMemoryAddress
	OperationalModeUnknown
)

func (o OperationalMode) String() string {
	modes := [...]string{
		"Write Through",
		"Write Back",
		"Varies With Memory Address",
		"Unknown",
	}
	return modes[o]
}

type CacheLocation byte

const (
	CacheLocationInternal CacheLocation = iota
	CacheLocationExternal
	CacheLocationReserved
	CacheLocationUnknown
)

func (c CacheLocation) String() string {
	locations := [...]string{
		"Internal",
		"External",
		"Reserved",
		"Unknown",
	}
	return locations[c]
}

type CacheLevel byte

const (
	Level1 CacheLevel = iota
	Level2
	Level3
)

func (c CacheLevel) String() string {
	levels := [...]string{
		"Level1",
		"Level2",
		"Level3",
	}
	return levels[c]
}

type CacheConfiguration struct {
	Mode     OperationalMode
	Enabled  bool
	Location CacheLocation
	Socketed bool
	Level    CacheLevel
}

func NewCacheConfiguration(u uint16) CacheConfiguration {
	var c CacheConfiguration
	c.Level = CacheLevel(byte(u & 0x7))
	c.Socketed = (u&0x10 == 1)
	c.Location = CacheLocation((u >> 5) & 0x3)
	c.Enabled = (u&(0x1<<7) == 1)
	c.Mode = OperationalMode((u >> 8) & 0x7)
	return c
}

func (c CacheConfiguration) String() string {
	return fmt.Sprintf("Cache Configuration:\n"+
		"\tLevel: %s\n"+
		"\tSocketed: %v\n"+
		"\tLocation: %s\n"+
		"\tEnabled: %v\n"+
		"\tMode:\n\t\t",
		c.Level,
		c.Socketed,
		c.Location,
		c.Enabled,
		c.Mode)
}

type CacheGranularity byte

const (
	CacheGranularity1K CacheGranularity = iota
	CacheGranularity64K
)

func (c CacheGranularity) String() string {
	grans := [...]string{
		"1K",
		"64K",
	}
	return grans[c]
}

type CacheSize struct {
	Granularity CacheGranularity
	Size        uint16
}

func NewCacheSize(u uint16) CacheSize {
	var c CacheSize
	c.Granularity = CacheGranularity(u >> 15)
	c.Size = u &^ (uint16(1) << 15)
	return c
}

func (c CacheSize) String() string {
	return fmt.Sprintf("%s * %s", c.Size, c.Granularity)
}

type SRAMType uint16

const (
	SRAMTypeOther SRAMType = 1 << iota
	SRAMTypeUnknown
	SRAMTypeNonBurst
	SRAMTypeBurst
	SRAMTypePipelineBurst
	SRAMTypeSynchronous
	SRAMTypeAsynchronous
	SRAMTypeReserved
)

func (s SRAMType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"Non-Burst",
		"Burst",
		"Pipeline Burst",
		"Synchronous",
		"Asynchronous",
		"Reserved",
	}
	return types[s/2]
}

type CacheSpeed byte

type ErrorCorrectionType byte

const (
	ErrorCorrectionTypeOther ErrorCorrectionType = 1 + iota
	ErrorCorrectionTypeUnknown
	ErrorCorrectionTypeNone
	ErrorCorrectionTypeParity
	ErrorCorrectionTypeSinglebitECC
	ErrorCorrectionTypeMultibitECC
)

func (e ErrorCorrectionType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"None",
		"Parity",
		"Single-bit ECC",
		"Multi-bit ECC",
	}
	return types[e-1]
}

type SystemCacheType byte

const (
	SystemCacheTypeOther SystemCacheType = 1 + iota
	SystemCacheTypeUnknown
	SystemCacheTypeInstruction
	SystemCacheTypeData
	SystemCacheTypeUnified
)

func (s SystemCacheType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"Instruction",
		"Data",
		"Unified",
	}
	return types[s-1]
}

type CacheAssociativity byte

const (
	CacheAssociativityOther CacheAssociativity = 1 + iota
	CacheAssociativityUnknown
	CacheAssociativityDirectMapped
	CacheAssociativity2waySetAssociative
	CacheAssociativity4waySetAssociative
	CacheAssociativityFullyAssociative
	CacheAssociativity8waySetAssociative
	CacheAssociativity16waySetAssociative
	CacheAssociativity12waySetAssociative
	CacheAssociativity24waySetAssociative
	CacheAssociativity32waySetAssociative
	CacheAssociativity48waySetAssociative
	CacheAssociativity64waySetAssociative
	CacheAssociativity20waySetAssociative
)

func (a CacheAssociativity) String() string {
	caches := [...]string{
		"Other",
		"Unknown",
		"Direct Mapped",
		"2-way Set-Associative",
		"4-way Set-Associative",
		"Fully Associative",
		"8-way Set-Associative",
		"16-way Set-Associative",
		"12-way Set-Associative",
		"24-way Set-Associative",
		"32-way Set-Associative",
		"48-way Set-Associative",
		"64-way Set-Associative",
		"20-way Set-Associative",
	}
	return caches[a]
}

type CacheInformation struct {
	infoCommon
	SocketDesignation   string
	Configuration       CacheConfiguration
	MaximumCacheSize    CacheSize
	InstalledSize       CacheSize
	SupportedSRAMType   SRAMType
	CurrentSRAMType     SRAMType
	CacheSpeed          CacheSpeed
	ErrorCorrectionType ErrorCorrectionType
	SystemCacheType     SystemCacheType
	Associativity       CacheAssociativity
}
