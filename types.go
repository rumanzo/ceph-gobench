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

type Moncommand struct {
	Prefix  string `json:"prefix"`
	Pool    string `json:"pool"`
	Format  string `json:"format"`
	Var     string `json:"var,omitempty"`
	Poolstr string `json:"poolstr,omitempty"`
}

type Monanswer struct {
	Pool   string `json:"pool,omitempty"`
	PoolId int    `json:"pool_id,omitempty"`
	Size   int    `json:"size,omitempty"`
}

func (times *PlacementGroups) StringsToTimes() {
	const LongForm = "2006-01-02 15:04:05.000000"
	times.Last_fresh, _ = time.Parse(LongForm, times.Last_fresh_str)
	times.Last_change, _ = time.Parse(LongForm, times.Last_fresh_str)
	times.Last_active, _ = time.Parse(LongForm, times.Last_active_str)
	times.Last_peered, _ = time.Parse(LongForm, times.Last_peered_str)
	times.Last_clean, _ = time.Parse(LongForm, times.Last_clean_str)
	times.Last_became_active, _ = time.Parse(LongForm, times.Last_became_active_str)
	times.Last_became_peered, _ = time.Parse(LongForm, times.Last_became_peered_str)
	times.Last_unstale, _ = time.Parse(LongForm, times.Last_unstale_str)
	times.Last_undegraded, _ = time.Parse(LongForm, times.Last_undegraded_str)
	times.Last_fullsized, _ = time.Parse(LongForm, times.Last_fullsized_str)
	times.Last_deep_scrub_stamp, _ = time.Parse(LongForm, times.Last_deep_scrub_stamp_str)
	times.Last_deep_scrub, _ = time.Parse(LongForm, times.Last_deep_scrub_str)
	times.Last_clean_scrub_stamp, _ = time.Parse(LongForm, times.Last_clean_scrub_stamp_str)
	times.Last_scrub_stamp, _ = time.Parse(LongForm, times.Last_scrub_stamp_str)
	times.Last_scrub, _ = time.Parse(LongForm, times.Last_scrub_str)
}

type PlacementGroups struct {
	Pgid                       string `json:"pgid"`
	Version                    string `json:"version"`
	Reported_seq               string `json:"reported_seq"`
	Reported_epoch             string `json:"reported_epoch"`
	State                      string `json:"state"`
	Last_fresh_str             string `json:"last_fresh"`
	Last_fresh                 time.Time
	Last_change_str            string `json:"last_change"`
	Last_change                time.Time
	Last_active_str            string `json:"last_active"`
	Last_active                time.Time
	Last_peered_str            string `json:"last_peered"`
	Last_peered                time.Time
	Last_clean_str             string `json:"last_clean"`
	Last_clean                 time.Time
	Last_became_active_str     string `json:"last_became_active"`
	Last_became_active         time.Time
	Last_became_peered_str     string `json:"last_became_peered"`
	Last_became_peered         time.Time
	Last_unstale_str           string `json:"last_unstale"`
	Last_unstale               time.Time
	Last_undegraded_str        string `json:"last_undegraded"`
	Last_undegraded            time.Time
	Last_fullsized_str         string `json:"last_fullsized"`
	Last_fullsized             time.Time
	Mapping_epoch              float64 `json:"mapping_epoch"`
	Log_start                  string  `json:"log_start"`
	Ondisk_log_start           string  `json:"ondisk_log_start"`
	Created                    float64 `json:"created"`
	Last_epoch_clean           float64 `json:"last_epoch_clean"`
	Parent                     string  `json:"parent"`
	Parent_split_bits          float64 `json:"parent_split_bits"`
	Last_scrub_str             string  `json:"last_scrub"`
	Last_scrub                 time.Time
	Last_scrub_stamp_str       string `json:"last_scrub_stamp"`
	Last_scrub_stamp           time.Time
	Last_deep_scrub_str        string `json:"last_deep_scrub"`
	Last_deep_scrub            time.Time
	Last_deep_scrub_stamp_str  string `json:"last_deep_scrub_stamp"`
	Last_deep_scrub_stamp      time.Time
	Last_clean_scrub_stamp_str string `json:"last_clean_scrub_stamp"`
	Last_clean_scrub_stamp     time.Time
	Log_size                   float64   `json:"log_size"`
	Ondisk_log_size            float64   `json:"ondisk_log_size"`
	Stats_invalid              bool      `json:"stats_invalid"`
	Dirty_stats_invalid        bool      `json:"dirty_stats_invalid"`
	Omap_stats_invalid         bool      `json:"omap_stats_invalid"`
	Hitset_stats_invalid       bool      `json:"hitset_stats_invalid"`
	Hitset_bytes_stats_invalid bool      `json:"hitset_bytes_stats_invalid"`
	Pin_stats_invalid          bool      `json:"pin_stats_invalid"`
	Manifest_stats_invalid     bool      `json:"manifest_stats_invalid"`
	Snaptrimq_len              float64   `json:"snaptrimq_len"`
	Stat_sum                   StatSum   `json:"stat_sum"`
	Up                         []float64 `json:"up"`
	Acting                     []float64 `json:"acting"`
	Blocked_by                 []float64 `json:"blocked_by"`
	Up_primary                 float64   `json:"up_primary"`
	Acting_primary             float64   `json:"acting_primary"`
	Purged_snaps               []string  `json:"purged_snaps"`
}

type StatSum struct {
	Num_objects                    float64 `json:"num_objects"`
	Num_object_clones              float64 `json:"num_object_clones"`
	Num_object_copies              float64 `json:"num_object_copies"`
	Num_objects_missing_on_primary float64 `json:"num_objects_missing_on_primary"`
	Num_objects_missing            float64 `json:"num_objects_missing"`
	Num_objects_degraded           float64 `json:"num_objects_degraded"`
	Num_objects_misplaced          float64 `json:"num_objects_misplaced"`
	Num_objects_unfound            float64 `json:"num_objects_unfound"`
	Num_objects_dirty              float64 `json:"num_objects_dirty"`
	Num_whiteouts                  float64 `json:"num_whiteouts"`
	Num_read                       float64 `json:"num_read"`
	Num_read_kb                    float64 `json:"num_read_kb"`
	Num_write                      float64 `json:"num_write"`
	Num_write_kb                   float64 `json:"num_write_kb"`
	Num_scrub_errors               float64 `json:"num_scrub_errors"`
	Num_shallow_scrub_errors       float64 `json:"num_shallow_scrub_errors"`
	Num_deep_scrub_errors          float64 `json:"num_deep_scrub_errors"`
	Num_objects_recovered          float64 `json:"num_objects_recovered"`
	Num_bytes_recovered            float64 `json:"num_bytes_recovered"`
	Num_keys_recovered             float64 `json:"num_keys_recovered"`
	Num_objects_omap               float64 `json:"num_objects_omap"`
	Num_objects_hit_set_archive    float64 `json:"num_objects_hit_set_archive"`
	Num_bytes_hit_set_archive      float64 `json:"num_bytes_hit_set_archive"`
	Num_flush                      float64 `json:"num_flush"`
	Num_flush_kb                   float64 `json:"num_flush_kb"`
	Num_evict                      float64 `json:"num_evict"`
	Num_evict_kb                   float64 `json:"num_evict_kb"`
	Num_promote                    float64 `json:"num_promote"`
	Num_flush_mode_high            float64 `json:"num_flush_mode_high"`
	Num_flush_mode_low             float64 `json:"num_flush_mode_low"`
	Num_evict_mode_some            float64 `json:"num_evict_mode_some"`
	Num_evict_mode_full            float64 `json:"num_evict_mode_full"`
	Num_objects_pinned             float64 `json:"num_objects_pinned"`
	Num_legacy_snapsets            float64 `json:"num_legacy_snapsets"`
	Num_large_omap_objects         float64 `json:"num_large_omap_objects"`
	Num_objects_manifest           float64 `json:"num_objects_manifest"`
}
