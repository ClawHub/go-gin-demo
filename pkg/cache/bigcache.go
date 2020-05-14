package gcache

import (
	"github.com/allegro/bigcache"
	"log"
	"time"
)

var Cache *bigcache.BigCache

func Setup() {
	config := bigcache.Config{
		// 分片的数量 (must be a power of 2)
		Shards: 1024,

		// 时间之后，entry可以被驱逐，在此之后，条目可以被称为已死，但不能被删除。
		LifeWindow: 10 * time.Minute,

		// 移除过期entry的间隔时间 (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
		CleanWindow: 5 * time.Minute,

		// rps * lifeWindow, 仅在初始内存分配中使用
		MaxEntriesInWindow: 1000 * 10 * 60,

		// 最大条目大小(以字节为单位)，仅在初始内存分配中使用
		MaxEntrySize: 500,

		// 打印关于额外内存分配的信息
		Verbose: true,

		// 缓存不会分配超过这个限制的内存, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: 8192,

		// 当最老的条目由于过期时间或没有剩余空间而被删除时触发的回调
		//对于新条目，或者因为调用了delete。返回一个表示原因的位掩码。
		//默认值为nil，这意味着没有回调，它阻止打开最旧的条目。
		OnRemove: nil,

		// OnRemoveWithReason是一个回调函数，当最老的条目由于过期时间或没有空间留给新条目或因为调用了delete而被删除时触发。
		//一个代表理性的常数将被传递。默认值为nil，这意味着没有回调，它阻止打开最旧的条目。如果指定OnRemove，则忽略。
		OnRemoveWithReason: nil,
	}
	var initErr error
	Cache, initErr = bigcache.NewBigCache(config)
	if initErr != nil {
		log.Fatal(initErr)
	}
}
