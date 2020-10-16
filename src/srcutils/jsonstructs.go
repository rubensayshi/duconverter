package srcutils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type scriptExportJson struct {
	Slots    map[string]*Slot `json:"slots"` // keys are quoted numbers
	Handlers []*handlerJson   `json:"handlers"`
	Methods  []*Method        `json:"methods"`
	Events   []*Event         `json:"events"`
}

func (e *ScriptExport) UnmarshalJSON(d []byte) error {
	// unmarshall into a tmp which resembles the actual json structure
	tmp := &scriptExportJson{}
	err := json.Unmarshal(d, tmp)
	if err != nil {
		return errors.WithStack(err)
	}

	// indexes are quoted strings, we need to convert them to ints
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
		// slotKey can be quoted string (or not)...
		slotKeyRaw := strings.Trim(string(v.Filter.SlotKey), `"`)
		slotKey, err := strconv.ParseInt(slotKeyRaw, 10, 64)
		if err != nil {
			return errors.WithStack(err)
		}

		// key can be quoted string (or not)...
		keyRaw := strings.Trim(string(v.Key), `"`)
		key, err := strconv.ParseInt(keyRaw, 10, 64)
		if err != nil {
			return errors.WithStack(err)
		}

		args := make([]string, len(v.Filter.Args))
		for k, v := range v.Filter.Args {
			args[k] = v.Value
		}

		handlers[k] = &Handler{
			Code: v.Code,
			Filter: &Filter{
				Args:      args,
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
		args := make([]argJson, len(v.Filter.Args))
		for k, v := range v.Filter.Args {
			args[k] = argJson{Value: v}
		}

		handlers[k] = &handlerJson{
			Code: v.Code,
			Filter: &filterJson{
				Args:      args,
				Signature: v.Filter.Signature,
				SlotKey:   []byte(fmt.Sprintf("\"%d\"", v.Filter.SlotKey)),
			},
			Key: []byte(fmt.Sprintf("\"%d\"", v.Key)),
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
	Code   string          `json:"code"`
	Filter *filterJson     `json:"filter"`
	Key    json.RawMessage `json:"key"` // can be quoted and unquoted number
}

type filterJson struct {
	Args      []argJson       `json:"args"`
	Signature string          `json:"signature"`
	SlotKey   json.RawMessage `json:"slotKey"` // can be quoted and unquoted number
}

type argJson struct {
	Value string `json:"value"`
}
