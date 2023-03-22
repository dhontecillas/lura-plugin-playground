package main

type CaseModifier struct {
	UseUpper bool
}

func NewCaseModifierConfig(extra map[string]interface{}) (*CaseModifier, error) {
	cfgRaw, ok := extra[pluginName].(map[string]interface{})
	if !ok {
		return nil, errRegistererNotFound
	}

	toUpper := false
	if v, ok := cfgRaw["useUpper"]; ok {
		if strV, ok := v.(string); ok {
			toUpper = (strV == "true")
		}
	}
	return &CaseModifier{
		UseUpper: toUpper,
	}, nil
}
