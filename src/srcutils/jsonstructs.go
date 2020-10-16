package srcutils

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type scriptExportJson struct {
	Slots    map[string]*Slot `json:"slots"` // keys are quoted numbers
	Handlers []*handlerJson   `json:"handlers"`
	Methods  []*Method        `json:"methods"`
	Events   []*Event         `json:"events"`
}

func (e *ScriptExport) UnmarshalJSON(d []byte) error {
	tmp := &scriptExportJson{}
	err := json.Unmarshal(d, tmp)
	if err != nil {
		return errors.WithStack(err)
	}

	slots := make(map[int]*Slot, len(tmp.Slots))
	for k, v := range tmp.Slots {
		v := v // we're referencing this so need to declare inside the loop
		kint, err := strconv.Atoi(k)
		if err != nil {
			return errors.WithStack(err)
		}

		slots[kint] = v
	}

	handlers := make([]*Handler, len(tmp.Handlers))
	for k, v := range tmp.Handlers {
		slotKey, _ := v.Filter.SlotKey.Int64()
		key, _ := v.Key.Int64()

		handlers[k] = &Handler{
			Code: v.Code,
			Filter: &Filter{
				Args:      v.Filter.Args,
				Signature: v.Filter.Signature,
				SlotKey:   int(slotKey),
			},
			Key: int(key),
		}
	}

	e.Slots = slots
	e.Handlers = handlers
	e.Methods = tmp.Methods
	e.Events = tmp.Events

	return nil
}

func (e *ScriptExport) MarshalJSON() ([]byte, error) {
	slots := make(map[string]*Slot, len(e.Slots))
	for k, v := range e.Slots {
		kstr := strconv.Itoa(k)
		slots[kstr] = v
	}

	handlers := make([]*handlerJson, len(e.Handlers))
	for k, v := range e.Handlers {
		handlers[k] = &handlerJson{
			Code: v.Code,
			Filter: &filterJson{
				Args:      v.Filter.Args,
				Signature: v.Filter.Signature,
				SlotKey:   json.Number(fmt.Sprintf("%d", v.Filter.SlotKey)),
			},
			Key: json.Number(fmt.Sprintf("%d", v.Key)),
		}
	}

	tmp := &scriptExportJson{
		Slots:    slots,
		Handlers: handlers,
		Methods:  e.Methods,
		Events:   e.Events,
	}
	res, err := json.Marshal(tmp)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

type handlerJson struct {
	Code   string      `json:"code"`
	Filter *filterJson `json:"filter"`
	Key    json.Number `json:"key"` // can be quoted and unquoted
}

type filterJson struct {
	Args      []Arg       `json:"args"`
	Signature string      `json:"signature"`
	SlotKey   json.Number `json:"slotKey"` // can be quoted and unquoted
}