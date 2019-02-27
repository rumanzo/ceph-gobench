package main

import (
	"github.com/ceph/go-ceph/rados"
	"time"
)

type Params struct {
	duration                                                   time.Duration
	threadsCount                                               int64
	blocksize, objectsize                                      int64
	parallel                                                   bool
	bs, os, cluster, user, keyring, config, pool, mode, define string
}

type Cephconnection struct {
	conn  *rados.Conn
	ioctx *rados.IOContext
}

type Poolinfo struct {
	Pool   string `json:"pool,omitempty"`
	PoolId int64  `json:"pool_id,omitempty"`
	Size   int64  `json:"size,omitempty"`
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
		Pos    int64   `json:"pos"`
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
		MaxSize  int64  `json:"max_size"`
		MinSize  int64  `json:"min_size"`
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
		Type int64 `json:"type"`
	} `json:"rules"`
	Tunables struct {
		AllowedBucketAlgs        int64  `json:"allowed_bucket_algs"`
		ChooseLocalFallbackTries int64  `json:"choose_local_fallback_tries"`
		ChooseLocalTries         int64  `json:"choose_local_tries"`
		ChooseTotalTries         int64  `json:"choose_total_tries"`
		ChooseleafDescendOnce    int64  `json:"chooseleaf_descend_once"`
		ChooseleafStable         int64  `json:"chooseleaf_stable"`
		ChooseleafVaryR          int64  `json:"chooseleaf_vary_r"`
		HasV2Rules               int64  `json:"has_v2_rules"`
		HasV3Rules               int64  `json:"has_v3_rules"`
		HasV4Buckets             int64  `json:"has_v4_buckets"`
		HasV5Rules               int64  `json:"has_v5_rules"`
		LegacyTunables           int64  `json:"legacy_tunables"`
		MinimumRequiredVersion   string `json:"minimum_required_version"`
		OptimalTunables          int64  `json:"optimal_tunables"`
		Profile                  string `json:"profile"`
		RequireFeatureTunables   int64  `json:"require_feature_tunables"`
		RequireFeatureTunables2  int64  `json:"require_feature_tunables2"`
		RequireFeatureTunables3  int64  `json:"require_feature_tunables3"`
		RequireFeatureTunables5  int64  `json:"require_feature_tunables5"`
		StrawCalcVersion         int64  `json:"straw_calc_version"`
	} `json:"tunables"`
	Types []struct {
		Name   string `json:"name"`
		TypeID int64  `json:"type_id"`
	} `json:"types"`
}

type OsdDump struct {
	BackfillfullRatio   float64  `json:"backfillfull_ratio"`
	Blacklist           struct{} `json:"blacklist"`
	ClusterSnapshot     string   `json:"cluster_snapshot"`
	Created             string   `json:"created"`
	CrushVersion        int64    `json:"crush_version"`
	Epoch               int64    `json:"epoch"`
	ErasureCodeProfiles struct {
		Default struct {
			K         string `json:"k"`
			M         string `json:"m"`
			Plugin    string `json:"plugin"`
			Technique string `json:"technique"`
		} `json:"default"`
	} `json:"erasure_code_profiles"`
	Flags           string        `json:"flags"`
	FlagsNum        int64         `json:"flags_num"`
	FlagsSet        []string      `json:"flags_set"`
	Fsid            string        `json:"fsid"`
	FullRatio       float64       `json:"full_ratio"`
	MaxOsd          int64         `json:"max_osd"`
	MinCompatClient string        `json:"min_compat_client"`
	Modified        string        `json:"modified"`
	NearfullRatio   float64       `json:"nearfull_ratio"`
	NewPurgedSnaps  []interface{} `json:"new_purged_snaps"`
	NewRemovedSnaps []interface{} `json:"new_removed_snaps"`
	OsdXinfo        []struct {
		DownStamp        string  `json:"down_stamp"`
		Features         int64   `json:"features"`
		LaggyInterval    int64   `json:"laggy_interval"`
		LaggyProbability float64 `json:"laggy_probability"`
		OldWeight        float64 `json:"old_weight"`
		Osd              int64   `json:"osd"`
	} `json:"osd_xinfo"`
	Osds []struct {
		ClusterAddr        string   `json:"cluster_addr"`
		DownAt             int64    `json:"down_at"`
		HeartbeatBackAddr  string   `json:"heartbeat_back_addr"`
		HeartbeatFrontAddr string   `json:"heartbeat_front_addr"`
		In                 int64    `json:"in"`
		LastCleanBegin     int64    `json:"last_clean_begin"`
		LastCleanEnd       int64    `json:"last_clean_end"`
		LostAt             int64    `json:"lost_at"`
		Osd                int64    `json:"osd"`
		PrimaryAffinity    float64  `json:"primary_affinity"`
		PublicAddr         string   `json:"public_addr"`
		State              []string `json:"state"`
		Up                 int64    `json:"up"`
		UpFrom             int64    `json:"up_from"`
		UpThru             int64    `json:"up_thru"`
		UUID               string   `json:"uuid"`
		Weight             float64  `json:"weight"`
	} `json:"osds"`
	PgTemp       []interface{} `json:"pg_temp"`
	PgUpmap      []interface{} `json:"pg_upmap"`
	PgUpmapItems []interface{} `json:"pg_upmap_items"`
	PoolMax      int64         `json:"pool_max"`
	Pools        []struct {
		ApplicationMetadata struct {
			Rbd struct{} `json:"rbd"`
			Rgw struct{} `json:"rgw"`
		} `json:"application_metadata"`
		Auid                           int64         `json:"auid"`
		CacheMinEvictAge               int64         `json:"cache_min_evict_age"`
		CacheMinFlushAge               int64         `json:"cache_min_flush_age"`
		CacheMode                      string        `json:"cache_mode"`
		CacheTargetDirtyHighRatioMicro int64         `json:"cache_target_dirty_high_ratio_micro"`
		CacheTargetDirtyRatioMicro     int64         `json:"cache_target_dirty_ratio_micro"`
		CacheTargetFullRatioMicro      int64         `json:"cache_target_full_ratio_micro"`
		CreateTime                     string        `json:"create_time"`
		CrushRule                      int64         `json:"crush_rule"`
		ErasureCodeProfile             string        `json:"erasure_code_profile"`
		ExpectedNumObjects             int64         `json:"expected_num_objects"`
		FastRead                       bool          `json:"fast_read"`
		Flags                          int64         `json:"flags"`
		FlagsNames                     string        `json:"flags_names"`
		GradeTable                     []interface{} `json:"grade_table"`
		HitSetCount                    int64         `json:"hit_set_count"`
		HitSetGradeDecayRate           int64         `json:"hit_set_grade_decay_rate"`
		HitSetParams                   struct {
			Type string `json:"type"`
		} `json:"hit_set_params"`
		HitSetPeriod                 int64         `json:"hit_set_period"`
		HitSetSearchLastN            int64         `json:"hit_set_search_last_n"`
		LastChange                   string        `json:"last_change"`
		LastForceOpResend            string        `json:"last_force_op_resend"`
		LastForceOpResendPreluminous string        `json:"last_force_op_resend_preluminous"`
		MinReadRecencyForPromote     int64         `json:"min_read_recency_for_promote"`
		MinSize                      int64         `json:"min_size"`
		MinWriteRecencyForPromote    int64         `json:"min_write_recency_for_promote"`
		ObjectHash                   int64         `json:"object_hash"`
		Options                      struct{}      `json:"options"`
		PgNum                        int64         `json:"pg_num"`
		PgPlacementNum               int64         `json:"pg_placement_num"`
		Pool                         int64         `json:"pool"`
		PoolName                     string        `json:"pool_name"`
		PoolSnaps                    []interface{} `json:"pool_snaps"`
		QuotaMaxBytes                int64         `json:"quota_max_bytes"`
		QuotaMaxObjects              int64         `json:"quota_max_objects"`
		ReadTier                     int64         `json:"read_tier"`
		RemovedSnaps                 string        `json:"removed_snaps"`
		Size                         int64         `json:"size"`
		SnapEpoch                    int64         `json:"snap_epoch"`
		SnapMode                     string        `json:"snap_mode"`
		SnapSeq                      int64         `json:"snap_seq"`
		StripeWidth                  int64         `json:"stripe_width"`
		TargetMaxBytes               int64         `json:"target_max_bytes"`
		TargetMaxObjects             int64         `json:"target_max_objects"`
		TierOf                       int64         `json:"tier_of"`
		Tiers                        []interface{} `json:"tiers"`
		Type                         int64         `json:"type"`
		UseGmtHitset                 bool          `json:"use_gmt_hitset"`
		WriteTier                    int64         `json:"write_tier"`
	} `json:"pools"`
	PrimaryTemp            []interface{} `json:"primary_temp"`
	RemovedSnapsQueue      []interface{} `json:"removed_snaps_queue"`
	RequireMinCompatClient string        `json:"require_min_compat_client"`
	RequireOsdRelease      string        `json:"require_osd_release"`
}

type PlacementGroup struct {
	Acting                  []int64       `json:"acting"`
	ActingPrimary           int64         `json:"acting_primary"`
	BlockedBy               []interface{} `json:"blocked_by"`
	Created                 int64         `json:"created"`
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
	LastEpochClean          int64  `json:"last_epoch_clean"`
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
	LogSize                 int64         `json:"log_size"`
	LogStart                string        `json:"log_start"`
	ManifestStatsInvalid    bool          `json:"manifest_stats_invalid"`
	MappingEpoch            int64         `json:"mapping_epoch"`
	OmapStatsInvalid        bool          `json:"omap_stats_invalid"`
	OndiskLogSize           int64         `json:"ondisk_log_size"`
	OndiskLogStart          string        `json:"ondisk_log_start"`
	Parent                  string        `json:"parent"`
	ParentSplitBits         int64         `json:"parent_split_bits"`
	Pgid                    string        `json:"pgid"`
	PinStatsInvalid         bool          `json:"pin_stats_invalid"`
	PurgedSnaps             []interface{} `json:"purged_snaps"`
	ReportedEpoch           string        `json:"reported_epoch"`
	ReportedSeq             string        `json:"reported_seq"`
	SnaptrimqLen            int64         `json:"snaptrimq_len"`
	StatSum                 struct {
		NumBytes                   int64 `json:"num_bytes"`
		NumBytesHitSetArchive      int64 `json:"num_bytes_hit_set_archive"`
		NumBytesRecovered          int64 `json:"num_bytes_recovered"`
		NumDeepScrubErrors         int64 `json:"num_deep_scrub_errors"`
		NumEvict                   int64 `json:"num_evict"`
		NumEvictKb                 int64 `json:"num_evict_kb"`
		NumEvictModeFull           int64 `json:"num_evict_mode_full"`
		NumEvictModeSome           int64 `json:"num_evict_mode_some"`
		NumFlush                   int64 `json:"num_flush"`
		NumFlushKb                 int64 `json:"num_flush_kb"`
		NumFlushModeHigh           int64 `json:"num_flush_mode_high"`
		NumFlushModeLow            int64 `json:"num_flush_mode_low"`
		NumKeysRecovered           int64 `json:"num_keys_recovered"`
		NumLargeOmapObjects        int64 `json:"num_large_omap_objects"`
		NumLegacySnapsets          int64 `json:"num_legacy_snapsets"`
		NumObjectClones            int64 `json:"num_object_clones"`
		NumObjectCopies            int64 `json:"num_object_copies"`
		NumObjects                 int64 `json:"num_objects"`
		NumObjectsDegraded         int64 `json:"num_objects_degraded"`
		NumObjectsDirty            int64 `json:"num_objects_dirty"`
		NumObjectsHitSetArchive    int64 `json:"num_objects_hit_set_archive"`
		NumObjectsManifest         int64 `json:"num_objects_manifest"`
		NumObjectsMisplaced        int64 `json:"num_objects_misplaced"`
		NumObjectsMissing          int64 `json:"num_objects_missing"`
		NumObjectsMissingOnPrimary int64 `json:"num_objects_missing_on_primary"`
		NumObjectsOmap             int64 `json:"num_objects_omap"`
		NumObjectsPinned           int64 `json:"num_objects_pinned"`
		NumObjectsRecovered        int64 `json:"num_objects_recovered"`
		NumObjectsUnfound          int64 `json:"num_objects_unfound"`
		NumPromote                 int64 `json:"num_promote"`
		NumRead                    int64 `json:"num_read"`
		NumReadKb                  int64 `json:"num_read_kb"`
		NumScrubErrors             int64 `json:"num_scrub_errors"`
		NumShallowScrubErrors      int64 `json:"num_shallow_scrub_errors"`
		NumWhiteouts               int64 `json:"num_whiteouts"`
		NumWrite                   int64 `json:"num_write"`
		NumWriteKb                 int64 `json:"num_write_kb"`
	} `json:"stat_sum"`
	State        string  `json:"state"`
	StatsInvalid bool    `json:"stats_invalid"`
	Up           []int64 `json:"up"`
	UpPrimary    int64   `json:"up_primary"`
	Version      string  `json:"version"`
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
	Acting        []int64 `json:"acting"`
	ActingPrimary int64   `json:"acting_primary"`
	Epoch         int64   `json:"epoch"`
	Objname       string  `json:"objname"`
	Pgid          string  `json:"pgid"`
	Pool          string  `json:"pool"`
	PoolID        int64   `json:"pool_id"`
	RawPgid       string  `json:"raw_pgid"`
	Up            []int64 `json:"up"`
	UpPrimary     int64   `json:"up_primary"`
}

//todo check types (int64 -> uint64)
