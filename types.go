package main

import (
	"github.com/ceph/go-ceph/rados"
	"time"
)

type Params struct {
	duration                                               time.Duration
	threadsCount                                           int
	blocksize                                              uint64
	parallel                                               bool
	bs, cluster, user, keyring, config, pool, mode, define string
}

type Cephconnection struct {
	conn  *rados.Conn
	ioctx *rados.IOContext
}

type Poolinfo struct {
	Pool   string `json:"pool,omitempty"`
	PoolId int    `json:"pool_id,omitempty"`
	Size   int    `json:"size,omitempty"`
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

type OsdCrushDump struct {
	Buckets []struct {
		Alg   string `json:"alg"`
		Hash  string `json:"hash"`
		ID    int64  `json:"id"`
		Items []struct {
			ID     int64 `json:"id"`
			Pos    int64 `json:"pos"`
			Weight int64 `json:"weight"`
		} `json:"items"`
		Name     string `json:"name"`
		TypeID   int64  `json:"type_id"`
		TypeName string `json:"type_name"`
		Weight   int64  `json:"weight"`
	} `json:"buckets"`
	ChooseArgs struct{} `json:"choose_args"`
	Devices    []struct {
		Class string `json:"class"`
		ID    int64  `json:"id"`
		Name  string `json:"name"`
	} `json:"devices"`
	Rules []struct {
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
	Buckets []struct {
		Alg   string `json:"alg"`
		Hash  string `json:"hash"`
		ID    int64  `json:"id"`
		Items []struct {
			ID     int64 `json:"id"`
			Pos    int64 `json:"pos"`
			Weight int64 `json:"weight"`
		} `json:"items"`
		Name     string `json:"name"`
		TypeID   int64  `json:"type_id"`
		TypeName string `json:"type_name"`
		Weight   int64  `json:"weight"`
	} `json:"buckets"`
	ChooseArgs struct{} `json:"choose_args"`
	Devices    []struct {
		Class string `json:"class"`
		ID    int64  `json:"id"`
		Name  string `json:"name"`
	} `json:"devices"`
	Rules []struct {
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
