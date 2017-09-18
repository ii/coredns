package file

import (
	"log"
	"os"
	"time"
)

// TickTime is the default time we use to reload zone. Exported to be tweaked in tests.
var TickTime = 1 * time.Minute

// Reload reloads a zone when it is changed on disk. If z.NoRoload is true, no reloading will be done.
func (z *Zone) Reload() error {
	if z.NoReload {
		return nil
	}

	tick := time.NewTicker(TickTime)

	go func() {

		for {
			select {

			case <-tick.C:
				reader, err := os.Open(z.file)
				if err != nil {
					log.Printf("[ERROR] Failed to open `%s' for `%s': %v", z.file, z.origin, err)
					continue
				}

				serial := z.SOASerialIfDefined()
				zone, err := Parse(reader, z.origin, z.file, serial)
				if err != nil {
					log.Printf("[WARNING] Parsing zone `%s': %v", z.origin, err)
					continue
				}

				// copy elements we need
				z.reloadMu.Lock()
				z.Apex = zone.Apex
				z.Tree = zone.Tree
				z.reloadMu.Unlock()

				log.Printf("[INFO] Successfully reloaded zone `%s'", z.origin)
				z.Notify()

			case <-z.ReloadShutdown:
				tick.Stop()
				return
			}
		}
	}()
	return nil
}

// SOASerialIfDefined returns the SOA's serial if the zone has a SOA record in the Apex, or
// -1 otherwise.
func (z *Zone) SOASerialIfDefined() int64 {
	z.reloadMu.Lock()
	defer z.reloadMu.Unlock()
	if z.Apex.SOA != nil {
		return int64(z.Apex.SOA.Serial)
	}
	return -1
}
