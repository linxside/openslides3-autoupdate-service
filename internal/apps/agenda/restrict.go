package agenda

import (
	"encoding/json"
	"fmt"

	"github.com/OpenSlides/openslides3-autoupdate-service/internal/restricter"
)

const (
	pCanSee         = "agenda.can_see"
	pCanManage      = "agenda.can_manage"
	pCanSeeInternal = "agenda.can_see_internal_items"

	// CanSeeListOfSpeakers is the permission string if a user can see the list
	// of speakers.
	CanSeeListOfSpeakers = "agenda.can_see_list_of_speakers"
)

// Restrict handels restrictions of agenda/item elements.
func Restrict(r restricter.HasPermer) restricter.ElementFunc {
	return func(uid int, element json.RawMessage) (json.RawMessage, error) {
		if !r.HasPerm(uid, pCanSee) {
			return nil, nil
		}

		var agenda struct {
			IsHidden   bool `json:"is_hidden"`
			IsInternal bool `json:"is_internal"`
		}
		if err := json.Unmarshal(element, &agenda); err != nil {
			return nil, fmt.Errorf("decoding item: %w", err)
		}

		canManage := r.HasPerm(uid, pCanManage)
		canSeeInternal := r.HasPerm(uid, pCanSeeInternal)

		if !canManage && agenda.IsHidden {
			return nil, nil
		}

		if !canSeeInternal && agenda.IsInternal {
			return nil, nil
		}

		if canManage && canSeeInternal {
			return element, nil
		}

		var agendaData map[string]json.RawMessage
		if err := json.Unmarshal(element, &agendaData); err != nil {
			return nil, fmt.Errorf("decoding itemdata: %w", err)
		}

		if !canSeeInternal {
			delete(agendaData, "duration")
		}

		if !canManage {
			delete(agendaData, "comment")
		}

		element, err := json.Marshal(agendaData)
		if err != nil {
			return nil, fmt.Errorf("encoding itemdata: %w", err)
		}
		return element, nil
	}
}
