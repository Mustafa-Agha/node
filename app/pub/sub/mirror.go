package sub

import (
	"github.com/cosmos/cosmos-sdk/pubsub"

	"github.com/Mustafa-Agha/node/plugins/bridge"
)

func SubscribeMirrorEvent(sub *pubsub.Subscriber) error {
	err := sub.Subscribe(bridge.MirrorTopic, func(event pubsub.Event) {
		switch event := event.(type) {
		case bridge.MirrorEvent:
			mirrorEvent := event
			if stagingArea.MirrorData == nil {
				stagingArea.MirrorData = make([]bridge.MirrorEvent, 0, 1)
			}
			stagingArea.MirrorData = append(stagingArea.MirrorData, mirrorEvent)
		default:
			sub.Logger.Info("unknown event type")
		}
	})
	return err
}

func commitMirror() {
	if len(stagingArea.MirrorData) > 0 {
		toPublish.EventData.MirrorData = append(toPublish.EventData.MirrorData, stagingArea.MirrorData...)
	}
}
