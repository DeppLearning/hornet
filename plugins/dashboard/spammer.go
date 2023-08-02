package dashboard

import (
	"github.com/iotaledger/hive.go/daemon"
	"github.com/iotaledger/hive.go/events"

	"github.com/iotaledger/hornet/pkg/shutdown"
	"github.com/iotaledger/hornet/pkg/spammer"
	spammerplugin "github.com/iotaledger/hornet/plugins/spammer"
)

func runSpammerMetricWorker() {

	onSpamPerformed := events.NewClosure(func(metrics *spammer.SpamStats) {
		hub.BroadcastMsg(&Msg{Type: MsgTypeSpamMetrics, Data: metrics})
	})

	onAvgSpamMetricsUpdated := events.NewClosure(func(metrics *spammer.AvgSpamMetrics) {
		hub.BroadcastMsg(&Msg{Type: MsgTypeAvgSpamMetrics, Data: metrics})
	})

	daemon.BackgroundWorker("Dashboard[SpammerMetricUpdater]", func(shutdownSignal <-chan struct{}) {
		spammerplugin.Events.SpamPerformed.Attach(onSpamPerformed)
		spammerplugin.Events.AvgSpamMetricsUpdated.Attach(onAvgSpamMetricsUpdated)
		<-shutdownSignal
		log.Info("Stopping Dashboard[SpammerMetricUpdater] ...")
		spammerplugin.Events.SpamPerformed.Detach(onSpamPerformed)
		spammerplugin.Events.AvgSpamMetricsUpdated.Detach(onAvgSpamMetricsUpdated)
		log.Info("Stopping Dashboard[SpammerMetricUpdater] ... done")
	}, shutdown.PriorityDashboard)
}
