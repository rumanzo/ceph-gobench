package main

import (
	"github.com/ceph/go-ceph/rados"
	"time"
)

type params struct {
	duration                                                                              time.Duration
	threadsCount                                                                          uint64
	blocksize, objectsize                                                                 uint64
	parallel, disablecheck                                                                bool
	bs, os, cluster, user, keyring, config, pool, define, rdefine, cpuprofile, memprofile string
}

type cephconnection struct {
	conn  *rados.Conn
	ioctx *rados.IOContext
}

type Poolinfo struct {
	Pool   string `json:"pool,omitempty"`
	PoolId uint64 `json:"pool_id,omitempty"`
	Size   uint64 `json:"size,omitempty"`
}

func (times *PlacementGroup) StringsToTimes() {
	const LongForm = "2006-01-02 15:04:05.000000"
	times.LastFreshT, _ = time.Parse(LongForm, times.LastFresh)
	times.LastChangeT, _ = time.Parse(LongForm, times.LastChange)
	times.LastActiveT, _ = time.Parse(LongForm, times.LastActive)
	times.LastPeeredT, _ = time.Parse(LongForm, times.LastPeered)
	times.LastCleanT, _ = time.Parse(LongForm, times.LastClean)
	times.LastBecameActiveT, _ = time.Parse(LongForm, times.LastBecameActive)
	times.LastBecamePeeredT, _ = time.Parse(LongForm, times.LastBecamePeered)
	times.LastUnstaleT, _ = time.Parse(LongForm, times.LastUnstale)
	times.LastUndegradedT, _ = time.Parse(LongForm, times.LastUndegraded)
	times.LastFullsizedT, _ = time.Parse(LongForm, times.LastFullsized)
	times.LastDeepScrubStampT, _ = time.Parse(LongForm, times.LastDeepScrubStamp)
	times.LastDeepScrubT, _ = time.Parse(LongForm, times.LastDeepScrub)
	times.LastCleanScrubStampT, _ = time.Parse(LongForm, times.LastCleanScrubStamp)
	times.LastScrubStampT, _ = time.Parse(LongForm, times.LastScrubStamp)
	times.LastScrubT, _ = time.Parse(LongForm, times.LastScrub)
}

type Bucket struct {
	Alg   string `json:"alg"`
	Hash  string `json:"hash"`
	ID    int64  `json:"id"`
	Items []struct {
		ID     int64   `json:"id"`
		Pos    uint64  `json:"pos"`
		Weight float64 `json:"weight"`
	} `json:"items"`
	Name     string  `json:"name"`
	TypeID   int64   `json:"type_id"`
	TypeName string  `json:"type_name"`
	Weight   float64 `json:"weight"`
}

type Device struct {
	Class string `json:"class"`
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Info  OsdMetadata
}

type OsdCrushDump struct {
	Buckets    []Bucket `json:"buckets"`
	ChooseArgs struct{} `json:"choose_args"`
	Devices    []Device `json:"devices"`
	Rules      []struct {
		MaxSize  uint64 `json:"max_size"`
		MinSize  uint64 `json:"min_size"`
		RuleID   int64  `json:"rule_id"`
		RuleName string `json:"rule_name"`
		Ruleset  int64  `json:"ruleset"`
		Steps    []struct {
			Item     int64  `json:"item"`
			ItemName string `json:"item_name"`
			Num      int64  `json:"num"`
			Op       string `json:"op"`
			Type     string `json:"type"`
		} `json:"steps"`
		Type uint64 `json:"type"`
	} `json:"rules"`
	Tunables struct {
		AllowedBucketAlgs        uint64 `json:"allowed_bucket_algs"`
		ChooseLocalFallbackTries uint64 `json:"choose_local_fallback_tries"`
		ChooseLocalTries         uint64 `json:"choose_local_tries"`
		ChooseTotalTries         uint64 `json:"choose_total_tries"`
		ChooseleafDescendOnce    uint64 `json:"chooseleaf_descend_once"`
		ChooseleafStable         uint64 `json:"chooseleaf_stable"`
		ChooseleafVaryR          uint64 `json:"chooseleaf_vary_r"`
		HasV2Rules               uint64 `json:"has_v2_rules"`
		HasV3Rules               uint64 `json:"has_v3_rules"`
		HasV4Buckets             uint64 `json:"has_v4_buckets"`
		HasV5Rules               uint64 `json:"has_v5_rules"`
		LegacyTunables           uint64 `json:"legacy_tunables"`
		MinimumRequiredVersion   string `json:"minimum_required_version"`
		OptimalTunables          uint64 `json:"optimal_tunables"`
		Profile                  string `json:"profile"`
		RequireFeatureTunables   uint64 `json:"require_feature_tunables"`
		RequireFeatureTunables2  uint64 `json:"require_feature_tunables2"`
		RequireFeatureTunables3  uint64 `json:"require_feature_tunables3"`
		RequireFeatureTunables5  uint64 `json:"require_feature_tunables5"`
		StrawCalcVersion         uint64 `json:"straw_calc_version"`
	} `json:"tunables"`
	Types []struct {
		Name   string `json:"name"`
		TypeID int64  `json:"type_id"`
	} `json:"types"`
}
type Osd struct {
	ClusterAddr        string   `json:"cluster_addr"`
	DownAt             uint64   `json:"down_at"`
	HeartbeatBackAddr  string   `json:"heartbeat_back_addr"`
	HeartbeatFrontAddr string   `json:"heartbeat_front_addr"`
	In                 uint64   `json:"in"`
	LastCleanBegin     uint64   `json:"last_clean_begin"`
	LastCleanEnd       uint64   `json:"last_clean_end"`
	LostAt             uint64   `json:"lost_at"`
	Osd                uint64   `json:"osd"`
	PrimaryAffinity    float64  `json:"primary_affinity"`
	PublicAddr         string   `json:"public_addr"`
	State              []string `json:"state"`
	Up                 uint64   `json:"up"`
	UpFrom             uint64   `json:"up_from"`
	UpThru             uint64   `json:"up_thru"`
	UUID               string   `json:"uuid"`
	Weight             float64  `json:"weight"`
}
type OsdDump struct {
	BackfillfullRatio   float64  `json:"backfillfull_ratio"`
	Blacklist           struct{} `json:"blacklist"`
	ClusterSnapshot     string   `json:"cluster_snapshot"`
	Created             string   `json:"created"`
	CrushVersion        uint64   `json:"crush_version"`
	Epoch               uint64   `json:"epoch"`
	ErasureCodeProfiles struct {
		Default struct {
			K         string `json:"k"`
			M         string `json:"m"`
			Plugin    string `json:"plugin"`
			Technique string `json:"technique"`
		} `json:"default"`
	} `json:"erasure_code_profiles"`
	Flags           string        `json:"flags"`
	FlagsNum        uint64        `json:"flags_num"`
	FlagsSet        []string      `json:"flags_set"`
	Fsid            string        `json:"fsid"`
	FullRatio       float64       `json:"full_ratio"`
	MaxOsd          uint64        `json:"max_osd"`
	MinCompatClient string        `json:"min_compat_client"`
	Modified        string        `json:"modified"`
	NearfullRatio   float64       `json:"nearfull_ratio"`
	NewPurgedSnaps  []interface{} `json:"new_purged_snaps"`
	NewRemovedSnaps []interface{} `json:"new_removed_snaps"`
	OsdXinfo        []struct {
		DownStamp        string  `json:"down_stamp"`
		Features         uint64  `json:"features"`
		LaggyInterval    uint64  `json:"laggy_interval"`
		LaggyProbability float64 `json:"laggy_probability"`
		OldWeight        float64 `json:"old_weight"`
		Osd              uint64  `json:"osd"`
	} `json:"osd_xinfo"`
	Osds         []Osd         `json:"osds"`
	PgTemp       []interface{} `json:"pg_temp"`
	PgUpmap      []interface{} `json:"pg_upmap"`
	PgUpmapItems []interface{} `json:"pg_upmap_items"`
	PoolMax      uint64        `json:"pool_max"`
	Pools        []struct {
		ApplicationMetadata struct {
			Rbd struct{} `json:"rbd"`
			Rgw struct{} `json:"rgw"`
		} `json:"application_metadata"`
		Auid                           uint64        `json:"auid"`
		CacheMinEvictAge               uint64        `json:"cache_min_evict_age"`
		CacheMinFlushAge               uint64        `json:"cache_min_flush_age"`
		CacheMode                      string        `json:"cache_mode"`
		CacheTargetDirtyHighRatioMicro uint64        `json:"cache_target_dirty_high_ratio_micro"`
		CacheTargetDirtyRatioMicro     uint64        `json:"cache_target_dirty_ratio_micro"`
		CacheTargetFullRatioMicro      uint64        `json:"cache_target_full_ratio_micro"`
		CreateTime                     string        `json:"create_time"`
		CrushRule                      int64         `json:"crush_rule"`
		ErasureCodeProfile             string        `json:"erasure_code_profile"`
		ExpectedNumObjects             uint64        `json:"expected_num_objects"`
		FastRead                       bool          `json:"fast_read"`
		Flags                          uint64        `json:"flags"`
		FlagsNames                     string        `json:"flags_names"`
		GradeTable                     []interface{} `json:"grade_table"`
		HitSetCount                    uint64        `json:"hit_set_count"`
		HitSetGradeDecayRate           uint64        `json:"hit_set_grade_decay_rate"`
		HitSetParams                   struct {
			Type string `json:"type"`
		} `json:"hit_set_params"`
		HitSetPeriod                 uint64        `json:"hit_set_period"`
		HitSetSearchLastN            uint64        `json:"hit_set_search_last_n"`
		LastChange                   string        `json:"last_change"`
		LastForceOpResend            string        `json:"last_force_op_resend"`
		LastForceOpResendPreluminous string        `json:"last_force_op_resend_preluminous"`
		MinReadRecencyForPromote     uint64        `json:"min_read_recency_for_promote"`
		MinSize                      uint64        `json:"min_size"`
		MinWriteRecencyForPromote    uint64        `json:"min_write_recency_for_promote"`
		ObjectHash                   uint64        `json:"object_hash"`
		Options                      struct{}      `json:"options"`
		PgNum                        uint64        `json:"pg_num"`
		PgPlacementNum               uint64        `json:"pg_placement_num"`
		Pool                         uint64        `json:"pool"`
		PoolName                     string        `json:"pool_name"`
		PoolSnaps                    []interface{} `json:"pool_snaps"`
		QuotaMaxBytes                uint64        `json:"quota_max_bytes"`
		QuotaMaxObjects              uint64        `json:"quota_max_objects"`
		ReadTier                     int64         `json:"read_tier"`
		RemovedSnaps                 string        `json:"removed_snaps"`
		Size                         uint64        `json:"size"`
		SnapEpoch                    uint64        `json:"snap_epoch"`
		SnapMode                     string        `json:"snap_mode"`
		SnapSeq                      uint64        `json:"snap_seq"`
		StripeWidth                  uint64        `json:"stripe_width"`
		TargetMaxBytes               uint64        `json:"target_max_bytes"`
		TargetMaxObjects             uint64        `json:"target_max_objects"`
		TierOf                       int64         `json:"tier_of"`
		Tiers                        []interface{} `json:"tiers"`
		Type                         uint64        `json:"type"`
		UseGmtHitset                 bool          `json:"use_gmt_hitset"`
		WriteTier                    int64         `json:"write_tier"`
	} `json:"pools"`
	PrimaryTemp            []interface{} `json:"primary_temp"`
	RemovedSnapsQueue      []interface{} `json:"removed_snaps_queue"`
	RequireMinCompatClient string        `json:"require_min_compat_client"`
	RequireOsdRelease      string        `json:"require_osd_release"`
}

type PlacementGroup struct {
	Acting                  []uint64      `json:"acting"`
	ActingPrimary           int64         `json:"acting_primary"`
	AvailNoMissing          []any  `json:"avail_no_missing"`
	BlockedBy               []interface{} `json:"blocked_by"`
	Created                 uint64        `json:"created"`
	DirtyStatsInvalid       bool          `json:"dirty_stats_invalid"`
	HitsetBytesStatsInvalid bool          `json:"hitset_bytes_stats_invalid"`
	HitsetStatsInvalid      bool          `json:"hitset_stats_invalid"`
	LastActive              string        `json:"last_active"`
	LastActiveT             time.Time
	LastBecameActive        string `json:"last_became_active"`
	LastBecameActiveT       time.Time
	LastBecamePeered        string `json:"last_became_peered"`
	LastBecamePeeredT       time.Time
	LastChange              string `json:"last_change"`
	LastChangeT             time.Time
	LastClean               string `json:"last_clean"`
	LastCleanT              time.Time
	LastCleanScrubStamp     string `json:"last_clean_scrub_stamp"`
	LastCleanScrubStampT    time.Time
	LastDeepScrub           string `json:"last_deep_scrub"`
	LastDeepScrubT          time.Time
	LastDeepScrubStamp      string `json:"last_deep_scrub_stamp"`
	LastDeepScrubStampT     time.Time
	LastEpochClean          uint64 `json:"last_epoch_clean"`
	LastFresh               string `json:"last_fresh"`
	LastFreshT              time.Time
	LastFullsized           string `json:"last_fullsized"`
	LastFullsizedT          time.Time
	LastPeered              string `json:"last_peered"`
	LastPeeredT             time.Time
	LastScrub               string `json:"last_scrub"`
	LastScrubT              time.Time
	LastScrubStamp          string `json:"last_scrub_stamp"`
	LastScrubStampT         time.Time
	LastUndegraded          string `json:"last_undegraded"`
	LastUndegradedT         time.Time
	LastUnstale             string `json:"last_unstale"`
	LastUnstaleT            time.Time
	LogSize                 uint64        `json:"log_size"`
	LogStart                string        `json:"log_start"`
	ObjectLocationCounts    []any  `json:"object_location_counts"`
	ManifestStatsInvalid    bool          `json:"manifest_stats_invalid"`
	MappingEpoch            uint64        `json:"mapping_epoch"`
	OmapStatsInvalid        bool          `json:"omap_stats_invalid"`
	OndiskLogSize           uint64        `json:"ondisk_log_size"`
	OndiskLogStart          string        `json:"ondisk_log_start"`
	Parent                  string        `json:"parent"`
	ParentSplitBits         uint64        `json:"parent_split_bits"`
	Pgid                    string        `json:"pgid"`
	PinStatsInvalid         bool          `json:"pin_stats_invalid"`
	PurgedSnaps             []interface{} `json:"purged_snaps"`
	ReportedEpoch           string        `json:"reported_epoch"`
	ReportedSeq             string        `json:"reported_seq"`
	SnaptrimqLen            uint64        `json:"snaptrimq_len"`
	StatSum                 struct {
		NumBytes                   uint64 `json:"num_bytes"`
		NumBytesHitSetArchive      uint64 `json:"num_bytes_hit_set_archive"`
		NumBytesRecovered          uint64 `json:"num_bytes_recovered"`
		NumDeepScrubErrors         uint64 `json:"num_deep_scrub_errors"`
		NumEvict                   uint64 `json:"num_evict"`
		NumEvictKb                 uint64 `json:"num_evict_kb"`
		NumEvictModeFull           uint64 `json:"num_evict_mode_full"`
		NumEvictModeSome           uint64 `json:"num_evict_mode_some"`
		NumFlush                   uint64 `json:"num_flush"`
		NumFlushKb                 uint64 `json:"num_flush_kb"`
		NumFlushModeHigh           uint64 `json:"num_flush_mode_high"`
		NumFlushModeLow            uint64 `json:"num_flush_mode_low"`
		NumKeysRecovered           uint64 `json:"num_keys_recovered"`
		NumLargeOmapObjects        uint64 `json:"num_large_omap_objects"`
		NumLegacySnapsets          uint64 `json:"num_legacy_snapsets"`
		NumObjectClones            uint64 `json:"num_object_clones"`
		NumObjectCopies            uint64 `json:"num_object_copies"`
		NumObjects                 uint64 `json:"num_objects"`
		NumObjectsDegraded         uint64 `json:"num_objects_degraded"`
		NumObjectsDirty            uint64 `json:"num_objects_dirty"`
		NumObjectsHitSetArchive    uint64 `json:"num_objects_hit_set_archive"`
		NumObjectsManifest         uint64 `json:"num_objects_manifest"`
		NumObjectsMisplaced        uint64 `json:"num_objects_misplaced"`
		NumObjectsMissing          uint64 `json:"num_objects_missing"`
		NumObjectsMissingOnPrimary uint64 `json:"num_objects_missing_on_primary"`
		NumObjectsOmap             uint64 `json:"num_objects_omap"`
		NumObjectsPinned           uint64 `json:"num_objects_pinned"`
		NumObjectsRecovered        uint64 `json:"num_objects_recovered"`
		NumObjectsRepaired         uint64 `json:"num_objects_repaired"`
		NumObjectsUnfound          uint64 `json:"num_objects_unfound"`
		NumOmapBytes               uint64 `json:"num_omap_bytes"`
		NumOmapKeys                uint64 `json:"num_omap_keys"`
		NumPromote                 uint64 `json:"num_promote"`
		NumRead                    uint64 `json:"num_read"`
		NumReadKb                  uint64 `json:"num_read_kb"`
		NumScrubErrors             uint64 `json:"num_scrub_errors"`
		NumShallowScrubErrors      uint64 `json:"num_shallow_scrub_errors"`
		NumWhiteouts               uint64 `json:"num_whiteouts"`
		NumWrite                   uint64 `json:"num_write"`
		NumWriteKb                 uint64 `json:"num_write_kb"`
	} `json:"stat_sum"`
	State        string   `json:"state"`
	StatsInvalid bool     `json:"stats_invalid"`
	Up           []uint64 `json:"up"`
	UpPrimary    int64    `json:"up_primary"`
	Version      string   `json:"version"`
}

type OsdMetadata struct {
	Arch                       string `json:"arch"`
	BackAddr                   string `json:"back_addr"`
	BackIface                  string `json:"back_iface"`
	Bluefs                     string `json:"bluefs"`
	BluefsSingleSharedDevice   string `json:"bluefs_single_shared_device"`
	BluestoreBdevAccessMode    string `json:"bluestore_bdev_access_mode"`
	BluestoreBdevBlockSize     string `json:"bluestore_bdev_block_size"`
	BluestoreBdevDev           string `json:"bluestore_bdev_dev"`
	BluestoreBdevDevNode       string `json:"bluestore_bdev_dev_node"`
	BluestoreBdevDriver        string `json:"bluestore_bdev_driver"`
	BluestoreBdevModel         string `json:"bluestore_bdev_model"`
	BluestoreBdevPartitionPath string `json:"bluestore_bdev_partition_path"`
	BluestoreBdevRotational    string `json:"bluestore_bdev_rotational"`
	BluestoreBdevSize          string `json:"bluestore_bdev_size"`
	BluestoreBdevType          string `json:"bluestore_bdev_type"`
	CephRelease                string `json:"ceph_release"`
	CephVersion                string `json:"ceph_version"`
	CephVersionShort           string `json:"ceph_version_short"`
	CPU                        string `json:"cpu"`
	DefaultDeviceClass         string `json:"default_device_class"`
	Devices                    string `json:"devices"`
	Distro                     string `json:"distro"`
	DistroDescription          string `json:"distro_description"`
	DistroVersion              string `json:"distro_version"`
	FrontAddr                  string `json:"front_addr"`
	FrontIface                 string `json:"front_iface"`
	HbBackAddr                 string `json:"hb_back_addr"`
	HbFrontAddr                string `json:"hb_front_addr"`
	Hostname                   string `json:"hostname"`
	ID                         int64  `json:"id"`
	JournalRotational          string `json:"journal_rotational"`
	KernelDescription          string `json:"kernel_description"`
	KernelVersion              string `json:"kernel_version"`
	MemSwapKb                  string `json:"mem_swap_kb"`
	MemTotalKb                 string `json:"mem_total_kb"`
	Os                         string `json:"os"`
	OsdData                    string `json:"osd_data"`
	OsdObjectstore             string `json:"osd_objectstore"`
	Rotational                 string `json:"rotational"`
}

type OsdMap struct {
	Acting        []uint64 `json:"acting"`
	ActingPrimary uint64   `json:"acting_primary"`
	Epoch         uint64   `json:"epoch"`
	Objname       string   `json:"objname"`
	Pgid          string   `json:"pgid"`
	Pool          string   `json:"pool"`
	PoolID        uint64   `json:"pool_id"`
	RawPgid       string   `json:"raw_pgid"`
	Up            []uint64 `json:"up"`
	UpPrimary     int64    `json:"up_primary"`
}

type avgLatencies struct {
	latencytotal int64
	len          int64
}

type placementGroupNautilus struct {
	PgReady bool             `json:"pg_ready"`
	PgStats []PlacementGroup `json:"pg_stats"`
}

type osdStatLine struct {
	num  int64
	line string
}
