package prompts

import (
	"github.com/FlagrantGarden/flfa/pkg/flfa/state/skirmish"
	"github.com/FlagrantGarden/flfa/pkg/tympan/state/instance"
	"github.com/erikgeiser/promptkit/selection"
)

func WhatNext(skirmish *instance.Instance[skirmish.Skirmish]) *selection.Selection {
	options := []string{"Check on Something", "Attack", "Move", "Cast", "Shoot", "Save & Quit"}
	return selection.New("What do you want to do now?", selection.Choices(options))
}
